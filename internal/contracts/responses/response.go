package responses

type Response struct {
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Total       int64 `json:"total"`
	CurrentPage int64 `json:"current_page"`
	TotalPages  int64 `json:"total_pages"`
}

type ResponseAI struct {
	Response    any `json:"response"`
	InputToken  int `json:"input_token"`
	OutputToken int `json:"output_token"`
}

type ResponseEmbedding struct {
	Embeddings  string `json:"embeddings"`
	InputToken  int    `json:"input_token"`
	OutputToken int    `json:"output_token"`
}
