### 字符串-string
* Go语言中有两种特殊的别名类型，是byte和rune，分别代表uint8和int32类型，即1个字节和4个字节。

* 字符串的底层实现:string是一个8bit字节的集合，且是不可变的
```

type stringStruct struct {
	str unsafe.Pointer // 指针，指向底层存储数据的[]byte
	len int            // 长度
}

```
* 字符串的不可变性:存储空间复用机制COW(copy on write)的体现；
* 把分配内存空间这种耗时操作推迟到最晚（也就是修改后必须分离）的时候才完成，减少了内存分配的次数、最大化复用同一个底层数组的时间。；
* 给多个字符串以及共享相同的底层数据结构带来了最大程度的优化。同时也保证了在Go的多协程状态下，操作字符串的安全性。






### 切片-slice

* 切片赋值会导致底层数据的变化，从而影响其它的切片值
```

func main() {
	var c = [4]int{1,2,3,4}
	var Aslice = c[0:2]
	Aslice = append(Aslice,5)
	fmt.Println(c) //[1 2 5 4] 改变了底层数组
	fmt.Println(Aslice) //[1 2 5]
	Bslice := append(Aslice,5,5,5) //扩容超过底层数组的容量
	fmt.Println(c) //[1 2 5 4]
	fmt.Println(Bslice) //[1 2 5 5 5 5] //指向了新的数组
}	

```
### 切片是引用类型
```
func main() {
	//a是一个数组，注意数组是一个固定长度的，初始化时候必须要指定长度，不指定长度的话就是切片了
	a := [3]int{1, 2, 3}
	//b是数组，是a的一份拷贝
	b := a
	//c是切片，是引用类型，底层数组是a
	c := a[:]
	for i := 0; i < len(a); i++ {
	a[i] = a[i] + 1
	}
     //改变a的值后，b是a的拷贝，b不变，c是引用，c的值改变
	fmt.Println(a) //[2,3,4]
	fmt.Println(b) //[1 2 3]
	fmt.Println(c) //[2,3,4]
}

```
* 在函数传参中是值传递，所以会copy一份原始的切片，但是指向底层数组的指针不变，如果我们在函数中对这个copy过的切片操作（非赋值），例如重新进行切片操作，这样不会影响原切片，但是如果我们在此进行例如a[0]=1此类的操作，会修改原数组 
* 对于slice来说来说，在Go语言当中，切片类型是不可比较的
### 切片，数组可进行赋值操作
```
func main() {
	//切片
	var a = make([]string,10)
	//a[0] = 1 //赋值其他类型均报错
	a[0] = "grape"
                    
	//数组
	var a = [3]int{}
	a[0] =1
	a[1] = "strin" //赋值其他类型均报错
}   
```
* 基本规则：对于每个赋值一定要类型一致，和其他一样，不同的类型不可以进行赋值操作。当然，interface{}例外
### 切片，数组和字符串的循环
* 切片数组字符串循环代码示例：
```
func main() {
	var a = [3]int{1,2,3}
	for i,v := range(a) { 
		fmt.Println(i,v)  // 0 1   1 2  2 3
     }
        
	var b = []int{3,4,5}
	for ide,v := range(b) {
		fmt.Println(i,v)  //0 3  1 4  2 5
	}
        
	var c = "hello world"
	hello := c[:5]  
	world := c[7:]
	fmt.Println(hello, world)  //hello  world
	for i,v := range(c) {
		fmt.Println(i,string(v))  // 'h', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd',
		fmt.Println(i,v)  //0 104 1 101 2 108 3 108 4 111 5 32 6 19990 9 30028 //range会转化底层byte为rune
	}
}

```
### 切片类型强转
```
func main() {
	var a = []float64{4, 2, 5, 7, 2, 1, 88, 1}
	//var c = ([]int)(a) //报错
	var b = make([]int, 8)
	for i,v := range a {
		b[i] = int(v)
	}
	fmt.Println(b)
}
```

### 字典-map
### map初始化与内存分配
* 首先，必须给map分配内存空间之后，才可以往map中添加元素：
```
func main() {
var m map[int]int // 使用var语法声明一个map，不会分配内存
m[1] = 1 // 报错：assignment to entry in nil map
}
```
* 如果你使用的是make来创建一个map，Go在声明的同时，会自动为map分配内存空间，不会报错：
```
func main() {
m := make(map[int]int) // make语法创建map
m[1] = 1 // ok
}
```

* map中get操作的返回值
我们直接看一个例子：
```
func main() {
m := make(map[int]int)
fmt.Println(m[1]) // 0
m[1] = 0
fmt.Println(m[1]) // 0
}
```
* 大家看到问题了吧，如果某个key-value对在map中并不存在，不像其他语言，我们访问这个key是并不会报错的，而是返回value的零值。如果是int，那就返回0。但是，如果我们真正的往map里添加一个key-value对，其值为0，那么我们如何区分是根本没有这个key-value对，还是有这个key-value对，但是值为0呢？其实，访问map中的元素这个表达式有两个返回值：
```
func main() {
m := make(map[int]int)
v, ok := m[1]
fmt.Println(v, ok) // 0, false
m[1] = 0
v, ok = m[1]
fmt.Println(v, ok) // 0, true
}
```

