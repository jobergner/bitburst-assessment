package persist

import (
	"assessment/pkg/object"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
	mu sync.RWMutex
}

func NewPostgres() *Postgres {
	return &Postgres{}
}

func (p *Postgres) Connect() error {
	connStr := "user=postgres password=postgrespassword dbname=assessment sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error while connecting to postgres: %s", err)
	}

	p.db = db

	return nil
}

func (p *Postgres) WriteObject(o object.Object) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	objectExists, err := p.objectExists(o.ObjectID)
	if err != nil {
		return fmt.Errorf("error while chicking whether object `%d` exists: %s", o.ObjectID, err)
	}

	if objectExists {
		return p.updateObject(o)
	}

	return p.createObject(o)
}

func (p *Postgres) updateObject(o object.Object) error {
	rows, err := p.db.Query("UPDATE object SET id=$1, online=$2, lastSeen=$3 WHERE id=$1", o.ObjectID, o.Online, o.LastSeen)
	if err != nil {
		return fmt.Errorf("error updating object with id `%d` in db: %s", o.ObjectID, err)
	}
	defer rows.Close()

	return rows.Err()
}

func (p *Postgres) createObject(o object.Object) error {
	rows, err := p.db.Query("INSERT INTO object (id, online, lastSeen) VALUES ($1, $2, $3)", o.ObjectID, o.Online, o.LastSeen)
	if err != nil {
		return fmt.Errorf("error inserting object with id `%d` into db: %s", o.ObjectID, err)
	}
	defer rows.Close()

	return rows.Err()
}

func (p *Postgres) DeleteObject(objectID int, lastSeen int64) error {
	p.mu.RLock()
	defer p.mu.Unlock()

	rows, err := p.db.Query("DELETE FROM object WHERE id=$1 AND lastSeen=$2", objectID, lastSeen)
	if err != nil {
		return fmt.Errorf("error deleting object with id `%d` from db: %s", objectID, err)
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

func (p *Postgres) objectExists(objectID int) (bool, error) {

	rows, err := p.db.Query("SELECT (id) FROM object WHERE id=$1", objectID)
	if err != nil {
		return false, fmt.Errorf("error finding specific object in db: %s", err)
	}
	defer rows.Close()

	objs := make(map[int]object.Object)

	for rows.Next() {
		var o object.Object

		if err := rows.Scan(&o.ObjectID, &o.Online, &o.LastSeen); err != nil {
			return false, fmt.Errorf("error scanning rows: %s", err)
		}

		objs[o.ObjectID] = o
	}

	if err = rows.Err(); err != nil {
		return false, err
	}

	_, exists := objs[objectID]
	return exists, nil
}
