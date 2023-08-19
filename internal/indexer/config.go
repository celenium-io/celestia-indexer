package indexer

type Config struct {
	Name    string `yaml:"name" validate:"omitempty"`
	Timeout int    `yaml:"timeout" validate:"omitempty"`
	Node    *Node  `yaml:"node" validate:"omitempty"`
}

type Node struct {
	Url string `yaml:"url" validate:"omitempty,url"`
	Rps int    `yaml:"requests_per_second" validate:"omitempty,min=1"`
}
