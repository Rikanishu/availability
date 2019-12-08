package main

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

type CSVDataSource struct {
}

func NewCSVDataSource() *CSVDataSource {
	return &CSVDataSource{}
}

func (s *CSVDataSource) LoadAvailability(c DataSourceConf) ([]Availability, error) {
	gStartRange := c.StartRangeTS
	gEndRange := c.EndRangeTS
	file, err := os.Open("./calendar.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	out := make([]Availability, 0)
	objectID := ""
	takenDates := make([]Range, 0)
	k := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if k != 0 && objectID != record[0] {
			availDates, err := ConvertIntoAvailDates(takenDates, gStartRange, gEndRange)
			if err != nil {
				return nil, err
			}
			out = append(out, Availability{
				ObjectID:   objectID,
				AvailDates: availDates,
			})
			objectID = ""
			takenDates = takenDates[:0]
		}
		objectID = record[0]
		startTS, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return nil, err
		}
		endTS, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return nil, err
		}
		takenDates = append(takenDates, Range{StartTS: startTS, EndTS: endTS})
		k++
	}
	availDates, err := ConvertIntoAvailDates(takenDates, gStartRange, gEndRange)
	if err != nil {
		return nil, err
	}
	out = append(out, Availability{
		ObjectID:   objectID,
		AvailDates: availDates,
	})

	return out, nil
}
