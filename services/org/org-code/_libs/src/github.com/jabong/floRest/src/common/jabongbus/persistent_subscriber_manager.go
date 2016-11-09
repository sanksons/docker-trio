package jabongbus

import ()

// data structure to be set in redis by each alive subscriber periodically
type statusData struct {
	Id   int64 `json:"id"`
	Time int64 `json:"time"`
}

// GetPersistentSubscriber get persistent subscriber implementation
func GetPersistentSubscriber(conf *Subscriberconfig) (PeristentSubscriber, error) {
	var err error
	if conf.Cluster {
		obj := new(persistSubCluster)
		err = obj.init(conf)
		return obj, err
	} else {
		obj := new(persistSubNonCluster)
		err = obj.init(conf)
		return obj, err
	}
}
