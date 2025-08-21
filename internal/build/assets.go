package build

import (
	"os"
	"path"
)

type Asset struct {
	Path string
	Data []byte
	Meta map[string]any
}

type Assets []*Asset

type Transformer interface {
	Transform(*Asset) error
}

type Filter func(Asset) bool

func (assets Assets) Transform(transformers ...Transformer) error {
	for _, transformer := range transformers {
		for _, asset := range assets {
			if err := transformer.Transform(asset); err != nil {
				return err
			}
		}
	}
	return nil
}

func (assets Assets) Write(outDir string) error {
	for _, asset := range assets {
		if err := os.MkdirAll(path.Dir(outDir+asset.Path), 0755); err != nil {
			return err
		}

		if err := os.WriteFile(outDir+asset.Path, asset.Data, 0644); err != nil {
			return err
		}
	}

	return nil
}

func (assets *Assets) Pop(filters ...Filter) Assets {
	if len(filters) == 0 {
		return nil
	}

	keep := make(Assets, 0, len(*assets))
	pop := make(Assets, 0, len(*assets))

	for _, asset := range *assets {
		willPop := true
		for _, filter := range filters {
			if !filter(*asset) {
				keep = append(keep, asset)
				willPop = false
				break
			}
		}
		if willPop {
			pop = append(pop, asset)
		}
	}

	*assets = keep
	return pop
}

// Returns a new Assets with only the assets that pass all the filters
func (assets Assets) Filter(filters ...Filter) Assets {
	if len(filters) == 0 {
		return assets
	}

	return assets.Pop(filters...)
}

func (assets Assets) ToMap(keyFromMeta string) map[string]*Asset {
	m := make(map[string]*Asset)
	for _, asset := range assets {
		if key, ok := asset.Meta[keyFromMeta].(string); ok {
			m[key] = asset
		}
	}
	return m
}

func (assets Assets) SetMetaFunc(metaKey string, fn func(Asset) string) Assets {
	for _, asset := range assets {
		// If the meta is nil, skip
		if asset.Meta == nil {
			continue
		}

		asset.Meta[metaKey] = fn(*asset)
	}

	return assets
}

func (assets Assets) AddToMeta(metaKey string, value string) Assets {
	return assets.SetMetaFunc(metaKey, func(asset Asset) string {
		return value
	})
}
