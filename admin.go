package main

import (
	"github.com/openvn/gocms/dbctx"
	"time"
)

func Admin(c *Controller) {
	p := c.Request().URL.Path

	if "/admin/post_submit.html" == p {
		entry := dbctx.Entry{}
		entry.Description = c.Post("description", true)
		entry.Title = c.Post("title", true)
		entry.Content = c.Post("content", true)
		entry.At = time.Now()
		entry.NumView = 0
		entry.Tags = []string{
			c.Post("tag0", true),
			c.Post("tag1", true),
			c.Post("tag2", true),
		}

		if c.db.IsValidId(c.Post("catid", false)) {
			entry.CatId = c.db.DecodeId(c.Post("catid", false))
			c.db.SaveEntry(&entry)
		}
	} else if "/admin/post.html" == p {
		data := struct {
			Title string
			Cats  []dbctx.Catergory
		}{
			"Post new",
			c.db.AllCats(),
		}
		c.View("post.tmpl", data)
	} else if "/admin/cat_submit.html" == p {
		if len(c.Post("submit", false)) > 0 {
			cat := dbctx.Catergory{}
			cat.Name = c.Post("name", true)
			if c.db.IsValidId(c.Post("parent", false)) {
				cat.Parent = c.db.DecodeId(c.Post("parent", false))
			}
			err := c.db.SaveCat(&cat)

			if err != nil {
				c.Print(err.Error())
			} else {
				c.Print("cat added")
			}
		}
	} else if "/admin/cat.html" == p {
		data := struct {
			Title    string
			MainCats []dbctx.Catergory
			Cats     []dbctx.Catergory
		}{
			"Add new Cat",
			c.db.AllMainCats(),
			c.db.CatTree(c.db.DecodeId("")),
		}
		c.View("cat.tmpl", &data)
	} else {
		c.Print("???")
	}
}
