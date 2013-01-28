package main

import (
	"github.com/openvn/gocms/dbctx"
	"github.com/openvn/toys"
	"github.com/openvn/toys/secure/membership"
	"github.com/openvn/toys/secure/membership/session"
	"github.com/openvn/toys/view"
	"labix.org/v2/mgo"
	"net/http"
	"path"
)

const (
	dbname string = "test"
)

type Controller struct {
	toys.Controller
	sess session.Provider
	auth membership.Authenticater
	tmpl *view.View
	db   *dbctx.DBCtx
}

func (c *Controller) NewViewData(title string) map[string]interface{} {
	m := make(map[string]interface{})
	m["Title"] = title
	m["MainCats"] = c.db.AllMainCats()
	m["AllTags"] = c.db.AllTags()
	m["DBCtx"] = c.db
	return m
}

func (c *Controller) View(page string, data interface{}) {
	c.tmpl.Load(c, page, data)
}

type Handler struct {
	fn     func(c *Controller)
	dbsess *mgo.Session
	tmpl   *view.View
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := Controller{}
	c.Init(w, r)

	dbsess := h.dbsess.Clone()
	defer dbsess.Close()

	//database collection (table)
	database := dbsess.DB(dbname)
	sessColl := database.C("toysSession")
	userColl := database.C("toysUser")
	rememberColl := database.C("toysUserRemember")
	entryColl := database.C("toysEntry")
	catColl := database.C("toysCat")
	commColl := database.C("toysComm")
	tagColl := database.C("toysTag")
	contColl := database.C("toysCont")
	houseColl := database.C("toysHouse")
	personColl := database.C("toysPerson")

	//web session
	c.sess = session.NewMgoProvider(w, r, sessColl)

	//web authenthicator
	c.auth = membership.NewAuthDBCtx(w, r, c.sess, userColl, rememberColl)

	//database context
	c.db = dbctx.NewDBCtx(catColl, entryColl, commColl, tagColl, contColl,
		houseColl, personColl)

	//view template
	c.tmpl = h.tmpl

	//process
	h.fn(&c)
}

// NewHandler receive a controll function and a mongodb session
func NewHandler(f func(c *Controller), dbsess *mgo.Session, tmpl *view.View) *Handler {
	h := &Handler{}
	h.dbsess = dbsess
	h.tmpl = tmpl
	h.fn = f

	return h
}

// Match is a wrapper function for path.Math
func Match(pattern, name string) bool {
	ok, err := path.Match(pattern, name)
	if err != nil {
		return false
	}
	return ok
}
