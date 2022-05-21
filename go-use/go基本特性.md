### 值方法和指针方法的区别
* 我们都知道，方法的接收者类型必须是某个自定义的数据类型，而且不能是接口类型或接口的指针类型。所谓的值方法，就是接收者类型是非指针的自定义数据类型的方法。那么，值方法和指针方法体现在哪里呢？我们看下边的代码：
```
func (cat *Cat) SetName(name string) {
cat.name = name
}
```

* 方法SetName的接收者类型是*Cat。Cat左边再加个*代表的就是Cat类型的指针类型,这时，Cat可以被叫做*Cat的基本类型。你可以认为这种指针类型的值表示的是指向某个基本类型值的指针。那么，这个SetName就是指针方法。那么什么是值方法呢？通俗的讲，把Cat前边的*去掉就是值方法。指针方法和值方法究竟有什么区别呢？请看下文。
* 值方法的接收者是该方法所属的那个类型值的一个副本。我们在该方法内对该副本的修改一般都不会体现在原值上，除非这个类型本身是某个引用类型（比如切片或字典）的别名类型。
* 而指针方法的接收者，是该方法所属的那个基本类型值的指针值的一个副本。我们在这样的方法内对该副本指向的值进行修改，却一定会体现在原值上。这块可能有点绕，但如果之前函数传切片那块理解的话这块也可以想明白，总之就是一个拷贝的是整个数据结构，一个拷贝的是指向数据结构的地址。

* 一个自定义数据类型的方法集合中仅会包含它的所有值方法，而该类型的指针类型的方法集合却囊括了前者的所有方法，包括所有值方法和所有指针方法。

* 严格来讲，我们在这样的基本类型的值上只能调用到它的值方法。但是，Go 语言会适时地为我们进行自动地转译，使得我们在这样的值上也能调用到它的指针方法。

例如下边这种也是可以调用的：
```
type Pet interface {
Name() string
}

type Dog struct {
Class string
}

func (dog Dog) Name() string{
return dog.Class
}

func (dog *Dog) SetName(name string) {
dog.Class = name
}

func main() {
a := Dog{"grape"}
a.SetName("nosay") //a会先取地址然后去调用指针方法
//Dog{"grape"}.SetName("nosay") //因为是值类型，调用失败，cannot call pointer method       on Dog literal，cannot take the address of Dog literal
(&Dog{"grape"}).SetName("nosay") //可以
}
```
在后边你会了解到，一个类型的方法集合中有哪些方法与它能实现哪些接口类型是息息相关的。如果一个基本类型和它的指针类型的方法集合是不同的，那么它们具体实现的接口类型的数量就也会有差异，除非这两个数量都是零。

比如，一个指针类型实现了某某接口类型，但它的基本类型却不一定能够作为该接口的实现类型。例如：
```
type Pet interface {
SetName(name string)
Name()string
Category()string
}

type Dog struct {
name string
}

func (dog *Dog) SetName(name string) {
dog.name = name
}

func(dog Dog) Name()string{
return dog.name
}

func (dog Dog)Category()string{
return "dog"
}

func main() {
dog:=Dog{"little pig"}

_,ok:=interface{}(dog).(Pet)
fmt.Printf("Dog implements interface Pet: %v\n", ok) //false
_, ok = interface{}(&dog).(Pet)
fmt.Printf("*Dog implements interface Pet: %v\n", ok)
fmt.Println() //true
}

```

## 基本语法——变量

## 一、变量的使用

### 1.1 什么是变量

变量是为存储特定类型的值而提供给内存位置的名称。在go中声明变量有多种语法。

所以变量的本质就是一小块内存，用于存储数据，在程序运行过程中数值可以改变



### 1.2 声明变量

var名称类型是声明单个变量的语法。

> 以字母或下划线开头，由一个或多个字母、数字、下划线组成

声明一个变量

第一种，指定变量类型，声明后若不赋值，使用默认值

```go
var name type
name = value
```

第二种，根据值自行判定变量类型(类型推断Type inference)

如果一个变量有一个初始值，Go将自动能够使用初始值来推断该变量的类型。因此，如果变量具有初始值，则可以省略变量声明中的类型。

```go
var name = value
```

第三种，省略var, 注意 :=左侧的变量不应该是已经声明过的(多个变量同时声明时，至少保证一个是新变量)，否则会导致编译错误(简短声明)



```go
name := value

// 例如
var a int = 10
var b = 10
c : = 10
```

> 这种方式它只能被用在函数体内，而不可以用于全局变量的声明与赋值

示例代码：

```go
package main

var a = "Hello"
var b = "World"
var c bool

func main() {
	println(a, b, c)
}
```

运行结果：

```go
Hello World false
```

#### 多变量声明

第一种，以逗号分隔，声明与赋值分开，若不赋值，存在默认值

```go
var name1, name2, name3 type
name1, name2, name3 = v1, v2, v3
```

第二种，直接赋值，下面的变量类型可以是不同的类型

```go
var name1, name2, name3 = v1, v2, v3
```

第三种，集合类型

```go
var (
    name1 type1
    name2 type2
)
```

### 1.3 注意事项

- 变量必须先定义才能使用
- go语言是静态语言，要求变量的类型和赋值的类型必须一致。
- 变量名不能冲突。(同一个作用于域内不能冲突)
- 简短定义方式，左边的变量名至少有一个是新的
- 简短定义方式，不能定义全局变量。
- 变量的零值。也叫默认值。
- 变量定义了就要使用，否则无法通过编译。

