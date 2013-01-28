package main

import (
	"github.com/openvn/gocms/dbctx"
)

func Contact(c *Controller) {
	p := c.Request().URL.Path
	if "/contact/post.html" == p {
		contact := dbctx.Contact{}
		c.View("", contact)
	} else {
		data := c.NewViewData("Contact")
		c.View("contact.tmpl", data)
	}
}
