package languages

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

var _ LanguageProcessor = (*GoProcessor)(nil)

// GoProcessor handles post-processing of Go FlatBuffer files.
type GoProcessor struct {
	flatbuffers []string
}

// NewGoProcessor creates a new Go processor.
func NewGoProcessor() *GoProcessor {
	return &GoProcessor{
		flatbuffers: []string{},
	}
}

// ProcessFile adds encryption to a Go FlatBuffer file.
func (p *GoProcessor) ProcessFile(filePath string) error {
	// Parse the Go file using go parser, go/parser is used to parse the file and return the AST.
	tree, err := parser.ParseFile(token.NewFileSet(), filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}
	// Check the imports for flatbuffers
	if !p.isFlatBufferFile(tree) {
		return ErrFlatBuffersNotImported
	}
	p.flatbuffers = append(p.flatbuffers, strings.TrimSuffix(filepath.Base(filePath), p.Extension()))

	return nil
}

// Extension returns the file extension for the language.
func (*GoProcessor) Extension() string {
	return ".go"
}

const (
	FlatDataHelperFileName = "flatdatas_helper.go"
)

func (*GoProcessor) PreProcess(context.Context, string) error {
	return nil
}

func (p *GoProcessor) PostProcess(_ context.Context, outputDir string) error {
	// Create a new file: flatdatas_helper.go
	// Create a global variable: fbs and set it to the flatbuffers package
	// Create a function: GetFlatDataByName(name string)
	f, osErr := os.Create(filepath.Join(outputDir, FlatDataHelperFileName))
	if osErr != nil {
		return fmt.Errorf("failed to create file: %w", osErr)
	}
	defer f.Close()

	// Write the file
	tmpl, err := template.New("flatbufferCode").Parse(flatbufferCode)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	// Execute the template. Output the result to the file
	sort.Strings(p.flatbuffers)
	if executeErr := tmpl.Execute(f, p.flatbuffers); executeErr != nil {
		return fmt.Errorf("failed to execute template: %w", executeErr)
	}
	return nil
}

func (*GoProcessor) isFlatBufferFile(file *ast.File) bool {
	// First check if the file imports the flatbuffers package
	if !hasFlatBuffersImport(file) {
		return false
	}
	return usesFlatBuffersTable(file)
}

// hasFlatBuffersImport checks if the file imports the flatbuffers package.
func hasFlatBuffersImport(file *ast.File) bool {
	for _, imp := range file.Imports {
		packagePath := strings.Trim(imp.Path.Value, `"`)
		if packagePath == "github.com/google/flatbuffers/go" {
			return true
		}
	}
	return false
}

// usesFlatBuffersTable checks if the file uses the flatbuffers.Table type.
func usesFlatBuffersTable(file *ast.File) bool {
	result := false
	ast.Inspect(file, func(n ast.Node) bool {
		field, ok := n.(*ast.Field)
		if !ok {
			return true
		}

		sel, ok := field.Type.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		x, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}

		if x.Name == "flatbuffers" && sel.Sel.Name == "Table" {
			result = true
			return false
		}

		return true
	})
	return result
}

const flatbufferCode = `package flatdata

import (
	"reflect"
)

var fbs = map[string]reflect.Type{
{{- range . }}
	"{{ . }}": reflect.TypeOf((*{{ . }})(nil)).Elem(),
{{- end }}
}

func GetFlatDataByName(name string) any {
	if data, ok := fbs[name]; ok {
		return reflect.New(data).Interface()
	}
	return nil
}
`
