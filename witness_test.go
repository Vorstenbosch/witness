package witness

import (
	"bytes"
	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	math_rand "math/rand"
	"os"
	"testing"
)

// In case you want to view the test images
const writeTestImagesToFile = true

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

func TestImageFindingWithRandomFingerprint(t *testing.T) {
	// Given
	randomSeed()
	contextImage := createContextImage()
	imageToFind := createImageToFind()

	fp, _ := CreateImageFingerprint(imageToFind, 4)

	// When
	found, rectangle := FindImage(contextImage, fp)

	// Then
	if !found {
		t.Errorf("Could not find image")
	} else {
		if rectangle.Min.X != 0 || rectangle.Min.Y != 3 || rectangle.Max.X != 5 || rectangle.Max.Y != 4 {
			t.Errorf("Matched on unexpected rectangle '%v'", rectangle)
		}
	}
}

func createContextImage() image.Image {
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

	if writeTestImagesToFile {
		writeImageToFile(contextImage, "contextImage")
	}

	return contextImage
}

func createImageToFind() image.Image {
	imageToFind := image.NewRGBA(image.Rect(0, 0, 5, 1))

	imageToFind.Set(0, 0, blackPixel)
	imageToFind.Set(1, 0, whitePixel)
	imageToFind.Set(2, 0, blackPixel)
	imageToFind.Set(3, 0, whitePixel)
	imageToFind.Set(4, 0, blackPixel)

	if writeTestImagesToFile {
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

func randomSeed() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])

	if err != nil {
		panic("Unable to setup the random seed")
	}
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}
