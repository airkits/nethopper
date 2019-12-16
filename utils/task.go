package utils

import (
	"errors"
	"fmt"
	"time"
)

// TaskManager accurate to seconds only, timed tasks less than 1s will fail
type TaskManager struct {
	tm *TimerManager
	// Mapping between task id and timerID, keep id globally unique
	taskTimerID map[uint64]TimerID
}

// NewTaskManager create task manager
func NewTaskManager() *TaskManager {
	schedule := &TaskManager{tm: NewTimerManager(time.Second)}
	return schedule
}

// Serve start time loop
func (t *TaskManager) Serve() {
	go func() {
		for {
			t.tm.DetectTimerInLock()
			time.Sleep(20 * time.Millisecond)
		}
	}()
}

//Tick use in your own loop
func (t *TaskManager) Tick() {
	t.tm.DetectTimerInLock()
}

//RunAt callback function should keep goruntine safe
func (t *TaskManager) RunAt(unix int64, f func(interface{})) (TimerID, error) {
	now := time.Now().Unix()
	interval := unix - now
	if interval < 0 {
		err := errors.New("can run at past time")
		return 0, err
	}
	tempTimer, ID := NewTimer(TimerTypeOnce)
	tempTimer.Start(uint64(interval), f, t.tm)
	return ID, nil
}

//RunAfter add timer run once,callback function should keep goruntine safe
func (t *TaskManager) RunAfter(d time.Duration, f func(interface{})) (TimerID, error) {
	interval := d / time.Second
	if interval <= 0 {
		err := fmt.Errorf("invalid interval time! %d", interval)
		return 0, err
	}
	tempTimer, ID := NewTimer(TimerTypeOnce)
	tempTimer.Start(uint64(interval), f, t.tm)
	return ID, nil
}

//RunLoop run loop timer,callback function should keep goruntine safe
func (t *TaskManager) RunLoop(d time.Duration, f func(interface{})) (TimerID, error) {
	interval := d / time.Second
	if interval <= 0 {
		err := fmt.Errorf("invalid interval time! %d", interval)
		return 0, err
	}
	tempTimer, ID := NewTimer(TimerTypeLoop)
	tempTimer.Start(uint64(interval), f, t.tm)
	return ID, nil
}

//Update update timer duration,callback function should keep goruntine safe
func (t *TaskManager) Update(ID TimerID, d time.Duration, f func(interface{})) error {
	interval := d / time.Second
	if interval <= 0 {
		err := fmt.Errorf("invalid interval time! %d", interval)
		return err
	}
	tempTimer := t.tm.GetTimerByID(ID)
	if tempTimer == nil {
		err := fmt.Errorf("find timer failed by  %v", ID)
		return err
	}
	tempTimer.Update(uint64(interval), f, t.tm)
	return nil
}

// Stop timer by id
func (t *TaskManager) Stop(ID TimerID) error {
	tempTimer := t.tm.GetTimerByID(ID)
	if tempTimer == nil {
		err := fmt.Errorf("find timer failed by  %v", ID)
		return err
	}
	tempTimer.Stop(t.tm)
	return nil
}
