package postgres

import (
	"bytes"
	"context"
	"os"
	"path/filepath"

	"github.com/dipdup-net/go-lib/database"
	"github.com/pkg/errors"
)

func (s Storage) createViews(ctx context.Context, conn *database.Bun) error {
	files, err := os.ReadDir(s.viewsDir)
	if err != nil {
		return err
	}

	for i := range files {
		if files[i].IsDir() {
			continue
		}

		path := filepath.Join(s.viewsDir, files[i].Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		queries := bytes.Split(raw, []byte{';'})
		if len(queries) == 0 {
			continue
		}

		for _, query := range queries {
			query = bytes.TrimLeft(query, "\n ")
			if len(query) == 0 {
				continue
			}
			if _, err := s.Connection().DB().NewRaw(string(query)).Exec(ctx); err != nil {
				return errors.Wrapf(err, "creating view '%s'", files[i].Name())
			}
		}

	}

	return nil
}
