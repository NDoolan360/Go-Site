// TODO write a string replacer transformer
package build

import "bytes"

type ReplacerTransformer struct {
	Replacements map[string]string
}

func (t ReplacerTransformer) Transform(asset *Asset, params map[string]any) error {
	for key, value := range t.Replacements {
		asset.Data = bytes.ReplaceAll(asset.Data, []byte(key), []byte(value))
	}
	return nil
}
