package store

//SecGroupInfo contains the information about what CF Security groups are
// associated with which service instances
type SecGroupInfo struct {
	ServiceInstanceGUID string
	SecGroupName        string
}

//WithGUID returns a copy of the receiver SecGroupInfo, except that the
// ServiceInstanceGUID is set to the given string
func (s SecGroupInfo) WithGUID(guid string) SecGroupInfo {
	s.ServiceInstanceGUID = guid
	return s
}

//WithGroupName returns a copy of the receiver SecGroupInfo, except that the
// SecGroupName is set to the given string
func (s SecGroupInfo) WithGroupName(name string) SecGroupInfo {
	s.SecGroupName = name
	return s
}
