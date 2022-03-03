package project

import (
	"flow-projects/mysql"
)

type Project struct {
	Id         uint64
	Name       string `json:"name"`
	ThemeColor string `json:"theme_color"`
	ParentId   uint64 `json:"parent_id"`
	Pinned     bool   `json:"pinned"`
	Hidden     bool   `json:"hidden"`
}

type Post struct {
	Name       string `json:"name" validate:"required"`
	ThemeColor string `json:"theme_color" validate:"required,hexcolor"`
	ParentId   uint64 `json:"parent_id" validate:"omitempty,number"`
	Pinned     bool   `json:"pinned" validate:"omitempty"`
}

type Patch struct {
	Name       *string `json:"name" validate:"omitempty"`
	ThemeColor *string `json:"theme_color" validate:"omitempty,hexcolor"`
	ParentId   *uint64 `json:"parent_id" validate:"omitempty,number"`
	Pinned     *bool   `json:"pinned" validate:"omitempty"`
	Hidden     *bool   `json:"hidden" validate:"omitempty"`
}

func Get(user_id uint64, id uint64) (u Project, notFound bool, err error) {
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

	var (
		name        string
		theme_color string
		parent_id   uint64
		pinned      bool
		hidden      bool
	)
	if !rows.Next() {
		// Not found
		return Project{}, true, nil
	}
	err = rows.Scan(&name, &theme_color, &parent_id, &pinned, &hidden)
	if err != nil {
		return Project{}, false, err
	}

	return Project{id, name, theme_color, parent_id, pinned, hidden}, false, nil
}

func GetByName(user_id uint64, name string) (u Project, notFound bool, err error) {
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

	var (
		id          uint64
		theme_color string
		parent_id   uint64
		pinned      bool
		hidden      bool
	)
	if !rows.Next() {
		// Not found
		return Project{}, true, nil
	}
	err = rows.Scan(&id, &theme_color, &parent_id, &pinned, &hidden)
	if err != nil {
		return Project{}, false, err
	}

	return Project{id, name, theme_color, parent_id, pinned, hidden}, false, nil
}

func Insert(user_id uint64, post Post) (p Project, usedName bool, err error) {
	_, notFound, err := GetByName(user_id, post.Name)
	if err != nil {
		return Project{}, false, err
	}
	if !notFound {
		return Project{}, true, nil
	}

	// Insert DB
	db, err := mysql.Open()
	if err != nil {
		return Project{}, false, err
	}
	defer db.Close()
	stmtIns, err := db.Prepare("INSERT INTO projects (user_id, name, theme_color, parent_id, pinned) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return Project{}, false, err
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(user_id, post.Name, post.ThemeColor, post.ParentId, post.Pinned)
	if err != nil {
		return Project{}, false, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Project{}, false, err
	}

	return Project{uint64(id), post.Name, post.ThemeColor, post.ParentId, post.Pinned, false}, false, nil
}

func Update(user_id uint64, id uint64, new Patch) (p Project, usedName bool, notFound bool, err error) {
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
		new.ParentId = &old.ParentId
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

	return Project{id, *new.Name, *new.ThemeColor, *new.ParentId, *new.Pinned, *new.Hidden}, false, false, nil
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

	stmtOut, err := db.Prepare("SELECT id, name, theme_color, parent_id, pinned, `hidden` FROM projects WHERE user_id = ? AND `hidden` = false AND `hidden` = ?")
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(user_id, show_hidden)
	if err != nil {
		return
	}

	for rows.Next() {
		var (
			id          uint64
			name        string
			theme_color string
			parent_id   uint64
			pinned      bool
			hidden      bool
		)

		err = rows.Scan(&id, &name, &theme_color, &parent_id, &pinned, &hidden)
		if err != nil {
			return
		}

		projects = append(projects, Project{id, name, theme_color, parent_id, pinned, hidden})
	}

	return
}
