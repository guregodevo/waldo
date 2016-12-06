package main

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"sync/atomic"
)

/**
 * Don't hard-code your credentials!
 * Export the following environment variables instead:
 *
 * export AWS_ACCESS_KEY_ID='AKID'
 * export AWS_SECRET_ACCESS_KEY='SECRET'
 * export AWS_REGION='us-east-1'
 */

//Number of workers
const readerworkerNum = 12 
const indexerworkerNum =  4 

func main() {
	db, err := leveldb.OpenFile("photo.db", nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

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

	for w := 1; w <= readerworkerNum; w++ {
	        r, _ := NewReader()
		go r.worker(w, readjobs, indexjobs)
	}
	
        for w := 1; w <= indexerworkerNum; w++ {
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
        close(notifications)
        close(indexjobs)

	fmt.Println("FINISHED", cnt, "photos have been processed")
}
