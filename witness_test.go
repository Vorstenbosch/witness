package see

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

func TestFingerPrint(t *testing.T) {
	// Given
	b, _ := ioutil.ReadFile("test-images/brave-icon.png")

	// When
	fp, err := CreateImageFingerprint(b, 10)

	// Then
	if err != nil {
		t.Errorf("Unexpected error '%v' occured", err)
	}

	if len(fp.points) != 10 {
		t.Errorf("Unexpected amount of points ('%v') found ", len(fp.points))
	}

	if fp.image == nil {
		t.Errorf("Image was nil")
	}
}

func TestImageFinding(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Given
	b, _ := ioutil.ReadFile("test-images/brave-icon.png")
	fp, _ := CreateImageFingerprint(b, 10)
	c, _ := ioutil.ReadFile("test-images/icons.png")

	// When
	found, x0, y0, x1, y1 := FindImage(c, fp)

	// Then
	if !found {
		t.Errorf("Failed to match")
	} else {
		if x0 != 24 && y0 != 508 && x1 != 154 && y0 != 154 {
			t.Errorf("Matched on unexpected coordinate")
		}
	}

	if found {
		cImage, _, _ := image.Decode(bytes.NewReader(c))
		foundImage := cImage.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(image.Rect(x0, y0, x1, y1))

		// create buffer
		foundImageBytes := new(bytes.Buffer)

		// encode image to buffer
		err := png.Encode(foundImageBytes, foundImage)
		if err != nil {
			fmt.Println("failed to create buffer", err)
		}

		f, _ := os.Create("test-images/found-image.png")
		defer f.Close()
		f.Write(foundImageBytes.Bytes())
	}
}

func TestImageFindingBasedOnSet(t *testing.T) {
	// Given
	b, _ := ioutil.ReadFile("test-images/brave-icon.png")
	fps := CreateImageFingerprintSet(b, 10, 50)
	c, _ := ioutil.ReadFile("test-images/icons.png")

	// When
	found, x0, y0, x1, y1 := FindImageBasedOnSet(c, fps)

	// Then
	if !found {
		t.Errorf("Failed to match")
	}

	if x0 != 24 && y0 != 508 && x1 != 154 && y0 != 154 {
		t.Errorf("Matched on unexpected coordinate")
	}

	if found {
		cImage, _, _ := image.Decode(bytes.NewReader(c))
		foundImage := cImage.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(image.Rect(x0, y0, x1, y1))

		// create buffer
		foundImageBytes := new(bytes.Buffer)

		// encode image to buffer
		err := png.Encode(foundImageBytes, foundImage)
		if err != nil {
			fmt.Println("failed to create buffer", err)
		}

		f, _ := os.Create("test-images/found-image.png")
		defer f.Close()
		f.Write(foundImageBytes.Bytes())
	}

}
