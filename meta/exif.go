package meta

import (
	"io"

	"github.com/dsoprea/go-exif/v3"
)

func ImageExif(reader io.Reader) ([]exif.ExifTag, error) {
	rawExif, err := exif.SearchAndExtractExifWithReader(reader)
	if err != nil {
		return nil, err
	}

	entries, _, err := exif.GetFlatExifData(rawExif, nil)

	return entries, err
}
