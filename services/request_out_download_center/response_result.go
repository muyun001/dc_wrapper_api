package request_out_download_center

import (
	"encoding/base64"
	"strings"
)

// splitHtmlCaptureResult 分割包含html和capture的result字符串
func splitHtmlCaptureResult(b64Result string) (html, capture string, err error) {
	resultBytes, err := base64.StdEncoding.DecodeString(b64Result)
	if err != nil {
		return html, capture, err
	}

	result := string(resultBytes)
	sep := "||||"
	if strings.Contains(result, sep) {
		splits := strings.SplitN(string(resultBytes), sep, 2)
		capture = splits[0]
		html = splits[1]
	} else {
		html = result
	}

	return html, capture, nil
}

// SplitResultFromRData 从RData分割Result内容
func SplitResultFromRData(rData RData) (header, body, capture string, err error) {
	body, capture, err = splitHtmlCaptureResult(rData.Result)
	return
}
