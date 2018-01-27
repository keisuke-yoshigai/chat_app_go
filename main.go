package main

import (
    "net/http"
    "log"
)

func Index(w http.ResponseWriter, req *http.Request) {
    html := `
    <!DOCTYPE html>
    <html>
        <head>
            <title>チャット</title>
        </head>
        <body>
            <h1>チャットを始めよう</h1>
        </body>
    </html>
    `
    w.Write([]byte(html))
}

func main() {
    http.HandleFunc("/", Index)
    if err := http.ListenAndServe(":8090", nil); err != nil {
        log.Fatal(err)
    }
}
