package main

import (
	"github.com/openvn/gocms/dbctx"
)

func Contact(c *Controller) {
	p := c.Request().URL.Path
	if "/contact/post.html" == p {
		contact := dbctx.Contact{}
		contact.Title = c.Post("title", true)
		contact.From = c.Post("from", true)
		contact.To = c.Post("to", true)
		contact.Content = c.Post("content", true)
		err := c.db.AddContact(&contact)
		if err != nil {
			c.Print(err.Error())
			return
		}
		c.Print("done")
	} else {
		data := c.NewViewData("Contact")
		c.View("contact.tmpl", data)
	}
}
