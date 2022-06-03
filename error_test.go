package lightinggo

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	// "errors"
)

var books = []string{
	"Book1",
	"Book222222",
	"Book3333333333",
}

func TestText(t *testing.T) {
	// books := []string{
	// 	"Book1",
	// 	"Book222222",
	// 	"Book3333333333",
	// }

	for _, bookName := range books {
		err := searchBook(bookName)

		// 特殊业务场景：如果发现书被借走了，下次再来就行了，不需要作为错误处理
		if err != nil {
			// 提取error这个interface底层的错误码，一般在API的返回前才提取
			// As - 获取错误的具体实现
			var customErr = new(CustomError)
			// As - 解析错误内容
			if errors.As(err, &customErr) {
				fmt.Printf("AS中的信息: 当前书为: %s, error code is %d, message is %s\n", bookName, customErr.Code, customErr.Message)
				if customErr.Code == ExpiredError {
					fmt.Printf("IS中的info信息: %s 已经失效了, 只需按Info处理!\n", bookName)
					err = nil
				} else if customErr.Code == NotFoundError {
					// 如果已有堆栈信息，应调用WithMessage方法
					newErr := errors.WithMessage(err, "WithMessage err1")
					// 使用%+v可以打印完整的堆栈信息
					fmt.Printf("IS中的error信息: %s 未找到, 应该按Error处理! ,newErr is: %+v\n", bookName, newErr)
				}
			}
		}
	}
}

func searchBook(bookName string) error {
	// 1 发现图书馆不存在这本书 - 认为是错误，需要打印详细的错误信息
	if len(bookName) > 10 {
		return InitError(NotFoundError)
	} else if len(bookName) > 6 {
		// 2 发现书被借走了 - 打印一下被接走的提示即可，不认为是错误
		return InitError(ExpiredError)
	}
	// 3 找到书 - 不需要任何处理
	return nil
}

func BenchmarkText(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = searchBook(books[i])
			}
		})
	}
}
