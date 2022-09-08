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
	intervalStr = flag.String("interval", "655,657-659,1631,1640-1645", "the range of progressive index files that will be renamed")
}

func main() {
	flag.Parse()

	files, err := os.ReadDir(*dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var fileset imgfile.ImgFileSet

	for _, fi := range files {
		cur_file := imgfile.New(fi, *file_prefix)
		if !cur_file.IsEmpty() {
			fileset.Add(*cur_file)
		}
	}
	fileset.RecoverEditorFiles()
	fileset.String()

	interval := ProcessSeq(*intervalStr)

	fileset.Rename(*dir, interval, *label)
}




//////////////////////////////
//////////////////////////////
//////////////////////////////
//////////////////////////////



func ProcessSeq(seq string) []int {
	var full_seq []int
	
	for _, val := range strings.Split(seq, ",") {
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
					panic("Lower bound is set higher than higher bound!")
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