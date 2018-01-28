package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	temp1    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t.once.Do(func() { //once.Do(f)は初めての呼び出しのみ実行されるメソッド
		t.temp1 =
			template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	if err := t.temp1.Execute(w, nil); err != nil { //temlateオブジェクトをio.Writerに書き込む
		log.Fatal(err)
	}
}

func main() {
	ptTmpHandler := templateHandler{filename: "chat.html"}
	http.Handle("/", &ptTmpHandler)
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}
