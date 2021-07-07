package main

import (
	"flag"
	"io"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
)

const (
	outFileName = "generated.js"

	twirpUtil = `function createRequest(url, body) {
	return new Request(url, {
		method: "POST",
		credentials: "same-origin",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(body),
	});
}
`
	methodTempl = `
{{.Comments.Leading}}export async function {{.JSName}}({{range $i, $v := .Input.Fields}}{{if $i}}, {{end}}{{$v.Desc.JSONName}}{{end}}) {
	const res = await fetch(createRequest("{{.PathPrefix}}/{{.Desc.ParentFile.Package}}.{{.Parent.GoName}}/{{.GoName}}", { {{range $i, $v := .Input.Fields}}{{if $i}}, {{end}}"{{$v.Desc.JSONName}}": {{$v.Desc.JSONName}}{{end}} }));
	const jsonBody = await res.json();
	if (res.ok) {
		return jsonBody;
	}
	throw new Error(jsonBody.msg);
}
`
)

func main() {
	// Set up our flags. The only one we care about for now is the server path prefix.
	var flags flag.FlagSet
	prefix := flags.String("pathPrefix", "/twirp", "the server path prefix to use, if modified from the Twirp default")

	// No special options for this generator
	opts := protogen.Options{ParamFunc: flags.Set}
	opts.Run(func(plugin *protogen.Plugin) error {
		gen, err := newGenerator(*prefix)
		if err != nil {
			return err
		}

		out := plugin.NewGeneratedFile(outFileName, "")
		out.Write([]byte(twirpUtil))

		for _, file := range plugin.Files {
			for _, svc := range file.Services {
				for _, method := range svc.Methods {
					if err := gen.writeMethod(out, method); err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
}

type jsMethod struct {
	*protogen.Method
	PathPrefix string
}

// JSName exists as a way to get our camelCase method name.
func (j jsMethod) JSName() string {
	if j.GoName == "" {
		return ""
	}
	return strings.ToLower(j.GoName[:1]) + j.GoName[1:]
}

func newGenerator(prefix string) (*generator, error) {
	methodTemplate, err := template.New("func").Parse(methodTempl)
	if err != nil {
		return nil, err
	}

	return &generator{
		Prefix:         prefix,
		MethodTemplate: methodTemplate,
	}, nil
}

type generator struct {
	Prefix         string
	MethodTemplate *template.Template
}

func (g generator) writeMethod(w io.Writer, method *protogen.Method) error {
	in := jsMethod{
		Method:     method,
		PathPrefix: g.Prefix,
	}
	return g.MethodTemplate.Execute(w, in)
}
