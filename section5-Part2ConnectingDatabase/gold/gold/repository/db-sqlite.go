package repository

import (
	"database/sql"
	"errors"
	"time"
)

// create a type for SQLiteRepository with connection of sql.DB for it
// Thi s type has one connection to pool connection of sqllite database
type SQLiteRepository struct {
	Conn *sql.DB
}

// we will create a function NewSQLiteRepository, use it to create a repository which is connected to SQLLite
func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	//From main func if we implement this NewSQLiteRepository func and provide db *sql.DB which connected then we get the repository
	return &SQLiteRepository{
		Conn: db,
	}
}

// type SQLiteRepository will implement the Repository interface type methods
func (repo *SQLiteRepository) Migrate() error {
	//Func to intialise the database
	query := `
	create table if not exists holdings(
		id integer primary key autoincrement,
		amount real not null,
		purchase_date integer not null,
		purchase_price integer not null);
	`
	//There is no concept in sqllite like time so it can be stored as integer/string.
	//Thats why purchase_date is integer here
	//sql.Result ignore and error := repo from receiver conn from connection.exec with query
	_, err := repo.Conn.Exec(query)
	return err
}

func (repo *SQLiteRepository) InsertHolding(holdings Holdings) (*Holdings, error) {
	stmt := "insert into holdings (amount, purchase_date, purchase_price) values (?, ?, ?)"
	//repo connection and execute stmt with holdings.Amount, holdings.PurchaseDate which will be in time.Time so use unix() func to convert it in integer format of int64 and holdings.PurchasePrice
	res, err := repo.Conn.Exec(stmt, holdings.Amount, holdings.PurchaseDate.Unix(), holdings.PurchasePrice)
	if err != nil {
		return nil, err
	}
	//save sql.Result LastInsertId() to id variable
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	//Assign this id variable to holdings.ID and return the holdings variable now which has been passed to function
	holdings.ID = id
	return &holdings, nil
}

func (repo *SQLiteRepository) AllHoldings() ([]Holdings, error) {
	query := "select id, amount, purchase_date, purchase_price from holdings order by purchase_date"
	//repo connection and execute query
	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Holdings
	for rows.Next() {
		var h Holdings
		var unixTime int64
		err := rows.Scan(
			&h.ID,
			&h.Amount,
			&unixTime,
			&h.PurchasePrice,
		)
		if err != nil {
			return nil, err
		}
		//From time package Unix() with unixtime we populated and 0.
		h.PurchaseDate = time.Unix(unixTime, 0)
		all = append(all, h)
	}

	return all, nil
}

func (repo *SQLiteRepository) GetHoldingByID(id int) (*Holdings, error) {
	row := repo.Conn.QueryRow("select id, amount, purchase_date, purchase_price from holdings where id = ?", id)

	var h Holdings
	var unixTime int64
	err := row.Scan(
		&h.ID,
		&h.Amount,
		&unixTime,
		&h.PurchasePrice,
	)

	if err != nil {
		return nil, err
	}

	h.PurchaseDate = time.Unix(unixTime, 0)

	return &h, nil
}

func (repo *SQLiteRepository) UpdateHolding(id int64, updated Holdings) error {
	if id == 0 {
		return errors.New("invalid updated id")
	}

	stmt := "Update holdings set amount = ?, purchase_date = ?, purchase_price = ? where id = ?"
	res, err := repo.Conn.Exec(stmt, updated.Amount, updated.PurchaseDate.Unix(), updated.PurchasePrice, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errUpdateFailed
	}

	return nil
}

func (repo *SQLiteRepository) DeleteHolding(id int64) error {
	res, err := repo.Conn.Exec("delete from holdings where id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errUpdateFailed
	}

	return nil
}
