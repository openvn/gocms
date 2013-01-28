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
	lang := locale.NewLang("language")
	lang.Parse("vi")
	lang.Parse("en")

	//view template system
	tmpl := view.NewView("template")
	tmpl.Resource = "//localhost:8080/statics"
	tmpl.SetLang(lang)
	tmpl.AddFunc("catlst", func() string { return "bbbbb" })
	tmpl.Parse("default")

	//routing
	http.Handle("/", NewHandler(Router, dbsess, tmpl))
	http.Handle("/statics/", http.StripPrefix("/statics/", http.FileServer(http.Dir("statics"))))
	http.ListenAndServe("localhost:8080", nil)
}

func Router(c *Controller) {
	p := c.Request().URL.Path

	if Match("/admin/data/*", p) {
		Data(c)
	} else if Match("/contact/*", p) {
		Contact(c)
	} else if Match("/admin/*", p) {
		Admin(c)
	} else if Match("/user/*", p) {
		Loggin(c)
	} else if Match("/view/*", p) {
		View(c)
	} else {
		Index(c)
	}
}
