package generate

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func RenderTemplates(srcDir, dstDir string, data any) error {
	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Only process .tmpl files.
		if filepath.Ext(path) != ".tmpl" {
			return nil
		}

		// Preserve the relative directory structure.
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// Remove the .tmpl suffix.
		rel = strings.TrimSuffix(rel, ".tmpl")
		outPath := filepath.Join(dstDir, rel)

		// Create the destination directory.
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}

		// Parse the template.
		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return err
		}

		// Create and render the output file.
		out, err := os.Create(outPath)
		if err != nil {
			return err
		}

		execErr := tmpl.Execute(out, data)
		closeErr := out.Close()

		if execErr != nil {
			return execErr
		}
		return closeErr
	})
}
