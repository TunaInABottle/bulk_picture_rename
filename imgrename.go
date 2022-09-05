package main

import (
	"os"
	"fmt"
	"flag"
	"strings"
	"regexp"
	"strconv"

	"file_rename/imgfile"
	)


var (
	file_prefix  *string
	intervalStr *string
	label *string
	dir *string
)

func init() {
	dir = flag.String("dir", "", "folder containing the pictures")
	label = flag.String("label", "", "name of the new files")
	file_prefix = flag.String("prefix", "TUNA", "prefix of the files to be renamed")
	intervalStr = flag.String("interval", "656-658,1159-1160", "the range of progressive index files that will be renamed")
}

func main() {
	// dir := `/home/tonno/Desktop/20220516/golang_test`
	// which_files := "656-658,1159-1160"
	// new_name := "testname"
	flag.Parse()

	fmt.Println("Flag prefix is ", *file_prefix)

	files, err := os.ReadDir(*dir)
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

	interval := ProcessSeq(*intervalStr)

	fileset.Print()
	fileset.Rename1(*dir, interval, *label)

}



//////////////////////////////
//////////////////////////////
//////////////////////////////
//////////////////////////////

func ProcessSeq(seq string) []int {
	var full_seq []int
	
	for _, val := range strings.Split(seq, ",") {
		fmt.Println(val)
		match, err := regexp.Match("[^\\d\\-]|^-|.*-.*-", []byte(val))
		if err != nil {
			panic(err)
		}
		if !match { //meaning valid interval
			if strings.Contains(val, "-") {
				splitted := strings.Split(val, "-")
				lower, _ := strconv.Atoi(splitted[0])
				higher, _ := strconv.Atoi(splitted[1])
				if lower > higher {
					panic("the lower bound is set higher than the higher bound!")
				}
				for lower <= higher {
					full_seq = append(full_seq, lower)
					lower++
				}
			} else {
				int_val, _ := strconv.Atoi(val)
				full_seq = append(full_seq, int_val)
			}
		}
	}
	return full_seq
}
	for _, value := range sl {
	   if value == name {
		  return true
	   }
	}
	return false
 }