// This file was automatically generated based on the contents of *.tmpl
// If you need to update this file, change the contents of those files
// (or add new ones) and run 'go generate'

package main

const (
	useConstantsTmpl = "// This file was automatically generated based on the contents of *.tmpl\n// If you need to update this file, change the contents of those files\n// (or add new ones) and run 'go generate'\n\npackage {{.PackageName}}\n\nconst (\n{{ range .Templates }}\n  {{.ConstantName}} = \"{{.Contents}}\"\n{{ end }}\n)\n"

	useMapFSTmpl = "// This file was automatically generated based on the contents of *.tmpl\n// If you need to update this file, change the contents of those files\n// (or add new ones) and run 'go generate'\n\npackage {{.PackageName}}\n\nimport \"golang.org/x/tools/godoc/vfs/mapfs\"\n\nvar {{.MapFSName}} = mapfs.New(map[string]string{\n{{ range .Templates -}}\n  `{{.SourceFile}}`: \"{{.Contents}}\",\n{{ end }}\n})\n"
)
