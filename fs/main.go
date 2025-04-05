package fs

import "github.com/toxyl/flo"

func SaveString(path, data string) error {
	f := flo.File(path)
	if err := f.Mkparent(0775); err != nil {
		return err
	}
	if err := f.StoreString(data); err != nil {
		return err
	}
	return nil
}

func LoadString(path string) string {
	return flo.File(path).AsString()
}

func LoadBytes(path string) []byte {
	return flo.File(path).AsBytes()
}

func LoadStringInto(path string, s *string) error {
	return flo.File(path).LoadString(s)
}

func Copy(src, dst string) error {
	return flo.File(src).Copy(dst)
}
