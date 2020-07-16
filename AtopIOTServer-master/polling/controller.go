package polling

import (
	"time"
)

// PollingJob polling jobs
type PollingJob struct {
	Url        string
	Collection string
}

// PollingController controller which start polling jobs
type PollingController struct {
	ticker *time.Ticker
	Jobs   []PollingJob
}

// func NewPollingController(sec int64) *PollingController {
// 	pc := &PollingController{
// 		ticker: time.NewTicker(time.Second * sec),
// 		Jobs:   make([]PollingJob, 10),
// 	}
// 	return pc
// }
