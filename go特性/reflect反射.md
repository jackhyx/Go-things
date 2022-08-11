### Golang的反射reflect

* 静态类型：每个变量都有一个静态类型，这个类型是在编译时（compile time）就已知且固定的。
* 动态类型：接口类型的变量还有一个动态类型，是在运行时（run time）分配给变量的值的一个非接口类型。（除非分配给变量的值是nil，因为nil没有类型）。
* 空接口类型interface{} (别名any)表示空的方法集，它可以是任何值的类型，因为任何值都满足有0或多个方法（有0个方法一定是任何值的子集）。
* 一个接口类型的变量存储一对内容：分配给变量的具体的值，以及该值的类型描述符。可以示意性地表示为(value, type)对，这里的type是具体的类型，而不是接口类型。


#### reflect的基本功能TypeOf和ValueOf
* 既然反射就是用来检测存储在接口变量内部(值value；类型concrete type) pair对的一种机制。
* 它提供了两种类型（或者说两个方法）让我们可以很容易的访问接口变量内容，分别是reflect.ValueOf() 和 reflect.TypeOf()，看看官方的解释
```
// ValueOf returns a new Value initialized to the concrete value
// stored in the interface i.  ValueOf(nil) returns the zero
func ValueOf(i interface{}) Value {...}
```
* ValueOf用来获取输入参数接口中的数据的值，如果接口为空则返回0

```
// TypeOf returns the reflection Type that represents the dynamic type of i.
// If i is a nil interface value, TypeOf returns nil.
func TypeOf(i interface{}) Type {...}
```
* TypeOf用来动态获取输入参数接口中的值的类型，如果接口为空则返回nil

* reflect.TypeOf()是获取pair中的type，reflect.ValueOf()获取pair中的value，示例如下：
``` 
package main

import (
"fmt"
"reflect"
)

func main() {
var num float64 = 1.2345

	fmt.Println("type: ", reflect.TypeOf(num))
	fmt.Println("value: ", reflect.ValueOf(num))
}

运行结果:
type:  float64
value:  1.2345
```


* reflect.TypeOf： 直接给到了我们想要的type类型，如float64、int、各种pointer、struct 等等真实的类型
* reflect.ValueOf：直接给到了我们想要的具体的值，如1.2345这个具体数值，或者类似&{1 "Allen.Wu" 25} 这样的结构体struct的值
* 也就是说明反射可以将“接口类型变量”转换为“反射类型对象”，反射类型指的是reflect.Type和reflect.Value这两种

##### 从relfect.Value中获取接口interface的信息
* 当执行reflect.ValueOf(interface)之后，就得到了一个类型为”relfect.Value”变量，可以通过它本身的Interface()方法获得接口变量的真实内容，然后可以通过类型判断进行转换，转换为原有真实类型。
* 不过，我们可能是已知原有类型，也有可能是未知原有类型，因此，下面分两种情况进行说明。
* 已知原有类型【进行“强制转换”】
* 已知类型后转换为其对应的类型的做法如下，直接通过Interface方法然后强制转换，如下：
```
realValue := value.Interface().(已知的类型)

package main

import (
"fmt"
"reflect"
)

func main() {
var num float64 = 1.2345

	pointer := reflect.ValueOf(&num)
	value := reflect.ValueOf(num)

	// 可以理解为“强制转换”，但是需要注意的时候，转换的时候，如果转换的类型不完全符合，则直接panic
	// Golang 对类型要求非常严格，类型一定要完全符合
	// 如下两个，一个是*float64，一个是float64，如果弄混，则会panic
	convertPointer := pointer.Interface().(*float64)
	convertValue := value.Interface().(float64)

	fmt.Println(convertPointer)
	fmt.Println(convertValue)
}

运行结果：
0xc42000e238
1.2345
```

* 转换的时候，如果转换的类型不完全符合，则直接panic，类型要求非常严格！
* 转换的时候，要区分是指针还是指
* 也就是说反射可以将“反射类型对象”再重新转换为“接口类型变量”

