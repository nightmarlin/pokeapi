// Command gettergen parses the pokeapi package and generates List* and  Get*
// methods for the pokeapi.Client on types with the pokeapi.Identifier and
// pokeapi.NamedIdentifier directly embedded into them.
//
// Usage:
//
//	gettergen path/to/file/in/package.go
package main

import (
	"context"
	"flag"
	"fmt"
	"go/types"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

func printErr(format string, args ...any) { _, _ = fmt.Fprintf(os.Stderr, format, args...) }

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		printErr("two args (package path, out file) are required, got args: %v\n", args)
		os.Exit(1)
	}

	packageDir := filepath.Dir(args[0])
	outFile := args[1]

	resourceDefinitions, err := loadResourceNames(ctx, packageDir)
	if err != nil {
		printErr("failed to load resource names: %s", err.Error())
		os.Exit(1)
	}

	tArgs := make([]resourceTemplateArgs, len(resourceDefinitions))
	for i, d := range resourceDefinitions {
		tArgs[i] = d.toTemplateArgs()
	}

	if err := generateFile(outFile, tArgs); err != nil {
		printErr("failed to generate file: %s", err.Error())
		os.Exit(1)
	}
}

type resourceDefinition struct {
	resourceName      string
	isUnnamedResource bool
}

const (
	unnamedIdentifierTypeName = `Identifier`
	namedIdentifierTypeName   = `NamedIdentifier`
)

var embedExclusions = map[string]struct{}{
	`NamedIdentifier`: {}, // Embeds Identifier, but is not a Resource
}

func loadResourceNames(ctx context.Context, dir string) ([]resourceDefinition, error) {
	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedTypes, Context: ctx}, dir)
	if err != nil {
		return nil, fmt.Errorf("loading package: %w", err)
	} else if len(pkgs) != 1 {
		return nil, fmt.Errorf("only expected to load 1 package directory, got %d", len(pkgs))
	}

	ts := pkgs[0].Types.Scope()
	var definitions []resourceDefinition

	for _, objName := range ts.Names() {
		if _, excluded := embedExclusions[objName]; excluded {
			continue
		}
		t := ts.Lookup(objName).Type()

		namedT, isNamed := t.(*types.Named)
		if !isNamed {
			continue
		}

		structT, isStruct := namedT.Underlying().(*types.Struct)
		if !isStruct {
			continue
		}

		for fieldNum := range structT.NumFields() {
			structF := structT.Field(fieldNum)

			if !structF.Embedded() {
				continue
			}

			switch structF.Name() {
			case unnamedIdentifierTypeName:
				definitions = append(
					definitions,
					resourceDefinition{resourceName: objName, isUnnamedResource: true},
				)

			case namedIdentifierTypeName:
				definitions = append(
					definitions,
					resourceDefinition{resourceName: objName, isUnnamedResource: false},
				)
			}
		}
	}

	return definitions, nil
}

func generateFile(path string, tArgs []resourceTemplateArgs) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer func() { _ = f.Close() }()

	if _, err := f.WriteString(prelude); err != nil {
		return fmt.Errorf("writing prelude: %w", err)
	}

	for _, arg := range tArgs {
		if err := resourceTemplate.Execute(f, arg); err != nil {
			return fmt.Errorf("executing template for %q: %w", arg.Name, err)
		}
	}

	return nil
}

// pluralise pluralises strings with the following rules:
// 1. if it ends in -y, it's converted to -ies
// 2. if it ends in -s, it's converted to -ses
// 3. else append 's'
func pluralise(name string) string {
	if name == "" {
		return ""
	} else if s, cut := strings.CutSuffix(name, "y"); cut {
		return fmt.Sprintf("%sies", s)
	} else if strings.HasSuffix(name, "s") {
		return fmt.Sprintf("%ses", name)
	}
	return fmt.Sprintf("%ss", name)
}

// pageType generates the correct type for the page returned by List* methods.
func pageType(name string, isUnnamed bool) string {
	n := "Named"
	if isUnnamed {
		n = ""
	}
	return fmt.Sprintf("[%sAPIResource[%s], %[2]s]", n, name)
}

// identName generates the identifier param name in Get* methods. For unnamed
// resources that only accept IDs, it returns "id" - otherwise, "ident".
func identName(isUnnamed bool) string {
	if isUnnamed {
		return "id"
	}
	return "ident"
}

var skewerRegexp = regexp.MustCompile(`([a-z])([A-Z])`)

// skewer converts a PascalCased string to a kebab-cased one.
func skewer(name string) string {
	return strings.ToLower(skewerRegexp.ReplaceAllString(name, `$1-$2`))
}

var resourceTemplate = template.Must(template.New("resource").Parse(resourceTemplateString))

func (rd resourceDefinition) toTemplateArgs() resourceTemplateArgs {
	return resourceTemplateArgs{
		Name:          rd.resourceName,
		Plural:        pluralise(rd.resourceName),
		IsUnnamed:     rd.isUnnamedResource,
		IdentName:     identName(rd.isUnnamedResource),
		KebabCase:     skewer(rd.resourceName),
		ReferenceType: pageType(rd.resourceName, rd.isUnnamedResource),
	}
}

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
