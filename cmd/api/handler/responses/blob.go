package responses

// structure is only for documentation
type Blob struct {
	Namespace    string `example:"AAAAAAAAAAAAAAAAAAAAAAAAAAAAs2bWWU6FOB0="     format:"base64"  json:"namespace"     swaggertype:"string"`
	Data         string `example:"b2sgZGVtbyBkYQ=="                             format:"base64"  json:"data"          swaggertype:"string"`
	ShareVersion int    `example:"0"                                            format:"integer" json:"share_version" swaggertype:"integer"`
	Commitment   string `example:"vbGakK59+Non81TE3ULg5Ve5ufT9SFm/bCyY+WLR3gg=" format:"base64"  json:"commitment"    swaggertype:"string"`
}
