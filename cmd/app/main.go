package main

import (
	"git.codenrock.com/tender/internal/app"
)

const configPath = "config"

func main() {
	app.Run(configPath)
}
