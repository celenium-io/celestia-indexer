package types

// swagger:enum Status
/*
	ENUM(
		success,
		failed
	)
*/
//go:generate go-enum --marshal --sql --values
type Status string
