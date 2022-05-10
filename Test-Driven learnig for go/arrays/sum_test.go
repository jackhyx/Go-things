package main

import (
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {

	t.Run("collectuib if 5 numbers", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}

		got := Sum(numbers)
		want := 15

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

	t.Run("colletion of any size", func(t *testing.T) {
		numbers := []int{1, 2, 3}
		got := Sum(numbers)
		want := 6

		if got != want {
			t.Errorf("got %d want givem. %v", got, want)
		}
	})

}
func Sum(numbers []int) int {
	sum := 0
	for _, number := range numbers {
		sum += number //ange 会迭代数组，每次迭代都会返回数组元素的索引和值。我们选择使用 _ 空白标志符 来忽略索引
	}

	return sum
}

func TestSumAll(t *testing.T) {

	got := SumAll([]int{1, 2}, []int{0, 9})
	want := []int{3, 9}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
func SumAll(numbersToSum ...[]int) (sums []int) {
	length0fNumbers := len(numbersToSum)
	sums = make([]int, length0fNumbers)

	for i, numbers := range numbersToSum {
		sums[i] = Sum(numbers)
	}
	return
}

/*
切片有容量的概念。如果你有一个容量为 2 的切片，但使用 mySlice[10]=1 进行赋值，会报运行时错误。
不过你可以使用 append 函数，它能为切片追加一个新值
func SumAll(numbersToSum ...[]int) []int {
    var sums []int
    for _, numbers := range numbersToSum {
        sums = append(sums, Sum(numbers))
    }

    return sums
}
*/

// 这里有一种创建切片的新方式。make 可以在创建切片的时候指定我们需要的长度和容量;我们可以使用切片的索引访问切片内的元素，使用 = 对切片元素进行赋值。

/*  判断两个值深度是否一直，即除了值相同，底层类型也相同
func DeepEqual(x, y any) bool {
	if x == nil || y == nil {
		return x == y
	}
	v1 := ValueOf(x)
	v2 := ValueOf(y)
	if v1.Type() != v2.Type() {
		return false
	}
	return deepValueEqual(v1, v2, make(map[visit]bool))
}

*/

// invalid operation: got != want (slice can only be compared to nil)
// 在 Go 中不能对切片使用等号运算符。你可以写一个函数迭代每个元素来检查它们的值。但是一种比较简单的办法是使用 reflect.DeepEqual，它在判断两个变量是否相等时十分有用。
// 需要注意的是 reflect.DeepEqual 不是「类型安全」的，所以有时候会发生比较怪异的行为：比较不同类型的参数

func TestSumAllTails(t *testing.T) {
	/*
		checkSums := func(t *testing.T, got, want []int) {
		        if !reflect.DeepEqual(got, want) {
		            t.Errorf("got %v want %v", got, want)
		        }
		    }
	*/

	t.Run("make the sums of some slices", func(t *testing.T) {
		got := SumAllTails([]int{1, 2}, []int{0, 9})
		want := []int{2, 9}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		} // checkSums(t, got, want)
	})

	t.Run("safely sum empty slices", func(t *testing.T) {
		got := SumAllTails([]int{}, []int{3, 4, 5})
		want := []int{0, 9}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		} // checkSums(t, got, want)
	})

}

func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		if len(numbers) == 0 {
			sums = append(sums, 0)
			// panic: runtime error: slice bounds out of range [1:0] [recovered]
			//panic: runtime error: slice bounds out of range [1:0]
		} else {
			tail := numbers[1:]
			sums = append(sums, Sum(tail))
		}
	}
	return sums
}

// 我们可以使用语法 slice[low:high] 获取部分切片。如果在冒号的一侧没有数字就会一直取到最边缘的元素。在我们的函数中，我们使用 numbers[1:] 取到从索引 1 到最后一个元素。

/*
func TestSumAllTails(t *testing.T) {

	checkSums := func(t *testing.T, got, want []int) {
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	}

	t.Run("make the sums of tails of", func(t *testing.T) {
		got := SumAllTails([]int{1, 2}, []int{0, 9})
		want := []int{2, 9}
		checkSums(t, got, want)
	})

	t.Run("safely sum empty slices", func(t *testing.T) {
		got := SumAllTails([]int{}, []int{3, 4, 5})
		want := []int{0, 9}
		checkSums(t, got, want)
	})

}

func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		if len(numbers) == 0 {
			sums = append(sums, 0)
		} else {
			tail := numbers[1:]
			sums = append(sums, Sum(tail))
		}
	}

	return sums
}

func Sum(numbers []int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}

/*
数组的容量是我们在声明它时指定的固定值。我们可以通过两种方式初始化数组：
[N]type{value1, value2, ..., valueN} e.g. numbers := [5]int{1, 2, 3, 4, 5}
[...]type{value1, value2, ..., valueN} e.g. numbers := [...]int{1, 2, 3, 4, 5}
在错误信息中打印函数的输入有时很有用。我们使用 %v（默认输出格式）占位符来打印输入，它非常适用于展示数组

range 会迭代数组，每次迭代都会返回数组元素的索引和值。
我们选择使用 _ 空白标志符 来忽略索引。
数组有一个有趣的属性，它的大小也属于类型的一部分，如果你尝试将 [4]int 作为 [5]int 类型的参数传入函数，是不能通过编译的。它们是不同的类型，就像尝试将 string 当做 int 类型的参数传入函数一样。
因为这个原因，所以数组比较笨重，大多数情况下我们都不会使用它。
Go 的切片（slice）类型不会将集合的长度保存在类型中，因此它的尺寸可以是不固定的。
下面我们会完成一个动态长度的 Sum 函数。

在 Go 中不能对切片使用等号运算符。你可以写一个函数迭代每个元素来检查它们的值。但是一种比较简单的办法是使用 reflect.DeepEqual，它在判断两个变量是否相等时十分有用。

这里有一种创建切片的新方式。make 可以在创建切片的时候指定我们需要的长度和容量。
我们可以使用切片的索引访问切片内的元素，使用 = 对切片元素进行赋值

如果你有一个容量为 2 的切片，但使用 mySlice[10]=1 进行赋值，会报运行时错误。
不过你可以使用 append 函数，它能为切片追加一个新值。

我们可以使用语法 slice[low:high] 获取部分切片。如果在冒号的一侧没有数字就会一直取到最边缘的元素。在我们的函数中，我们使用 numbers[1:] 取到从索引 1 到最后一个元素
*/
