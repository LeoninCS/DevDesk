package service

import (
	"fmt"
	"sync"
)

type WorkPlan struct {
	mu   sync.RWMutex
	Plan map[string]*PersonalPlan
}

type PersonalPlan struct {
	mu    sync.RWMutex
	Hash  string    `json:"hash"`
	TODOs [200]TODO `json:"todos"`
	sta   []int     `json:"-"` // 空闲槽位池
}

type TODO struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

func NewWorkPlan() *WorkPlan {
	return &WorkPlan{
		Plan: make(map[string]*PersonalPlan),
	}
}

func (wp *WorkPlan) NewPersonalPlan() *PersonalPlan {
	pp := &PersonalPlan{
		Hash: GetHash(10),
		sta:  make([]int, 200),
	}

	// 初始化空闲槽位列表
	for i := 0; i < 200; i++ {
		pp.sta[i] = i
	}

	wp.mu.Lock()
	wp.Plan[pp.Hash] = pp
	wp.mu.Unlock()

	return pp
}

func (wp *WorkPlan) GetPlan(hash string) *PersonalPlan {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.Plan[hash]
}

// ---------------- TODO 逻辑 ----------------

func (pp *PersonalPlan) AddTODO(content string) error {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	if len(pp.sta) == 0 {
		return fmt.Errorf("TODO is too many")
	}

	id := pp.sta[len(pp.sta)-1]
	pp.sta = pp.sta[:len(pp.sta)-1]

	pp.TODOs[id] = TODO{
		Id:      id,
		Content: content,
		Done:    false,
	}
	return nil
}

func (pp *PersonalPlan) DeleteTODO(id int) error {
	if id < 0 || id >= 200 {
		return fmt.Errorf("invalid TODO id")
	}

	pp.mu.Lock()
	defer pp.mu.Unlock()

	// 回收槽位
	pp.sta = append(pp.sta, id)

	// 清空内容
	pp.TODOs[id] = TODO{}
	return nil
}

func (pp *PersonalPlan) EditTODO(id int, content string) error {
	if id < 0 || id >= 200 {
		return fmt.Errorf("invalid TODO id")
	}

	pp.mu.Lock()
	defer pp.mu.Unlock()

	if pp.TODOs[id].Content == "" {
		return fmt.Errorf("TODO not exists")
	}

	pp.TODOs[id].Content = content
	return nil
}

func (pp *PersonalPlan) SetTODODone(id int) error {
	if id < 0 || id >= 200 {
		return fmt.Errorf("invalid TODO id")
	}

	pp.mu.Lock()
	defer pp.mu.Unlock()

	if pp.TODOs[id].Content == "" {
		return fmt.Errorf("TODO not exists")
	}

	pp.TODOs[id].Done = !pp.TODOs[id].Done
	return nil
}

func (pp *PersonalPlan) GetTODOs() []TODO {
	pp.mu.RLock()
	defer pp.mu.RUnlock()

	todos := make([]TODO, 0, 200)

	// 未完成的在前
	for _, t := range pp.TODOs {
		if t.Content != "" && !t.Done {
			todos = append(todos, t)
		}
	}
	// 已完成的在后
	for _, t := range pp.TODOs {
		if t.Content != "" && t.Done {
			todos = append(todos, t)
		}
	}

	return todos
}
