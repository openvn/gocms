package main

import (
	"time"
)

func Post(c *Controller) {
	_, err := c.auth.GetUser()
	if err != nil {
		//not logged user, redirect
		c.Redirect("/", 403)
		return
	}
	if "/post/submit" == c.Request().URL.Path {
		entry := Entry{}
		entry.Description = c.Post("description", true)
		entry.Title = c.Post("title", true)
		entry.Content = c.Post("content", true)
		entry.At = time.Now()
		entry.NumView = 0

		if c.db.IsValidId(c.Post("catid", false)) {
			entry.CatId = c.db.DecodeId(c.Post("catid", false))
			c.db.SaveEntry(&entry)
		}
	} else {
		c.Print("html form go here")
	}
}
