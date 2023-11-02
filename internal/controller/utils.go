package controller

func getPortOrDefault(port int32, defaultVal int) int32 {
	if port == 0 {
		return int32(defaultVal)
	} else {
		return port
	}
}
