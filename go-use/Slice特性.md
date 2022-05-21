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