package display

import (
	"image"

	"github.com/kbinani/screenshot"
	"github.com/nfnt/resize"
)

func GetDisplayCapture(displayID int, size int) (image.Image, error) {
	bounds := screenshot.GetDisplayBounds(displayID)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}

	resizedImage := resize.Resize(uint(size), uint(size), img, resize.Bilinear)

	return resizedImage, nil
}

func GetDisplayAmount() int {
	return screenshot.NumActiveDisplays()
}
