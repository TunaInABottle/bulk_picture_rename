package main

import (
	"os"
	"fmt"
	"syscall"
	"time"
	// "log"
	"strconv"
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

	var fileset FileSet

	for _, fi := range files {
		cur_file := Create_file_class(fi)

		// ara_slice = append(ara_slice, cur_file)
		fileset = fileset.Add(cur_file)
	}

	fileset.Print()

	fileset.Rename(dir)

	fmt.Println("Renaming done!")

}




/// class FileSet
type FileSet struct {
	set map[int][]File
}

// TODO change name after renaming in FileSet ?
func (fs FileSet) Rename(directory string) {
	for key, files := range fs.set {

		for _, file := range files{
			// @TODO check if renaming file that exists23
			os.Rename(
				directory + "/" + file.full_name, 
				directory + "/" + file.TidyName("", key))
		}

	}	
}

func (fileset FileSet) Add(file File) FileSet {
	if fileset.set == nil {
		fmt.Println("fileset map empty, filling...") //@TODO logging
		fileset.set = make(map[int][]File)
	}
	

	fileset.set[file.suffix] = append(fileset.set[file.suffix], file)

	return fileset
}

func (fileset FileSet) Print() {
	fmt.Println("Fileset tree:")
	for common_name, file_slice := range fileset.set {
		
		fmt.Println(common_name)
		for _, val := range file_slice {
			fmt.Println("\t ", val.Print())
		}

	}
}



/// class File

// @TODO check that is not dir 

type File struct {
	name string
	prefix string
	suffix int
	extension string
	full_name string
	generated_date string //time.Time
}

func (file File) Print() string {
	return file.full_name
}

func Create_file_class(path os.DirEntry) File {

	fileinfo, err := path.Info()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fname := fileinfo.Name()
	fprefix := fname[:strings.LastIndex(fname, ".")][:4]
	fsuffix, _ := strconv.Atoi(fname[:strings.LastIndex(fname, ".")][4:])

	newFile := File{ 
		name: fname[:strings.LastIndex(fname, ".")],
		prefix: fprefix,
		suffix: fsuffix,
		extension: fname[strings.LastIndex(fname, ".")+1:],
		full_name: fname,
		generated_date: unix_filetime(fileinfo),
	}

	//print(newFile.name, "\n")
	//print(newFile.full_name, "\n")

	return newFile
}

func (file File) TidyName(term string, idx int) string {
	return file.generated_date + "_" + term + strconv.Itoa(idx) + "." + file.extension
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