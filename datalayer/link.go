package datalayer

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/chytilp/links/model"
)

// Links type wrapps database methods above link table.
type Links struct {
	records         *records
	fieldsForSelect []string
	selectPattern   string
	insertPattern   string
	updatePattern   string
	deletePattern   string
}

// CreateLinks creates and returns instance of Links struct.
func CreateLinks(db *sql.DB) *Links {
	if db == nil {
		db = getDb()
	}
	links := &Links{
		records: newRecords(db),
		fieldsForSelect: []string{"l_id", "link", "l_name", "l_active", "l_created",
			"c_id", "c_name", "parent_id", "c_active", "c_created"},
		selectPattern: "SELECT l.id AS l_id, l.link, l.name AS l_name, l.active AS l_active, " +
			"l.created AS l_created, c.id AS c_id, c.name AS c_name, c.parent_id, c.active AS c_active, " +
			" c.created AS c_created " +
			"FROM link l " +
			"JOIN category c on l.category_id = c.id ",
		insertPattern: "INSERT INTO link(link, name, category_id) " +
			"VALUES(?, ?, ?)",
		updatePattern: "UPDATE link SET link=?, name=?, category_id=? WHERE id=?",
		deletePattern: "UPDATE link SET active=? WHERE id=?",
	}
	return links
}

// Get method returns link record from link table by id.
func (l *Links) Get(id int) (*model.Link, error) {
	row := l.records.db.QueryRow(l.selectPattern+" WHERE l.id = ?", id)
	link, err := l.scanRow(row.Scan)
	if err != nil {
		return nil, err
	}
	return link, nil
}

// Save method insert/update record in link table.
func (l *Links) Save(link model.Link) (*model.Link, error) {
	var id int
	var err error
	if link.ID > 0 {
		err := l.update(link)
		if err != nil {
			return nil, err
		}
		id = link.ID
	} else {
		id, err = l.insert(link)
		if err != nil {
			return nil, err
		}
	}
	return l.Get(id)
}

// insert new record to link table.
func (l *Links) insert(link model.Link) (int, error) {
	values := []interface{}{
		link.Link,
		link.Name,
		link.Category.ID,
	}
	id, err := l.records.insert(values, l.insertPattern)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// update record in link table.
func (l *Links) update(link model.Link) error {
	values := []interface{}{
		link.Link,
		link.Name,
		link.Category.ID,
		link.ID,
	}
	err := l.records.update(values, l.updatePattern)
	if err != nil {
		return err
	}
	return nil
}

// Delete method archives record in link table by id.
func (l *Links) Delete(id int, time time.Time) (*model.Link, error) {
	values := []interface{}{
		time,
		id,
	}
	err := l.records.update(values, l.deletePattern)
	if err != nil {
		return nil, err
	}
	return l.Get(id)
}

// Retrieve method selects from link table records by sended filers.
func (l *Links) Retrieve(filters map[string][]string) ([]*model.Link, error) {
	whereClause, values, err := l.buildFilters(filters)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if len(whereClause) != 0 {
		query := l.selectPattern + "WHERE " + whereClause
		rows, err = l.records.db.Query(query, values...)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = l.records.db.Query(l.selectPattern)
		if err != nil {
			return nil, err
		}
	}
	var result []*model.Link
	for rows.Next() {
		link, err := l.scanRow(rows.Scan)
		if err != nil {
			return nil, err
		}
		result = append(result, link)
	}
	return result, nil
}

// buildFilters creates sql filter and prepare values from request input.
func (l *Links) buildFilters(filters map[string][]string) (string, []interface{}, error) {
	if len(filters) == 0 {
		return "", nil, nil
	}

	var filter string
	var outValues []interface{}
	operator := " "
	for field, values := range filters {
		if len(filter) > 0 {
			operator = " AND "
		}
		field = strings.ToLower(field)
		if !l.isFieldAllowed(field) {
			return "", nil, fmt.Errorf("Field %s is not allowed", field)
		}
		if len(values) == 1 {
			filter += operator + correctFieldName(field) + " = ?"
		} else {
			inClauseVals := strings.Repeat("?, ", len(values))
			filter += operator + correctFieldName(field) + " IN (" +
				inClauseVals[:len(inClauseVals)-2] + ")"
			for _, val := range values {
				typedValue, err := l.convertValue(field, val)
				if err != nil {
					return "", nil, err
				}
				outValues = append(outValues, typedValue)
			}
		}
	}
	return filter, outValues, nil
}

// correctFieldName repairs field name, replace _ by .
func correctFieldName(field string) string {
	return strings.Replace(field, "_", ".", 1)
}

// isFieldAllowed checks if allowed field comes from request.
func (l *Links) isFieldAllowed(field string) bool {
	sort.Strings(l.fieldsForSelect)
	idx := sort.SearchStrings(l.fieldsForSelect, field)
	return idx != len(l.fieldsForSelect) && l.fieldsForSelect[idx] == field
}

// scanRow fills link structure with values from db record.
func (l *Links) scanRow(fn scanner) (*model.Link, error) {
	link := &model.Link{}
	category := &model.Category{}
	err := fn(&link.ID, &link.Link, &link.Name, &link.Active, &link.Created,
		&category.ID, &category.Name, &category.ParentID, &category.Active,
		&category.Created)
	if err != nil {
		return nil, err
	}
	link.Category = category
	return link, nil
}

// convertValues converts value from request to correct type.
func (l *Links) convertValue(field string, value string) (interface{}, error) {
	switch field {
	case "l_id", "c_id", "parent_id":
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		return intVal, nil
	case "l_name", "c_name", "link":
		return value, nil
	case "l_active", "l_created", "c_active", "c_created":
		// 2014-11-12T11:45:26.371Z
		t, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return nil, err
		}
		return t, nil
	}
	return nil, fmt.Errorf("Unknown field %s", field)
}

// Close db connection
func (l *Links) Close() error {
	return l.records.close()
}
