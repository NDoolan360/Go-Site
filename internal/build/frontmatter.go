package build

import (
	"bytes"
	"maps"

	"github.com/adrg/frontmatter"
)

type CollectFrontMatter struct{}

func (p CollectFrontMatter) Transform(asset *Asset, _ map[string]any) error {
	var meta map[string]any

	rest, err := frontmatter.Parse(bytes.NewReader(asset.Data), &meta)
	if err != nil {
		return err
	}

	maps.Copy(asset.Meta, meta)
	asset.Data = rest

	return nil
}
