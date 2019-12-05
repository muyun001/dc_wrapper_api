package jobs

import (
	"dc-wrapper-api/channels"
	"dc-wrapper-api/databases"
	"dc-wrapper-api/databases/entities"
	"dc-wrapper-api/services/request_out_download_center"
	"dc-wrapper-api/structs/models"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// DcGetResult 从下载中心获取查询任务结果并处理
// 支持多个Go程
func DcGetResult() {
	mgoSession := databases.MgoSession.Copy()
	defer mgoSession.Close()
	co := mgoSession.DB("fxt").C("fxt_requests")

	dcRequest := <-channels.DcGetResultChan

	dcRequestMap := map[string]models.DcRequest{dcRequest.UniqueMD5: *dcRequest}
	dcResponseBytes, err := request_out_download_center.ApiGetResults(dcRequest.Config.Priority, dcRequestMap)
	if err != nil {
		request_out_download_center.ApiResetIp()
		time.Sleep(time.Second * 5)
		return
	}

	if dcResponseBytes == nil {
		return
	}

	dcGetResultResponse := &request_out_download_center.DcApiResponse{}
	err = json.Unmarshal(dcResponseBytes, &dcGetResultResponse)
	if err != nil {
		fmt.Println(err)
		return
	}

	rDataMap, err := request_out_download_center.ResponseRDataMap(dcGetResultResponse.RData)
	if err != nil {
		fmt.Println(err)
		return
	}

	for uniqueMd5 := range rDataMap {
		uniqueKey := dcRequestMap[uniqueMd5].UniqueKey

		if rDataMap[uniqueMd5].Status == request_out_download_center.RESPONSE_STATUS_抓取失败 {
			fmt.Println("返回状态为3，重置状态")
			_ = co.Update(
				bson.M{"unique_key": uniqueKey},
				bson.M{"$set": bson.M{"status": entities.DCURL_STATUS_未抓取, "updated_at": time.Now()}},
			)
			continue
		}

		header, body, capture, err := request_out_download_center.SplitResultFromRData(rDataMap[uniqueMd5])
		if err != nil {
			fmt.Println(err)
			_ = co.Update(
				bson.M{"unique_key": uniqueKey},
				bson.M{"$set": bson.M{"status": entities.DCURL_STATUS_未抓取, "updated_at": time.Now()}},
			)
			continue
		}

		if body == "" {
			if time.Now().Sub(dcRequestMap[uniqueMd5].UpdatedAt) > time.Minute*30 {
				_ = co.Update(
					bson.M{"unique_key": uniqueKey},
					bson.M{"$set": bson.M{"status": entities.DCURL_STATUS_未抓取, "updated_at": time.Now()}},
				)
			}
			continue
		}

		redisSession := databases.RedisSession
		jsonBytes, _ := json.Marshal(models.DcResponse{
			Header:  header,
			Body:    body,
			Capture: capture,
		})
		redisSession.Do("set", uniqueKey, jsonBytes)
		redisSession.Expire(uniqueKey, time.Hour)

		_ = co.Update(
			bson.M{"unique_key": uniqueKey},
			bson.M{"$set": bson.M{"status": entities.DCURL_STATUS_已抓取, "updated_at": time.Now()}},
		)
	}
}
