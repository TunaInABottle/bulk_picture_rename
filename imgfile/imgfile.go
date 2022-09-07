package imgfile

import (

	"os"
	"fmt"
	"syscall"
	"time"
	"strconv"
	"strings"
	"regexp"
	// "errors"
	
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	)


var editor_extensions []string = []string{".ORF.dop", ".JPG.dop"}
var allowedExtensions []string = append([]string{".ORF", ".JPG"}, editor_extensions...)
const allow_rename = false //disabled during debugging

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
}

// TODO interface to decouple ImgFile and a future struct

///////////////////////////////
/////////  ImgFileSet /////////
///////////////////////////////

type ImgFileSet struct {
	date_idx_set 	map[string]map[string][]ImgFile
	idx_date 	 	map[string]string					// imgidx1: date
	dLI				dateLabelIdx						// Date > Label > Idx
}

//@TODO put in Add the following
// ifs.dLI.SetIdx(ymd, label, 42) 
// handle difference between TUNA and pre-named

// suppose all maps are up-to-date
func (ifs ImgFileSet) Rename1(directory string, interval []int, label string) {
	for _, fileIdx := range interval { // iterate over [605, 606, 607, 625] (TUNA files)
		strFileIdx := strconv.Itoa(fileIdx)
		fileDate := ifs.idx_date[strFileIdx]
		ymd := ymdDate{fileDate}

		if (ymd != ymdDate{}) { // date exists in map
			ifs.dLI.increaseIdx(ymd, label)
			labelIdx := ifs.dLI.getIdx(ymd, label)
			log.Debug().Int("idx_label", labelIdx).Send()
			// files := ifs.date_idx_set[fileDate][strFileIdx] 
			// for _, file := range files {
			// 	newFileName := file.TidyName(label, labelIdx)

			// 	if _, err := os.Stat(directory + "/" + newFileName);
			// 	err == nil {
			// 		log.Fatal().Str("oldName", file.full_name).Str("newName", newFileName).Msg("file already existing")
			// 	} else if os.IsNotExist(err) {
			// 		log.Info().Str("dir", directory).Str("oldName", file.full_name).Str("newName", newFileName).Msg("renaming file")
			// 		if allow_rename {
			// 			os.Rename(
			// 				directory + "/" + file.full_name, 
			// 				directory + "/" + newFileName)
			// 		}
			// 	} else {
			// 		// something else happened
			// 		panic(err)
			// 	}
			// }
		}
	}
}

// suppose all maps are up-to-date
func (ifs ImgFileSet) Rename(directory string, interval []int, label string) {
	dateIdxMap := make(map[string]int)

	for _, fileIdx := range interval {
		strFileIds := strconv.Itoa(fileIdx)
		fileDate := ifs.idx_date[strFileIds]

		//ymd := ymdDate{fileDate}
		if fileDate != "" { // file exists
			dateIdxMap[fileDate] += 1
			fileNameIdx := dateIdxMap[fileDate]
			files := ifs.date_idx_set[fileDate][strFileIds] 
			for _, file := range files {
				newFileName := file.TidyName(label, fileNameIdx)

				if _, err := os.Stat(directory + "/" + newFileName);
				err == nil {
					log.Fatal().Str("oldName", file.full_name).Str("newName", newFileName).Msg("file already existing")
				} else if os.IsNotExist(err) {
					log.Info().Str("dir", directory).Str("oldName", file.full_name).Str("newName", newFileName).Msg("renaming file")
					if allow_rename {
						os.Rename(
							directory + "/" + file.full_name, 
							directory + "/" + newFileName)
					}
				} else {
					// something else happened
					panic(err)
				}
			}
		}
	}
}

