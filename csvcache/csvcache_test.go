package csvcache

import (
	"os"
	"path/filepath"
	"testing"
)

const fixtureRoot = "./testdata"
const tmpFileRoot = "./testdata"

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

func createAndLoadCSVCache(path string, t *testing.T) *CSVCache {
	c := NewCSVCache()

	r, err := os.Open(path)
	if err != nil {
		t.Errorf("problem opening %s", path)
	}
	defer r.Close()

	err = c.LoadCache(r)
	if err != nil {
		t.Errorf("problem loading the cache")
	}

	return c
}

//------------------------------------------------------------------------------
// begin tests
//------------------------------------------------------------------------------

func TestHeader(t *testing.T) {
	var want, got string

	cache := createAndLoadCSVCache(filepath.Join(fixtureRoot, "basic.csv"), t)
	sut := cache.Header

	want = "unique_id"
	got = sut[0]
	assertStringsEqual(want, got, t)

	want = "do_type"
	got = sut[1]
	assertStringsEqual(want, got, t)

	want = "count"
	got = sut[2]
	assertStringsEqual(want, got, t)

	want = "width"
	got = sut[3]
	assertStringsEqual(want, got, t)

	want = "height"
	got = sut[4]
	assertStringsEqual(want, got, t)
}

func TestGetRecordPresent(t *testing.T) {
	var want, got []string

	sut := createAndLoadCSVCache(filepath.Join(fixtureRoot, "basic.csv"), t)

	want = []string{"ghx3fpf7", "image_set", "2", "136", "108"}
	got = sut.GetRecord("ghx3fpf7")
	assertRecordsEqual(want, got, t)
}

func TestGetRecordMissing(t *testing.T) {

	sut := createAndLoadCSVCache(filepath.Join(fixtureRoot, "basic.csv"), t)

	got := sut.GetRecord("this-key-does-not-have-a-record")
	assertRecordsEqual(nil, got, t)
}

func TestAddRecord(t *testing.T) {

	sut := createAndLoadCSVCache(filepath.Join(fixtureRoot, "basic.csv"), t)

	got := sut.GetRecord("9ec2c7f5d0c4")
	assertRecordsEqual(nil, got, t)

	record := []string{"9ec2c7f5d0c4", "whoop", "97", "100", "1000"}
	sut.AddRecord(record)

	want := []string{"9ec2c7f5d0c4", "whoop", "97", "100", "1000"}
	got = sut.GetRecord("9ec2c7f5d0c4")
	assertRecordsEqual(want, got, t)

}

func TestIsModified(t *testing.T) {

	sut := createAndLoadCSVCache(filepath.Join(fixtureRoot, "basic.csv"), t)

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

func TestAddRecordToCSVCThatWasNotLoadedFromFile(t *testing.T) {

	sut := NewCSVCache()

	record := []string{"9ec2c7f5d0c4", "whoop", "97"}
	sut.AddRecord(record)
	assertBoolsEqual(true, sut.IsModified(), t)

	want := []string{"9ec2c7f5d0c4", "whoop", "97"}
	got := sut.GetRecord("9ec2c7f5d0c4")
	assertRecordsEqual(want, got, t)
}

func TestWriteCache(t *testing.T) {

	sut1 := createAndLoadCSVCache(filepath.Join(fixtureRoot, "basic.csv"), t)

	// assert baseline
	got := sut1.GetRecord("9ec2c7f5d0c4")
	assertRecordsEqual(nil, got, t)
	assertBoolsEqual(false, sut1.IsModified(), t)

	// add record
	record := []string{"9ec2c7f5d0c4", "whoop", "97", "100", "1000"}
	sut1.AddRecord(record)
	assertBoolsEqual(true, sut1.IsModified(), t)

	// assert that record is now in cache
	want := []string{"9ec2c7f5d0c4", "whoop", "97", "100", "1000"}
	got = sut1.GetRecord("9ec2c7f5d0c4")
	assertRecordsEqual(want, got, t)

	// open the target file
	w, err := os.Create(filepath.Join(tmpFileRoot, "tmp-basic.csv"))
	if err != nil {
		t.Errorf("problem creating %s", filepath.Join(tmpFileRoot, "tmp-basic.csv"))
	}
	defer w.Close()

	// write the target file
	err = sut1.WriteCache(w)
	if err != nil {
		t.Errorf("problem writing %s", filepath.Join(tmpFileRoot, "tmp-basic.csv"))
	}

	sut2 := createAndLoadCSVCache(filepath.Join(tmpFileRoot, "tmp-basic.csv"), t)

	assertRecordsEqual(sut1.Header, sut2.Header, t)
	assertRecordsEqual(sut1.GetRecord("m63xss7g"), sut2.GetRecord("m63xss7g"), t)
	assertRecordsEqual(sut1.GetRecord("ghx3fpf7"), sut2.GetRecord("ghx3fpf7"), t)
	assertRecordsEqual(sut1.GetRecord("zkh18f2c"), sut2.GetRecord("zkh18f2c"), t)
	assertRecordsEqual(sut1.GetRecord("xgxd28gq"), sut2.GetRecord("xgxd28gq"), t)
	assertRecordsEqual(sut1.GetRecord("9ec2c7f5d0c4"), sut2.GetRecord("9ec2c7f5d0c4"), t)

	// cleanup
	err = os.Remove(filepath.Join(tmpFileRoot, "tmp-basic.csv"))
	if err != nil {
		t.Errorf("problem removing %s", filepath.Join(tmpFileRoot, "tmp-basic.csv"))
	}
}

func TestHeaderInitialized(t *testing.T) {

	sut := NewCSVCache()

	want := []string{"unique_id", "do_type", "count", "width", "height"}
	got := sut.Header
	assertRecordsEqual(want, got, t)
}

/* func TestIncompatibleFile(t *testing.T) {

	sut := NewCSVCache()

	want := []string{"unique_id", "do_type", "count", "width", "height"}
	got := sut.Header
	assertRecordsEqual(want, got, t)
}
*/
