package responses

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
)

type Validator struct {
	Id          uint64 `example:"321"                                                    json:"id"           swaggertype:"integer"`
	Delegator   string `example:"celestia1un77nfm6axkhkupe8fk4xl6fd4adz3y5qk7ph6"        json:"delegator"    swaggertype:"string"`
	Address     string `example:"celestiavaloper1un77nfm6axkhkupe8fk4xl6fd4adz3y59fucpu" json:"address"      swaggertype:"string"`
	ConsAddress string `example:"E641C7A2C964833E556AEF934FBF166B712874B6"               json:"cons_address" swaggertype:"string"`

	Moniker  string `example:"Easy 2 Stake"                   json:"moniker"  swaggertype:"string"`
	Website  string `example:"https://www.easy2stake.com/"    json:"website"  swaggertype:"string"`
	Identity string `example:"2C877AC873132C91"               json:"identity" swaggertype:"string"`
	Contacts string `example:"security@0xfury.com"            json:"contacts" swaggertype:"string"`
	Details  string `example:"Some long text about validator" json:"details"  swaggertype:"string"`

	Rate              string `example:"0.03" json:"rate"                swaggertype:"string"`
	MaxRate           string `example:"0.1"  json:"max_rate"            swaggertype:"string"`
	MaxChangeRate     string `example:"0.01" json:"max_change_rate"     swaggertype:"string"`
	MinSelfDelegation string `example:"1"    json:"min_self_delegation" swaggertype:"string"`
}

func NewValidator(val storage.Validator) *Validator {
	if val.Id == 0 { // for genesis block
		return nil
	}
	return &Validator{
		Id:                val.Id,
		Delegator:         val.Delegator,
		Address:           val.Address,
		ConsAddress:       val.ConsAddress,
		Moniker:           val.Moniker,
		Website:           val.Website,
		Identity:          val.Identity,
		Contacts:          val.Contacts,
		Details:           val.Details,
		Rate:              val.Rate.String(),
		MaxRate:           val.MaxRate.String(),
		MaxChangeRate:     val.MaxChangeRate.String(),
		MinSelfDelegation: val.MinSelfDelegation.String(),
	}
}

type ShortValidator struct {
	Id          uint64 `example:"321"                                      json:"id"           swaggertype:"integer"`
	ConsAddress string `example:"E641C7A2C964833E556AEF934FBF166B712874B6" json:"cons_address" swaggertype:"string"`
	Moniker     string `example:"Easy 2 Stake"                             json:"moniker"      swaggertype:"string"`
}

func NewShortValidator(val storage.Validator) *ShortValidator {
	if val.Id == 0 { // for genesis block
		return nil
	}
	return &ShortValidator{
		Id:          val.Id,
		ConsAddress: val.ConsAddress,
		Moniker:     val.Moniker,
	}
}
