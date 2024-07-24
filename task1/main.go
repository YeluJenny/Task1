package main

import (
	"encoding/json"
	"time"
)

// Task 任务数据模型
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	DueDate     time.Time `json:"due_date"`
}

// 自定义 JSON 序列化
func (t Task) MarshalJSON() ([]byte, error) {
	type Alias Task // 防止递归调用
	return json.Marshal(&struct {
		DueDate string `json:"due_date"`
	}{
		DueDate: t.DueDate.Format("2006-01-02 15:04"), // 自定义格式
	})
}

// 自定义 JSON 反序列化
func (t *Task) UnmarshalJSON(data []byte) error {
	type Alias Task // 防止递归调用
	aux := &struct {
		DueDate string `json:"due_date"`
		*Alias
	}{
		Alias: (*Alias)(t), // 将当前 Task 赋值给 Alias
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// 解析 DueDate
	layout := "2006-01-02 15:04"
	dueDate, err := time.Parse(layout, aux.DueDate)
	if err != nil {
		return err
	}
	t.DueDate = dueDate

	return nil
}

type TaskRet struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   string `json:"completed"`
	DueDate     string `json:"due_date"`
}

type Response struct {
	Message string `json:"message"`
}
