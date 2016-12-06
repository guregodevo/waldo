package main

import (
	"flag"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"strings"
	"sync/atomic"
)

/**
 * Export the following environment variables :
 *
 * export AWS_ACCESS_KEY_ID='AKID'
 * export AWS_SECRET_ACCESS_KEY='SECRET'
 * export AWS_REGION='us-east-1'
 */

var (
	readerworkerNum  = flag.Int("readers", 12, "Number of workers reading photos")
	indexerworkerNum = flag.Int("writers", 12, "Number of workers indexing photos")
	photoKey         = flag.String("key", "", "Key of the photo")
)

func main() {
	db, err := leveldb.OpenFile("photo.db", nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	flag.Parse()
	if *photoKey != "" {
		iter := db.NewIterator(util.BytesPrefix([]byte(*photoKey)), nil)
		for iter.Next() {
			key := iter.Key()
			value := iter.Value()
			fmt.Println(strings.SplitAfter(string(key), ":")[1], "=", string(value))
		}
		iter.Release()
		err = iter.Error()

		return
	}
	//concurrent read, concurrent write
	readjobs := make(chan string, 100)
	indexjobs := make(chan *PhotoResult, 100)
	notifications := make(chan bool, 100)

	reader, _ := NewReader()
	photoObjects, err := reader.List()
	fmt.Println(len(photoObjects), "photos")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for w := 1; w <= *readerworkerNum; w++ {
		r, _ := NewReader()
		go r.worker(w, readjobs, indexjobs)
	}

	for w := 1; w <= *indexerworkerNum; w++ {
		i := NewIndexer(db)
		go i.worker(w, indexjobs, notifications)
	}

	for i := range photoObjects {
		readjobs <- *photoObjects[i].Key
	}
	close(readjobs)

	cnt := int32(0)
	for a := 0; a < len(photoObjects); a++ {
		<-notifications
		cnt = atomic.AddInt32(&cnt, 1)
	}

	fmt.Println("FINISHED", cnt, "photos have been processed")
}
