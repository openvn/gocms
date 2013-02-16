package gocms

func Router(c *Controller) {
	p := c.Request().URL.Path
	if Match("/contact/*", p) {
		Contact(c)
	} else if Match("/admin/*", p) {
		Admin(c)
	} else if Match("/user/*", p) {
		Loggin(c)
	} else if Match("/view/*", p) {
		View(c)
	} else {
		Index(c)
	}
}
