package repository

/* IN this course we wont use mysql/postgres as done in webserver,
we want to embed the database in our fyne excutable, for that mysql lite is the best one.
It is best database in the work as it is embedded, if we use android device it certainly uses sqllite. Probably many applications in laptop will use it as well
SQL lite is small, lite and reliable. If powergoes, app crashes also but still sqllite inside database sqllite wont destroy
We will use database repository pattern which is popular and good to use even though embedding database into our app
This makes testing easy and switch the databases if we need
*/

import (
	"errors"
	"time"
)

// create var errUpdateFailed and errDeleteFailed with errors errors.New("update failed") and errors.New("delete failed")
var (
	errUpdateFailed = errors.New("update failed")
	errDeleteFailed = errors.New("delete failed")
)

// create a type Repository of interface and place funcs which interact with databases
type Repository interface {
	Migrate() error //this will be created by begining when there is no databse, it will create tables or watever else we need in database
	InsertHolding(h Holdings) (*Holdings, error)
	AllHoldings() ([]Holdings, error)
	GetHoldingByID(id int) (*Holdings, error)
	UpdateHolding(id int64, updated Holdings) error
	DeleteHolding(id int64) error
}

type Holdings struct {
	ID            int64     `json:"id"`
	Amount        int       `json:"amount"`
	PurchaseDate  time.Time `json:"purchase_date"`
	PurchasePrice int       `json:"purchase_price"`
}
