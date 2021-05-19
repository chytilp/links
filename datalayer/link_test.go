package datalayer

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/chytilp/links/model"
	"github.com/google/go-cmp/cmp"
)

func createLink(id int, name string) *model.Link {
	created := time.Date(2021, 3, 11, 0, 0, 0, 0, time.UTC)
	expectedCategory := &model.Category{
		ID:       id,
		Name:     "Category 1",
		ParentID: 0,
		Active:   nil,
		Created:  &created,
	}
	expectedLink := &model.Link{
		ID:       id,
		Link:     "https://link1.cz",
		Name:     name,
		Category: expectedCategory,
		Active:   nil,
		Created:  &created,
	}
	return expectedLink
}

func createMockGetExpectedQuery(mock sqlmock.Sqlmock, link *model.Link, id int) {
	columns := []string{"l_id", "link", "l_name", "l_active", "l_created", "c_id",
		"c_name", "parent_id", "c_active", "c_created"}
	mock.ExpectQuery("^SELECT (.+) FROM link l JOIN category c on l.category_id = c.id WHERE l.id = ?").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(
			link.ID, link.Link, link.Name, link.Active,
			link.Created, link.Category.ID, link.Category.Name,
			link.Category.ParentID, link.Category.Active, link.Category.Created))
}

func createMockRetrieveExpectedQuery(mock sqlmock.Sqlmock, links []*model.Link) {
	columns := []string{"l_id", "link", "l_name", "l_active", "l_created", "c_id",
		"c_name", "parent_id", "c_active", "c_created"}
	rows := sqlmock.NewRows(columns)
	for _, link := range links {
		rows.AddRow(link.ID, link.Link, link.Name, link.Active,
			link.Created, link.Category.ID, link.Category.Name,
			link.Category.ParentID, link.Category.Active, link.Category.Created)
	}
	mock.ExpectQuery("^SELECT (.+) FROM link l JOIN category c on l.category_id = c.id "+
		"WHERE l.id IN \\(\\?, \\?\\) AND l.name IN \\(\\?, \\?\\)").
		WithArgs(1, 3, "tenis", "fotbal").
		WillReturnRows(rows)
}

func createMockInsertExpectedQuery(mock sqlmock.Sqlmock, link *model.Link, id int) {
	mock.ExpectPrepare("^INSERT INTO link\\(link, name, category_id\\) VALUES\\(\\?, \\?, \\?\\)").
		ExpectExec().
		WithArgs(link.Link, link.Name, link.Category.ID).
		WillReturnResult(sqlmock.NewResult(int64(id), 1))
}

func createMockUpdateExpectedQuery(mock sqlmock.Sqlmock, link *model.Link) {
	mock.ExpectPrepare("^UPDATE link SET link=\\?, name=\\?, category_id=\\? WHERE id=\\?").
		ExpectExec().
		WithArgs(link.Link, link.Name, link.Category.ID, link.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func createMockDeleteExpectedQuery(mock sqlmock.Sqlmock, link *model.Link) {
	mock.ExpectPrepare("^UPDATE link SET active=\\? WHERE id=\\?").
		ExpectExec().
		WithArgs(link.Active, link.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func TestLinkGetShouldReturnRecord(t *testing.T) {
	db, mock, _ := sqlmock.New()
	id := 1
	expectedLink := createLink(id, "link 1")
	createMockGetExpectedQuery(mock, expectedLink, id)
	links := CreateLinks(db)
	defer links.Close()
	link, err := links.Get(id)
	if err != nil {
		t.Errorf("Links.Get[%d] should return result, but error: %v", id, err)
	}
	same := cmp.Equal(expectedLink, link)
	if !same {
		t.Errorf("Links object are different: %#v, %#v", expectedLink, link)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLinkSaveShouldInsertRecord(t *testing.T) {
	db, mock, _ := sqlmock.New()
	link := createLink(0, "link 1")
	id := 1
	createMockInsertExpectedQuery(mock, link, id)
	createMockGetExpectedQuery(mock, link, id)
	links := CreateLinks(db)
	defer links.Close()
	outputLink, err := links.Save(*link)
	if err != nil {
		t.Errorf("Links.Save[%#v] should insert record, but error: %v", link, err)
	}
	link.ID = outputLink.ID
	same := cmp.Equal(outputLink, link)
	if !same {
		t.Errorf("Links object are different: %#v, %#v", outputLink, link)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLinkSaveShouldUpdateRecord(t *testing.T) {
	db, mock, _ := sqlmock.New()
	id := 1
	link := createLink(id, "link 1")
	createMockUpdateExpectedQuery(mock, link)
	createMockGetExpectedQuery(mock, link, id)
	links := CreateLinks(db)
	defer links.Close()
	outputLink, err := links.Save(*link)
	if err != nil {
		t.Errorf("Links.Save[%#v] should update record, but error: %v", link, err)
	}
	same := cmp.Equal(outputLink, link)
	if !same {
		t.Errorf("Links object are different: %#v, %#v", outputLink, link)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLinkDeleteShouldArchiveRecord(t *testing.T) {
	db, mock, _ := sqlmock.New()
	id := 1
	link := createLink(id, "link 1")
	now := time.Now()
	link.Active = &now
	createMockDeleteExpectedQuery(mock, link)
	createMockGetExpectedQuery(mock, link, id)
	links := CreateLinks(db)
	defer links.Close()
	outputLink, err := links.Delete(id, now)
	if err != nil {
		t.Errorf("Links.Delete[%d, %s] should archive record, but error: %v",
			id, now.Format("2006-01-02 15:04:05"), err)
	}
	same := cmp.Equal(outputLink, link)
	if !same {
		t.Errorf("Links object are different: %#v, %#v", outputLink, link)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLinkRetriveShouldReturnRecords(t *testing.T) {
	db, mock, _ := sqlmock.New()
	links := make([]*model.Link, 2)
	links[0] = createLink(1, "tenis")
	links[1] = createLink(3, "fotbal")
	createMockRetrieveExpectedQuery(mock, links)
	linksObj := CreateLinks(db)
	defer linksObj.Close()
	filters := make(map[string][]string)
	filters["l_id"] = []string{"1", "3"}
	filters["l_name"] = []string{"tenis", "fotbal"}
	outputLinks, err := linksObj.Retrieve(filters)
	if err != nil {
		t.Errorf("Links.Retrieve[%v] should retrieve records, but error: %v",
			filters, err)
	}
	same := cmp.Equal(outputLinks, links)
	if !same {
		t.Errorf("Links object are different: %#v, %#v", outputLinks, links)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
