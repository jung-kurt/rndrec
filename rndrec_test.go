package rndrec

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestSrcType_general(t *testing.T) {
	var err error
	var r *SrcType

	// empty list
	_, err = NewRandomRecordSource([][]string{}, 0, 42)
	if err == nil {
		t.Fatal("empty record list unexpectedly accepted")
	}

	// zero weight
	_, err = NewRandomRecordSource([][]string{
		{"a", "0"},
		{"b", "0"},
		{"c", "0"},
	}, 1, 0)
	if err == nil {
		t.Fatal("record list with zero cumulative weight unexpectedly accepted")
	}

	// invalid weight column
	_, err = NewRandomRecordSource([][]string{
		{"a", "1"},
		{"b", "2"},
		{"c", "3"},
	}, 2, 0)
	if err == nil {
		t.Fatal("record list with invalid weight column indicator unexpectedly accepted")
	}

	// exercise Stringer method for test coverage
	r, _ = NewRandomRecordSource([][]string{
		{"a", "1"},
		{"b", "2"},
		{"c", "3"},
	}, 1, 0)
	if len(r.String()) == 0 {
		t.Fatal("invalid Stringer implementation")
	}

}

func keyList(m map[string]int) (list []string) {
	for k := range m {
		list = append(list, k)
	}
	gensort(len(list), func(a, b int) bool {
		return list[a] < list[b]
	}, func(a, b int) {
		list[a], list[b] = list[b], list[a]
	})
	return
}

// For test purposes, but not a requirement of the package, the first field of
// each data set is assumed to be a unique key
func srcReport(r *SrcType, weightCol int) {
	const loopCount = 100000
	var mp map[string]int
	var key string
	var fields []string
	mp = make(map[string]int)
	for j := 0; j < loopCount; j++ {
		fields = r.Record()
		key = fields[0]
		mp[key] = mp[key] + 1
	}
	fields = keyList(mp)
	for _, key = range fields {
		fmt.Printf("%s: %.2f\n", key, float64(mp[key])/float64(loopCount))
	}
}

func report(list [][]string, weightCol int) {
	var r *SrcType
	var err error
	r, err = NewRandomRecordSource(list, weightCol, 42)
	if err == nil {
		srcReport(r, weightCol)
	}
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleSrcType_simple() {
	var list = [][]string{
		{"20%", "20"},
		{"30%", "30"},
		{"10%", "10"},
		{"40%", "40"},
	}
	report(list, 1)
	// Output:
	// 10%: 0.10
	// 20%: 0.20
	// 30%: 0.30
	// 40%: 0.40
}

func ExampleSrcType_equalWeight() {
	var list = [][]string{
		{"red"},
		{"green"},
		{"blue"},
	}
	report(list, -1)
	// Output:
	// blue: 0.33
	// green: 0.33
	// red: 0.33
}

func ExampleSrcType_readme() {
	var r *SrcType
	var err error
	var rec []string

	r, err = NewRandomRecordSourceFromFile("continent_population.csv", 1, '|', 0)
	if err == nil {
		for row := 0; row < 8; row++ {
			for col := 0; col < 8; col++ {
				if col > 0 {
					fmt.Printf(" | ")
				}
				rec = r.Record()
				fmt.Printf("%s", rec[0])
			}
			fmt.Println("")
		}
	} else {
		fmt.Printf("%s\n", err)
	}
	// Output:
	// South America | Asia | Asia | Africa | Asia | Asia | Asia | Asia
	// North America | Asia | Asia | North America | Europe | Asia | Asia | Asia
	// Europe | Africa | Europe | Europe | Asia | Asia | Asia | Asia
	// Asia | Asia | Asia | Asia | Asia | Asia | Africa | Asia
	// Asia | Asia | Asia | Asia | Asia | Asia | Asia | Africa
	// Asia | Africa | Asia | Asia | Europe | Africa | North America | North America
	// Asia | Europe | Africa | Europe | Asia | South America | Africa | Europe
	// Asia | Europe | Africa | Asia | Asia | Asia | Asia | Africa
}

func ExampleSrcType_file() {
	var r *SrcType
	var err error

	r, err = NewRandomRecordSourceFromFile("continent_population.csv", 1, '|', 0)
	if err == nil {
		srcReport(r, 1)
	} else {
		fmt.Printf("%s\n", err)
	}
	// Output:
	// Africa: 0.15
	// Asia: 0.61
	// Australia: 0.01
	// Europe: 0.11
	// North America: 0.07
	// South America: 0.06
}

func ExampleSrcType_population() {
	var list = [][]string{
		{"Africa", "1,030,400,000"},
		{"Antarctica", "0"},
		{"Asia", "4,157,300,000"},
		{"Australia", "36,700,000"},
		{"Europe", "738,600,000"},
		{"North America", "461,114,000"},
		{"South America", "390,700,000"},
	}
	report(list, 1)
	// Output:
	// Africa: 0.15
	// Asia: 0.61
	// Australia: 0.01
	// Europe: 0.11
	// North America: 0.07
	// South America: 0.06
}

// Generate dummy names based on 1990 US census data
func ExampleSrcType_names() {
	const (
		cnLast = iota
		cnFemale
		cnMale
		cnCount
	)
	var filenameList = [cnCount]string{
		"data/us/name_last.csv",
		"data/us/name_first_female.csv",
		"data/us/name_first_male.csv",
	}
	var srcList [cnCount]*SrcType
	var err error
	var rnd *rand.Rand
	var first, mid, last []string
	var j, k int

	for j = 0; j < cnCount && err == nil; j++ {
		srcList[j], err = NewRandomRecordSourceFromFile(filenameList[j], 1, '|', 0)
	}
	if err == nil {
		rnd = rand.New(rand.NewSource(0))
		for j = 0; j < 16; j++ {
			if rnd.Intn(5) < 4 {
				k = cnFemale
			} else {
				k = cnMale
			}
			first = srcList[k].Record()
			mid = srcList[k].Record()
			last = srcList[cnLast].Record()
			fmt.Printf("%s %s %s\n", first[0], mid[0][0:1], last[0])
		}
	}
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	// Output:
	// Kendall T Creel
	// Earl J Cox
	// Jasmin A Stein
	// Yolanda L Brown
	// Evelyn M Perkins
	// Sharon Y Foster
	// Lea S Carter
	// Martha A Potts
	// Jeannie V Ayres
	// Veronica B Wright
	// Harriet M Simmons
	// Janie L Colburn
	// Anthony P Pulliam
	// Teresa D Coleman
	// Florence C Sweeney
	// Sarah B Ramirez
}
