package imgfile

import "fmt"

type ymdDate string

func (d ymdDate) String() string {
	return fmt.Sprintf("%s", string(d))
}

// type ymdDate struct {
// 	ymd string
// }