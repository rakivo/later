package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ADDR = "127.0.0.1:6969"
)

var (
	TrustedProxies = []string{
		"127.0.0.1",
	}
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.SetTrustedProxies(TrustedProxies)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello, world")
	})

	done := make(chan error, 1)
	defer close(done)

	go func() {
		err := r.Run(ADDR); if err != nil {
			fmt.Println("ERROR LAUNCHING SERVER: ", err)
			done <- err
		}
	}()

	err := <- done
	if err != nil {
		fmt.Println("STOPPING SERVER, ERROR: ", err)
	}
}
