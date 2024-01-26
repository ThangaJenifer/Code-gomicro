package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// getPriceText() method will return special type of *canvas.Text from fyne library
func (app *Config) getPriceText() (*canvas.Text, *canvas.Text, *canvas.Text) {
	//We get current price of gold by using this technique, var g will hold the data coming from json
	//open, current and change will be type of *canvas.Text which we return
	var g Gold
	var open, current, change *canvas.Text
	//call GetPrices() method to get pointer to *price type and error which we used
	gold, err := g.GetPrices()
	//If we wont able to get response or internet problem to api website
	//we need to return some meaningful error indicating network not reachable
	if err != nil {
		//Grey color will be unviserval in our app which something went wrong
		//throughout our app lets use that grey for error by using image/color package
		grey := color.NRGBA{R: 155, G: 155, B: 155, A: 255}
		//set open, current and change to Open: Unreachable, Current: Unreachable, Change: Unreachable with canvas package withText
		open = canvas.NewText("Open: Unreachable", grey)
		current = canvas.NewText("Current: Unreachable", grey)
		change = canvas.NewText("Change: Unreachable", grey)
	} else { //else statement will be success , color.NRGBA{R: 0, G: 180, B: 0, A: 255} which is nice color
		//Red set to 0,Green set to 180, blue set to 0 and Alpha set to 255
		displayColor := color.NRGBA{R: 0, G: 180, B: 0, A: 255}
		//gold.Price < gold.PreviousClose will change the color to RED more 180
		if gold.Price < gold.PreviousClose {
			displayColor = color.NRGBA{R: 180, G: 0, B: 0, A: 255}
		}
		//format opentxt, currentTxt, Changetxt in Open: $%.4f %s which shows string formatted value and currency
		openTxt := fmt.Sprintf("Open: $%.4f %s", gold.PreviousClose, currency)
		currentTxt := fmt.Sprintf("Current: $%.4f %s", gold.Price, currency)
		changeTxt := fmt.Sprintf("Change: $%.4f %s", gold.Change, currency)
		//now we need to build our returns open, current, change with canvas.Newtext function
		//func canvas.NewText(text string, color color.Color) *canvas.Text , NewText returns a new Text implementation
		//For open we will give color nil, which will default to system color
		open = canvas.NewText(openTxt, nil)
		current = canvas.NewText(currentTxt, displayColor)
		change = canvas.NewText(changeTxt, displayColor)
	}
	//Set alignment, left alignment using TextAlignLeading, center to .TextAlignCenter, right one TextAlignTrailing
	open.Alignment = fyne.TextAlignLeading
	current.Alignment = fyne.TextAlignCenter
	change.Alignment = fyne.TextAlignTrailing

	return open, current, change
}
