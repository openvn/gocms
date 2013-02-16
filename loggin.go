package gocms

// Loggin is the controller for showing loggin form and checking for the data
func Loggin(c *Controller) {
	uripath := c.Request().URL.Path
	if "/user/loggin.html" == uripath {
		remember := 0
		if c.Post("remember", false) == "1" {
			remember = 43200
		}

		u, err := c.auth.ValidateUser(c.Post("email", false), c.Post("password", false))
		if err != nil {
			c.Print("worng", err.Error())
			return
		}

		err = c.auth.LogginUser(u.Id.Hex(), remember)
		if err != nil {
			c.Print("err loggin", err.Error())
		}
	} else if "/user/adduser.html" == uripath {
		var name string
		u, err := c.auth.GetUser()
		if err != nil {
			name = "not logged" + err.Error()
		} else {
			name = u.Email
		}

		data := c.NewViewData("Add new user")
		data["Email"] = name

		c.View("adduser.tmpl", &data)
	} else if "/user/add.html" == uripath {
		_, err := c.auth.GetUser()
		if err != nil {
			c.Print("Loogin required!")
			return
		}
		if len(c.Post("submit", false)) > 0 {
			err := c.auth.AddUser(
				c.Post("email", true),
				c.Post("password", false),
				false,
				true,
			)
			if err != nil {
				println(err.Error())
			}
		}
		c.Redirect("/user/", 307)
	} else {
		data := c.NewViewData("Loggin")
		c.View("loggin.tmpl", data)
	}
}
