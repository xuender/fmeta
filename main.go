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
		val, err := meta.FileMeta(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", arg, err)

			continue
		}

		fmt.Fprintf(os.Stdout, "%s: %v\n", arg, val.GetType())

		switch val.GetType() {
		case meta.MetaType_Document, meta.MetaType_Image:
			subtype(val)
		case meta.MetaType_Directory, meta.MetaType_Archive, meta.MetaType_Audio, meta.MetaType_Unknown, meta.MetaType_Video:
			other(val)
		default:
			other(val)
		}
	}
}

func other(val *meta.Meta) {
	fmt.Fprintf(os.Stdout, "DateTime: %s\n", val.GetDatetime())
}

func subtype(val *meta.Meta) {
	fmt.Fprintf(os.Stdout, "Subtype: %s\n", val.GetSubtype())
	fmt.Fprintf(os.Stdout, "Extension: %s\n", val.GetExtension())
	fmt.Fprintf(os.Stdout, "DateTime: %s\n", val.GetDatetime())
}

func usage() {
	fmt.Fprintf(os.Stderr, "fmeta\n\n")
	fmt.Fprintf(os.Stderr, "File Meta.\n\n")
	fmt.Fprintf(os.Stderr, "Usage: %s [flags] file\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}
