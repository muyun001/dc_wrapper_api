package request_out_download_center

// PriorityNum: 获取优先级的对应序号
func PriorityNum(priority string) int {
	switch priority {
	case "low":
		return PRIORITY_LOW
	case "normal":
		return PRIORITY_NORMAL
	case "high":
		return PRIORITY_HIGH
	default:
		return PRIORITY_NORMAL
	}
}
