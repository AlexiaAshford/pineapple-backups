package multi

import (
	"sync"
)

type GoLimit struct {
	max       uint             //max concurrent goroutine
	count     uint             //now concurrent goroutine
	isAddLock bool             //is lock add
	zeroChan  chan interface{} //the channel when count is 0
	addLock   sync.Mutex       //add lock
	dataLock  sync.Mutex       //add data lock
}

func NewGoLimit(max uint) *GoLimit {
	return &GoLimit{max: max, count: 0, isAddLock: false, zeroChan: nil}
}

func (g *GoLimit) Add() {
	g.addLock.Lock()
	g.dataLock.Lock()

	g.count += 1

	if g.count < g.max { //if count < max, can unlock
		g.addLock.Unlock()
	} else { //if count >= max, can't unlock
		g.isAddLock = true
	}

	g.dataLock.Unlock()
}

// Done then count - 1
//若计数<max_num, 可以使原阻塞的Add()快速解除阻塞
func (g *GoLimit) Done() {
	g.dataLock.Lock()

	g.count -= 1

	//解锁
	if g.isAddLock == true && g.count < g.max {
		g.isAddLock = false
		g.addLock.Unlock()
	}

	//0广播
	if g.count == 0 && g.zeroChan != nil {
		close(g.zeroChan)
		g.zeroChan = nil
	}

	g.dataLock.Unlock()
}

// SetMax 更新最大并发计数为, 若是调大, 可以使原阻塞的Add()快速解除阻塞
func (g *GoLimit) SetMax(n uint) {
	g.dataLock.Lock()
	g.max = n

	//解锁
	if g.isAddLock == true && g.count < g.max {
		g.isAddLock = false
		g.addLock.Unlock()
	}

	//加锁
	if g.isAddLock == false && g.count >= g.max {
		g.isAddLock = true
		g.addLock.Lock()
	}

	g.dataLock.Unlock()
}

// WaitZero 若当前并发计数为0, 则快速返回; 否则阻塞等待,直到并发计数为0
func (g *GoLimit) WaitZero() {
	g.dataLock.Lock()

	//无需等待
	if g.count == 0 {
		g.dataLock.Unlock()
		return
	}

	//无广播通道, 创建一个
	if g.zeroChan == nil {
		g.zeroChan = make(chan interface{})
	}

	//复制通道后解锁, 避免从nil读数据
	c := g.zeroChan
	g.dataLock.Unlock()

	<-c
}

// Count 获取并发计数
func (g *GoLimit) Count() uint {
	return g.count
}

// Max 获取最大并发计数
func (g *GoLimit) Max() uint {
	return g.max
}
