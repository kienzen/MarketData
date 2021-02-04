package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"golang.org/x/text/unicode/norm"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"sync"
)

type MarketData struct {
	ID	string
	WPSOA string
	LegNr	string
	Kursdatum	string
	Zeit	string
	Kursart	string
	Kurs	string
	Währung	string
	Notierung	string
	Kursquelle	string
	MarktbereichsIdentifikation	string
	Datenquelle string
}

type MarketDataNoTime struct {
	ID	string
	WPSOA string
	LegNr	string
	Kursdatum	string
	Kursart	string
	Kurs	string
	Währung	string
	Notierung	string
	Kursquelle	string
	MarktbereichsIdentifikation	string
	Datenquelle string
}

type IDCurrency struct{
	ID string
	Währung string
	Kursquelle string
}

var wg sync.WaitGroup

func main() {

	// Load the data

	GruppenMaster := loadData("D:\\go\\src\\marketdata\\MDL-WPTS_20210113.txt")
	GruppenMasterNoTime := loadDataNoTime("D:\\go\\src\\marketdata\\MDL-WP_20210113.txt")
	request := loadDataRequest("D:\\go\\src\\marketdata\\INST01_MDA_20210114_1244.txt")


	wg.Add(2)
	go GetMatchesNoTIme(request,GruppenMasterNoTime)
	go GetMatches(request,GruppenMaster)
	wg.Wait()

}

func GetMatchesNoTIme(request []IDCurrency, data []MarketDataNoTime)  {

	//Get matches for the small file
	matchesNoTime := make([]MarketDataNoTime,0)
	for _, record := range request{
		for _, gruppenrecord := range data{
			if record.ID==gruppenrecord.ID && record.Währung==gruppenrecord.Währung {
				matchesNoTime = append(matchesNoTime, gruppenrecord)
			}
		}
	}
	//Output the smaller file

	f, err := os.Create("D:\\go\\src\\marketdata\\OutputNoTime.csv")
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	for _, record := range matchesNoTime {
		toPrint := []string{record.String()}
		if err := w.Write(toPrint); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
	wg.Done()
}

func GetMatches(request []IDCurrency, data []MarketData)  {
	//Get matches for a large file
	matches := make([]MarketData,0)
	for _, record := range request{
		for _, gruppenrecord := range data{
			if record.ID==gruppenrecord.ID && record.Währung==gruppenrecord.Währung {
				matches = append(matches, gruppenrecord)
			}
		}
	}

	//Output the large file

	f, err := os.Create("D:\\go\\src\\marketdata\\Output.csv")
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	for _, record := range matches {
		toPrint := []string{record.String()}
		if err := w.Write(toPrint); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
	wg.Done()
}





func (d MarketData) String() string {
	v := reflect.ValueOf(d)
	var csvString string
	for i := 0; i < v.NumField(); i++ {
		csvString = fmt.Sprintf("%v%v,", csvString, v.Field(i).Interface())
	}

	return csvString
}

func (d MarketDataNoTime) String() string {
	v := reflect.ValueOf(d)
	var csvString string
	for i := 0; i < v.NumField(); i++ {
		csvString = fmt.Sprintf("%v%v,", csvString, v.Field(i).Interface())
	}

	return csvString
}

func loadDataRequest(fileLocation string) []IDCurrency {
	file, _ := os.Open(fileLocation)
	defer file.Close()

	var data IDCurrency
	request := make([]IDCurrency, 0)
	parser, _ := NewParser(file, &data)

	for {
		eof, err := parser.Next()
		if eof {
			break
		}
		if err != nil {
			panic(err)
		}
		request = append(request, data)


	}
	return request
}

func loadDataNoTime(fileLocation string) []MarketDataNoTime {
	file, _ := os.Open(fileLocation)
	defer file.Close()

	var data MarketDataNoTime
	request := make([]MarketDataNoTime, 0)
	parser, _ := NewParser(file, &data)

	for {
		eof, err := parser.Next()
		if eof {
			break
		}
		if err != nil {
			panic(err)
		}
		request = append(request, data)


	}
	return request
}

func loadData(fileLocation string) []MarketData {
	file, _ := os.Open(fileLocation)
	defer file.Close()

	var data MarketData
	request := make([]MarketData, 0)
	parser, _ := NewParser(file, &data)

	for {
		eof, err := parser.Next()
		if eof {
			break
		}
		if err != nil {
			panic(err)
		}
		request = append(request, data)


	}
	return request
}

// Parser has information for parser
type Parser struct {
	Headers    []string
	Reader     *csv.Reader
	Data       interface{}
	ref        reflect.Value
	indices    []int // indices is field index list of header array
	structMode bool
	normalize  norm.Form
}

// NewStructModeParser creates new TSV parser with given io.Reader as struct mode
func NewParser(reader io.Reader, data interface{}) (*Parser, error) {
	r := csv.NewReader(reader)
	r.Comma = '\t'

	// first line should be fields
	headers, err := r.Read()

	if err != nil {
		return nil, err
	}

	for i, header := range headers {
		headers[i] = header
	}

	p := &Parser{
		Reader:     r,
		Headers:    headers,
		Data:       data,
		ref:        reflect.ValueOf(data).Elem(),
		indices:    make([]int, len(headers)),
		structMode: false,
		normalize:  -1,
	}

	// get type information
	t := p.ref.Type()

	for i := 0; i < t.NumField(); i++ {
		// get TSV tag
		tsvtag := t.Field(i).Tag.Get("tsv")
		if tsvtag != "" {
			// find tsv position by header
			for j := 0; j < len(headers); j++ {
				if headers[j] == tsvtag {
					// indices are 1 start
					p.indices[j] = i + 1
					p.structMode = true
				}
			}
		}
	}

	if !p.structMode {
		for i := 0; i < len(headers); i++ {
			p.indices[i] = i + 1
		}
	}

	return p, nil
}

// Next puts reader forward by a line
func (p *Parser) Next() (eof bool, err error) {

	// Get next record
	var records []string

	for {
		// read until valid record
		records, err = p.Reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				return true, nil
			}
			return false, err
		}
		if len(records) > 0 {
			break
		}
	}

	if len(p.indices) == 0 {
		p.indices = make([]int, len(records))
		// mapping simple index
		for i := 0; i < len(records); i++ {
			p.indices[i] = i + 1
		}
	}

	// record should be a pointer
	for i, record := range records {
		idx := p.indices[i]
		if idx == 0 {
			// skip empty index
			continue
		}
		// get target field
		field := p.ref.Field(idx - 1)
		switch field.Kind() {
		case reflect.String:
			// Normalize text
			if p.normalize >= 0 {
				record = p.normalize.String(record)
			}
			field.SetString(record)
		case reflect.Bool:
			if record == "" {
				field.SetBool(false)
			} else {
				col, err := strconv.ParseBool(record)
				if err != nil {
					return false, err
				}
				field.SetBool(col)
			}
		case reflect.Int:
			if record == "" {
				field.SetInt(0)
			} else {
				col, err := strconv.ParseInt(record, 10, 0)
				if err != nil {
					return false, err
				}
				field.SetInt(col)
			}
		default:
			return false, errors.New("Unsupported field type")
		}
	}

	return false, nil
}