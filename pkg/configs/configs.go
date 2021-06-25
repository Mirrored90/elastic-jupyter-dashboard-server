package configs

func GetServerConfig() (serverConfig map[string]string) {
	serverConfig = make(map[string]string)

	serverConfig["HOST"] = "0.0.0.0"
	serverConfig["PORT"] = "9090"
	serverConfig["ENV"] = "debug"
	return
}
