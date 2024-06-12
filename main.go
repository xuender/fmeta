package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xuender/fmeta/meta"
)

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		usage()

		return
	}

	for _, arg := range flag.Args() {
		readPath(arg)
	}
}

func readPath(path string) {
	val, err := meta.FileMeta(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)

		return
	}

	fmt.Fprintf(os.Stdout, "%s: %v\n", path, val.GetType())

	switch val.GetType() {
	case meta.MetaType_Document:
		subtype(val)
		file(val)
	case meta.MetaType_Audio:
		subtype(val)
		audio(val)
		file(val)
	case meta.MetaType_Video:
		subtype(val)
		video(val)
		file(val)
	case meta.MetaType_Image:
		subtype(val)
		image(val)
		file(val)
	case meta.MetaType_Directory:
		dir(val)
	case meta.MetaType_Archive, meta.MetaType_Unknown:
		file(val)
	default:
		file(val)
	}
}

func dir(val *meta.Meta) {
	fmt.Fprintf(os.Stdout, "DateTime: %s\n", val.GetDatetime())
}

func file(val *meta.Meta) {
	fmt.Fprintf(os.Stdout, "Size: %d\n", val.GetSize())
	dir(val)
}

func audio(val *meta.Meta) {
	fmt.Fprintf(os.Stdout, "Duration: %.3f\n", val.GetDuration())
}

func video(val *meta.Meta) {
	image(val)
	audio(val)
}

func image(val *meta.Meta) {
	fmt.Fprintf(os.Stdout, "Width: %d\n", val.GetWidth())
	fmt.Fprintf(os.Stdout, "Height: %d\n", val.GetHeight())
}

func subtype(val *meta.Meta) {
	fmt.Fprintf(os.Stdout, "Subtype: %s\n", val.GetSubtype())
	fmt.Fprintf(os.Stdout, "Extension: %s\n", val.GetExtension())
}

func usage() {
	fmt.Fprintf(os.Stderr, "fmeta\n\n")
	fmt.Fprintf(os.Stderr, "File Meta.\n\n")
	fmt.Fprintf(os.Stderr, "Usage: %s [flags] file\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}
