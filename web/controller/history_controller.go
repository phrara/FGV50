package controller

import (
	"fgv50/err"
	"fgv50/tools/storage"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type Query struct {
	Time string `json:"time"`
}


func HistoryQuery(c *gin.Context) {
	hdb, ok := c.Get("histDB")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.ErrLevelDBInit.Error(),
			"data": nil,
		})
		return
	}
	var q Query
	if err1 := c.ShouldBind(&q); err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err1.Error(),
			"data": nil,
		})
		return
	} else {
		println("+++" , q.Time)
		histDB := hdb.(*storage.HisDB)
		r, _ := regexp.Compile(`^((([0-9]{3}[1-9]|[0-9]{2}[1-9][0-9]{1}|[0-9]{1}[1-9][0-9]{2}|[1-9][0-9]{3})-(((0[13578]|1[02])-(0[1-9]|[12][0-9]|3[01]))|((0[469]|11)-(0[1-9]|[12][0-9]|30))|(02-(0[1-9]|[1][0-9]|2[0-8]))))|((([0-9]{2})(0[48]|[2468][048]|[13579][26])|((0[48]|[2468][048]|[3579][26])00))-02-29))\s+([0-1]?[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])$`)
		if r.MatchString(q.Time) {
			vRes, vVul := histDB.GetRecord([]byte(q.Time))
			resp := formatJson(q.Time, vRes, vVul)
			c.String(http.StatusOK, resp)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.ErrIllFormedTime.Error(),
				"data": nil,
			})
		}
	}

}