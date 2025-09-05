// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package responses

type Relayer struct {
	Name      string   `example:"cosmosrescue"                                                                               format:"string"  json:"name"                swaggertype:"string"`
	Logo      string   `example:"https://raw.githubusercontent.com/irisnet/iob-registry/main/relayers/cosmosrescue/logo.png" format:"string"  json:"logo"                swaggertype:"string"`
	Addresses []string `example:"celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60"                                            json:"addresses" swaggertype:"array,string"`

	Contact *Contact `json:"contact,omitempty"`
}

type Contact struct {
	Website  string `example:"https://cosmosrescue.com"         json:"website,omitempty"  swaggertype:"string"`
	Github   string `example:"https://github.com/cosmosrescue"  json:"github,omitempty"   swaggertype:"string"`
	Twitter  string `example:"https://x.com/cosmosrescue"       json:"twitter,omitempty"  swaggertype:"string"`
	Telegram string `example:"https://t.me/cosmosrescue"        json:"telegram,omitempty" swaggertype:"string"`
	Discord  string `example:"https://discord.cosmosrescue.com" json:"discord,omitempty"  swaggertype:"string"`
	Medium   string `example:"https://medium.com/cosmosrescue"  json:"medium,omitempty"   swaggertype:"string"`
}
