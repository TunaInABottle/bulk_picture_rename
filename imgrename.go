package main

import (
	"os"
	"fmt"
	"flag"

	"file_rename/imgfile"
	)


var (
	file_prefix  *string
	port *int
)

func init() {
	file_prefix = flag.String("prefix", "TUNA", "prefix of the files to be renamed")
}

// @TODO flag to set date string manually (creation date is not copied)
// @TODO flag to override name prefix
// @TODO flag to override name prefix length
// @TODO flag for interval
// @TODO prefix check
// @TODO strategy to prevent renaming if a file cant be renamed

// example: main.go --interval:400-600,606 --name:landscape
// renames files from 400 to 600 and 606 by putting "landscape" in the name

func main() {
	dir := `/home/tonno/Desktop/20220516/golang_test`
	flag.Parse()

	fmt.Println("Flag prefix is ", *file_prefix)

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var fileset imgfile.ImgFileSet

	for _, fi := range files {
		cur_file := imgfile.Create(fi, *file_prefix)

		if !cur_file.IsEmpty() {
			fileset = fileset.Add(*cur_file)
		}
	}

	fileset.Print()
	fileset.Rename(dir)

}



//////////////////////////////
//////////////////////////////
//////////////////////////////
//////////////////////////////


func Contains(sl []string, name string) bool {
	for _, value := range sl {
	   if value == name {
		  return true
	   }
	}
	return false
 }