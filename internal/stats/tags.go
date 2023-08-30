package stats

import (
	"reflect"
	"strings"

	"github.com/dipdup-net/go-lib/hasura"
	"github.com/pkg/errors"
)

var (
	Tables = make(map[string]Table)

	ValidFuncs = map[string]struct{}{
		"count": {},
		"avg":   {},
		"max":   {},
		"min":   {},
		"sum":   {},
	}
)

type Table struct {
	Columns map[string]Column
}

type Column struct {
	Functions  map[string]struct{}
	Filterable bool
}

func Init(models ...any) error {
	for i := range models {
		if err := InitModel(models[i]); err != nil {
			return err
		}
	}

	return nil
}

const (
	statsTagName = "stats"
	bunTagName   = "bun"
)

func InitModel(model any) error {
	typ := reflect.TypeOf(model)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	var (
		tableName string
		table     = Table{
			Columns: make(map[string]Column, 0),
		}
	)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		bunTag := field.Tag.Get(bunTagName)
		if bunTag == "" || bunTag == "-" {
			continue
		}
		bunTagValues := strings.Split(bunTag, ",")
		if len(bunTagValues) == 0 {
			continue
		}

		switch field.Name {
		case "BaseModel":
			tableName = strings.TrimPrefix(bunTagValues[0], "table:")
			if tableName == "" {
				tableName = hasura.ToSnakeCase(typ.Name())
			}
		default:
			statsTag := field.Tag.Get(statsTagName)
			if statsTag == "" || statsTag == "-" {
				continue
			}
			statsTagValues := strings.Split(statsTag, ",")
			if len(statsTagValues) == 0 {
				continue
			}

			column := Column{
				Functions: make(map[string]struct{}),
			}
			for _, value := range statsTagValues {
				switch {
				case strings.HasPrefix(value, "func:"):
					funcs := strings.Split(strings.TrimPrefix(value, "func:"), " ")

					for _, fnc := range funcs {
						if _, ok := ValidFuncs[fnc]; !ok {
							return errors.Errorf("unknown stats function: %s", value)
						}
						column.Functions[fnc] = struct{}{}
					}
				case value == "filterable":
					column.Filterable = true
				}
			}
			columnName := bunTagValues[0]
			if columnName == "" {
				columnName = hasura.ToSnakeCase(field.Name)
			}
			table.Columns[columnName] = column
		}
	}

	Tables[tableName] = table

	return nil
}
