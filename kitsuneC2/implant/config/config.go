package config

const (
	ImplantName           string = "BabyImplant"
	MaxRegisterRetryCount int    = 50
	PublicKey             string = "MIIBCgKCAQEA6DytLw66kgK+M4AljiutUpYygmKCo7nLevDRA0Oa4myxcQHIJRM09ARigqx4nlK9Tah4Czf4UvWTXzkpwYmLX8LPuQEks059hqRQuegieJ6UqPBFymuwPhy8P4Ml59tUlSAcXcgppWu0eeaHrnR06PAx0Ae2omad0/95qVZpWIfasBVOVWLz+2T+S7PErlaspXAd7QVO/eEfeDU4WSGbt/VhiiMTg2oKT6NSCC3GcfLCPGgPr9jf1/KIBeER39j1KyAe2Rji4+oMltSmPujn70tiiR2YdtfvxJI6bCMKHhauqK6d4Ps4rK2JQ9Ht+7KF2c7Cb5MCHymGz8e8+eQGhwIDAQAB"
)

var (
	ServerIp         string = "127.0.0.1"
	ServerPort       int    = 4444
	CallbackInterval int    = 10
	CallbackJitter   int    = 2
)
