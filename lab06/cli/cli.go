package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"sync"
)

// Struct CLI quản lý các giá trị và phương thức liên quan đến giao diện command line
type CLI struct {
	fileName string
	line     int
}

// Khởi tạo CLI với giá trị mặc định
func NewCLI() *CLI {
	return &CLI{}
}

// Entry yêu cầu người dùng nhập tên file và số dòng từ dòng lệnh
func (cli *CLI) Entry() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: search <filename> <line_number>")
	}

	cli.fileName = os.Args[1]

	var err error
	cli.line, err = strconv.Atoi(os.Args[2])
	if err != nil || cli.line <= 0 {
		return fmt.Errorf("please enter a valid positive integer for line number")
	}

	return nil
}

// countLinesInChunk đếm số dòng trong một chunk
func countLinesInChunk(filename string, offset, size int, wg *sync.WaitGroup, resultChan chan<- int) {
	defer wg.Done()
	file, err := os.Open(filename)
	if err != nil {
		resultChan <- 0
		return
	}
	defer file.Close()

	// Tạo slice byte với kích thước đủ lớn sau đó đọc dữ liệu file bắt đầu từ offset vào buffer
	buffer := make([]byte, size)
	_, err = file.ReadAt(buffer, int64(offset))
	if err != nil {
		resultChan <- 0
		return
	}

	// Đếm số dòng của chunk
	scanner := bufio.NewScanner(bytes.NewReader(buffer))
	lineCount := 0

	for scanner.Scan() {
		lineCount++
	}
	if scanner.Err() != nil {
		resultChan <- 0
		return
	}

	resultChan <- lineCount
}

// ReadLine đọc dòng thứ Y từ file bằng cách sử dụng binary search và chunks
func (cli *CLI) ReadLine() (string, error) {
	// Mở file
	file, err := os.Open(cli.fileName)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Đọc kích thước của file
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("error getting file stats: %w", err)
	}
	fileSize := fileInfo.Size()

	// Tạo danh sách các chunk offsets
	chunkSize := 1024 // Kích thước khởi đầu của chunk (1 KB)
	var chunkOffsets []int64
	for offset := int64(0); offset < fileSize; offset += int64(chunkSize) {
		chunkOffsets = append(chunkOffsets, offset)
		// if chunkSize < 1024*1024 {
		// 	chunkSize *= 2 // Tăng kích thước chunk theo dạng số mũ
		// }
	}

	// Đếm số dòng trong từng chunk song song
	var wg sync.WaitGroup
	resultChan := make(chan int, len(chunkOffsets))
	for _, offset := range chunkOffsets {
		wg.Add(1)
		go countLinesInChunk(cli.fileName, int(offset), chunkSize, &wg, resultChan)
	}
	// Đợi đến khi WaitGroup không còn goroutine nào chưa hoàn thành
	wg.Wait()

	// Đóng channel lại
	close(resultChan)

	// Tổng hợp số dòng trong từng chunk
	var totalLines int
	lineCounts := make([]int, len(chunkOffsets))
	i := 0
	for count := range resultChan {
		lineCounts[i] = count
		totalLines += count
		i++
	}

	// Sử dụng binary search để tìm chunk chứa dòng Y
	lower := 0
	upper := len(chunkOffsets) - 1
	var targetChunk int

	for lower <= upper {
		mid := (lower + upper) / 2
		linesBeforeCurrentChunk := 0
		if mid > 0 {
			for i := 0; i < mid; i++ {
				linesBeforeCurrentChunk += lineCounts[i]
			}
		}

		linesInChunk := lineCounts[mid]
		if linesBeforeCurrentChunk < cli.line && linesBeforeCurrentChunk+linesInChunk >= cli.line {
			targetChunk = mid
			break
		} else if linesBeforeCurrentChunk+linesInChunk < cli.line {
			lower = mid + 1
		} else {
			upper = mid - 1
		}
	}

	if targetChunk < 0 || targetChunk >= len(chunkOffsets) {
		return "", fmt.Errorf("line number %d not found", cli.line)
	}

	// Đọc dữ liệu từ chunk chứa dòng Y
	targetOffset := chunkOffsets[targetChunk]
	file.Seek(targetOffset, 0)
	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
		if lineCount == cli.line {
			return scanner.Text(), nil
		}
	}
	if scanner.Err() != nil {
		return "", scanner.Err()
	}

	return "", fmt.Errorf("line number %d not found", cli.line)
}

// Thực thi việc đọc file và in ra error nếu xuất hiện
func (cli *CLI) Run() {
	line, err := cli.ReadLine()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Line %d: %s\n", cli.line, line)
}
