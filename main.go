package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	if command == "status" {
		err := statusRepository()

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		return
	}

	fmt.Println("Unknown command:", command)
}

func initRepository() (bool, error) {

	currentDir, err := os.Getwd()
	if err != nil {
		return false, fmt.Errorf(
			"could not get current directory: %w",
			err,
		)
	}


	mugitPath := filepath.Join(currentDir, ".mugit")

	
	info, err := os.Stat(mugitPath)

	if err == nil {


		if !info.IsDir() {
			return false, fmt.Errorf(
				".mugit exists but is not a directory",
			)
		}

		return false, nil
	}

	if !os.IsNotExist(err) {
		return false, fmt.Errorf(
			"could not check repository: %w",
			err,
		)
	}

	
	err = os.Mkdir(mugitPath, 0755)
	if err != nil {
		return false, fmt.Errorf(
			"could not create .mugit: %w",
			err,
		)
	}

	
	objectsPath := filepath.Join(
		mugitPath,
		"objects",
	)

	err = os.Mkdir(objectsPath, 0755)
	if err != nil {
		return false, fmt.Errorf(
			"could not create objects directory: %w",
			err,
		)
	}

	
	headsPath := filepath.Join(
		mugitPath,
		"refs",
		"heads",
	)

	err = os.MkdirAll(headsPath, 0755)
	if err != nil {
		return false, fmt.Errorf(
			"could not create refs/heads directory: %w",
			err,
		)
	}

	
	indexPath := filepath.Join(
		mugitPath,
		"index",
	)

	err = os.WriteFile(
		indexPath,
		[]byte{},
		0644,
	)

	if err != nil {
		return false, fmt.Errorf(
			"could not create index: %w",
			err,
		)
	}


	headPath := filepath.Join(
		mugitPath,
		"HEAD",
	)

	headContent := []byte(
		"ref: refs/heads/main\n",
	)

	err = os.WriteFile(
		headPath,
		headContent,
		0644,
	)

	if err != nil {
		return false, fmt.Errorf(
			"could not create HEAD: %w",
			err,
		)
	}

	
	mainBranchPath := filepath.Join(
		headsPath,
		"main",
	)

	err = os.WriteFile(
		mainBranchPath,
		[]byte{},
		0644,
	)

	if err != nil {
		return false, fmt.Errorf(
			"could not create main branch: %w",
			err,
		)
	}

	return true, nil
}

func findRepositoryRoot() (string, error) {
	
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf(
			"could not get current directory: %w",
			err,
		)
	}

	current := currentDir

	
	for {
		
		mugitPath := filepath.Join(
			current,
			".mugit",
		)

		
		_, err := os.Stat(mugitPath)

		if err == nil {
		
			return current, nil
		}

		if !os.IsNotExist(err) {
			return "", fmt.Errorf(
				"could not check %s: %w",
				mugitPath,
				err,
			)
		}


		parent := filepath.Dir(current)

		
		if parent == current {
			return "", fmt.Errorf(
				"not a Mugit repository",
			)
		}

		current = parent
	}
}

func readHEAD(root string) (string, error) {
	headPath := filepath.Join(
		root,
		".mugit",
		"HEAD",
	)

	data, err := os.ReadFile(headPath)
	if err != nil {
		return "", fmt.Errorf(
			"could not read HEAD: %w",
			err,
		)
	}

	return string(data), nil
}

func getCurrentBranch(head string) (string, error) {

	head = strings.TrimSpace(head)

	const prefix = "ref: refs/heads/"

	
	if !strings.HasPrefix(head, prefix) {
		return "", fmt.Errorf(
			"unsupported HEAD format: %s",
			head,
		)
	}


	branch := strings.TrimPrefix(
		head,
		prefix,
	)

	return branch, nil
}

func readBranch(root string, branch string) (string, error) {
	branchPath := filepath.Join(
		root,
		".mugit",
		"refs",
		"heads",
		branch,
	)

	data, err := os.ReadFile(branchPath)
	if err != nil {
		return "", fmt.Errorf(
			"could not read branch %s: %w",
			branch,
			err,
		)
	}

	return strings.TrimSpace(string(data)), nil
}

func statusRepository() error {

	root, err := findRepositoryRoot()
	if err != nil {
		return err
	}

	head, err := readHEAD(root)
	if err != nil {
		return err
	}

	branch, err := getCurrentBranch(head)
	if err != nil {
		return err
	}

	fmt.Println("On branch", branch)

	
	commitID, err := readBranch(
		root,
		branch,
	)

	if err != nil {
		return err
	}

	
	if commitID == "" {
		fmt.Println("No commits yet")
	} else {
		fmt.Println("Current commit:", commitID)
	}

	return nil
}

