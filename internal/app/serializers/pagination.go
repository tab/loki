package serializers

type PaginationMeta struct {
	Page  uint64 `json:"page"`
	Per   uint64 `json:"per"`
	Total uint64 `json:"total"`
}

type PaginationResponse[T interface{}] struct {
	Data []T            `json:"data"`
	Meta PaginationMeta `json:"meta"`
}
