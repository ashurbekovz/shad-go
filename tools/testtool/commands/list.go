package commands

import (
	"bytes"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

// getPackageFiles returns absolute paths for all files in rootPackage and it's subpackages
// including tests and non-go files.
func getPackageFiles(rootPackage string, buildFlags []string) map[string]struct{} {
	cfg := &packages.Config{
		Dir:        rootPackage,
		Mode:       packages.NeedFiles,
		BuildFlags: buildFlags,
		Tests:      true,
	}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		log.Fatalf("unable to load packages %s: %s", rootPackage, err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	files := make(map[string]struct{})
	for _, p := range pkgs {
		for _, f := range p.GoFiles {
			if strings.HasSuffix(f, ".go") {
				files[f] = struct{}{}
			}
		}
	}

	return files
}

// listTestFiles returns absolute paths for all _test.go files of the package
// including the ones with "private" build tag.
func listTestFiles(rootPackage string) []string {
	files := getPackageFiles(rootPackage, []string{"-tags", "private"})
	var tests []string
	for f := range files {
		if strings.HasSuffix(f, "_test.go") {
			tests = append(tests, f)
		}
	}

	sort.Strings(tests)
	return tests
}

// listProtectedFiles returns absolute paths for all files of the package
// protected by "!change" build tag.
func listProtectedFiles(rootPackage string) []string {
	allFiles := getPackageFiles(rootPackage, nil)
	allFilesWithoutProtected := getPackageFiles(rootPackage, []string{"-tags", "change"})

	var protectedFiles []string
	for f := range allFiles {
		if _, ok := allFilesWithoutProtected[f]; !ok {
			protectedFiles = append(protectedFiles, f)
		}
	}

	sort.Strings(protectedFiles)
	return protectedFiles
}

// listPrivateFiles returns absolute paths for all files of the package
// protected by "private,solution" build tag.
func listPrivateFiles(rootPackage string) []string {
	allFiles := getPackageFiles(rootPackage, []string{})
	allWithPrivate := getPackageFiles(rootPackage, []string{"-tags", "private,solution"})

	var files []string
	for f := range allWithPrivate {
		if _, isPublic := allFiles[f]; !isPublic {
			files = append(files, f)
		}
	}

	if err := filepath.WalkDir(rootPackage, func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, ".proto") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			if bytes.Contains(content, []byte("//go:build solution")) {
				absPath, err := filepath.Abs(path)
				if err != nil {
					return err
				}

				files = append(files, absPath)
			}
		}

		return nil
	}); err != nil {
		log.Fatalf("filewalk failed: %v", err)
	}

	config, err := os.ReadFile(".private")
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed: %v", err)
	}

	for line := range strings.SplitSeq(string(config), "\n") {
		line = strings.Trim(line, " ")
		if line != "" {
			fname, _ := filepath.Abs(line)
			_, err := os.Stat(fname)
			if err == nil {
				files = append(files, fname)
			}
		}
	}

	sort.Strings(files)
	return files
}

func listTestsAndBinaries(rootDir string, buildFlags []string) (binaries, tests map[string]struct{}) {
	cfg := &packages.Config{
		Dir:        rootDir,
		Mode:       packages.NeedName | packages.NeedFiles,
		BuildFlags: buildFlags,
		Tests:      true,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		log.Fatalf("unable to load packages %s: %s", rootDir, err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	tests = map[string]struct{}{}
	binaries = map[string]struct{}{}

	for _, p := range pkgs {
		if p.Name != "main" {
			continue
		}

		if strings.HasSuffix(p.PkgPath, ".test") {
			tests[strings.TrimSuffix(p.PkgPath, ".test")] = struct{}{}
		} else {
			binaries[p.PkgPath] = struct{}{}
		}
	}

	return
}
