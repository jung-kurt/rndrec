# rndrec

[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/jung-kurt/rndrec/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/jung-kurt/rndrec?status.svg)](https://godoc.org/github.com/jung-kurt/rndrec)
[![Language](https://img.shields.io/badge/language-go-blue.svg)](https://golang.org/)
[![Report Card](https://goreportcard.com/badge/github.com/jung-kurt/rndrec)](https://goreportcard.com/report/github.com/jung-kurt/rndrec)

Package rndrec is used to randomly select records from a pool based on their
relative weight. For example, if the relative weight of one record is 50, on
average it will be selected five times more often than a record with a relative
weight of 10. This is useful for generating plausible data sets for testing
purposes, for example names based on frequency or regions based on population.

## Example
Given a file named "continent_population.csv" with the following contents,

```
Africa|1,030,400,000
Antarctica|0
Asia|4,157,300,000
Australia|36,700,000
Europe|738,600,000
North America|461,114,000
South America|390,700,000
```

the following call will create a weighted record sample source:

```
var r *SrcType
var err error

r, err = NewRandomRecordSourceFromFile("continent_population.csv", 1, '|', 0)
```

The integer argument following the filename is the zero-based column that
contains the relative weights in numeric form. Note that the commas in these
values are disregarded. The rune argument following the weight column specifies
the field separator. All input records are assumed to be delimited with
newlines. The final argument is the seed value for the instance's random number
source. This can be used to generate repeatable sequences. time.Now().Unix()
can be used if repeatable sequences are not desired.

Call `r.Record()` to randomly retrieve weighted records:

```go
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
```

This will generate the following ouput:

```
South America | Asia | Asia | Africa | Asia | Asia | Asia | Asia
North America | Asia | Asia | North America | Europe | Asia | Asia | Asia
Europe | Africa | Europe | Europe | Asia | Asia | Asia | Asia
Asia | Asia | Asia | Asia | Asia | Asia | Africa | Asia
Asia | Asia | Asia | Asia | Asia | Asia | Asia | Africa
Asia | Africa | Asia | Asia | Europe | Africa | North America | North America
Asia | Europe | Africa | Europe | Asia | South America | Africa | Europe
Asia | Europe | Africa | Asia | Asia | Asia | Asia | Africa
```

## Installation
To install the package on your system, run

```
go get github.com/jung-kurt/rndrec
```

## License
rndrec is released under the MIT License.
