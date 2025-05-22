package main

import (
	service "forum/internal/service"
)

func main() {
	service.InitDependencies()
	service.StartServer()
}