func (fileset *ImgFileSet) Add(file ImgFile) {

	// @TODO check file validity in its entirety first!

	if file.IsEmpty() {
		log.Debug().Msg("Empty file, continuing")
		//return fileset
	}

	fileset.createMissingMaps(file)
	
	// add to idx_date map BEGIN
	editorFile := false
	for _, val := range editor_extensions {
		if file.extension == val {
			editorFile = true
		}
	}
	if !editorFile && fileset.idx_date[strconv.Itoa(file.enum)] == "" {
		fileset.idx_date[strconv.Itoa(file.enum)] = file.generated_date
	}
	// add to idx_date map END

	// add to date_idx_set map BEGIN
	var dateKey string
	if editorFile {
		dateKey = "0" //ASSUMPTION: no file can have index 0
	} else {
		dateKey = file.generated_date
	}
	slice_to_append := fileset.date_idx_set[dateKey][strconv.Itoa(file.enum)]
	fileset.date_idx_set[dateKey][strconv.Itoa(file.enum)] = append(slice_to_append, file)
	// add to date_idx_set map END
}
func (fileset ImgFileSet) Print() {
	fmt.Println("ImgFileSet tree:")
	for common_date, scd_map := range fileset.date_idx_set {
		fmt.Println(common_date)
		for common_idx, file_slice := range scd_map{
			fmt.Println(" ↳", common_idx)
			// fmt.Println("   ", scd_map)
			for _, val := range file_slice {
				fmt.Println(" | ↳", val.Print())
			}
		}
	}

	fmt.Println("\nthe other map\n", fileset.idx_date)
}
func (fileset *ImgFileSet) createMissingMaps(file ImgFile) {
	// main date_idx_map
	if fileset.date_idx_set == nil {
		log.Debug().Str("func", "createMissingMaps").Msg("fileset.date_idx_set empty, creating...")
		fileset.date_idx_set = make(map[string]map[string][]ImgFile)
	}
	// nested date map in date_idx_set
	if fileset.date_idx_set[file.generated_date] == nil {
		log.Debug().Str("func", "createMissingMaps").Str("date", file.generated_date).Msg("fileset fileset.date_idx_set of date empty, creating...")
		fileset.date_idx_set[file.generated_date] = make(map[string][]ImgFile)
	}
	// nested date 0 map in date_idx_set
	if fileset.date_idx_set["0"] == nil {
		log.Debug().Str("date", "0").Msg("fileset fileset.date_idx_set of date empty, creating...")
		fileset.date_idx_set["0"] = make(map[string][]ImgFile)
	}
	// map in idx_date
	if fileset.idx_date == nil {
		log.Debug().Str("func", "createMissingMaps").Msg("idx_date Smap empty, creating...")
		fileset.idx_date = make(map[string]string)
	}
}
// put editor files in the same map of the original file
func (fileset *ImgFileSet) RecoverEditorFiles() {
	recoveryMaps := fileset.date_idx_set["0"]

	if recoveryMaps != nil {
		for _, editorFileSlice := range recoveryMaps {
			for _, editorFile := range editorFileSlice {
				fileIndex := strconv.Itoa(editorFile.enum)
				originalFileDate := fileset.idx_date[fileIndex]
				if originalFileDate != "" {
					log.Debug().Str("file", editorFile.full_name).Str("orgDate", originalFileDate).Msg("Adjusted editor file date successfully")
					editorFile.generated_date = originalFileDate
					fileset.date_idx_set[originalFileDate][fileIndex] = append(fileset.date_idx_set[originalFileDate][fileIndex], editorFile)
				} else {
					log.Debug().Str("file", editorFile.full_name).Str("orgDate", originalFileDate).Msg("Adjusting editor file date failed")
				}
			}
		}
		delete(fileset.date_idx_set, "0")
	}
}



////////////////////////////////
///////// dateLabelIdx /////////
////////////////////////////////
//helps to find the existing indexes of labels

type dateLabelIdx struct {
	mapDLI map[ymdDate]map[string]int 
}

func (dli *dateLabelIdx) createMissingMaps(dateKey ymdDate) {
	if dli.mapDLI == nil {
		dli.mapDLI = make(map[ymdDate]map[string]int)
	}
	if dli.mapDLI[dateKey] == nil {
		dli.mapDLI[dateKey] = make(map[string]int)
	}
}
func (dli *dateLabelIdx) SetIdx(dateKey ymdDate, label string, val int) {
	dli.createMissingMaps(dateKey)
	log.Debug().Str("date", dateKey.ymd).Str("label", label).Int("val", val).Str("func", "dateLabelIdx").Msg("setting value")
	dli.mapDLI[dateKey][label] = max(dli.mapDLI[dateKey][label], val)
}
func (dli *dateLabelIdx) increaseIdx(dateKey ymdDate, label string) {
	dli.createMissingMaps(dateKey)
	// @TODO how to distinguish TUNA from already named?
	dli.mapDLI[dateKey][label] += 1
}
func (dli dateLabelIdx) getIdx(dateKey ymdDate, label string) int {
	dli.createMissingMaps(dateKey)
	// @TODO what if 0 (not initialised)?
	log.Debug().Str("date", dateKey.ymd).Str("label", label).Str("func", "dateLabelIdx").Msg("getting value")
	return dli.mapDLI[dateKey][label]
}




////////////////////////////
///////// ImgFiler /////////
////////////////////////////

type ImgFiler interface {
	Print() string //@TODO rename String()
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
	labels string			// the labels in the file name eg. landscape_lake
	prefix string			// the first component of a filename
	enum int				// progressive number written in the filename
	extension string		// extension of the file
	full_name string		// full file name
	generated_date string	// date in which the picture has been taken
}

func (file ImgFile) Print() string {
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
	label := strings.Join(splitString[1:len(splitString)-2], "_") //join labels in the middle

	newFile := ImgFile{ 
		labels: label, //here it is the label
		prefix: splitString[0],
		enum: fenum, //
		extension: fExtension,
		full_name: fname,
		generated_date: unix_filetime(fileinfo),
	}
	log.Debug().Str("func", "newWithPrefix").Msg("creating file from date " + newFile.labels + newFile.prefix + strconv.Itoa(newFile.enum) + newFile.extension + newFile.full_name + newFile.generated_date)
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
	return file.labels == "" && file.extension == "" && file.prefix == ""
}

///////////////////////////
///////// ymdDate /////////
///////////////////////////

///////////////////////////
///////// ymdDate /////////
///////////////////////////

type ymdDate struct {
	ymd string
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

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}