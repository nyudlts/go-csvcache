package csvcache

import (
	"encoding/csv"
	"fmt"
	"io"
)

const Version = "v0.3.0"

var HEADER_ROW = []string{"unique_id", "do_type", "count", "width", "height"}

type CSVCache struct {
	modified bool
	Header   []string
	cache    map[string][]string
}

func ensureCacheInit(csvc *CSVCache) {
	// initialize the map variable if it is nil
	if csvc.cache == nil {
		csvc.cache = make(map[string][]string)
	}
}

func ensureHeaderInit(csvc *CSVCache) {
	// initialize the Header if it is empty
	if len(csvc.Header) == 0 {
		// add default header
		csvc.Header = HEADER_ROW
	}
}

func NewCSVCache() *CSVCache {
	csvc := new(CSVCache)
	ensureCacheInit(csvc)
	ensureHeaderInit(csvc)
	return csvc
}

func (csvc *CSVCache) LoadCache(r io.Reader) error {

	csvr := csv.NewReader(r)

	headerRow := true
	for {
		record, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		// skip the header record
		if headerRow {
			// compare header rows to make sure that they align
			if len(record) != len(HEADER_ROW) {
				return fmt.Errorf("incompatible csv file: input file header field count does not match current header row configuration")
			}

			for idx, value := range HEADER_ROW {
				if value != record[idx] {
					return fmt.Errorf("header fields do not match: idx: %d, want: '%s', got: '%s'. expecting: %v", idx, value, record[idx], HEADER_ROW)
				}
			}

			headerRow = false
			csvc.Header = record
			continue
		}
		// load the data into the cache
		csvc.cache[record[0]] = record
	}

	csvc.modified = false
	return nil
}

func (csvc *CSVCache) GetRecord(key string) []string {
	return csvc.cache[key]
}

func (csvc *CSVCache) AddRecord(record []string) {
	csvc.modified = true
	csvc.cache[record[0]] = record
}

func (csvc *CSVCache) IsModified() bool {
	return csvc.modified
}

func (csvc *CSVCache) WriteCache(w io.Writer) error {
	csvw := csv.NewWriter(w)
	csvw.Write(csvc.Header)

	for _, record := range csvc.cache {
		err := csvw.Write(record)
		if err != nil {
			return err
		}
	}

	csvw.Flush()
	if csvw.Error() != nil {
		return csvw.Error()
	}

	return nil
}
