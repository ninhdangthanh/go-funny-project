Here’s the complete `README.md` file covering all aspects of your project:

---

```markdown
# Task Manager CLI

A simple command-line task manager written in Go. This project allows you to manage tasks stored in a CSV file, supporting operations such as adding tasks, viewing tasks, marking tasks as done, and removing tasks.

---

## Features

- Add tasks with a unique ID and timestamp.
- View all tasks in a formatted table.
- Mark tasks as done.
- Remove tasks by ID.
- Automatically creates a `data.csv` file if it doesn’t exist.

---

## Installation

### Prerequisites

- [Go](https://go.dev/dl/) must be installed.

### Clone the Repository
```bash
git clone <repository-url>
cd <repository-name>
```

### Run the Program
Use the `go run` command to execute the program.

---

## Usage

### 1. **Initialize the CSV File**
The program creates a `data.csv` file in the current directory if it doesn't exist. The file will contain the following columns:

| **ID**                                | **Task**          | **Created**           | **Done** |
|---------------------------------------|-------------------|-----------------------|----------|
| A unique UUID for each task           | Task description  | Timestamp in RFC3339  | `true` or `false` |

---

### 2. **Commands**

#### Add a Task
Add a task to the list.

```bash
go run . add "Task description"
```

**Example**:
```bash
go run . add "Learn Go programming"
```

**Output**:
```
Task added: Learn Go programming ID: 123e4567-e89b-12d3-a456-426614174000
```

---

#### Show Tasks
Display all tasks in a formatted table.

```bash
go run . show
```

**Example Output**:
```
ID                                    Task                     Created               Done
123e4567-e89b-12d3-a456-426614174000  Learn Go programming     1 hour ago           false
```

---

#### Mark a Task as Done
Mark a task as done by its ID.

```bash
go run . done <task-id>
```

**Example**:
```bash
go run . done 123e4567-e89b-12d3-a456-426614174000
```

**Output**:
```
Task with ID 123e4567-e89b-12d3-a456-426614174000 marked as done.
```

**Updated Task**:
```
ID                                    Task                     Created               Done
123e4567-e89b-12d3-a456-426614174000  Learn Go programming     1 hour ago           true
```

---

#### Remove a Task
Remove a task from the list by its ID.

```bash
go run . remove <task-id>
```

**Example**:
```bash
go run . remove 123e4567-e89b-12d3-a456-426614174000
```

**Output**:
```
Task with ID 123e4567-e89b-12d3-a456-426614174000 removed.
```

---

## File Format

The tasks are stored in a `data.csv` file with the following structure:

```csv
ID,Task,Created,Done
123e4567-e89b-12d3-a456-426614174000,Learn Go programming,2024-11-20T14:23:45Z,false
223e4567-e89b-12d3-a456-426614174001,Practice problem-solving,2024-11-20T15:10:12Z,false
```

---

## Error Handling

- **Task Not Found**: If a task ID is not found for `done` or `remove` commands, the program will return an error:
  ```
  Error: task with ID <id> not found.
  ```

- **No Tasks to Show**: If no tasks are available when running `show`, it will display:
  ```
  No tasks to show.
  ```

- **File I/O Errors**: The program handles file-related errors and outputs appropriate messages if something goes wrong.

---

## Contributing

Feel free to fork this repository and make improvements. Contributions are welcome!

1. Fork the project.
2. Create a feature branch (`git checkout -b feature-name`).
3. Commit your changes (`git commit -am 'Add feature'`).
4. Push to the branch (`git push origin feature-name`).
5. Create a pull request.

---

## License

This project is licensed under the MIT License. See `LICENSE` for more information.

---

## Examples

### Adding a Task

Command:
```bash
go run . add "Write documentation"
```

Output:
```
Task added: Write documentation ID: 223e4567-e89b-12d3-a456-426614174001
```

---

### Showing Tasks

Command:
```bash
go run . show
```

Output:
```
ID                                    Task                     Created               Done
223e4567-e89b-12d3-a456-426614174001  Write documentation      2 minutes ago        false
```

---

### Marking a Task as Done

Command:
```bash
go run . done 223e4567-e89b-12d3-a456-426614174001
```

Output:
```
Task with ID 223e4567-e89b-12d3-a456-426614174001 marked as done.
```

---

### Removing a Task

Command:
```bash
go run . remove 223e4567-e89b-12d3-a456-426614174001
```

Output:
```
Task with ID 223e4567-e89b-12d3-a456-426614174001 removed.
```
```

---

### Save this as `README.md` in your project directory.

This file provides detailed usage instructions, examples, and contribution guidelines, making it comprehensive for anyone using or contributing to your project.