package request_out_download_center

import (
	"dc-wrapper-api/settings"
	"dc-wrapper-api/structs/models"
	"dc-wrapper-api/utils/json_util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/panwenbin/ghttpclient"
	"github.com/panwenbin/ghttpclient/header"
	"net/url"
	"strings"
	"time"
)

const POST_SET_TASK = "/download/setTask"

// sendSetTask 发送任务到下载中心(网络请求)
func sendSetTask(formData url.Values) ([]byte, error) {
	apiUrl := settings.DcApi + POST_SET_TASK

	client := ghttpclient.NewClient().Timeout(time.Second * 30).Url(apiUrl).
		Headers(nil).Body(strings.NewReader(formData.Encode())).
		ContentType(header.CONTENT_TYPE_FORM_URLENCODED).Post()

	res, err := client.Response()
	if err != nil {
		return nil, err
	}

	body, _ := client.ReadBodyClose()

	if res.StatusCode > 200 {
		return body, errors.New(fmt.Sprintf("unexpect status: %d", res.StatusCode))
	}

	return body, nil
}

// ApiSetTask 发送任务到下载中心
// 包含一次重试和错误判断
func ApiSetTask(dcRequest models.DcRequest) (string, error) {
	dcSetTaskRequest := &DcSetTaskRequest{
		UserID:  settings.DcUserId,
		Headers: fmt.Sprintf(`{"User-Agent": "%s", "Cookie": "%s"}`, dcRequest.Request.UserAgent, dcRequest.Request.Cookie),
		Config:  fmt.Sprintf(`{"redirect": 0, "priority": %d}`, PriorityNum(dcRequest.Config.Priority)),
		Urls:    fmt.Sprintf(`[{"url": "%s", "type": %d, "unique_key": "%s"}]`, dcRequest.Request.Url, ResponseType(dcRequest.Config.ResponseTypes), dcRequest.UniqueKey),
	}

	var dcSetTaskRequestMap map[string]string
	_ = json_util.StructToStringMap(dcSetTaskRequest, &dcSetTaskRequestMap)

	formData := url.Values{}
	for key, value := range dcSetTaskRequestMap {
		formData.Set(key, value)
	}

	body, err := sendSetTask(formData)
	if err != nil {
		time.Sleep(time.Second * 2)
		body, err = sendSetTask(formData)
		if err != nil {
			return "", errors.New("request failed")
		}
	}

	dcSetTaskResponse := &DcApiResponse{}
	_ = json.Unmarshal(body, &dcSetTaskResponse)

	if strings.Contains(string(body), "params error") {
		return "", errors.New("params error")
	}

	uniqueMd5 := ResponseUniqueMd5(dcSetTaskResponse.RData)

	if strings.Contains(string(body), "task insert error") {
		fmt.Println("重复插入数据，已处理")
		return uniqueMd5, nil
	}

	return uniqueMd5, nil
}
