package main

import "github.com/IgorKravtsov/esport_server_go/internal/app"

const configsDir = "configs"

func main() {
	app.Run(configsDir)
}
