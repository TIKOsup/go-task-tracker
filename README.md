# Task Tracker CLI

A simple task management CLI app written in Go and using Cobra CLI.

## Features

- Add, Update, and Delete tasks
- Mark a task as in progress or done
- List all tasks
- List all tasks that are done
- List all tasks that are in progress
- List all tasks that are todo



## Installation

1. Clone the repository

```bash
git clone https://github.com/TIKOsup/go-task-tracker.git
cd go-task-tracker
```

2. Build an app

```bash
go build -o task-cli
```


## Usage

The list of commands and their usage is given below:

```bash
# Adding a new task
./task-cli add "Buy groceries"

# Listing all tasks
./task-cli list

# Listing tasks by status
./task-cli list done
./task-cli list todo
./task-cli list in-progress

# Updating tasks' description
./task-cli update 1 "Buy groceries and cook dinner"

# Marking a task as in progress or done
task-cli mark-in-progress 1
task-cli mark-done 1

# Deleting tasks
./task-cli delete 1
```

## Data Structure

```json
{
 "tasks": [
  {
   "id": 1,
   "description": "Buy groceries",
   "status": "todo",
   "createdAt": "2026-06-05T10:00:00",
   "updatedAt": "2026-06-05T12:00:00"
  }
 ]
}
```

## Possible improvements

- [ ] SQLite storage
- [ ] Unit tests
- [ ] Task sorting
- [ ] Search tasks
- [ ] Task deadlines
- [ ] Colored output

