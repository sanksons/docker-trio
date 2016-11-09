package versionmanager

import (
	"testing"
)

/*
Test Versionable interface implementation
*/
type testVersionableImpl struct {
}

func (o testVersionableImpl) GetInstance() interface{} {
	return o
}

/*
Test version manager
*/
func TestVersionManager(t *testing.T) {
	resource := "TEST_RESOURCE"
	version := "TEST_VERSION"
	action := "TEST_ACTION"
	bucketId := "TEST_BUCKET_ID"

	vmap := VersionMap{
		Version{
			Resource: resource,
			Version:  version,
			Action:   action,
			BucketId: bucketId,
		}: *new(testVersionableImpl),
	}

	Initialize(vmap)

	versionableInstance, verr := Get(resource, version, action, bucketId)
	if verr != nil {
		t.Error("Failed to get versionable from version manager")
	}

	_, ok := versionableInstance.(testVersionableImpl)

	if !ok {
		t.Error("Returned versionable instance mismatch")
	}
}
