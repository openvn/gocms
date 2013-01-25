package main

import (
	"github.com/openvn/gocms/dbctx"
)

func Index(c *Controller) {
	data := struct {
		Title string
		Inf   []struct {
			CatId, CatName, LastEntryId, LastEntryName, LastEntryDescription string
		}
		MainCats []dbctx.Catergory
	}{
		"index",
		c.db.AllCatsLastEntry(),
		c.db.AllMainCats(),
	}

	c.View("index.tmpl", &data)
}
