### 01. 数组和切片有什么区别？

* 切片是对数组的抽象，因为数组的长度是不可变的 切片"动态数组"

```
var slice []int // 直接声明
slice := []int{1,2,3,4,5} // 字面量方式
slice := make([]int, 5, 10) // make创建
slice := array[1:5] // 截取下标的方式
slice := *new([]int) // new一个
切片可以使用append追加元素，当cap不足时进行动态扩容。
```

### 02. 拷贝大切片一定比拷贝小切片代价大吗？

```
实际上不会，因为切片本质内部结构如下：
// SliceHeader 切片动态时表现
type SliceHeader struct {
 Data uintptr
 Len  int
 Cap  int
}
切片中的第一个字是指向切片底层数组的指针，这是切片的存储空间，第二个字段是切片的长度，第三个字段是容量。
将一个切片变量分配给另一个变量只会复制三个机器字，大切片跟小切片的区别无非就是 Len 和 Cap的值比小切片的这两个值大一些，如果发生拷贝，本质上就是拷贝上面的三个字段。
```

### 03. 切片的深浅拷贝

* 深浅拷贝都是进行复制，区别在于复制出来的新对象与原来的对象在它们发生改变时，是否会相互影响，本质区别就是复制出来的对象与原对象是否会指向同一个地址。在Go语言，切片拷贝有三种方式：

* 使用=操作符拷贝切片，这种就是浅拷贝

* 使用[:]下标的方式复制切片，这种也是浅拷贝

* 使用Go语言的内置函数copy()进行切片拷贝，这种就是深拷贝

### 04. 零切片、空切片、nil切片是什么

* 为什么问题中这么多种切片呢？因为在Go语言中切片的创建方式有五种，不同方式创建出来的切片也不一样；

#### 零切片

* 我们把切片内部数组的元素都是零值或者底层数组的内容就全是 nil的切片叫做零切片，使用make创建的、长度、容量都不为0的切片就是零值切片：

```
slice := make([]int,5) // 0 0 0 0 0
slice := make([]*int,5) // nil nil nil nil nil

```

#### nil切片

* nil切片的长度和容量都为0，并且和nil比较的结果为true，采用直接创建切片的方式、new创建切片的方式都可以创建nil切片：

```
var slice []int
var slice = *new([]int) new出来的是一个地址 所以需要*
```

#### 空切片

* 空切片的长度和容量也都为0，但是和nil的比较结果为false，因为所有的空切片的数据指针都指向同一个地址
  0xc42003bda0；使用字面量、make可以创建空切片：

```
var slice = []int{}
var slice = make([]int, 0)
空切片指向的 zerobase 内存地址是一个神奇的地址，从 Go 语言的源代码中可以看到它的定义：

// base address for all 0-byte allocations
var zerobase uintptr

// 分配对象内存
func mallocgc(size uintptr, typ *_type, needzero bool) unsafe.Pointer {
...
if size == 0 {
return unsafe.Pointer(&zerobase)
}
...
```

### 05. 切片的扩容策略

这个问题是一个高频考点，我们通过源码来解析一下切片的扩容策略，切片的扩容都是调用growslice方法，截取部分重要源代码：

```go
// runtime/slice.go
// et：表示slice的一个元素；old：表示旧的slice；cap：表示新切片需要的容量；
func growslice(et *_type, old slice, cap int) slice {
if cap < old.cap {
panic(errorString("growslice: cap out of range"))
}

if et.size == 0 {
// append should not create a slice with nil pointer but non-zero len.
// We assume that append doesn't need to preserve old.array in this case.
return slice{unsafe.Pointer(&zerobase), old.len, cap}
}

newcap := old.cap
// 两倍扩容
doublecap := newcap + newcap
// 新切片需要的容量大于两倍扩容的容量，则直接按照新切片需要的容量扩容
if cap > doublecap {
newcap = cap
} else {
// 原 slice 容量小于 1024 的时候，新 slice 容量按2倍扩容
if old.cap < 1024 {
newcap = doublecap
} else { // 原 slice 容量超过 1024，新 slice 容量变成原来的1.25倍。
// Check 0 < newcap to detect overflow
// and prevent an infinite loop.
for 0 < newcap && newcap < cap {
newcap += newcap / 4
}
// Set newcap to the requested cap when
// the newcap calculation overflowed.
if newcap <= 0 {
newcap = cap
}
}
}

// 后半部分还对 newcap 作了一个内存对齐，这个和内存分配策略相关。进行内存对齐之后，新 slice 的容量是要 大于等于 老 slice 容量的 2倍或者1.25倍。
var overflow bool
var lenmem, newlenmem, capmem uintptr
// Specialize for common values of et.size.
// For 1 we don't need any division/multiplication.
// For sys.PtrSize, compiler will optimize division/multiplication into a shift by a constant.
// For powers of 2, use a variable shift.
switch {
case et.size == 1:
lenmem = uintptr(old.len)
newlenmem = uintptr(cap)
capmem = roundupsize(uintptr(newcap))
overflow = uintptr(newcap) > maxAlloc
newcap = int(capmem)
case et.size == sys.PtrSize:
lenmem = uintptr(old.len) * sys.PtrSize
newlenmem = uintptr(cap) * sys.PtrSize
capmem = roundupsize(uintptr(newcap) * sys.PtrSize)
overflow = uintptr(newcap) > maxAlloc/sys.PtrSize
newcap = int(capmem / sys.PtrSize)
case isPowerOfTwo(et.size):
var shift uintptr
if sys.PtrSize == 8 {
// Mask shift for better code generation.
shift = uintptr(sys.Ctz64(uint64(et.size))) & 63
} else {
shift = uintptr(sys.Ctz32(uint32(et.size))) & 31
}
lenmem = uintptr(old.len) << shift
newlenmem = uintptr(cap) << shift
capmem = roundupsize(uintptr(newcap) << shift)
overflow = uintptr(newcap) > (maxAlloc >> shift)
newcap = int(capmem >> shift)
default:
lenmem = uintptr(old.len) * et.size
newlenmem = uintptr(cap) * et.size
capmem, overflow = math.MulUintptr(et.size, uintptr(newcap))
capmem = roundupsize(capmem)
newcap = int(capmem / et.size)
}
}

```

