package main

import (
	"path"

	"internal/build"
)

// Filter utils

func withParentDir(parent string) func(build.Asset) bool {
	return func(asset build.Asset) bool {
		dir := path.Dir(asset.Path)
		if path.Base(asset.Path) == path.Base(parent) && dir == path.Dir(parent) {
			return true
		} else if dir == "/" || dir == "." {
			return false
		} else {
			return withParentDir(parent)(build.Asset{Path: dir})
		}
	}
}

func withPath(filepath string) func(build.Asset) bool {
	return func(asset build.Asset) bool {
		return path.Clean(filepath) == path.Clean(asset.Path)
	}
}

func withExtensions(exts ...string) func(page build.Asset) bool {
	return func(page build.Asset) bool {
		for _, ext := range exts {
			if path.Ext(page.Path) == ext {
				return true
			}
		}
		return false
	}
}

func withoutExtensions(exts ...string) func(page build.Asset) bool {
	return func(page build.Asset) bool {
		for _, ext := range exts {
			if path.Ext(page.Path) == ext {
				return false
			}
		}
		return true
	}
}

func withMeta(key string) func(page build.Asset) bool {
	return func(page build.Asset) bool {
		val, ok := page.Meta[key]
		if ok {
			return val != false
		}
		return false
	}
}

func withoutMeta(key string) func(page build.Asset) bool {
	return func(page build.Asset) bool {
		val, ok := page.Meta[key]
		if ok {
			return val == false
		}
		return true
	}
}
