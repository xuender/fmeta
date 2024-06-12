package meta

import (
	"context"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/dsoprea/go-exif/v3"
	"github.com/h2non/filetype"
	"github.com/samber/lo"
	"gopkg.in/vansante/go-ffprobe.v2"
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

	val, err := getMetaByReader(file)
	if err != nil {
		return val, err
	}

	val.Datetime = info.ModTime().Format(time.DateTime)
	val.Size = info.Size()

	if isMedia(val.GetType()) {
		if err := readMedia(file, val); err != nil {
			return val, err
		}
	}

	if val.GetType() != MetaType_Image {
		return val, err
	}

	return readImage(file, val)
}

func readMedia(file *os.File, val *Meta) error {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if data, err := ffprobe.ProbeReader(ctx, file); err == nil {
		for _, stream := range data.Streams {
			if stream.Width > 0 && stream.Height > 0 {
				val.Width = int32(stream.Width)
				val.Height = int32(stream.Height)
			}

			if stream.Duration != "" {
				val.Duration, _ = strconv.ParseFloat(stream.Duration, 64)
			}

			if stream.Channels > 0 {
				val.Channels = int32(stream.Channels)
			}
		}
	}

	return nil
}

func isMedia(mediaType MetaType) bool {
	return mediaType == MetaType_Image ||
		mediaType == MetaType_Audio ||
		mediaType == MetaType_Video
}

func readImage(file *os.File, val *Meta) (*Meta, error) {
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

	if val.GetWidth() > 0 && val.GetHeight() > 0 {
		return val, nil
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return val, err
	}

	img, err := imaging.Decode(file)
	if err != nil {
		return val, err
	}

	val.Width, val.Height = int32(img.Bounds().Dx()), int32(img.Bounds().Dy())

	return val, err
}

func getMetaByReader(reader io.Reader) (*Meta, error) {
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
