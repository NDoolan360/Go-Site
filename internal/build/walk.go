package build

import (
	"bytes"
	"io/fs"
	"strings"
)

type Build struct {
	Assets
}

func (build *Build) WalkDir(fsys fs.FS, root string, includeRoot bool) error {
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

			if includeRoot {
				filepath = "/" + filepath
			} else {
				filepath = strings.TrimPrefix(filepath, root)
			}

			build.Assets = append(build.Assets, &Asset{
				Path:       filepath,
				SourceRoot: root,
				Meta:       map[string]any{},
				Data:       bytes.TrimSpace(data),
			})

			return nil
		},
	)
}
