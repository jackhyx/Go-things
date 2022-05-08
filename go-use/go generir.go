package go_use

import "fmt"

/* 泛型语法详解
小结：
类型参数  T
约束：利用interface


type Addable interface {
    type int, int8, int16, int32, int64,
        uint, uint8, uint16, uint32, uint64, uintptr,
        float32, float64, complex64, complex128,
        string
}


package main

import (
    "fmt"
)

type Addable interface {
    type int, int8, int16, int32, int64,
        uint, uint8, uint16, uint32, uint64, uintptr,
        float32, float64, complex64, complex128,
        string
}

func add[T Addable](a, b T) T {
    return a + b
}

func main() {
    fmt.Println(add(1,2))

    fmt.Println(add("foo","bar"))
}

## 类比C++







下面开始详细介绍泛型的语法

MyType[T1 constraint1 | constraint2, T2 constraint3...] ...

泛型的语法非常简单, 就类似于上面这样, 其中:

MyType可以是函数名, 结构体名, 类型名…
T1, T2…是泛型名, 可以随便取
constraint的意思是约束, 也是泛型中最重要的概念, 接下来会详解constraint
使用 | 可以分隔多个constraint, T满足其中之一即可(如T1可以是constraint1和constraint2中的任何一个)

*/




// Constraint(约束)是什么

约束的意思是限定范围, constraint的作用就是限定范围, 将T限定在某种范围内

而常用的范围, 我们自然会想到的有:

any(interface{}, 任何类型都能接收, 多方便啊!)
Interger(所有int, 多方便啊, int64 int32…一网打尽)
Float(同上)
comparable(所有可以比较的类型, 我们可以给所有可以比较的类型定制一些方法)
…
这些约束, 不是被官方定义为内置类型, 就是被涵盖在了constraints包内!!!

下面是builtin.go的部分官方源码:

// any is an alias for interface{} and is equivalent to interface{} in all ways.
type any = interface{}

// comparable is an interface that is implemented by all comparable types
// (booleans, numbers, strings, pointers, channels, interfaces,
// arrays of comparable types, structs whose fields are all comparable types).
// The comparable interface may only be used as a type parameter constraint,
// not as the type of a variable.

type comparable comparable

// 下面是constraints.go的部分官方源码:

// Integer is a constraint that permits any integer type.
// If future releases of Go add new predeclared integer types,
// this constraint will be modified to include them.

type Integer interface {
	Signed | Unsigned
}

// Float is a constraint that permits any floating-point type.
// If future releases of Go add new predeclared floating-point types,
// this constraint will be modified to include them.
type Float interface {
	~float32 | ~float64
}
//......

// 可以看到, 官方还是非常贴心的, 很多轮子已经帮我们造好了

// 而通过观察constraints包和阅读官方文档, 我也掌握了如何自定义约束

// 自定义constraint(约束): 下面是constraints包中的官方源码:

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// 泛型中的"~"符号是什么 : 符号"~"都是与类型一起出现的，用来表示支持该类型的衍生类型
// int8的衍生类型
type int8A int8
type int8B = int8

// 不仅支持int8, 还支持int8的衍生类型int8A和int8B
type MyInt interface {
	~int8
}

// 泛型的进阶使用
// 创建一个带有泛型的结构体User，提供两个获取age和name的方法
// 注意：只有在结构体上声明了泛型，结构体方法中才可以使用泛型
type AgeT interface {
	int8 | int16
}

type NameE interface {
	string
}

type User[T AgeT, E NameE] struct {
	age  T
	name E
}

// 获取age
func (u *User[T, E]) GetAge() T {
	return u.age
}


// 获取name
func (u *User[T, E]) GetName() E {
	return u.name
}

我们可以通过声明结构体对象时，声明泛型的类型来使用带有泛型的结构体
// 声明要使用的泛型的类型
var u User[int8, string]

// 赋值
u.age = 18
u.name = "weiwei"

// 调用方法
age := u.GetAge()
name := u.GetName()

// 输出结果 18 weiwei
fmt.Println(age, name)


/* Signed约束就是这样被写出来的, 其中需要我们get的点有如下几个:

使用interface{}就可以自定义约束
使用 | 就可以在该约束中包含不同的类型, 例如int, int8, int64均满足Signed约束
你可能会有疑问, ~是什么??? int我认识, ~int我可不认识呀??? 没关系, 实际上~非常简单, 它的意思就是模糊匹配, 例如:
type MyInt int64
此时 MyInt并不等同于int64类型(Go语言特性)
若我们使用int64来约束MyInt, 则Myint不满足该约束
若我们使用~int64来约束MyInt, 则Myint满足该约束(也就是说, ~int64只要求该类型的底层是int64, 也就是模糊匹配了)
官方为了鲁棒性, 自然把所有的类型前面都加上了~
下面我们自定义一个约束

type My_64_Bits_Long_Num interface {
	~int64 | ~float64
}
1
2
3
是不是很简单?
*/

// 泛型综合使用案例
// 自定义类型:这里自定义一个map :

// map的key必须要可以比较, 也就是可以被 == 和 != 比较(用于处理哈希冲突)
type MyMap[K comparable, V constraints.Integer | constraints.Float] map[K]V

func main() {
	m := make(MyMap[string, int])
	m["表哥"] = 100
	m["小张"] = 0
	for k, v := range m{
		fmt.Printf("key: %v, val: %v\n", k, v)
	}
}
// 自定义结构体: 这里以一个手写链表(只能存储整数)做示范

type MyIntergerNode[T constraints.Integer] struct {
	Next *MyIntergerNode[T] //注意这里一定要加类型声明(和C++一样)
	Data T
}

func main() {
	head := &MyIntergerNode[int64]{Next: nil, Data: 1}
	head.Next = &MyIntergerNode[int64]{Next: nil, Data: 2}

	for p := head; p != nil; p = p.Next{
		fmt.Printf("%d ", p.Data)
	}
}

// 自定义函数:这里自定义一个比较64比特大小的类型的函数

//刚才自定义的约束
type My_64_Bits_Long_Num interface {
	~int64 | ~float64
}

func MyCompare[T My_64_Bits_Long_Num](a, b T) bool {
	return a < b
}

func main() {
	var a int64 = 1
	var b int64 = 8

	//函数可以省略不写参数类型(语法糖)
	ans := MyCompare(a, b)
	if ans{
		fmt.Printf("%v小于%v", a, b)
	}else{
		fmt.Printf("%v大于%v", a, b)
	}
}

// 泛型的限制或缺陷 :无法直接和switch配合使用,将泛型和switch配合使用时，无法通过编译
func Get[T any]() T {
	var t T

	switch T {
	case int:
		t = 18
	}

	return t
}
// 只能先将泛型赋值给interface才可以和switch配合使用
func Get[T any]() T {
	var t T

	var ti interface{} = &t
	switch v := ti.(type) {
	case *int:
		*v = 18
	}

	return t
}