package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Locations []Location `json:"locations"`
	LogFile   string     `json:"log_file"`
}

type Location struct {
	SourceDir string `json:"path"`
	Days      int    `json:"days"`
	Action    string `json:"action"`
	TargetDir string `json:"target,omitempty"`
}

func main() {
	configFile := flag.String("c", "config.json", "Path to configuration file")
	flag.Parse()

	config, err := loadConfig(*configFile)
	if err != nil {
		slog.Error("Error loading configuration:", err)
		os.Exit(1)
	}

	logFile, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Failed to open log file %s", err)
		os.Exit(1)
	}

	defer logFile.Close()

	log.SetOutput(logFile)

	for _, location := range config.Locations {
		processFiles(location)
	}
}

func loadConfig(filename string) (Config, error) {
	var config Config
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	return config, err
}

func processFiles(location Location) {
	slog.Info(fmt.Sprintf("Processing location: %s", location.SourceDir))
	fileTree := os.DirFS(location.SourceDir)

	fs.WalkDir(fileTree, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			slog.Error("Some error occurred", err)
			os.Exit(1)
		}

		if !d.IsDir() && shouldClean(d, location.Days) {
			sourceFile := filepath.Join(location.SourceDir, path)
			if location.Action == "delete" {
				deleteFiles(sourceFile)
			} else if location.Action == "move" {
				targetFile := filepath.Join(location.TargetDir, path)
				moveFiles(sourceFile, targetFile)
			}
		}
		return nil
	})
}

func shouldClean(file fs.DirEntry, days int) bool {
	fileInfo, err := file.Info()
	if err != nil {
		slog.Error(fmt.Sprintf("Could not retrieve file info for: %s", file.Name()), err)
		return false
	}

	now := time.Now()
	modTime := fileInfo.ModTime()
	diff := now.Sub(modTime).Hours() / 24
	return int(diff) > days
}

func deleteFiles(fullPath string) {
	err := os.Remove(fullPath)
	if err != nil {
		slog.Error(fmt.Sprintf("Error deleting file: %s", fullPath), err)
	} else {
		slog.Info(fmt.Sprintf("Deleted file: %s", fullPath))
	}
}

func moveFiles(sourceFile string, targetFile string) {

	err := os.MkdirAll(filepath.Dir(targetFile), 0755)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating target directory: %s", filepath.Dir(targetFile)), err)
	}

	err = os.Rename(sourceFile, targetFile)
	if err != nil {
		slog.Error(fmt.Sprintf("Error moving file: %s", sourceFile), err)
	} else {
		slog.Info(fmt.Sprintf("Moved file from %s to %s", sourceFile, targetFile))
	}
}
