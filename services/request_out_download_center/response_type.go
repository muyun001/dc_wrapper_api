package request_out_download_center

import "dc-wrapper-api/utils/strings_util"

// ResponseType: 获取返回值类型的对应序号
// 未匹配时默认为只拿body，此场景常用
func ResponseType(responseTypes []string) int {
	switch len(responseTypes) {
	case 1:
		if responseTypes[0] == "body" {
			return RESPONSE_TYPE_BODY
		}

	case 2:
		if strings_util.InSlice(responseTypes, "body") && strings_util.InSlice(responseTypes, "header") {
			return RESPONSE_TYPE_BODY_HEADER
		}
		if strings_util.InSlice(responseTypes, "body") && strings_util.InSlice(responseTypes, "capture") {
			return RESPONSE_TYPE_BODY_CAPTURE
		}
	}

	return RESPONSE_TYPE_BODY
}
