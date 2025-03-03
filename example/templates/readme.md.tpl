# {{ .project.name }}

{{ .project.description }}

Version: {{ .project.version }}
Author: {{ .project.author }}

## API Endpoints

{{- range $route := .routes }}
### {{ $route.method }} {{ $route.path }}
{{ $route.description }}
{{- end }}

## Getting Started

```bash
go run main.go
```

The server will start on :8080
