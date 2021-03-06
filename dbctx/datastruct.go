package dbctx

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type Catergory struct {
	CatId     bson.ObjectId `bson:"_id"`
	Name      string
	Ancestors []bson.ObjectId `bson:",omitempty"`
	Children  []bson.ObjectId `bson:",omitempty"`
	Parent    bson.ObjectId   `bson:",omitempty"`
	LastEntry bson.ObjectId   `bson:",omitempty"`
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

type Tag struct {
	Name  string `bson:"_id"`
	Usage int
}

type Comment struct {
	CommId  bson.ObjectId `bson:"_id"`
	Content string
	By      string
	At      time.Time
	EntryId bson.ObjectId `bson:",omitempty"`
}

type CatergorySummary struct {
	CatId                bson.ObjectId
	CatName              string
	LastEntryId          bson.ObjectId
	LastEntryName        string
	LastEntryDescription string
}

type Contact struct {
	ContactId bson.ObjectId `bson:"_id"`
	Title     string
	Content   string
	From      string
	To        string
	At        time.Time
	New       bool
}

type CatLstTmpl struct {
	Deep      int
	Root      Catergory
	Relatives []Catergory
}

func NewCatLstTmpl(root Catergory, relatives []Catergory) *CatLstTmpl {
	return &CatLstTmpl{0, root, relatives}
}

func (lst *CatLstTmpl) New(rootid bson.ObjectId) *CatLstTmpl {
	var root Catergory
	for _, v := range lst.Relatives {
		if v.CatId == rootid {
			root = v
			break
		}
	}
	return &CatLstTmpl{lst.Deep + 1, root, lst.Relatives}
}
