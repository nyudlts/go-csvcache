package csvcache

import (
	"encoding/csv"
	"io"
)

type CSVCache struct {
	modified  bool
	HeaderRow []string
	cache     map[string][]string
}

func (csvc *CSVCache) LoadCache(r io.Reader) error {

	csvc.cache = make(map[string][]string)
	csvr := csv.NewReader(r)

	headerRow := true
	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		// skip the header row
		if headerRow {
			headerRow = false
			csvc.HeaderRow = row
			continue
		}
		// load the data into the cache
		csvc.cache[row[0]] = row
	}

	csvc.modified = false
	return nil
}

func (csvc *CSVCache) GetRecord(key string) []string {
	return csvc.cache[key]
}

func (csvc *CSVCache) AddRecord(row []string) {
	csvc.modified = true

	csvc.cache[row[0]] = row
}

func (csvc *CSVCache) IsModified() bool {
	return csvc.modified
}
