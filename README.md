## go-csvcache


## PACKAGE STATUS: IN DEVELOPMENT

## Overview
This package provides a simplistic way to cache data in
comma-separate-value (csv) files.


#### Use case
In some applications, time consuming API calls are required to
retrieve small amounts of infrequently updated, or even unchanging,
data.  This impacts overall application performance, and can lead to
application failure in the event that network connectivity is lost.

This package allows a applications to cache data in a CSV file. 


#### Cache usage pseudo code
* instantiate a `CSVCache` object using `csvcache.NewCSVCache()`
* specify the path to a cache file
* if the cache file exists, load the cache using `LoadCache()`
* during program execution, query the cache using `GetRecord()`
* on a cache miss, collect the data you wish to cache from data source
* update the cache using `AddRecord()`
* when execution is complete, check if the cache has been modified
  using `IsModified()`
* if the cache was modified, write the updated cache back to disk
  using `WriteCache()`


#### Notes on the CSV File Structure
The cache contents are persisted on disk as a CSV file.  

The unique identifier for a record (row) in the CSV file MUST always
be in the first column.

The consuming application needs to know the order and number of fields
in a given record.


#### General Warning:
* The entire contents of the cache file are loaded into memory, so
  please take into account any memory limitations in your computing
  environment before using this package.


#### Specific Example
```

//------------------------------------------------------------------------------
// INITIALIZE THE CACHE 
//------------------------------------------------------------------------------

// instantiate a CSVCache object
myCache := csvcache.NewCSVCache()

// specify a path to the cache file
cachePath := "foo.csv"

// if the cache file exists...
if _, err := os.Stat(cachePath); err == nil {

	// open the cache file for reading...
	r, err := os.Open(cachePath)
	if err != nil {
		return fmt.Errorf("error: %s", err.Error())
	}
	defer r.Close()

	// and load the cache into memory
	err = DOCache.LoadCache(io.Reader(r))
	if err != nil {
		return fmt.Errorf("error: %s", err.Error())
	}
}

...

//------------------------------------------------------------------------------
// USE THE CACHE
//------------------------------------------------------------------------------

// query the cache
record := myCache.GetRecord("some-unique-id")

// check if you have a cache hit
if len(record) != 0 {
	// extract data from the record
	bar := record[1]
	baz := record[2]
	...

} else {

	// the data you were looking for is not in the cache
	// so access your data source to gather the data
	bar := ...
	baz := ...
	
	
	// update the cache
	newRecord := make([]string, 3)
	newRecord[0] = "some-unique-id"
	newRecord[1] = bar
	newRecord[2] = baz
	
	myCache.AddRecord(newRecord)

}

...
// your application does more stuff
...


//------------------------------------------------------------------------------
// SAVE THE CACHE CONTENTS TO DISK
//------------------------------------------------------------------------------

// your application is about to exit, and you want to 
// save the cache to disk if the cache was updated

if myCache.IsModified() {

	// open the cache file for writing
	w, err := os.Create(cachePath)

	if err != nil {
		return fmt.Errorf("error: %s", err.Error())
	}
	defer w.Close()


	// write the cache to disk
	err = DOCache.WriteCache(w)
	if err != nil {
		return fmt.Errorf("error: %s", err.Error())
	}

}
```
