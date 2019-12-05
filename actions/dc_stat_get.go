package actions

import (
	"dc-wrapper-api/databases"
	"dc-wrapper-api/databases/entities"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TypeStatusCount struct {
	Type   int `json:"type"`
	Status int `json:"status"`
	Count  int `json:"count"`
}

type StatusCount struct {
	UnQueried int `json:"un_queried"`
	Querying  int `json:"querying"`
	Queried   int `json:"queried"`
	Failed    int `json:"failed"`
}

// DcStatGet 下载中心任务统计
func DcStatGet(c *gin.Context) {

	typeStatusCounts := make([]TypeStatusCount, 0)
	databases.Db.Model(entities.DcUrl{}).Group("type, status").Select("type, status, count(*) as count").Scan(&typeStatusCounts)

	countMapHtmlOnly := StatusCount{}
	countMapHtmlWithHeader := StatusCount{}
	countMapHtmlWithRender := StatusCount{}
	countMapHtmlWithCapture := StatusCount{}
	for _, typeStatusCount := range typeStatusCounts {
		switch typeStatusCount.Type {
		case entities.DCURL_TYPE_抓HTML不带头:
			setCount(&countMapHtmlOnly, typeStatusCount)
		case entities.DCURL_TYPE_抓HTML带头:
			setCount(&countMapHtmlWithHeader, typeStatusCount)
		case entities.DCURL_TYPE_抓HTML渲染:
			setCount(&countMapHtmlWithRender, typeStatusCount)
		case entities.DCURL_TYPE_抓HTML截图:
			setCount(&countMapHtmlWithCapture, typeStatusCount)
		}
	}

	countMap := make(map[string]StatusCount)
	countMap["html_only"] = countMapHtmlOnly
	countMap["html_with_header"] = countMapHtmlWithHeader
	countMap["html_with_render"] = countMapHtmlWithRender
	countMap["html_with_capture"] = countMapHtmlWithCapture

	c.JSON(http.StatusOK, countMap)
}

func setCount(statusCount *StatusCount, typeStatusCount TypeStatusCount) {
	switch typeStatusCount.Status {
	case entities.DCURL_STATUS_未抓取:
		statusCount.UnQueried = typeStatusCount.Count
	case entities.DCURL_STATUS_抓取中:
		statusCount.Querying = typeStatusCount.Count
	case entities.DCURL_STATUS_已抓取:
		statusCount.Queried = typeStatusCount.Count
	case entities.DCURL_STATUS_抓取失败:
		statusCount.Failed = typeStatusCount.Count
	}
}
