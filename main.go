package main

import (
	"os"
	"fmt"
	"syscall"
	"time"
	// "log"
	// "strconv"
	"strings"
	)

// TODO flag to set date string manually (creation date is not copied)

func main() {
	dir := `/home/tonno/Desktop/20220516/golang_test`

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var ara_slice []File
	for _, fi := range files {
		cur_file := create_file_class(fi)

		ara_slice = append(ara_slice, cur_file)
	}

	for _, el := range ara_slice {
		fmt.Println(el.full_name)
	}
}

func Contains(sl []string, name string) bool {
	for _, value := range sl {
	   if value == name {
		  return true
	   }
	}
	return false
 }


// func main() {
// 	dir := `/home/tonno/Desktop/20220516/golang_test`

// 	files, err := os.ReadDir(dir)
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 		os.Exit(1)
// 	}

// 	i := 0
//  	for _, fi := range files {
// 		i += 1
// 		idx := strconv.Itoa(i)
		

// 		cur_file := create_file_class(fi)
// 		e := os.Rename(dir + "/" + cur_file.full_name, 
// 					   dir + "/" + cur_file.generated_date + "_" + idx + "." + cur_file.extension ) // 
// 		if e != nil {
// 			log.Fatal(e)
// 		}
// 	}
//    }
   

func unix_filetime(path os.FileInfo) string {
	// source: https://developpaper.com/getting-access-creation-modification-time-of-files-on-linux-using-golang/ //

	// Sys () returns interface {}, so you need a type assertion. Different platforms need different types. On linux, * syscall. Stat_t
	stat_t := path.Sys().(*syscall.Stat_t)
	//fmt.Println(stat_t.Ctim.Sec)
	
	unix_time := time.Unix(int64(stat_t.Mtim.Sec), int64(stat_t.Mtim.Nsec) )
	//fmt.Println( unix_time )
	//fmt.Println( unix_time.Format("20060102") ) //.Format("2006-01-02")

	return unix_time.Format("20060115")

}
   

/// class FileSet
type FileSet struct {
	set map[string][]File
}

func add(fileset FileSet, file File) FileSet {
	fileset.set[file.name] = append(fileset.set[file.name], file)

	return fileset
}

/// class File

// @TODO check that is not dir 

type File struct {
	name string
	extension string
	full_name string
	generated_date string //time.Time
}

func create_file_class(path os.DirEntry) File {

	fileinfo, err := path.Info()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fname := fileinfo.Name()

	newFile := File{ 
		name: fname[:strings.LastIndex(fname, ".")],
		extension: fname[strings.LastIndex(fname, ".")+1:],
		full_name: fname,
		generated_date: unix_filetime(fileinfo),
	}

	print(newFile.name, "\n")
	print(newFile.full_name, "\n")

	return newFile
}