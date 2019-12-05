package jobs

import (
	"dc-wrapper-api/channels"
	"dc-wrapper-api/databases"
	"dc-wrapper-api/services/request_out_download_center"
	"dc-wrapper-api/structs/models"
	"dc-wrapper-api/structs/models/dc_request_logics"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// AddDcGetResultChan 往DcGetResultChan中添加dcRequest
// 不支持多个Go程
func AddDcGetResultChan() {
	mgoSession := databases.MgoSession.Copy()
	defer mgoSession.Close()
	co := mgoSession.DB("fxt").C("fxt_requests")

	count, err := co.Find(bson.M{"status": bson.M{"$in": []int{dc_request_logics.STATUS_查询中}}}).Count()
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Minute)
		return
	}

	if count == 0 {
		fmt.Println("无查询中的任务，等1分钟")
		time.Sleep(time.Minute)
		return
	}

	iter := co.Find(bson.M{"status": bson.M{"$in": []int{dc_request_logics.STATUS_查询中}}}).Sort("updated_at").Iter()
	for {
		dcRequest := models.DcRequest{}
		if iter.Next(&dcRequest) == false {
			break
		}
		if dcRequest.UniqueMD5 == "" {
			dcRequest.UniqueMD5 = request_out_download_center.UniqueMd5(dcRequest)
		}
		channels.DcGetResultChan <- &dcRequest
	}
	time.Sleep(time.Second * 5)
}
