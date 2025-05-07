package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/spf13/pflag"
)

var version = "1.0.0"

// readLines reads a file line by line and returns a slice of strings.
// If the file is not found, returns nil, nil — this is convenient for
// simply ignoring missing .gitignore / exclusions files.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func main() {
	var dir string
	var exclusionsFilePath string
	var extensionsCSV string
	var openResult bool

	pflag.StringVarP(&dir, "directory", "d", ".", "Directory to scan")
	pflag.StringVarP(&exclusionsFilePath, "exclusions", "e", "exclusions.txt", "Path to an additional ignore file (gitignore syntax)")
	pflag.StringVarP(&extensionsCSV, "extensions", "x", "", "Comma‑separated list of file extensions to include (without dots). Empty means all files.")
	pflag.BoolVarP(&openResult, "open", "o", false,
		"Open the generated file in the system’s default editor/viewer")
	pflag.Parse()

	// Parse the list of extensions into a map for O(1) checks.
	extSet := map[string]struct{}{}
	if extensionsCSV != "" {
		for _, ext := range strings.Split(extensionsCSV, ",") {
			ext = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(ext, ".")))
			if ext != "" {
				extSet[ext] = struct{}{}
			}
		}
	}

	// Collect ignore rules from <dir>/.gitignore and the -e file.
	var patterns []string
	if lines, _ := readLines(filepath.Join(dir, ".gitignore")); lines != nil {
		patterns = append(patterns, lines...)
	}
	if lines, _ := readLines(exclusionsFilePath); lines != nil {
		patterns = append(patterns, lines...)
	}

	ig := ignore.CompileIgnoreLines(patterns...)

	outputFile := "go_files.txt"
	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer out.Close()

	// Traverse the file tree.
	err = filepath.Walk(dir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		// Do not add the root .gitignore to the result
		if relPath == ".gitignore" {
			return nil
		}

		// gitignore / exclusions filter.
		if ig != nil && ig.MatchesPath(relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Only files are of interest (we do not write directories to output).
		if info.IsDir() {
			return nil
		}

		// Filter by extensions if set.
		if len(extSet) > 0 {
			ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
			if _, ok := extSet[ext]; !ok {
				return nil
			}
		}

		// Write the path and the content.
		if _, err := out.WriteString(relPath + "\n"); err != nil {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if _, err := out.WriteString("```\n" + string(content) + "\n```\n"); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the path:", err)
		return
	}

	fmt.Println("File list has been written to", outputFile)

	if openResult {
		if err := openFile(outputFile); err != nil {
			fmt.Println("Cannot open file:", err)
		}
	}
}

func openFile(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default: // linux, *bsd
		cmd = exec.Command("xdg-open", path)
	}
	// Start asynchronously so we don't wait for the editor to close
	return cmd.Start()
}
