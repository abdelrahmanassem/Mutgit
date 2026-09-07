package main

import (
	"fmt"
	"os"
	"path/filepath"
)
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mugit <command>")
		return
	}
	command := os.Args[1]
	if command == "init" {
		created, err := initRepository()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if created {
			fmt.Println("Initialized empty Mugit repository")
		} else {
			fmt.Println("Mugit repository already exists")
		}
		return
	}
	fmt.Println("Unknown command:", command)
}
func initRepository() (bool, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return false, fmt.Errorf("could not get current directory: %w", err)
	}
	mugitPath := filepath.Join(currentDir, ".mugit")
	info, err := os.Stat(mugitPath)
	if err == nil {
		if !info.IsDir() {
			return false, fmt.Errorf(".mugit exists but is not a directory")
		}
		return false, nil
	}
	if !os.IsNotExist(err) {
		return false, fmt.Errorf("could not check repository: %w", err)
	}
	err = os.Mkdir(mugitPath, 0755)
	if err != nil {
		return false, fmt.Errorf("could not create .mugit: %w", err)
	}
	return true, nil
}