* 第一个返回值和之前的例子相同，而第二个返回值就可以被用来判断，是否map中存在这个key-value对。如果存在，返回true；反之返回false，我们通常可以与if联合进行使用：
```
func main() {
m := make(map[int]int)
if _, ok := m[1]; !ok {
fmt.Println("key不存在")
}
}
``` 
* map遍历的无序性
* 在Go语言中，多次遍历相同的map，得到的结果是不一样的：
```
func main() {
m := make(map[int]int)
m[0] = 1
m[1] = 2
m[3] = 5
for k, v := range m {
fmt.Println(k, v)
}
// 第一次遍历结果：
0 1
1 2
3 5
// 第二次遍历结果：
3 5
0 1
1 2
}
``` 
### 为什么map是引用类型
* 为什么我们常常把map视为引用类型？我们先看一个简单的例子：
```
func main() {
m := make(map[int]int)
m[1] = 1 // 赋一个初始值
test(m) // 函数调用
fmt.Println(m[1]) // 2
}

func test(m map[int]int) {
m[1] = 2 // 修改值
}
``` 
* 我们看到，当map作为函数参数传递的时候，在外部函数对map的修改，会影响到原来map的值，为什么会这样呢？
* 大家都知道，Go语言只有值传递，那么为什么我们还会有把指针传过去的错觉呢？这还要从字典get与set操作的底层实现说起。Go语言的map在底层是用hashtable来实现的。在我们用var语法声明一个map的时候，实际上就创建了一个hmap结构体：
```
type hmap struct {
count     int // 元素个数，调用 len(map) 时，直接返回此值
buckets    unsafe.Pointer // 指向一个bucket数组
...
}
```

* 我们主要关注count和buckets这两个字段。count就是指map元素的个数；而buckets是真正存储map的key-value对的地方。这也就可以解释为什么我们一开始那个坑的报错问题。我们用var m map[int]int声明的map，只是分配了一个hmap结构体而已，而buckets这个字段并没有分配内存空间。
* 所以，最后解答我们为什么是引用类型的问题。其实我们传给test函数的值，只是一个hmap结构体；而这个结构体里面又包含了一个bucket数组的指针，也就相当于，表面上我们传了个结构体值过去，而内部却是传了一个指针，这个指针所存储的地址，也就是指针指向的bucket数组结构并没有改变。我们如果对存储key-value对的bucket进行修改，如m[1] = 2这种操作，实际上修改的就是改变了外部函数的bucket值。
* 每一个bucket数组中存储的元素结构为bmap，这里真正存储着key与value的值：
```
  type bmap struct {
  tophash  [8]uint8   // tophash，在hash计算过程中会用到
  keys     [8]keyType // 存储key
  values   [8]keyType // 存储value
  pad      uintptr    // 填充，用于内存对齐
  overflow uintptr    // 溢出bucket，hash值相同时会用到
  }
```
### 为什么key有类型约束
* Go 语言字典的键类型不可以是函数类型、字典类型和切片类型，但是value可以为任意类型 原因：哈希冲突需要比较
* 哈希冲突的解决
* 如果插入之后当前bucket无法容纳这个元素，Go就会新分配一个bucket，用当前bucket的overflow字段指向这个新的bucket，然后往新的bucket里插入当前key-value对即可
* 如果overflow bucket数量过多，在get操作时，对这个overflow链表进行遍历的时间复杂度会大大升高，为了避免溢出bucket数量过多，Go语言会在超过某一个阈值的时候，触发扩容操作。Go语言bucket的扩容操作也是渐进式的
* Go语言结合了链地址法和开放定址法这两种方案

### 错误与异常处理
* 多返回值
```
func main() {
	res, err := json.Marshal(payload)
	if err != nil {
		return "", errors.New("序列化请求参数失败")
	}
}

```
### try-catch
* Java、PHP等语言提供了try-catch-finally的解决方案。
* try-catch彻底完成了对错误与正常代码逻辑的分离。我们用try代码块中包裹可能出现问题的代码，在catch中对这些问题代码统一进行错误处理。
```
try {
// 正常代码逻辑
} catch(\Exception $e) {
// 错误处理逻辑
} finally {
// 释放资源逻辑
}
```
### 资源的释放
* finally代码块比较特殊，它被常常用来做一些资源及句柄的释放工作。如果没有finally，我们的代码可能会像这样
``` 
  func main() {
  mutex := sync.Mutex{}
  // 加锁
  mutex.Lock()
  res, err := json.Marshal("abc")
  if err != nil {
  // 释放锁资源
  mutex.Unlock()
  // ....其余错误处理逻辑
  }
  file, err := os.Open("abvc")
  if err != nil {
  // 释放锁资源
  mutex.Unlock()
  // ....其余错误处理逻辑
  }
  mutex.Unlock()
  }
```
* 为了确保锁资源在代码结束之前一定要被释放，我们每次在错误处理逻辑中，都需要写一次mutex.Unlock代码，导致大量的代码冗余。finally代码块内的语句会在代码返回或者退出之前执行，而且是百分百会执行。这样，我们就可以把释放锁资源这一行代码放到finally块即可，且只用写一次，这样就解决了之前代码冗余率高的问题。
* 在Go语言中，defer()也同样解决了这个问题。我们用Go中的defer语句改写一下上述代码：
```
  func main() {
  mutex := sync.Mutex{}
  defer mutex.Unlock()
  mutex.Lock()
  res, err := json.Marshal("abc")
  if err != nil {
  // 错误处理
  }
  file, err := os.Open("abvc")
  if err != nil {
  // 错误处理
  }
  }
```
### Go错误处理的实现
* 接下来我们深入讲解Go语言中的错误处理实现。我们看一下之前讲过的例子中，json.Marshal方法的签名：
```
func Marshal(v interface{}) ([]byte, error)
```
* 我们重点关注最后一个error类型的参数，它是一个Go语言内置的接口类型。那么，我们为什么要用接口类型来抽象所有的错误类型呢？先别急，我们先自己想想。

