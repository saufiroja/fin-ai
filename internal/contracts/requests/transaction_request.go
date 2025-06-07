package requests

type GetAllTransactionsQuery struct {
	Limit    int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Offset   int    `query:"offset" validate:"omitempty,min=0"`
	Category string `query:"category" validate:"omitempty"`
	Search   string `query:"search" validate:"omitempty"`
}
