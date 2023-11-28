package db

type Implant_info struct {
	Id           string
	Name         string
	Public_ip    string
	Session_key  string
	Os           string
	Arch         string
	Last_checkin int
	Username     string
	Uid          string
	Gid          string
	Hostname     string
}

type Implant_task struct {
	Task_id    string
	Implant_id string
	Task_type  string
	Task_data  string
}
