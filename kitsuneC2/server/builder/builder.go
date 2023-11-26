package builder

type BuilderConfig struct {
	C2ServerIp              string
	C2ServerPort            int
	ImplantName             string
	ImplantCallbackInterval int
	ImplantCallbackJitter   int
}
