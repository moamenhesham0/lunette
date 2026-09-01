package main

import "github.com/gin-gonic/gin"

func main() {
	gateway := gin.Default()
	InitGateway(gateway)
}
