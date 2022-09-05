package interval

import(
	"fmt"
	"strings"
	"regexp"
	"strconv"
)

func main() {
	//inter := Process("606-609,625,aaa,6-a,4-4-4,66-50") //gives error as expected
	inter := ProcessSeq("606-609,625,aaa,6-a,4-4-4")

	fmt.Println(inter)
}



// seq is a string made of numbers, , and -
func ProcessSeq(seq string) []int {
	var full_seq []int
	
	for _, val := range strings.Split(seq, ",") {
		fmt.Println(val)
		match, err := regexp.Match("[^\\d\\-]|^-|.*-.*-", []byte(val))
		if err != nil {
			panic(err)
		}
		if !match { //meaning valid interval
			if strings.Contains(val, "-") {
				splitted := strings.Split(val, "-")
				lower, _ := strconv.Atoi(splitted[0])
				higher, _ := strconv.Atoi(splitted[1])
				if lower > higher {
					panic("the lower bound is set higher than the higher bound!")
				}
				for lower <= higher {
					full_seq = append(full_seq, lower)
					lower++
				}
			} else {
				int_val, _ := strconv.Atoi(val)
				full_seq = append(full_seq, int_val)
			}
		}
	}
	return full_seq
}