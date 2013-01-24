package main

import (
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

//var (
//	decodeForm = schema.NewDecoder()
//)

type Controller struct {
	toys.Controller
	sess session.Provider
	auth membership.Authenticater
	tmpl *view.View
	db   *DBCtx
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

	//web session
	c.sess = session.NewMgoProvider(w, r, sessColl)

	//web authenthicator
	c.auth = membership.NewAuthDBCtx(w, r, c.sess, userColl, rememberColl)

	//view template
	c.tmpl = h.tmpl

	//database context
	c.db = NewDBCtx(entryColl, catColl, commColl)

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
