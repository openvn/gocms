package gocms

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
	} else if "/view/find.html" == uripath {
		tag := c.Get("tag", false)
		println(tag)
		data := c.NewViewData("Post with tag: " + tag)
		if len(tag) > 0 {
			data["Entrys"] = c.db.EntryByTag(tag, c.db.DecodeId(""), 10)
		}
		c.View("entrytag.tmpl", data)
	}
}
