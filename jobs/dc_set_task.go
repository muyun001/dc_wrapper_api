package jobs

import (
	"dc-wrapper-api/channels"
	"dc-wrapper-api/databases"
	"dc-wrapper-api/services/request_out_download_center"
	"dc-wrapper-api/structs/models/dc_request_logics"
	"gopkg.in/mgo.v2/bson"
)

// DcSetTask 发送任务到下载中心
// 支持多个Go程
func DcSetTask() {
	mgoSession := databases.MgoSession.Copy()
	defer mgoSession.Close()
	co := mgoSession.DB("fxt").C("fxt_requests")

	dcRequest := <-channels.DcSetTaskChan

	selector := bson.M{"unique_key": dcRequest.UniqueKey}

	uniqueMd5, err := request_out_download_center.ApiSetTask(*dcRequest)

	if err != nil {
		updateData := bson.M{"$set": bson.M{"status": dc_request_logics.STATUS_发送失败}}
		_ = co.Update(selector, updateData)
		request_out_download_center.ApiResetIp()

		return
	}

	updateData := bson.M{"$set": bson.M{"status": dc_request_logics.STATUS_查询中, "unique_md5": uniqueMd5}}
	_ = co.Update(selector, updateData)
}
