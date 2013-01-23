package main

// Loggin is the controller for showing loggin form and checking for the data
func Loggin(c *Controller) {
	if "/loggin/submit" == c.Request().URL.Path {
		//checking
		c.auth.ValidateUser(c.Post("email", false), c.Post("password", false))
		//TODO: loggin
	} else {
		c.Print("html form go here")
	}
}
