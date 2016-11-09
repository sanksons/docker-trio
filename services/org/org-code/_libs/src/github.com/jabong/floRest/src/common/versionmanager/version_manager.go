package versionmanager

import (
	"errors"
)

/*
Instance of the Version Manager
*/
type versionManager struct {
	//Mapping of the resource, version, action and Executable type
	mapping VersionMap
}

var vmgr *versionManager = nil

/*
Constructor for the version manager
*/
func Initialize(m VersionMap) {
	if vmgr != nil {
		return
	}
	vmgr = new(versionManager)
	vmgr.mapping = m
}

/*
Get the executable for the resource, version, action, bucketId
*/
func Get(resource string, version string, action string, bucketId string) (Versionable, error) {
	if vmgr == nil {
		return nil, errors.New("Version manager not initialized")
	}

	if vmgr.mapping == nil {
		return nil, errors.New("No version mapping present")
	}

	ver := Version{
		Resource: resource,
		Version:  version,
		Action:   action,
		BucketId: bucketId,
	}

	versionable, ok := vmgr.mapping[ver]
	if !ok {
		return nil, errors.New("Versionable not found in version manager")
	}
	return versionable, nil
}
