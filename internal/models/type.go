// package models contains all database model data for this service
// including off the wire representations of on-wire messages.
package models

import (
	pb "github.com/tehzwen/real_estate-service/proto/golang"
)

// ListingType is a system representation of a type of real estate listing.
// example: {ID: 1, Name: "Apartment"}
type ListingType struct {
	ID   int
	Name string
}

// FromProto populates a system ListingType from an on-wire listing type.
func (n *ListingType) FromProto(p *pb.ListingType) {
	n.ID = int(p.Id)
	n.Name = p.Name
}

// ToProto converts system ListingType to an on-wire listing type.
func (n *ListingType) ToProto() *pb.ListingType {
	return &pb.ListingType{
		Id:   int64(n.ID),
		Name: n.Name,
	}
}
