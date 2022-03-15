package project

import (
	"database/sql"
	"flow-projects/mysql"
)

func Get(userId uint64, id uint64) (p Project, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return Project{}, false, err
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT name, theme_color, parent_id, pinned, `hidden` FROM projects WHERE user_id = ? AND id = ?")
	if err != nil {
		return Project{}, false, err
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId, id)
	if err != nil {
		return Project{}, false, err
	}

	// TODO: uint64に対応
	var (
		name       string
		themeColor string
		parentId   sql.NullInt64
		pinned     bool
		hidden     bool
	)
	if !rows.Next() {
		// Not found
		return Project{}, true, nil
	}
	err = rows.Scan(&name, &themeColor, &parentId, &pinned, &hidden)
	if err != nil {
		return Project{}, false, err
	}

	p.Id = id
	p.Name = name
	p.ThemeColor = themeColor
	if parentId.Valid {
		parentIdTmp := uint64(parentId.Int64)
		p.ParentId = &parentIdTmp
	}
	p.Pinned = pinned
	p.Hidden = hidden

	return
}

func GetByName(userId uint64, name string) (p Project, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return Project{}, false, err
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT id, theme_color, parent_id, pinned, `hidden` FROM projects WHERE user_id = ? AND name = ?")
	if err != nil {
		return Project{}, false, err
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId, name)
	if err != nil {
		return Project{}, false, err
	}

	// TODO: uint64に対応
	var (
		id         uint64
		themeColor string
		parentId   sql.NullInt64
		pinned     bool
		hidden     bool
	)
	if !rows.Next() {
		// Not found
		return Project{}, true, nil
	}
	err = rows.Scan(&id, &themeColor, &parentId, &pinned, &hidden)
	if err != nil {
		return Project{}, false, err
	}

	p.Name = name
	p.Id = id
	p.ThemeColor = themeColor
	if parentId.Valid {
		parentIdTmp := uint64(parentId.Int64)
		p.ParentId = &parentIdTmp
	}
	p.Pinned = pinned
	p.Hidden = hidden

	return
}
