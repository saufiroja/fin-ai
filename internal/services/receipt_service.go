package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"math"
	"mime/multipart"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/oklog/ulid/v2"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/saufiroja/fin-ai/internal/constants/prompt"
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/domains/categories"
	"github.com/saufiroja/fin-ai/internal/domains/log_message"
	"github.com/saufiroja/fin-ai/internal/domains/receipt"
	"github.com/saufiroja/fin-ai/internal/domains/transaction"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/pkg/llm"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
	"github.com/saufiroja/fin-ai/pkg/minio"
	"google.golang.org/genai"
)

type receiptService struct {
	receiptRepository  receipt.ReceiptStorer
	transactionService transaction.TransactionManager
	logMessageService  log_message.LogMessageManager
	categoryService    categories.CategoryManager
	minioClient        minio.MinioManager
	logging            logging.Logger
	openaiClient       llm.OpenAI
	geminiClient       llm.Gemini
	bucketName         string
	objectName         string
}

func NewReceiptService(
	receiptRepository receipt.ReceiptStorer,
	transactionService transaction.TransactionManager,
	logMessageService log_message.LogMessageManager,
	categoryService categories.CategoryManager,
	minioClient minio.MinioManager,
	logging logging.Logger,
	openaiClient llm.OpenAI,
	geminiClient llm.Gemini,
) receipt.ReceiptManager {
	return &receiptService{
		receiptRepository:  receiptRepository,
		transactionService: transactionService,
		logMessageService:  logMessageService,
		categoryService:    categoryService,
		minioClient:        minioClient,
		logging:            logging,
		openaiClient:       openaiClient,
		bucketName:         "receipts",
		objectName:         "receipt",
		geminiClient:       geminiClient,
	}
}

func (s *receiptService) UploadReceipt(filePath *multipart.FileHeader, userId string) (*models.Receipt, error) {
	s.logging.LogInfo(fmt.Sprintf("Uploading receipt for user %s from file %s", userId, filePath.Filename))

	if err := s.validateFileSize(filePath); err != nil {
		return nil, err
	}

	optimizedImageBytes, err := s.processImage(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to process image: %w", err)
	}

	categoriesOfString, err := s.getCategoriesString()
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	if err := s.uploadToMinIO(filePath, userId); err != nil {
		return nil, fmt.Errorf("failed to upload to MinIO: %w", err)
	}

	extractedData, responseString, responseAi, err := s.processReceiptWithGemini(optimizedImageBytes, categoriesOfString, filePath.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to process receipt with AI: %w", err)
	}

	if err := s.logAIResponse(responseString, responseAi, userId); err != nil {
		return nil, fmt.Errorf("failed to log AI response: %w", err)
	}

	receipt, err := s.saveReceipt(extractedData, filePath, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to save receipt to database: %w", err)
	}

	s.logging.LogInfo(fmt.Sprintf("Receipt for user %s uploaded and processed successfully", userId))
	return receipt, nil
}

func (s *receiptService) validateFileSize(filePath *multipart.FileHeader) error {
	if filePath.Size > 10*1024*1024 {
		s.logging.LogError(fmt.Sprintf("File %s exceeds the maximum size limit of 10MB", filePath.Filename))
		return fmt.Errorf("file exceeds the maximum size limit of 10MB")
	}
	return nil
}

func (s *receiptService) processImage(filePath *multipart.FileHeader) ([]byte, error) {
	file, err := filePath.Open()
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to open multipart file %s: %v", filePath.Filename, err))
		return nil, fmt.Errorf("failed to open multipart file: %w", err)
	}
	defer file.Close()

	// Decode the image from the multipart file
	img, err := imaging.Decode(file)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to decode image file %s: %v", filePath.Filename, err))
		return nil, fmt.Errorf("failed to decode image file: %w", err)
	}

	img = s.optimizeImageForAIVision(img)

	optimizedImageBytes, err := s.imageToBytes(img, filePath.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to convert image to bytes: %w", err)
	}

	return optimizedImageBytes, nil
}

