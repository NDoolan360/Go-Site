package build

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func (build *Build) FromGit(url string, branch string, outDir string) error {
	root := "tmp/" + outDir

	_, err := git.PlainClone(root, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.NoRecurseSubmodules,
		SingleBranch:      true,
		ReferenceName:     plumbing.NewBranchReferenceName(branch),
	})
	if err != nil && err != git.ErrRepositoryAlreadyExists {
		return fmt.Errorf("could not clone repository: %v", err)
	}

	return build.WalkDir(os.DirFS("tmp"), outDir, true)

	// wt, err := repo.Worktree()
	// if err != nil {
	// 	return fmt.Errorf("could not get worktree: %v", err)
	// }

	// if err := wt.Checkout(&git.CheckoutOptions{
	// 	Branch: plumbing.NewBranchReferenceName(branch),
	// }); err != nil {
	// 	return fmt.Errorf("could not checkout tag: %v", err)
	// }

	// // Walk the file system and add files to the build
	// return util.Walk(wt.Filesystem, ".",
	// 	func(filepath string, file fs.FileInfo, err error) error {
	// 		if file.IsDir() {
	// 			return nil
	// 		}

	// 		// Get the file content
	// 		content, err := util.ReadFile(wt.Filesystem, filepath)
	// 		if err != nil {
	// 			return fmt.Errorf("could not read git worktree file %s: %v", filepath, err)
	// 		}

	// 		// Add the file to the build
	// 		build.Assets = append(build.Assets, &Asset{
	// 			Path: "/" + outDir + "/" + filepath,
	// 			Meta: map[string]any{},
	// 			Data: bytes.TrimSpace(content),
	// 		})

	// 		return nil
	// 	},
	// )
}
