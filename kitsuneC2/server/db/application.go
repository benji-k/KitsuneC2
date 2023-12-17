package db

import (
	"database/sql"
	"errors"
)

const (
	MAX_PENDING_TASKS int = 100
)

var (
	ErrNoResults error = errors.New("no results for query") //Used in all Get* functions
)

// Returns information about all active implants
func GetAllImplants() ([]*Implant_info, error) {
	stmt, err := dbConn.Prepare("SELECT * FROM implant_info")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var output []*Implant_info
	for rows.Next() {
		info := new(Implant_info)
		rows.Scan(&info.Id, &info.Name, &info.Public_ip, &info.Os, &info.Arch, &info.Last_checkin, &info.Username, &info.Uid, &info.Gid, &info.Hostname, &info.Active)
		output = append(output, info)
	}
	if len(output) == 0 {
		return nil, ErrNoResults
	}

	return output, nil
}

// Given an implant ID, returns all data from the implant_info table. Returns a db.ErrNoResults error if there were no results
// for the specified ID
func GetImplantInfo(implantId string) (*Implant_info, error) {
	stmt, err := dbConn.Prepare("SELECT * FROM implant_info WHERE id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(implantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var info *Implant_info = new(Implant_info)
	hasResult := rows.Next()
	if !hasResult {
		return nil, ErrNoResults
	}

	err = rows.Scan(&info.Id, &info.Name, &info.Public_ip, &info.Os, &info.Arch, &info.Last_checkin, &info.Username, &info.Uid, &info.Gid, &info.Hostname, &info.Active)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// Given info about an implant, registers an entry in the implant_info table.
func AddImplant(info *Implant_info) error {
	stmt, err := dbConn.Prepare("INSERT INTO implant_info VALUES (?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(info.Id, info.Name, info.Public_ip, info.Os, info.Arch, info.Last_checkin, info.Username, info.Uid, info.Gid, info.Hostname, info.Active)
	if err != nil {
		return err
	}
	return nil
}

// Given an implant ID, removes all information related to said implant from ALL tables.
func RemoveImplant(implantId string) error {
	stmt1, _ := dbConn.Prepare("DELETE FROM implant_info WHERE id=?")
	stmt2, _ := dbConn.Prepare("DELETE FROM implant_tasks WHERE implant_id=?")
	defer stmt1.Close()
	defer stmt2.Close()

	res1, err1 := stmt1.Exec(implantId)
	res2, err2 := stmt2.Exec(implantId)
	if err1 != nil || err2 != nil {
		return errors.New("error for implant_info table: " + err1.Error() + ". Error for implant_tasks table: " + err2.Error())
	}

	entriesDeleted1, _ := res1.RowsAffected()
	entriesDeleted2, _ := res2.RowsAffected()

	if entriesDeleted1 == 0 && entriesDeleted2 == 0 {
		return errors.New("no database entries for implant with ID: " + implantId)
	}

	return nil
}

// Changes the "active" status of an implant to "status"
func SetImplantStatus(implantId string, status bool) error {
	stmt, err := dbConn.Prepare("UPDATE implant_info SET active=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(status, implantId)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("No implant with id: " + implantId)
	}

	return nil
}

// Gets the "active" status of an implant
func GetImplantStatus(implantId string) (bool, error) {
	stmt, err := dbConn.Prepare("SELECT active FROM implant_info WHERE id=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(implantId)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	hasResult := rows.Next()
	if !hasResult {
		return false, ErrNoResults
	}
	var output bool
	rows.Scan(&output)

	return output, nil
}

// Given an implant ID, returns tasks belonging to that implant. The completed paramter dictates whether the tasks returned are
// completed or pending
func GetTasks(implantId string, completed bool) ([]*Implant_task, error) {
	var stmt *sql.Stmt
	var err error
	if completed {
		stmt, err = dbConn.Prepare("SELECT * FROM implant_tasks WHERE implant_id=? AND completed=TRUE")
	} else {
		stmt, err = dbConn.Prepare("SELECT * FROM implant_tasks WHERE implant_id=? AND completed=FALSE")
	}
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(implantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Implant_task = make([]*Implant_task, MAX_PENDING_TASKS)
	var i int = 0
	for rows.Next() {
		tasks[i] = new(Implant_task)
		rows.Scan(&tasks[i].Task_id, &tasks[i].Implant_id, &tasks[i].Task_type, &tasks[i].Task_data, &tasks[i].Completed, &tasks[i].Task_result)
		i++
	}
	if i == 0 {
		return nil, ErrNoResults
	}

	tasks = tasks[:i]

	return tasks, nil
}

// Given a task, adds it to the implant_tasks table
func AddTask(task *Implant_task) error {
	stmt, err := dbConn.Prepare("INSERT INTO implant_tasks VALUES (?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(task.Task_id, task.Implant_id, task.Task_type, task.Task_data, task.Completed, task.Task_result)
	if err != nil {
		return err
	}
	return nil
}

// removes a non-completed task for a specific implant
func RemovePendingTaskForImplant(implantId string, taskId string) error {
	stmt, err := dbConn.Prepare("DELETE FROM implant_tasks WHERE implant_id=? AND task_id=? AND completed=false")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(implantId, taskId)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("No task with ID: " + taskId + " for implant with ID: " + implantId)
	}
	return nil
}

// Given a taskId, returns all information related to said task
func GetTask(taskId string) (*Implant_task, error) {
	stmt, err := dbConn.Prepare("SELECT * FROM implant_tasks WHERE task_id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(taskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	task := new(Implant_task)
	hasResult := rows.Next()
	if !hasResult {
		return nil, ErrNoResults
	}
	rows.Scan(&task.Task_id, &task.Implant_id, &task.Task_type, &task.Task_data, &task.Completed, &task.Task_result)

	return task, nil
}

// Given a task ID, removes it from the implant_tasks table
func RemoveTask(taskId string) error {
	stmt, _ := dbConn.Prepare("DELETE FROM implant_tasks WHERE task_id=?")
	defer stmt.Close()

	result, err := stmt.Exec(taskId)
	if err != nil {
		return err
	}
	entriesDeleted, _ := result.RowsAffected()
	if entriesDeleted == 0 {
		return errors.New("No database entries for task with ID: " + taskId)
	}
	return nil
}

// Changes status of a task from pending to complete. The taskResult parameter is optional.
func CompleteTask(taskId string, taskResult []byte) error {
	stmt, err := dbConn.Prepare("UPDATE implant_tasks SET completed=true, task_result=? WHERE task_id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(taskResult, taskId)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("No task with id: " + taskId)
	}

	return nil
}

// Given an implant ID and time of last checkin in Unix time format, updates the database entry.
func UpdateLastCheckin(implantId string, time int) error {
	stmt, err := dbConn.Prepare("UPDATE implant_info SET last_checkin=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(time, implantId)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("No implant with id: " + implantId)
	}
	return nil
}
