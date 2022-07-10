package thread

import (
	"sync"
)

type GoLimit struct {
	max       uint             //并发最大数量
	count     uint             //当前已有并发数
	isAddLock bool             //是否已锁定增加
	zeroChan  chan interface{} //为0时广播
	addLock   sync.Mutex       //(增加并发数的)锁
	dataLock  sync.Mutex       //(修改数据的)锁
}

func NewGoLimit(max uint) *GoLimit {
	return &GoLimit{max: max, count: 0, isAddLock: false, zeroChan: nil}
}

func (g *GoLimit) Add() {
	g.addLock.Lock()
	g.dataLock.Lock()

	g.count += 1

	if g.count < g.max { //未超并发时解锁,后续可以继续增加
		g.addLock.Unlock()
	} else { //已到最大并发数, 不解锁并标记. 等数量减少后解锁
		g.isAddLock = true
	}

	g.dataLock.Unlock()
}

// Done 并发计数减1
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
