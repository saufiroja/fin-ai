package prompt

const (
	// ReceiptExtractionSystemPrompt is the system prompt for receipt extraction
	ReceiptExtractionSystemPrompt = `You are an expert Indonesian receipt analysis AI that specializes in extracting precise data from Indonesian retail receipts (Alfamart, Indomaret, Hypermart, etc.). You must analyze receipt images methodically and extract all financial information with 100% accuracy. Your response must be valid JSON only, without any additional text, formatting, or code blocks.

CRITICAL: You must read every single character on the receipt carefully, especially numbers and prices. Indonesian receipts often have specific formatting patterns. All monetary amounts must be in Indonesian Rupiah (Rp) as integers without decimal places.`

	// ReceiptExtractionUserPromptTemplate is the template for user prompt in receipt extraction
	ReceiptExtractionUserPromptTemplate = `
<categories>
%s
</categories>

<rules>
- extract all receipt data from the image, don't miss any details
- discounts are represented as negative numbers (e.g., -5000)
- all monetary amounts must be in Indonesian Rupiah as integers without decimal places
- ensure accuracy in number reading as Indonesian receipts use specific formatting
</rules>

RESPONSE FORMAT (JSON only, no code blocks):
{
    "extracted_receipt": {
        "merchant_name": "string", // Name of the merchant
        "sub_total": 0, // Subtotal amount in Rupiah (integer, no decimal)
        "total_discount": 0, // Total discount amount in Rupiah (integer, no decimal)
        "total_shopping": 0, // Total shopping amount after discounts in Rupiah (integer, no decimal)
        "transaction_date": "2024-01-01T00:00:00Z", // Date of the transaction
        "items": [
            {
                "category_id": "string", // Category ID from the available categories
                "item_name": "string", // Name of the item
                "item_quantity": 1, // Quantity of the item purchased
                "item_price": 0, // Price of the item in Rupiah (integer, no decimal)
                "item_price_total": 0, // Total price for the item (quantity * item_price) in Rupiah (integer, no decimal)
                "item_discount": 0, // Discount applied to the item in Rupiah (integer, no decimal)
                "ai_category_confidence": 0.0 // Confidence score for the AI's category prediction
            }
        ]
    }
}
    
MANDATORY:
- Response must be valid JSON, no additional text or formatting or code blocks
- Don't use backticks or any other formatting
`
)
