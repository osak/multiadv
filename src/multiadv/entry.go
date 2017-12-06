package multiadv

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

type Entry struct {
	Year  int
	Month int
	Day   int
	Title string
	Url   string
}

type EntryDAO interface {
	FetchEntriesOfMonth(year, month int) ([]Entry, error)
	Create(entry Entry) error
}

type entryDAO struct {
	db *mgo.Database
}

func NewEntryDAO(db *mgo.Database) EntryDAO {
	return &entryDAO{
		db: db,
	}
}

func (e *entryDAO) FetchEntriesOfMonth(year, month int) ([]Entry, error) {
	var result []Entry

	err := e.db.C("entries").Find(bson.M{
		"year":  year,
		"month": month,
	}).All(&result)
	if err != nil {
		return nil, fmt.Errorf("entry.fetchEntriesOfMonth: cannot fetch entries of year=%d, month=%d\n%s", year, month, err.Error())
	}
	return result, nil
}

func (e *entryDAO) Create(entry Entry) error {
	err := e.db.C("entries").Insert(entry)
	if err != nil {
		return fmt.Errorf("entry.Create: failed to insert entry %s\n%s", entry.Title, err.Error())
	}
	return nil
}
