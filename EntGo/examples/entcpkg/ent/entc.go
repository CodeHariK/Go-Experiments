// Copyright 2019-present Facebook Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

func main() {
	ex, err := NewExtension("json")
	if err != nil {
		log.Fatalf("creating extension: %v", err)
	}
	// A usage for custom options to configure the code generator to use
	// an extension and inject external dependencies in the generated API.
	opts := []entc.Option{
		entc.Extensions(ex),
		entc.Dependency(
			entc.DependencyType(&http.Client{}),
		),
		entc.Dependency(
			entc.DependencyName("Writer"),
			entc.DependencyTypeInfo(&field.TypeInfo{
				Ident:   "io.Writer",
				PkgPath: "io",
			}),
		),
	}
	err = entc.Generate("./schema", &gen.Config{
		Header: `
			// Copyright 2019-present Facebook Inc. All rights reserved.
			// This source code is licensed under the Apache 2.0 license found
			// in the LICENSE file in the root directory of this source tree.

			// Code generated by ent, DO NOT EDIT.
		`,
		Templates: []*gen.Template{
			// Custom templates can be provided by entc.Extension (see below),
			// or by setting templates on the gen.Config object.
			//
			//	gen.MustParse(gen.NewTemplate("static").
			//		Funcs(template.FuncMap{"title": strings.ToTitle}).
			//		ParseFiles("template/static.tmpl")),
			//
		},
		Hooks: []gen.Hook{
			// Hooks can be provided by entc.Extension (see below),
			// or by setting hooks on the gen.Config object.
			//
			// 	CustomHook1("config 1"),
			// 	CustomHook2("config 2"),
			//
		},
	}, opts...)
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}

// Extension is an example implementation of entc.Extension.
type Extension struct {
	entc.DefaultExtension
	tag       string
	templates []*gen.Template
}

// NewExtension creates a new entc.Extension.
func NewExtension(tag string) (*Extension, error) {
	ex := &Extension{tag: tag}
	t, err := gen.NewTemplate("static").
		Funcs(template.FuncMap{"title": strings.ToTitle}).
		ParseFiles("template/static.tmpl")
	if err != nil {
		return nil, err
	}
	ex.templates = append(ex.templates, t)
	t, err = gen.NewTemplate("debug").
		Funcs(template.FuncMap{"byName": byName}).
		ParseFiles("template/debug.tmpl")
	if err != nil {
		return nil, err
	}
	ex.templates = append(ex.templates, t)
	return ex, nil
}

// Templates of the extension.
func (e *Extension) Templates() []*gen.Template {
	return e.templates
}

// Hooks of the extension.
func (e *Extension) Hooks() []gen.Hook {
	return []gen.Hook{
		TagFields(e.tag),
	}
}

// Annotations of the extension.
func (e *Extension) Annotations() []entc.Annotation {
	return []entc.Annotation{
		Annotation{StructTag: "rql"},
	}
}

// Options provides additional options for the extension.
func (e *Extension) Options() []entc.Option {
	return []entc.Option{
		entc.TemplateFiles("template/stringer.tmpl"),
	}
}

var _ entc.Extension = (*Extension)(nil)

// byName returns a node in the graph by its label/name.
func byName(g *gen.Graph, name string) (*gen.Type, error) {
	for _, n := range g.Nodes {
		if n.Name == name {
			return n, nil
		}
	}
	return nil, fmt.Errorf("node %q was not found in the graph", name)
}

// TagFields tags all fields defined in the schema with the given struct-tag.
func TagFields(name string) gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			for _, node := range g.Nodes {
				for _, field := range node.Fields {
					field.StructTag = fmt.Sprintf("%s:%q", name, field.Name)
				}
			}
			return next.Generate(g)
		})
	}
}

const AnnotationName = "RQL"

// Annotation defines a custom annotation
// to be inject globally to all templates.
type Annotation struct {
	StructTag string
}

func (Annotation) Name() string {
	return AnnotationName
}

var _ schema.Annotation = (*Annotation)(nil)