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


var file_extensions [4]string = [4]string{".ORF", ".JPG", ".ORF.dop", ".JPG.dop"}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
}

////////////////////////
/// class ImgFileSet ///
////////////////////////

type ImgFileSet struct {
	set map[int][]ImgFile
	descr_idx map[string]map[string]int // date: [description1: idx1, descr2: idx2] 
}

// @TODO change name after renaming in ImgFileSet ?
func (ifs ImgFileSet) Rename(directory string) {
	if err := ifs.checkFutureNamesUniqueness(directory);
	err != nil {
		log.Fatal().Err(err)
	}

	for key, files := range ifs.set {
		for _, file := range files {
			newfilename := file.TidyName("", key)
			if _, err := os.Stat(directory + "/" + newfilename);
			err == nil {
				//file exist
				log.Warn().Str("file", newfilename).Msg("Rename: file name already exists")
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
	log.Info().Msg("Renaming done!")	
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
		log.Debug().Msg("fileset map empty, creating...")
		fileset.set = make(map[int][]ImgFile)
	}
	
	if file.IsEmpty() {
		log.Debug().Msg("Empty file, continuing")
		return fileset
	}

	fileset.set[file.suffix] = append(fileset.set[file.suffix], file)
	return fileset
}
func (fileset ImgFileSet) Print() {
	fmt.Println("ImgFileSet tree:")
	for common_name, file_slice := range fileset.set {
		
		fmt.Println(common_name)
		for _, val := range file_slice {
			fmt.Println("      ", val.Print())
		}

	}
}

/////////////////////
/// class ImgFile ///
/////////////////////


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
