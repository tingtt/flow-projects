package project

import "flow-projects/mysql"

type PostBody struct {
	Name       string  `json:"name" validate:"required"`
	ThemeColor string  `json:"theme_color" validate:"required,hexcolor"`
	ParentId   *uint64 `json:"parent_id" validate:"omitempty,gte=1"`
	Pinned     bool    `json:"pinned" validate:"omitempty"`
	Hidden     bool    `json:"hidden" validate:"omitempty"`
}

func Post(userId uint64, post PostBody) (p Project, parentNotFound bool, parentHasParent bool, err error) {
	// Check parent id
	if post.ParentId != nil {
		var parent Project
		parent, parentNotFound, err = Get(userId, *post.ParentId)
		if err != nil {
			return
		}
		if parentNotFound {
			return
		}
		if parent.ParentId != nil {
			parentHasParent = true
			return
		}
	}

	// Insert DB
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmtIns, err := db.Prepare("INSERT INTO projects (user_id, name, theme_color, parent_id, pinned, `hidden`) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(userId, post.Name, post.ThemeColor, post.ParentId, post.Pinned, post.Hidden)
	if err != nil {
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	p.Id = uint64(id)
	p.Name = post.Name
	p.ThemeColor = post.ThemeColor
	if post.ParentId != nil {
		p.ParentId = post.ParentId
	}
	p.Pinned = post.Pinned
	p.Hidden = post.Hidden

	return
}
