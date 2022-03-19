package project

import (
	"flow-projects/mysql"
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
		p := Project{}
		err = rows.Scan(&p.Id, &p.Name, &p.ThemeColor, &p.ParentId, &p.Pinned, &p.Hidden)
		if err != nil {
			return
		}

		if p.ParentId != nil {
			// Embedding sub project
			tmpParent.SubProjects = append(tmpParent.SubProjects, p)
		} else {
			if tmpParent.Id != 0 {
				// Finish embedding
				projects = append(projects, tmpParent)
			}
			// Start embedding to next parent
			tmpParent = ProjectSubEmbed{p.Id, p.Name, p.ThemeColor, p.Pinned, p.Hidden, []Project{}}
		}
	}
	projects = append(projects, tmpParent)

	return
}
