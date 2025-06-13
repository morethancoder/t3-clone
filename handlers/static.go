package handlers

import (
    "mime"
    "net/http"
    "path/filepath"
    "strings"
)

func StaticGET() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // get path and check if it is gzipped
        path := r.URL.Path
        gzipped := strings.HasSuffix(path, ".gz")
        if gzipped {
            path = strings.TrimSuffix(path, ".gz")
        }

        // get the extension of file
        ext := filepath.Ext(path)

        // get content type
        contentType := mime.TypeByExtension(ext)

        // default fallbacks
        if contentType == "" {
            switch ext {
            case ".js":
                contentType = "application/javascript"
            case ".css":
                contentType = "text/css"
            case ".html":
                contentType = "text/html"
            default:
                contentType = "application/octet-stream"
            }
        }

        // set content type to let browser know what to expect
        w.Header().Set("Content-Type", contentType)

        // handle file compression
        if gzipped {
            w.Header().Set("Content-Encoding", "gzip")
            w.Header().Set("Vary", "Accept-Encoding")
        }

        // cache control 1 year
        w.Header().Set("Cache-Control", "public, max-age=31536000")

        // strip prefix to allow file server to find file
        http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, r)

    })
}

