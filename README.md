# waldo project


##Install : 
go install

##Getting started : 

$GOPATH/bin/waldo -h

Usage of /home/grego/go/bin/waldo:
  -key string
    	Key of the photo
  -readers int
    	Number of workers reading photos (default 12)
  -writers int
    	Number of workers indexing photos (default 12)

Use either -key if you want to query a photo or without this key value if you want to run the whole process.


./waldo

worker 3 indexing photo 01a2539a-9e53-4050-a6b5-d94f0ee6cc55.10ab5118-b1e5-4ac9-922e-125138dbdf87.jpg
worker 7 indexing photo 01a11242-35d0-4865-8f90-5db01a30ed51.e8ae7a45-8b4c-4142-b3d2-0f631d543b20.jpg
worker 5 indexing photo 01DD6669-1AF1-4A46-9D32-75FDD2882D17.ede96cc7-5500-4b3a-8828-26aabcaa2f4c.jpg
...
worker 4 indexing photo 01a8f6cd-3239-43ce-b756-95abf64a1b12.bf54b6d2-542e-4c48-8a8a-53b47e6b91d5.jpg
worker 12 indexing photo 01b819c4-c765-4dca-a407-609b64954126.a17c6591-de20-4b75-a5de-0bb11a34a116.jpg
worker 8 indexing photo 0188017b-0d90-4cab-9009-bbb74501c3d5.ede96cc7-5500-4b3a-8828-26aabcaa2f4c.jpg

FINISHED 129 photos have been processed

Query photo 

$GOPATH/bin/waldo -key 0009fcfe-376e-42fe-85a2-85ee7d2193d0.0649232a-b406-4ec1-b175-ba0d91aa3e7c.jpg
Artist = 
Copyright = 
DateTime = 2016:06:09 17:24:35
DateTimeDigitized = 2016:06:09 17:24:35
DateTimeOriginal = 2016:06:09 17:24:35
InteroperabilityIndex = R98
Make = NIKON CORPORATION
Model = NIKON D750
Software = Ver.1.02 
SubSecTime = 32
SubSecTimeDigitized = 32
SubSecTimeOriginal = 32


##Concurrency

Since I/O network resources and Disk I/O are the main bottleneck I have separated consumer from producer and made each process concurrent. The producer process reads and parses EXIF data. I called it Reader. The consumer process (Indexer) indexes EXIF data into a LevelDB data store. 

This is a parallel memory model in which multiple go routine can read simultaneously , and multiple indexers can write simultaneously to a single LevelDB datastore. 

The number of consumes rand producers are configurable.


