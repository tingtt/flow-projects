package project

import (
	"flow-projects/mysql"
)

func Get(userId uint64, id uint64) (p Project, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT name, theme_color, parent_id, pinned, `hidden` FROM projects WHERE user_id = ? AND id = ?")
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId, id)
	if err != nil {
		return
	}

	if !rows.Next() {
		// Not found
		notFound = true
		return
	}
	err = rows.Scan(&p.Name, &p.ThemeColor, &p.ParentId, &p.Pinned, &p.Hidden)
	if err != nil {
		return
	}

	p.Id = id
	return
}

func GetByName(userId uint64, name string) (p Project, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT id, theme_color, parent_id, pinned, `hidden` FROM projects WHERE user_id = ? AND name = ?")
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId, name)
	if err != nil {
		return
	}
	defer rows.Close()

	if !rows.Next() {
		// Not found
		notFound = true
		return
	}
	err = rows.Scan(&p.Id, &p.ThemeColor, &p.ParentId, &p.Pinned, &p.Hidden)
	if err != nil {
		return
	}

	p.Name = name
	return
}
