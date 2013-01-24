package main

import (
	"time"
)

func Admin(c *Controller) {
	//_, err := c.auth.GetUser()
	//if err != nil {
	//	//not logged user, redirect
	//	c.Redirect("/", 307)
	//	return
	//}
	p := c.Request().URL.Path

	if "/admin/post_submit.html" == p {
		entry := Entry{}
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
			Cats  []Catergory
		}{
			"Post new",
			c.db.AllCats(),
		}
		c.View("post.tmpl", data)
	} else if "/admin/cat_submit.html" == p {
		if len(c.Post("submit", false)) > 0 {
			cat := Catergory{}
			cat.Name = c.Post("name", true)
			err := c.db.SaveCat(&cat)
			if err != nil {
				c.Print(err.Error())
			} else {
				c.Print("cat added")
			}
		}
	} else if "/admin/cat.html" == p {
		println("admin cat")
		c.View("cat.tmpl", nil)
	} else {
		c.Print("???")
	}
}
