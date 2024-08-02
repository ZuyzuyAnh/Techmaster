package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

// Struct CLI quản lý các giá trị và phương thức liên quan đến giao diện command line
type CLI struct {
	dirPath string
}

// Khởi tạo CLI với giá trị mặc định
func NewCLI() *CLI {
	return &CLI{}
}

// Entry yêu cầu người dùng nhập đường dẫn thư mục
func (cli *CLI) Entry() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("Usage: tree <directory_path>")
	}
	cli.dirPath = os.Args[1]

	// Kiểm tra xem thư mục có tồn tại không
	info, err := os.Stat(cli.dirPath)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("error: path is not a directory")
	}

	return nil
}

// printTree in thư mục và tệp trong cấu trúc cây
func (cli *CLI) printTree(root string, level int, prefix string) error {
	// Từ thư mục gốc, lưu entry của tất cả các thư mục con trong nó lại, DirEntry là interface chứa các thông tin về thư mục, file
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	//Duyệt qua các entry
	for i, entry := range entries {
		//Kiểm tra xem liệu directory hiện tại đã là cuối cùng chưa
		isLast := i == len(entries)-1
		prefixLine := prefix
		if isLast {
			prefixLine += "└── "
		} else {
			prefixLine += "├── "
		}

		//In ra tên của thư mục hiện tại
		fmt.Printf("%s%s\n", prefixLine, entry.Name())

		// Nếu entry này là 1 thư mục, thì sẽ thực hiện gọi đệ quy hàm printTree với root là thư mục con hiện tại
		if entry.IsDir() {
			newPrefix := prefix
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			//Đệ quy
			err = cli.printTree(filepath.Join(root, entry.Name()), level+1, newPrefix)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Run thực thi việc hiển thị cấu trúc cây thư mục
func (cli *CLI) Run() {
	fmt.Println(cli.dirPath)
	err := cli.printTree(cli.dirPath, 0, "")
	if err != nil {
		fmt.Println("Error:", err)
	}
}
