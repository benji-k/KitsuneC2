package db

type Implant_info struct {
	Id           string
	Name         string
	Public_ip    string
	Os           string
	Arch         string
	Last_checkin int64
	Username     string
	Uid          string
	Gid          string
	Hostname     string
	Active       bool
}

type Implant_task struct {
	Task_id     string
	Implant_id  string
	Task_type   int
	Task_data   []byte
	Completed   bool
	Task_result []byte
}
