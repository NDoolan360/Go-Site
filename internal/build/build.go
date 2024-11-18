package build

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"text/template"
)

type Build struct {
	Dirs []string
	Pages
	Templates map[string][]byte
	OutDir    string
}

type Page struct {
	FileType string
	Data     *[]byte
	PageMeta map[string]any
}

type Pages map[string]Page

type Content struct {
	Page map[string]any
	Site map[string]any
}

func NewBuild(out string) Build {
	return Build{
		Dirs:      []string{},
		Pages:     Pages{},
		Templates: map[string][]byte{},
		OutDir:    out,
	}
}

func (build Build) StartWith(fs fs.FS) error {
	if rmErr := os.RemoveAll(build.OutDir); rmErr != nil {
		return rmErr
	}
	if cpErr := os.CopyFS(build.OutDir, fs); cpErr != nil {
		return cpErr
	}
	return nil
}

func (build *Build) CollectFiles(files fs.FS, path string) error {
	if walkErr := fs.WalkDir(files, path, build.Walk(files, path)); walkErr != nil {
		return walkErr
	}
	return nil
}

func (build *Build) Walk(fileSystem fs.FS, root string) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			path = strings.Replace(path, root, "", 1)
			build.Dirs = append(build.Dirs, path)
			return nil
		}

		fileContent, readErr := fs.ReadFile(fileSystem, path)
		if readErr != nil {
			return readErr
		}
		path = strings.Replace(path, root, "", 1)

		// TODO - handle other file types like html
		if strings.HasSuffix(path, "_template.html") {
			build.Templates[path] = fileContent
		} else if strings.HasSuffix(path, ".md") {
			meta, err := GetMeta(fileContent)
			if err != nil {
				return err
			}
			if meta == nil {
				meta = map[string]any{}
			}

			var ok bool
			if path, ok = strings.CutSuffix(path, ".md"); ok {
				path = path + ".html"
			}

			// Default meta values
			if _, ok := meta["Title"]; !ok {
				meta["Title"] = path
			}
			if _, ok := meta["Template"]; !ok {
				meta["Template"] = "/base_template.html"
			}
			if _, ok := meta["Scripts"]; !ok {
				meta["Scripts"] = []interface{}{}
			}

			build.Pages[path] = Page{
				Data:     &fileContent,
				FileType: "md",
				PageMeta: meta,
			}
		}

		return nil
	}
}

func (build *Build) ApplyTemplates(siteMeta map[string]any) error {
	content := Content{
		Site: siteMeta,
	}
	for _, page := range build.Pages {
		content.Page = page.PageMeta
		if content.Site["Env"] == "dev" {
			content.Page["Scripts"] = append(content.Page["Scripts"].([]interface{}), "/static/scripts/reload.js")
		}

		if page.FileType == "md" {
			if err := ApplyTemplate(page.Data, content); err != nil {
				return err
			}

			if err := ToHtml(page.Data); err != nil {
				return err
			}
		}

		tmpl, err := template.New("main").Parse(string(*page.Data))
		if err != nil {
			return err
		}

		baseTemplateStr, ok := page.PageMeta["Template"].(string)
		if !ok {
			return errors.New("Template meta not found")
		}

		baseTemplate, ok := build.Templates[baseTemplateStr]
		if !ok {
			return errors.New(fmt.Sprintf("Template %s not found", baseTemplateStr))
		}

		layout, err := tmpl.New("index").Parse(string(baseTemplate))
		if err != nil {
			return err
		}

		buf := &bytes.Buffer{}
		if err := layout.ExecuteTemplate(buf, "index", content); err != nil {
			return err
		}

		*page.Data = buf.Bytes()
	}
	return nil
}

func (build Build) WriteFiles() error {
	for _, dir := range build.Dirs {
		os.MkdirAll(build.OutDir+dir, 0777)
	}

	for path, page := range build.Pages {
		os.WriteFile(build.OutDir+path, *page.Data, 0777)
	}
	return nil
}

func (pages Pages) GetBlogPosts() (blogPosts BlogPosts) {
	for path, page := range pages {
		if blogPost, ok := NewBlogPost(path, page.PageMeta); ok {
			blogPosts = append(blogPosts, blogPost)
		}
	}
	return blogPosts
}
