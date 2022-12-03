// package models contains all database model data for this service
// including off the wire representations of on-wire messages.
package models

import (
	"time"

	pb "github.com/tehzwen/real_estate-service/proto/golang"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Listing is the system representation of a real estate listing.
type Listing struct {
	ID            int
	Address       string
	Price         int
	City          City
	Neighbourhood Neighbourhood
	Type          ListingType
	AddedDate     time.Time
	LastUpdated   time.Time
	MlsID         string
	URL           string
}

// FromProto takes an on-wire listing and populates a system Listing from it.
func (l *Listing) FromProto(p *pb.Listing) {
	c := City{}
	c.FromProto(p.City)
	l.City = c

	n := Neighbourhood{}
	n.FromProto(p.Neighbourhood)
	l.Neighbourhood = n

	t := ListingType{}
	t.FromProto(p.Type)
	l.Type = t

	l.ID = int(p.Id)
	l.Address = p.Address
	l.Price = int(p.Price)

	l.AddedDate = p.AddedDate.AsTime()
	l.LastUpdated = p.LastUpdated.AsTime()
	l.MlsID = p.MlsId
	l.URL = p.Url
}

// ToProto converts our system listing to an on-wire representation.
func (l *Listing) ToProto() *pb.Listing {
	return &pb.Listing{
		Id:            int32(l.ID),
		Address:       l.Address,
		Price:         int32(l.Price),
		City:          l.City.ToProto(),
		Neighbourhood: l.Neighbourhood.ToProto(),
		Type:          l.Type.ToProto(),
		AddedDate:     timestamppb.New(l.AddedDate),
		LastUpdated:   timestamppb.New(l.LastUpdated),
		MlsId:         l.MlsID,
		Url:           l.URL,
	}
}
