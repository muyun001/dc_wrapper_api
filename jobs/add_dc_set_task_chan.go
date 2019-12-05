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

// AddDcSetTask 往DcSetTaskChan添加dcRequest
// 不支持多个Go程
func AddDcSetTask() {
	if !request_out_download_center.IS_SETTING_TASK {
		fmt.Println("未开启任务，等待5秒")
		time.Sleep(time.Second * 5)
		return
	}

	mgoSession := databases.MgoSession.Copy()
	defer mgoSession.Close()
	co := mgoSession.DB("fxt").C("fxt_requests")

	count, err := co.Find(bson.M{"status": bson.M{"$in": []int{dc_request_logics.STATUS_未查询, dc_request_logics.STATUS_发送失败}}}).Count()

	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Second * 5)
		return
	}

	if count == 0 {
		fmt.Println("未获取到任务，等待5秒")
		time.Sleep(time.Second * 5)
		return
	}

	iter := co.Find(bson.M{"status": bson.M{"$in": []int{dc_request_logics.STATUS_未查询, dc_request_logics.STATUS_发送失败}}}).Sort("updated_at").Iter()

	for {
		dcRequest := models.DcRequest{}
		if iter.Next(&dcRequest) == false {
			break
		}
		channels.DcSetTaskChan <- &dcRequest
	}

	time.Sleep(time.Second)
}
