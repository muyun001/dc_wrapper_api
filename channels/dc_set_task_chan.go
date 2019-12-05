package channels

import "dc-wrapper-api/structs/models"

var DcSetTaskChan chan *models.DcRequest

func init() {
	DcSetTaskChan = make(chan *models.DcRequest, 2)
}
