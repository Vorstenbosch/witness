package witness

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

// In case you want to view the test images
const writeTestImagesToFile = false

func TestFingerprintCreation(t *testing.T) {
	// Given
	imageToFind := createImageToFind()

	// When
	fp, err := CreateImageFingerprint(imageToFind, 5)

	// Then
	if err != nil {
		t.Errorf("Unexpected error '%v' occured", err)
	}

	if len(fp.pixels) != 5 {
		t.Errorf("Unexpected amount of points ('%v') found ", len(fp.pixels))
	}

	if fp.image == nil {
		t.Errorf("Image was nil")
	}
}

func TestImageFinding(t *testing.T) {
	// Given
	contextImage := createContextImage()
	imageToFind := createImageToFind()

	fp, _ := CreateImageFingerprint(imageToFind, 5)

	// When
	found, _ := FindImage(contextImage, fp)

	// Then
	if !found {
		t.Errorf("Could not find image")
	}

}

func createContextImage() image.Image {
	var whitePixel = color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}

	var blackPixel = color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}

	contextImage := image.NewRGBA(image.Rect(0, 0, 5, 5))

	for x := 0; x < contextImage.Bounds().Dx(); x++ {
		for y := 0; y < contextImage.Bounds().Dy(); y++ {
			if y == 3 && x%2 == 0 {
				contextImage.Set(x, y, blackPixel)
			} else {
				contextImage.Set(x, y, whitePixel)
			}
		}
	}

	return contextImage
}

func createImageToFind() image.Image {
	contextImage := createContextImage()

	imageToFind := contextImage.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(0, 3, 5, 4))

	if writeTestImagesToFile {
		writeImageToFile(contextImage, "contextImage")
		writeImageToFile(imageToFind, "imageToFind")
	}

	return imageToFind
}

func writeImageToFile(i image.Image, imageName string) {
	b := new(bytes.Buffer)
	png.Encode(b, i)

	f, _ := os.Create(fmt.Sprintf("%s.png", imageName))
	defer f.Close()
	f.Write(b.Bytes())
}
