package main

func main() {
	config := LoadMetricsServerConfig()
	server := NewMetricsServer(config)
	server.Serve()
}
