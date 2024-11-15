package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/addvanced/go-enum/internal/generator"
	enumparser "github.com/addvanced/go-enum/internal/parser"
)

func main() {
	// CLI flags
	outputDir := flag.String("o", "", "Output directory for generated files")
	packageName := flag.String("pkg", "", "Package name for generated files")
	inputFiles := flag.String("i", "*", "Input files")
	flag.Parse()

	// Get the list of files to parse
	files, err := filepath.Glob(*inputFiles)
	if err != nil {
		fmt.Printf("Error finding files: %v\n", err)
		os.Exit(1)
	}

	// Determine the directory and package if not provided
	if *outputDir == "" || *packageName == "" {
		fileDir, pkgName, err := inferDefaults(files[0])
		if err != nil {
			fmt.Printf("Error inferring defaults: %v\n", err)
			os.Exit(1)
		}

		if *outputDir == "" {
			*outputDir = fileDir
		}
		if *packageName == "" {
			*packageName = pkgName
		}
	}

	// Parse Go files for enums
	enums, err := enumparser.ParseEnums(files)
	if err != nil {
		fmt.Printf("Error parsing files: %v\n", err)
		os.Exit(1)
	}

	// Generate files for each enum
	for _, enum := range enums {
		err := generator.GenerateEnum(*outputDir, *packageName, enum)
		if err != nil {
			fmt.Printf("Error generating file for %s: %v\n", enum.TypeName, err)
			os.Exit(1)
		}
	}

	fmt.Println("Enum generation completed!")
}

// inferDefaults infers the output directory and package name from the first input file.
func inferDefaults(file string) (string, string, error) {
	// Get the absolute path of the file
	absPath, err := filepath.Abs(file)
	if err != nil {
		return "", "", err
	}

	// Get the directory of the file
	fileDir := filepath.Dir(absPath)

	// Parse the Go file to infer the package name
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.PackageClauseOnly)
	if err != nil {
		return "", "", err
	}

	return fileDir, node.Name.Name, nil
}
