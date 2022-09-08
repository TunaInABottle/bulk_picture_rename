package imgfile

import (
	// "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	)

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
