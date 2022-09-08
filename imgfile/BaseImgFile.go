package imgfile


import (
	"os"
	"strconv"
	"strings"

	// "github.com/rs/zerolog"
	// "github.com/rs/zerolog/log"
)

////////////////////////////
////// BaseImgFile //////
////////////////////////////

type BaseImgFile struct {
	ImgFile
}

func newBaseImgFile(fileinfo os.FileInfo, fname string, fExtension string, prefix string) (ImgFiler) {
	temp_str := strings.TrimPrefix(fname, prefix)
	fenum, _ := strconv.Atoi(strings.TrimSuffix(temp_str, fExtension))

	newFile := BaseImgFile{ImgFile: ImgFile{ 
		//labels: fname[:strings.LastIndex(fname, ".")],
		prefix: prefix,
		enum: fenum,
		extension: fExtension,
		full_name: fname,
		generated_date: unix_filetime(fileinfo),
	}}
	return newFile
}