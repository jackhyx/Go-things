### unsafe.Pointer
* 为了解决这个问题，需要用到 unsafe.Pointer。它表示任意类型且可寻址的指针值，可以在不同的指针类型之间进行转换（类似 C 语言的 void * 的用途）。

* 其包含四种核心操作：

* 任何类型的指针值都可以转换为 Pointer。

* Pointer 可以转换为任何类型的指针值。

* uintptr 可以转换为 Pointer。

* Pointer 可以转换为 uintptr。

* 在这一部分，重点看第一点、第二点。你再想想怎么修改 “错误的例子” 让它运行起来？
```
修改如下：

func main(){
 num := 5
 numPointer := &num
 
 flnum := (*float32)(unsafe.Pointer(numPointer))
 fmt.Println(flnum)
}
输出结果：

0xc4200140b0
```
* 在上述代码中，我们小加改动。通过 unsafe.Pointer 的特性对该指针变量进行了修改，就可以完成任意类型（*T）的指针转换。

* 需要注意的是，这时还无法对变量进行操作或访问，因为不知道该指针地址指向的东西具体是什么类型。不知道是什么类型，又如何进行解析呢？

* 无法解析也就自然无法对其变更了。

unsafe.Offsetof
在上小节中，我们对普通的指针变量进行了修改。那么它是否能做更复杂一点的事呢？
```
type Num struct{
 i string
 j int64
}
 
func main(){
 n := Num{i: "EDDYCJY", j: 1}
 nPointer := unsafe.Pointer(&n)
 
 niPointer := (*string)(unsafe.Pointer(nPointer))
 *niPointer = "煎鱼"
 
 njPointer := (*int64)(unsafe.Pointer(uintptr(nPointer) + unsafe.Offsetof(n.j)))
 *njPointer = 2
 
 fmt.Printf("n.i: %s, n.j: %d", n.i, n.j)
}
 
输出结果：

n.i: 煎鱼, n.j: 2
```
* 在剖析这段代码做了什么事之前，我们需要了解结构体的一些基本概念：

* 结构体的成员变量在内存存储上是一段连续的内存。

* 结构体的初始地址就是第一个成员变量的内存地址。

* 基于结构体的成员地址去计算偏移量。就能够得出其他成员变量的内存地址。

* 再回来看看上述代码，得出执行流程：

* 修改 n.i 值：i 为第一个成员变量。因此不需要进行偏移量计算，直接取出指针后转换为 Pointer，再强制转换为字符串类型的指针值即可。

* 修改 n.j 值：j 为第二个成员变量。需要进行偏移量计算，才可以对其内存地址进行修改。在进行了偏移运算后，当前地址已经指向第二个成员变量。接着重复转换赋值即可。

### 细节分析
* 需要注意的是，这里使用了如下方法（来完成偏移计算的目标）：

* 1、uintptr：uintptr 是 Go 的内置类型。返回无符号整数，可存储一个完整的地址。后续常用于指针运算

* type uintptr uintptr
* 2、unsafe.Offsetof：返回成员变量 x 在结构体当中的偏移量。更具体的讲，就是返回结构体初始位置到 x 之间的字节数。需要注意的是入参 ArbitraryType 表示任意类型，并非定义的 int。它实际作用是一个占位符

* func Offsetof(x ArbitraryType) uintptr
在这一部分，其实就是巧用了 Pointer 的第三、第四点特性。这时候就已经可以对变量进行操作了。

糟糕的例子
```
func main(){
 n := Num{i: "EDDYCJY", j: 1}
 nPointer := unsafe.Pointer(&n)
    ...
 
 ptr := uintptr(nPointer)
 njPointer := (*int64)(unsafe.Pointer(ptr + unsafe.Offsetof(n.j)))
 ...
}
```
这里存在一个问题，uintptr 类型是不能存储在临时变量中的。因为从 GC 的角度来看，uintptr 类型的临时变量只是一个无符号整数，并不知道它是一个指针地址。

因此当满足一定条件后，ptr 这个临时变量是可能被垃圾回收掉的，那么接下来的内存操作，岂不成迷？

### 总结
* 简洁回顾两个知识点，如下：

* 第一是 unsafe.Pointer 可以让你的变量在不同的指针类型转来转去，也就是表示为任意可寻址的指针类型。

* 第二是 uintptr 常用于与 unsafe.Pointer 打配合，用于做指针运算，巧妙地很。

* 最后还是那句，没有特殊必要的话。是不建议使用 unsafe 标准库，它并不安全。虽然它常常能让你眼前一亮。