#### 通过源代码可以总结切片扩容策略：

* 切片在扩容时会进行内存对齐，这个和内存分配策略相关。
* 进行内存对齐之后，新 slice 的容量是要 大于等于老 slice 容量的 2倍或者1.25倍;
* 当原 slice 容量小于 1024 的时候，新 slice 容量变成原来的 2 倍； 大于 1024 的时候，新 slice 容量变成原来的 1.25 倍

```
newcap := old.cap
  // 两倍扩容
 doublecap := newcap + newcap
  // 新切片需要的容量大于两倍扩容的容量，则直接按照新切片需要的容量扩容
 if cap > doublecap {
  newcap = cap
 } else {
    // 原 slice 容量小于 1024 的时候，新 slice 容量按2倍扩容
  if old.cap < 1024 {
   newcap = doublecap
```

* 原 slice 容量超过 1024，新 slice 容量变成原来的1.25倍。

```
else { // 原 slice 容量超过 1024，新 slice 容量变成原来的1.25倍。
   // Check 0 < newcap to detect overflow
   // and prevent an infinite loop.
   for 0 < newcap && newcap < cap {
    newcap += newcap / 4

```

### 07. 参数传递切片和切片指针有什么区别？

我们都知道切片底层就是一个结构体，里面有三个元素：

```
type SliceHeader struct {
Data uintptr
Len  int
Cap  int
}

分别表示切片底层数据的地址，切片长度，切片容量。
```

当切片作为参数传递时，其实就是一个结构体的传递，因为Go语言参数传递只有值传递，传递一个切片就会浅拷贝原切片，但因为底层数据的地址没有变，所以在函数内对切片的修改，也将会影响到函数外的切片，举例：

```
func modifySlice(s []string)  {
s[0] = "song"
s[1] = "Golang"
fmt.Println("out slice: ", s)
}

func main()  {
s := []string{"asong", "Golang梦工厂"}
modifySlice(s)
fmt.Println("inner slice: ", s)
}
// 运行结果
out slice:  [song Golang]
inner slice:  [song Golang]
不过这也有一个特例，先看一个例子：

func appendSlice(s []string)  {
s = append(s, "快关注！！")
fmt.Println("out slice: ", s)
}

func main()  {
s := []string{"asong", "Golang梦工厂"}
appendSlice(s)
fmt.Println("inner slice: ", s)
}
// 运行结果
out slice:  [asong Golang梦工厂 快关注！！]
inner slice:  [asong Golang梦工厂]
```

* 因为切片发生了扩容，函数外的切片指向了一个新的底层数组，所以函数内外不会相互影响
* 当参数直接传递切片时，如果指向底层数组的指针被覆盖或者修改（copy、重分配、append触发扩容），此时函数内部对数据的修改将不再影响到外部的切片，代表长度的len和容量cap也均不会被修改。

* 参数传递切片指针就很容易理解了，如果你想修改切片中元素的值，并且更改切片的容量和底层数组，则应该按指针传递。

### 08. range遍历切片有什么要注意的？

* Go语言提供了range关键字用于for 循环中迭代数组(array)、切片(slice)、通道(channel)或集合(map)的元素，有两种使用方式：

* 第一种是遍历下标和对应值，第二种是只遍历下标，使用range遍历切片时会先拷贝一份，然后在遍历拷贝数据：

```
for k,v := range _ { }
for k := range _ { }


s := []int{1, 2}
for k, v := range s {

}
会被编译器认为是
for_temp := s
len_temp := len(for_temp)
for index_temp := 0; index_temp < len_temp; index_temp++ {
value_temp := for_temp[index_temp]
_ = index_temp
value := value_temp

}
不知道这个知识点的情况下很容易踩坑，例如下面这个例子：

package main

import (
"fmt"
)

type user struct {
name string
age uint64
}

func main()  {
u := []user{
{"asong",23},
{"song",19},
{"asong2020",18},
}
for _,v := range u{
if v.age != 18{
v.age = 20
}
}
fmt.Println(u)
}
// 运行结果
[{asong 23} {song 19} {asong2020 18}]

```

* 因为使用range遍历切片u，变量v是拷贝切片中的数据，修改拷贝数据不会对原切片有影响。