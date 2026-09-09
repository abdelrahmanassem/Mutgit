package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

	if command == "hash-object" {
	if len(os.Args) < 3 {
		fmt.Println("Usage: mugit hash-object <file>")
		return
	}

	err := hashObject(os.Args[2])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	return
}

if command == "cat-file" {
	if len(os.Args) < 3 {
		fmt.Println("Usage: mugit cat-file <objectID>")
		return
	}

	err := catFile(os.Args[2])
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

func hashObject(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf(
			"could not read file %s: %w",
			filePath,
			err,
		)
	}

	header := fmt.Sprintf("blob %d\x00", len(data))

	blobData := append([]byte(header), data...)

	hash := sha1.Sum(blobData)

	objectID := fmt.Sprintf("%x", hash)

	fmt.Println("Object ID:", objectID)
	root, err := findRepositoryRoot()
	if err != nil {
		return err
	}

	objectPath := filepath.Join(
		root,
		".mugit",
		"objects",
		objectID,
	)

	_, err = os.Stat(objectPath)

	if err == nil {
		fmt.Println(objectID)
		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf(
			"could not check object: %w",
			err,
		)
	}

	err = os.WriteFile(
		objectPath,
		blobData,
		0644,
	)

	if err != nil {
		return fmt.Errorf(
			"could not store object: %w",
			err,
		)
	}

	fmt.Println(objectID)
	fmt.Println("File bytes:", data)
	fmt.Println("Number of bytes:", len(data))

	return nil
}


func catFile(objectID string) error {
	root, err := findRepositoryRoot()
	if err != nil {
		return err
	}

	objectPath := filepath.Join(
		root,
		".mugit",
		"objects",
		objectID,
	)

	data, err := os.ReadFile(objectPath)
	if err != nil {
		return fmt.Errorf(
			"could not read object %s: %w",
			objectID,
			err,
		)
	}

	separator := bytes.IndexByte(data, 0)

	if separator == -1 {
		return fmt.Errorf("invalid object format: missing separator")
	}

	header := string(data[:separator])
	content := data[separator+1:]

	parts := strings.SplitN(header, " ", 2)

	if len(parts) != 2 {
		return fmt.Errorf("invalid object header")
	}

	objectType := parts[0]
	sizeText := parts[1]

	if objectType != "blob" {
		return fmt.Errorf(
			"unsupported object type: %s",
			objectType,
		)
	}

	expectedSize, err := strconv.Atoi(sizeText)
	if err != nil {
		return fmt.Errorf(
			"invalid object size: %s",
			sizeText,
		)
	}

	actualSize := len(content)

	if expectedSize != actualSize {
		return fmt.Errorf(
			"object size mismatch: expected %d, got %d",
			expectedSize,
			actualSize,
		)
	}

	fmt.Print(string(content))

	return nil
}