package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"mime/multipart"
	"os"
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
)

type receiptService struct {
	receiptRepository  receipt.ReceiptStorer
	transactionService transaction.TransactionManager
	logMessageService  log_message.LogMessageManager
	categoryService    categories.CategoryManager
	minioClient        minio.MinioManager
	logging            logging.Logger
	openaiClient       llm.OpenAI
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
	}
}

func (s *receiptService) UploadReceipt(filePath *multipart.FileHeader, userId string) error {
	s.logging.LogInfo(fmt.Sprintf("Uploading receipt for user %s from file %s", userId, filePath.Filename))

	file, err := filePath.Open()
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to open multipart file %s: %v", filePath.Filename, err))
		return fmt.Errorf("failed to open multipart file: %w", err)
	}
	defer file.Close()
	// Decode the image from the multipart file
	img, err := imaging.Decode(file)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to decode image file %s: %v", filePath.Filename, err))
		return fmt.Errorf("failed to decode image file: %w", err)
	}
	// Optimize image for OCR
	img = s.optimizeImageForOCR(img)

	// Convert optimized image to bytes for base64 encoding
	var optimizedImageBytes []byte
	buf := new(bytes.Buffer)

	// Determine format based on original filename
	var format imaging.Format = imaging.JPEG
	if strings.HasSuffix(strings.ToLower(filePath.Filename), ".png") {
		format = imaging.PNG
	}

	err = imaging.Encode(buf, img, format)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to encode optimized image: %v", err))
		return fmt.Errorf("failed to encode optimized image: %w", err)
	}
	optimizedImageBytes = buf.Bytes()

	// get all categories
	reqCategoryQuery := &requests.GetAllCategoryQuery{
		Limit:  100, // Set a reasonable limit for categories
		Offset: 0,   // Start from the beginning
	}
	categories, _ := s.categoryService.FindAllCategories(reqCategoryQuery)
	dateNow := time.Now()
	// size < 10MB
	if filePath.Size > 10*1024*1024 {
		s.logging.LogError(fmt.Sprintf("File %s exceeds the maximum size limit of 10MB", filePath.Filename))
		return fmt.Errorf("file exceeds the maximum size limit of 10MB")
	}

	categoriesStrings := make([]string, len(categories.Categories))
	for i, category := range categories.Categories {
		categoriesStrings[i] = fmt.Sprintf("%s (%s)", category.Name, category.CategoryId)
	}
	categoriesOfString := strings.Join(categoriesStrings, ", ")

	// Ensure bucket exists before proceeding
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

	// Set the object name for the file in MinIO
	s.objectName = fmt.Sprintf("%s/%s", userId, filePath.Filename)

	// check if the file exists, if it does, skip the upload
	exists, err := s.minioClient.FileExists(s.bucketName, s.objectName)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Error checking file existence: %v", err))
		return err
	}
	if !exists {
		// Upload the file to MinIO
		err = s.minioClient.UploadFileFromMultipart(s.bucketName, s.objectName, filePath)
		if err != nil {
			s.logging.LogError(fmt.Sprintf("Failed to upload file %s to bucket %s: %v", filePath.Filename, s.bucketName, err))
			return fmt.Errorf("failed to upload file: %w", err)
		}
		s.logging.LogInfo(fmt.Sprintf("File %s uploaded successfully to bucket %s", filePath.Filename, s.bucketName))
	}

	// Use optimized image bytes instead of reading from MinIO
	base64Image := base64.StdEncoding.EncodeToString(optimizedImageBytes)

	// Save optimized image to images folder for debugging/verification
	imagesDir := "images"
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to create images directory: %v", err))
		// Don't return error here as this is optional debugging feature
	} else {
		// Generate unique filename for the processed image
		outputFileName := fmt.Sprintf("optimized_%s_%s", userId, filePath.Filename)
		outputPath := fmt.Sprintf("%s/%s", imagesDir, outputFileName)

		if err := imaging.Save(img, outputPath); err != nil {
			s.logging.LogError(fmt.Sprintf("Failed to save optimized image: %v", err))
			// Don't return error here as this is optional debugging feature
		} else {
			s.logging.LogInfo(fmt.Sprintf("Optimized image saved to: %s", outputPath))
		}
	}

	// image type
	imageType := "image/png" // default
	if strings.HasSuffix(s.objectName, ".jpg") || strings.HasSuffix(s.objectName, ".jpeg") {
		imageType = "image/jpeg"
	}

	imageBase64 := fmt.Sprintf("data:%s;base64,%s", imageType, base64Image)

	// Process the receipt using LLM
	messagePrompt := []openai.ChatCompletionMessageParamUnion{
		// System message
		{
			OfSystem: &openai.ChatCompletionSystemMessageParam{
				Content: openai.ChatCompletionSystemMessageParamContentUnion{
					OfString: param.Opt[string]{
						Value: prompt.ReceiptExtractionSystemPrompt,
					},
				},
			},
		},
		// User message with image and text
		{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfArrayOfContentParts: []openai.ChatCompletionContentPartUnionParam{
						{
							OfImageURL: &openai.ChatCompletionContentPartImageParam{
								ImageURL: openai.ChatCompletionContentPartImageImageURLParam{
									URL:    imageBase64,
									Detail: "high",
								},
							},
						},
						{
							OfText: &openai.ChatCompletionContentPartTextParam{
								Text: fmt.Sprintf(prompt.ReceiptExtractionUserPromptTemplate, categoriesOfString),
							},
						},
					},
				},
			},
		},
	}
	responseAi, err := s.openaiClient.SendChat(context.Background(), "gpt-4o-mini", messagePrompt)
	fmt.Printf("LLM response: %s\n", responseAi.Response)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to process receipt with LLM: %v", err))
		return fmt.Errorf("failed to process receipt: %w", err)
	}

	// Convert response to string
	var responseString string
	if resp, ok := responseAi.Response.(string); ok {
		responseString = resp
	} else {
		s.logging.LogError("Failed to convert LLM response to string")
		return fmt.Errorf("failed to convert LLM response to string")
	}

	var extractedData responses.ReceiptExtractionResponse
	err = json.Unmarshal([]byte(responseString), &extractedData)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to parse JSON response: %v", err))
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Convert messagePrompt to JSON string for logging
	messagePromptJSON, err := json.Marshal(messagePrompt)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to marshal message prompt: %v", err))
		return fmt.Errorf("failed to marshal message prompt: %w", err)
	}

	// log the AI response
	logMessage := &models.LogMessage{
		LogMessageId: ulid.Make().String(),
		UserId:       userId,
		Message:      string(messagePromptJSON),
		Response:     responseString,
		InputToken:   responseAi.InputToken,
		OutputToken:  responseAi.OutputToken,
		Topic:        "receipt_extraction",
		Model:        "gpt-4o-mini",
		CreatedAt:    dateNow,
		UpdatedAt:    dateNow,
	}

	err = s.logMessageService.InsertLogMessage(logMessage)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to insert log message: %v", err))
		return fmt.Errorf("failed to insert log message: %w", err)
	}

	// create embedding for the extracted receipt text
	extractedReceiptJSON, err := json.Marshal(extractedData.ExtractedReceipt)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to marshal extracted receipt: %v", err))
		return fmt.Errorf("failed to marshal extracted receipt: %w", err)
	}

	input := openai.EmbeddingNewParamsInputUnion{
		OfString: param.NewOpt(string(extractedReceiptJSON)), // Use the extracted receipt text for embedding
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
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Create a new receipt model
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
		Confirmed:                 false, // Assuming new receipts are not confirmed by default
		TransactionDate:           dateNow,
		CreatedAt:                 dateNow,
		UpdatedAt:                 dateNow,
	}

	// Insert the receipt into the database
	err = s.insertReceipt(receiptModel)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to insert receipt: %v", err))
		return fmt.Errorf("failed to insert receipt: %w", err)
	}

	// insert receipt items
	for _, item := range extractedData.ExtractedReceipt.Items {
		if *item.CategoryId == "" {
			// null category ID
			s.logging.LogInfo("Category ID is empty, setting it to nil")
			item.CategoryId = nil
		}
		newTransaction := &requests.TransactionRequest{
			UserId:               userId,
			Amount:               item.ItemPriceTotal,
			Description:          item.ItemName,
			CategoryId:           *item.CategoryId,
			Type:                 "expense", // Assuming all transactions from receipts are expenses
			Source:               "receipt", // Assuming the source is 'receipt'
			TransactionDate:      dateNow,
			IsAutoCategorized:    true,
			AiCategoryConfidence: item.AiCategoryConfidence,
			CreatedAt:            dateNow,
			UpdatedAt:            dateNow,
			Confirmed:            false, // Assuming new transactions are not confirmed by default
		}
		err = s.transactionService.InsertTransaction(newTransaction)
		if err != nil {
			s.logging.LogError(fmt.Sprintf("Failed to insert transaction: %v", err))
			return fmt.Errorf("failed to insert transaction: %w", err)
		}
		receiptItem := &models.ReceiptItem{
			ReceiptItemId:  ulid.Make().String(),
			ReceiptId:      receiptModel.ReceiptId,
			ItemName:       item.ItemName,
			ItemQuantity:   item.ItemQuantity,
			ItemPrice:      item.ItemPrice,
			ItemPriceTotal: item.ItemPriceTotal,
			ItemDiscount:   item.ItemDiscount,
			CreatedAt:      dateNow,
			UpdatedAt:      dateNow,
		}

		err = s.receiptRepository.InsertReceiptItem(receiptItem)
		if err != nil {
			s.logging.LogError(fmt.Sprintf("Failed to insert receipt item: %v", err))
			return fmt.Errorf("failed to insert receipt item: %w", err)
		}
		s.logging.LogInfo(fmt.Sprintf("Receipt item %s inserted successfully for receipt %s", item.ItemName, receiptModel.ReceiptId))
	}

	s.logging.LogInfo(fmt.Sprintf("Receipt for user %s uploaded and processed successfully", userId))

	return nil
}

func (s *receiptService) optimizeImageForOCR(img image.Image) image.Image {
	// 1. Auto-rotate if needed (detect orientation)
	img = s.autoRotateImage(img)

	// 2. Resize to optimal size for both OCR and AI vision
	img = imaging.Resize(img, 1600, 0, imaging.Lanczos)

	// 3. Convert to grayscale for better text recognition
	img = imaging.Grayscale(img)

	// 4. Apply gamma correction for better contrast
	img = imaging.AdjustGamma(img, 1.2)

	// 5. Enhance contrast for text clarity
	img = imaging.AdjustContrast(img, 35)

	// 6. Brightness adjustment for white background
	img = imaging.AdjustBrightness(img, 10)

	// 7. Sharpen for text clarity
	img = imaging.Sharpen(img, 1.0)

	// 8. Apply threshold for binary image (optional)
	// img = s.applyThreshold(img)

	return img
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
