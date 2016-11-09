package logger

import (
	"github.com/jabong/floRest/src/common/logger"
	"github.com/jabong/floRest/src/common/monitor"
)

//StartProfile starts a profiling using profiler instance p for key. Key should have the
//following name '<package-file-method>'. For example a method 'Initialise'
//in the file 'search_dao.go' file inside package accessor should have a profile
//key as 'accessor-search_dao-Initialise'
func StartProfile(p *logger.Profiler, key string) {
	if p == nil {
		return
	}
	p.StartProfile(key)
}

//EndProfile ends the profiling starting using profiler instance p for key k
func EndProfile(p *logger.Profiler, k string, a ...interface{}) int64 {
	if p == nil {
		return 0
	}
	t := p.EndProfile(k)
	if t != logger.NilTime {
		monitor.GetInstance().Histogram(k, float64(t), nil, float64(1))
	}
	return t
}

//EndProfileCustomMetric ends the profiling starting using profiler instance p for key k
//And sends the mettic with the name n
func EndProfileCustomMetric(p *logger.Profiler, k string, n string, a ...interface{}) int64 {
	if p == nil {
		return 0
	}
	t := p.EndProfile(k)
	if t != logger.NilTime {
		monitor.GetInstance().Histogram(n, float64(t), nil, float64(1))
	}
	return t
}

//NewProfiler returns a new instance of Profiler
func NewProfiler() *logger.Profiler {
	p := logger.NewProfiler()
	if p != nil {
		return p
	}
	return nil
}
