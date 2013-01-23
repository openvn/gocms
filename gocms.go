// demoweb2 project main.go
package main

import (
	"github.com/openvn/toys/locale"
	"github.com/openvn/toys/view"
	"labix.org/v2/mgo"
	"net/http"
)

func main() {
	//database session
	dbsess, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer dbsess.Close()

	//multi language support
	lang := locale.NewLang("/home/nvcnvn/WorkSpace/test/temp/language")
	lang.Parse("vi")
	lang.Parse("en")

	//view template system
	tmpl := view.NewView("/home/nvcnvn/WorkSpace/test/temp/template")
	tmpl.Resource = "//localhost:8080/static"
	tmpl.SetLang(lang)
	tmpl.Parse("default")

	//routing
	http.Handle("/", NewHandler(IndexRouter, dbsess, tmpl))
	http.ListenAndServe("localhost:8080", nil)
}
