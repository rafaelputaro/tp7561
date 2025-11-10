package main

const MESSAGE_START = "Starting client..."

func main() {
	client, err := NewClient()
	// Corre el cliente subiendo archivos, descargando archivos y chequeando su llegada.
	// Respalda métricas
	if err == nil {
		client.Start()
	}
}
