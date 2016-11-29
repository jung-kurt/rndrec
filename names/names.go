// This command reads various United States census files ands generates files
// that are compatible with the rndrec package. The census bureau publishes the
// data with the following statement:
//
//   Copyright protection is not available for any work of the United States
//   Government (Title 17 U.S.C., Section 105). Thus you are free to reproduce
//   census materials as you see fit. We would ask, however, that you cite the
//   Census Bureau as the source.
//
// The following files are used as input to this command:
//
//   http://www2.census.gov/topics/genealogy/1990surnames/dist.all.last
//   http://www2.census.gov/topics/genealogy/1990surnames/dist.female.first
//   http://www2.census.gov/topics/genealogy/1990surnames/dist.male.first
//
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// JAMES          3.318  3.318      1
var reData = regexp.MustCompile("^(\\S+)\\s+(\\d+\\.\\d+)\\s+(\\d+\\.\\d+)\\s+(\\d+)$")

var reApostrophe = regexp.MustCompile("'")

func process(inFileStr, outFileStr string, correct map[string]string) (err error) {
	var inFile, outFile *os.File
	var scanner *bufio.Scanner
	var strList []string
	var nameStr, subStr, wtStr, str string
	var wr *bufio.Writer
	var wt float64
	var wtErr error
	var ok bool

	inFile, err = os.Open(inFileStr)
	if err == nil {
		outFile, err = os.Create(outFileStr)
		if err == nil {
			wr = bufio.NewWriter(outFile)
			scanner = bufio.NewScanner(inFile)
			for scanner.Scan() {
				str = scanner.Text()
				strList = reData.FindStringSubmatch(str)
				if strList != nil {
					nameStr = strList[1]
					wtStr = strList[2]
					wt, wtErr = strconv.ParseFloat(wtStr, 64)
					if wtErr == nil && wt > 0.002 {
						subStr, ok = correct[nameStr]
						if ok {
							nameStr = subStr
						} else {
							if len(nameStr) > 1 {
								nameStr = nameStr[:1] + strings.ToLower(nameStr[1:])
							}
						}
						fmt.Fprintf(wr, "%s|%.3f\n", nameStr, wt)
					}
				} else {
					fmt.Printf("Match error: [%s]\n", str)
				}
			}
			err = scanner.Err()
			wr.Flush()
			outFile.Close()
		}
		inFile.Close()
	}
	return
}

func corrections(fileStr string) (mp map[string]string, err error) {
	var f *os.File
	var scanner *bufio.Scanner
	var str, keyStr string

	f, err = os.Open(fileStr)
	if err == nil {
		mp = make(map[string]string)
		scanner = bufio.NewScanner(f)
		for scanner.Scan() {
			str = scanner.Text()
			keyStr = strings.ToUpper(reApostrophe.ReplaceAllString(str, ""))
			mp[keyStr] = str
		}
		err = scanner.Err()
		f.Close()
	}
	return

}

func main() {
	type ioType struct {
		in, out string
	}
	var list = []ioType{
		{in: "dist.all.last", out: "../data/us/name_last.csv"},
		{in: "dist.female.first", out: "../data/us/name_female_first.csv"},
		{in: "dist.male.first", out: "../data/us/name_male_first.csv"},
	}
	var rec ioType
	var err error
	var correct map[string]string

	correct, err = corrections("corrections.txt")
	if err == nil {
		for _, rec = range list {
			if err == nil {
				err = process(rec.in, rec.out, correct)
			}
		}
	}
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
