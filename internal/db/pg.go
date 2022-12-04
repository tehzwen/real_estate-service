// package db holds all structs and functions used for interacting with the
// database portion of the project.
package db

import (
	"context"
	"database/sql"
	"encoding/base64"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/tehzwen/real_estate-service/internal/models"
)

var (
	selectListingsQuery = `SELECT 
	l.id as listing_id, l.address, l.price, 
	c.id as city_id, c.name as city_name,
	n.id as neighbourhood_id, n.name as neighbourhood_name,
	t.id as listing_type_id, t.name as listing_type_name,
	added_date,
	last_updated,
	mls_id,
	url
	FROM listings l
		inner join cities c on l.city_id = c.id
		inner join neighbourhoods n on l.neighbourhood_id = n.id
		inner join listing_types t on l.type_id = t.id
	WHERE 
		added_date >= $1
		AND
		added_date < $2
		AND
		(l.price between $4 AND $5)
	ORDER BY added_date, l.id ASC
	LIMIT $3;`
	selectListingsPaginateQuery = `SELECT 
	l.id as listing_id, l.address, l.price, 
	c.id as city_id, c.name as city_name,
	n.id as neighbourhood_id, n.name as neighbourhood_name,
	t.id as listing_type_id, t.name as listing_type_name,
	added_date,
	last_updated,
	mls_id,
	url
	FROM listings l
		inner join cities c on l.city_id = c.id
		inner join neighbourhoods n on l.neighbourhood_id = n.id
		inner join listing_types t on l.type_id = t.id
	WHERE 
		(added_date, l.id) > ($1, $3)
		AND
		added_date < $2
		AND
		(l.price between $5 AND $6)
	ORDER BY added_date, l.id ASC
	LIMIT $4;`
)

// PostgresWorker is the concrete type that implements db.Worker to work with
// postgresql databases.
type PostgresWorker struct {
	db *sql.DB
}

// NewPostgresWorker takes credentials and connects to a postgres database, returning
// the Worker or error.
func NewPostgresWorker(ctx context.Context, c Credentials) (Worker, error) {
	d, err := sql.Open("postgres", c.GetConnectionString())
	if err != nil {
		return nil, err
	}

	tc, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	if err = d.PingContext(tc); err != nil {
		return nil, err
	}

	return &PostgresWorker{
		db: d,
	}, nil
}

// GetListings takes a filter and retrieves listings data based off of that filter. Postgres
// implementation.
func (p *PostgresWorker) GetListings(ctx context.Context, filter *GetListingsFilter) ([]models.Listing, error) {
	var listings []models.Listing
	var query string = selectListingsQuery
	var sinceNum int64 = filter.Since.Unix()
	var rows *sql.Rows
	var err error
	var args []any

	tc, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	args = []any{
		filter.Since.Unix(),
		filter.Until.Unix(),
		filter.Limit,
		filter.MinPrice,
		filter.MaxPrice,
	}

	if filter.PageToken != "" {
		decoded, err := base64.StdEncoding.DecodeString(filter.PageToken)
		if err != nil {
			log.Printf("error: %v", err)
			return nil, ErrDecode
		}

		// split on the comma to get the id & the date
		values := strings.Split(string(decoded), ",")
		t, err := strconv.Atoi(values[0])
		if err != nil {
			log.Printf("error: %v", err)
			return nil, ErrDecode
		}

		id, err := strconv.Atoi(values[1])
		if err != nil {
			log.Printf("error: %v", err)
			return nil, ErrDecode
		}

		query = selectListingsPaginateQuery
		sinceNum = int64(t)
		args = []any{
			sinceNum,
			filter.Until.Unix(),
			id,
			filter.Limit,
			filter.MinPrice,
			filter.MaxPrice,
		}
	}

	log.Printf("GetListings Filter: %+v", filter)
	rows, err = p.db.QueryContext(
		tc,
		query,
		args...,
	)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, ErrQuery
	}
	defer rows.Close()

	for rows.Next() {
		v := models.Listing{}
		var addedTimeUnix int64
		var updatedTimeUnix int64

		if err = rows.Scan(
			&v.ID,
			&v.Address,
			&v.Price,
			&v.City.ID,
			&v.City.Name,
			&v.Neighbourhood.ID,
			&v.Neighbourhood.Name,
			&v.Type.ID,
			&v.Type.Name,
			&addedTimeUnix,
			&updatedTimeUnix,
			&v.MlsID,
			&v.URL,
		); err != nil {
			log.Printf("error: %v", err)
			return nil, ErrScan
		}

		// necessary to convert between time.Time & int64.
		v.AddedDate = time.Unix(addedTimeUnix, 0)
		v.LastUpdated = time.Unix(updatedTimeUnix, 0)
		listings = append(listings, v)
	}

	return listings, nil
}
