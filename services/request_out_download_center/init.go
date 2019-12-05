package request_out_download_center

import "os"

var IS_SETTING_TASK bool

func init() {
	settingTask := os.Getenv("SETTING_TASK_ON")
	if settingTask != "" && settingTask != "false" && settingTask != "0" {
		IS_SETTING_TASK = true
	}
}
