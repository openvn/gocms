package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type DBCtx struct {
	entryColl *mgo.Collection
	catColl   *mgo.Collection
	commColl  *mgo.Collection
}

func NewDBCtx(entry, cat, comm *mgo.Collection) *DBCtx {
	ctx := &DBCtx{}
	ctx.entryColl = entry
	ctx.catColl = cat
	ctx.commColl = comm

	return ctx
}

func (ctx *DBCtx) IsValidId(hex string) bool {
	return bson.IsObjectIdHex(hex)
}

func (ctx *DBCtx) DecodeId(hex string) bson.ObjectId {
	return bson.ObjectIdHex(hex)
}

func (ctx *DBCtx) SaveCat(c *Catergory) error {
	if bson.IsObjectIdHex(c.CatId.Hex()) {
		//update
		return ctx.catColl.UpdateId(c.CatId, c)
	} else {
		//save
		c.CatId = bson.NewObjectId()
		return ctx.catColl.Insert(c)
	}
	return nil
}

func (ctx *DBCtx) SaveEntry(e *Entry) error {
	return nil
}

func (ctx *DBCtx) SaveComment(c *Comment) error {
	return nil
}

func (ctx *DBCtx) IncView(catId bson.ObjectId) error {
	return nil
}

type Catergory struct {
	CatId     bson.ObjectId `bson:"_id"`
	Name      string
	LastEntry bson.ObjectId `bson:"_id"`
}

type Entry struct {
	EntryId     bson.ObjectId `bson:"_id"`
	Description string
	Title       string
	Content     string
	At          time.Time
	NumView     int
	Tags        []string
	CatId       bson.ObjectId `bson:"_id"`
}

type Comment struct {
	CommId  bson.ObjectId `bson:"_id"`
	Content string
	By      string
	At      time.Time
	EntryId bson.ObjectId `bson:"_id"`
}
