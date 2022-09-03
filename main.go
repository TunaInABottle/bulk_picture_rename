package main

import (
	"os"
	"fmt"
	"syscall"
	"time"
	"log"
	"strconv"
	"strings"
	)

// TODO flag to set date string manually (creation date is not copied)

func main() {
	//dir, _ := os.Stat(`/home/tonno/Desktop/20220516/Cuna Island`) // /TUNA1159.ORF
	dir := `/home/tonno/Desktop/20220516/golang_test` // /TUNA1159.ORF

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	i := 0
 	for _, fi := range files {
		i += 1
		fileinfo, err := fi.Info()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		
		ph := create_file_class(fi)

		fmt.Println( ph.name, " I did it" ) 


		filetime := unix_filetime(fileinfo)    
		
		//fmt.Println( fileinfo.Sys() )
		
		idx := strconv.Itoa(i)

		fmt.Println( filetime, idx ) 

		e := os.Rename(dir + "/" + fileinfo.Name(), dir + "/" + filetime + "_" + idx + ".ORF" ) // 
		if e != nil {
			log.Fatal(e)
		}
	}
   }
   

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