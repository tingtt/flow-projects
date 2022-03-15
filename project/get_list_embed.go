package project

import (
	"database/sql"
	"flow-projects/mysql"
	"sort"
)

type ProjectSubEmbed struct {
	Id          uint64    `json:"id"`
	Name        string    `json:"name"`
	ThemeColor  string    `json:"theme_color"`
	Pinned      bool      `json:"pinned"`
	Hidden      bool      `json:"hidden"`
	SubProjects []Project `json:"sub_projects,omitempty"`
}

func GetListEmbed(userId uint64, show_hidden bool) (projects []ProjectSubEmbed, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	var queryStr string
	if show_hidden {
		queryStr =
			`WITH RECURSIVE cte AS (
				SELECT id, name, theme_color, parent_id, pinned, ` + "`hidden`" + `
				FROM projects
				WHERE user_id = ? AND parent_id is NULL
			UNION ALL
				SELECT childs.id, childs.name, childs.theme_color, childs.parent_id, childs.pinned, childs.` + "`hidden`" + `
				FROM projects AS childs, cte
				WHERE childs.user_id = ? AND cte.id = childs.parent_id
			)
			SELECT * FROM cte ORDER BY COALESCE(parent_id, id), parent_id IS NULL DESC, pinned DESC`
	} else {
		queryStr =
			`WITH RECURSIVE cte AS (
				SELECT id, name, theme_color, parent_id, pinned, ` + "`hidden`" + `
				FROM projects
				WHERE user_id = ? AND parent_id is NULL AND ` + "`hidden`" + ` = false
			UNION ALL
				SELECT childs.id, childs.name, childs.theme_color, childs.parent_id, childs.pinned, childs.` + "`hidden`" + `
				FROM projects AS childs, cte
				WHERE childs.user_id = ? AND ` + "childs.`hidden`" + ` = false AND cte.id = childs.parent_id
			)
			SELECT * FROM cte ORDER BY COALESCE(parent_id, id), parent_id IS NULL DESC, pinned DESC`
	}
	stmtOut, err := db.Prepare(queryStr)
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId, userId)
	if err != nil {
		return
	}

	// Emdedding sub projects
	var tmpParent ProjectSubEmbed
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

		if parentId.Valid {
			// Embedding sub project
			tmpParentId := uint64(parentId.Int64)
			tmpParent.SubProjects = append(tmpParent.SubProjects, Project{id, name, themeColor, &tmpParentId, pinned, hidden})
		} else {
			if tmpParent.Id != 0 {
				// Finish embedding
				projects = append(projects, tmpParent)
			}
			// Start embedding to next parent
			tmpParent = ProjectSubEmbed{id, name, themeColor, pinned, hidden, []Project{}}
		}
	}
	projects = append(projects, tmpParent)

	// Sort by pinned of parent projects
	sort.Slice(projects, func(i, j int) bool {
		return !projects[j].Pinned && projects[i].Pinned
	})

	return
}