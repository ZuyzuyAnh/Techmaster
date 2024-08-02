package cli

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Struct CLI quản lý các giá trị và phương thức liên quan đến giao diện command line
type CLI struct {
	dirPath string
	keyword string
}

/*
Struct SearchResult lưu trữ kết quả tìm kiếm, vì các giá trị khi đọc file sẽ được gửi vào channel,
channel cần biết chính xác cấu trúc dữ liệu để lưu trữ các giá trị này. Tránh việc in ra các kết quả tìm kiếm lộn xộn
*/
type SearchResult struct {
	FileName   string
	LineNumber int
	Content    string
}

// Khởi tạo CLI với giá trị mặc định
func NewCLI() *CLI {
	return &CLI{}
}

// Entry yêu cầu người dùng nhập đường dẫn thư mục và từ khóa từ command line
func (cli *CLI) Entry() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: search <directory_path> <keyword>")
	}

	cli.dirPath = os.Args[1]
	cli.keyword = os.Args[2]

	if cli.dirPath == "" || cli.keyword == "" {
		return fmt.Errorf("both directory path and keyword are required")
	}

	return nil
}

// searchInFile tìm từ khóa trong một file và gửi kết quả qua channel
func searchInFile(filePath string, keyword string, wg *sync.WaitGroup, resultChan chan<- SearchResult) {
	defer wg.Done()
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Sử dụng biến đếm để lưu dòng mà chứa từ khóa cần tìm kiếm
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		lineContent := scanner.Text()
		if strings.Contains(lineContent, keyword) {

			resultChan <- SearchResult{
				FileName:   filePath,
				LineNumber: lineNumber,
				Content:    lineContent,
			}
		}
	}
	if scanner.Err() != nil {
		fmt.Println("Error reading file:", scanner.Err())
	}
}

// searchInDirectory tìm từ khóa trong tất cả các file .txt trong thư mục
func (cli *CLI) SearchInDirectory() {
	var wg sync.WaitGroup
	resultChan := make(chan SearchResult)

	// Tìm tất cả các file .txt trong thư mục
	// Hàm Walk sẽ đi qua tất cả các file trong thư mục và thực hiện 1 hàm call back trong argument lên các file dó
	err := filepath.WalkDir(cli.dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".txt") {
			wg.Add(1)
			// Tạo goroutine đọc các file
			go searchInFile(path, cli.keyword, &wg, resultChan)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking the directory:", err)
		return
	}

	// Đóng channel khi tất cả goroutines hoàn thành
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Sử dụng map để lưu trữ kết quả tìm kiếm theo tên file, in ra theo đúng yêu cầu của đề bài
	results := make(map[string][]SearchResult)
	for result := range resultChan {
		results[result.FileName] = append(results[result.FileName], result)
	}

	// In kết quả tìm kiếm theo định dạng yêu cầu
	for fileName, searchResults := range results {
		fmt.Println(fileName)
		for _, result := range searchResults {
			fmt.Printf("\t%d: %s\n", result.LineNumber, result.Content)
		}
	}
}
