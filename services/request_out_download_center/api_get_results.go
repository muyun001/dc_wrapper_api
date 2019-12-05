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

const GET_DC_GET_RESULT = "/download/getResult"

// sendGetResult 从下载中心获取查询结果(网络请求)
func sendGetResult(formData url.Values) ([]byte, error) {
	apiUrl := settings.DcApi + GET_DC_GET_RESULT

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

// ApiGetResult 从下载中心获取多个查询结果
func ApiGetResults(priority string, dcRequestMap map[string]models.DcRequest) ([]byte, error) {
	if len(dcRequestMap) == 0 {
		return nil, errors.New("dcRequests为空")
	}

	dcConfig := DcConfig{
		Redirect: 0,
		Priority: PriorityNum(priority),
	}
	configJsonBytes, _ := json.Marshal(dcConfig)

	dcUrls := make([]DcUrl, 0)
	for uniqueMd5 := range dcRequestMap {
		dcUrls = append(dcUrls, DcUrl{
			Url:       dcRequestMap[uniqueMd5].Request.Url,
			Type:      ResponseType(dcRequestMap[uniqueMd5].Config.ResponseTypes),
			UniqueKey: dcRequestMap[uniqueMd5].UniqueKey,
			UniqueMd5: uniqueMd5,
		})
	}
	urlsJsonBytes, _ := json.Marshal(dcUrls)

	dcGetResultRequest := &DcGetResultRequest{
		UserID: settings.DcUserId,
		Config: string(configJsonBytes),
		Urls:   string(urlsJsonBytes),
	}

	formData := url.Values{}
	err := json_util.StructToFormData(dcGetResultRequest, &formData)
	if err != nil {
		return nil, err
	}

	body, err := sendGetResult(formData)
	if err != nil {
		time.Sleep(time.Second * 2)
		body, err = sendGetResult(formData)
		if err != nil {
			return nil, err
		}
	}

	if !strings.Contains(string(body), "result") {
		return nil, nil
	}

	return body, nil
}
