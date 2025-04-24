# {{ .project.name }}

{{ .project.description }}

## Features

{{ range .project.features }}
- {{ . }}
{{ end }}

## API Endpoints

{{ range .routes }}
### {{ .path }}
- Handler: `{{ .handler }}`
- Description: {{ .description }}
{{ end }}
## Getting Started

### Prerequisites

- Go {{ .project.goVersion }}
- PostgreSQL

### Installation

1. Clone the repository
2. Install dependencies
3. Set up the database
4. Run the server

## Configuration

### Database

- Driver: {{ .database.driver }}
- Host: {{ .database.host }}
- Port: {{ .database.port }}
- Name: {{ .database.name }}

### Server

- Port: {{ .server.port }}
- Host: {{ .server.host }}

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.