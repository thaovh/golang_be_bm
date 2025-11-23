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

	// Serve Swagger UI index
	srv.HandleFunc("/api/docs", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/docs" || r.URL.Path == "/api/docs/" {
			w.Header().Set("Content-Type", "text/html")
			html := createSwaggerUIHTML()
			w.Write([]byte(html))
			return
		}
	})

	// Serve static assets with /api/docs/ prefix using HandlePrefix
	srv.HandlePrefix("/api/docs/", http.StripPrefix("/api/docs", fileServer))
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
  <script src="/api/docs/swagger-ui-bundle.js" charset="UTF-8"></script>
  <script src="/api/docs/swagger-ui-standalone-preset.js" charset="UTF-8"></script>
  <script>
    window.onload = function() {
      if (typeof SwaggerUIBundle === 'undefined') {
        console.error('SwaggerUIBundle is not loaded. Check if swagger-ui-bundle.js is accessible.');
        document.getElementById('swagger-ui').innerHTML = '<div style="padding: 20px; color: red;">Error: Swagger UI JavaScript files failed to load. Please check the console for details.</div>';
        return;
      }
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
