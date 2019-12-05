package entities

const (
	DCURL_TYPE_抓HTML不带头 = 1
	DCURL_TYPE_抓HTML带头  = 2
	DCURL_TYPE_抓HTML渲染  = 3
	DCURL_TYPE_抓HTML截图  = 4
)

const (
	DCURL_STATUS_未抓取  = 0
	DCURL_STATUS_抓取中  = 1
	DCURL_STATUS_已抓取  = 2
	DCURL_STATUS_抓取失败 = 3
)
