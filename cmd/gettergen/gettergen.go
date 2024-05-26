// Command gettergen parses the pokeapi package and generates List* and  Get*
// methods for the pokeapi.Client on types with the pokeapi.Identifier and
// pokeapi.NamedIdentifier directly embedded into them.
//
// Usage:
//
//	gettergen "output-path.go"
//
// Directives (must be written as part of doc comments preceding the identifier
// embed declaration):
//
//	gettergen:plural {{Plural}}
//		When present, this will be used as the pluralized form of the resource's name.
//	gettergen:ignore
//		When present, will ignore the struct when generating getters.
package main

import (
	"context"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"os/signal"
	"regexp"
	"slices"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		_, _ = fmt.Fprintf(os.Stderr, "one args (out file) is required, got args: %v", args)
		os.Exit(1)
	}
	outFile := args[0]

	rds, err := loadResourceDefs(ctx)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to load resource names: %s", err.Error())
		os.Exit(1)
	}

	if err := generateFile(outFile, rds); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to generate file: %s", err.Error())
		os.Exit(1)
	}
}

type directive string

const (
	pluralDirective directive = "plural"
	ignoreDirective directive = "ignore"
)

var directiveRegexp = regexp.MustCompile(`gettergen:(.+)`)

type resourceDefinition struct {
	Name    string
	Plural  string
	Unnamed bool
}

const (
	unnamedIdentifierTypeName = `github.com/nightmarlin/pokeapi.Identifier`
	namedIdentifierTypeName   = `github.com/nightmarlin/pokeapi.NamedIdentifier`
)

func loadResourceDefs(ctx context.Context) ([]resourceDefinition, error) {
	pkgs, err := packages.Load(
		&packages.Config{
			Context: ctx,
			// - type info for embedded struct field checks
			// - syntax & compiled files for directive checks
			Mode: packages.NeedTypes | packages.NeedSyntax | packages.NeedCompiledGoFiles,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("loading package: %w", err)
	} else if len(pkgs) != 1 {
		return nil, fmt.Errorf("only expected to load 1 package directory, got %d", len(pkgs))
	}

	var (
		pkg         = pkgs[0]
		typeScope   = pkg.Types.Scope()
		definitions []resourceDefinition
	)

	for _, objName := range typeScope.Names() {
		var (
			obj = typeScope.Lookup(objName)
			t   = obj.Type()
		)

		if !obj.Exported() {
			continue
		}

		namedT, isNamed := t.(*types.Named)
		if !isNamed {
			continue
		}

		structT, isStruct := namedT.Underlying().(*types.Struct)
		if !isStruct {
			continue
		}

		for fieldNum := range structT.NumFields() {
			f := structT.Field(fieldNum)
			if !f.Embedded() {
				continue
			}

			var rd resourceDefinition

			switch f.Type().String() {
			case unnamedIdentifierTypeName:
				rd = resourceDefinition{Name: objName, Unnamed: true}

			case namedIdentifierTypeName:
				rd = resourceDefinition{Name: objName, Unnamed: false}

			default:
				continue
			}

			d, val := extractDirectiveForDecl(pkg, f.Pos())
			switch d {
			case ignoreDirective:
				continue
			case pluralDirective:
				rd.Plural = val
			}

			definitions = append(definitions, rd)
		}
	}

	return definitions, nil
}

// extractDirectiveForDecl extracts a gettergen directive, if present, from the
// comment group that precedes the declaration at the provided token.Pos
func extractDirectiveForDecl(
	pkg *packages.Package,
	lpos token.Pos,
) (directive, string) {
	var (
		pos               = pkg.Fset.Position(lpos)
		expectCommentLine = pos.Line - 1
		compiledFileIDX   = slices.Index(pkg.CompiledGoFiles, pos.Filename)

		commentGroupIDX = slices.IndexFunc(
			pkg.Syntax[compiledFileIDX].Comments,
			func(cg *ast.CommentGroup) bool {
				for _, cLine := range cg.List {
					if pkg.Fset.Position(cLine.Pos()).Line == expectCommentLine {
						return true
					}
				}
				return false
			},
		)
	)

	if commentGroupIDX == -1 {
		return "", ""
	}

	for _, cLine := range pkg.Syntax[compiledFileIDX].Comments[commentGroupIDX].List {
		match := directiveRegexp.FindStringSubmatch(cLine.Text)
		if len(match) == 2 {
			d, s, _ := strings.Cut(match[1], " ")
			return directive(d), s
		}
	}
	return "", ""
}

// generateFile creates a file, writes the prelude to it, and then executes the
// resourceTemplate repeatedly with the provided resourceDefinition slice.
func generateFile(path string, rds []resourceDefinition) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer func() { _ = f.Close() }()

	if _, err := f.WriteString(prelude); err != nil {
		return fmt.Errorf("writing prelude: %w", err)
	}

	for _, rd := range rds {
		tArg := rd.toTemplateArgs()
		if err := resourceTemplate.Execute(f, tArg); err != nil {
			return fmt.Errorf("executing template for %q: %w", tArg.Name, err)
		}
	}

	return nil
}

var skewerRegexp = regexp.MustCompile(`([a-z])([A-Z])`)

func (rd resourceDefinition) toTemplateArgs() resourceTemplateArgs {
	plural := fmt.Sprintf("%ss", rd.Name)
	if rd.Plural != "" {
		plural = rd.Plural
	}

	identName := "ident"
	pageTypePrefix := "Named"
	if rd.Unnamed {
		identName = "id"
		pageTypePrefix = ""
	}

	return resourceTemplateArgs{
		Name:          rd.Name,
		Plural:        plural,
		IsUnnamed:     rd.Unnamed,
		IdentName:     identName,
		KebabCase:     strings.ToLower(skewerRegexp.ReplaceAllString(rd.Name, `$1-$2`)),
		ReferenceType: fmt.Sprintf("[%sAPIResource[%s], %[2]s]", pageTypePrefix, rd.Name),
	}
}

var resourceTemplate = template.Must(template.New("resource").Parse(resourceTemplateString))

type resourceTemplateArgs struct {
	Name          string
	Plural        string
	IsUnnamed     bool
	IdentName     string
	KebabCase     string
	ReferenceType string
}

const prelude = `// Code generated by github.com/nightmarlin/pokeapi/cmd/gettergen@v0; DO NOT EDIT.

package pokeapi

import "context"
`

const resourceTemplateString = `
const {{ .Name }}Resource ResourceName{{ .ReferenceType }} = "{{ .KebabCase }}"
{{ if .IsUnnamed }}
// Get{{ .Name }} only accepts the ID of the desired {{ .Name }}.{{ end }}
func (c *Client) Get{{ .Name }}(ctx context.Context, {{ .IdentName }} string) (*{{ .Name }}, error) {
	return {{ .Name }}Resource.Get(ctx, c, {{ .IdentName }})
}
func (c *Client) List{{ .Plural }}(ctx context.Context, opts *ListOpts) (*Page{{ .ReferenceType }}, error) {
	return {{ .Name }}Resource.List(ctx, c, opts)
}
`
