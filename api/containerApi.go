package api

import (
	"cft/log"
	"cft/model"
	"cft/monitor"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddContainer(c *gin.Context) {
	id := c.PostForm("id")
	name := c.PostForm("name")
	host := c.PostForm("host")
	if id == "" || name == "" {
		log.Errorf("bad request id : %v, name : %v, host : %v", id, name, host)
		c.JSON(http.StatusBadRequest, gin.H{
			"id":   id,
			"name": name,
			"host": host,
		})
		return
	}
	monitor.AddContainer(id, name, host, model.StateMaster)
	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
