package build

import (
	"bytes"
	"fmt"
	"maps"
	"path"
	"text/template"
)

type TemplateTransformer struct {
	WrapperTemplate *WrapperTemplate
	Components      map[string]*Asset
}

type WrapperTemplate struct {
	Template       *Asset
	ChildBlockName string
}

func (t TemplateTransformer) Transform(asset *Asset, params map[string]any) error {
	if asset.Path == "/base_template.html" {
		for name, component := range t.Components {
			fmt.Println(component.Path, name)
		}
		return nil
	}

	var primarySource []byte
	templateMeta := map[string]any{
		"Global": params,
		"Asset":  asset.Meta,
	}

	if t.WrapperTemplate != nil {
		primarySource = t.WrapperTemplate.Template.Data
		templateMeta["WrapperTemplate"] = t.WrapperTemplate.Template.Meta
		maps.Copy(templateMeta, t.WrapperTemplate.Template.Meta)
	} else {
		primarySource = asset.Data
	}

	maps.Copy(templateMeta, asset.Meta)

	tmpl, err := template.New("").Parse(string(primarySource))
	if err != nil {
		return err
	}

	for name, component := range t.Components {
		if path.Ext(component.Path) != path.Ext(asset.Path) {
			continue
		}
		tmpl, err = tmpl.New(name).Parse(string(component.Data))
		if err != nil {
			return err
		}
	}

	if t.WrapperTemplate != nil {
		tmpl, err = tmpl.New(t.WrapperTemplate.ChildBlockName).Parse(string(asset.Data))
		if err != nil {
			return err
		}
	}

	buf := &bytes.Buffer{}
	tmpl.ExecuteTemplate(buf, "", templateMeta)

	asset.Data = buf.Bytes()
	return nil
}
