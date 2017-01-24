package store

//SecGroupInfo contains the information about what CF Security groups are
// associated with which service instances
type SecGroupInfo struct {
	ServiceInstanceGUID string
	SecGroupName        string
}
