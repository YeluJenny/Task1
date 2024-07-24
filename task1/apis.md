Auth:
admin
password


* add a task
localhost:8080/add_task
{
  "title": "任务标题",
  "description": "任务描述",
  "completed": false,
  "due_date": "2023-11-12 00:00"
}

* delete a task
localhost:8080/delete_task_by_id/1

* list all tasks by due day
localhost:8080/list-tasks

* Modify a task
localhost:8080/update-task/1
body:
{
  "title": "任务标题",
  "description": "任务描述",
  "completed": false,
  "due_date": "2023-11-12 00:00"
}

* Read details of a task
    localhost:8080/task/1


go version go1.22.5 windows/amd64