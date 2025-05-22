package build

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func newTestAsset(path string, data string, meta map[string]any) *Asset {
	return &Asset{Path: path, Data: []byte(data), Meta: meta}
}

type MockTransformer struct {
	TransformFunc func(*Asset) error
	CalledCount   int
}

func (m *MockTransformer) Transform(asset *Asset) error {
	m.CalledCount++
	if m.TransformFunc != nil {
		return m.TransformFunc(asset)
	}
	return nil
}

func TestAssets_Transform(t *testing.T) {
	asset1 := newTestAsset("1.txt", "content", nil)
	asset2 := newTestAsset("2.txt", "content", nil)
	assets := Assets{asset1, asset2}

	transformer1 := &MockTransformer{
		TransformFunc: func(a *Asset) error {
			if a.Meta == nil {
				a.Meta = make(map[string]any)
			}
			a.Meta["transformedBy"] = "transformer1"
			return nil
		},
	}

	transformer2 := &MockTransformer{
		TransformFunc: func(a *Asset) error {
			if a.Meta == nil {
				a.Meta = make(map[string]any)
			}
			if val, ok := a.Meta["transformedBy"].(string); ok {
				a.Meta["transformedBy"] = val + ", transformer2"
			} else {
				a.Meta["transformedBy"] = "transformer2"
			}
			return nil
		},
	}

	err := assets.Transform(transformer1, transformer2)
	if err != nil {
		t.Fatalf("Transform() returned error: %v", err)
	}

	if transformer1.CalledCount != 2 {
		t.Errorf("transformer1 was called %d times, want 2", transformer1.CalledCount)
	}
	if transformer2.CalledCount != 2 {
		t.Errorf("transformer2 was called %d times, want 2", transformer2.CalledCount)
	}

	for i, asset := range assets {
		expectedMeta := "transformer1, transformer2"
		if meta, ok := asset.Meta["transformedBy"].(string); !ok || meta != expectedMeta {
			t.Errorf("Asset %d Meta[\"transformedBy\"] = %v, want %v", i, asset.Meta["transformedBy"], expectedMeta)
		}
	}

	failingTransformer := &MockTransformer{
		TransformFunc: func(a *Asset) error {
			return os.ErrPermission
		},
	}
	assetsSingle := Assets{newTestAsset("fail.txt", "", nil)}
	err = assetsSingle.Transform(failingTransformer)
	if err == nil {
		t.Errorf("Transform() with failing transformer did not return an error")
	} else if err != os.ErrPermission {
		t.Errorf("Transform() with failing transformer returned wrong error: got %v, want %v", err, os.ErrPermission)
	}
}

func TestAssets_Write(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "assets_test_write")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	asset1 := newTestAsset("/file1.txt", "content1", nil)
	asset2 := newTestAsset("/subdir/file2.txt", "content2", nil)

	assets := Assets{asset1, asset2}

	err = assets.Write(tmpDir)
	if err != nil {
		t.Fatalf("Write() returned error: %v", err)
	}

	diskPath1 := tmpDir + asset1.Path
	content1, err := os.ReadFile(diskPath1)
	if err != nil {
		t.Errorf("Failed to read %s: %v", diskPath1, err)
	}
	if string(content1) != string(asset1.Data) {
		t.Errorf("Content of %s = %s, want %s", diskPath1, string(content1), string(asset1.Data))
	}

	diskPath2 := tmpDir + asset2.Path
	content2, err := os.ReadFile(diskPath2)
	if err != nil {
		t.Errorf("Failed to read %s: %v", diskPath2, err)
	}
	if string(content2) != string(asset2.Data) {
		t.Errorf("Content of %s = %s, want %s", diskPath2, string(content2), string(asset2.Data))
	}

	if len(assets) > 0 {
		filePathAsDir := tmpDir + "/blocked_dir"
		err = os.WriteFile(filePathAsDir, []byte("i am a file"), 0644)
		if err != nil {
			t.Fatalf("Could not create blocking file: %v", err)
		}

		asset3 := newTestAsset("/blocked_dir/file3.txt", "content3", nil)
		assetsBlocking := Assets{asset3}
		err = assetsBlocking.Write(tmpDir)
		if err == nil {
			t.Errorf("Write() did not return error when path component (blocked_dir) is a file")
		}
	}
}