#### 未知原有类型【遍历探测其Filed】
* 很多情况下，我们可能并不知道其具体类型，那么这个时候，该如何做呢？需要我们进行遍历探测其Filed来得知，示例如下:
```
package main

import (
"fmt"
"reflect"
)

type User struct {
Id   int
Name string
Age  int
}

func (u User) ReflectCallFunc() {
fmt.Println("Allen.Wu ReflectCallFunc")
}

func main() {

	user := User{1, "Allen.Wu", 25}

	DoFiledAndMethod(user)

}

// 通过接口来获取任意参数，然后一一揭晓
func DoFiledAndMethod(input interface{}) {

	getType := reflect.TypeOf(input)
	fmt.Println("get Type is :", getType.Name())

	getValue := reflect.ValueOf(input)
	fmt.Println("get all Fields is:", getValue)

	// 获取方法字段
	// 1. 先获取interface的reflect.Type，然后通过NumField进行遍历
	// 2. 再通过reflect.Type的Field获取其Field
	// 3. 最后通过Field的Interface()得到对应的value
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i).Interface()
		fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
	}

	// 获取方法
	// 1. 先获取interface的reflect.Type，然后通过.NumMethod进行遍历
	for i := 0; i < getType.NumMethod(); i++ {
		m := getType.Method(i)
		fmt.Printf("%s: %v\n", m.Name, m.Type)
	}
}

运行结果：
get Type is : User
get all Fields is: {1 Allen.Wu 25}
Id: int = 1
Name: string = Allen.Wu
Age: int = 25
ReflectCallFunc: func(main.User)
```

* 通过运行结果可以得知获取未知类型的interface的具体变量及其类型的步骤为： 
* 先获取interface的reflect.Type，然后通过NumField进行遍历
* 再通过reflect.Type的Field获取其Field
* 最后通过Field的Interface()得到对应的value

* 通过运行结果可以得知获取未知类型的interface的所属方法（函数）的步骤为：

* 先获取interface的reflect.Type，然后通过NumMethod进行遍历
* 再分别通过reflect.Type的Method获取对应的真实的方法（函数）
* 最后对结果取其Name和Type得知具体的方法名
* 也就是说反射可以将“反射类型对象”再重新转换为“接口类型变量”
* struct 或者 struct 的嵌套都是一样的判断处理方式
```
通过reflect.Value设置实际变量的值
reflect.Value是通过reflect.ValueOf(X)获得的，只有当X是指针的时候，才可以通过reflec.Value修改实际变量X的值，即：要修改反射类型的对象就一定要保证其值是“addressable”的。
示例如下：
package main

import (
"fmt"
"reflect"
)

func main() {

	var num float64 = 1.2345
	fmt.Println("old value of pointer:", num)

	// 通过reflect.ValueOf获取num中的reflect.Value，注意，参数必须是指针才能修改其值
	pointer := reflect.ValueOf(&num)
	newValue := pointer.Elem()

	fmt.Println("type of pointer:", newValue.Type())
	fmt.Println("settability of pointer:", newValue.CanSet())

	// 重新赋值
	newValue.SetFloat(77)
	fmt.Println("new value of pointer:", num)

	////////////////////
	// 如果reflect.ValueOf的参数不是指针，会如何？
	pointer = reflect.ValueOf(num)
	//newValue = pointer.Elem() // 如果非指针，这里直接panic，“panic: reflect: call of reflect.Value.Elem on float64 Value”
}

运行结果：
old value of pointer: 1.2345
type of pointer: float64
settability of pointer: true
new value of pointer: 77
```

* 需要传入的参数是* float64这个指针，然后可以通过pointer.Elem()去获取所指向的Value，注意一定要是指针。
* 如果传入的参数不是指针，而是变量，那么

  通过Elem获取原始值对应的对象则直接panic
  通过CanSet方法查询是否可以设置返回false


* newValue.CantSet()表示是否可以重新设置其值，如果输出的是true则可修改，否则不能修改，修改完之后再进行打印发现真的已经修改了。
* reflect.Value.Elem() 表示获取原始值对应的反射对象，只有原始对象才能修改，当前反射对象是不能修改的
* 也就是说如果要修改反射类型对象，其值必须是“addressable”【对应的要传入的是指针，同时要通过Elem方法获取原始值对应的反射对象】
* struct 或者 struct 的嵌套都是一样的判断处理方式

