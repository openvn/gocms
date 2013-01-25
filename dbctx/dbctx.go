package dbctx

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

// NewDBCtx receive 3 *mgo.Collection for catergorys, entrys and comments
func NewDBCtx(cat, entry, comm *mgo.Collection) *DBCtx {
	ctx := &DBCtx{}
	ctx.entryColl = entry
	ctx.catColl = cat
	ctx.commColl = comm

	return ctx
}

// IsValidId returns true if hex can convert to bson.ObjectId
func (ctx *DBCtx) IsValidId(hex string) bool {
	return bson.IsObjectIdHex(hex)
}

// DecodeId returns a valid bson.ObjectId if the given hex is valid. If not it
// will return an invalid bson.ObjectId
func (ctx *DBCtx) DecodeId(hex string) bson.ObjectId {
	if bson.IsObjectIdHex(hex) {
		return bson.ObjectIdHex(hex)
	}
	var invalid bson.ObjectId
	return invalid
}

func (ctx *DBCtx) SaveCat(c *Catergory) error {
	if c.CatId.Valid() {
		//update
		return ctx.catColl.UpdateId(c.CatId, c)
	}
	//save
	c.CatId = bson.NewObjectId()
	if c.Parent.Valid() {
		acts := Catergory{}
		err := ctx.catColl.FindId(c.Parent).Select(bson.M{
			"path":      1,
			"ancestors": 1,
		}).One(&acts)
		if err == nil {
			c.Ancestors = acts.Ancestors
		}
		c.Ancestors = append(c.Ancestors, c.Parent)
	}
	return ctx.catColl.Insert(c)

}

func (ctx *DBCtx) SaveEntry(e *Entry) error {
	if e.EntryId.Valid() {
		//update
		return ctx.entryColl.UpdateId(e.EntryId, e)
	}
	//save
	e.EntryId = bson.NewObjectId()
	//update lastentry for parent
	err := ctx.catColl.UpdateId(e.CatId, bson.M{"$set": bson.M{
		"lastentry": e.EntryId,
	}})
	if err != nil {
		return err
	}
	//update lastentry for parent ancestors
	parent := Catergory{}
	err = ctx.catColl.FindId(e.CatId).Select(bson.M{"ancestors": 1}).One(&parent)
	if err == nil {
		if len(parent.Ancestors) > 0 {
			ctx.catColl.Update(
				bson.M{"ancestors": bson.M{"$in": parent}},
				bson.M{"$set": bson.M{"lastentry": e.EntryId}},
			)
		}
	}

	return ctx.entryColl.Insert(e)
}

func (ctx *DBCtx) SaveComment(c *Comment) error {
	c.At = time.Now()
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

func (ctx *DBCtx) AllMainCats() []Catergory {
	var cats []Catergory
	err := ctx.catColl.Find(bson.M{"path": bson.M{"$exists": false}}).All(&cats)
	if err != nil {
		return nil
	}
	return cats
}

func (ctx *DBCtx) CatTree(root bson.ObjectId) []Catergory {
	all := []Catergory{}
	var query *mgo.Query

	if root.Valid() {
		query = ctx.catColl.Find(bson.M{"parent": root})
	} else {
		query = ctx.catColl.Find(bson.M{"parent": bson.M{"$exists": false}})
	}
	query.All(&all)

	return all
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
