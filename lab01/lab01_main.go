package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	readFromFile("sample.txt", 1024)
}

// Tạo 1 hàm để đọc file theo tùy chọn: tên file và kích cỡ của buffer.
func readFromFile(fileName string, bufferSize int) {

	/*hàm Open mở file cho việc đọc, giá trị trả về sẽ là 1 con trỏ kiểu File và giá trị error,
	thực hiện kiểm tra xem liệu biến err có nhận giá trị nil, nếu có sẽ in ra console và kết thúc hàm */
	f, err := os.Open(fileName)
	if err != nil {
		log.Printf("unable to read file: %v", err)
		return
	}
	//Thực hiện dóng File f sau khi kết thúc hàm readFromFile.
	defer f.Close()

	//2 Block dưới đây sẽ thực hiện việc truy xuất thông tin về file và kiểm tra xem file có rỗng không
	info, err := f.Stat()
	if err != nil {
		log.Fatalf("unable to get file stat: %v", err)
	}

	// Nếu file rỗng, in ra thông báo và dừng hàm
	if info.Size() == 0 {
		fmt.Println("this file is empty")
		return
	}

	//Tạo buffer có kích thước là bufferSize.
	buffer := make([]byte, bufferSize)
	for {
		n, err := f.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n > 0 {
			fmt.Print(string(buffer[:n]))
		}
	}
}
