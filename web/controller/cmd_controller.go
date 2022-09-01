package controller

import (
	"fgv50/err"
	"fgv50/flag"
	"fgv50/scanner"
	"fgv50/tools"
	"fgv50/tools/storage"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)


type (
	Resp struct {
		ResList []tools.Result `json:"res_list"`

	}
	
	Cmd struct {
		CmdType string `json:"cmd_type"`
		Url string `json:"url"`
		Ip string `json:"ip"`
		Port int `json:"port"`
		NetworkSegment string `json:"network_segment"`
	}

)
const (
	SingleIP = "i"
	NetworkSegment = "ns"
	URL = "url"
)

func CommandExec(c *gin.Context) {
	hdb, ok := c.Get("histDB")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.ErrLevelDBInit.Error(),
			"data": nil,
		})
		return
	}
	var cmd Cmd
	if err1 := c.ShouldBind(&cmd); err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err1.Error(),
			"data": nil,
		})
		return
	} else {
		switch cmd.CmdType {
		case SingleIP:
			args, err2 := flag.NewArgs("", "", cmd.Ip, cmd.Port, 3)
			if err2 != nil {
				c.JSON(http.StatusOK, gin.H{
					"msg": err2.Error(),
					"data": nil,
				})
				return
			} else {
				args.HistDB = hdb.(*storage.HisDB)
				kTime := scanner.RunCli(args)
				
				// read results and vuls from histDB
				vRes, vVul := args.HistDB.GetRecord(kTime)
				
				c.String(http.StatusOK, fmt.Sprintf("%s@%s", string(vRes), string(vVul)))
				return

			}

		case NetworkSegment:
		case URL:
		default:
			c.JSON(http.StatusOK, gin.H{
				"msg": err.ErrUnknownCmd.Error(),
				"data": nil,
			})
			return
		}
	}
}