func (s *receiptService) optimizeImageForAIVision(img image.Image) image.Image {
	// 1. Auto-rotate if needed
	img = s.autoRotateImage(img)

	// 2. Resize to optimal size for AI Vision
	img = imaging.Resize(img, 1600, 0, imaging.Lanczos)

	// 3. Enhanced gamma correction for AI Vision
	img = imaging.AdjustGamma(img, 1.1)

	// 4. Enhanced contrast for better text clarity
	img = imaging.AdjustContrast(img, 40)

	// 5. Brightness adjustment for optimal reading
	img = imaging.AdjustBrightness(img, 15)

	// 6. Sharpen for text clarity
	img = imaging.Sharpen(img, 1.2)

	return img
}

func (s *receiptService) imageToBytes(img image.Image, filename string) ([]byte, error) {
	buf := new(bytes.Buffer)

	var format imaging.Format = imaging.JPEG
	if strings.HasSuffix(strings.ToLower(filename), ".png") {
		format = imaging.PNG
	}

	err := imaging.Encode(buf, img, format)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to encode optimized image: %v", err))
		return nil, fmt.Errorf("failed to encode optimized image: %w", err)
	}

	return buf.Bytes(), nil
}

func (s *receiptService) getCategoriesString() (string, error) {
	reqCategoryQuery := &requests.GetAllCategoryQuery{
		Limit:  100,
		Offset: 0,
	}

	categories, err := s.categoryService.FindAllCategories(reqCategoryQuery)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to get categories: %v", err))
		return "", fmt.Errorf("failed to get categories: %w", err)
	}

	categoriesStrings := make([]string, len(categories.Categories))
	for i, category := range categories.Categories {
		categoriesStrings[i] = fmt.Sprintf("%s (%s)", category.Name, category.CategoryId)
	}

	return strings.Join(categoriesStrings, ", "), nil
}

func (s *receiptService) uploadToMinIO(filePath *multipart.FileHeader, userId string) error {
	bucketExists, err := s.minioClient.BucketExists(s.bucketName)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Error checking bucket existence: %v", err))
		return fmt.Errorf("error checking bucket existence: %w", err)
	}

	if !bucketExists {
		err = s.minioClient.CreateBucket(s.bucketName)
		if err != nil {
			s.logging.LogError(fmt.Sprintf("Error creating bucket %s: %v", s.bucketName, err))
			return fmt.Errorf("error creating bucket: %w", err)
		}
		s.logging.LogInfo(fmt.Sprintf("Bucket %s created successfully", s.bucketName))
	}

	s.objectName = fmt.Sprintf("%s/%s", userId, filePath.Filename)

	exists, err := s.minioClient.FileExists(s.bucketName, s.objectName)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Error checking file existence: %v", err))
		return err
	}

	if !exists {
		err = s.minioClient.UploadFileFromMultipart(s.bucketName, s.objectName, filePath)
		if err != nil {
			s.logging.LogError(fmt.Sprintf("Failed to upload file %s to bucket %s: %v", filePath.Filename, s.bucketName, err))
			return fmt.Errorf("failed to upload file: %w", err)
		}
		s.logging.LogInfo(fmt.Sprintf("File %s uploaded successfully to bucket %s", filePath.Filename, s.bucketName))
	}

	return nil
}

func (s *receiptService) processReceiptWithGemini(optimizedImageBytes []byte, categoriesOfString, filename string) (*responses.ReceiptExtractionResponse, string, *responses.ResponseAI, error) {
	imageType := "image/png"
	if strings.HasSuffix(strings.ToLower(filename), ".jpg") || strings.HasSuffix(strings.ToLower(filename), ".jpeg") {
		imageType = "image/jpeg"
	}

	messagePrompt := fmt.Sprintf(prompt.ReceiptExtractionUserPromptTemplate, categoriesOfString)
	parts := []*genai.Part{
		genai.NewPartFromBytes(optimizedImageBytes, imageType),
		genai.NewPartFromText(messagePrompt),
	}

	messages := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	responseAi, err := s.geminiClient.Run(context.Background(), "gemini-2.5-flash", messages)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to process receipt with AI: %v", err))
		return nil, "", nil, fmt.Errorf("failed to process receipt with AI: %w", err)
	}

	responseString, ok := responseAi.Response.(string)
	if !ok {
		s.logging.LogError("Failed to convert AI response to string")
		return nil, "", nil, fmt.Errorf("failed to convert AI response to string")
	}

	// Clean the response string to extract JSON content
	cleanedResponse := s.cleanAIResponse(responseString)

	var extractedData responses.ReceiptExtractionResponse
	err = json.Unmarshal([]byte(cleanedResponse), &extractedData)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to parse JSON response: %v", err))
		s.logging.LogError(fmt.Sprintf("Raw response: %s", responseString))
		s.logging.LogError(fmt.Sprintf("Cleaned response: %s", cleanedResponse))
		return nil, "", nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &extractedData, responseString, responseAi, nil
}

