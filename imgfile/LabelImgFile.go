package imgfile

import (
	"os"
	"strconv"
	"strings"

	// "github.com/rs/zerolog"
	// "github.com/rs/zerolog/log"
)

////////////////////////////
////// LabelImgFile //////
////////////////////////////

type LabelImgFile struct {
	label string
	ImgFile
}
// expect string in a format eg. "20220102_landscape_1.JPEG"
func newLabelImgFile(fileinfo os.FileInfo, fname string, fExtension string) (ImgFiler) {
	splitString := strings.Split(fname, "_") // [20220102, landscape, 1.JPEG]
	tempEnum := splitString[len(splitString)-1]
	fenum, _ := strconv.Atoi(tempEnum[:strings.LastIndex(tempEnum, ".")]) // img index
	label := strings.Join(splitString[1:len(splitString)-2], "_") //join labels in the middle

	newBaseFile := ImgFile{
		prefix: splitString[0],
		enum: fenum,
		extension: fExtension,
		full_name: fname,
		generated_date: unix_filetime(fileinfo),
		// isLabeled: true,
	}

	newFile := LabelImgFile{
		label: label,
		ImgFile: newBaseFile,
	}

	// log.Debug().Str("func", "newWithPrefix").Msg("creating file from date " + newFile.labels + newFile.prefix + strconv.Itoa(newFile.enum) + newFile.extension + newFile.full_name + newFile.generated_date)
	// log.Debug().Str("func", "newWithDate").Msg("creating file from date " + newFile.prefix + strconv.Itoa(newFile.enum) + newFile.extension + newFile.full_name + newFile.generated_date)
	return newFile
}