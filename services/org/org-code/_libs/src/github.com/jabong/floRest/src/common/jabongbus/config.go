package jabongbus

import ()

// PublisherConfig config for publisher
type PublisherConfig struct {
	Url string // e.g "http://localhost:8090/omnibus/"
}

// Subscriberconfig config for subscriber
type Subscriberconfig struct {
	Publisher  string
	RoutingKey string
	RedisCon   string
	Cluster    bool
	Persistent bool
}
