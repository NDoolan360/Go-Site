package build

import (
	"os"
	"path"
	"slices"
	"strings"
)

type Asset struct {
	Path       string
	Data       []byte
	Meta       map[string]any
	SourceRoot string
}

type Assets []*Asset

type Transformer interface {
	Transform(*Asset, map[string]any) error
}

type Filter func(Asset) bool

func (assets Assets) Transform(params map[string]any, transformers ...Transformer) error {
	for _, transformer := range transformers {
		for _, asset := range assets {
			if err := transformer.Transform(asset, params); err != nil {
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

func (assets Assets) SortBy(metaDataField string) Assets {
	slices.SortFunc(assets, func(i, j *Asset) int {
		stringI, okI := i.Meta[metaDataField].(string)
		stringJ, okJ := j.Meta[metaDataField].(string)
		if okI && okJ {
			return strings.Compare(stringJ, stringI)
		}

		stringArrayI, okI := i.Meta[metaDataField].([]any)
		stringArrayJ, okJ := j.Meta[metaDataField].([]any)
		if okI && okJ {
			return len(stringArrayJ) - len(stringArrayI)
		}

		return 0
	})

	return assets
}

func (assets Assets) AddToMeta(metaKey string, value string) Assets {
	for _, asset := range assets {
		// If the meta is nil, skip
		if asset.Meta == nil {
			continue
		}

		asset.Meta[metaKey] = value
	}

	return assets
}

func (assets Assets) AddToMetaArray(metaKey string, value string) Assets {
	for _, asset := range assets {
		// If the meta is nil, skip
		if asset.Meta == nil {
			continue
		}

		if asset.Meta[metaKey] == nil {
			asset.Meta[metaKey] = []string{value}
		} else if _, ok := asset.Meta[metaKey].([]string); ok {
			asset.Meta[metaKey] = append(asset.Meta[metaKey].([]string), value)
		}
	}

	return assets
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