如果在相同的代码块中，我们不可以再次对于相同名称的变量使用初始化声明，例如：a := 20 就是不被允许的，编译器会提示错误 no new variables on left side of :=，但是 a = 20 是可以的，因为这是给相同的变量赋予一个新的值。

如果你在定义变量 a 之前使用它，则会得到编译错误 undefined: a。如果你声明了一个局部变量却没有在相同的代码块中使用它，同样会得到编译错误，例如下面这个例子当中的变量 a：

```go
func main() {
   var a string = "abc"
   fmt.Println("hello, world")
}
```

尝试编译这段代码将得到错误 a declared and not used

此外，单纯地给 a 赋值也是不够的，这个值必须被使用，所以使用

在同一个作用域中，已存在同名的变量，则之后的声明初始化，则退化为赋值操作。但这个前提是，最少要有一个新的变量被定义，且在同一作用域，例如，下面的y就是新定义的变量

```go
package main

import (
	"fmt"
)

func main() {
	x := 140
	fmt.Println(&x)
	x, y := 200, "abc"
	fmt.Println(&x, x)
	fmt.Print(y)
}
```

运行结果：

```go
0xc04200a2b0
0xc04200a2b0 200
abc
```





# 基本语法——常量constant

## 一、常量的使用

### 1.1 常量声明

常量是一个简单值的标识符，在程序运行时，不会被修改的量。

```go
const identifier [type] = value
```

```go
显式类型定义： const b string = "abc"
隐式类型定义： const b = "abc"
```
```go
package main

import "fmt"

func main() {
   const LENGTH int = 10
   const WIDTH int = 5   
   var area int
   const a, b, c = 1, false, "str" //多重赋值

   area = LENGTH * WIDTH
   fmt.Printf("面积为 : %d", area)
   println()
   println(a, b, c)   
}
```
运行结果：

```go
面积为 : 50
1 false str
```

常量可以作为枚举，常量组

```go
const (
    Unknown = 0
    Female = 1
    Male = 2
)
```
常量组中如不指定类型和初始化值，则与上一行非空常量右值相同

```go
package main

import (
	"fmt"
)

func main() {
	const (
		x uint16 = 16
		y
		s = "abc"
		z
	)
	fmt.Printf("%T,%v\n", y, y)
	fmt.Printf("%T,%v\n", z, z)
}
```
运行结果：

```go
uint16,16
string,abc
```

常量的注意事项：

- 常量中的数据类型只可以是布尔型、数字型（整数型、浮点型和复数）和字符串型

- 不曾使用的常量，在编译的时候，是不会报错的

- 显示指定类型的时候，必须确保常量左右值类型一致，需要时可做显示类型转换。这与变量就不一样了，变量是可以是不同的类型值



### 1.2 iota

iota，特殊常量，可以认为是一个可以被编译器修改的常量

iota 可以被用作枚举值：

```go
const (
    a = iota
    b = iota
    c = iota
)
```
第一个 iota 等于 0，每当 iota 在新的一行被使用时，它的值都会自动加 1；所以 a=0, b=1, c=2 可以简写为如下形式：

```go
const (
    a = iota
    b
    c
)
```
**iota 用法**

```go
package main

import "fmt"

func main() {
    const (
            a = iota   //0
            b          //1
            c          //2
            d = "ha"   //独立值，iota += 1
            e          //"ha"   iota += 1
            f = 100    //iota +=1
            g          //100  iota +=1
            h = iota   //7,恢复计数
            i          //8
    )
    fmt.Println(a,b,c,d,e,f,g,h,i)
}
```
运行结果：

```
0 1 2 ha ha 100 100 7 8
```

如果中断iota自增，则必须显式恢复。且后续自增值按行序递增

自增默认是int类型，可以自行进行显示指定类型

数字常量不会分配存储空间，无须像变量那样通过内存寻址来取值，因此无法获取地址

### 格式化打印中的常用占位符：

```
格式化打印占位符：
			%v,原样输出
			%T，打印类型
			%t,bool类型
			%s，字符串
			%f，浮点
			%d，10进制的整数
			%b，2进制的整数
			%o，8进制
			%x，%X，16进制
				%x：0-9，a-f
				%X：0-9，A-F
			%c，打印字符
			%p，打印地址
			...
```
### channel

```

type hchan struct {
// 通道里元素的数量
qcount   uint
// 循环队列的长度
dataqsiz uint
// 指针，指向存储缓冲通道数据的循环队列
buf      unsafe.Pointer
// 通道中元素的大小
elemsize uint16
// 通道是否关闭的标志
closed   uint32
// 通道中元素的类型
elemtype *_type
// 已接收元素在循环队列的索引
sendx    uint  
// 已发送元素在循环队列的索引
recvx    uint
// 等待接收的协程队列
recvq    waitq
// 等待发送的协程队列
sendq    waitq
// 互斥锁，保护hchan的并发读写，下文会讲
lock mutex
}
```

* waitq:这两个字段就是使用两个双向链表，来存储所有等待的协程的。一旦通道不再为空或者不再为满，那么Go协程的调度器就会去这个双向链表中唤醒一个链表中的协程，允许这个协程往通道中写入/接收数据也同理
``` 
  type waitq struct {
  first *sudog
  last  *sudog
  }
```