package prompt

const (
	// ChatSystemPrompt is the main system prompt for AI chat assistant
	ChatSystemPrompt = "You are a financial assistant. Provide helpful and accurate responses to user queries."

	// TitleGenerationSystemPrompt is the system prompt for generating chat titles
	TitleGenerationSystemPrompt = "You are a helpful assistant that creates concise, descriptive titles for conversations. Respond with only the title, no additional text."

	// TitleGenerationUserPromptTemplate is the template for title generation user prompt
	TitleGenerationUserPromptTemplate = `Generate a concise, descriptive title (max 50 characters) for a conversation that starts with this message:

"%s"

Requirements:
- Maximum 50 characters
- Clear and descriptive
- No quotes or special formatting
- Summarize the main topic or intent
- Professional tone

Title:`
)
