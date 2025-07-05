package requests

type GetAllReceiptsQuery struct {
	Limit     int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Offset    int    `query:"offset" validate:"omitempty,min=0"`
	Search    string `query:"search" validate:"omitempty"`
	SortBy    string `query:"sort_by"`
	SortOrder string `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}
