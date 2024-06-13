package utils

import (
	"sync"
	"time"

	"github.com/olahol/melody"
)

var M *melody.Melody

var MessageCache *sync.Map

func InitWS() {
	M = melody.New()
	MessageCache = new(sync.Map)
	M.Config = &melody.Config{
		WriteWait:         10 * time.Second,
		PongWait:          60 * time.Second,
		PingPeriod:        (60 * time.Second * 9) / 10,
		MessageBufferSize: 256,
		MaxMessageSize:    1024 * 1024,
	}
}
