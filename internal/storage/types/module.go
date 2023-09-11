package types

// swagger:enum ModuleName
/*
	ENUM(
		auth,
		blob,
		crisis,
		distribution,
		indexer,
		gov,
		slashing,
		staking
	)
*/
//go:generate go-enum --marshal --sql --values
type ModuleName string
