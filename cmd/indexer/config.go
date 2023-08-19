package main

import (
	"bytes"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Version  string   `yaml:"version" validate:"required"`
	Database Database `yaml:"database" validate:"required"`
	Indexer  Indexer  `yaml:"indexer"`
}

type Database struct {
	Path       string `yaml:"path,omitempty"`
	Kind       string `yaml:"kind" validate:"required,oneof=sqlite postgres mysql clickhouse elasticsearch"`
	Host       string `yaml:"host" validate:"required_with=Port User Database"`
	Port       int    `yaml:"port" validate:"required_with=Host User Database,gt=-1,lt=65535"`
	User       string `yaml:"user" validate:"required_with=Host Port Database"`
	Password   string `yaml:"password"`
	Database   string `yaml:"database" validate:"required_with=Host Port User"`
	SchemaName string `yaml:"schema_name"`
}

type Indexer struct {
	Name    string `yaml:"name" validate:"omitempty"`
	Timeout uint64 `yaml:"timeout" validate:"omitempty"`
	Node    *Node  `yaml:"node" validate:"omitempty"`
}

type Node struct {
	Url string `yaml:"url" validate:"omitempty,url"`
	Rps int    `yaml:"requests_per_second" validate:"omitempty,min=1"`
}

// Parse -
func Parse(filename string, output Config) error {
	buf, err := readFile(filename)
	if err != nil {
		return err
	}

	if err := yaml.NewDecoder(buf).Decode(output); err != nil {
		return err
	}

	return validator.New().Struct(output)
}

func readFile(filename string) (*bytes.Buffer, error) {
	if filename == "" {
		return nil, errors.Errorf("you have to provide configuration filename")
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "reading file %s", filename)
	}
	//expanded, err := expandVariables(data)
	//if err != nil {
	//	return nil, err
	//}
	return bytes.NewBuffer(data), nil
}
