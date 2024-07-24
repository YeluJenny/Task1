package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// TaskList 任务列表
var tasks = []Task{}

// LoadTasks 加载文件并返回struct[]
func LoadTasks() ([]Task, error) {

	file, err := os.Open(tasksFile)

	if err != nil {
		dir, err := os.Getwd()
		fmt.Println("os.Getwd() :", dir)
		return tasks, err
	}
	defer file.Close()
	// 获取文件信息
	info, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return nil, nil
	}

	// 判断文件是否为空
	if info.Size() == 0 {
		fmt.Println("The file is empty.")
		return nil, nil
	}

	scanner := bufio.NewScanner(file)
	var currentTask Task
	tasks = []Task{}
	layout := "2006-01-02 15:04:05 -0700 MST"
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ID:") {
			idStr := strings.TrimSpace(strings.TrimPrefix(line, "ID:"))
			currentTask.ID, err = strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("ID 解析错误:", err)
				return nil, err
			}
		} else if strings.HasPrefix(line, "Title:") {
			currentTask.Title = strings.TrimSpace(strings.TrimPrefix(line, "Title:"))
		} else if strings.HasPrefix(line, "Description:") {
			currentTask.Description = strings.TrimSpace(strings.TrimPrefix(line, "Description:"))
		} else if strings.HasPrefix(line, "DueDate:") {
			var dueDateStr = strings.TrimSpace(strings.TrimPrefix(line, "DueDate:"))
			dueDate, err := time.Parse(layout, dueDateStr)
			if err != nil {
				fmt.Println("解析错误:", err)
				return nil, err
			}
			currentTask.DueDate = dueDate
		} else if strings.HasPrefix(line, "Completed:") {
			completedStr := strings.TrimSpace(strings.TrimPrefix(line, "Completed:"))
			completedInt, err := strconv.Atoi(completedStr)
			if err != nil {
				fmt.Println("ID 解析错误:", err)
				return nil, err
			}
			currentTask.Completed = completedInt != 0
			//到底了添加进数组
			tasks = append(tasks, currentTask)
		}
	}

	if err := scanner.Err(); err != nil {
		return tasks, err
	}
	return tasks, nil
}

// AddTask 添加任务
func AddTask(title string, description string, dueDate time.Time, completed bool) error {
	var completedStr string

	nextID, err := GetNextId()
	if err != nil {
		return err
	}

	//更新全局变量
	task := Task{
		ID:          nextID,
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Completed:   completed,
	}
	tasks = append(tasks, task)

	//bool转str存入文件
	if completed {
		completedStr = "1"
	} else {
		completedStr = "0"
	}

	//字符串拼接
	fullText := fmt.Sprintf("ID: %s\nTitle: %s\nDescription: %s\nDueDate: %s\nCompleted: %s\n",
		strconv.Itoa(nextID), title, description, dueDate, completedStr)

	//写入文件
	line := fmt.Sprintf(fullText)
	file, err := os.OpenFile(tasksFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return err
	}
	defer file.Close()
	_, err = file.WriteString(line)
	if err != nil {
		log.Fatalf("写入失败: %v", err)
		return err
	}
	return nil
}

// 获取下一个ID
func GetNextId() (int, error) {
	length := len(tasks)
	// 判断文件是否为空
	if length == 0 {
		return one, nil
	}
	lastID := tasks[length-1].ID
	lastID++
	return lastID, nil
}

// FindTask 根据 ID 查找任务
func FindTask(id int) *Task {
	if len(tasks) == 0 {
		_, err := LoadTasks()
		if err != nil {
			return nil
		}
	}
	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i]
		}
	}
	return nil
}

// DeleteTask 根据 ID 删除任务
func DeleteTask(id int) (bool, error) {
	for i := range tasks {
		if tasks[i].ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}

	err := refreshFile()
	if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateTask 更新任务
func UpdateTask(id int, title, description string, dueDate time.Time, completed bool) (bool, error) {
	var index int
	for i := range tasks {
		if i == id {
			index = i
		}
	}

	//改tasks中的值
	if index != -1 {
		tasks[index].Title = title
		tasks[index].Description = description
		tasks[index].DueDate = dueDate
		tasks[index].Completed = completed
	}

	//写入
	err := refreshFile()
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// 直接替换文件内容
// emm一种偷懒写法，意思意思
func refreshFile() error {
	// 创建一个新文件替换
	file, err := os.Create(tasksFile)
	if err != nil {
		log.Fatalf("Error CREATE file: %v", err)
		return err
	}
	defer file.Close()

	var fullText string

	for i := range tasks {
		//bool转str存入文件
		var completedStr string
		if tasks[i].Completed {
			completedStr = "1"
		} else {
			completedStr = "0"
		}
		//字符串拼接
		fullText += fmt.Sprintf("ID: %s\nTitle: %s\nDescription: %s\nDueDate: %s\nCompleted: %s\n",
			strconv.Itoa(tasks[i].ID), tasks[i].Title, tasks[i].Description, tasks[i].DueDate, completedStr)
	}

	//写入
	_, err = file.WriteString(fullText)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
		return err
	}

	//落盘
	err = file.Sync()
	if err != nil {
		log.Fatalf("Error syncing file: %v", err)
		return err
	}

	return nil
}
