package indexer

type Config struct {
	Name    string `validate:"omitempty" yaml:"name"`
	Timeout int    `validate:"omitempty" yaml:"timeout"`
	Node    *Node  `validate:"omitempty" yaml:"node"`
}

type Node struct {
	Url string `validate:"omitempty,url"   yaml:"url"`
	Rps int    `validate:"omitempty,min=1" yaml:"requests_per_second"`
}
