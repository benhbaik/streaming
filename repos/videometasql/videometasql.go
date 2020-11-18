package videometasql

import (
	"database/sql"
	videometamodel "streaming/models/videometamodel"
)

// SQLRepo Handles database logic for video metadata.
type SQLRepo struct {
	db *sql.DB
}

// New videometadata.Repo constructor.
func New() *SQLRepo {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	vmr := new(SQLRepo)
	vmr.db = db
	return vmr
}

// Insert Inserts video metadata into DB.
func (r *SQLRepo) Insert(vm videometamodel.Model) {

}

func executeSQLQuery(sqlQuery string) (*sql.Rows, error) {
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
