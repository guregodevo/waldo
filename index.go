package main

import (
	"github.com/syndtr/goleveldb/leveldb"
        "fmt"
)

type Indexer struct {
	db *leveldb.DB
}

// NewIndexer will return an intialized Indexer struct
func NewIndexer(db *leveldb.DB) *Indexer {
	return &Indexer{db: db}
}

//Index photo EXIF data
func (i *Indexer) Index(key string, tags map[string]string) {
	batch := new(leveldb.Batch)
	for k := range tags {
		var b []byte
		s := key
		b = append(b, s...)
		b = append(b, k...)
		batch.Put(b, []byte(tags[k]))
	}
	i.db.Write(batch, nil)
}

func (i *Indexer) worker(id int, jobs chan *PhotoResult, out chan bool) {
        for j := range jobs {
                if j == nil {
                    out <- false
                } else {
                  fmt.Println("worker", id, "indexing photo", j.key)
                  i.Index(j.key,j.tags)
                  out <- true 
                }
        }
}



