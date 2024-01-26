package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// Will create App config and beginings of our main routines
type Config struct {
	App      fyne.App
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	//This will be main window for our app
	MainWindow fyne.Window
	//it is reference fyne.Container variable so we can refresh it anytime we need. THis will refresh container
	PriceContainer *fyne.Container
}

var myApp Config

func main() {
	// create a fyne application
	//Difference between app.New and app.NewWithID is, newwithID we use when we distribute app to playstore/Appstore or to be disturbed with people
	//basically It allows to have unique identifer for your application, the convention means naming our application use your domain name inverse ca.gocode and name of app goldwatchers and put prefereneces at end
	fyneApp := app.NewWithID("ca.gocode.goldwatcher.preferences")
	//Assign fyneApp to our Config which is of type myApp config type
	myApp.App = fyneApp

	// create our loggers
	//Create INFO logger, by setting os.Stdout which is to terminal and prefix by Info with tab space and Longtime and LongDate
	myApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	//Create error logger, by setting os.Stdout which is to terminal and prefix by Info with tab space and Longtime and LongDate
	//We will write the line where error took place log.Lshortfile
	myApp.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// open a connection to the database

	// create a database repository
	//repository is simple pattern to interact with things

	// create and size a fyne window
	myApp.MainWindow = fyneApp.NewWindow("GoldWatcher")
	myApp.MainWindow.Resize(fyne.NewSize(770, 410))
	myApp.MainWindow.SetFixedSize(true)
	//Set this as main window of app using method SetMaster()
	myApp.MainWindow.SetMaster()

	myApp.makeUI()

	// show and run the application
	myApp.MainWindow.ShowAndRun()
}
