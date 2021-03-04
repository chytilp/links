package datalayer

import (
	"github.com/chytilp/links/logging"
	"github.com/chytilp/links/model"
)

const selectPattern = "SELECT l.id AS l_id, l.link, l.name AS l_name, l.active AS l_active, " +
	"l.created AS l_created, c.id AS c_id, c.name AS c_name, c.parent_id, c.active AS c_active, " +
	" c.created AS c_created " +
	"FROM link l " +
	"JOIN category c on l.category_id = c.id "

// Links type wrapps database methods above link table.
type Links struct {
}

// Get method returns link record from link table by id.
func (l *Links) Get(id int) (*model.Link, error) {
	db, err := getDb()
	if err != nil {
		logging.L.Error("Failed to get DB connection. err: %s", err)
		return nil, err
	}

	query := selectPattern + " WHERE l.id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		logging.L.Error("Prepare db query failed. err: %s", err)
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	link := &model.Link{}
	category := &model.Category{}
	err = row.Scan(&link.ID, &link.Link, &link.Name, &link.Active, &link.Created,
		&category.ID, &category.Name, &category.ParentID, &category.Active,
		&category.Created)
	if err != nil {
		logging.L.Error("Row Scan values failed. err: %s", err)
		return nil, err
	}
	link.Category = category
	return link, nil
}

// Save method insert/update record in link table.
func (l *Links) Save(link model.Link) (int, error) {
	return 0, nil
}

// Delete method archive record in link table by id.
func (l *Links) Delete(id int) (int, error) {
	return 0, nil
}

// Retrieve method selects from link table records by sended filers.
func (l *Links) Retrieve(filters map[string]string) ([]*model.Link, error) {
	return nil, nil
}
