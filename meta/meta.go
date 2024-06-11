package meta

import (
	"io"
	"os"
	"time"

	"github.com/dsoprea/go-exif/v3"
	"github.com/h2non/filetype"
	"github.com/samber/lo"
)

const (
	_headSize = 265
)

func FileMeta(path string) (*Meta, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return DirMeta(path)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	val, err := GetMetaByReader(file)
	if err != nil {
		return val, err
	}

	val.Datetime = info.ModTime().Format(time.DateTime)

	if val.GetType() != MetaType_Image {
		return val, err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return val, err
	}

	if exifs, err := ImageExif(file); err == nil {
		lo.Find(exifs, func(exif exif.ExifTag) bool {
			if exif.TagName == "DateTimeOriginal" {
				if datetime, ok := exif.Value.(string); ok {
					val.Datetime = datetime
				}

				return true
			}

			return false
		})
	}

	return val, err
}

func GetMetaByReader(reader io.Reader) (*Meta, error) {
	head := make([]byte, _headSize)

	if _, err := reader.Read(head); err != nil {
		return nil, err
	}

	kind, err := filetype.Match(head)
	if err != nil {
		return nil, err
	}

	metaType := MetaType_Unknown

	switch {
	case filetype.IsImage(head):
		metaType = MetaType_Image
	case filetype.IsDocument(head), kind.MIME.Subtype == "pdf":
		metaType = MetaType_Document
	case filetype.IsVideo(head):
		metaType = MetaType_Video
	case filetype.IsAudio(head):
		metaType = MetaType_Audio
	case filetype.IsArchive(head):
		metaType = MetaType_Archive
	}

	return &Meta{
		Type:      metaType,
		Subtype:   kind.MIME.Subtype,
		Extension: kind.Extension,
	}, nil
}

func DirMeta(path string) (*Meta, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, ErrNotDir
	}

	return &Meta{
		Type:     MetaType_Directory,
		Datetime: info.ModTime().Format(time.DateTime),
	}, nil
}