#### 通过reflect.ValueOf来进行方法的调用
这算是一个高级用法了，前面我们只说到对类型、变量的几种反射的用法，包括如何获取其值、其类型、如果重新设置新值。但是在工程应用中，另外一个常用并且属于高级的用法，就是通过reflect来进行方法【函数】的调用。比如我们要做框架工程的时候，需要可以随意扩展方法，或者说用户可以自定义方法，那么我们通过什么手段来扩展让用户能够自定义呢？关键点在于用户的自定义方法是未可知的，因此我们可以通过reflect来搞定
```
package main

import (
"fmt"
"reflect"
)

type User struct {
Id   int
Name string
Age  int
}

func (u User) ReflectCallFuncHasArgs(name string, age int) {
fmt.Println("ReflectCallFuncHasArgs name: ", name, ", age:", age, "and origal User.Name:", u.Name)
}

func (u User) ReflectCallFuncNoArgs() {
fmt.Println("ReflectCallFuncNoArgs")
}

// 如何通过反射来进行方法的调用？
// 本来可以用u.ReflectCallFuncXXX直接调用的，但是如果要通过反射，那么首先要将方法注册，也就是MethodByName，然后通过反射调动mv.Call

func main() {
user := User{1, "Allen.Wu", 25}

	// 1. 要通过反射来调用起对应的方法，必须要先通过reflect.ValueOf(interface)来获取到reflect.Value，得到“反射类型对象”后才能做下一步处理
	getValue := reflect.ValueOf(user)

	// 一定要指定参数为正确的方法名
	// 2. 先看看带有参数的调用方法
	methodValue := getValue.MethodByName("ReflectCallFuncHasArgs")
	args := []reflect.Value{reflect.ValueOf("wudebao"), reflect.ValueOf(30)}
	methodValue.Call(args)

	// 一定要指定参数为正确的方法名
	// 3. 再看看无参数的调用方法
	methodValue = getValue.MethodByName("ReflectCallFuncNoArgs")
	args = make([]reflect.Value, 0)
	methodValue.Call(args)
}


运行结果：
ReflectCallFuncHasArgs name:  wudebao , age: 30 and origal User.Name: Allen.Wu
ReflectCallFuncNoArgs
```

* 要通过反射来调用起对应的方法，必须要先通过reflect.ValueOf(interface)来获取到reflect.Value，得到“反射类型对象”后才能做下一步处理


* reflect.Value.MethodByName这.MethodByName，需要指定准确真实的方法名字，如果错误将直接panic，MethodByName返回一个函数值对应的reflect.Value方法的名字。


* []reflect.Value，这个是最终需要调用的方法的参数，可以没有或者一个或者多个，根据实际参数来定。

* reflect.Value的 Call 这个方法，这个方法将最终调用真实的方法，参数务必保持一致，如果reflect.Value'Kind不是一个方法，那么将直接panic。


* 本来可以用u.ReflectCallFuncXXX直接调用的，但是如果要通过反射，那么首先要将方法注册，也就是MethodByName，然后通过反射调用methodValue.Call



#### Golang的反射reflect性能
* Golang的反射很慢，这个和它的API设计有关。在 java 里面，我们一般使用反射都是这样来弄的。
```
Field field = clazz.getField("hello");
field.get(obj1);
field.get(obj2);
这个取得的反射对象类型是 java.lang.reflect.Field。它是可以复用的。只要传入不同的obj，就可以取得这个obj上对应的 field。
但是Golang的反射不是这样设计的:
type_ := reflect.TypeOf(obj)
field, _ := type_.FieldByName("hello")
这里取出来的 field 对象是 reflect.StructField 类型，但是它没有办法用来取得对应对象上的值。如果要取值，得用另外一套对object，而不是type的反射
type_ := reflect.ValueOf(obj)
fieldValue := type_.FieldByName("hello")
这里取出来的 fieldValue 类型是 reflect.Value，它是一个具体的值，而不是一个可复用的反射对象了，每次反射都需要malloc这个reflect.Value结构体，并且还涉及到GC。
```

#### Golang reflect慢主要有两个原因


* 涉及到内存分配以及后续的GC；


* reflect实现里面有大量的枚举，也就是for循环，比如类型之类的。


总结
上述详细说明了Golang的反射reflect的各种功能和用法，都附带有相应的示例，相信能够在工程应用中进行相应实践，总结一下就是：


反射可以大大提高程序的灵活性，使得interface{}有更大的发挥余地

反射必须结合interface才玩得转
变量的type要是concrete type的（也就是interface变量）才有反射一说



反射可以将“接口类型变量”转换为“反射类型对象”

反射使用 TypeOf 和 ValueOf 函数从接口中获取目标对象信息



反射可以将“反射类型对象”转换为“接口类型变量

