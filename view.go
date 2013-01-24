package main

func View(c *Controller) {
	uripath := c.Request().URL.Path
	if "/view/post.html" == uripath {
		id := c.Get("id", false)
		entry, err := c.db.FindPost(id)
		if err == nil {
			data := struct {
				Title string
				Entry *Entry
			}{
				"View",
				entry,
			}
			c.View("viewpost.tmpl", data)
		} else {
			c.Print(err.Error())
		}
	} else if "/view/cat.html" == uripath {

	}
}
