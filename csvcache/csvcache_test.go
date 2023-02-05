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

func assertBoolsEqual(want, got bool, t *testing.T) {
	if want != got {
		t.Errorf("want: %v , got: %v", want, got)
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
	defer r.Close()

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
	defer r.Close()

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

func TestGetRecordPresent(t *testing.T) {
	var want, got []string

	sut := new(CSVCache)

	fixturePath := "./testdata/basic.csv"
	r, err := os.Open(fixturePath)
	if err != nil {
		t.Errorf("problem opening %s", fixturePath)
	}
	defer r.Close()

	err = sut.LoadCache(r)
	if err != nil {
		t.Errorf("problem loading the cache")
	}

	want = []string{"ghx3fpf7", "image_set", "2"}
	got = sut.GetRecord("ghx3fpf7")
	assertRecordsEqual(want, got, t)
}

func TestGetRecordMissing(t *testing.T) {
	sut := new(CSVCache)

	fixturePath := "./testdata/basic.csv"
	r, err := os.Open(fixturePath)
	if err != nil {
		t.Errorf("problem opening %s", fixturePath)
	}
	defer r.Close()

	err = sut.LoadCache(r)
	if err != nil {
		t.Errorf("problem loading the cache")
	}

	got := sut.GetRecord("this-key-does-not-have-a-record")
	assertRecordsEqual(nil, got, t)
}

func TestAddRecord(t *testing.T) {
	sut := new(CSVCache)

	fixturePath := "./testdata/basic.csv"
	r, err := os.Open(fixturePath)
	if err != nil {
		t.Errorf("problem opening %s", fixturePath)
	}
	defer r.Close()

	err = sut.LoadCache(r)
	if err != nil {
		t.Errorf("problem loading the cache")
	}

	got := sut.GetRecord("9ec2c7f5d0c4")
	assertRecordsEqual(nil, got, t)

	record := []string{"9ec2c7f5d0c4", "whoop", "97"}
	sut.AddRecord(record)

	want := []string{"9ec2c7f5d0c4", "whoop", "97"}
	got = sut.GetRecord("9ec2c7f5d0c4")
	assertRecordsEqual(want, got, t)

}

func TestIsModified(t *testing.T) {
	sut := new(CSVCache)

	fixturePath := "./testdata/basic.csv"
	r, err := os.Open(fixturePath)
	if err != nil {
		t.Errorf("problem opening %s", fixturePath)
	}
	defer r.Close()

	err = sut.LoadCache(r)
	if err != nil {
		t.Errorf("problem loading the cache")
	}

	got := sut.GetRecord("9ec2c7f5d0c4")
	assertRecordsEqual(nil, got, t)
	assertBoolsEqual(false, sut.IsModified(), t)

	record := []string{"9ec2c7f5d0c4", "whoop", "97"}
	sut.AddRecord(record)
	assertBoolsEqual(true, sut.IsModified(), t)

	want := []string{"9ec2c7f5d0c4", "whoop", "97"}
	got = sut.GetRecord("9ec2c7f5d0c4")
	assertRecordsEqual(want, got, t)
}

func TestWriteCache(t *testing.T) {
	sut1 := new(CSVCache)

	fixturePath := "./testdata/basic.csv"
	tmpFilePath := "./testdata/tmp.csv"

	r, err := os.Open(fixturePath)
	if err != nil {
		t.Errorf("problem opening %s", fixturePath)
	}
	defer r.Close()

	err = sut1.LoadCache(r)
	if err != nil {
		t.Errorf("problem loading the cache")
	}

	got := sut1.GetRecord("9ec2c7f5d0c4")
	assertRecordsEqual(nil, got, t)
	assertBoolsEqual(false, sut1.IsModified(), t)

	record := []string{"9ec2c7f5d0c4", "whoop", "97"}
	sut1.AddRecord(record)
	assertBoolsEqual(true, sut1.IsModified(), t)

	want := []string{"9ec2c7f5d0c4", "whoop", "97"}
	got = sut1.GetRecord("9ec2c7f5d0c4")
	assertRecordsEqual(want, got, t)

	w, err := os.Create(tmpFilePath)
	if err != nil {
		t.Errorf("problem creating %s", tmpFilePath)
	}
	defer w.Close()

	err = sut1.WriteCache(w)
	if err != nil {
		t.Errorf("problem writing %s", tmpFilePath)
	}

}
