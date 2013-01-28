package main

import (
	"github.com/openvn/gocms/dbctx"
	"strconv"
	"time"
)

func Data(c *Controller) {
	p := c.Request().URL.Path
	var err error
	if "/admin/data/submit.html" == p {
		h := dbctx.House{}
		h.Group, err = strconv.Atoi(c.Post("Group", true))
		if err != nil {
			return
		}

		h.Block, err = strconv.Atoi(c.Post("Block", true))
		if err != nil {
			return
		}

		h.Address = c.Post("Address", true)

		m := dbctx.Person{}
		m.FullName = c.Post("FullName", true)
		if len(c.Post("Gender", false)) > 0 {
			m.Gender = true
		}

		m.Birth, err = time.Parse("02/01/2006", c.Post("Birth", false))
		if err != nil {
			return
		}

		m.Quals, err = strconv.Atoi(c.Post("Quals", false))
		if err != nil {
			return
		}

		m.Area = c.Post("Area", true)
		if err != nil {
			return
		}

		if len(c.Post("AttendingSchool", false)) > 0 {
			m.AttendingSchool = true
			m.Class.Title = c.Post("SchoolTitle", true)
			m.Class.School = c.Post("School", true)
		}

		if len(c.Post("Working", false)) > 0 {
			m.Working = true
			m.Work.Title = c.Post("WorkTitle", true)
			m.Work.Office = c.Post("Office", true)
		}

		amount, err := strconv.Atoi(c.Post("Amount", true))
		if err != nil {
			return
		}
		m.Incomes = []dbctx.Income{dbctx.Income{amount, c.Post("Form", true)}}

		m.HI, err = strconv.Atoi(c.Post("HI", true))
		if err != nil {
			return
		}

		m.Desire = c.Post("Desire", true)
		m.Note = c.Post("Note", true)
	} else if "/admin/data/edit.html" == p {
		c.View("houseadd.tmpl", nil)
	} else if "/admin/data/print.html" == p {

	} else {

	}
}
