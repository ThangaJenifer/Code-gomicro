package repository

import (
	"testing"
	"time"
)

// Test Migrate function
func TestSQLiteRepository_Migrate(t *testing.T) {
	//Now testRepo is SQLiteRepository type so can utilizw all the methods like Migrate
	//We need to have a file sql.db in our testdata folder now
	err := testRepo.Migrate()
	if err != nil {
		t.Error("migrate failed:", err)
	}
}

func TestSQLiteRepository_InsertHolding(t *testing.T) {
	//this function requires a holding of h data
	h := Holdings{
		Amount:        1,
		PurchaseDate:  time.Now(),
		PurchasePrice: 1000,
	}
	//pass that h to InsertHolding method
	result, err := testRepo.InsertHolding(h)
	if err != nil {
		t.Error("insert failed:", err)
	}
	//make sure result.ID > 0
	if result.ID <= 0 {
		t.Error("invalid id sent back:", result.ID)
	}
}

func TestSQLiteRepository_AllHoldings(t *testing.T) {
	//Pass through all holdings
	h, err := testRepo.AllHoldings()
	if err != nil {
		t.Error("get all failed:", err)
	}

	if len(h) != 1 {
		t.Error("wrong number of rows returned; expected 1, but got", len(h))
	}
}

func TestSQLiteRepository_GetHoldingByID(t *testing.T) {
	//we only have 1 holding in database so ID is 1, we try with one
	h, err := testRepo.GetHoldingByID(1)
	if err != nil {
		t.Error("get by id failed:", err)
	}
	//check PurchasePrice check if 1000 or not
	if h.PurchasePrice != 1000 {
		t.Error("wrong purchase price returned; expected 1000 but got", h.PurchasePrice)
	}
	//Try non-exisiting holding as 2 does not exist
	_, err = testRepo.GetHoldingByID(2)
	if err == nil {
		t.Error("get one returned value for non-existent id")
	}
}

func TestSQLiteRepository_UpdateHolding(t *testing.T) {
	//geting holidng id 1 in variable h
	h, err := testRepo.GetHoldingByID(1)
	if err != nil {
		t.Error(err)
	}
	//Updating holding 1 with PurchasePrice 1000 to 1001
	h.PurchasePrice = 1001

	err = testRepo.UpdateHolding(1, *h)
	if err != nil {
		t.Error("update failed:", err)
	}
}

func TestSQLiteRepository_DeleteHolding(t *testing.T) {
	err := testRepo.DeleteHolding(1)
	if err != nil {
		t.Error("failed to delete holding", err)
		if err != errDeleteFailed {
			t.Error("wrong error returned")
		}
	}

	err = testRepo.DeleteHolding(2)
	if err == nil {
		t.Error("no error when trying to delete non-existent record")
	}

}
