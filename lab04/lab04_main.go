package main

import (
	"bufio"
	"fmt"
	"lab04/internal"
	"log"
	"os"
)

func main() {
	pq := internal.New()

	// Đọc từ console
	scanner := bufio.NewScanner(os.Stdin)

	// Thêm câu hỏi
	pq.AddQuestion("Gia đình bạn có mấy người?")

	// Chạy qua các câu hỏi và in ra màn hình, người dùng sẽ nhập câu trả lời để lưu vào map
	for question := range pq.Qa {
		fmt.Println(question)
		scanner.Scan()

		// Lấy câu trả lời từ console
		answer := scanner.Text()
		// Kiểm tra lỗi
		err := pq.AddAnswer(question, answer)
		if err != nil {
			log.Fatal("Có lỗi khi thêm câu trả lời:", err)
		}
	}

	// Lưu các câu hỏi và câu trả lời vào file
	err := pq.WriteToFile()
	if err != nil {
		log.Fatal("Có lỗi khi lưu dữ liệu:", err)
	}
	fmt.Println("Lưu dữ liệu thành công vào file person.txt")
}
