package model

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type Filter struct {
	cond string
	sync.RWMutex
}

func NewFilter(cond string) *Filter {
	f := new(Filter)
	f.cond = cond
	return f
}

func (f *Filter) And(cond string) *Filter {
	f.Lock()
	defer f.Unlock()
	if f.cond == "" {
		f.cond = cond
	} else {
		f.cond = fmt.Sprintf("%s and %s", f.cond, cond)
	}
	return f
}

func (f *Filter) Or(cond string) *Filter {
	f.Lock()
	defer f.Unlock()
	if f.cond == "" {
		f.cond = cond
	} else {
		f.cond = fmt.Sprintf("%s or %s", f.cond, cond)
	}
	return f
}

func (f *Filter) AndNot(cond string) *Filter {
	f.Lock()
	defer f.Unlock()
	if f.cond == "" {
		f.cond = cond
	} else {
		f.cond = fmt.Sprintf("%s and not %s", f.cond, cond)
	}
	return f
}

func (f *Filter) OrNot(cond string) *Filter {
	f.Lock()
	defer f.Unlock()
	if f.cond == "" {
		f.cond = cond
	} else {
		f.cond = fmt.Sprintf("%s or not %s", f.cond, cond)
	}
	return f
}

func In(key string, list interface{}) string {
	if reflect.TypeOf(list).Kind() != reflect.Slice &&
		reflect.TypeOf(list).Kind() != reflect.Array {
		return ""
	}
	values := reflect.ValueOf(list)
	if values.Len() == 0 {
		return ""
	}
	var b strings.Builder
	for i := 0; i < values.Len(); i++ {
		b.WriteString(",")
		v := values.Index(i)
		b.WriteString(fmt.Sprintf("%v", v.Interface()))
	}
	if b.Len() != 0 {
		return fmt.Sprintf("%s in (%s)", key, b.String()[1:])
	}
	return ""
}

func (f *Filter) Cond() string {
	f.RLock()
	defer f.RUnlock()
	return f.cond
}
