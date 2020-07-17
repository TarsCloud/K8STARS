package recentuse

import (
	"sync"
	"time"
)

type node struct {
	time       time.Time
	value      string
	prev, next *node
}

// RecentUse implements algorithm of recent use
type RecentUse struct {
	recentTime time.Duration
	m          map[string]*node
	head, tail *node
	lock       sync.Mutex
}

// NewRecentUse returns instance of RecentUse
func NewRecentUse(recentTime time.Duration) *RecentUse {
	return &RecentUse{
		m:          make(map[string]*node),
		recentTime: recentTime,
		lock:       sync.Mutex{},
	}
}

// KeepAlive returns true if key is used recently
func (ru *RecentUse) KeepAlive(key string) bool {
	ru.lock.Lock()
	defer ru.lock.Unlock()
	// remove old data
	for n := ru.tail; n != nil; n = n.prev {
		if time.Since(n.time) < ru.recentTime {
			break
		}
		delete(ru.m, n.value)
		ru.tail = ru.tail.prev
		if ru.tail != nil {
			ru.tail.next = nil
		}
		if ru.head == n {
			ru.head = nil
		}
	}
	// check active
	bActive := false
	activeTime := time.Now()
	if v, ok := ru.m[key]; ok {
		bActive = true
		activeTime = v.time
		// remove node
		if v.prev != nil {
			v.prev.next = v.next
		}
		if v.next != nil && v.prev != nil {
			v.prev.next = v.next
		}
		if v == ru.head {
			ru.head = ru.head.next
		}
		if v == ru.tail {
			ru.tail = ru.tail.prev
		}
	}
	// add to head
	n := &node{
		value: key,
		time:  activeTime,
	}
	n.next = ru.head
	if ru.head != nil {
		ru.head.prev = n
	}
	if ru.tail == nil {
		ru.tail = n
	}
	ru.head = n
	ru.m[key] = n
	return bActive
}

// String returns string of node list
func (n *node) String() string {
	s := ""
	for nn := n; n != nil; n = n.next {
		if s != "" {
			s += "->"
		}
		s += nn.value
	}
	return s
}
