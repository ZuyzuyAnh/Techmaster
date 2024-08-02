package cli

import (
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
)

const (
	lineLength = 256
	chars      = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

// struct CLI quản lý các giá trị và phương thức liên quan đến giao diện command line
type CLI struct {
	fileName string
	size     int
}

// Hàm New để khởi tạo CLI
func New() *CLI {
	return &CLI{}
}

// Bắt đầu chương trình cmd, yêu cầu người dùng nhập tên file và kích thước file từ command line arguments
func (cli *CLI) Entry() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("Usage: %s <filename> <size_in_MB>", os.Args[0])
	}

	cli.fileName = os.Args[1]

	inputSize, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return fmt.Errorf("please enter a valid integer for file size")
	}

	if inputSize <= 10 {
		return fmt.Errorf("please enter a positive integer greater than 10 for file size")
	}

	cli.size = inputSize
	return nil
}

// Tạo file với kích cỡ tương ứng với size trong CLI
func (cli *CLI) CreateFile() error {
	file, err := os.Create(cli.fileName)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// Vì 1 MB = 1024*1024 bytes, nên tổng số byte = size*1024*1024
	totalBytes := cli.size * 1024 * 1024
	// 1 dòng có 256 kí tự => Chiếm 256 bytes, số dòng = tổng số bytes / số bytes mỗi dòng
	numLines := totalBytes / lineLength
	// Phần dư để làm đủ thành X MB theo đề bài yêu cầu
	remainder := totalBytes % lineLength

	// Lặp lại số hàng cần thiết để đủ X MB, mỗi hàng generate ra các string có độ dài 256 kí tự và viết thành 1 dòng ở trong file
	for i := 0; i < numLines; i++ {
		err = writeToFile(file, lineLength)
		if err != nil {
			return err
		}
	}

	// Với số lượng byte còn lại, tạo 1 string với kích thước tương ứng, sau đó thêm vào dòng cuối cùng
	if remainder > 0 {
		err = writeToFile(file, remainder)
		if err != nil {
			return err
		}
	}

	return nil
}

// Hàm helper giup tạo và ghi string vào file
func writeToFile(file *os.File, lineLength int) error {
	line, err := randomString(lineLength)
	if err != nil {
		return fmt.Errorf("error generating random string: %w", err)
	}

	_, err = file.WriteString(line + "\n")
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

// Hàm helper để tạo ngẫu nhiên 1 string có length kí tự (cụ thể ở đề bài là 256). Tuy nhiên sẽ cần giá trị length đầu vào để sử dụng cho việc gen ra số lượng byte còn thiếu
func randomString(length int) (string, error) {
	// Khởi tạo slice bytes có length phần tử
	bytes := make([]byte, length)

	// Fill slice với các giá trị byte ngẫu nhiên
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Do khi sinh ngẫu nhiên các byte trong bytes, sẽ dễ dẫn đến việc các byte đó có thể nhận giá trị lớn hơn index của bytes, cụ thể là 256.
	// Phép chia dư đảm bảo giá trị sẽ nằm trong khoảng 0 -> 256
	for i := 0; i < length; i++ {
		bytes[i] = chars[int(bytes[i])%len(chars)]
	}
	return string(bytes), nil
}