reflect.value.Interface().(已知的类型)
遍历reflect.Type的Field获取其Field



反射可以修改反射类型对象，但是其值必须是“addressable”

想要利用反射修改对象状态，前提是 interface.data 是 settable,即 pointer-interface



通过反射可以“动态”调用方法


因为Golang本身不支持模板，因此在以往需要使用模板的场景下往往就需要使用反射(reflect)来实现



























### 反射是什么
#### 在计算机学中，反射是指计算机程序在运行时（runtime）可以访问、检测和修改它本身状态或行为的一种能力。
#### 用比喻来说，反射就是程序在运行的时候能够 “观察” 并且修改自己的行为（来自维基百科）。
#### 使用


```
v := reflect.ValueOf(&i)
v.Elem().SetFloat(6.66)
log.Println("value: ", i)

```

* 简单来讲就是，应用程序能够在运行时观察到变量的值，并且能够修改他。
```
一个例子
最常见的 reflect 标准库例子，如下：

package main

import (
	"fmt"
	"reflect"
)

func main() {
	rv := []interface{}{"hi", 42, func() {}}
	for _, v := range rv {
		switch v := reflect.ValueOf(v); v.Kind() {
		case reflect.String:
			fmt.Println(v.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fmt.Println(v.Int())
		default:
			fmt.Printf("unhandled kind %s", v.Kind())
		}
	}
}

输出结果：

hi
42
unhandled kind func
```
* 在程序中主要是声明了 rv 变量，变量类型为 interface{}，其包含 3 个不同类型的值，分别是字符串、数字、闭包。

* 而在使用 interface{} 时常见于不知道入参者具体的基本类型是什么，那么我们就会用 interface{} 类型来做一个伪 “泛型”。

#### 此时又会引出一个新的问题，既然入参是 interface{}，那么出参时呢？


* Go 语言是强类型语言，入参是 interface{}，出参也肯定是跑不了的，因此必然离不开类型的判断，这时候就要用到反射，也就是 reflect 标准库。反射过后又再进行 (type) 的类型断言。
* 这就是我们在编写程序时最常遇见的一个反射使用场景。

#### Go reflect
* reflect 标准库中，最核心的莫过于 reflect.Type 和 reflect.Value 类型。而在反射中所使用的方法都围绕着这两者进行，其方法主要含义如下：


* TypeOf 方法：用于提取入参值的类型信息。

* ValueOf 方法：用于提取存储的变量的值信息。

#### reflect.TypeOf
```
演示程序：

func main() {
blog := Blog{"煎鱼"}
typeof := reflect.TypeOf(blog)
fmt.Println(typeof.String())
}
输出结果：

main.Blog
从输出结果中，可得出 reflect.TypeOf 成功解析出 blog 变量的类型是 main.Blog，也就是连 package 都知道了。
```
* 通过人识别的角度来看似乎很正常，但程序就不是这样了。他是怎么知道 “他” 是哪个 package 下的什么呢？

* 我们一起追一下源码看看：
```
func TypeOf(i interface{}) Type {
eface := *(*emptyInterface)(unsafe.Pointer(&i))
return toType(eface.typ)
}
```
从源码层面来看，TypeOf 方法中主要涉及三块操作，分别如下：

* 使用 unsafe.Pointer 方法获取任意类型且可寻址的指针值。 
* 利用 emptyInterface 类型进行强制的 interface 类型转换。 
* 调用 toType 方法转换为可供外部使用的 Type 类型。

而这之中信息量最大的是 emptyInterface 结构体中的 rtype 类型：
```
type rtype struct {
size       uintptr
ptrdata    uintptr
hash       uint32
tflag      tflag
align      uint8  
fieldAlign uint8  
kind       uint8   
equal     func(unsafe.Pointer, unsafe.Pointer) bool
gcdata    *byte  
str       nameOff
ptrToThis typeOff
}
```
在使用上最重要的是 rtype 类型，其实现了 Type 类型的所有接口方法，因此他可以直接作为 Type 类型返回。

而 Type 本质上是一个接口实现，其包含了获取一个类型所必要的所有方法：

type Type interface {
// 适用于所有类型
// 返回该类型内存对齐后所占用的字节数
Align() int

// 仅作用于 strcut 类型
// 返回该类型内存对齐后所占用的字节数
FieldAlign() int

// 返回该类型的方法集中的第 i 个方法
Method(int) Method

// 根据方法名获取对应方法集中的方法
MethodByName(string) (Method, bool)

// 返回该类型的方法集中导出的方法的数量。
NumMethod() int

// 返回该类型的名称
Name() string
}

