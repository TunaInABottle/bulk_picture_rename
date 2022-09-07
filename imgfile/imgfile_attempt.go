package imgfile

import (
	"os"
	"fmt"
	"syscall"
	"time"
	"strconv"
	"strings"
	"regexp"
	
	// "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//////////////////////////////////
///////// ImgFileFactory /////////
//////////////////////////////////

func getImgFile() (*ImgFiler, error) {
	return nil, nil
}


////////////////////////////
///////// ImgFiler /////////
////////////////////////////

type ImgFiler interface {
	String() string
	//New(os.DirEntry, string) *T
	IsEmpty() bool
	TidyName(string, int) string
}

///////////////////////////
///////// ImgFile /////////
///////////////////////////

// example of expected filenames:
// TUNA0707.JPG
// 20200906_coding_1.JPG
type ImgFile struct {
	//label string			// the labels in the file name eg. landscape_lake
	prefix string			// the first component of a filename
	enum int				// progressive number written in the filename
	extension string		// extension of the file
	full_name string		// full file name
	generated_date string	// date in which the picture has been taken
	//isLabeled bool			// tells if the image has been previously labeled
}

func (file ImgFile) String() string {
	return file.full_name
}

func New(path os.DirEntry, prefix string) *ImgFile {
	fileinfo, err := path.Info() // os.FileInfo type

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fname := fileinfo.Name()
	
	var fExtension string
	for _, val := range allowedExtensions {
		if strings.HasSuffix(fname, val) {
			fExtension = val
		}
	}
	return newImgFile(fileinfo, fname, fExtension, prefix)
}

func newImgFile(fileinfo os.FileInfo, fname string, fExtension string, prefix string) (*ImgFile) {
	if fExtension == "" {
		log.Warn().Str("file", fname).Msg("unexpected file extension, skipped")
		return new(ImgFile)
	}

	// skip if the file does not have the expected prefix
	if strings.HasPrefix(fname, prefix) {
		return newWithPrefix(fileinfo, fname, fExtension, prefix)
	} 
	if match, err := regexp.Match("^\\d{8}_", []byte(fname)); //string starts with 8 numbers
	err != nil {
		panic(err)
	} else if match {
		return newWithDate(fileinfo, fname, fExtension)
	}
	log.Warn().Str("file", fname).Msg("something went wrong with the creation")
	return new(ImgFile)
}

// expect string in a format eg. "20220102_landscape_1.JPEG"
func newWithDate(fileinfo os.FileInfo, fname string, fExtension string) (*ImgFile) {
	splitString := strings.Split(fname, "_") // [20220102, landscape, 1.JPEG]
	tempEnum := splitString[len(splitString)-1]
	fenum, _ := strconv.Atoi(tempEnum[:strings.LastIndex(tempEnum, ".")]) // img index
	//label := strings.Join(splitString[1:len(splitString)-2], "_") //join labels in the middle

	newFile := ImgFile{ 
		// labels: label
		prefix: splitString[0],
		enum: fenum,
		extension: fExtension,
		full_name: fname,
		generated_date: unix_filetime(fileinfo),
		// isLabeled: true,
	}
	// log.Debug().Str("func", "newWithPrefix").Msg("creating file from date " + newFile.labels + newFile.prefix + strconv.Itoa(newFile.enum) + newFile.extension + newFile.full_name + newFile.generated_date)
	log.Debug().Str("func", "newWithPrefix").Msg("creating file from date " + newFile.prefix + strconv.Itoa(newFile.enum) + newFile.extension + newFile.full_name + newFile.generated_date)
	return &newFile
}

func newWithPrefix(fileinfo os.FileInfo, fname string, fExtension string, prefix string) (*ImgFile) {
	temp_str := strings.TrimPrefix(fname, prefix)
	fenum, _ := strconv.Atoi(strings.TrimSuffix(temp_str, fExtension))


	newFile := ImgFile{ 
		//labels: fname[:strings.LastIndex(fname, ".")],
		prefix: prefix,
		enum: fenum,
		extension: fExtension,
		full_name: fname,
		generated_date: unix_filetime(fileinfo),
	}

	return &newFile
}

func (file ImgFile) TidyName(term string, idx int) string {
	return file.generated_date + "_" + term + "_" + strconv.Itoa(idx) + file.extension
}
func (file ImgFile) IsEmpty() bool {
	return file.extension == "" && file.prefix == "" && file.enum == 0 //&& file.labels == ""
}

////////////////////////////
////// util functions //////
////////////////////////////

func unix_filetime(path os.FileInfo) string {
	// source: https://developpaper.com/getting-access-creation-modification-time-of-files-on-linux-using-golang/ //

	// Sys () returns interface {}, so you need a type assertion. Different platforms need different types. On linux, * syscall. Stat_t
	stat_t := path.Sys().(*syscall.Stat_t)
	//fmt.Println(stat_t.Ctim.Sec)

	unix_time := time.Unix(int64(stat_t.Mtim.Sec), int64(stat_t.Mtim.Nsec) )
	// fmt.Println( unix_time )
	// fmt.Println( unix_time.Format("20060102") ) //.Format("2006-01-02")

	return unix_time.Format("20060102") //yyyymmdd

}