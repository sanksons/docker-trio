package service

import (
	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
)

func getServiceVersion(data workflow.WorkFlowData) (resource string,
	version string, action string, orchBucket string) {
	resourceData, _ := data.IOData.Get(constants.RESOURCE)
	versionData, _ := data.IOData.Get(constants.VERSION)
	actionData, _ := data.IOData.Get(constants.ACTION)
	bucketsMap, _ := data.ExecContext.GetBuckets()

	resource = ""
	if v, ok := resourceData.(string); ok {
		resource = v
	}

	version = ""
	if v, ok := versionData.(string); ok {
		version = v
	}

	action = ""
	if v, ok := actionData.(string); ok {
		action = v
	}

	orchBucket = getBucketValue(resource, bucketsMap)
	return resource, version, action, orchBucket
}

func getBucketValue(resource string, bucketsMap map[string]string) string {
	orchBucket := constants.ORCHESTRATOR_BUCKET_DEFAULT_VALUE
	resourceBucketKey, present := resourceBucketMapping[resource]
	if !present {
		return orchBucket
	}

	orchbucketId, found := bucketsMap[resourceBucketKey]
	if found {
		orchBucket = orchbucketId
	}
	return orchBucket
}
