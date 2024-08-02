- Code để chạy cli: 
    ./myprogram.exe <file_name> <line_to_read>
    Ví dụ: ./myprogram.exe sample.txt 5

- Khác với việc ghi file, việc đọc file có dung lượng lớn sẽ không ảnh hưởng đến bộ nhớ trong trường hợp file có kích thước cực lớn.
Vậy nên theo em chỉ cần quan tâm đến vấn đề cải thiện hiệu suất đọc file. việc đọc file tuần tự theo dòng sẽ tốn nhiều thời gian

- Ý tưởng của em như sau: kết hợp đọc file theo chunk và thuật toán binary search. 

    + Trước hết, chúng ta chia file thành các chunk có kích thước cố định sử dụng offset.
    + Thực hiện việc đọc các chunk trong các goroutine, mỗi goroutine sẽ có 1 biến để lưu số dòng đọc được.
    + Số dòng của 1 chunk sẽ được goroutine chuyển vào trong channel.

    + Sau khi đã có tổng số dòng và danh sách các offset rồi thì thực hiện thuật toán binary search để tìm kiếm dòng thứ Y theo đề bài.
    + Do chúng ta đã chia file thành các chunk có kích thước nhỏ hơn, tức là dòng thứ Y sẽ nằm đầu đó trong số các chunk này
    + Chúng ta sẽ bắt đầu từ chunk ở giữa, tính tổng số dòng từ chunk trong [0, mid) ( nửa khoảng) bằng cách lấy giá trị từ channel.

    + So sánh để thu gọn phạm vi tìm kiếm :
        * Trường hợp tổng số dòng từ chunk 0 đến chunk mid - 1 bé hơn Y và tổng số dòng từ chunk 0 đến chunk mid lớn hơn bằng Y => Y nằm trong chunk mid
        * Trường hợp tổng số dòng từ chunk 0 đến chunk mid - 1 lớn hơn Y: cận trên đổi thành mid - 1
        * Còn lại thì chuyển cận dưới thành mid + 1

