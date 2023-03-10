package server

import (
	"github.com/gin-gonic/gin"
	"hash"
	"net/http"
)

type deduper struct {
	h  hash.Hash32
	h1 uint32
	h2 uint32
}

func (server *HttpServer) registerVideo() {
	server.instance.POST("/video", func(c *gin.Context) {
		c.Header("Content-Type", "multipart/x-mixed-replace;boundary=endofsection")
		c.Status(http.StatusOK)
	})
}
