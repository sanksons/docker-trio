package versionmanager

/*
Data structure to store the API that is versioned
*/
type Version struct {
	Resource string
	Version  string
	Action   string
	BucketId string
}

type VersionMap map[Version]Versionable
