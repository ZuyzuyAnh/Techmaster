package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	content := `Hello Word
My name is Duy Anh
I want to become a software developer
`

	writeToFile(content, "sample.txt")
}

// Hàm ghi nội dung từ string vào file
func writeToFile(content, fileName string) {

	//Thực hiện mở file ở chế độ append và chỉ ghi, nếu không tồn tại file thì tạo file mới vaf
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error creating or opening file: %v", err)
	}
	defer file.Close()

	/*Hàm WriteString trả về kiểu int và error. Giá trị int là số lượng byte đã được viết vào file.
	Do không cần thiết nên sẽ sử dụng blank */
	_, err = file.WriteString(content)
	if err != nil {
		log.Fatalf("Error writing to file:: %v", err)
	}

	fmt.Println("Content written to file successfully")
}
