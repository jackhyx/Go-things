* Go语言中type关键字用于定义类型，因此又称为类型别名。

* Go语言中的type并不对应着C/C++语言中的typedef关键字。

* 理解了type关键字就会很容易理解Go语言中的函数、结构体、接口等。

### 类型定义
#### 使用type关键字定义类型

type NewType BaseType

* 类型定义是基于底层类型（BaseType）创建全新的类型（NewType），类型定义创建的是一个全新的类型，与其所基于的类型是两个不同的类型。

* 新类型不会继承底层类型的方法，但会继承底层类型的元素。比如底层类型是interface则其方法会被保留。

* 每个type都拥有其底层类型

* 由于Go语言是强类型静态语言，使用type创建出的新类型虽然与旧类型存储范围相同，但二者并不能相互赋值，只能类型转换。
```
package main

import "fmt"

type interger int

var n interger = interger(100)
//var i int = n // cannot use n (type interger) as type int in assignment
var i int = int(n)

func main() {
	fmt.Printf("n = %v, type = %T\n", n, n)//n = 100, type = main.interger
	fmt.Printf("i = %v, type = %T\n", i, i)//i = 100, type = int
}
/* 类型别名
   使用type关键字定义类型别名

type IntAlias = int

类型定义和类型别名表面上只有一个等号的差异，区别在于类型定义会形成一种新的类型，新类型本身依然具备原始类型的特性。而类型别名只是为类型取别名，别名与原始类型仍旧是同一种类型。


*/
package main

import "fmt"

//定义类型
type NewInt int

//定义类型别名
type IntAlias = int

func main(){
	var x NewInt
	var y IntAlias

	fmt.Printf("x type is %T\n", x)//x type is main.NewInt
	fmt.Printf("y type is %T\n", y)//y type is int
}
// 若为非本地类型（不是在当前包中定义的类型）定义别名，则不能为其定义方法。

package main

import "time"

//为非本地类型定义别名
type Duration = time.Duration

//为类型添加函数
func (d Duration) test(str string){

}

// 编译时报错

cannot define new methods on non-local type time.Duration

// time.Duration并非在当前main包中定义，而是在time包中定义的，与main包并非同一个包，因此不能为非相同包内的类型定义方法。
```
