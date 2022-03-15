package project

import (
	"flow-projects/mysql"
)

type Patch struct {
	Name       *string `json:"name" validate:"omitempty"`
	ThemeColor *string `json:"theme_color" validate:"omitempty,hexcolor"`
	ParentId   *uint64 `json:"parent_id" validate:"omitempty,gte=1"`
	Pinned     *bool   `json:"pinned" validate:"omitempty"`
	Hidden     *bool   `json:"hidden" validate:"omitempty"`
}

func Update(userId uint64, id uint64, new Patch) (_ Project, usedName bool, notFound bool, err error) {
	// Get old
	old, notFound, err := Get(userId, id)
	if err != nil {
		return Project{}, false, false, err
	}
	if notFound {
		return Project{}, false, true, nil
	}

	// Set no update values
	if new.Name == nil {
		new.Name = &old.Name
	}
	if new.ThemeColor == nil {
		new.ThemeColor = &old.ThemeColor
	}
	if new.ParentId == nil {
		new.ParentId = old.ParentId
	}
	if new.Pinned == nil {
		new.Pinned = &old.Pinned
	}
	if new.Hidden == nil {
		new.Hidden = &old.Hidden
	}

	// Update row
	db, err := mysql.Open()
	if err != nil {
		return Project{}, false, false, err
	}
	defer db.Close()
	stmtIns, err := db.Prepare("UPDATE projects SET name = ?, theme_color = ?, parent_id = ?, pinned = ?, `hidden` = ? WHERE user_id = ? AND id = ?")
	if err != nil {
		return Project{}, false, false, err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(new.Name, new.ThemeColor, new.ParentId, new.Pinned, new.Hidden, userId, id)
	if err != nil {
		return Project{}, false, false, err
	}

	return Project{id, *new.Name, *new.ThemeColor, new.ParentId, *new.Pinned, *new.Hidden}, false, false, nil
}
