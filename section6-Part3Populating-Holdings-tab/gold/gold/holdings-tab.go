package main

import (
	"fmt"
	"goldwatcher/repository"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (app *Config) holdingsTab() *fyne.Container {
	//holdingsTab should have to data of holding and table where getHoldingsTable() func has access app.Holdings = data and itself is a widget table
	app.HoldingsTable = app.getHoldingsTable()
	//for this table we need to create a container of New vertical box and pass the holdings table to it
	holdingsContainer := container.NewVBox(app.HoldingsTable)

	return holdingsContainer
}

// func of receiver of app and pointer to config, getHoldingsTable has no args and return pointer to the table
// This will take the holding slice and show that in widget take
func (app *Config) getHoldingsTable() *widget.Table {
	//data will have holding slice
	data := app.getHoldingSlice()
	//put the data to app.Holdings
	app.Holdings = data

	t := widget.NewTable(
		//func returns int, int which is len(data) which is no. of rows and len(data[0]) no. of columns
		func() (int, int) {
			return len(data), len(data[0])
		},
		//second one returns a template which is used by below function func(i widget.TableCellID, o fyne.CanvasObject)
		//it will return container ctr
		func() fyne.CanvasObject {
			//we will create a new VerticalBox with widget empty label in it and will return container ctr
			ctr := container.NewVBox(widget.NewLabel(""))
			return ctr
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			//If the current column is equal to length of my data (interface of slice) at index 0 minus 1 && i.Row != 0
			//If the condition is true, then we are at last element of the table
			if i.Col == (len(data[0])-1) && i.Row != 0 {
				// last cell - put in a button
				//Creating delete trash icon using widget.NewButtonWithIcon func "Delete", theme.DeleteIcon() and inline func
				w := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
					//we need confirmation to delete it which requires DELETE title
					dialog.ShowConfirm("Delete?", "", func(deleted bool) {
						//data[i.Row][0] our ID will be stored here we convert that to interger
						id, _ := strconv.Atoi(data[i.Row][0].(string))
						//We will delete the holding in the database
						err := app.DB.DeleteHolding(int64(id))
						if err != nil {
							app.ErrorLog.Println(err)
						}
						// TODO: refresh the holdings table (anytime we do content of the table, we need to refresh that)
						// To refresh table we need reference two things, both the table widget itself and also data for the table
						//So go and add Holding [][]interface{} and HoldingsTable *widget.Table to the Config struct main.go, so we can save reference to app when we need to
						app.refreshHoldingsTable()
					}, app.MainWindow) //add what window this delete dialog attached to is mainWindow
				})
				//Created button but still not shown yet so to make button standout put w.Importance  = widget.HighImportance, this will make it nice blue
				w.Importance = widget.HighImportance
				//finally we had to do somehitng with w, o.(*fyne.Container).Objects = []fyne.CanvasObject{ w,}
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					w,
				}
			} else {
				// we're just putting in textual information
				//we will refer to o and cast it to (*fyne.Container).Objects equal to
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					widget.NewLabel(data[i.Row][i.Col].(string)),
				}
			}
		})
	//we need to tell table exactly its width should be using slice of float32
	colWidths := []float32{50, 200, 200, 200, 110}
	for i := 0; i < len(colWidths); i++ {
		//setting the column width to the i
		t.SetColumnWidth(i, colWidths[i])
	}

	return t

}

// func of receiver of app and pointer to config, this is going to get slice of current holdings of user has so name of func is getHoldingSlice()
// this will return slice of slice of interface [][]interface{}, it will hold absolutely anything, slice of the slice of anything any type
// This will have one row and second slice will correspond to column for information [][]interface{}
func (app *Config) getHoldingSlice() [][]interface{} {
	var slice [][]interface{}
	//we get all holding using app.currentHoldings() written just below
	holdings, err := app.currentHoldings()
	if err != nil {
		app.ErrorLog.Println(err)
	}
	//So lets some heading for our table here
	//We have used [][]interface{} as "ID", "Amount", "Price", "Date" are labels and Delete is a button type
	slice = append(slice, []interface{}{"ID", "Amount", "Price", "Date", "Delete?"})
	//range through all holdings
	for _, x := range holdings {
		//here itself we want to populate our table, but it will be useful to begin table with headings
		//otherwise ppl wont have idea what appears in what column
		var currentRow []interface{}
		//first 4 are valid data and last one is a dummy button
		//First one is ID from database which is integer so we need to convert into string with strconv FormatInt(x.ID, 10)) valuew is x.ID and 10 is base 10 for strconv
		currentRow = append(currentRow, strconv.FormatInt(x.ID, 10))
		//Second one is Amount, convert with fmt.Sprintf
		currentRow = append(currentRow, fmt.Sprintf("%d toz", x.Amount))
		//from database we get PurchasePrice as integer so cast it into float32
		currentRow = append(currentRow, fmt.Sprintf("$%2f", float32(x.PurchasePrice/100)))
		//Format time using Format ISO standard
		currentRow = append(currentRow, x.PurchaseDate.Format("2006-01-02"))
		//Dummy button right now
		currentRow = append(currentRow, widget.NewButton("Delete", func() {}))

		slice = append(slice, currentRow)
	}

	return slice
}

// func of receiver of app and pointer to config, it will return of slice holdings from the database and potential error
func (app *Config) currentHoldings() ([]repository.Holdings, error) {
	//Inside app we have access to DB so we will call app.DB.AllHoldings() to get all current holdings
	holdings, err := app.DB.AllHoldings()
	if err != nil {
		app.ErrorLog.Println(err)
		return nil, err
	}

	return holdings, nil
	//Next we use holdings returned by DB converting that into a slice so we can use as we build our table using getHoldingSlice() method
}

/*
Check fyne documentation for table, we will have three functions,
1st func for size of the table rows and columns, second for known for template kind of information we are putting in the table
3rd func will build the table one row at a time. Thats the way widget.NewTable will work

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

var data = [][]string{[]string{"top left", "top right"},
	[]string{"bottom left", "bottom right"}}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Table Widget")

	list := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Row][i.Col])
		})

	myWindow.SetContent(list)
	myWindow.ShowAndRun()
}
*/
