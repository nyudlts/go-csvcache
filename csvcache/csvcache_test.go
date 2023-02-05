package csvcache

import (
	"os"
	"testing"
)

func assertStringsEqual(want, got string, t *testing.T) {
	if want != got {
		t.Errorf("want: %s , got: %s", want, got)
	}
}

func assertRecordsEqual(want, got []string, t *testing.T) {
	if len(want) != len(got) {
		t.Errorf("record lengths do not match: len(want): %d , len(got): %d", len(want), len(got))
	}

	for idx, value := range want {
		if value != got[idx] {
			t.Errorf("record fields do not match: idx: %d, want: %s, got: %s", idx, value, got[idx])
		}
	}
}

func TestLoadCache(t *testing.T) {
	sut := new(CSVCache)

	fixturePath := "./testdata/basic.csv"
	r, err := os.Open(fixturePath)
	if err != nil {
		t.Errorf("problem opening %s", fixturePath)
	}

	err = sut.LoadCache(r)
	if err != nil {
		t.Errorf("problem loading the cache")
	}
}
func TestHeaderRow(t *testing.T) {
	var want, got string

	cache := new(CSVCache)

	fixturePath := "./testdata/basic.csv"
	r, err := os.Open(fixturePath)
	if err != nil {
		t.Errorf("problem opening %s", fixturePath)
	}

	err = cache.LoadCache(r)
	if err != nil {
		t.Errorf("problem loading the cache")
	}

	sut := cache.HeaderRow

	want = "unique_id"
	got = sut[0]
	assertStringsEqual(want, got, t)

	want = "do_type"
	got = sut[1]
	assertStringsEqual(want, got, t)

	want = "count"
	got = sut[2]
	assertStringsEqual(want, got, t)
}

func TestGetRecord(t *testing.T) {
	var want, got []string

	sut := new(CSVCache)

	fixturePath := "./testdata/basic.csv"
	r, err := os.Open(fixturePath)
	if err != nil {
		t.Errorf("problem opening %s", fixturePath)
	}

	err = sut.LoadCache(r)
	if err != nil {
		t.Errorf("problem loading the cache")
	}

	want = []string{"ghx3fpf7", "image_set", "2"}
	got = sut.GetRecord("ghx3fpf7")
	assertRecordsEqual(want, got, t)
}
