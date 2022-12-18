package torips

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// This file and its functions are coded by ChatGPT.

func RefreshList() (size int64) {
	// Set the GitHub repository URL and file path
	repoURL := "https://github.com/SecOps-Institute/Tor-IP-Addresses"
	fileName := "tor-nodes.lst"

	// Clone the repository to a temporary directory
	tempDir, err := cloneRepo(repoURL)
	if err != nil {
		fmt.Printf("Error cloning repository: %s\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir) // Clean up the temporary directory when we're done

	// Open the file for reading
	file, err := os.Open(filepath.Join(tempDir, fileName))
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close() // Close the file when we're done

	// Copy the file to standard output
	f, err := os.Create(fileName)
	if err != nil {
		panic(err.Error())
	}
	size, err = io.Copy(f, file)
	if err != nil {
		fmt.Printf("Error copying file: %s\n", err)
		os.Exit(1)
	}
	return
}

func cloneRepo(repoURL string) (string, error) {
	// Create a temporary directory to store the repository
	tempDir, err := filepath.Abs(os.TempDir())
	if err != nil {
		return "", err
	}
	tempDir = filepath.Join(tempDir, "repo")

	// Clone the repository
	cmd := exec.Command("git", "clone", repoURL, tempDir)
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	return tempDir, nil
}
