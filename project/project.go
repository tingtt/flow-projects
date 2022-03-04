package project

import (
	"database/sql"
	"flow-projects/mysql"
)

type Project struct {
	Id         uint64
	Name       string  `json:"name"`
	ThemeColor string  `json:"theme_color"`
	ParentId   *uint64 `json:"parent_id,omitempty"`
	Pinned     bool    `json:"pinned"`
	Hidden     bool    `json:"hidden"`
}

type Post struct {
	Name       string  `json:"name" validate:"required"`
	ThemeColor string  `json:"theme_color" validate:"required,hexcolor"`
	ParentId   *uint64 `json:"parent_id" validate:"omitempty,number"`
	Pinned     bool    `json:"pinned" validate:"omitempty"`
	Hidden     bool    `json:"hidden" validate:"omitempty"`
}

type Patch struct {
	Name       *string `json:"name" validate:"omitempty"`
	ThemeColor *string `json:"theme_color" validate:"omitempty,hexcolor"`
	ParentId   *uint64 `json:"parent_id" validate:"omitempty,number"`
	Pinned     *bool   `json:"pinned" validate:"omitempty"`
	Hidden     *bool   `json:"hidden" validate:"omitempty"`
}

func Get(user_id uint64, id uint64) (p Project, notFound bool, err error) {
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

	rows, err := stmtOut.Query(user_id, id)
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

func GetByName(user_id uint64, name string) (p Project, notFound bool, err error) {
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

	rows, err := stmtOut.Query(user_id, name)
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

func Insert(user_id uint64, post Post) (p Project, invalidParentId bool, err error) {
	// Check parent id
	if post.ParentId != nil {
		_, notFound, err := Get(user_id, *post.ParentId)
		if err != nil {
			return Project{}, false, err
		}
		if notFound {
			return Project{}, true, nil
		}
	}

	// Insert DB
	db, err := mysql.Open()
	if err != nil {
		return Project{}, false, err
	}
	defer db.Close()
	stmtIns, err := db.Prepare("INSERT INTO projects (user_id, name, theme_color, parent_id, pinned, `hidden`) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return Project{}, false, err
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(user_id, post.Name, post.ThemeColor, post.ParentId, post.Pinned, post.Hidden)
	if err != nil {
		return Project{}, false, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Project{}, false, err
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

func Update(user_id uint64, id uint64, new Patch) (_ Project, usedName bool, notFound bool, err error) {
	// Get old
	old, notFound, err := Get(user_id, id)
	if err != nil {
		return Project{}, false, false, err
	}
	if notFound {
		return Project{}, false, true, nil
	}

	// Check deplicate name
	if new.Name != nil {
		named, notFound, err := GetByName(user_id, *new.Name)
		if err != nil {
			return Project{}, false, false, err
		}

		if !notFound && named.Id != id {
			// Duplicated name
			return Project{}, true, false, nil
		}
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
	_, err = stmtIns.Exec(new.Name, new.ThemeColor, new.ParentId, new.Pinned, new.Hidden, user_id, id)
	if err != nil {
		return Project{}, false, false, err
	}

	return Project{id, *new.Name, *new.ThemeColor, new.ParentId, *new.Pinned, *new.Hidden}, false, false, nil
}

func Delete(user_id uint64, id uint64) (notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return false, err
	}
	defer db.Close()
	stmtIns, err := db.Prepare("DELETE FROM projects WHERE user_id = ? AND id = ?")
	if err != nil {
		return false, err
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(user_id, id)
	if err != nil {
		return false, err
	}
	affectedRowCount, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if affectedRowCount == 0 {
		// Not found
		return true, nil
	}

	return false, nil
}

func GetList(user_id uint64, show_hidden bool) (projects []Project, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	// TODO: 子をネスト表示

	queryStr := "SELECT id, name, theme_color, parent_id, pinned, `hidden` FROM projects WHERE user_id = ?"
	if !show_hidden {
		queryStr += " AND `hidden` = false"
	}
	stmtOut, err := db.Prepare(queryStr)
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(user_id)
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
