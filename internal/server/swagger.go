package server

import (
	"embed"
	"io/fs"
	"net/http"
	"os"

	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

//go:embed swagger-ui/*
var swaggerUIFiles embed.FS

// RegisterSwaggerUI registers Swagger UI endpoints
func RegisterSwaggerUI(srv *kratoshttp.Server) {
	// Serve OpenAPI spec file
	srv.HandleFunc("/api/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		// Read openapi.yaml from root directory
		data, err := os.ReadFile("openapi.yaml")
		if err != nil {
			http.Error(w, "Failed to read OpenAPI spec", http.StatusInternalServerError)
			return
		}
		w.Write(data)
	})

	// Serve Swagger UI static files
	swaggerUI, err := fs.Sub(swaggerUIFiles, "swagger-ui")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(swaggerUI))

	// Serve Swagger UI index with custom OpenAPI spec URL
	srv.HandleFunc("/api/docs", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/docs" || r.URL.Path == "/api/docs/" {
			w.Header().Set("Content-Type", "text/html")
			html := createSwaggerUIHTML()
			w.Write([]byte(html))
			return
		}
		
		// Serve static files
		// Remove /api/docs prefix
		r.URL.Path = r.URL.Path[len("/api/docs"):]
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})

	// Handle static assets with /api/docs/ prefix
	srv.HandleFunc("/api/docs/", func(w http.ResponseWriter, r *http.Request) {
		// Remove /api/docs prefix
		path := r.URL.Path[len("/api/docs"):]
		if path == "" {
			path = "/"
		}
		r.URL.Path = path
		fileServer.ServeHTTP(w, r)
	})
}

// createSwaggerUIHTML creates Swagger UI HTML with our OpenAPI spec URL
func createSwaggerUIHTML() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>API Documentation - Backend Service</title>
  <link rel="stylesheet" type="text/css" href="/api/docs/swagger-ui.css" />
  <link rel="icon" type="image/png" href="/api/docs/favicon-32x32.png" sizes="32x32" />
  <link rel="icon" type="image/png" href="/api/docs/favicon-16x16.png" sizes="16x16" />
  <style>
    html {
      box-sizing: border-box;
      overflow: -moz-scrollbars-vertical;
      overflow-y: scroll;
    }
    *, *:before, *:after {
      box-sizing: inherit;
    }
    body {
      margin:0;
      background: #fafafa;
    }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="/api/docs/swagger-ui-bundle.js"></script>
  <script src="/api/docs/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function() {
      window.ui = SwaggerUIBundle({
        url: "/api/openapi.yaml",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout",
        validatorUrl: null
      });
    };
  </script>
</body>
</html>`
}
