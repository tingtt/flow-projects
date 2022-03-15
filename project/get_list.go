package project

import (
	"database/sql"
	"flow-projects/mysql"
)

func GetList(userId uint64, show_hidden bool) (projects []Project, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	queryStr := "SELECT id, name, theme_color, parent_id, pinned, `hidden` FROM projects WHERE user_id = ?"
	if !show_hidden {
		queryStr += " AND `hidden` = false"
	}
	queryStr += " ORDER BY pinned DESC"
	stmtOut, err := db.Prepare(queryStr)
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId)
	if err != nil {
		return
	}

	for rows.Next() {
		// TODO: uint64に対応
		var (
			id         uint64
			name       string
			themeColor string
			parentId   sql.NullInt64
			pinned     bool
			hidden     bool
		)
		err = rows.Scan(&id, &name, &themeColor, &parentId, &pinned, &hidden)
		if err != nil {
			return
		}

		p := Project{Id: id, Name: name, ThemeColor: themeColor, Pinned: pinned, Hidden: hidden}
		if parentId.Valid {
			parentIdTmp := uint64(parentId.Int64)
			p.ParentId = &parentIdTmp
		}

		projects = append(projects, p)
	}

	return
}
