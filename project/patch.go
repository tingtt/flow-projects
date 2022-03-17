package project

import (
	"flow-projects/mysql"
	"strings"
)

type PatchBody struct {
	Name       *string `json:"name" validate:"omitempty"`
	ThemeColor *string `json:"theme_color" validate:"omitempty,hexcolor"`
	ParentId   *uint64 `json:"parent_id" validate:"omitempty,gte=1"`
	Pinned     *bool   `json:"pinned" validate:"omitempty"`
	Hidden     *bool   `json:"hidden" validate:"omitempty"`
}

func Patch(userId uint64, id uint64, new PatchBody) (p Project, usedName bool, notFound bool, err error) {
	// Get old
	p, notFound, err = Get(userId, id)
	if err != nil {
		return
	}
	if notFound {
		return
	}

	// Generate query
	queryStr := "UPDATE projects SET"
	var queryParams []interface{}
	// Set no update values
	if new.Name != nil {
		queryStr += " name = ?,"
		queryParams = append(queryParams, new.Name)
		p.Name = *new.Name
	}
	if new.ThemeColor != nil {
		queryStr += " theme_color = ?,"
		queryParams = append(queryParams, new.ThemeColor)
		p.ThemeColor = *new.ThemeColor
	}
	if new.ParentId != nil {
		queryStr += " parent_id = ?,"
		queryParams = append(queryParams, new.ParentId)
		p.ParentId = new.ParentId
	}
	if new.Pinned != nil {
		queryStr += " pinned = ?,"
		queryParams = append(queryParams, new.Pinned)
		p.Pinned = *new.Pinned
	}
	if new.Hidden != nil {
		queryStr += " `hidden` = ?"
		queryParams = append(queryParams, new.Hidden)
		p.Hidden = *new.Hidden
	}
	queryStr = strings.TrimRight(queryStr, ",")
	queryStr += " WHERE user_id = ? AND id = ?"
	queryParams = append(queryParams, userId, id)

	// Update row
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmtIns, err := db.Prepare(queryStr)
	if err != nil {
		return
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(queryParams...)
	if err != nil {
		return
	}

	return
}
