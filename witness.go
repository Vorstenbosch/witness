package witness

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	math_rand "math/rand"

	"github.com/rs/zerolog/log"
)

type ImageFingerprint struct {
	image  image.Image
	pixels []pixel
}

type pixel struct {
	point image.Point
	color color.Color
}

// CreateImageFingerprint creates a fingerprint of 'image' based on 'numberOfPixels'
// FIXME: it is now possible for to get multiple pixels with the same coordinates
func CreateImageFingerprint(imageToFind image.Image, numberOfPixels int) (ImageFingerprint, error) {
	var err error
	var fp ImageFingerprint

	randomSeed()

	var pixels []pixel
	for i := 0; i < numberOfPixels; i++ {
		x := math_rand.Intn(imageToFind.Bounds().Dx())
		y := math_rand.Intn(imageToFind.Bounds().Dy())

		pixel := pixel{
			point: image.Point{x, y},
			color: imageToFind.At(x, y),
		}

		pixels = append(pixels, pixel)
	}

	fp = ImageFingerprint{
		image:  imageToFind,
		pixels: pixels,
	}

	log.Info().
		Int("number of points", len(pixels)).
		Str("fingerprint", fmt.Sprintf("%v", fp)).
		Msg("image fingerprint created")

	return fp, err
}

// FIXME: is a bit flaky based on what fingerprint is used
func FindImage(contextImage image.Image, ImageFingerprint ImageFingerprint) (bool, image.Rectangle) {
	var found bool
	var rectangle image.Rectangle

	for x := 0; x < contextImage.Bounds().Dx(); x++ {
		for y := 0; y < contextImage.Bounds().Dy(); y++ {
			color := contextImage.At(x, y)

			// Collect the set of pixels that match the color of this pixel of the context image
			matchedPixels := matchPixelColor(color, ImageFingerprint)
			for i := range matchedPixels {
				pixel := matchedPixels[i]

				// Verify if we can match the entire fingerprint
				//   from this base pixel from the context image
				if matchFingerprint(contextImage, ImageFingerprint, x, y, pixel) {
					found = true
					x0 := x - pixel.point.X
					x1 := x0 + ImageFingerprint.image.Bounds().Dx()
					y0 := y - pixel.point.Y
					y1 := y0 + ImageFingerprint.image.Bounds().Dy()

					rectangle = image.Rect(x0, y0, x1, y1)

					log.Debug().
						Str("rectangle", rectangle.String()).
						Msg("fingerprint match found")

					return found, rectangle
				}
			}
		}
	}

	return found, rectangle
}

func equals(c1, c2 color.Color) bool {
	c1r, c1g, c1b, _ := c1.RGBA()
	c2r, c2g, c2b, _ := c2.RGBA()

	// Alpha values seem to differ, not comparing them
	return c1r == c2r && c1g == c2g && c1b == c2b
}

func matchPixelColor(c color.Color, imageFingerprint ImageFingerprint) []pixel {
	var matchingPoints []pixel

	for i := range imageFingerprint.pixels {
		pixel := imageFingerprint.pixels[i]
		color := pixel.color

		if equals(color, c) {
			matchingPoints = append(matchingPoints, pixel)
		}
	}

	return matchingPoints
}

func matchFingerprint(contextImage image.Image, imageFingerprint ImageFingerprint, contextX, contextY int, matchedPixel pixel) bool {
	match := true

	for i := range imageFingerprint.pixels {
		pixel := imageFingerprint.pixels[i]

		// Translate the coordinates based on the matching pixel, the current pixel and the context image coordinates
		xTranslated := contextX - (matchedPixel.point.X - pixel.point.X)
		yTranslated := contextY - (matchedPixel.point.Y - pixel.point.Y)

		match = xTranslated >= 0 && yTranslated >= 0 && equals(contextImage.At(xTranslated, yTranslated), pixel.color)
	}

	return match
}

func randomSeed() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])

	if err != nil {
		panic("Unable to setup the random seed")
	}
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}
