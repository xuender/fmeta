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
