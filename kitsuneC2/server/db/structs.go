package db

type Implant_info struct {
	Id           string
	Name         string
	Public_ip    string
	Session_key  string
	Os           string
	Arch         string
	Last_checkin int
}

type Implant_task struct {
	Task_id    string
	Implant_id string
	Task_type  string
	Task_data  string
}
