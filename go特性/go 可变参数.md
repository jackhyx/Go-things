package go_use

import (
	"fmt"
	"reflect"
	"testing"
)

// c语言实现
#include <stdio.h>
#include <stdarg.h>

int sum(int count, ...)
{
int sum=0;
int val=0;

// 定义一个可变参数列表，可以看成是一种特殊的指针类型，list指向的对象是栈上的数据。
va_list list;

// 初始化list，list指向第一个被压栈的参数，即函数的最后一个参数，而count则是栈上最后一个参数，系统由此确定栈上参数内存的范围。
va_start(list, count);

while(count--)
{
// 通过va_arg依次获取参数值，两个参数，一个是指向栈上参数的指针list，这个指针每取出一个数据移动一次，总是指向栈上第一个未取出的参数，int指需要取出的参数的类型，CPU根据这个类型所占的地址空间来进行寻址。
val = va_arg(list,int);
// printf("%d at %X\n", val, &val);
sum += val;
}

//释放list，释放后list置为空。
va_end(list);

return sum;
}

int main()
{
printf("Sum of 1,2,3,4,5 = %d\n", sum(5, 1, 2, 3, 4, 5));
printf("Sum of 10,20,30 = %d\n", sum(3, 10,20,30));
}

/*
实例一：
func1使用的是Go语言的语法糖，按照内部机制来说，
...type本质是一个切片，也就是[]type，params被看作是类型为[] int的切片传入func1中，func1可接收任意个int值，返回sum结果。
虽然在可变参数函数内部，...int型参数的行为看起来类似slice，实际上，可变参数函数和切片作为参数的函数是不相同的。

 */

// 可变参数
func func1(params ...int) {

	sum := 0
	for _, param := range params {
		sum += param
	}

	fmt.Println("params : ", params, "\tsum : ", sum)
}
// 调用一
func TestFunc1(t *testing.T) {
	var params = []int{1, 2, 3}
	func1(params...)
	func1(2, 5)
	func1(2, 5, 8)
}
//结果一
=== RUN   TestFunc1
params :  [1 2 3]       sum :  6
params :  [2 5]         sum :  7
params :  [2 5 8]       sum :  15
--- PASS: TestFunc1 (0.00s)

/*
实例二：

func2虽然同样实现了不定参数的功能，但是使用起来比较繁琐，需要[]type{}来构造切片实例。
我们可以看到传递的数据是slice，但是在参数传递的时候，我们需要手工初始化slice再传入函数。

 */
// 切片
func func2(params []int) {

	sum := 0
	for _, param := range params {
		sum += param
	}

	fmt.Println("params : ", params, "\tsum : ", sum)
}

// 调用2
func TestFunc2(t *testing.T) {
	func2([]int{3})
	func2([]int{3, 6})
	func2([]int{3, 6, 9})
}
// 结果2
=== RUN   TestFunc2
params :  [3]   sum :  3
params :  [3 6]         sum :  9
params :  [3 6 9]       sum :  18
--- PASS: TestFunc2 (0.00s)

// 实例三：

//我们再看一下可变类型的可变参数，见func3：
// 可变类型的可变参数
func func3(params ...interface{}) {

	for _, param := range params {
		switch reflect.TypeOf(param).Kind().String() {
		case "int":
			fmt.Printf("param:%d is an int value!\n", param)
		case "int32":
			fmt.Printf("param:%v is an int32 value!\n", param)
		case "int64":
			fmt.Printf("param:%v is an int64 value!\n", param)
		case "float32":
			fmt.Printf("param:%v is an float32 value!\n", param)
		case "float64":
			fmt.Printf("param:%v is an float64 value!\n", param)
		case "string":
			fmt.Printf("param:%s is an string value!\n", param)
		case "func":
			fmt.Printf("param:%v is an func value!\n", param)
		case "map":
			fmt.Printf("param:%v is an map value!\n", param)
		default:
			fmt.Printf("param:%v is an unknown type.\n", param)
		}
	}
}
//
func TestFunc3(t *testing.T) {
	var p1 int = 100 //传递int值
	func3(p1)

	var p2 int32 = 200 //传递int32
	func3(p2)

	var p3 int64 = 300 //传递int64
	func3(p3)

	var p4 = "test string" //传递string
	func3(p4)

	var p5 float32 = 1.11 //传递float32
	func3(p5)

	var p6 float64 = 2.22 //传递float64
	func3(p6)

	var p7 = func(a, b int) int { return a + b } //传递func
	func3(p7)

	var p8 = map[string]string{} 传递map
	func3(p8)
}
//
=== RUN   TestFunc3
param:100 is an int value!
param:200 is an int32 value!
param:300 is an int64 value!
param:test string is an string value!
param:1.11 is an float32 value!
param:2.22 is an float64 value!
param:0x506b50 is an func value!
param:map[] is an map value!
--- PASS: TestFunc3 (0.00s)
PASS