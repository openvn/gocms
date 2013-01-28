package dbctx

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type Role string

var DefaultRole [28]Role = [28]Role{
	"Chủ hộ",
	"Ông",
	"Bà",
	"Ông cố",
	"Bà cố",
	"Ông nội",
	"Bà nội",
	"Ông ngoại",
	"Bà ngoại",
	"Cha",
	"Mẹ",
	"Cha chồng",
	"Mẹ chồng",
	"Cha vợ",
	"Mẹ vợ",
	"Anh",
	"Chị",
	"Em",
	"Anh rễ",
	"Chị dâu",
	"Em họ",
	"Anh họ",
	"Chị họ",
	"Em rễ",
	"Em dâu",
	"Con",
	"Cháu",
	"Khác",
}

type Class struct {
	Title  string
	School string
}

type Work struct {
	Title  string
	Office string
}

type Income struct {
	Amount int
	Form   string
}

type Person struct {
	PersonId        bson.ObjectId `bson:"_id"`
	Role            Role
	FullName        string
	Gender          bool
	Birth           time.Time
	Quals           int
	Area            string
	Orgs            []string
	AttendingSchool bool
	Class           Class `bson:",omitempty"`
	Working         bool
	Work            Work `bson:",omitempty"`
	AvgIncome       int
	Incomes         []Income `bson:",omitempty"`
	Desire          string   `bson:",omitempty"`
	HI              int
	Note            string `bson:",omitempty"`
}

type House struct {
	HouseId  bson.ObjectId `bson:"_id"`
	Group    int
	Block    int
	Ward     string `bson:",omitempty"`
	District string `bson:",omitempty"`
	City     string `bson:",omitempty"`
	Street   string `bson:",omitempty"`
	Address  string `bson:",omitempty"`
	Member   []Person
}
