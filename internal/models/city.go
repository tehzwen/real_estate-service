// package models contains all database model data for this service
// including off the wire representations of on-wire messages.
package models

import (
	pb "github.com/tehzwen/real_estate-service/proto/golang"
)

// City is the system representation of a City.
type City struct {
	ID   int
	Name string
}

// FromProto takes on-wire City and converts to system City.
func (c *City) FromProto(p *pb.City) {
	c.ID = int(p.Id)
	c.Name = p.Name
}

// ToProto converts system City to on the wire City.
func (c *City) ToProto() *pb.City {
	return &pb.City{
		Id:   int64(c.ID),
		Name: c.Name,
	}
}
