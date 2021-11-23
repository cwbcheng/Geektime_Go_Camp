package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	Up = iota
	Circle_Open
)


type Hystrix struct {
	windowPoint time.Time
	windowSpan time.Duration
	rollingWindowPoint time.Time
	rollingWindowSpan time.Duration
	errorThresholdPercentage int32
	stats int8
	normalCount int32
	errorCount    int32
	minTotalCount int32
	retried       bool
}

func InitHystrix() *Hystrix {
	return &Hystrix {
		windowPoint:              time.Time{},
		windowSpan:               10 * time.Second,
		rollingWindowPoint:       time.Time{},
		rollingWindowSpan:        5 * time.Second,
		errorThresholdPercentage: 50,
		stats:                    Up,
		minTotalCount:            5,
		retried:                  false,
	}
}

func (h *Hystrix) SetWindowPoint() {
	h.windowPoint = time.Now()
}

func (h *Hystrix) ResetWindowPoint()  {
	h.windowPoint = time.Time{}
}

func (h *Hystrix) SetRollingWindowPoint()  {
	h.rollingWindowPoint = time.Now()
}

func (h *Hystrix) ResetRollingWindowPoint()  {
	h.rollingWindowPoint = time.Time{}
}

func (h *Hystrix) OnCall() bool {
	if h.stats == Circle_Open {
		if time.Now().Sub(h.rollingWindowPoint) < h.rollingWindowSpan {
			println(time.Now().Sub(h.rollingWindowPoint).Seconds())
			// 滑动窗口内可以放行一个请求
			if h.retried == false {
				// 放行
				h.retried = true
				return true
			} else {
				return false
			}
		} else {
			// 滑动窗口外，重置滑动窗口
			h.SetRollingWindowPoint()
			h.retried = true
			return true
		}
	} else {
		return true
	}
}

func (h *Hystrix) OnSuccess()  {
	if h.stats == Up {
		if h.windowPoint.Equal(time.Time{}) == false &&
			time.Now().Sub(h.windowPoint) <= h.windowSpan {
			// 在时间窗内则计数
			atomic.AddInt32(&h.normalCount,1)
		}
	} else if h.stats == Circle_Open {
		// 访问成功则恢复状态
		h.ResetWindowPoint()
		h.ResetRollingWindowPoint()
		h.stats = Up
		h.retried = false
		h.errorCount = 0
		h.normalCount = 0
	}
}

func (h *Hystrix) OnError() {
	atomic.AddInt32(&h.errorCount,1)
	if h.stats == Up {
		if h.windowPoint.Equal(time.Time{}) {
			h.SetWindowPoint()
		} else if time.Now().Sub(h.windowPoint) > h.windowSpan {
			// 超出了时间窗且之前没有跳闸，清除计数
			atomic.StoreInt32(&h.errorCount, 0)
			atomic.StoreInt32(&h.normalCount, 0)
			h.ResetWindowPoint()
		} else if (h.errorCount + h.normalCount) > h.minTotalCount &&
			(h.errorCount / (h.normalCount + h.errorCount) ) * 100 > h.errorThresholdPercentage {
			// 在时间窗内超出最小请求数并且错误比例达到阈值，跳闸
			h.stats = Circle_Open
			h.SetRollingWindowPoint()
		}
	}
}

func (h *Hystrix) Wrap(f func(w http.ResponseWriter, r *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		if h.OnCall() {
		   if err := f(writer, request); err == nil {
			   h.OnSuccess()
		   } else {
			   h.OnError()
			   fmt.Fprintf(writer, "err: %#v", h)
		   }
		} else {
			fmt.Fprintf(writer, "拒绝服务")
		}
	}
}