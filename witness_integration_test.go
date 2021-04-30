package witness

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rs/zerolog"
)

// If you want to view the image found
const writeFoundImageToFile = true

func TestScreenCaptureFinding(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	randomSeed()

	// Given
	b, _ := ioutil.ReadFile("test-images/brave-icon.png")
	c, _ := ioutil.ReadFile("test-images/icons.png")
	imageToFind, _, _ := image.Decode(bytes.NewReader(b))
	contextImage, _, _ := image.Decode(bytes.NewReader(c))
	fp, _ := CreateImageFingerprint(imageToFind, 7)

	// When
	found, rectangle := FindImage(contextImage, fp)

	// Then
	if !found {
		t.Errorf("Failed to match")
	} else {
		if rectangle.Min.X != 24 && rectangle.Min.Y != 508 && rectangle.Max.X != 154 && rectangle.Max.Y != 154 {
			t.Errorf("Matched on unexpected coordinate")
		}
	}

	if found && writeFoundImageToFile {
		cImage, _, _ := image.Decode(bytes.NewReader(c))
		foundImage := cImage.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(rectangle)

		foundImageBytes := new(bytes.Buffer)

		err := png.Encode(foundImageBytes, foundImage)
		if err != nil {
			fmt.Println("failed to create buffer", err)
		}

		f, _ := os.Create("test-images/found-image.png")
		defer f.Close()
		f.Write(foundImageBytes.Bytes())
	}
}
