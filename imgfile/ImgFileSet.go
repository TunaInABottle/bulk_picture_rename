package imgfile

import (
	"os"
	"fmt"
	"time"
	"strconv"
	
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
// as editor files have beeen created in a different day
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





////////////////////////////
////// util functions //////
////////////////////////////

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}