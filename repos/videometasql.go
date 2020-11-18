package videometasql

import (
	"database/sql"
	"streaming/models/videometamodel"
)

// Repo Handles database logic for video metadata.
type Repo struct {
	db *sql.DB
}

// New videometadata.Repo constructor.
func New(db *sql.DB) *Repo {
	vmr := new(Repo)
	vmr.db = db
	return vmr
}

// Insert Inserts video metadata into DB.
func (*Repo) Insert(vmm videometamodel.Model) {

}

func ExecuteSQLQuery(sqlQuery string) (*sql.Rows, error) {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	results, err := db.Query(sqlQuery)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return results, err
}
