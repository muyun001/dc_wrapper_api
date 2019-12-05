package actions

import (
	"dc-wrapper-api/databases"
	"dc-wrapper-api/structs/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"strconv"
)

// 请求任务列表
func RequestsGet(c *gin.Context) {
	mgoSession := databases.MgoSession.Copy()
	defer mgoSession.Close()
	co := mgoSession.DB("fxt").C("fxt_requests")

	mgoFilter := bson.M{}
	queryStatus:= c.Query("status")
	if queryStatus != "" {
		status, err := strconv.Atoi(queryStatus)
		if err != nil {
			log.Fatal(err)
		}
		mgoFilter = bson.M{"status": status}
	}


	dcRequests := make([]models.DcRequest, 0)
	dcRequest := models.DcRequest{}
	iter := co.Find(mgoFilter).Iter()
	for iter.Next(&dcRequest) {
		dcRequests = append(dcRequests, dcRequest)
	}

	c.JSON(http.StatusOK, dcRequests)
}
