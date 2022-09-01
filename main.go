package main

import (
	"github.com/ivmello/fakehub/cmd/api"
)

// @title           Fake Hub
// @version         1.0
// @description     This is a API for get news from fact checking websites

// @contact.name   Igor Vieira de Mello
// @contact.email  ivmello@gmail.com

// @host      localhost:3000
// @BasePath  /api/v1
func main() {
	api.Initialize()
}
