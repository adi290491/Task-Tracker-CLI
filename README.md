# Task-Tracker-CLI
This is a basic Task Tracker CLI developed in Go. This project is a minimalistic one in the sense that there are no external libraries, or interaction with databases to manage your tasks. All the tasks are stored in a tasks.json file which is created when you first add a task.

## FEATURES
1. Add
   ```bash
   task-cli add <description>
   ```
2. Update
   ```bash
   task-cli update <id> <new-description>
   ```
4. Delete
   ```bash
   task-cli delete <id>
   ```
5. Change progress status
   ```bash
   task-cli mark-in-progress <id>
   ```
   ```bash
   task-cli mark-done <id>
   ```
6. List Tasks
   ```bash
   task-cli list [optional status]
   options=[todo|in-progress|done]
   ```


   
