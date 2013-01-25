package main

// Loggin is the controller for showing loggin form and checking for the data
func Loggin(c *Controller) {
	uripath := c.Request().URL.Path
	if "/user/loggin" == uripath {
		u, err := c.auth.ValidateUser(c.Post("email", false), c.Post("password", false))
		if err != nil {
			c.Print("worng", err.Error())
			return
		}

		err = c.auth.LogginUser(u.Id.Hex(), 0)
		if err != nil {
			c.Print("err loggin", err.Error())
		}
	} else if "/user/add" == uripath {
		var name string
		u, err := c.auth.GetUser()
		if err != nil {
			name = "not logged" + err.Error()
		} else {
			name = u.Email
		}
		data := struct {
			Title string
			Email string
		}{
			"loggin",
			name,
		}
		c.View("adduser.tmpl", &data)
	} else if "/user/add2" == uripath {
		if len(c.Post("submit", false)) > 0 {
			c.auth.AddUser(
				c.Post("email", true),
				c.Post("password", false),
				false,
				true,
			)
		}
		c.Redirect("/user/", 307)
	} else {
		c.View("loggin.tmpl", nil)
	}
}
