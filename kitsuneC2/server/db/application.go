package db

import (
	"database/sql"
	"errors"
	"log"
)

const (
	MAX_PENDING_TASKS int = 100
)

var (
	ErrNoResults error = errors.New("no results for query") //Error that indicates no results where found for a given database query.
)

// Returns information about all active implants
func GetAllImplants() ([]*Implant_info, error) {
	log.Printf("[INFO] db: executing statement: SELECT * FROM implant_info")
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
	log.Printf("[INFO] db: executing statement: SELECT * FROM implant_info WHERE id=%s", implantId)
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
	log.Printf("[INFO] db: executing statement: INSERT INTO implant_info VALUES (?,?,?,?,?,?,?,?,?,?,?)")
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
	log.Printf("[INFO] db: executing statement: DELETE FROM implant_info WHERE id=%s", implantId)
	log.Printf("[INFO] db: executing statement: DELETE FROM implant_tasks WHERE implant_id=%s", implantId)
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
	log.Printf("[INFO] db: executing statement: UPDATE implant_info SET active=%t WHERE id=%s", status, implantId)
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
	log.Printf("[INFO] db: executing statement: SELECT active FROM implant_info WHERE id=%s", implantId)
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
func GetTasksForImplant(implantId string, completed bool) ([]*Implant_task, error) {
	log.Printf("[INFO] db: executing statement: SELECT * FROM implant_tasks WHERE implant_id=%s AND completed=%t", implantId, completed)
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

// Fetches all tasks for every implant from the database. The completed parameter dictates whether the tasks that are fetched are
// completed or not.
func GetAllTasks(completed bool) ([]*Implant_task, error) {
	log.Printf("[INFO] db: executing statement: SELECT * FROM implant_tasks WHERE completed=%t", completed)
	var stmt *sql.Stmt
	var err error
	if completed {
		stmt, err = dbConn.Prepare("SELECT * FROM implant_tasks WHERE completed=TRUE")
	} else {
		stmt, err = dbConn.Prepare("SELECT * FROM implant_tasks WHERE completed=FALSE")
	}
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
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
	log.Printf("[INFO] db: executing statement: INSERT INTO implant_tasks VALUES (?,?,?,?,?,?)")
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
	log.Printf("[INFO] db: executing statement: DELETE FROM implant_tasks WHERE implant_id=%s AND task_id=%s AND completed=false", implantId, taskId)

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
	log.Printf("[INFO] db: executing statement: SELECT * FROM implant_tasks WHERE task_id=%s", taskId)
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
	log.Printf("[INFO] db: executing statement: DELETE FROM implant_tasks WHERE task_id=%s", taskId)
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
	log.Printf("[INFO] db: executing statement: UPDATE implant_tasks SET completed=true, task_result=? WHERE task_id=%s", taskId)
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
	log.Printf("[INFO] db: executing statement: UPDATE implant_info SET last_checkin=%d WHERE id=%s", time, implantId)
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

// Should only be called once. If there isn't an existing keypair, this function can be used to add a keypair to the db.
func InitKeypair(private string, public string) error {
	log.Printf("[INFO] db: executing statement: INSERT INTO secrets VALUES (?,?)")
	stmt1, err := dbConn.Prepare("INSERT INTO secrets VALUES (?,?)")
	if err != nil {
		return err
	}
	stmt2, err := dbConn.Prepare("INSERT INTO secrets VALUES (?,?)")
	if err != nil {
		return err
	}
	defer stmt1.Close()
	defer stmt2.Close()

	_, err = stmt1.Exec("private_key", private)
	if err != nil {
		return err
	}

	_, err = stmt2.Exec("public_key", public)
	if err != nil {
		return err
	}

	return nil
}

// Fetches the private key used for encryption/decryption of implant comms.
func GetPrivateKey() (string, error) {
	log.Printf("[INFO] db: executing statement: SELECT value FROM secrets WHERE key=private_key")
	stmt, err := dbConn.Prepare("SELECT value FROM secrets WHERE key=?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	rows, err := stmt.Query("private_key")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	hasResult := rows.Next()
	if !hasResult {
		return "", ErrNoResults
	}

	var privkey string

	err = rows.Scan(&privkey)
	if err != nil {
		return "", err
	}

	return privkey, nil
}

// Fetches the public key that is used by implants to encrypt session keys.
func GetPublicKey() (string, error) {
	log.Printf("[INFO] db: executing statement: SELECT value FROM secrets WHERE key=public_key")
	stmt, err := dbConn.Prepare("SELECT value FROM secrets WHERE key=?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	rows, err := stmt.Query("public_key")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	hasResult := rows.Next()
	if !hasResult {
		return "", ErrNoResults
	}
	var publicKey string

	err = rows.Scan(&publicKey)
	if err != nil {
		return "", err
	}

	return publicKey, nil
}
