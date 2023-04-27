package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/russross/blackfriday/v2"
)

func main() {
	htmlTemplate, err := readTemplateFile("file.html")
	if err != nil {
		log.Fatalf("Error reading HTML template: %v", err)
	}

	outputDir := "output"
	if err := clearOutputDirectory(outputDir); err != nil {
		log.Fatalf("Error clearing output directory: %v", err)
	}

	if err := copyFile("template/styles.css", filepath.Join(outputDir, "styles.css")); err != nil {
		log.Fatalf("Error copying styles.css: %v", err)
	}

	if err := copyFile("template/script.js", filepath.Join(outputDir, "script.js")); err != nil {
		log.Fatalf("Error copying script.js: %v", err)
	}

	generatedFiles := []string{}

	if err := processMarkdownFiles("input", outputDir, htmlTemplate, &generatedFiles); err != nil {
		log.Fatalf("Error processing Markdown files: %v", err)
	}

	if err := writeJSONFile(outputDir, generatedFiles); err != nil {
		log.Fatalf("Error writing JSON file: %v", err)
	}
}

func copyFile(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

func writeJSONFile(outputDir string, data []string) error {
	// Remove the output directory prefix from the file paths
	for i, path := range data {
		data[i] = strings.TrimPrefix(path, outputDir+"\\")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	jsonFilePath := filepath.Join(outputDir, "files.json")

	if err := ioutil.WriteFile(jsonFilePath, jsonData, 0644); err != nil {
		return err
	}

	return nil
}

func readTemplateFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filepath.Join("template", filename))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func processMarkdownFiles(inputDir, outputDir, htmlTemplate string, generatedFiles *[]string) error {
	err := filepath.Walk(inputDir, func(inputPath string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if file.IsDir() {
			return nil
		}

		if filepath.Ext(inputPath) != ".md" {
			return nil
		}

		// Get the output directory path and create the directory if it doesn't exist
		dir := filepath.Dir(inputPath)
		outputDirPath := filepath.Join(outputDir, strings.TrimPrefix(dir, inputDir))
		if err := os.MkdirAll(outputDirPath, os.ModePerm); err != nil {
			return err
		}

		// Sanitize the file name and change the extension to .html
		outputFile := sanitizeFileName(strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))) + ".html"
		outputPath := filepath.Join(outputDir, strings.TrimPrefix(dir, inputDir), outputFile)

		if err := convertMarkdownToHTML(inputPath, outputPath, htmlTemplate); err != nil {
			return err
		}

		*generatedFiles = append(*generatedFiles, outputPath)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func clearOutputDirectory(outputDir string) error {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		// If the output directory does not exist, create it
		return os.MkdirAll(outputDir, os.ModePerm)
	}

	// Delete all files and subdirectories within the output directory
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(outputDir, entry.Name())
		if entry.IsDir() {
			if err := os.RemoveAll(path); err != nil {
				return err
			}
		} else {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	}

	return nil
}

func sanitizeFileName(name string) string {
	// Remove special characters and spaces
	reg, err := regexp.Compile("[^a-zA-Z0-9 ]+")
	if err != nil {
		log.Fatalf("Error compiling regular expression: %v", err)
	}
	name = reg.ReplaceAllString(name, "")

	// Convert to lowercase and replace spaces with underscores
	name = strings.ReplaceAll(strings.ToLower(name), " ", "_")
	return name
}

func convertMarkdownToHTML(inputPath, outputPath, htmlTemplate string) error {
	md, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	// Convert Markdown to HTML
	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{})
	html := blackfriday.Run(md, blackfriday.WithRenderer(renderer))

	// Replace placeholders in the HTML template
	output := strings.Replace(htmlTemplate, "{{content}}", string(html), 1)

	// Write the output to a file
	if err := ioutil.WriteFile(outputPath, []byte(output), 0644); err != nil {
		return err
	}

	return nil
}
