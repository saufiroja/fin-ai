package prompt

const (
	// TransactionConfidenceSystemPrompt is the system prompt for AI confidence scoring
	TransactionConfidenceSystemPrompt = "You are an AI assistant that analyzes transaction descriptions and provides confidence scores for category assignments. Respond with only a decimal number between 0.0 and 1.0, where 1.0 means very confident and 0.0 means not confident at all."

	// TransactionConfidenceUserPromptTemplate is the template for user prompt in confidence scoring
	TransactionConfidenceUserPromptTemplate = "How confident are you that the category '%s' is correct for this transaction: '%s'? Respond with only a number between 0.0 and 1.0."
)
