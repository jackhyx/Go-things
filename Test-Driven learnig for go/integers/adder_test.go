package integers

import (
	"fmt"
	"testing"
)

func TestAdder(t *testing.T) {
	sum := Add(2, 2)
	expected := 4

	if sum != expected {
		t.Errorf("expected '%d' but got '%d'", expected, sum) // %d 整数 %s 字符串
	}

}

func Add(x, y int) int {
	return x + y

}

func ExampleAdd() {
	sum := Add(1, 5)
	fmt.Println(sum)
	// Output: 6
}

// 我们在学习了# 命名返回值 #但没有在这里使用。它通常应该在结果的含义在上下文中不明确时使用，在我们的例子中，Add 函数将参数相加是非常明显的。

/*
请注意，如果删除注释 「//Output: 6」，示例函数将不会执行。虽然函数会被编译，但是它不会执行。
通过添加这段代码，示例将出现在 godoc 的文档中，这将使你的代码更容易理解。
为了验证这一点，运行 godoc -http=:6060 并访问 http://localhost:6060/pkg/。在这里你能看到 $GOPATH 下所有包的列表，假如你是在 $GOPATH/src/github.com/{your_id} 下编写的这些代码，你就能在文档中找到它
*/
