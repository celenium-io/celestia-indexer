package postgres

import (
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

		if _, err := s.Connection().DB().NewRaw(string(raw)).Exec(ctx); err != nil {
			return errors.Wrapf(err, "creating view '%s'", files[i].Name())
		}
	}

	return nil
}
