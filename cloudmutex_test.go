package cloudmutex

import (
	"sync"
	"testing"
	"time"
)

const (
	PROJECT = "marc-general"
	BUCKET  = "cloudmutex"
	OBJECT  = "lock"
)

var (
	limit      = 1
	lockHolder = -1
)

func locker(done chan struct{}, t *testing.T, i int, m sync.Locker) {
	var lockHolderMutex sync.Mutex
	m.Lock()
	lockHolderMutex.Lock()
	if lockHolder != -1 {
		t.Errorf("%d trying to lock, but already held by %d",
			i, lockHolder)
	}
	lockHolder = i
	lockHolderMutex.Unlock()
	t.Logf("locked by %d", i)
	time.Sleep(10 * time.Millisecond)
	m.Unlock()
	lockHolderMutex.Lock()
	lockHolder = -1
	lockHolderMutex.Unlock()
	done <- struct{}{}
}

func TestParallel(t *testing.T) {
	m, err := New(nil, PROJECT, BUCKET, OBJECT)
	if err != nil {
		t.Errorf("unable to allocate a cloudmutex global object")
		return
	}
	runParallelTest(t, m)
}

func runParallelTest(t *testing.T, m sync.Locker) {
	done := make(chan struct{}, 1)
	total := 0
	for i := 0; i < limit; i++ {
		total++
		go locker(done, t, i, m)
	}
	for ; total > 0; total-- {
		<-done
	}
}

/* TODO: add testing for timed lock (both success and timeout cases)
func TestLockTimeout(t *testing.T) {
	m, err := New(nil, PROJECT, BUCKET, OBJECT)
	if err != nil {
		t.Errorf("unable to allocate a cloudmutex global object")
		return
	}
	Lock(m, 3*time.Second)
}
*/

/* TODO: add testing for timed unlock (both success and timeout cases)
func TestUnlockTimeout(t *testing.T) {
	m, err := New(nil, PROJECT, BUCKET, OBJECT)
	if err != nil {
		t.Errorf("unable to allocate a cloudmutex global object")
		return
	}
	Unlock(m, 3*time.Second)
}
*/
