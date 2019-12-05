package channels

import "dc-wrapper-api/structs/models"

var DcGetResultChan chan *models.DcRequest

func init() {
	DcGetResultChan = make(chan *models.DcRequest, 5)
}
