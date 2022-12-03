// package models contains all database model data for this service
// including off the wire representations of on-wire messages.
package models

// Model interface is a WIP idea behind being able to include custom models that
// implement this interface.
type Model interface {
	ToProto() any
	FromProto(any) any
}
