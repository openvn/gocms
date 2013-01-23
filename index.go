package main

func IndexRouter(c *Controller) {
	p := c.Request().URL.Path
	//BUG
	if "/" == p {
		Index(c)
	} else if "/view/" == p {

	} else if "/post/" == p {
		Post(c)
	} else if "/loggin/" == p {
		Loggin(c)
	} else {
		Index(c)
	}
}

func Index(c *Controller) {
	c.Print("index here")
}