建议大致过一遍，了解清楚有哪些方法，再针对向看就好。

主体思想是给自己大脑建立一个索引，便于后续快速到 pkg.go.dev 上查询即可。
```
reflect.ValueOf
演示程序：

func main() {
var x float64 = 3.4
fmt.Println("value:", reflect.ValueOf(x))
}
输出结果：

value: 3.4
```
从输出结果中，可得知通过 reflect.ValueOf 成功获取到了变量 x 的值为 3.4。与 reflect.TypeOf 形成一个相匹配，一个负责获取类型，一个负责获取值。

那么 reflect.ValueOf 是怎么获取到值的呢，核心源码如下：
```
func ValueOf(i interface{}) Value {
if i == nil {
return Value{}
}

escapes(i)

return unpackEface(i)
}

func unpackEface(i interface{}) Value {
e := (*emptyInterface)(unsafe.Pointer(&i))
t := e.typ
if t == nil {
return Value{}
}
f := flag(t.Kind())
if ifaceIndir(t) {
f |= flagIndir
}
return Value{t, e.word, f}
}
```
从源码层面来看，ValueOf 方法中主要涉及如下几个操作：

调用 escapes 让变量 i 逃逸到堆上。

将变量 i 强制转换为 emptyInterface 类型。

将所需的信息（其中包含值的具体类型和指针）组装成 reflect.Value 类型后返回。

何时类型转换
在调用 reflect 进行一系列反射行为时，Go 又是在什么时候进行的类型转换呢？

毕竟我们传入的是 float64，而函数如参数是 inetrface 类型。

查看汇编如下:

$ go tool compile -S main.go                         
...
0x0058 00088 ($GOROOT/src/reflect/value.go:2817) LEAQ type.float64(SB), CX
0x005f 00095 ($GOROOT/src/reflect/value.go:2817) MOVQ CX, reflect.dummy+8(SB)
0x0066 00102 ($GOROOT/src/reflect/value.go:2817) PCDATA $0, $-2
0x0066 00102 ($GOROOT/src/reflect/value.go:2817) CMPL runtime.writeBarrier(SB), $0
0x006d 00109 ($GOROOT/src/reflect/value.go:2817) JNE 357
0x0073 00115 ($GOROOT/src/reflect/value.go:2817) MOVQ AX, reflect.dummy+16(SB)
0x007a 00122 ($GOROOT/src/reflect/value.go:2348) PCDATA $0, $-1
0x007a 00122 ($GOROOT/src/reflect/value.go:2348) MOVQ CX, reflect.i+64(SP)
0x007f 00127 ($GOROOT/src/reflect/value.go:2348) MOVQ AX, reflect.i+72(SP)
...
显然，Go 语言会在编译阶段就会完成分析，且进行类型转换。这样子 reflect 真正所使用的就是 interface 类型了。

reflect.Set
演示程序：

func main() {
i := 2.33
v := reflect.ValueOf(&i)
v.Elem().SetFloat(6.66)
log.Println("value: ", i)
}
输出结果：

value:  6.66

从输出结果中，我们可得知在调用 reflect.ValueOf 方法后，我们利用 SetFloat 方法进行了值变更。

核心的方法之一就是 Setter 相关的方法，我们可以一起看看其源码是怎么实现的：

func (v Value) Set(x Value) {
v.mustBeAssignable()
x.mustBeExported() // do not let unexported x leak
var target unsafe.Pointer
if v.kind() == Interface {
target = v.ptr
}
x = x.assignTo("reflect.Set", v.typ, target)
if x.flag&flagIndir != 0 {
typedmemmove(v.typ, v.ptr, x.ptr)
} else {
*(*unsafe.Pointer)(v.ptr) = x.ptr
}
}

* 检查反射对象及其字段是否可以被设置。

* 检查反射对象及其字段是否导出（对外公开）。

* 调用 assignTo 方法创建一个新的反射对象并对原本的反射对象进行覆盖。

* 根据 assignTo 方法所返回的指针值，对当前反射对象的指针进行值的修改。

简单来讲就是，检查是否可以设置，接着创建一个新的对象，最后对其修改。是一个非常标准的赋值流程。

