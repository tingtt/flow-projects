package project

import (
	"flow-projects/mysql"
)

func GetList(userId uint64, show_hidden bool, name *string) (projects []Project, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	queryStr := "SELECT id, name, theme_color, parent_id, pinned, `hidden` FROM projects WHERE user_id = ?"
	queryParams := []interface{}{userId}
	if !show_hidden {
		queryStr += " AND `hidden` = false"
	}
	if name != nil {
		queryStr += " AND name = ?"
		queryParams = append(queryParams, *name)
	}
	queryStr += " ORDER BY pinned DESC"
	stmtOut, err := db.Prepare(queryStr)
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(queryParams...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		p := Project{}
		err = rows.Scan(&p.Id, &p.Name, &p.ThemeColor, &p.ParentId, &p.Pinned, &p.Hidden)
		if err != nil {
			return
		}
		projects = append(projects, p)
	}

	return
}
