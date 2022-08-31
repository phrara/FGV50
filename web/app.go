package web

import (
	"fgv50/web/controller"

	"github.com/gin-gonic/gin"
)


func StartUpWeb() {

	r := gin.Default()
	initRouter(r)

	r.Run()
}



func initRouter(r *gin.Engine) {
	r.GET("/",controller.Index)
	r.POST("/cmd", controller.CommandExec)
}