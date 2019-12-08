package main

import (
	"fmt"
	"log"
	"time"
)

func main() {

	st := NewStorage(NewCSVDataSource())

	startTSNano := time.Now().UnixNano()
	s := time.Now().Add(-90 * 24 * time.Hour).Unix()
	e := time.Now().Add(-70 * 24 * time.Hour).Unix()
	fmt.Println(st.FindAvail(s, e))
	log.Printf("search took %v sec", float64(time.Now().UnixNano()-startTSNano)/float64(time.Second))
}