func TestAssets_Pop(t *testing.T) {
	asset1 := newTestAsset("/path/to/file1.txt", "content1", map[string]any{"tag": "odd"})
	asset2 := newTestAsset("/path/to/file2.txt", "content2", map[string]any{"tag": "even"})
	asset3 := newTestAsset("/path/to/file3.txt", "content3", map[string]any{"tag": "odd"})

	assets := Assets{asset1, asset2, asset3}
	originalAssets := Assets{asset1, asset2, asset3}

	// Filter to pop assets with tag "odd"
	oddTagFilter := func(a Asset) bool {
		tag, ok := a.Meta["tag"].(string)
		return ok && tag == "odd"
	}

	popped := assets.Pop(oddTagFilter)

	// Check popped assets
	expectedPopped := Assets{asset1, asset3}
	if !reflect.DeepEqual(popped, expectedPopped) {
		t.Errorf("Pop() popped = %v, want %v", popped, expectedPopped)
	}

	// Check remaining assets in the original slice
	expectedRemaining := Assets{asset2}
	if !reflect.DeepEqual(assets, expectedRemaining) {
		t.Errorf("Pop() remaining assets = %v, want %v", assets, expectedRemaining)
	}

	// Test with no filters (should pop nothing)
	assets = Assets{asset1, asset2, asset3}
	popped = assets.Pop()
	if len(popped) != 0 {
		t.Errorf("Pop() with no filters popped = %v, want empty", popped)
	}
	if !reflect.DeepEqual(assets, originalAssets) {
		t.Errorf("Pop() with no filters modified original assets = %v, want %v", assets, originalAssets)
	}

	// Test with a filter that matches nothing
	assets = Assets{asset1, asset2, asset3}
	neverPopFilter := func(a Asset) bool { return false }
	popped = assets.Pop(neverPopFilter)
	if len(popped) != 0 {
		t.Errorf("Pop() with neverPopFilter popped = %v, want empty", popped)
	}
	if !reflect.DeepEqual(assets, originalAssets) {
		t.Errorf("Pop() with neverPopFilter modified original assets = %v, want %v", assets, originalAssets)
	}

	// Test with a filter that matches everything
	assets = Assets{asset1, asset2, asset3}
	alwaysPopFilter := func(a Asset) bool { return true }
	popped = assets.Pop(alwaysPopFilter)
	if !reflect.DeepEqual(popped, originalAssets) {
		t.Errorf("Pop() with alwaysPopFilter popped = %v, want %v", popped, originalAssets)
	}
	if len(assets) != 0 {
		t.Errorf("Pop() with alwaysPopFilter remaining assets = %v, want empty", assets)
	}
}

func TestAssets_Filter(t *testing.T) {
	asset1 := newTestAsset("/path/to/file1.txt", "content1", map[string]any{"tag": "odd", "type": "A"})
	asset2 := newTestAsset("/path/to/file2.txt", "content2", map[string]any{"tag": "even", "type": "B"})
	asset3 := newTestAsset("/path/to/file3.txt", "content3", map[string]any{"tag": "odd", "type": "A"})

	assets := Assets{asset1, asset2, asset3}
	originalAssets := Assets{asset1, asset2, asset3}

	// Filter to get assets with tag "odd"
	oddTagFilter := func(a Asset) bool {
		tag, ok := a.Meta["tag"].(string)
		return ok && tag == "odd"
	}

	filtered := assets.Filter(oddTagFilter)

	// Check filtered assets
	expectedFiltered := Assets{asset1, asset3}
	if !reflect.DeepEqual(filtered, expectedFiltered) {
		t.Errorf("Filter() filtered = %v, want %v", filtered, expectedFiltered)
	}

	// Check that the original slice is unchanged
	if !reflect.DeepEqual(assets, originalAssets) {
		t.Errorf("Filter() modified original assets = %v, want %v", assets, originalAssets)
	}

	// Test with multiple filters
	typeAFilter := func(a Asset) bool {
		typeVal, ok := a.Meta["type"].(string)
		return ok && typeVal == "A"
	}
	filteredMultiple := assets.Filter(oddTagFilter, typeAFilter)
	expectedFilteredMultiple := Assets{asset1, asset3}
	if !reflect.DeepEqual(filteredMultiple, expectedFilteredMultiple) {
		t.Errorf("Filter() with multiple filters = %v, want %v", filteredMultiple, expectedFilteredMultiple)
	}

	// Test with no filters (should return all assets)
	filteredNone := assets.Filter()
	if !reflect.DeepEqual(filteredNone, originalAssets) {
		t.Errorf("Filter() with no filters = %v, want %v", filteredNone, originalAssets)
	}

	// Test with a filter that matches nothing
	neverFilter := func(a Asset) bool { return false }
	filteredNever := assets.Filter(neverFilter)
	if len(filteredNever) != 0 {
		t.Errorf("Filter() with neverFilter = %v, want empty", filteredNever)
	}
}

