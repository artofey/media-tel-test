package csv

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	rr "github.com/artofey/media-tel-test"
)

type CSVFile struct {
	fileName string
	csv      *csv.Reader
	records  [][]string
	h        []string
}

func New(fileName string) (*CSVFile, error) {
	newCSV := CSVFile{}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("open file err: %w", err)
	}
	defer file.Close()

	str, err := bufio.NewReader(file).ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("read file err: %w", err)
	}
	// set offset to the beginning of the file
	file.Seek(0, 0)

	csv := csv.NewReader(file)
	csv.Comma = newCSV.commaFromLine(str)
	csv.LazyQuotes = true
	csv.TrimLeadingSpace = true
	records, err := csv.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read csv file err: %w", err)
	}
	if len(records) < 1 {
		return nil, fmt.Errorf("file is empty")
	}
	headers, generated := newCSV.headers(records)
	if !generated {
		records = append(records[:0], records[1:]...)
	}

	// file name without dir and extention
	newCSV.fileName = strings.Split(path.Base(fileName), ".")[0]
	newCSV.csv = csv
	newCSV.records = records
	newCSV.h = headers
	return &newCSV, nil
}

func (c *CSVFile) GetRecords() (*rr.Records, error) {
	return &rr.Records{
		Headers: c.h,
		R:       c.records,
	}, nil
}

func (c *CSVFile) FileName() string {
	return c.fileName
}

func (c *CSVFile) String() string {
	var resStr []string
	resStr = append(resStr, fmt.Sprintf("%v\n", c.h))
	for _, str := range c.records {
		resStr = append(resStr, fmt.Sprintf("%v\n", str))
	}
	return fmt.Sprintf("%v", resStr)
}

// commaFromLine is detect comma rune in CSV-file
func (*CSVFile) commaFromLine(s string) rune {
	runes := make(map[rune]int)
	runes[','] = strings.Count(s, ",")
	runes[';'] = strings.Count(s, ";")
	runes[':'] = strings.Count(s, ":")
	fmt.Println("1 line: ", s)
	var max int
	var res rune
	var i int
	for k, v := range runes {
		if i == 0 {
			max = v
			res = k
			i++
			continue
		} else if max < v {
			max = v
			res = k
		}
	}
	return res
}

/*
Файл не имеет заголовка если:
- TODO: Первая строка содержит даты или другие распространенные форматы данных (например, xx-xx-xx)
*/
func (c *CSVFile) headers(rec [][]string) ([]string, bool) {
	set := make(map[string]struct{})
	for _, col := range rec[0] {
		set[col] = struct{}{}
		// если первая строка имеет столбцы, которые не являются строками
		var i int
		if _, err := fmt.Sscan(col, &i); err == nil {
			return c.makeHeaders(len(rec[0]))
		}
		// если первая строка имеет столбцы, которые пусты
		if col == "" {
			return c.makeHeaders(len(rec[0]))
		}
	}
	// если столбцы первой строки не все уникальны
	if len(set) != len(rec[0]) {
		return c.makeHeaders(len(rec[0]))
	}

	return rec[0], false
}

func (*CSVFile) makeHeaders(len int) ([]string, bool) {
	var res []string
	for i := 1; i <= len; i++ {
		res = append(res, "col"+strconv.Itoa(i))
	}
	return res, true
}