* 简单版的实现
*在我们对字符串进行marshal操作的过程中，可能会产生好多种类型的错误。为了在marshal函数内部区分不同的错误类型，我们简单粗暴一点，可能会进行如下的处理：
```
func (e *encodeState) marshal(v interface{}, opts encOpts) (errorMsg string) {
// 操作1可能的错误
if errType1 := doOp1(), errType1 != nil {
err1 := errType1.getErrorMessage() // 获取errorType1的错误信息
return err1
}
// 操作2可能的错误
if errType2 := doOp2(), errType2 != nil {
err2 := errType2.getErrMsg() // 方法名和errorType1不同
return err2
}
return ""
}
```
* 们分析一下上面这段代码，操作doOp1可能会发生errorType1类型的错误，我们要返回给调用者errorType1类型中错误的字符串信息；doOp2也同理。这样做确实可以，但是还是有一些麻烦，我们看看还有没有其他方案来优化一下。

### 抽象一下试试
* 我们先简单介绍一下，Go语言用一个接口类型抽象了所有错误类型：
```
type error interface {
Error() string
}
``` 
* 这个接口定义了一个Error()方法，用于返回错误信息，我们先记下来，等会要用。同上个例子，我们给之前自定义的两种错误类型加点料，实现这个error接口：
```
type errType1 struct {}

// 实现接口方法
func (*errType1) Error() {
fmt.Println("我是错误类型1的信息")
}

type errType2 struct {}

// 实现接口方法
func (*errType2) Error() {
fmt.Println("我是错误类型1的信息")
}
``` 
* 然后在marshal()函数上稍作改动，使用这两种实现接口的错误类型：
```
func (e *encodeState) marshal(v interface{}, opts encOpts) (errorMsg string) {
// 操作1可能的错误
if errType1 := doOp1(), errType1 != nil {
return errType1.Error()
}
// 操作2可能的错误
if errType2 := doOp2(), errType2 != nil {
return errType2.Error()
}
return ""
} 
```
* 大家看到优势在哪里了吗？在我们调用每个错误类型的返回信息方法的时候，如果用我们一开始的方式，我们需要进入每一个错误类型的实现类中去翻看他的API，看看函数名是什么；而在第二种实现方案中，由于两种错误的实现类型均实现了Error()方法，这样，在marshal函数中如果想进行错误信息的获取，我们统一调用Error()函数，即可返回对应错误实现类的错误信息。
* 这其实就是一种依赖的倒置。调用方marshal()函数不再关注错误类型的具体实现类，里面有哪些方法，而转为依赖抽象的接口 
### panic和recover
* Go语言的panic和其他语言的error有点像。如果调用了panic，代码会立刻停止运行，一层一层向上冒泡并积累堆栈信息，直到调用栈顶结束，并打印出所有堆栈信息。
* panic没什么好说的，而recover我们需要好好聊一聊。recover专门用来恢复panic。也就是说，如果你在panic之前声明了recover语句，那么你就可以在panic之后使用recover接收到panic的信息。但是问题又来了，我们panic不是直接就退出程序了吗，就算声明了recover也执行不了呀。这个时候，我们就需要配合defer来使用了。defer能够让程序在panic之后，仍然执行一段收尾的代码逻辑。这样一来，我们就可以通过recover获得panic的信息，并对信息作出识别与处理了。仍然举上述的marshal的源码的例子，这次是真的源码了，不是我编的：
```
func (e *encodeState) marshal(v interface{}, opts encOpts) (err error) {
defer func() { // defer收尾
if r := recover(); r != nil { // recover恢复案发现场
if je, ok := r.(jsonError); ok { // 拿到panic的值，并转为错误来返回
err = je.error
} else {
panic(r)
}
}
}()
e.reflectValue(reflect.ValueOf(v), opts)
return nil
}
```
* 我们看到，源码中将defer与recover配合使用，直接改变了panic的运行逻辑。原本是panic之后会直接退出程序，这样一来，现在程序并不会直接退出，而是被转为了jsonError类型，并返回。
* 通过使用recover捕获运行时的panic，可以让代码继续运行下去而不至于直接停止。

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