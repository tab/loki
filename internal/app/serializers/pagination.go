package serializers

type PaginationMeta struct {
	Page  int `json:"page"`
	Per   int `json:"per"`
	Total int `json:"total"`
}

type PaginationResponse[T interface{}] struct {
	Data []T            `json:"data"`
	Meta PaginationMeta `json:"meta"`
}
