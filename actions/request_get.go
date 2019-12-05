package actions

import (
	"crypto/md5"
	"dc-wrapper-api/databases"
	"dc-wrapper-api/databases/entities"
	"dc-wrapper-api/structs/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

// 查看单个请求任务
func RequestGet(c *gin.Context) {
	uniqueKey := c.Param("unique-key")

	mgoSession := databases.MgoSession.Copy()
	defer mgoSession.Close()
	co := mgoSession.DB("fxt").C("fxt_requests")

	dcRequest := models.DcRequest{}
	err := co.Find(bson.M{"unique_key": uniqueKey}).One(&dcRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	dcMd5 := fmt.Sprintf("%x", md5.Sum([]byte(dcRequest.Request.Url+dcRequest.UniqueKey)))
	dcUrl := &entities.DcUrl{}
	if databases.Db.Where("md5 = ?", dcMd5).First(dcUrl).RecordNotFound() {
		dcUrl = nil
	}

	c.JSON(http.StatusOK, gin.H{
		"dcRequest": dcRequest,
		"dcUrl":     dcUrl,
	})
}
