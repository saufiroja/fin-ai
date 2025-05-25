package models

type LogMessageModel struct {
	LogMessageId string `json:"log_message_id"`
	UserId       string `json:"user_id"`
	Message      string `json:"message"`
	Response     string `json:"response"`
	InputToken   int    `json:"input_token"`
	OutputToken  int    `json:"output_token"`
	Topic        string `json:"topic"`
	Model        string `json:"model"`
}
