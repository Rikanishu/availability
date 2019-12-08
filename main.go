package main

import (
	"log"
	"time"
)

func main() {

	st := NewStorage(NewCSVDataSource())

	startTSNano := time.Now().UnixNano()
	s, _ := time.Parse("2006-01-02 15:04:05", "2019-09-12 10:00:00")
	e, _ := time.Parse("2006-01-02 15:04:05", "2019-10-12 09:59:59")
	log.Printf("found %d results", len(st.FindAvail(s.Unix(), e.Unix())))
	log.Printf("search took %v sec", float64(time.Now().UnixNano()-startTSNano)/float64(time.Second))

	startTSNano = time.Now().UnixNano()
	s, _ = time.Parse("2006-01-02 15:04:05", "2019-10-09 12:00:00")
	e, _ = time.Parse("2006-01-02 15:04:05", "2019-10-08 23:59:01")
	log.Printf("found %d results", len(st.FindAvail(s.Unix(), e.Unix())))
	log.Printf("search took %v sec", float64(time.Now().UnixNano()-startTSNano)/float64(time.Second))
}
