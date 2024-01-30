package main

import "testing"

func TestConfig_getHoldings(t *testing.T) {
	//Currentholdings() func will call the AllHoldings() func written in db-test.go which has 2 elements holding prepopulated
	all, err := testApp.currentHoldings()
	if err != nil {
		t.Error("failed to get current holdings from database:", err)
	}

	if len(all) != 2 {
		t.Error("wrong number of rows returned")
	}
}

func TestConfig_getHoldingSlice(t *testing.T) {
	slice := testApp.getHoldingSlice()
	//This should written two records along with one table heading, so length should be 3
	if len(slice) != 3 {
		t.Error("wrong number of rows returned; expected 3 but got", len(slice))
	}
}
