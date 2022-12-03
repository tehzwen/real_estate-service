// package models contains all database model data for this service
// including off the wire representations of on-wire messages.
package models

import (
	pb "github.com/tehzwen/real_estate-service/proto/golang"
)

// Neighbourhood is the system representation of a neighbourhood.
type Neighbourhood struct {
	ID   int
	Name string
}

// FromProto populates a system neighbourhood from an on-wire neighbourhood.
func (n *Neighbourhood) FromProto(p *pb.Neighbourhood) {
	n.ID = int(p.Id)
	n.Name = p.Name
}

// ToProto converts system neighbourhood to on-wire neighbourhood.
func (n *Neighbourhood) ToProto() *pb.Neighbourhood {
	return &pb.Neighbourhood{
		Id:   int64(n.ID),
		Name: n.Name,
	}
}
