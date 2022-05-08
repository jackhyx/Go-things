package iteration

import "testing"

func TestRepeat(t *testing.T) {
	repeated := Repeat("a")
	expected := "aaaaa"

	if repeated != expected {
		t.Errorf("expected '%q' but got '%q'", expected, repeated)
	}
}

func Repeat(character string) string {
	var repeated string
	for i := 0; i < 5; i++ {
		repeated += character
	}
	return repeated
}

// 基准测试
func BenchmarkRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Repeat("a")
	}
}

/*
总结
更多的 TDD 练习
学习了 for 循环
学习了如何编写基准测试
*/

/*
我们目前都是使用 := 来声明和初始化变量
然后 := 只是两个步骤的简写
这里我们使用显式的版本来声明一个 string 类型的变量
我们还可以使用 var 来声明函数，稍后我们将看到这一点
*/
