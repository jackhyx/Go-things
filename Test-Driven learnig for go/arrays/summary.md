
## 总结
我们学习了：
### 数组
## 切片
* 多种方式的切片初始化
  * 切片的容量是 固定 的，但是你可以使用 append 从原来的切片中创建一个新切片
```go
func SumAll(numbersToSum ...[]int) []int {
    var sums []int
    for _, numbers := range numbersToSum {
    sums = append(sums, Sum(numbers))
    }

    return sums
} 
```
* make 可以在创建切片的时候指定我们需要的长度和容量;我们可以使用切片的索引访问切片内的元素，使用 = 对切片元素进行赋值。
* 使用 len 获取数组和切片的长度
```go
length0fNumbers := len(numbersToSum)
sums = make([]int, length0fNumbers)

for i, numbers := range numbersToSum {
sums[i] = Sum(numbers)
}

```
* 如何获取部分切片: 获取部分切片。如果在冒号的一侧没有数字就会一直取到最边缘的元素。在我们的函数中，我们使用 numbers[1:] 取到从索引 1 到最后一个元素
```go
 slice[low:high]
```
* 使用测试代码覆盖率的工具
* reflect.DeepEqual 的妙用和对代码类型安全性的影响: 判断两个值深度是否一直，即除了值相同，底层类型也相同
```go
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

```

* invalid operation: got != want (slice can only be compared to nil)
* 在 Go 中不能对切片使用等号运算符。你可以写一个函数迭代每个元素来检查它们的值。但是一种比较简单的办法是使用 reflect.DeepEqual，它在判断两个变量是否相等时十分有用。
* 需要注意的是 reflect.DeepEqual 不是「类型安全」的，所以有时候会发生比较怪异的行为：比较不同类型的参数

