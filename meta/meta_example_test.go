package meta_test

import (
	"fmt"

	"github.com/xuender/fmeta/meta"
)

func ExampleFileMeta() {
	val, err := meta.FileMeta("green.jpg")

	fmt.Println(err)
	fmt.Println(val.GetType())

	// Output:
	// <nil>
	// Image
}

func ExampleDirMeta() {
	val, err := meta.DirMeta(".")

	fmt.Println(err)
	fmt.Println(val.GetType())

	_, err = meta.DirMeta("nofound")
	fmt.Println(err)

	// Output:
	// <nil>
	// Directory
	// stat nofound: no such file or directory
}

func ExampleDateFormat() {
	fmt.Println(meta.DateFormat("2006/01/02 15:04:05"))
	fmt.Println(meta.DateFormat("2006:01:02 15:04:05"))

	// Output:
	// 2006-01-02 15:04:05
	// 2006-01-02 15:04:05
}