func (s *receiptService) logAIResponse(responseString string, responseAi *responses.ResponseAI, userId string) error {
	messagePromptJSON, err := json.Marshal(responseString)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to marshal message prompt: %v", err))
		return fmt.Errorf("failed to marshal message prompt: %w", err)
	}

	dateNow := time.Now()
	logMessage := &models.LogMessage{
		LogMessageId: ulid.Make().String(),
		UserId:       userId,
		Message:      string(messagePromptJSON),
		Response:     responseString,
		InputToken:   responseAi.InputToken,
		OutputToken:  responseAi.OutputToken,
		Topic:        "receipt_extraction",
		Model:        "gemini-2.5-flash",
		CreatedAt:    dateNow,
		UpdatedAt:    dateNow,
	}

	err = s.logMessageService.InsertLogMessage(logMessage)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to insert log message: %v", err))
		return fmt.Errorf("failed to insert log message: %w", err)
	}

	return nil
}

func (s *receiptService) saveReceipt(extractedData *responses.ReceiptExtractionResponse, filePath *multipart.FileHeader, userId string) (*models.Receipt, error) {
	dateNow := time.Now()

	extractedReceiptJSON, err := json.Marshal(extractedData.ExtractedReceipt)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to marshal extracted receipt: %v", err))
		return nil, fmt.Errorf("failed to marshal extracted receipt: %w", err)
	}

	input := openai.EmbeddingNewParamsInputUnion{
		OfString: param.NewOpt(string(extractedReceiptJSON)),
	}
	embedding := s.openaiClient.CreateEmbedding(context.Background(), input)

	metaData := models.MetaData{
		FileName: filePath.Filename,
		FileSize: filePath.Size,
		FileType: filePath.Header.Get("Content-Type"),
	}

	metaDataJSON, err := json.Marshal(metaData)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to marshal metadata: %v", err))
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Create receipt model
	receiptModel := &models.Receipt{
		ReceiptId:                 ulid.Make().String(),
		UserId:                    userId,
		MerchantName:              extractedData.ExtractedReceipt.MerchantName,
		SubTotal:                  extractedData.ExtractedReceipt.SubTotal,
		TotalDiscount:             extractedData.ExtractedReceipt.TotalDiscount,
		TotalShopping:             extractedData.ExtractedReceipt.TotalShopping,
		MetaData:                  metaDataJSON,
		ExtractedReceipt:          extractedReceiptJSON,
		ExtractedReceiptEmbedding: embedding.Embeddings,
		Confirmed:                 false,
		TransactionDate:           dateNow,
		CreatedAt:                 dateNow,
		UpdatedAt:                 dateNow,
	}

	// Insert receipt
	err = s.insertReceipt(receiptModel)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to insert receipt: %v", err))
		return nil, fmt.Errorf("failed to insert receipt: %w", err)
	}

	// Insert receipt items
	err = s.insertReceiptItems(extractedData.ExtractedReceipt.Items, receiptModel.ReceiptId, userId, dateNow)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to insert receipt items: %v", err))
		return nil, fmt.Errorf("failed to insert receipt items: %w", err)
	}

	return receiptModel, nil
}

