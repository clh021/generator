package main

import (
    "log"
    "net/http"

    "{{ $.project.package }}/handler"
)

func main() {
    log.Printf("Starting {{ $.project.name }} v{{ $.project.version }}")

    // 注册路由
    {{- range $.routes }}
    http.HandleFunc("{{ .path }}", handler.{{ .handler }})
    {{- end }}

    log.Fatal(http.ListenAndServe(":8080", nil))
}
