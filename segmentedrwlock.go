package eventbus

import (
	"sync"
)

type SegmentedRWLock struct {
	segments int
	locks    []sync.RWMutex
}

func NewSegmentedRWLock(segments int) *SegmentedRWLock {
	return &SegmentedRWLock{
		segments: segments,
		locks:    make([]sync.RWMutex, segments),
	}
}

func (self *SegmentedRWLock) locker(key string) *sync.RWMutex {
	hash := abs(hash([]byte(key)))
	return &self.locks[hash%self.segments]
}

func abs(x int) int {
	if x < 0 {
		return x * -1
	}
	return x
}

func hash(bytes []byte) int {
	var h int
	for i := 0; i < len(bytes); i++ {
		h = 31*h ^ int(bytes[i])
	}
	return h
}
