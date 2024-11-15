package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// Enum represents metadata for an enum type.
type Enum struct {
	TypeName string  // Name of the type (e.g., Color)
	BaseType string  // Base type (e.g., string, int)
	Values   []Value // Enum values
}

// Value represents a single value in the enum.
type Value struct {
	Name  string // Enum name (e.g., RED)
	Value string // Associated value (e.g., "#FF0000")
}

// ParseEnums parses Go files and extracts enums based on `// enum:` comments.
func ParseEnums(files []string) ([]Enum, error) {
	var enums []Enum

	for _, file := range files {
		if !strings.HasSuffix(file, ".go") || strings.HasSuffix(file, "_enum.go") {
			continue
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, file, nil, parser.AllErrors|parser.ParseComments)
		if err != nil {
			return nil, err
		}

		// Process declarations in the file
		for _, decl := range node.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE || genDecl.Doc == nil {
				continue
			}

			// Check if the declaration has a `// enum:` comment
			for _, comment := range genDecl.Doc.List {
				if commentText := strings.ReplaceAll(strings.TrimSpace(strings.ToLower(comment.Text)), " ", ""); strings.HasPrefix(commentText, "//enum:") {
					enum, err := parseEnumComment(comment.Text, genDecl)
					if err != nil {
						return nil, err
					}
					enums = append(enums, enum)
				}
			}
		}
	}

	return enums, nil
}

// parseEnumComment processes the `// enum:` comment and extracts enum details.
func parseEnumComment(comment string, genDecl *ast.GenDecl) (Enum, error) {
	comment = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(comment)[7:], ":"))

	// Extract the type name and base type
	if len(genDecl.Specs) != 1 {
		return Enum{}, errors.New("unexpected type declaration")
	}

	typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
	if !ok {
		return Enum{}, errors.New("invalid type spec")
	}

	baseType := ""
	switch t := typeSpec.Type.(type) {
	case *ast.Ident:
		baseType = t.Name
	default:
		return Enum{}, errors.New("unsupported base type")
	}

	// Extract enum values from the comment
	values := []Value{}
	iotaCounter := 0
	for _, part := range strings.Split(comment, "|") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		parts := strings.Split(part, "=")
		name := strings.TrimSpace(parts[0])
		val := ""
		if len(parts) > 1 {
			val = strings.TrimSpace(parts[1])
		}
		if baseType == "string" {
			val = fmt.Sprintf("\"%s\"", strings.Trim(strings.TrimSpace(val), "\""))
		} else if strings.Contains(baseType, "int") && val == "" {
			val = fmt.Sprintf("%d", iotaCounter)
			iotaCounter++
		}
		values = append(values, Value{Name: name, Value: val})
	}

	return Enum{
		TypeName: typeSpec.Name.Name,
		BaseType: baseType,
		Values:   values,
	}, nil
}
