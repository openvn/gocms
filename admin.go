package main

import (
	"encoding/json"
	"github.com/openvn/gocms/dbctx"
	"strings"
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

		tags := strings.Split(c.Post("tags", true), ",")
		entry.Tags = make([]string, 0, len(tags))
		for i := range tags {
			if len(tags[i]) > 0 {
				entry.Tags = append(entry.Tags, tags[i])
			}
		}

		if c.db.IsValidId(c.Post("catid", false)) {
			entry.CatId = c.db.DecodeId(c.Post("catid", false))
			err := c.db.SaveEntry(&entry)
			if err == nil {
				c.db.AddTags(entry.Tags...)
			}
		}
	} else if "/admin/post.html" == p {
		data := c.NewViewData("New Post")
		data["Cats"] = c.db.AllCats()
		c.View("post.tmpl", data)
	} else if "/admin/cat_submit.html" == p {
		if len(c.Post("name", false)) > 0 {
			cat := dbctx.Catergory{}
			cat.Name = c.Post("name", true)
			if c.db.IsValidId(c.Post("parent", false)) {
				cat.Parent = c.db.DecodeId(c.Post("parent", false))
			}
			err := c.db.SaveCat(&cat)

			if err != nil {
				c.Print(err.Error())
			} else {
				b, err := json.Marshal(&[1]dbctx.Catergory{cat})
				if err == nil {
					c.Write(b)
				} else {
					println(err.Error())
				}
			}
		}
	} else if "/admin/cat.html" == p {
		data := c.NewViewData("New Caterogry")
		data["Cats"] = c.db.CatTree(c.db.DecodeId(""))
		c.View("cat.tmpl", &data)
	} else if "/admin/theme.html" == p {
		c.tmpl.Parse("default")
	} else {
		c.Print("???")
	}
}
