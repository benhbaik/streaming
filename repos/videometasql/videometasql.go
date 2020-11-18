package videometasql

import (
	"database/sql"
	"fmt"
	videometamodel "streaming/models/videometamodel"
)

// SQLRepo Handles database logic for video metadata.
type SQLRepo struct {
	// db *sql.DB
}

// New videometadata.Repo constructor.
func New() *SQLRepo {
	// db, err := sql.Open("mysql", "streaming:3XC3!!ence@/streaming")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer db.Close()
	sr := new(SQLRepo)
	// sr.db = db
	return sr
}

// Insert Inserts video metadata into DB.
func (sr *SQLRepo) Insert(vm videometamodel.Model) error {
	fmt.Println("Opening connection to DB.")

	db, err := sql.Open("mysql", "streaming:3XC3!!ence@/streaming")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	insertStatement := `INSERT INTO videometadata (name, filename, extension, directory, filepath)
						VALUES (?, ?, ?, ?, ?)`

	fmt.Println(insertStatement)

	result, err := db.Exec(insertStatement, vm.Filename, vm.Fullfilename, vm.Extension, vm.Directory, vm.Filepath)
	if err != nil {
		return err
	}

	fmt.Println(result.RowsAffected())
	fmt.Println(result.LastInsertId())

	fmt.Println("Closing connection to DB.")
	return nil
}

func executeSQLQuery(sqlQuery string) (*sql.Rows, error) {
	db, err := sql.Open("mysql", "streaming:3XC3!!ence@/streaming")
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
