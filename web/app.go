package web

import (
	"fgv50/tools/storage"
	"fgv50/web/controller"

	"github.com/gin-gonic/gin"
)


func StartUpWeb(histDB *storage.HisDB) {

	r := gin.Default()
	r.Use()
	initRouter(r, histDB)

	r.Run()
}



func initRouter(r *gin.Engine, histDB *storage.HisDB) {
	r.GET("/",controller.Index)
	r.POST("/cmd", SetHistDB(histDB), controller.CommandExec)
	r.POST("/hist", SetHistDB(histDB), controller.HistoryQuery)
}

// SetHistDB 是一个传递HistDB的中间件
func SetHistDB(histDB *storage.HisDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("histDB", histDB)
		c.Next()
	}
}