func TestAssets_ToMap(t *testing.T) {
	asset1 := newTestAsset("/path/to/file1.txt", "content1", map[string]any{"id": "id1", "tag": "odd"})
	asset2 := newTestAsset("/path/to/file2.txt", "content2", map[string]any{"id": "id2", "tag": "even"})
	asset3 := newTestAsset("/path/to/file3.txt", "content3", map[string]any{"tag": "odd"})
	asset4 := newTestAsset("/path/to/file4.txt", "content4", map[string]any{"id": 123})

	assets := Assets{asset1, asset2, asset3, asset4}

	assetMap := assets.ToMap("id")

	expectedMap := map[string]*Asset{
		"id1": asset1,
		"id2": asset2,
	}

	if !reflect.DeepEqual(assetMap, expectedMap) {
		t.Errorf("ToMap() map = %v, want %v", assetMap, expectedMap)
	}

	// Test with a key that doesn't exist
	emptyMap := assets.ToMap("nonexistent_key")
	if len(emptyMap) != 0 {
		t.Errorf("ToMap() with nonexistent key = %v, want empty map", emptyMap)
	}
}

func TestAssets_SetMetaFunc(t *testing.T) {
	asset1 := newTestAsset("/path/file1.txt", "data1", map[string]any{"id": "1"})
	asset2 := newTestAsset("/another/file2.txt", "data2", nil)
	asset3 := newTestAsset("/short.txt", "data3", map[string]any{})

	assets := Assets{asset1, asset2, asset3}

	pathPrefixFunc := func(a Asset) string {
		return filepath.Dir(a.Path)
	}

	assets.SetMetaFunc("dir", pathPrefixFunc)

	// Check asset1
	if dir, ok := asset1.Meta["dir"].(string); !ok || dir != "/path" {
		t.Errorf("SetMetaFunc() asset1.Meta[\"dir\"] = %v, want %v", asset1.Meta["dir"], "/path")
	}

	// Check asset2 (meta was nil, should be skipped)
	if asset2.Meta != nil {
		if _, ok := asset2.Meta["dir"]; ok {
			t.Errorf("SetMetaFunc() asset2.Meta[\"dir\"] was added, but Meta was initially nil and should be skipped. Meta: %v", asset2.Meta)
		}
	}

	// Check asset3
	if dir, ok := asset3.Meta["dir"].(string); !ok || dir != "/" { // filepath.Dir("/short.txt") is "/"
		t.Errorf("SetMetaFunc() asset3.Meta[\"dir\"] = %v, want %v", asset3.Meta["dir"], "/")
	}
}

func TestAssets_AddToMeta(t *testing.T) {
	asset1 := newTestAsset("1.txt", "", map[string]any{"key1": "val1"})
	asset2 := newTestAsset("2.txt", "", nil)              // Meta is nil
	asset3 := newTestAsset("3.txt", "", map[string]any{}) // Meta is empty

	assets := Assets{asset1, asset2, asset3}
	assets.AddToMeta("newKey", "newValue")

	// Check asset1
	if val, ok := asset1.Meta["newKey"].(string); !ok || val != "newValue" {
		t.Errorf("AddToMeta() asset1.Meta[\"newKey\"] = %v, want %v", asset1.Meta["newKey"], "newValue")
	}
	if val, ok := asset1.Meta["key1"].(string); !ok || val != "val1" {
		t.Errorf("AddToMeta() asset1.Meta[\"key1\"] was changed or removed")
	}

	// Check asset2 (meta was nil, should be skipped)
	if asset2.Meta != nil {
		if _, ok := asset2.Meta["newKey"]; ok {
			t.Errorf("AddToMeta() asset2.Meta[\"newKey\"] was added, but Meta was initially nil and should be skipped. Meta: %v", asset2.Meta)
		}
	}

	// Check asset3
	if val, ok := asset3.Meta["newKey"].(string); !ok || val != "newValue" {
		t.Errorf("AddToMeta() asset3.Meta[\"newKey\"] = %v, want %v", asset3.Meta["newKey"], "newValue")
	}
}