func (s *receiptService) insertReceiptItems(items []responses.ReceiptItemResponse, receiptId, userId string, dateNow time.Time) error {
	for _, item := range items {
		// Create receipt item
		receiptItem := &models.ReceiptItem{
			ReceiptItemId:        ulid.Make().String(),
			ReceiptId:            receiptId,
			ItemName:             item.ItemName,
			ItemQuantity:         item.ItemQuantity,
			ItemPrice:            item.ItemPrice,
			ItemPriceTotal:       item.ItemPriceTotal,
			ItemDiscount:         item.ItemDiscount,
			CreatedAt:            dateNow,
			UpdatedAt:            dateNow,
			CategoryId:           item.CategoryId,
			AiCategoryConfidence: item.AiCategoryConfidence,
		}

		err := s.receiptRepository.InsertReceiptItem(receiptItem)
		if err != nil {
			s.logging.LogError(fmt.Sprintf("Failed to insert receipt item: %v", err))
			return fmt.Errorf("failed to insert receipt item: %w", err)
		}

		s.logging.LogInfo(fmt.Sprintf("Receipt item %s inserted successfully for receipt %s", item.ItemName, receiptId))
	}

	return nil
}

func (s *receiptService) autoRotateImage(img image.Image) image.Image {
	// Simple orientation detection - can be enhanced with more sophisticated algorithms
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// If width > height significantly, might need rotation
	if float64(width)/float64(height) > 1.5 {
		img = imaging.Rotate90(img)
	}

	return img
}

func (s *receiptService) insertReceipt(receipt *models.Receipt) error {
	s.logging.LogInfo(fmt.Sprintf("Inserting receipt for user %s", receipt.UserId))
	err := s.receiptRepository.InsertReceipt(receipt)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to insert receipt: %v", err))
		return fmt.Errorf("failed to insert receipt: %w", err)
	}
	s.logging.LogInfo("Receipt inserted successfully")
	return nil
}

func (s *receiptService) GetReceiptsByUserId(userId string) ([]*models.Receipt, error) {
	s.logging.LogInfo(fmt.Sprintf("Fetching all receipts for user %s", userId))

	receipts, err := s.receiptRepository.GetReceiptsByUserId(userId)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to fetch receipts for user %s: %v", userId, err))
		return nil, fmt.Errorf("failed to fetch receipts: %w", err)
	}

	s.logging.LogInfo(fmt.Sprintf("Fetched %d receipts for user %s", len(receipts), userId))
	return receipts, nil
}

func (s *receiptService) GetAllReceiptsByUserId(userId string, req *requests.GetAllReceiptsQuery) (*responses.ReceiptResponse, error) {
	s.logging.LogInfo(fmt.Sprintf("Fetching all receipts for user %s", userId))

	switch {
	case req.SortBy == "":
		req.SortBy = "created_at"
	case req.SortBy == "merchant_name":
		req.SortBy = "merchant_name"
	case req.SortBy == "total_shopping":
		req.SortBy = "total_shopping"
	}

	offset := 0
	if req.Offset > 1 {
		offset = (req.Offset - 1) * req.Limit
	}

	queryReq := &requests.GetAllReceiptsQuery{
		Offset:    offset,
		Limit:     req.Limit,
		Search:    req.Search,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	receipts, err := s.receiptRepository.GetAllReceiptsByUserId(userId, queryReq)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to fetch receipts for user %s: %v", userId, err))
		return nil, fmt.Errorf("failed to fetch receipts: %w", err)
	}

	count, err := s.receiptRepository.CountReceiptsByUserId(userId, queryReq)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Error counting receipts: %v", err))
		return nil, err
	}

	totalPages := math.Ceil(float64(count) / float64(req.Limit))
	currentPage := math.Min(float64(req.Offset), float64(totalPages))

	res := &responses.ReceiptResponse{
		TotalPages:  int64(totalPages),
		CurrentPage: int64(currentPage),
		Total:       int64(count),
		Receipts:    receipts,
	}

	s.logging.LogInfo(fmt.Sprintf("Fetched %d receipts for user %s", len(receipts), userId))
	return res, nil
}

