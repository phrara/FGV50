package web

import (
	"fgv50/tools/storage"
	"fgv50/web/controller"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)


func StartUpWeb(histDB *storage.HisDB) {

	r := gin.Default()
	r.Static("/", "./web/static")
	r.Use(Cors())
	initRouter(r, histDB)

	go func ()  {
		addr := "Access address is http://localhost:8989/dingzhen.html"
		time.Sleep(time.Second * 3)
		fmt.Printf("\x1b[1;34mHost: %s\x1b[0m\n", addr)
	}()

	r.Run(":8989")
}



func initRouter(r *gin.Engine, histDB *storage.HisDB) {
	// r.GET("/indexx",func(c *gin.Context) {
	// 	c.Request.URL.Path = "/dingzhen.html"
	// 	r.HandleContext(c)
	// })
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

// 解决跨域问题
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}