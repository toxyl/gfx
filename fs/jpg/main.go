package jpg

import (
	"bytes"
	"image"
	"image/jpeg"
	"os"
)

func Save(img image.Image, path string) {
	outFile, err := os.Create(path) // #nosec G304
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	err = jpeg.Encode(outFile, img, nil)
	if err != nil {
		panic(err)
	}
}

func FromFile(filename string) (image.Image, error) {
	file, err := os.Open(filename) // #nosec G304
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return jpeg.Decode(file)
}

func FromBytes(data []byte) (image.Image, error) {
	return jpeg.Decode(bytes.NewReader(data))
}
