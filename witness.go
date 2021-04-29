package see

import (
	"bytes"
	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	math_rand "math/rand"

	"github.com/rs/zerolog/log"
)

type ImageFingerprint struct {
	image  image.Image
	points []ImagePoint
}

type ImagePoint struct {
	x int
	y int
	r uint32
	g uint32
	b uint32
	a uint32
}

func CreateImageFingerprint(imageAsByteArray []byte, numberOfPoints int) (ImageFingerprint, error) {
	var err error
	var fp ImageFingerprint

	var b [8]byte
	_, err = crypto_rand.Read(b[:])
	if err == nil {

		math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

		image, format, err := image.Decode(bytes.NewReader(imageAsByteArray))
		if err == nil {
			log.Debug().
				Str("type", format).
				Msg("image created from byte array")

			var points []ImagePoint
			for i := 0; i < numberOfPoints; i++ {
				x := math_rand.Intn(image.Bounds().Dx())
				y := math_rand.Intn(image.Bounds().Dy())
				r, g, b, a := image.At(x, y).RGBA()

				point := ImagePoint{
					x: x,
					y: y,
					r: r,
					g: g,
					b: b,
					a: a,
				}

				points = append(points, point)
			}

			fp = ImageFingerprint{
				image:  image,
				points: points,
			}

			log.Info().
				Int("number of points", len(points)).
				Str("fingerprint", fmt.Sprintf("%v", fp)).
				Msg("image fingerprint created")
		}
	}

	return fp, err
}

// TODO: Understand why it is flaky. It seems to be related to the fingerprint different fingerprints have different results
// One: same colored points fail on first match and never go through the rest of the list
func FindImage(contextImageInBytes []byte, ImageFingerprint ImageFingerprint) (bool, int, int, int, int) {
	var found bool
	var x0, y0, x1, y1 int

	contextImage, format, err := image.Decode(bytes.NewReader(contextImageInBytes))

	if err == nil {
		log.Debug().
			Str("type", format).
			Msg("Context image created from byte array")

		for x := 0; x < contextImage.Bounds().Dx(); x++ {
			for y := 0; y < contextImage.Bounds().Dy(); y++ {
				points := matchPoints(x, y, contextImage, ImageFingerprint)
				for i := range points {
					point := points[i]
					if matchFingerprint(contextImage, ImageFingerprint, x, y, point.x, point.y) {
						found = true
						x0 = x - point.x
						x1 = x0 + ImageFingerprint.image.Bounds().Dx()
						y0 = y - point.y
						y1 = y0 + ImageFingerprint.image.Bounds().Dy()

						log.Debug().
							Str("coordinates", fmt.Sprintf("%v,%v:%v,%v", x0, y0, x1, y1)).
							Msg("fingerprint match found")

						return found, x0, y0, x1, y1
					}
				}
			}
		}
	}

	return found, x0, y0, x1, y1
}

func matchPoints(x, y int, contextImage image.Image, imageFingerprint ImageFingerprint) []ImagePoint {
	var matchingPoints []ImagePoint

	for i := range imageFingerprint.points {
		point := imageFingerprint.points[i]
		r, g, b, a := contextImage.At(x, y).RGBA()
		if r == point.r && g == point.g && b == point.b && a == point.a {
			log.Debug().Msgf("match found between context coordinate '%v,%v' and image fingerprint point '%v'", x, y, point)
			matchingPoints = append(matchingPoints, point)
		}
	}

	return matchingPoints
}

func matchFingerprint(contextImage image.Image, imageFingerprint ImageFingerprint, contextX, contextY, imageX, imageY int) bool {
	match := true

	for i := range imageFingerprint.points {
		point := imageFingerprint.points[i]
		xTranslated := contextX - (imageX - point.x)
		yTranslated := contextY - (imageY - point.y)

		if xTranslated >= 0 && yTranslated >= 0 {
			r, g, b, a := contextImage.At(xTranslated, yTranslated).RGBA()
			match = (r == point.r && g == point.g && b == point.b && a == point.a)
			if !match {
				break
			}
		} else {
			match = false
			break
		}
	}

	if !match {
		log.Info().
			Str("imageFingerprint", fmt.Sprintf("%v", imageFingerprint.points)).
			Str("matchOrigin", fmt.Sprintf("contextX[%v]:imageX[%v],contextY[%v]:imageY[%v]", contextX, imageX, contextY, imageY)).
			Msg("No match found for fingerprint")
	}

	return match
}

////
// These 'set' functions should not be necessary and are very bad for the performance
// TODO: make it not flaky and remove 'set' functions
////

func CreateImageFingerprintSet(imageAsByteArray []byte, numberOfPoints int, setSize int) []ImageFingerprint {
	var imageFingerprintSet []ImageFingerprint

	for i := 0; i < setSize; i++ {
		fp, err := CreateImageFingerprint(imageAsByteArray, numberOfPoints)
		if err == nil {
			imageFingerprintSet = append(imageFingerprintSet, fp)
		}

	}

	return imageFingerprintSet
}

func FindImageBasedOnSet(contextImageInBytes []byte, imageFingerprintSet []ImageFingerprint) (bool, int, int, int, int) {
	var found bool
	var x0, y0, x1, y1 int

	for i := range imageFingerprintSet {
		found, x0, y0, x1, y1 = FindImage(contextImageInBytes, imageFingerprintSet[i])
		if found {
			break
		}

		log.Info().
			Int("image fingerprint", i+1).
			Int("amount of fingerprints", len(imageFingerprintSet)).
			Msg("No image found based on image fingerprint")
	}

	return found, x0, y0, x1, y1
}
