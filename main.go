package main

import (
	"dc-wrapper-api/actions"
	"dc-wrapper-api/databases"
	"dc-wrapper-api/jobs"
	"dc-wrapper-api/services/request_out_download_center"
	"github.com/gin-gonic/gin"
)

func foreverGo(run func(), routineLimits int) {
	for i := 0; i < routineLimits; i++ {
		go func() {
			for {
				run()
			}
		}()
	}
}

func main() {
	defer databases.MgoSession.Close()
	defer databases.RedisSession.Close()

	request_out_download_center.ApiResetIp()

	foreverGo(jobs.AddDcSetTask, 1)
	foreverGo(jobs.DcSetTask, 1)
	foreverGo(jobs.AddDcGetResultChan, 1)
	foreverGo(jobs.DcGetResult, 1)

	r := gin.Default()
	r.PUT("/requests/:unique-key", actions.RequestPut)
	r.PUT("/request-reset/:unique-key", actions.RequestResetPut)
	r.GET("/requests/:unique-key", actions.RequestGet)
	r.GET("/responses/:unique-key", actions.ResponseGet)
	r.GET("/requests", actions.RequestsGet)
	r.GET("/requests-stat", actions.RequestsStatGet)
	r.POST("/responses-check", actions.ResponsesCheckPost)
	r.GET("/dc-stat", actions.DcStatGet)
	r.PUT("/dc-set-task/on", actions.DcStartSetTaskPut)
	r.PUT("/dc-set-task/off", actions.DcStopSetTaskPut)

	r.Run(":8090") // listen and serve on 0 .0.0.0:8080
}
