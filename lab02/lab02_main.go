package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	sensitiveWords = []string{"sex", "fuck", "drug", "kill"}
	vowels         = "ueoai"
)

func main() {
	readFromFile("sample.txt")
}

// Hàm dể đọc và filter, sửa đổi các từ nhạy cảm.
// Hàm sẽ sử dụng scanner để đọc file theo từng dòng
func readFromFile(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("unable to open file: %v", err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		log.Fatalf("unable to get file stat: %v", err)
	}

	if info.Size() == 0 {
		fmt.Println("This file is empty")
		return
	}

	//Kiểu Scanner hỗ trợ đọc input theo dòng
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		filteredLine := replaceSensitiveWords(line, sensitiveWords)
		fmt.Println(filteredLine)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("error reading file: %v", err)
	}
}

// Hàm thay thế nguyên âm trong các từ nhạy cảm
func replaceSensitiveWords(line string, sensitiveWords []string) string {
	//tạo []string words bằng cách chia string thành các từ theo khoảng cách
	words := strings.Fields(line)

	for i, word := range words {
		//Nếu như có xuất hiện từ nhạy cảm trong word, tiến hành thay thế nguyên âm của từ đó thành *
		if containsSensitiveWord(word, sensitiveWords) {
			words[i] = replaceVowels(word)
		}
	}

	//Nối các từ trong words lại thành string như ban đầu
	return strings.Join(words, " ")
}

// Kiểm tra xem từ đã cho có phải là từ nhạy cảm không
func containsSensitiveWord(word string, sensitiveWords []string) bool {
	//Để tránh việc bỏ qua trường hợp SEX và sex, chuyển input thành viết thường
	lowerWord := strings.ToLower(word)

	//Kiểm tra xem liệu nó có trùng bất kì từ nào trong danh sách các từ nhạy cảm không, nếu có trả về true, nếu không trả về false
	for _, sensitiveWord := range sensitiveWords {
		if lowerWord == sensitiveWord {
			return true
		}
	}

	return false
}

// Thay thế nguyên âm bằng ký tự '*'
func replaceVowels(word string) string {
	//String builder để xây dựng string mới, do string trong Golang là immutable
	result := strings.Builder{}

	//Nếu xuất hiện nguyên âm thì thêm vào *, còn không thì giữ nguyên
	for _, ch := range word {
		if strings.ContainsRune(vowels, ch) {
			result.WriteRune('*')
		} else {
			result.WriteRune(ch)
		}
	}

	return result.String()
}
