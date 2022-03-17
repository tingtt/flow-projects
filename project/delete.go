package project

import (
	"flow-projects/mysql"
)

func Delete(userId uint64, id uint64) (notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmtIns, err := db.Prepare("DELETE FROM projects WHERE user_id = ? AND id = ?")
	if err != nil {
		return
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(userId, id)
	if err != nil {
		return
	}
	affectedRowCount, err := result.RowsAffected()
	if err != nil {
		return
	}
	notFound = affectedRowCount == 0

	return
}
