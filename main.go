package main

import (
	auth "AmbiTradfri/authentication"
	display "AmbiTradfri/display"
	"flag"
	"fmt"
	"os"
	"time"

	tradfri "github.com/eriklupander/tradfri-go/tradfri"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.WarnLevel)

	var clientPsk string
	var reauth bool

	flag.StringVar(&clientPsk, "psk", "", "Specify the psk of your device, this info is found on the back of your device")
	flag.BoolVar(&reauth, "reauth", false, "Tell the program to reauthenticate with the gateway")
	flag.Parse()

	fmt.Println("Searching gateway")

	client, err := auth.ConnectToClient("Ambilight-"+auth.RandStringBytes(6), clientPsk, reauth)
	if err != nil {
		fail(err)
	}

	// The lamps to execute ambilight on
	lamps := []int{
		65537,
		65542,
	}

	doAmbilight(client, lamps, 0, 1, 300)

}

func doAmbilight(client *tradfri.Client, lamps []int, screenID int, screenshotResolution int, speedMS int) {
	for true {
		img, err := display.GetDisplayCapture(screenID, screenshotResolution)
		if err != nil {
			fmt.Println(err)
			continue
		}

		rRaw, gRaw, bRaw, _ := img.At(0, 0).RGBA()

		r := int(rRaw / 257)
		g := int(gRaw / 257)
		b := int(bRaw / 257)

		for _, lampID := range lamps {
			_, err := client.PutDeviceColorRGBIntTimed(lampID, r, g, b, speedMS)
			if err != nil {
				fmt.Println(err)
			}
		}

		time.Sleep(time.Millisecond * time.Duration(speedMS))
	}
}

func fail(err error) {
	fmt.Println(err)
	os.Exit(1)
}
