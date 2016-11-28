// Package rndrec is used to randomly select records from a pool based on their
// relative weight. For example, if the relative weight of one record is 50, on
// average it will be selected five times more often than a record with a
// relative weight of 10. This is useful for generating plausible data sets for
// testing purposes, for example names based on frequency or regions based on
// population.
package rndrec

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strconv"
)

var reIntDelimiter = regexp.MustCompile("[,_]")

type recType struct {
	// cumulative frequency
	cf float64
	// list of record fields
	fields []string
}

// SrcType is used to generate plausible random records based on a list of
// weighted records.
type SrcType struct {
	// list of records with ascending cumulative frequency
	list []recType
	// maximum cumulative frequency
	cfMax float64
	// random number generator
	rand *rand.Rand
}

// String implements the fmt.Stringer interface
func (r SrcType) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "Cumulative frequency maximum: %.2f\n", r.cfMax)
	for j, rec := range r.list {
		fmt.Fprintf(&b, "%2d: [%10.2f] %v\n", j, rec.cf, rec.fields)
	}
	return b.String()
}

// NewRandomRecordSource processes a list of multi-field records in which each
// field is a string. With one exception, one column must be an integer weight.
// In this column, specified by weightColPos, each occurrence of an underscore,
// comma or period is removed and the remaining string is parsed as an integer.
// The values in this column are relative weights; that is, a record that has a
// weight twice that of some other record will be selected by Record() on
// average twice as often. The sum of these weights does not have to be any
// special value. The exception to the requirement that one field be a weight
// is when all records are weighted equally. In this case, weightColPos can be
// set to -1 and records do not need to have a weight column. Records returned
// by the Record() method depend on a local pseudo-random number generator;
// seed is used to seed this generator. If any value in the column specified by
// weightColPos can not be parsed as an integer, or the cumulative value of
// weights is zero, or the number of records is zero, an error is returned.
// Otherwise, err is nil and src may be used to retrieve records that are
// distributed according to their relative weights.
func NewRandomRecordSource(recs [][]string, weightColPos int, seed int64) (src *SrcType, err error) {
	var rec recType
	var weight string
	src = new(SrcType)
	for _, fields := range recs {
		if err == nil {
			if weightColPos == -1 {
				rec.cf = 1
			} else if weightColPos >= 0 && weightColPos < len(fields) {
				weight = fields[weightColPos]
				rec.cf, err = strconv.ParseFloat(reIntDelimiter.ReplaceAllString(weight, ""), 64)
				// rec.cf, err = strconv.ParseInt(reIntDelimiter.ReplaceAllString(weight, ""), 10, 64)
			} else {
				err = fmt.Errorf("specified weight column (%d) is out of range", weightColPos)
			}
			if err == nil {
				rec.fields = fields
				src.list = append(src.list, rec)
			}
		}
	}
	if err == nil {
		if len(src.list) > 0 {
			for j, rec := range src.list {
				src.cfMax += rec.cf
				src.list[j].cf = src.cfMax
			}
			if src.cfMax > 0 {
				src.rand = rand.New(rand.NewSource(seed))
			} else {
				err = fmt.Errorf("cumulative frequency must be greater than zero")
			}
		} else {
			err = fmt.Errorf("number of records must be greater than zero")
		}
	}
	if err != nil {
		src = nil
	}
	return
}

// NewRandomRecordSourceFromReader processes a list of multi-field records in
// the form of a comma-separated-value buffer that can be read with the
// io.Reader r. Each record must be separated by a newline. Each field is
// separated by the value specified by fieldSep. For more information on the
// return value and the other arguments, see NewRandomRecordSource().
func NewRandomRecordSourceFromReader(r io.Reader, weightColPos int, fieldSep rune, seed int64) (src *SrcType, err error) {
	var rdr *csv.Reader
	var recs [][]string
	rdr = csv.NewReader(r)
	rdr.Comma = fieldSep
	recs, err = rdr.ReadAll()
	if err == nil {
		src, err = NewRandomRecordSource(recs, weightColPos, seed)
	}
	return
}

// NewRandomRecordSourceFromFile processes a list of multi-field records in the
// form of a comma-separated-value file with the filename specified by fileStr.
// Each record must be separated by a newline. Each field is separated by the
// value specified by fieldSep. For more information on the return value and
// the other arguments, see NewRandomRecordSource().
func NewRandomRecordSourceFromFile(fileStr string, weightColPos int, fieldSep rune, seed int64) (src *SrcType, err error) {
	var f *os.File
	f, err = os.Open(fileStr)
	if err == nil {
		src, err = NewRandomRecordSourceFromReader(f, weightColPos, fieldSep, seed)
		f.Close()
	}
	return
}

// Record returns a random record based on its relative weight. For example, a
// record with a relative weight of 40 will be returned, on average, four times
// as often as a record with the relative weight of 10. The returned record
// will be in the form of a slice of strings taken directly from the original
// list used to initialize the SrcType instance.
func (r *SrcType) Record() []string {
	var cf float64
	var pos int
	cf = r.rand.Float64() * r.cfMax // 0 <= cf < r.cfMax
	pos = sort.Search(len(r.list), func(j int) bool {
		return r.list[j].cf > cf
	})
	return r.list[pos].fields
}
