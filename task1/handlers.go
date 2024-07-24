package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// ByDateDesc 定义一个自定义类型来实现 sort.Interface 接口
type ByDateDesc []Task

func (a ByDateDesc) Len() int           { return len(a) }
func (a ByDateDesc) Less(i, j int) bool { return a[i].DueDate.After(a[j].DueDate) } // 从新到旧排序
func (a ByDateDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// BasicAuth 验证用户名和密码
func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Auth")
		user, pass, ok := r.BasicAuth()
		if !ok {
			// 未提供基本身份验证，返回 401 Unauthorized
			w.Header().Set("WWW-Authenticate", `Basic realm="Unauthorized"`)
			http.Error(w, "请先输入用户名和密码！", http.StatusUnauthorized)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if user != username || pass != password {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, "用户名或密码不正确！", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

/*
AddTaskHandler 处理添加任务请求
*/
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "非法请求！", http.StatusMethodNotAllowed)
		return
	}

	var task Task
	//json解码不匹配task数据模型
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "请求格式有误！", http.StatusBadRequest)
		return
	}

	//更改时区
	task.DueDate = task.DueDate.UTC()
	if err := AddTask(task.Title, task.Description, task.DueDate, task.Completed); err != nil {
		http.Error(w, "添加任务失败！", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = statusOKResponse(w)
	if err != nil {
		fmt.Println("写入响应体失败:", err)
	}
}

/*
DeleteTaskHandler 处理删除任务请求
*/
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "非法请求！", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/delete_task_by_id/")
	id, err := strconv.Atoi(idStr)
	if err != nil || FindTask(id) == nil {
		http.Error(w, "找不到该任务！", http.StatusNotFound)
		return
	}

	task, err := DeleteTask(id)
	if !task || err != nil {
		http.Error(w, "删除失败！", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = statusOKResponse(w)
	if err != nil {
		fmt.Println("写入响应体失败:", err)
	}
}

/*
ListTasksByDueDateHandler 处理按截止日期列出任务请求
*/
func ListTasksByDueDateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ListTasksByDueDateHandler")
	if r.Method != http.MethodGet {
		http.Error(w, "非法请求！", http.StatusMethodNotAllowed)
		return
	}

	var result []TaskRet
	var completed string
	sort.Sort(ByDateDesc(tasks))

	for _, t := range tasks {
		if t.Completed {
			completed = "yes"
		} else {
			completed = "no"
		}
		taskRet := TaskRet{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			DueDate:     t.DueDate.Format("2006-01-02 15:04"),
			Completed:   completed,
		}
		result = append(result, taskRet)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

/*
ModifyTaskHandler 处理更新任务请求
*/
func ModifyTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "非法请求！", http.StatusMethodNotAllowed)
		return
	}
	var task Task
	idStr := strings.TrimPrefix(r.URL.Path, "/update-task/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "请输入正确的ID！", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "请求格式有误！", http.StatusBadRequest)
		_, err2 := statusBadResponse(w, "请求格式有误！")
		if err2 != nil {
			return
		}
		return
	}

	if FindTask(id) == nil {
		http.Error(w, "任务不存在！", http.StatusNotFound)
		_, _ = statusBadResponse(w, "任务不存在！")
		return
	}

	task.DueDate = task.DueDate.UTC()
	updateTask, err := UpdateTask(id, task.Title, task.Description, task.DueDate, task.Completed)
	if !updateTask || err != nil {
		http.Error(w, "任务更改失败！", http.StatusBadRequest)
		_, err := statusBadResponse(w, "任务更改失败！")
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = statusOKResponse(w)
	if err != nil {
		fmt.Println("写入响应体失败:", err)
	}
}

/*
GetTaskHandler 处理读取任务的请求
*/
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "非法请求！", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/task/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id格式有误！", http.StatusBadRequest)
		return
	}

	task := FindTask(id)
	if task == nil {
		http.Error(w, "任务不存在！", http.StatusNotFound)
		return
	}

	var completedStr string
	if task.Completed {
		completedStr = "yes"
	} else {
		completedStr = "no"
	}
	taskRet := TaskRet{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate.Format("2006-01-02 15:04"),
		Completed:   completedStr,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(taskRet)
}

// 构造请求成功返回
func statusOKResponse(w http.ResponseWriter) ([]byte, error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Response{Message: "操作成功！"}
	responseBytes, _ := json.Marshal(response)
	_, err := w.Write(responseBytes)
	if err != nil {
		fmt.Println("写入响应体失败:", err)
		return nil, err
	}
	return responseBytes, nil
}

// 构造请求成功返回
func statusBadResponse(w http.ResponseWriter, message string) ([]byte, error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Response{Message: message}
	responseBytes, _ := json.Marshal(response)
	_, err := w.Write(responseBytes)
	if err != nil {
		fmt.Println("写入响应体失败:", err)
		return nil, err
	}
	return responseBytes, nil
}

func main() {
	if _, err := LoadTasks(); err != nil {
		fmt.Println("任务加载失败！:", err)
		return
	}

	http.HandleFunc("/add_task", BasicAuth(AddTaskHandler))
	http.HandleFunc("/delete_task_by_id/", BasicAuth(DeleteTaskHandler))
	http.HandleFunc("/list-tasks", BasicAuth(ListTasksByDueDateHandler))
	http.HandleFunc("/update-task/", BasicAuth(ModifyTaskHandler))
	http.HandleFunc("/task/", BasicAuth(GetTaskHandler))

	fmt.Println("服务启动", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println("服务启动报错:", err)
	}
}
