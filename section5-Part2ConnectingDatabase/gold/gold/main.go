package main

import (
	"database/sql"
	"goldwatcher/repository"
	"log"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"

	_ "github.com/glebarez/go-sqlite"
)

// Config is the type used to share data with various parts of our application.
// It includes the parts of our GUI that are dynamic and will need to be updated,
// such as the holdings table, gold price info, and the chart. In order to refresh
// those things, we need a reference to them, and this is a convenient place to put
// them, instead of package level variables.
type Config struct {
	App                 fyne.App
	InfoLog             *log.Logger
	ErrorLog            *log.Logger
	DB                  repository.Repository
	MainWindow          fyne.Window
	PriceContainer      *fyne.Container
	ToolBar             *widget.Toolbar
	PriceChartContainer *fyne.Container
	HTTPClient          *http.Client
}

func main() {

	var myApp Config

	// create a fyne application
	fyneApp := app.NewWithID("ca.gocode.goldwatcher.preferences")
	myApp.App = fyneApp
	myApp.HTTPClient = &http.Client{}

	// create our loggers
	myApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	myApp.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// open a connection to the database
	sqlDB, err := myApp.connectSQL()
	if err != nil {
		log.Panic(err)
	}

	// create a database repository
	myApp.setupDB(sqlDB)
	// create and size a fyne window
	myApp.MainWindow = fyneApp.NewWindow("GoldWatcher")
	myApp.MainWindow.Resize(fyne.NewSize(770, 410))
	myApp.MainWindow.SetFixedSize(true)
	myApp.MainWindow.SetMaster()

	myApp.makeUI()

	// show and run the application
	myApp.MainWindow.ShowAndRun()
}

func (app *Config) connectSQL() (*sql.DB, error) {
	//If we create a sqlite database for the first time, it needs to store somewhere in path storage
	//If we run go run . it will work and create folder for us but if
	//build the app and excuete binary, if we start that way, this will fail because execute has no access to write the information to it
	//So instead we will take advantage of fyne storage package, fyne will decide where to store the sqllite database
	//declare a variable path which is an empty string to start with
	path := ""
	//check my env DB_PATH is set to something than empty, then we set path var to env DB_PATH
	if os.Getenv("DB_PATH") != "" {
		path = os.Getenv("DB_PATH")
	} else {
		//Otherwise fyne will do for us, app.App.Storage().RootURI().Path() gives full path for running application then append that to /sql.db
		path = app.App.Storage().RootURI().Path() + "/sql.db"
		app.InfoLog.Println("db in:", path)
	}
	//we need use sql.Open with sqlite with sqlite path
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// After connecting to sqllite, we need to setup our database repository create setupDB receiver of app
func (app *Config) setupDB(sqlDB *sql.DB) {
	//we will use app.DB to repository.NewSQLiteRepository(sqlDB)
	//This will only create empty database connection without database and tables so we call Migrate func to create them
	app.DB = repository.NewSQLiteRepository(sqlDB)
	//call Migrate func to create database and tables
	err := app.DB.Migrate()
	if err != nil {
		app.ErrorLog.Println(err)
		log.Panic()
	}
}
