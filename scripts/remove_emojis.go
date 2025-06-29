package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Configurable, hardcoded path: Windows style
const scanPath = "C:\\Users\\adilm\\repositories\\General\\blueops"

// Stats holds scan and clean statistics
type Stats struct {
	TotalFilesScanned  int
	FilesWithEmojis    int
	FilesCleaned       int
	TotalEmojisRemoved int
}

// EmojiRemover handles emoji removal from files
type EmojiRemover struct {
	Verbose   bool
	FileTypes []string
	Stats     Stats
}

// NewEmojiRemover creates a new emoji remover instance
func NewEmojiRemover(verbose bool, fileTypes []string) *EmojiRemover {
	return &EmojiRemover{
		Verbose:   verbose,
		FileTypes: fileTypes,
	}
}

// isEmoji checks if a rune is an emoji (covering most Unicode emoji ranges)
func (er *EmojiRemover) isEmoji(r rune) bool {
	return (r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols and Pictographs
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map
		(r >= 0x1F1E6 && r <= 0x1F1FF) || // Regional indicators
		(r >= 0x2600 && r <= 0x26FF) || // Misc symbols
		(r >= 0x2700 && r <= 0x27BF) || // Dingbats
		(r >= 0x1F900 && r <= 0x1F9FF) || // Supplemental Symbols and Pictographs
		(r >= 0x1FA70 && r <= 0x1FAFF) || // Extended
		(r >= 0x1F018 && r <= 0x1F270) || // Various symbols
		(r >= 0x238C && r <= 0x2454) || // Misc technical
		(r >= 0x20D0 && r <= 0x20FF) // Combining marks
}

// removeEmojisFromText removes emojis from text, returning cleaned string and number of emojis removed
func (er *EmojiRemover) removeEmojisFromText(text string) (string, int) {
	var result strings.Builder
	removed := 0
	for _, r := range text {
		if !er.isEmoji(r) {
			result.WriteRune(r)
		} else {
			removed++
		}
	}
	return result.String(), removed
}

// shouldProcessFile checks if file should be processed based on extension
func (er *EmojiRemover) shouldProcessFile(filename string) bool {
	if len(er.FileTypes) == 0 {
		return true // Process all files if no specific types specified
	}
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
	for _, fileType := range er.FileTypes {
		if ext == strings.ToLower(strings.TrimPrefix(fileType, ".")) {
			return true
		}
	}
	return false
}

// processFile removes emojis from a single file, updates stats
func (er *EmojiRemover) processFile(filePath string) error {
	if !er.shouldProcessFile(filePath) {
		return nil
	}

	er.Stats.TotalFilesScanned++

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	originalText := string(content)
	cleanedText, emojisRemoved := er.removeEmojisFromText(originalText)

	if emojisRemoved > 0 {
		er.Stats.FilesWithEmojis++
		er.Stats.TotalEmojisRemoved += emojisRemoved

		if originalText != cleanedText {
			er.Stats.FilesCleaned++
			err = os.WriteFile(filePath, []byte(cleanedText), 0644)
			if err != nil {
				return fmt.Errorf("failed to write cleaned file %s: %w", filePath, err)
			}
			fmt.Printf("CLEANED: %s (removed %d emoji(s))\n", filePath, emojisRemoved)
		}
	} else if er.Verbose {
		fmt.Printf("No emojis found in: %s\n", filePath)
	}
	return nil
}

// processDirectory recursively processes all files in a directory
func (er *EmojiRemover) processDirectory(dirPath string) error {
	return filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip hidden directories and common ignore patterns
		if d.IsDir() {
			name := d.Name()
			if strings.HasPrefix(name, ".") ||
				name == "node_modules" ||
				name == "__pycache__" ||
				name == "vendor" ||
				name == "target" {
				return filepath.SkipDir
			}
			return nil
		}
		return er.processFile(path)
	})
}

func main() {
	fmt.Println("BlueOps Emoji Remover")
	fmt.Printf("Hardcoded scan path: %s\n", scanPath)
	fmt.Println(strings.Repeat("-", 50))

	startTime := time.Now()

	// You can configure allowed file types here, or leave empty to process all types
	fileTypes := []string{} // e.g. []string{"txt", "md"}
	verbose := false

	remover := NewEmojiRemover(verbose, fileTypes)
	err := remover.processDirectory(scanPath)
	if err != nil {
		log.Fatalf("Error during scan: %v", err)
	}

	elapsed := time.Since(startTime)
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("Scan complete!\n")
	fmt.Printf("Total files scanned:        %d\n", remover.Stats.TotalFilesScanned)
	fmt.Printf("Files with emojis found:    %d\n", remover.Stats.FilesWithEmojis)
	fmt.Printf("Files cleaned:              %d\n", remover.Stats.FilesCleaned)
	fmt.Printf("Total emojis removed:       %d\n", remover.Stats.TotalEmojisRemoved)
	fmt.Printf("Time taken:                 %s\n", elapsed.Round(time.Millisecond))
}
