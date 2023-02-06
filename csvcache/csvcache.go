package csvcache

import (
	"encoding/csv"
	"io"
)

type CSVCache struct {
	modified  bool
	Header []string
	cache     map[string][]string
}

func (csvc *CSVCache) LoadCache(r io.Reader) error {

	csvc.cache = make(map[string][]string)
	csvr := csv.NewReader(r)

	headerRecord := true
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
		if headerRecord {
			headerRecord = false
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
