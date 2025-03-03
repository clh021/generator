package handler

import (
    "encoding/json"
    "net/http"
)

{{ range $route := .routes }}
// {{ $route.handler }} {{ $route.description }}
func {{ $route.handler }}(w http.ResponseWriter, r *http.Request) {
    if r.Method != "{{ $route.method }}" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // TODO: Implement {{ $route.handler }} logic
    response := map[string]string{
        "message": "{{ $route.description }}",
        "status": "not implemented",
    }

    json.NewEncoder(w).Encode(response)
}
{{ end }}
