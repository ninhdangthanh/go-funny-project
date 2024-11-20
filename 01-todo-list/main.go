package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/google/uuid"
)

func loadFile(filepath string) (*os.File, error) {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading")
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		_ = f.Close()
		return nil, err
	}

	return f, nil
}

func closeFile(f *os.File) error {
	syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	return f.Close()
}

func isFileEmpty(file *os.File) (bool, error) {
	info, err := file.Stat()
	if err != nil {
		return false, fmt.Errorf("failed to get file info: %v", err)
	}

	return info.Size() == 0, nil
}

func writeHeader(file *os.File) error {
	_, err := file.WriteString("ID, Task, Created, Done\n")
	if err != nil {
		return fmt.Errorf("failed to write header: %v", err)
	}
	return nil
}

func addTask(filepath, task string) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open file for appending: %v", err)
	}
	defer file.Close()

	timestamp := time.Now().Format(time.RFC3339)
	id := uuid.New().String()

	line := fmt.Sprintf("%s,%s,%s,%s\n", id, task, timestamp, "false")

	if _, err := file.WriteString(line); err != nil {
		return fmt.Errorf("failed to add task: %v", err)
	}
	fmt.Println("Task added:", task, "ID:", id)
	return nil
}

func readTasks(filepath string) ([][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var records [][]string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		records = append(records, strings.Split(line, ","))
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	return records, nil
}

func showTasks(filepath string) error {
	records, err := readTasks(filepath)
	if err != nil {
		return fmt.Errorf("failed to read tasks: %v", err)
	}

	if len(records) <= 1 {
		fmt.Println("No tasks to show.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "ID\tTask\tCreated\tDone\t")

	for i, record := range records {
		if i == 0 {
			continue
		}
		id, task, created, done := record[0], record[1], record[2], record[3]
		createdTime, _ := time.Parse(time.RFC3339, created)
		relativeTime := formatRelativeTime(createdTime)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", id, task, relativeTime, done)
	}
	w.Flush()
	return nil
}

func formatRelativeTime(t time.Time) string {
	duration := time.Since(t)
	switch {
	case duration < time.Minute:
		return fmt.Sprintf("%d seconds ago", int(duration.Seconds()))
	case duration < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	default:
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	}
}

func doneTask(filepath, id string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file for reading: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	var updatedLines []string
	taskFound := false
	for _, line := range lines {
		fields := strings.Split(line, ",")
		if len(fields) < 4 {
			updatedLines = append(updatedLines, line)
			continue
		}

		if fields[0] == id {
			fields[3] = "true"
			taskFound = true
		}
		updatedLines = append(updatedLines, strings.Join(fields, ","))
	}

	if !taskFound {
		return fmt.Errorf("task with ID %s not found", id)
	}

	file, err = os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range updatedLines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}
	writer.Flush()

	fmt.Println("Task with ID", id, "marked as done.")
	return nil
}

func removeTask(filepath, id string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file for reading: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	var updatedLines []string
	taskFound := false
	for _, line := range lines {
		fields := strings.Split(line, ",")
		if len(fields) < 4 {
			updatedLines = append(updatedLines, line)
			continue
		}

		if fields[0] == id {
			taskFound = true
			continue
		}

		updatedLines = append(updatedLines, line)
	}

	if !taskFound {
		return fmt.Errorf("task with ID %s not found", id)
	}

	file, err = os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range updatedLines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}
	writer.Flush()

	fmt.Println("Task with ID", id, "removed.")
	return nil
}

func main() {

	var path_file_csv string
	path_file_csv = "./data.csv"

	file, err := loadFile(path_file_csv)
	if err != nil {
		fmt.Printf("Error when loading file: %v\n", err)
		return
	}

	defer closeFile(file)

	isEmpty, err := isFileEmpty(file)
	if err != nil {
		fmt.Printf("Error checking if file is empty: %v\n", err)
		return
	}

	if isEmpty {
		err = writeHeader(file)
		if err != nil {
			// fmt.Printf("Error writing header: %v\n", err)
			return
		}
		// fmt.Println("Header added to the empty file.")
	} else {
		// fmt.Println("File is not empty; no header added.")
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . [add <task>|remove <id>]")
		return
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run . add <task>")
			return
		}
		task := strings.Join(os.Args[2:], " ")
		if err := addTask(path_file_csv, task); err != nil {
			fmt.Printf("Error adding task: %v\n", err)
		}
	case "show":
		if err := showTasks(path_file_csv); err != nil {
			fmt.Printf("Error showing tasks: %v\n", err)
		}
	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run . done <id>")
			return
		}
		id := os.Args[2]
		if err := doneTask(path_file_csv, id); err != nil {
			fmt.Printf("Error marking task as done: %v\n", err)
		}
	case "remove":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run . remove <id>")
			return
		}
		id := os.Args[2]
		if err := removeTask(path_file_csv, id); err != nil {
			fmt.Printf("Error removing task: %v\n", err)
		}

	default:
		fmt.Println("Unknown command. Usage: go run . [add <task>|remove <id>]")
	}

}
