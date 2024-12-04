package png

import (
	"bytes"
	"image"
	"image/png"
	"os"
)

func Save(img image.Image, path string) {
	outFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, img)
	if err != nil {
		panic(err)
	}
}

func FromFile(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return png.Decode(file)
}

func FromBytes(data []byte) (image.Image, error) {
	return png.Decode(bytes.NewReader(data))
}
