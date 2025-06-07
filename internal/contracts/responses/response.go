package responses

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponseAI struct {
	Response    string `json:"response"`
	InputToken  int    `json:"input_token"`
	OutputToken int    `json:"output_token"`
}

type ResponseEmbedding struct {
	Embeddings  string `json:"embeddings"`
	InputToken  int    `json:"input_token"`
	OutputToken int    `json:"output_token"`
}
