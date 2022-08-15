



#### 区别
本质上在于 make 函数在初始化时，会初始化 slice、chan、map 类型的内部数据结构，new 函数并不会。
例如：在 map 类型中，合理的长度（len）和容量（cap）可以提高效率和减少开销。
更进一步的区别：

* make 函数：
  能够创建类型所需的内存空间，返回引用类型的本身。
  具有使用范围的局限性，仅支持 channel、map、slice 三种类型。
  具有独特的优势，make 函数会对三种类型的内部数据结构（长度、容量等）赋值。
* new 函数：
  能够创建并分配类型所需的内存空间，返回指针引用（指向内存的指针）。
  可被替代，能够通过字面值快速初始化。









#### Go语言中，开发者仅需声明变量，Go语言就会根据变量的类型自动分配响应内存。内存分为两部分：

* 栈内存，go语言管理，分配和释放
* 堆内存，开发者需要关注
```
声明变量
var s string // 零值是 “”
var sp *string // 零值是 nil
```
* 仅声明变量，没有初始化，则值都默认为零值
变量赋值
// 声明直接初始化
var s string ="hello"
// 声明后再初始化
var s string
s = "hello"
// 简单声明
s := "hello"
```
// 指针类型 -->编译时将报错：panic: runtime error: invalid memory address or nil pointer dereference
var sp * string
*sp = "hello"

```
* 值类型，没有初始化时直接赋值，没有问题 
* 指针类型，如果没有分配内存，则默认零值是nil，没有指向内存将无法使用。所以以上方式将会报错 
* 指针类型的变量必须要经过声明、内存分配才能赋值，才可以在声明时进行初始化。GO语言在指针类型声明时，不会自动分配内存，所以不能赋值操作。 
* 分配内存可以使用new函数或者make函数
* new函数 
* 通过new函数分配内存并返回指向该内存的指针，就可以通过指针对这块内存进行赋值、取值等操作 
* new函数只用于分配内存，并把内存清零，返回一个指向对应类型的零值的指针。 一般用于需显示返回指针的情况
```
  type person struc{
  name string
  age  int
  }

// 工厂函数,通过不同的参数构建不同的*person变量
func New(name string, age int) *person{
p := new(person)
p.name = name
p.age = age
return &p
}
```
* make函数
* 只用于slice、chan和map这三种内置类型的创建和初始化
m := make(map[string]int,10)

* make函数就是map类型的工厂函数，可根据传递它的K-V键值对类型，创建不同类型的map，同时可以初始化map的大小。



### 基本特性
#### make
* 在 Go 语言中，内置函数 make 仅支持 slice、map、channel 三种数据类型的内存创建，其返回值是所创建类型的本身，而不是新的指针引用。
```
函数签名如下：
func make(t Type, size ...IntegerType) Type
复制代码
具体使用示例：
func main() {
v1 := make([]int, 1, 5)
v2 := make(map[int]bool, 5)
v3 := make(chan int, 1)

	fmt.Println(v1, v2, v3)
}
```
在代码中，我们分别对三种类型调用了 make 函数进行了初始化。你会发现有的入参是有多个长度指定，有的没有。
这块的区别主要是长度（len）和容量（cap）的指定，有的类型是没有容量这一说法，因此自然也就无法指定。
输出结果：
[0] map[] 0xc000044070

* 有一个细节点要注意，调用 make 函数去初始化切片（slice）的类型时，会带有零值，需要明确是否需要。
见过不少的小伙伴在这上面踩坑。
#### new
在 Go 语言中，内置函数 new 可以对类型进行内存创建和初始化。其返回值是所创建类型的指针引用，与 make 函数在实质细节上存在区别。
```
函数签名如下：
func new(Type) *Type

具体使用示例：
type T struct {
Name string
}

func main() {
v := new(T)
v.Name = "煎鱼"
}

从上面的例子的效果来看，是不是似曾相似？其实与下面这种方式的一样的：
func main() {
v := T{}
v.Name = "煎鱼"
}

输出结果均是：
&{Name:煎鱼}
```
其实 new 函数在日常工程代码中是比较少见的，因为他可被替代。
一般会直接用快捷的 T{} 来进行初始化，因为常规的结构体都会带有结构体的字面属性：
func NewT() *T {
return &T{Name: "煎鱼"}
}
复制代码
这种初始化方式更方便。
区别在哪里
可能会有的小伙伴会疑惑一点，就是 new 函数也能初始化 make 的三种类型：
v1 := new(chan bool)
v2 := new(map[string]struct{})


