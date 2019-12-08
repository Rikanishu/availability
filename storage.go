package main

import (
	"log"
	"sync"
	"time"

	"github.com/google/btree"
)

type DataSourceConf struct {
	StartRangeTS int64
	EndRangeTS   int64
}

type DataSource interface {
	LoadAvailability(c DataSourceConf) ([]Availability, error)
}

type NodeRangeStart struct {
	StartTS int64
	Tree    *btree.BTree
}

func (n *NodeRangeStart) Less(than btree.Item) bool {
	t := than.(*NodeRangeStart)
	return n.StartTS < t.StartTS
}

type NodeRangeEnd struct {
	EndTS     int64
	ObjectIDs []string
}

func (n *NodeRangeEnd) Less(than btree.Item) bool {
	t := than.(*NodeRangeEnd)
	return n.EndTS < t.EndTS
}

type Storage struct {
	tree       *btree.BTree
	lock       sync.RWMutex
	dataSource DataSource
}

func NewStorage(ds DataSource) *Storage {
	s := &Storage{
		dataSource: ds,
	}

	s.lock.Lock()
	go func() {
		defer s.lock.Unlock()
		err := s.refreshTree()
		if err != nil {
			log.Fatal(err)
		}

	}()

	return s
}

func (s *Storage) Refresh() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.refreshTree()
}

func (s *Storage) FindAvail(startTS int64, endTS int64) (out []string) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.tree == nil {
		return nil
	}

	sNodes := make([]*NodeRangeStart, 0)
	s.tree.AscendGreaterOrEqual(&NodeRangeStart{
		StartTS: startTS,
	}, func(item btree.Item) bool {
		sNodes = append(sNodes, item.(*NodeRangeStart))
		return true
	})

	out = make([]string, 0)
	for _, sn := range sNodes {
		sn.Tree.AscendLessThan(&NodeRangeEnd{
			EndTS: endTS,
		}, func(item btree.Item) bool {
			i := item.(*NodeRangeEnd)
			out = append(out, i.ObjectIDs...)
			return true
		})
	}
	return
}

func (s *Storage) refreshTree() error {
	log.Print("refreshing the tree...")
	// start building for -120 daus from now (due the dataset)
	gStartRange := time.Now().Add(-24 * 120 * time.Hour).Unix()
	// building for +30 days from now
	gEndRange := time.Now().Add(24 * 30 * time.Hour).Unix()

	log.Print("extracting the data...")
	startTSNano := time.Now().UnixNano()
	avail, err := s.dataSource.LoadAvailability(DataSourceConf{
		StartRangeTS: gStartRange,
		EndRangeTS:   gEndRange,
	})
	if err != nil {
		return err
	}
	log.Printf("done, took %v sec", float64(time.Now().UnixNano()-startTSNano)/float64(time.Second))
	log.Print("building the tree...")

	startTSNano = time.Now().UnixNano()

	m := make(map[int64]map[int64][]string)
	for _, objAv := range avail {
		for _, dRange := range objAv.AvailDates {
			if _, ok := m[dRange.StartTS]; !ok {
				m[dRange.StartTS] = make(map[int64][]string)
			}
			if _, ok := m[dRange.StartTS][dRange.EndTS]; !ok {
				m[dRange.StartTS][dRange.EndTS] = make([]string, 0, 1)
			}
			m[dRange.StartTS][dRange.EndTS] = append(m[dRange.StartTS][dRange.EndTS], objAv.ObjectID)
		}
	}
	t := btree.New(2)
	for startTS, ends := range m {
		et := btree.New(2)
		for endTS, ids := range ends {
			et.ReplaceOrInsert(&NodeRangeEnd{
				EndTS:     endTS,
				ObjectIDs: ids,
			})
		}
		rs := &NodeRangeStart{
			StartTS: startTS,
			Tree:    et,
		}
		t.ReplaceOrInsert(rs)
	}
	s.tree = t
	log.Printf("done, took %v sec", float64(time.Now().UnixNano()-startTSNano)/float64(time.Second))

	return nil
}
