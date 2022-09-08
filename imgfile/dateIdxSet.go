package imgfile

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
