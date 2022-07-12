package project

import (
	"encoding/json"
	"flow-projects/mysql"
	"strings"
)

type PatchBody struct {
	Name       *string             `json:"name" validate:"omitempty"`
	ThemeColor *string             `json:"theme_color" validate:"omitempty,hexcolor"`
	ParentId   PatchNullJSONUint64 `json:"parent_id" validate:"dive"`
	Pinned     *bool               `json:"pinned" validate:"omitempty"`
	Hidden     *bool               `json:"hidden" validate:"omitempty"`
}

type PatchNullJSONUint64 struct {
	UInt64 **uint64 `validate:"omitempty,gte=1"`
}

func (p *PatchNullJSONUint64) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	var valueP *uint64 = nil
	if string(data) == "null" {
		// key exists and value is null
		p.UInt64 = &valueP
		return nil
	}

	var tmp uint64
	tmpP := &tmp
	if err := json.Unmarshal(data, &tmp); err != nil {
		// invalid value type
		return err
	}
	// valid value
	p.UInt64 = &tmpP
	return nil
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

	// Check name
	if new.Name != nil {
		var projects []Project
		projects, err = GetList(userId, false, new.Name)
		if err != nil {
			return
		}
		if len(projects) != 0 && projects[0].Id != id {
			usedName = true
			return
		}
	}
	if new.Hidden != nil && p.Hidden != *new.Hidden && !*new.Hidden {
		// project.Hidden change to false
		// check duplicate of project.name
		var projects []Project
		if new.Name != nil {
			projects, err = GetList(userId, true, new.Name)
		} else {
			projects, err = GetList(userId, true, &p.Name)
		}
		if err != nil {
			return
		}
		if len(projects) != 0 && projects[0].Id != id {
			usedName = true
			return
		}
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
	if new.ParentId.UInt64 != nil {
		if *new.ParentId.UInt64 != nil {
			queryStr += " parent_id = ?,"
			queryParams = append(queryParams, **new.ParentId.UInt64)
			p.ParentId = *new.ParentId.UInt64
		} else {
			queryStr += " parent_id = ?,"
			queryParams = append(queryParams, nil)
			p.ParentId = nil
		}
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
