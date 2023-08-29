package responses

type Searchable interface {
	Block | Address | Namespace | Tx

	SearchType() string
}

type SearchResponse[T Searchable] struct {
	// Search result. Can be one of folowwing types: Block, Address, Namespace, Tx
	Result T `json:"result" swaggertype:"object"`
	// Result type which is in the result. Can be 'block', 'address', 'namespace', 'tx'
	Type string `json:"type"`
} //	@name	SearchResponse

func NewSearchResponse[T Searchable](val T) SearchResponse[T] {
	return SearchResponse[T]{
		Result: val,
		Type:   val.SearchType(),
	}
}
