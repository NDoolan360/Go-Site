package build

import (
	"bytes"
	"io/fs"
	"strings"
)

type Build struct {
	Assets
}

func (build *Build) WalkDir(fsys fs.FS, root string) error {
	return fs.WalkDir(fsys, root,
		func(filepath string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			data, err := fs.ReadFile(fsys, filepath)
			if err != nil {
				return err
			}

			build.Assets = append(build.Assets, &Asset{
				Path:       strings.TrimPrefix(filepath, root),
				SourceRoot: root,
				Meta:       map[string]any{},
				Data:       bytes.TrimSpace(data),
			})

			return nil
		},
	)
}
