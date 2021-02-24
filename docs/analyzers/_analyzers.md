{{- range .Analyzers }}

## *{{ .Name }}*

**Description**

{{ .Doc }}

{{- with .Flags }}

**Flags**

|Flag|Default|Description|
|-|-|-|
{{- range flags . }}
|```{{.Name}}```|```{{.DefValue}}```|{{.Usage}}|
{{- end }}
{{- end }}
{{- end }}
