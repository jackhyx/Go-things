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

* 我们主要关注count和buckets这两个字段。count就是指map元素的个数；而buckets是真正存储map的key-value对的地方。
* 这也就可以解释为什么我们一开始那个坑的报错问题。我们用var m map[int]int声明的map，只是分配了一个hmap结构体而已，而buckets这个字段并没有分配内存空间。
* 所以，最后解答我们为什么是引用类型的问题。其实我们传给test函数的值，只是一个hmap结构体；
* 而这个结构体里面又包含了一个bucket数组的指针，也就相当于，表面上我们传了个结构体值过去，而内部却是传了一个指针，这个指针所存储的地址，也就是指针指向的bucket数组结构并没有改变。我们如果对存储key-value对的bucket进行修改，如m[1] = 2这种操作，实际上修改的就是改变了外部函数的bucket值。
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
* 如果overflow bucket数量过多，在get操作时，对这个overflow链表进行遍历的时间复杂度会大大升高，为了避免溢出bucket数量过多，Go语言会在超过某一个阈值的时候，触发扩容操作。
* Go语言bucket的扩容操作也是渐进式的
* Go语言结合了链地址法和开放定址法这两种方案