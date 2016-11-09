package logger

import (
	"time"
)

//Niltime used when time is nil or no time captured
const (
	NilTime        = -1
	timeUnit int64 = int64(time.Millisecond)
)

type Profiler struct {

	//Key specifies a key for a profiling event
	key string

	//StartTime specifies the timestamp in nano second when a profile event
	//was started
	startTime time.Time
}

//StartProfile starts a profiling using profiler instance p for key. Key should have the
//following name '<package-file-method>'. For example a method 'Initialise'
//in the file 'search_dao.go' file inside package accessor should have a profile
//key as 'accessor-search_dao-Initialise'
func (p *Profiler) StartProfile(key string) {
	if conf.ProfilerEnabled == false {
		return
	}
	p.startTime = time.Now()
	p.key = key
}

//EndProfile ends the profiling using profiler instance p for key k. Return time in MicroSeconds
func (p *Profiler) EndProfile(k string) int64 {
	if conf.ProfilerEnabled == false {
		return NilTime
	}
	duration := time.Now().Sub(p.startTime).Nanoseconds()
	t := duration / timeUnit
	p = nil
	return t
}

//NewProfiler returns a new instance of Profiler
func NewProfiler() *Profiler {
	if conf.ProfilerEnabled == false {
		return nil
	}
	return new(Profiler)
}
