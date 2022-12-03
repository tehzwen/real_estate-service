// package db holds all structs and functions used for interacting with the
// database portion of the project.
package db

import (
	"context"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/tehzwen/real_estate-service/internal/models"
	pb "github.com/tehzwen/real_estate-service/proto/golang"
)

const (
	defaultQueryTimeout = time.Duration(time.Millisecond * 8000)
	defaultMaxPrice     = 1000000
	defaultMinPrice     = 0
	defaultGetLimit     = 100
)

// Worker is an interface that can be used to work with different database models
// ie: postgres/memory/mysql etc.
type Worker interface {
	GetListings(context.Context, *GetListingsFilter) ([]models.Listing, error)
}

// GetListingsFilter is the filter used for filtering down listings data.
type GetListingsFilter struct {
	Since          time.Time
	Until          time.Time
	MaxPrice       int
	MinPrice       int
	Cities         []string
	Neighbourhoods []string
	Types          []string
	Limit          int
	PageToken      string
}

func (gf *GetListingsFilter) defaultValues() {
	sinceTime := time.Now().Add(-30 * (time.Hour * 24))
	untilTime := time.Now()

	gf.Since = sinceTime
	gf.Until = untilTime

	gf.Limit = defaultGetLimit
	gf.MaxPrice = defaultMaxPrice
	gf.MinPrice = defaultMinPrice
	gf.Cities = []string{}
	gf.Neighbourhoods = []string{}
	gf.Types = []string{}
}

// FromProto takes a proto GetListingsFilter and generates a server side representation
// of it to be used while querying.
func (gf *GetListingsFilter) FromProto(f *pb.GetListingsFilter) {
	gf.defaultValues()

	if f == nil {
		return
	}

	// default the filter times
	if f.TimeSpan != nil && f.TimeSpan.Since != nil {
		gf.Since = f.TimeSpan.Since.AsTime()
	}

	if f.TimeSpan != nil && f.TimeSpan.Until != nil {
		gf.Until = f.TimeSpan.Until.AsTime()
	}

	gf.MaxPrice = int(f.MaxPrice)
	if f.MaxPrice <= 0 {
		gf.MaxPrice = defaultMaxPrice
	}

	gf.MinPrice = int(f.MinPrice)
	if f.MinPrice <= 0 {
		gf.MinPrice = defaultMinPrice
	}

	gf.Cities = append(gf.Cities, f.Cities...)
	gf.Neighbourhoods = append(gf.Neighbourhoods, f.Neighbourhoods...)
	gf.Types = append(gf.Types, f.Types...)
}

// Credentials is a struct used for connecting to databases. Will likely be expanded more to
// allow SSL/other database query strings, for now was lazy.
type Credentials struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

// TODO -> change this to not be postgres specific.
// GetConnectionString returns the connection string generated from credentials in order to connect.
func (c *Credentials) GetConnectionString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Database)
}
