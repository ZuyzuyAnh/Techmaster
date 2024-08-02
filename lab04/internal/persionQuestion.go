package internal

import (
	"errors"
	"fmt"
	"os"
)

/*
PersonQuestion chứa các câu hỏi và câu trả lời
Map là 1 hash table nên các giá trị sẽ không thể theo thứ tự từ điển hay thứ tự nhập vào.
Thế nhưng chúng ta có thể cài đặt thêm 1 slice lưu các câu hỏi. Thay vì đọc các key ở map thì sẽ đọc các câu hỏi ở slice, từ đó giúp đảm bảo thứ tự của các câu hỏi.
Tuy nhiên do yêu cầu của bài không nêu rõ nên em xin phép chỉ được nêu ý tưởng và không cài đặt.
*/
type PersonQuestion struct {
	Qa map[string]string
}

// Hàm khởi tạo cho struct
func New() *PersonQuestion {
	return &PersonQuestion{
		Qa: map[string]string{
			"Tên bạn là gì?":     "",
			"Bạn sinh ngày nào?": "",
			"Bạn làm nghề gì?":   "",
		},
	}
}

// Hàm thêm câu hỏi
func (p *PersonQuestion) AddQuestion(question string) error {
	if isContain(question, p.Qa) {
		return errors.New("câu hỏi đã tồn tại")
	}
	p.Qa[question] = ""
	return nil
}

// Hàm thêm câu trả lời
func (p *PersonQuestion) AddAnswer(question, answer string) error {
	if !isContain(question, p.Qa) {
		return errors.New("câu hỏi không tồn tại")
	}
	p.Qa[question] = answer
	return nil
}

// Hàm ghi câu hỏi và trả lời vào file person.txt
func (p *PersonQuestion) WriteToFile() error {
	file, err := os.Create("person.txt")
	if err != nil {
		return fmt.Errorf("không thể tạo tệp: %v", err)
	}
	defer file.Close()

	for question, answer := range p.Qa {
		_, err := fmt.Fprintf(file, "%s\n%s\n\n", question, answer)
		if err != nil {
			return fmt.Errorf("lỗi khi ghi dữ liệu vào tệp: %v", err)
		}
	}

	return nil
}

// Hàm kiểm tra sự tồn tại của câu hỏi
func isContain(question string, m map[string]string) bool {
	_, exists := m[question]
	return exists
}
