# {{ .Name }} - {{ .Description }}

{{ partial "_gravl" }}

## Global Flags
|Name|Aliases|Description|
|-|-|-|
{{- range $f := .GlobalFlags }}
|```{{ $f.Name }}```|{{ if names $f }}```{{ names $f }}```{{ end }}|{{ description $f }}|
{{- end }}

## Commands

{{- range .Commands }}
* [{{ fullname . " " }}](#{{ fullname . "-" }})
{{- end }}

{{- range .Commands }}

## *{{ fullname . " " }}*

**Description**

{{ if .Cmd.Description }}{{ .Cmd.Description }}{{ else }}{{ .Cmd.Usage }}{{ end }}

{{ if .Cmd.Action }}
**Syntax**

```sh
$ gravl {{ fullname . " " }} [flags] {{- if .Cmd.ArgsUsage }} {{.Cmd.ArgsUsage}}{{ end }}
```
{{ end }}

{{- if aliases . -}}
{{- end -}}

{{- with .Cmd.Flags }}

**Flags**

|Name|Aliases|Description|
|-|-|-|
{{- range $f := . }}
|```{{ $f.Name }}```|{{ if names $f }}```{{ names $f }}```{{ end }}|{{ description $f }}|
{{- end }}
{{- end }}

{{- $fn := fullname . "-" }}
{{- $x := partial $fn }}
{{ if $x }}
{{ if .Cmd.Action }}**Example**{{ else }}**Overview**{{ end }}

{{ $x }}
{{- end }}
{{- end }}
