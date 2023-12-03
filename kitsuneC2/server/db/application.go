package db

import (
	"errors"
)

const MAX_PENDING_TASKS int = 100

// Given an implant ID, returns all data from the implant_info table.
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
		return nil, errors.New("No results for implant with id: " + implantId)
	}

	err = rows.Scan(&info.Id, &info.Name, &info.Public_ip, &info.Os, &info.Arch, &info.Last_checkin)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// Given info about an implant, registers an entry in the implant_info table.
func AddImplant(info *Implant_info) error {
	stmt, err := dbConn.Prepare("INSERT INTO implant_info VALUES (?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(info.Id, info.Name, info.Public_ip, info.Os, info.Arch, info.Last_checkin)
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

// Given an implant ID, returns a list of pending tasks from the implant_tasks table.
func GetTasks(implantId string) ([]*Implant_task, error) {
	stmt, err := dbConn.Prepare("SELECT * FROM implant_tasks WHERE implant_id=?")
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
		rows.Scan(&tasks[i].Task_id, &tasks[i].Implant_id, &tasks[i].Task_type, &tasks[i].Task_data)
		i++
	}

	tasks = tasks[:i]

	return tasks, nil
}

// Given a task, adds it to the implant_tasks table
func AddTask(task *Implant_task) error {
	stmt, err := dbConn.Prepare("INSERT INTO implant_tasks VALUES (?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(task.Task_id, task.Implant_id, task.Task_type, task.Task_data)
	if err != nil {
		return err
	}
	return nil
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

func AddListener() {

}

func RemoveListener() {

}

func AddPayload() {

}

func RemovePayload() {

}
