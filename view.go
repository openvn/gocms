package main

import (
	"encoding/json"
	"github.com/openvn/gocms/dbctx"
)

func View(c *Controller) {
	uripath := c.Request().URL.Path
	if "/view/post.html" == uripath {
		id := c.Get("id", false)
		entry, err := c.db.FindPost(id)
		if err == nil {
			data := struct {
				Title string
				Entry *dbctx.Entry
			}{
				"View",
				entry,
			}
			c.View("viewpost.tmpl", &data)
		} else {
			c.Print(err.Error())
		}
	} else if "/view/catlst.html" == uripath {
		id := c.Get("id", false)
		b, err := json.Marshal(c.db.CatTree(c.db.DecodeId(id)))
		if err != nil {
			c.Print(err.Error())
			return
		}
		c.Write(b)
	}
}
