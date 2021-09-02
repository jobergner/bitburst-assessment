package persist

import (
	"assessment/pkg/object"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres() *Postgres {
	return &Postgres{}
}

func (p *Postgres) Connect() error {
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable connect_timeout=5",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error while connecting to postgres: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("could not successfully connect to db (%s)", err)
	}

	p.db = db

	return nil
}

func (p *Postgres) WriteObject(o object.Object) error {
	rows, err := p.db.Query("INSERT INTO object (id, online, lastSeen) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET online=$2, lastSeen=$3 ", o.ObjectID, o.Online, o.LastSeen)
	if err != nil {
		return fmt.Errorf("error updating object with id `%d` in db: %s", o.ObjectID, err)
	}
	defer rows.Close()

	return rows.Err()
}

func (p *Postgres) DeleteObject(objectID int, lastSeen int64) error {
	rows, err := p.db.Query("DELETE FROM object WHERE id=$1 AND lastSeen=$2", objectID, lastSeen)
	if err != nil {
		return fmt.Errorf("error deleting object with id `%d` from db: %s", objectID, err)
	}
	defer rows.Close()

	return rows.Err()
}

func (p *Postgres) DeleteObjectsOlderThan(maxValidAge time.Duration) error {
	earliestLastSeen := time.Now().Add(maxValidAge * -1).UnixNano()

	rows, err := p.db.Query("DELETE FROM object WHERE lastSeen < $1", earliestLastSeen)
	if err != nil {
		return fmt.Errorf("error deleting objects older than %d: %s", earliestLastSeen, err)
	}
	defer rows.Close()

	return rows.Err()
}

func (p *Postgres) GetObjects() (map[int]object.Object, error) {

	rows, err := p.db.Query("SELECT * FROM object")
	if err != nil {
		return nil, fmt.Errorf("error selecting all objects in db: %s", err)
	}
	defer rows.Close()

	objs := make(map[int]object.Object)

	for rows.Next() {
		var o object.Object

		if err := rows.Scan(&o.ObjectID, &o.Online, &o.LastSeen); err != nil {
			return nil, fmt.Errorf("error scanning rows: %s", err)
		}

		objs[o.ObjectID] = o
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return objs, nil
}
