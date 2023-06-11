package subfix

import (
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gertm/hardsub/srt"
	"github.com/sanity-io/litter"
)

var (
	VERBOSE          = false
	DEFAULT_FONTSIZE = 22
)

func Log(msg ...interface{}) {
	if VERBOSE {
		log.Println(msg...)
	}
}

func Logf(str string, args ...interface{}) {
	if VERBOSE {
		log.Printf(str, args...)
	}
}

func FixSubs(subfile string, default_fontsize int, inplace, verbose bool) {
	VERBOSE = verbose
	sub, err := srt.ParseSrt(subfile)
	if err != nil {
		log.Fatal(err)
	}
	// first find out what the maximum size is in the subtitle.
	// then the most used size, so we'll take that as base.
	// the adjust to the size we want to see.
	fontsizes := analyzeFontSizes(sub)
	Logf("Most used font size: %d, largest font size: %d\n", fontsizes.MostUsed(), fontsizes.Largest())
	scaleFactor := float64(DEFAULT_FONTSIZE) / float64(fontsizes.MostUsed())
	if fixSubSizes(sub, scaleFactor) == nil {
		if inplace {
			Log("Sub fixing done, writing to file.")
			srt.WriteSrt(sub, subfile)
		} else {
			srt.WriteSrtToWriter(sub, os.Stdout)
		}
	}
}

type FontSizes map[int]int

func analyzeFontSizes(sub *srt.SubRip) FontSizes {
	r := regexp.MustCompile(`size="\d+"`)
	maxSize := 0
	sizes := make(map[int]int)
	for _, s := range sub.Subtitle.Content {
		for _, a := range s.Line {
			match := r.FindString(a)
			if match == "" {
				continue
			}
			size, err := strconv.Atoi(match[6:strings.LastIndex(match, "\"")])
			if err != nil {
				continue
			}
			sizes[size] += 1
			if size > maxSize {
				maxSize = size
			}
		}
	}
	Logf("max subtitle size found: %d\n", maxSize)
	if VERBOSE {
		litter.Dump("Font sizes:", sizes)
	}
	return sizes
}

func (fs FontSizes) Largest() int {
	largest := 0
	for k := range fs {
		if k > largest {
			largest = k
		}
	}
	return largest
}

func (fs FontSizes) MostUsed() int {
	mostused := 0
	highestcount := 0
	for k, v := range fs {
		if v > highestcount {
			highestcount = v
			mostused = k
		}
	}
	return mostused
}

func fixSubSizes(sub *srt.SubRip, scaleFactor float64) error {
	r := regexp.MustCompile(`size="\d+"`)
	for _, s := range sub.Subtitle.Content {
		for j, a := range s.Line {
			match := r.FindString(a)
			if match == "" {
				continue
			}
			size, err := strconv.Atoi(match[6:strings.LastIndex(match, "\"")])
			if err != nil {
				continue
			}
			newSize := float64(size) * scaleFactor
			oldstr := fmt.Sprintf("size=\"%d\"", size)
			newstr := strings.Replace(oldstr, strconv.Itoa(size), strconv.Itoa(int(math.Round(newSize))), -1)
			s.Line[j] = strings.Replace(a, oldstr, newstr, -1)
		}
	}
	return nil
}
