package actions

import (
	"dc-wrapper-api/databases"
	"dc-wrapper-api/structs/models/dc_request_logics"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

// RequestsStatGet 输出所有请求中各状态的数量
func RequestsStatGet(c *gin.Context) {
	mgoSession := databases.MgoSession.Copy()
	defer mgoSession.Close()
	co := mgoSession.DB("fxt").C("fxt_requests")

	unQueriedCount, _ := co.Find(bson.M{"status": dc_request_logics.STATUS_未查询}).Count()
	queryingCount, _ := co.Find(bson.M{"status": dc_request_logics.STATUS_查询中}).Count()
	queriedCount, _ := co.Find(bson.M{"status": dc_request_logics.STATUS_已查询}).Count()
	readCount, _ := co.Find(bson.M{"status": dc_request_logics.STATUS_已读取}).Count()

	c.JSON(http.StatusOK, gin.H{
		"un_queried_count": unQueriedCount,
		"querying_count":   queryingCount,
		"queried_count":    queriedCount,
		"read_count":       readCount,
	})
}
