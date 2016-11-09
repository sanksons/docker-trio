package jabongbus

import (
	"time"
)

type Publisher interface {
	PublishMessage(*PubRequest, time.Duration, bool) (*PubResponse, error)
}
