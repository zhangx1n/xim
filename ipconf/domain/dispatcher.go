package domain

import (
	"github.com/zhangx1n/xim/ipconf/source"
	"sort"
	"sync"
)

// Dispatcher 调度器
type Dispatcher struct {
	candidateTable map[string]*Endport
	sync.RWMutex
}

// dp 全局的Dispatcher
var dp *Dispatcher

func Init() {
	dp = &Dispatcher{}
	dp.candidateTable = make(map[string]*Endport)
}

func Dispatch(ctx *IpConfContext) []*Endport {
	// step1: 获取候选 ip
	eds := dp.getCandidateEndport(ctx)
	// step2: 逐一计算得分
	for _, ed := range eds {
		ed.CalculateScore(ctx)
	}
	// step3: 全局排序, 动静结合的排序策略
	sort.Slice(eds, func(i, j int) bool {
		// 优先基于活跃分数进行排序
		if eds[i].ActiveSorce > eds[j].ActiveSorce {
			return true
		}
		// 如果活跃分数相同，则使用静态分数排序
		if eds[i].ActiveSorce == eds[j].ActiveSorce {
			return eds[i].StaticSorce > eds[j].StaticSorce
		}
		return false
	})
	return eds

}

func (d *Dispatcher) getCandidateEndport(ctx *IpConfContext) []*Endport {
	dp.RLock()
	defer dp.RUnlock()
	candidateList := make([]*Endport, 0, len(dp.candidateTable))
	for _, ed := range dp.candidateTable {
		candidateList = append(candidateList, ed)
	}
	return candidateList
}

func (dp *Dispatcher) delNode(event *source.Event) {
	dp.Lock()
	defer dp.Unlock()
	delete(dp.candidateTable, event.Key())
}
func (dp *Dispatcher) addNode(event *source.Event) {
	dp.Lock()
	defer dp.Unlock()
	var (
		ed *Endport
		ok bool
	)
	if ed, ok = dp.candidateTable[event.Key()]; !ok { // 不存在
		ed = NewEndport(event.IP, event.Port)
		dp.candidateTable[event.Key()] = ed
	}
	ed.UpdateStat(&Stat{
		ConnectNum:   event.ConnectNum,
		MessageBytes: event.MessageBytes,
	})

}
