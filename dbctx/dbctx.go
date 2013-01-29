package dbctx

import (
	"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type DBCtx struct {
	entryColl  *mgo.Collection
	catColl    *mgo.Collection
	commColl   *mgo.Collection
	tagColl    *mgo.Collection
	contColl   *mgo.Collection
	houseColl  *mgo.Collection
	personColl *mgo.Collection
}

// NewDBCtx receive 3 *mgo.Collection for catergorys, entrys and comments
func NewDBCtx(cat, entry, comm, tag, cont, house, person *mgo.Collection) *DBCtx {
	ctx := &DBCtx{}
	ctx.entryColl = entry
	ctx.catColl = cat
	ctx.commColl = comm
	ctx.tagColl = tag
	ctx.contColl = cont
	ctx.houseColl = house
	ctx.personColl = person
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
	//save
	c.CatId = bson.NewObjectId()
	if c.Parent.Valid() {
		acts := Catergory{}
		err := ctx.catColl.FindId(c.Parent).Select(bson.M{"ancestors": 1}).One(&acts)
		if err == nil {
			c.Ancestors = acts.Ancestors
		}
		c.Ancestors = append(c.Ancestors, c.Parent)
		ctx.catColl.UpdateId(c.Parent, bson.M{"$push": bson.M{"children": c.CatId}})
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
				bson.M{"_id": bson.M{"$in": parent.Ancestors}},
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
	err := ctx.catColl.Find(bson.M{"parent": bson.M{"$exists": false}}).Sort("_id").All(&cats)
	if err != nil {
		return nil
	}
	return cats
}

func (ctx *DBCtx) AddTags(tags ...string) int {
	n := 0
	for i := range tags {
		_, err := ctx.tagColl.UpsertId(tags[i], bson.M{"$inc": bson.M{"usage": 1}})
		if err == nil {
			n++
		}
	}
	return n
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

func (ctx *DBCtx) AllTags() []Tag {
	tags := []Tag{}
	ctx.tagColl.Find(nil).Sort("usage").All(&tags)
	return tags
}

func (ctx *DBCtx) CatSummary(parent bson.ObjectId) ([]CatergorySummary, error) {
	cats := []Catergory{}
	err := ctx.catColl.Find(bson.M{"parent": parent}).Select(bson.M{
		"_id":       true,
		"name":      true,
		"lastentry": true,
	}).All(&cats)
	if err != nil {
		return nil, err
	}

	n := len(cats)
	summary := make([]CatergorySummary, n, n)
	for i := range cats {
		summary[i].CatId = cats[i].CatId
		summary[i].CatName = cats[i].Name
		if cats[i].LastEntry.Valid() {
			entry := Entry{}
			err = ctx.entryColl.FindId(cats[i].LastEntry).Select(bson.M{
				"title":       true,
				"description": true,
			}).One(&entry)
			if err == nil {
				summary[i].LastEntryId = cats[i].LastEntry
				summary[i].LastEntryName = entry.Title
				summary[i].LastEntryDescription = entry.Description
			}
		}

	}
	return summary, nil
}

func (ctx *DBCtx) IdRoot(id bson.ObjectId) bson.ObjectId {
	cat := Catergory{}
	ctx.catColl.FindId(id).Select(bson.M{"ancestors": true}).One(&cat)
	if len(cat.Ancestors) > 0 {
		return cat.Ancestors[0]
	}
	return cat.CatId
}

func (ctx *DBCtx) CatString(id bson.ObjectId) ([]Catergory, error) {
	cat := Catergory{}
	err := ctx.catColl.FindId(id).Select(bson.M{
		"ancestors": true,
		"name":      true,
		"_id":       true,
	}).One(&cat)
	if err != nil {
		return nil, err
	}

	ancestors := make([]Catergory, 0, len(cat.Ancestors)+1)
	err = ctx.catColl.Find(bson.M{"_id": bson.M{"$in": cat.Ancestors}}).Select(bson.M{
		"ancestors": true,
		"name":      true,
		"_id":       true,
	}).Sort("ancestors").All(&ancestors)
	ancestors = append(ancestors, cat)
	return ancestors, nil
}

func (ctx *DBCtx) AddContact(c *Contact) error {
	c.ContactId = bson.NewObjectId()
	c.At = time.Now()
	return ctx.contColl.Insert(c)
}

func (ctx *DBCtx) AllContact(offset bson.ObjectId, limit int) []Contact {
	conts := []Contact{}
	if offset.Valid() {
		ctx.contColl.Find(bson.M{"_id": bson.M{"$gt": offset}}).Limit(limit).All(&conts)
	}
	ctx.contColl.Find(nil).Limit(limit).All(&conts)
	return conts
}
