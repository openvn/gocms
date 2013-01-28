package main

import (
	"encoding/json"
)

func View(c *Controller) {
	uripath := c.Request().URL.Path
	if "/view/post.html" == uripath {
		id := c.Get("id", false)
		entry, err := c.db.FindPost(id)
		if err == nil {
			data := c.NewViewData(entry.Title)
			data["Entry"] = entry
			data["CatString"], _ = c.db.CatString(entry.CatId)
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
