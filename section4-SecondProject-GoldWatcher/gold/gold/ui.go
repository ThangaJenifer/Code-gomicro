package main

import "fyne.io/fyne/container"

func (app *Config) makeUI() {
	// get the current price of gold
	//we need openPrice, currentPrice, priceChange to get the current value of gold, we get this by calling function getPriceText() on app method
	openPrice, currentPrice, priceChange := app.getPriceText()

	//I got openPrice, currentPrice, priceChange from previous step, Now I need to put them in its own container
	// put price information into a container
	//using NewGridWithColumns with 1 row and 3 columns and content openprice, currentPrice and priceChange
	priceContent := container.NewGridWithColumns(3,
		openPrice,
		currentPrice,
		priceChange,
	)
	//Lets save priceContent in our app Config type app.PriceContainer for refresh purpose
	app.PriceContainer = priceContent

	//Above we have a container and we can place that anywhere we need to in our app
	// add container to window
	//we use NewVBox container because that stack things top to bottom
	finalContent := container.NewVBox(priceContent)
	//We need to SetContent of the MainWindow to finalContent
	app.MainWindow.SetContent(finalContent)
}
