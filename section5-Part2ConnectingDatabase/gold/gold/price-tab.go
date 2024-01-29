package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func (app *Config) pricesTab() *fyne.Container {
	//Creat a chart using app.getChart() which is func mentioned below which returns graph image
	chart := app.getChart()
	//Create a new Vbox container for it
	chartContainer := container.NewVBox(chart)
	//Create a new PriceChartContainer type in app type of Config in main.go and assign it to app config here
	app.PriceChartContainer = chartContainer

	return chartContainer
}

// getChart methods calls the download file
func (app *Config) getChart() *canvas.Image {
	//apiURL will compose the URL to the image png location of graph
	apiURL := fmt.Sprintf("https://goldprice.org/charts/gold_3d_b_o_%s_x.png", strings.ToLower(currency))
	//Create a img var with fyne lib canvas.Image which is not a normal image
	var img *canvas.Image
	//Call the downloadfile method
	err := app.downloadFile(apiURL, "gold.png")
	if err != nil {
		// use bundled image png which shows as a error instead of price chart
		//we use bundle feature of fyne and golang here. This will convert image to go code. copy Unreachable.Png to local path
		//use command fyne bundle unreachable.png >> bundled.go
		//Check bundled.go file, we will see package level variable called resourceUnreachablePng and gives this as gocode
		//canvas.NewImageFromResource(resourceUnreachablePng)
		//Bundle many images then check documentation of fyne use fyne bundle --append unreachable1.png >> bundled.go for MAC
		//fyne bundle -o bundled.go unreachable.png for windows
		img = canvas.NewImageFromResource(resourceUnreachablePng)
	} else {
		img = canvas.NewImageFromFile("gold.png")
	}
	//we need to get this image back from remote resource. we need to tell our app how big it is
	//We can do that by calling SetMinSize on our canvas.Image variable img
	img.SetMinSize(fyne.Size{
		Width:  770,
		Height: 410,
	})
	//means of filling available windows space.canvas.ImageFillOriginal- Keep the orginial image size and fill all other spaces
	img.FillMode = canvas.ImageFillOriginal

	return img
}

// This func is pointer receiver to app type *config and downloadFile takes URL and filename as string and can return potential error
func (app *Config) downloadFile(URL, fileName string) error {
	// get the response bytes from calling a url
	response, err := app.HTTPClient.Get(URL)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New("received wrong response code when downloading image")
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	//use image package from std library and decode func on it using bytes.NewReader(b)
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return err
	}
	//Once control is here, we have our image, we need to save that to our file system
	//Declare var out, create a file using os.Create fmt.Sprintf("./%s", fileName) , save it in current directory, when we run the app we can see the file in our location
	out, err := os.Create(fmt.Sprintf("./%s", fileName))
	if err != nil {
		return err
	}
	//we want to encode out file to png format
	//func png.Encode(w io.Writer, m image.Image) error
	//Encode writes the Image m to w in PNG format. Any Image may be encoded, but images that are not image.NRGBA might be encoded lossily.
	err = png.Encode(out, img)
	if err != nil {
		return err
	}

	return nil
}
