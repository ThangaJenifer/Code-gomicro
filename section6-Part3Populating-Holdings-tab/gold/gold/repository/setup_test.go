package repository

import (
	"database/sql"
	"log"
	"os"
	"testing"

	//Importing third party driver for go-sqllite
	_ "github.com/glebarez/go-sqlite"
)

// create variable testRepo of type pointer to *SQLiteRepository, that package level but ignored by compiler in the actual application
var testRepo *SQLiteRepository

// Create a test for main function and create a testing environment
func TestMain(m *testing.M) {
	//ignoring err and using os package removing file if exists in path ./testdata/sql.db
	//This is the place we store temporary database while testing
	_ = os.Remove("./testdata/sql.db")
	//function level var path will be assigned to temporary database location
	path := "./testdata/sql.db"
	//creating a database connection using the temporary database location
	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Println(err)
	}
	//using testRepo to store the sqllite db connection pool
	testRepo = NewSQLiteRepository(db)
	//we will run the test
	os.Exit(m.Run())
}
