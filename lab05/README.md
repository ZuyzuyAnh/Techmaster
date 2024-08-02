- Code để chạy cli: 
    ./myapp.exe <file_name> <size>
    Ví dụ: ./myapp.exe sample 20
    Nếu tên file đã tồn tại sẽ ghi đè lên file cũ
- Tương tự như bài số 4, tạo 1 struct quản lý các yêu cầu liên quan đến command line.

- Khảo sát thời gian chạy đối với các kích thước:
    + 11MB: 4.0217044s
    + 50MB: 4.7699275s
    + 100MB: 5.7196274s
    + 200MB: 13.6122764s

    + 1024MB: 26.294093s, 
    + 2048MB: 37.4810056s
    + 5120MB: 1m37.4332013s

    + 10240MB: 2m48.5118117s
    + 20480MB: 5m16.1150027s
    + 51200MB: 14m39.445179s

    + 524288MB: hiện lỗi

- Với kích thước file càng lớn, thời gian để ghi file càng lâu.
- Em đã thử kích thước file lớn với 2 trường hợp:
    1. Kích thước file vượt quá dung lượng ổ cứng (trong trường hợp của em là 525GB)
    2. Kích thước file vượt quá dung lượng hệ điều hành cho phép (đối với window 11 là 16TB)
    
- Kết quả của trường hợp 1 đúng như mong đợi, quá trình ghi file diễn ra rất lâu và chương trình báo lỗi: 
"Error: error writing to file: write sample: There is not enough space on the disk".

- Kết quả của trường hợp 2 thì không được như mong đợi, không có báo lỗi trước từ hệ điều hành và chương trình chạy như trường 
hợp 1. Chạy rất lâu và đến khi ổ đĩa đầy mới có báo lỗi từ hệ điều hành.

=> Như vậy chúng ta có 2 vấn đề chính khi thực hiện với file có kích thước lớn
    + Thời gian thực thi lâu
    + Dễ dẫn đến lỗi từ hệ điều hành vì dung lượng vượt quá mức cho phép

- Các giải pháp:
    1. Thời gian thực thi lâu:

    + Cách em nghĩ đến là việc đọc ghi file song song bằng goroutine, tuy nhiên em chưa chắc chắn về vấn đề
    race conditioning khi 2 goroutine cùng thay đổi 1 file. Sau khi tra cứu bằng chat gpt và google thì giải pháp họ đưa
    ra là sử dụng sync.Mutex
    + Chúng ta sẽ tạo nhiều goroutine cho quá trình generate ra chuỗi 256 kí tự bất kì và ghi nó vào file. Các goroutine khi tạo xong chuỗi sẽ 
    gọi Lock từ mutex và tiến hành ghi file, sau khi ghi file xong thì Unlock.
    + Như vậy, cách ghi trên vẫn chính là ghi tuần tự nhưng sẽ cải thiện thêm chút hiệu suất từ việc các goroutine sinh ra các chuỗi 256 kí tự đồng thời.

    + Từ các hạn chế của cách trên: lưu lượng đọc ghi file nhiều, lock và unlock liên tục dễ dẫn đến overhead, các goroutine chưa thực hiện song song việc ghi vào file
    + Để cải thiện thêm về hiệu suất, chúng ta có thể sử dụng buffer. Ý tưởng chính sẽ là các goroutine sẽ tạo 1 buffer của riêng nó, và thực hiện việc ghi 
    dữ liệu vào buffer. Việc lưu dữ liệu vào buffer của riêng goroutine đó sẽ không bị race conditioning. Sau khi 1 goroutine hoàn thành xong việc ghi vào buffer,
    nó sẽ gọi Lock và tiến hành ghi từ buffer vào file, và unlock khi ghi xong.
    + Lý do hiệu suất được cải thiện:
        * Các goroutine thục hiện ghi đồng thời
        * mỗi goroutine chỉ gọi lock và unlock 1 lần duy nhất => tránh overhead của mutex
        * chia tải đọc ghi ra các goroutine khác nhau
        * mỗi lần ghi dữ liệu trực tiếp vào file sẽ thực hiện system call. Nếu sử dụng Buffer, 5 goroutine sẽ chỉ có 5 system call => giảm tốn kém về mặt hiệu suất và tài nguyên

    2. Lỗi kích thước vượt quá dung lượng ổ đĩa hay hệ điều hành

    + Cách đơn giản nhất là giới hạn kích thước file người dùng có thể nhập, tuy nhiên nó lại chưa đúng với nghiệp vụ của đề bài
    + Chúng ta có thể khai báo hằng số maxSizeOS = 16,777,216MB. Kiểm tra giá trị nhập vào, nếu nó chạm đến giá trị kia thì sẽ in ra lỗi và không tiến hành ghi file
    + Tính toán trước dung lượng còn lại của ổ cứng, và xem liệu kích thước file nhập vào có thỏa mãn không. Nếu không thì in ra lỗi và không tiến hành ghi file.
    + Giới hạn thời gian timeout. Nếu chương trình chạy quá thời gian sẽ trả về lỗi và dừng.

=> Nhìn chung các vấn đề về kích thước file dễ xử lý hơn vấn đề về thời gian chạy.
    