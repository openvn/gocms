package main

import (
	"errors"
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
	}
	//save
	c.CatId = bson.NewObjectId()
	return ctx.catColl.Insert(c)

}

func (ctx *DBCtx) SaveEntry(e *Entry) error {
	if bson.IsObjectIdHex(e.EntryId.Hex()) {
		//update
		return ctx.entryColl.UpdateId(e.EntryId, e)
	}
	//save
	e.EntryId = bson.NewObjectId()
	err := ctx.entryColl.Insert(e)
	if err != nil {
		return err
	}

	return ctx.catColl.UpdateId(e.CatId, bson.M{"$set": bson.M{
		"lastentry": e.EntryId,
	}})
}

func (ctx *DBCtx) SaveComment(c *Comment) error {
	return ctx.commColl.Insert(c)
}

func (ctx *DBCtx) FindPost(id string) (*Entry, error) {
	if bson.IsObjectIdHex(id) {
		oid := bson.ObjectIdHex(id)
		entry := Entry{}
		err := ctx.entryColl.FindId(oid).One(&entry)
		if err != nil {
			return nil, err
		}
		return &entry, nil
	}
	return nil, errors.New("db: invalid id")
}

func (ctx *DBCtx) IncView(catId bson.ObjectId) error {
	return nil
}

func (ctx *DBCtx) AllCats() []Catergory {
	var cats []Catergory
	err := ctx.catColl.Find(nil).All(&cats)
	if err != nil {
		return nil
	}
	return cats
}

func (ctx *DBCtx) AllCatsLastEntry() []struct {
	CatId, CatName, LastEntryId, LastEntryName, LastEntryDescription string
} {
	cats := ctx.AllCats()
	s := make([]struct {
		CatId, CatName, LastEntryId, LastEntryName,
		LastEntryDescription string
	}, 0, len(cats))

	for i := range cats {
		entry := Entry{}
		err := ctx.entryColl.FindId(cats[i].LastEntry).One(&entry)
		if err == nil {
			d := struct {
				CatId, CatName, LastEntryId, LastEntryName, LastEntryDescription string
			}{
				cats[i].CatId.Hex(),
				cats[i].Name,
				entry.EntryId.Hex(),
				entry.Title,
				entry.Description,
			}
			s = append(s, d)
		}
	}
	return s
}

type Catergory struct {
	CatId     bson.ObjectId `bson:"_id"`
	Name      string
	LastEntry bson.ObjectId `bson:",omitempty"`
}

type Entry struct {
	EntryId     bson.ObjectId `bson:"_id"`
	Description string
	Title       string
	Content     string
	At          time.Time
	NumView     int
	Tags        []string
	CatId       bson.ObjectId `bson:",omitempty"`
}

type Comment struct {
	CommId  bson.ObjectId `bson:"_id"`
	Content string
	By      string
	At      time.Time
	EntryId bson.ObjectId `bson:",omitempty"`
}
