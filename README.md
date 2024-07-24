This is a repository for the interview task.
Please read carefully for the below task requirements.

# How to submit your works
1. Fork this repository.
2. Implement the requirement and put your code to the corresponding folder. for example, the codes for task 1 should put into the task1 folder.
3. Send the Github link of your repository to the people who contacted you.

# Task 1
Please implement a todo list HTTP web service.
**No UI** is required.
You can choose any programming language **except for** Java and Python.
The app should provide HTTP Restful API to:
* add a task
* delete a task
* list all tasks by due day
* Modify a task
* Read details of a task

No database is required. You can just have a simple local text file to store the task data.
There should be a proper username and password checking in the API.
There is no requirement for multiple users.
This HTTP server should be easy to run in local.

# Task 2
There is a IP-location csv file in the folder.
For example, there is such one line in the file
```csv
13.104.184.80/28,CN,,BEIJING,
```
the `13.104.184.80/28` is the IP range, `CN` is the country and the `BEIJING` is the city.
There are lots of China city in this file.
Please write a Bash script to read this file and output the all the Chinese cities as a list.
One city can only appear once in the output. 

> (optional requirement) In the Chinese city list, put a number alongside withe the city name to indicate that how many lines appear for this City. For example, 
99 BEIJING
18 SHANGHAI
...
(format is not important, only need the number to count the lines)


# Where to ask questions for the tasks?
You can send email to kennethzhang@outlook.com