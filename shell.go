package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorBlue   = "\033[34m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		// Get current directory
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		// Get username
		username := os.Getenv("USER")
		if username == "" {
			username = "user"
		}

		// Format prompt with colors
		prompt := fmt.Sprintf("%s%s%s:%s%s%s$ ",
			colorGreen, username, colorReset,
			colorBlue, shortenPath(currentDir), colorReset)
		fmt.Print(prompt)

		// Read the keyboard input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// Handle the execution of the input
		if err = execInput(input); err != nil {
			fmt.Fprintf(os.Stderr, "%sError: %v%s\n", colorYellow, err, colorReset)
		}
	}
}

// shortenPath replaces home directory with ~ and shortens the path
func shortenPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	// Replace home directory with ~
	if strings.HasPrefix(path, home) {
		path = "~" + path[len(home):]
	}
	return path
}

var ErrNoPath = errors.New("path required")

func execInput(input string) error {

	input = strings.TrimSuffix(input, "\n")

	args := strings.Split(input, " ")

	if len(args) == 0 || args[0] == "" {
		return nil
	}

	switch args[0] {
	case "cd":
		if len(args) < 2 {
			// CD to home directory if no path is given
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			return os.Chdir(homeDir)
		}
		// Change the directory and return the error.
		return os.Chdir(args[1])
	case "exit":
		os.Exit(0)
	case "pwd":
		// Add pwd command to print current directory
		currentDir, err := os.Getwd()
		if err != nil {
			return err
		}
		fmt.Println(currentDir)
		return nil
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
