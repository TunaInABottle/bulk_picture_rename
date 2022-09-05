package imgfile

import (
	"os"
	"fmt"
	"syscall"
	"time"
	"strconv"
	"strings"
	"errors"
	
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	)

var editor_extensions []string = []string{".ORF.dop", ".JPG.dop"}
var file_extensions []string = append([]string{".ORF", ".JPG"}, editor_extensions...)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
}

////////////////////////////////////
///////// class ImgFileSet /////////
////////////////////////////////////

type ImgFileSet struct {
	//set map[int][]ImgFile
	set 		 map[string][]ImgFile
	date_idx_set map[string]map[string][]ImgFile
	descr_idx	 map[string]map[string]int // date: [description1: idx1, descr2: idx2] 
	idx_date 	 map[string]string					// imgidx1: date
}

// @TODO change name after renaming in ImgFileSet ?
func (ifs ImgFileSet) Rename(directory string, interval string, label string) {
	if err := ifs.checkFutureNamesUniqueness(directory);
	err != nil {
		log.Fatal().Err(err)
	}

	for _, files := range ifs.set {
		for key, file := range files {
			newfilename := file.TidyName(label, key)
			if _, err := os.Stat(directory + "/" + newfilename);
			err == nil {
				//file exist
				log.Warn().Str("old_name", file.full_name).Str("new_name", newfilename).Msg("Rename: file name already exists")
				continue
			} else if os.IsNotExist(err) {
				log.Debug().Str("dir", directory).Str("old_name", file.full_name).Str("new_name", newfilename).Msg("renaming file")
				os.Rename(
					directory + "/" + file.full_name, 
					directory + "/" + newfilename)
			} else {
				// something else happened
				panic(err)
			}
		}
	}
	log.Info().Msg("Renaming done!")	
}

// @TODO check that also the new names are unique from each other
func (ifs ImgFileSet) checkFutureNamesUniqueness(directory string) error{
	futureNames := make(map[string]string) // key is new filename, value is old filename

	for key, files := range ifs.set {
		for _, file := range files {
			newfilename := file.TidyName1("", key)
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
	if file.IsEmpty() {
		log.Debug().Msg("Empty file, continuing")
		return fileset
	}

	fileset = fileset.createMissingMaps(file)
	
	// add to date_idx_set map BEGIN
	slice_to_append := fileset.date_idx_set[file.generated_date][strconv.Itoa(file.suffix)]
	fileset.date_idx_set[file.generated_date][strconv.Itoa(file.suffix)] = append(slice_to_append, file)
	// add to date_idx_set map END

	// add to idx_date map BEGIN
	doIt := true
	for _, val := range editor_extensions {
		if file.extension == val {
			doIt = false
		}
	}
	if doIt && fileset.idx_date[strconv.Itoa(file.suffix)] == "" {
		fileset.idx_date[strconv.Itoa(file.suffix)] = file.generated_date
	}
	// add to idx_date map END

	return fileset
}
func (fileset ImgFileSet) Print() {
	fmt.Println("ImgFileSet tree:")
	for common_date, scd_map := range fileset.date_idx_set {
		fmt.Println(common_date)
		for common_idx, file_slice := range scd_map{
			fmt.Println("  ", common_idx)
			// fmt.Println("   ", scd_map)
			for _, val := range file_slice {
				fmt.Println("    ", val.Print())
			}
		}
	}
	fmt.Println(fileset.idx_date)
}


func (fileset ImgFileSet) createMissingMaps(file ImgFile) ImgFileSet {
	// main date_idx_map
	if fileset.date_idx_set == nil {
		log.Debug().Msg("@TODO better text, fileset.date_idx_set empty, creating...")
		fileset.date_idx_set = make(map[string]map[string][]ImgFile)
	}
	// nested date map in date_idx_set
	if fileset.date_idx_set[file.generated_date] == nil {
		log.Debug().Str("date", file.generated_date).Msg("fileset fileset.date_idx_set of date empty, creating...")
		fileset.date_idx_set[file.generated_date] = make(map[string][]ImgFile)
	}
	// map in idx_date
	if fileset.idx_date == nil {
		log.Debug().Msg("@TODO better text, fileset map empty, creating...")
		fileset.idx_date = make(map[string]string)
	}

	return fileset
}
/////////////////////////////////
///////// class ImgFile /////////
/////////////////////////////////

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

func Create(path os.DirEntry, prefix string) *ImgFile {
	fileinfo, err := path.Info()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fname := fileinfo.Name()

	if path.IsDir() {
		log.Warn().Str("name", fname).Msg("directory detected, continuing")
		return new(ImgFile)
	}

	var fExtension string
	for _, val := range file_extensions {
		if strings.HasSuffix(fname, val) {
			fExtension = val
		}
	}
	if fExtension == "" {
		log.Warn().Str("file", fname).Msg("unexpected file extension, continuing")
		return new(ImgFile)
	}


	// skip if the file does not have the expected prefix
	if !strings.HasPrefix(fname, prefix) {
		log.Warn().Str("file", fname).Msg("prefix different from expected, continuing")
		return new(ImgFile)
	}

	temp_str := strings.TrimPrefix(fname, prefix)
	fsuffix, _ := strconv.Atoi(strings.TrimSuffix(temp_str, fExtension))


	newFile := ImgFile{ 
		name: fname[:strings.LastIndex(fname, ".")],
		prefix: prefix,
		suffix: fsuffix,
		extension: fExtension,
		full_name: fname,
		generated_date: unix_filetime(fileinfo),
	}

	return &newFile
}

func (file ImgFile) TidyName1(term string, idx string) string {
	return file.generated_date + "_" + term + idx + file.extension
}

func (file ImgFile) TidyName(term string, idx int) string {
	return file.generated_date + "_" + term + strconv.Itoa(idx) + file.extension
}

func (file ImgFile) IsEmpty() bool {
	return file.name == "" && file.extension == ""
}

func unix_filetime(path os.FileInfo) string {
	// source: https://developpaper.com/getting-access-creation-modification-time-of-files-on-linux-using-golang/ //

	// Sys () returns interface {}, so you need a type assertion. Different platforms need different types. On linux, * syscall. Stat_t
	stat_t := path.Sys().(*syscall.Stat_t)
	//fmt.Println(stat_t.Ctim.Sec)

	unix_time := time.Unix(int64(stat_t.Mtim.Sec), int64(stat_t.Mtim.Nsec) )
	// fmt.Println( unix_time )
	// fmt.Println( unix_time.Format("20060102") ) //.Format("2006-01-02")

	return unix_time.Format("20060115")

}
