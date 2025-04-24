package main

import (
    "fmt"
    "net/http"
)

func main() {
    fmt.Printf("Starting {{ .project.name }} v{{ .project.version }}\n")

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Welcome to {{ .project.name }}!")
    })

    {{- range .routes }}
    http.HandleFunc("{{ .path }}", {{ .handler }})
    {{- end }}

    fmt.Println("Server is running on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}

{{- range .routes }}
func {{ .handler }}(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "This is the {{ .handler }} handler")
}
{{- end }}
