package main

func Index(c *Controller) {
	data := struct {
		Title string
		Inf   []struct {
			CatId, CatName, LastEntryId, LastEntryName, LastEntryDescription string
		}
	}{
		"index",
		c.db.AllCatsLastEntry(),
	}

	c.View("index.tmpl", &data)
}
