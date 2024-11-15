package generator

import (
	"bufio"
	_ "embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/addvanced/go-enum/internal/parser"
)

// Embed the template into the binary
//
//go:embed enum.go.tmpl
var enumTemplate string

// contains is a helper function to check if a string contains another string.
func contains(base, substr string) bool {
	return strings.Contains(base, substr)
}

// defaultFor returns the default value for a given type.
func defaultFor(baseType string) string {
	switch baseType {
	case "string":
		return `""`
	case "float64", "float32":
		return "0.0"
	case "bool":
		return "false"
	default:
		return "0"
	}
}

func lower(s string) string {
	return strings.ToLower(s)
}

// GenerateEnum generates a Go file for the given Enum.
func GenerateEnum(outputDir, packageName string, enum parser.Enum) error {
	filePath := filepath.Join(outputDir, strings.TrimSpace(strings.ToLower(enum.TypeName))+"_enum.go")

	// Check if the file already exists and contains the header
	if fileContainsHeader(filePath, packageName, enum.TypeName) {
		// Skip generation
		return nil
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Parse the embedded template
	tmpl, err := template.New("enum").
		Funcs(template.FuncMap{
			"contains":   contains,
			"lower":      lower,
			"defaultFor": defaultFor,
		}).
		Parse(enumTemplate)
	if err != nil {
		return err
	}

	// Execute the template
	return tmpl.Execute(file, struct {
		PackageName string
		Enum        parser.Enum
	}{
		PackageName: packageName,
		Enum:        enum,
	})
}

// fileContainsHeader checks if a file contains the generated header.
func fileContainsHeader(filePath, packageName, typeName string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		// If the file doesn't exist, we won't skip it
		return false
	}
	defer file.Close()

	header := "// Package " + packageName + " adds an enum value and parsing functions for the enum type " + typeName + "."
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), header) {
			return true
		}
	}
	return false
}
