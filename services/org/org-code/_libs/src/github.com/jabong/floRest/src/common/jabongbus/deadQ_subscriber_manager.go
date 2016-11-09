package jabongbus

import ()

// GetPersistentSubscriber get persistent subscriber implementation
func GetDeadQSubscriber(conf *Subscriberconfig) (DeadQSubscriber, error) {
	var err error
	if !conf.Persistent {
		return nil, ErrNoDeadQForNonPers
	}
	if conf.Cluster {
		obj := new(deadQSubCluster)
		err = obj.init(conf)
		return obj, err
	} else {
		obj := new(deadQSubNonCluster)
		err = obj.init(conf)
		return obj, err
	}
}
