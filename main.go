package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"gopkg.in/yaml.v2"
)

var yamlFile string

func init() {
	flag.StringVar(&yamlFile, "f", "", "YAML file containing tasks")
	flag.Parse()
}

type Task struct {
	Command struct {
		Val    interface{} `yaml:"val"`
		Python bool        `yaml:"python"`
	} `yaml:"command"`
}

func executeCommand(command string, python bool) error {
	if python {
		command = "poetry run python" + command

	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func executeCommands(commands []string, python bool) error {
	var wg sync.WaitGroup
	for _, command := range commands {
		wg.Add(1)
		go func(command string) {
			defer wg.Done()
			if err := executeCommand(command, python); err != nil {
				fmt.Println("Error:", err)
			}
		}(command)
	}
	wg.Wait()
	return nil
}

func runTask(task Task) {
	if val, ok := task.Command.Val.(string); ok {
		if err := executeCommand(val, task.Command.Python); err != nil {
			fmt.Println("Error:", err)
		}
	} else if valList, ok := task.Command.Val.([]interface{}); ok {
		var commands []string
		for _, val := range valList {
			if val, ok := val.(string); ok {
				commands = append(commands, val)
			}
		}
		if err := executeCommands(commands, task.Command.Python); err != nil {
			fmt.Println("Error:", err)
		}

	}
}

func main() {
	if yamlFile == "" {
		fmt.Println("Please specify a YAML file using the -f flag.")
		return
	}

	data, err := os.ReadFile(yamlFile)
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		return
	}

	var tasks map[string][]Task
	if err := yaml.Unmarshal(data, &tasks); err != nil {
		fmt.Println("Error unmarshaling YAML:", err)
		return
	}

	if run, ok := tasks["run"]; ok {
		for _, task := range run {
			runTask(task)
		}
	}
}
