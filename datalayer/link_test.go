package datalayer

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/chytilp/links/model"
	"github.com/google/go-cmp/cmp"
)

func TestShouldGetLink(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	id := 1
	rows := sqlmock.NewRows([]string{"l_id", "link", "l_name", "l_active", "l_created",
		"c_id", "c_name", "parent_id", "c_active", "c_created"}).
		AddRow(id, "post 1", "hello")
	queryTemplate := "^SELECT (.+) FROM link JOIN category c on l.category_id = c.id WHERE l.id ="
	mock.ExpectQuery(queryTemplate).WillReturnRows(rows)

	links := &Links{}
	link, _ := links.Get(id)
	expected := &model.Link{}
	cmp.Equal(expected, link)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
