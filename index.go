package main

func Index(c *Controller) {
	data := c.NewViewData("Index")
	c.View("index.tmpl", &data)
}
