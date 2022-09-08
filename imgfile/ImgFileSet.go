package imgfile

import (

	"os"
	"fmt"
	// "syscall"
	"time"
	"strconv"
	// "strings"
	// "regexp"
	// "errors"
	
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	)


var editor_extensions []string = []string{".ORF.dop", ".JPG.dop"}
var base_extensions []string = []string{".ORF", ".JPG"}
var allowedExtensions []string = append(base_extensions, editor_extensions...)
const allow_rename = false //disabled during debugging
const impDate = ymdDate("20000231") // 31st february 2000, used for mapping editor files

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
}

// TODO interface to decouple ImgFile and a future struct

///////////////////////////////
/////////  ImgFileSet /////////
///////////////////////////////

type ImgFileSet struct {
	dIS				dateIdxSet 							// Date > ImgIdx > setImg
	dLI				dateLabelIdx						// Date > Label > Idx
	idx_date 	 	map[string]string					// imgidx1 > date
}

//@TODO put in Add the following
// handle difference between TUNA and pre-named
// eg. ifs.dLI.SetIdx(ymd, label, 42) 

// suppose all maps are up-to-date
func (ifs ImgFileSet) Rename(directory string, interval []int, label string) {
	for _, fileIdx := range interval { // iterate over [605, 606, 607, 625] (TUNA files)
		strFileIdx := strconv.Itoa(fileIdx)
		fileDate := ymdDate( ifs.idx_date[strFileIdx] )

		if (fileDate != "") { // date exists in map 
							  //@TODO is check good? (infering underlying type)
			ifs.dLI.increaseIdx(fileDate, label)
			labelIdx := ifs.dLI.getIdx(fileDate, label)
			files := ifs.dIS.GetImgFiles(fileDate, strFileIdx)

			for _, file := range files {
				newFileName := file.TidyName(label, labelIdx)

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
	
	// editor file check
	editorFile := false
	for _, val := range editor_extensions {
		if file.extension == val {
			editorFile = true
		}
	}

	// add to idx_date map BEGIN
	if fileset.idx_date == nil { 
		fileset.idx_date = make(map[string]string)
	}

	if !editorFile && fileset.idx_date[strconv.Itoa(file.enum)] == "" {
		fileset.idx_date[strconv.Itoa(file.enum)] = file.generated_date
	}
	// add to idx_date map END

	// add to dIS map BEGIN
	var dateKey ymdDate
	if editorFile {
		dateKey = impDate
	} else {
		dateKey = ymdDate(file.generated_date)
		log.Debug().Str("date", file.generated_date).Str("date_obj", dateKey.String()).Send()
	}
	fileset.dIS.Add(dateKey, strconv.Itoa(file.enum), file )
	// add to date_idx_set map END
}

func (fileset ImgFileSet) String() {
	fmt.Println("ImgFileSet tree in Date > Index > ImgSet class:")
	fmt.Println(fileset.dIS)
	fmt.Println("\nMap Index > Date\n", fileset.idx_date)
}

// put editor files in the same map of the original file
// change with change on editor file
func (fileset *ImgFileSet) RecoverEditorFiles() {
	recoveryMaps := fileset.dIS.GetIndexMaps(impDate)

	if recoveryMaps != nil {
		for _, editorFileSlice := range recoveryMaps {
			for _, editorFile := range editorFileSlice {
				fileIndex := strconv.Itoa(editorFile.enum)
				originalFileDate := fileset.idx_date[fileIndex]
				if originalFileDate != "" {
					editorFile.generated_date = originalFileDate
					fileset.dIS.Add(ymdDate(originalFileDate), fileIndex, editorFile)
					log.Debug().Str("file", editorFile.full_name).Str("orgDate", originalFileDate).Msg("Adjusted editor file date successfully")
				} else {
					log.Debug().Str("file", editorFile.full_name).Str("orgDate", originalFileDate).Msg("Adjusting editor file date failed")
				}
			}
		}
		fileset.dIS.deleteDate(impDate)
	}
}

//////////////////////////////
///////// dateIdxSet /////////
//////////////////////////////

// stores for each date and index, a slice of images
type dateIdxSet struct {
	path_DIS map[ymdDate]map[string][]ImgFile
}

func (dis *dateIdxSet) createDISMaps(dateKey ymdDate) {
	if dis.path_DIS == nil {
		dis.path_DIS = make(map[ymdDate]map[string][]ImgFile)
	}
	if dis.path_DIS[dateKey] == nil {
		dis.path_DIS[dateKey] = make(map[string][]ImgFile)
	}
}
func (dis *dateIdxSet) Add(dateKey ymdDate, idx string, file ImgFile) {
	dis.createDISMaps(dateKey)
	dis.path_DIS[dateKey][idx] = append(dis.path_DIS[dateKey][idx], file)
}
func (dis dateIdxSet) GetIndexMaps(date ymdDate) map[string][]ImgFile {
	return dis.path_DIS[date]
}
func (dis dateIdxSet) GetImgFiles(date ymdDate, idx string) []ImgFile {
	// @TODO check empty maps?
	return dis.path_DIS[date][idx]
}
func (dis dateIdxSet) String() string {
	var retString string
	for common_date, scd_map := range dis.path_DIS {
		retString += common_date.String() + "\n"
		for common_idx, file_slice := range scd_map{
			retString += " ↳" + common_idx + "\n"
			for _, val := range file_slice {
				retString += " | ↳" + val.String() + "\n"
			}
		}
	}
	return retString
}
func (dis *dateIdxSet) deleteDate(dateToDel ymdDate) {
	delete(dis.path_DIS, dateToDel)
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
	log.Debug().Str("date", dateKey.String()).Str("label", label).Int("val", val).Str("func", "dateLabelIdx").Msg("setting value")
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
	// log.Debug().Str("date", dateKey.String()).Str("label", label).Str("func", "dateLabelIdx").Msg("getting value")
	return dli.mapDLI[dateKey][label]
}

////////////////////////////
////// util functions //////
////////////////////////////

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}