package main

import (
	"os"
	"fmt"
	"syscall"
	"time"
	"log"
	"strconv"
	"strings"
	"errors"
	)

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

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var fileset ImgFileSet

	for _, fi := range files {
		cur_file := Create(fi)
		fileset = fileset.Add(cur_file)
	}

	fileset.Print()
	fileset.Rename(dir)

}


/// class ImgFileSet
type ImgFileSet struct {
	set map[int][]ImgFile
	descr_idx map[string]map[string]int // date: [descr1: idx1, descr2: idx2] 
}

// @TODO change name after renaming in ImgFileSet ?
func (ifs ImgFileSet) Rename(directory string) {
	if err := ifs.checkFutureNamesUniqueness(directory);
	err != nil {
		log.Fatal(err)
	}

	for key, files := range ifs.set {
		for _, file := range files {
			newfilename := file.TidyName("", key)
			if _, err := os.Stat(directory + "/" + newfilename);
			err == nil {
				//file exist
				fmt.Println("Rename: the file already exists: ", newfilename )
				continue
			} else if os.IsNotExist(err) {
				// our case, do the stuff
				os.Rename(
					directory + "/" + file.full_name, 
					directory + "/" + newfilename)
			} else {
				// something else happened
				panic(err)
			}
		}
	}
	fmt.Println("Renaming done!")	
}

// @TODO check that also the new names are unique from each other
func (ifs ImgFileSet) checkFutureNamesUniqueness(directory string) error{
	futureNames := make(map[string]string) // key is new filename, value is old filename

	for key, files := range ifs.set {
		for _, file := range files {
			newfilename := file.TidyName("", key)
			if _, err := os.Stat(directory + "/" + newfilename);
			err == nil {
				//file exist
				return errors.New("nameUniqueness: file \"" + file.full_name + "\" was about to be renamed to \"" + newfilename + "\" but it already exist, script aborted")
			} else if _, val :=futureNames[newfilename]; val {
				return errors.New("nameUniqueness: file \"" + file.full_name + "\" has a future name collision with \"" + futureNames[newfilename] + "\" , script aborted")
			} else if os.IsNotExist(err) {
				futureNames[newfilename] = file.full_name
				continue
			} else {
				// something else happened
				panic(err)
			}
		}
	}
	return nil
}

func (fileset ImgFileSet) Add(file ImgFile) ImgFileSet {
	if fileset.set == nil {
		fmt.Println("fileset map empty, filling...") //@TODO logging
		fileset.set = make(map[int][]ImgFile)
	}
	

	fileset.set[file.suffix] = append(fileset.set[file.suffix], file)

	return fileset
}

func (fileset ImgFileSet) Print() {
	fmt.Println("Fileset tree:")
	for common_name, file_slice := range fileset.set {
		
		fmt.Println(common_name)
		for _, val := range file_slice {
			fmt.Println("\t ", val.Print())
		}

	}
}



/// class ImgFile

// @TODO check that is not dir 

type ImgFile struct {
	name string
	prefix string
	suffix int
	extension string
	full_name string
	generated_date string //time.Time
}

func (file ImgFile) Print() string {
	return file.full_name
}

func Create(path os.DirEntry) ImgFile {

	fileinfo, err := path.Info()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fname := fileinfo.Name()
	fprefix := fname[:strings.LastIndex(fname, ".")][:4]
	fsuffix, _ := strconv.Atoi(fname[:strings.LastIndex(fname, ".")][4:])

	newFile := ImgFile{ 
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

func (file ImgFile) TidyName(term string, idx int) string {
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