反射三大定律
Go 语言中的反射，其归根究底都是在实现三大定律：

Reflection goes from interface value to reflection object.

Reflection goes from reflection object to interface value.

To modify a reflection object, the value must be settable.

我们将针对这核心的三大定律进行介绍和说明，以此来理解 Go 反射里的各种方法是基于什么理念实现的。

第一定律
反射的第一定律是：“反射可以从接口值（interface）得到反射对象”。
```
示例代码：

func main() {
var x float64 = 3.4
fmt.Println("type:", reflect.TypeOf(x))
}
输出结果：

type: float64

```
可能有读者就迷糊了，我明明在代码中传入的变量 x，他的类型是 float64。怎么就成从接口值得到反射对象了。

其实不然，虽然在代码中我们所传入的变量基本类型是 float64，但是 reflect.TypeOf 方法入参是 interface{}，本质上 Go 语言内部对其是做了类型转换的。这一块会在后面会进一步展开说明。

第二定律
反射的第二定律是：“可以从反射对象得到接口值（interface）”。其与第一条定律是相反的定律，可以是互相补充了。
```
示例代码：

func main() {
vo := reflect.ValueOf(3.4)
vf := vo.Interface().(float64)
log.Println("value:", vf)
}
输出结果：

value: 3.4
```
可以看到在示例代码中，变量 vo 已经是反射对象，然后我们可以利用其所提供的的 Interface 方法获取到接口值（interface），并最后强制转换回我们原始的变量类型。

第三定律
反射的第三定律是：“要修改反射对象，该值必须可以修改”。第三条定律看上去与第一、第二条均无直接关联，但却是必不可少的，因为反射在工程实践中，目的一就是可以获取到值和类型，其二就是要能够修改他的值。

否则反射出来只能看，不能动，就会造成这个反射很鸡肋。例如：应用程序中的配置热更新，必然会涉及配置项相关的变量变动，大多会使用到反射来变动初始值。
```
示例代码：

func main() {
i := 2.33
v := reflect.ValueOf(&i)
v.Elem().SetFloat(6.66)
log.Println("value: ", i)
}
输出结果：

value:  6.66
单从结果来看，变量 i 的值确实从 2.33 变成了 6.66，似乎非常完美。

但是单看代码，似乎有些 “问题”，怎么设置一个反射值这么 ”麻烦“：

为什么必须传入变量 i 的指针引用？

为什么变量 v 在设置前还需要 Elem 一下？

本叛逆的 Gophper 表示我就不这么设置，行不行呢，会不会出现什么问题：

func main() {
i := 2.33
reflect.ValueOf(i).SetFloat(6.66)
log.Println("value: ", i)
}
报错信息：

panic: reflect: reflect.Value.SetFloat using unaddressable value

goroutine 1 [running]:
reflect.flag.mustBeAssignableSlow(0x8e)
/usr/local/Cellar/go/1.15/libexec/src/reflect/value.go:259 +0x138
reflect.flag.mustBeAssignable(...)
/usr/local/Cellar/go/1.15/libexec/src/reflect/value.go:246
reflect.Value.SetFloat(0x10b2980, 0xc00001a0b0, 0x8e, 0x401aa3d70a3d70a4)
/usr/local/Cellar/go/1.15/libexec/src/reflect/value.go:1609 +0x37
main.main()
/Users/eddycjy/go-application/awesomeProject/main.go:10 +0xc5

```


* 根据上述提示可知，由于使用 “使用不可寻址的值”，因此示例程序无法正常的运作下去。并且这是一个 reflect 标准库本身就加以防范了的硬性要求。

* 这么做的原因在于，Go 语言的函数调用的传递都是值拷贝的，因此若不传指针引用，单纯值传递，那么肯定是无法变动反射对象的源值的。因此 Go 标准库就对其进行了逻辑判断，避免出现问题。

* 因此期望变更反射对象的源值时，我们必须主动传入对应变量的指针引用，并且调用 reflect 标准库的 Elem 方法来获取指针所指向的源变量，并且最后调用 Set 相关方法来进行设置。

#### 总结
* 通过本文我们学习并了解了 Go 反射是如何使用，又是基于什么定律设计的。另外我们稍加关注，不难发现 Go 的反射都是基于接口（interface）来实现的，更进一步来讲，Go 语言中运行时的功能很多都是基于接口来实现的。