func (s *receiptService) GetDetailReceiptUserById(userId string, receiptId string) (*responses.DetailReceiptUserResponse, error) {
	s.logging.LogInfo(fmt.Sprintf("Fetching detail receipt for user %s and receipt ID %s", userId, receiptId))

	receipt, err := s.receiptRepository.GetDetailReceiptUserById(userId, receiptId)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to fetch detail receipt for user %s and receipt ID %s: %v", userId, receiptId, err))
		return nil, fmt.Errorf("failed to fetch detail receipt: %w", err)
	}

	items, err := s.receiptRepository.GetReceiptItemsByReceiptId(receipt.ReceiptId)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to fetch receipt items for receipt ID %s: %v", receipt.ReceiptId, err))
		return nil, fmt.Errorf("failed to fetch receipt items: %w", err)
	}

	detailResponse := &responses.DetailReceiptUserResponse{
		ReceiptId:       receipt.ReceiptId,
		UserId:          receipt.UserId,
		MerchantName:    receipt.MerchantName,
		SubTotal:        receipt.SubTotal,
		TotalDiscount:   receipt.TotalDiscount,
		TotalShopping:   receipt.TotalShopping,
		TransactionDate: receipt.TransactionDate,
		CreatedAt:       receipt.CreatedAt,
		UpdatedAt:       receipt.UpdatedAt,
		Items:           items,
		Confirmed:       receipt.Confirmed,
	}

	s.logging.LogInfo(fmt.Sprintf("Fetched detail receipt for user %s and receipt ID %s successfully", userId, receiptId))
	return detailResponse, nil
}

func (s *receiptService) UpdateReceiptConfirmed(userId, receiptId string, confirmed bool) error {
	s.logging.LogInfo(fmt.Sprintf("Updating receipt confirmation status for receipt ID %s to %t", receiptId, confirmed))

	err := s.receiptRepository.UpdateReceiptConfirmed(receiptId, confirmed)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to update receipt confirmation status for receipt ID %s: %v", receiptId, err))
		return fmt.Errorf("failed to update receipt confirmation status: %w", err)
	}

	// insert transaction if confirmed
	if confirmed {
		receipt, err := s.GetDetailReceiptUserById(userId, receiptId)
		if err != nil {
			s.logging.LogError(fmt.Sprintf("Failed to fetch receipt for ID %s: %v", receiptId, err))
			return fmt.Errorf("failed to fetch receipt: %w", err)
		}

		for _, item := range receipt.Items {
			// Handle empty category ID
			if item.CategoryId != nil && *item.CategoryId == "" {
				s.logging.LogInfo("Category ID is empty, setting it to nil")
				item.CategoryId = nil
			}

			// Create transaction
			categoryId := ""
			if item.CategoryId != nil {
				categoryId = *item.CategoryId
			}

			dateNow := time.Now()
			newTransaction := &requests.TransactionRequest{
				UserId:               userId,
				Amount:               item.ItemPriceTotal,
				Description:          item.ItemName,
				CategoryId:           categoryId,
				Type:                 "expense",
				Source:               "receipt",
				TransactionDate:      receipt.TransactionDate,
				IsAutoCategorized:    true,
				AiCategoryConfidence: item.AiCategoryConfidence,
				CreatedAt:            dateNow,
				UpdatedAt:            dateNow,
				Confirmed:            false,
				Discount:             item.ItemDiscount,
			}

			err := s.transactionService.InsertTransaction(newTransaction)
			if err != nil {
				s.logging.LogError(fmt.Sprintf("Failed to insert transaction: %v", err))
				return fmt.Errorf("failed to insert transaction: %w", err)
			}
		}
	}

	s.logging.LogInfo(fmt.Sprintf("Receipt confirmation status for receipt ID %s updated successfully", receiptId))
	return nil
}

// cleanAIResponse removes markdown formatting and extracts JSON content from AI response
func (s *receiptService) cleanAIResponse(response string) string {
	// Remove common markdown code block patterns
	response = strings.ReplaceAll(response, "```json", "")
	response = strings.ReplaceAll(response, "```", "")

	// Trim whitespace
	response = strings.TrimSpace(response)

	// Find the first '{' and last '}' to extract JSON object
	startIndex := strings.Index(response, "{")
	if startIndex == -1 {
		return response // No JSON object found, return as is
	}

	lastIndex := strings.LastIndex(response, "}")
	if lastIndex == -1 || lastIndex <= startIndex {
		return response // Invalid JSON structure, return as is
	}

	// Extract the JSON content
	jsonContent := response[startIndex : lastIndex+1]

	return jsonContent
}
