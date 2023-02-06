package csvcache

import (
	"encoding/csv"
	"io"
)

type CSVCache struct {
	modified bool
	Header   []string
	cache    map[string][]string
}

func NewCSVCache() *CSVCache {
	csvc := new(CSVCache)
	assertCacheInit(csvc)
	return csvc
}

func assertCacheInit(csvc *CSVCache) {
	// initialize the map variable if it is nil
	if csvc.cache == nil {
		csvc.cache = make(map[string][]string)
	}

	// // initialize the Header if it is empty
	// if len(csvc.Header) == 0 {
	// 	// add default header
	// 	csvc.Header = []string{"unique_id", "do_type", "count"}
	// }
}

func (csvc *CSVCache) LoadCache(r io.Reader) error {

	assertCacheInit(csvc)

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
	assertCacheInit(csvc)
	return csvc.cache[key]
}

func (csvc *CSVCache) AddRecord(record []string) {
	assertCacheInit(csvc)

	csvc.modified = true

	csvc.cache[record[0]] = record
}

func (csvc *CSVCache) IsModified() bool {
	return csvc.modified
}

func (csvc *CSVCache) WriteCache(w io.Writer) error {
	assertCacheInit(csvc)

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
