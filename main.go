package main

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday/v2"
)

func main() {
	htmlTemplate, err := readTemplateFile("file.html")
	if err != nil {
		log.Fatalf("Error reading HTML template: %v", err)
	}

	outputDir := "output"
	// Create the output folder if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
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
	err := filepath.Walk(inputDir, func(inputPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(inputPath, ".md") {
			relPath, err := filepath.Rel(inputDir, inputPath)
			if err != nil {
				return err
			}

			outputPath := filepath.Join(outputDir, strings.TrimSuffix(relPath, ".md")+".html")
			outputDirPath := filepath.Dir(outputPath)

			if err := os.MkdirAll(outputDirPath, os.ModePerm); err != nil {
				return err
			}

			if err := convertMarkdownToHTML(inputPath, outputPath, htmlTemplate); err != nil {
				return err
			}

			*generatedFiles = append(*generatedFiles, strings.TrimPrefix(outputPath, outputDir+"/"))

			log.Printf("Processed %s -> %s", inputPath, outputPath)
		}

		return nil
	})

	return err
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
