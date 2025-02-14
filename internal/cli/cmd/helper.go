package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"gitlab.com/zigal0/architect/internal/cli/logger"
)

// Errors
const (
	formatErrFileAction = "failed to %s file %s: %w"

	statFileAction   = "get info for"
	createFileAction = "create"
	openFileAction   = "open"
	writeFileAction  = "write to"
	scanFileAction   = "scan"
)

const (
	goModFileName = "go.mod"
)

type projectPartInfo struct {
	absPath        string
	pathParts      []string
	tmplt          string
	tmpltData      any
	needToOverride bool
}

func checkFileExist(path string) bool {
	if _, err := os.Stat(filepath.Clean(path)); errors.Is(err, fs.ErrNotExist) {
		return false
	}

	return true
}

func createProjectPart(info projectPartInfo) {
	if len(info.pathParts) == 0 {
		logger.Fatal("Incorrect info to create poject part")
	}

	pathParts := append([]string{info.absPath}, info.pathParts...)

	logger.Infof("Creating %s...", pathParts[len(pathParts)-1])

	filePath := filepath.Join(pathParts...)

	if !info.needToOverride && checkFileExist(filePath) {
		return
	}

	content, err := createContentFromTemplate(info.tmplt, info.tmpltData)
	logger.FatalIfErr(err)

	logger.FatalIfErr(writeStringToFile(filePath, content))
}

func appendToProjectPart(info projectPartInfo) {
	if len(info.pathParts) == 0 {
		logger.Fatal("Incorrect info to append to poject part")
	}

	pathParts := append([]string{info.absPath}, info.pathParts...)

	logger.Infof("Appending %s...", pathParts[len(pathParts)-1])

	filePath := filepath.Join(pathParts...)

	if !checkFileExist(filePath) {
		return
	}

	content, err := createContentFromTemplate(info.tmplt, info.tmpltData)
	logger.FatalIfErr(err)

	logger.FatalIfErr(appendStringToFile(filePath, content))
}

func createContentFromTemplate(templateSrc string, data any) (string, error) {
	tmpl, err := template.New("").Parse(templateSrc)
	if err != nil {
		return "", fmt.Errorf("failed to parse source template: %w", err)
	}

	buf := bytes.Buffer{}

	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute teplate: %w", err)
	}

	return buf.String(), nil
}

func writeStringToFile(rawFilePath, content string) error {
	cleanPath := filepath.Clean(rawFilePath)

	if dir := filepath.Dir(cleanPath); dir != "" {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to make dir %s: %w", dir, err)
		}
	}

	file, err := os.Create(cleanPath)
	if err != nil {
		return fmt.Errorf(formatErrFileAction, createFileAction, cleanPath, err)
	}

	defer func() { _ = file.Close() }()

	_, err = file.Write([]byte(content))
	if err != nil {
		return fmt.Errorf(formatErrFileAction, writeFileAction, cleanPath, err)
	}

	return nil
}

func appendStringToFile(rawFilePath, content string) error {
	cleanPath := filepath.Clean(rawFilePath)

	file, err := os.OpenFile(cleanPath, os.O_APPEND|os.O_WRONLY, 0600) // nolint: gomnd
	if err != nil {
		return fmt.Errorf(formatErrFileAction, openFileAction, cleanPath, err)
	}

	defer func() { _ = file.Close() }()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf(formatErrFileAction, writeFileAction, cleanPath, err)
	}

	return nil
}

func moduleFromGoMod() (string, error) {
	const (
		modulePrefix = "module "
	)

	if _, err := os.Stat(goModFileName); err != nil {
		return "", fmt.Errorf(formatErrFileAction, statFileAction, goModFileName, err)
	}

	goMod, err := os.Open(goModFileName)
	if err != nil {
		return "", fmt.Errorf(formatErrFileAction, openFileAction, goModFileName, err)
	}

	defer func() { _ = goMod.Close() }()

	scanner := bufio.NewScanner(goMod)
	scanner.Split(bufio.ScanLines)

	var module string

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), modulePrefix) {
			module = strings.TrimSpace(strings.TrimPrefix(scanner.Text(), modulePrefix))

			break
		}
	}

	if err = scanner.Err(); err != nil {
		return "", fmt.Errorf(formatErrFileAction, scanFileAction, goModFileName, err)
	}

	return module, nil
}

func execute(commandName string, args ...string) error {
	logger.Info(fmt.Sprintf("Executing command '%s' with args: %q...", commandName, args))

	cmd := exec.Command(commandName, args...)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute command '%s' with args: %q: %w", commandName, args, err)
	}

	return nil
}

func executeMake(target, path string) error {
	return execute("make", "-C", path, target)
}

func executeGoModTidy() error {
	return execute("go", "mod", "tidy")
}
