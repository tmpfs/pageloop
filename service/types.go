package service

// A helper type for service methods to declare as the
// reply (second) argument. When this type is used the
// Result assigned to the ServiceReply is used as the
// Reply on the Response pbject.
//
// This enables easily passing references to existing objects
// in service methods and returning to actual Result object
// rather than a wrapper with a struct field name.
type ServiceReply struct {
  Status int
  Reply interface{}
}
