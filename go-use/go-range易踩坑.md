### 1. 指针数据坑
```
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
 n := make([]*user,0,len(u))
 for _,v := range u{
  n = append(n, &v)
 }
 fmt.Println(n)
 for _,v := range n{
  fmt.Println(v)
 }
}
//
[0xc0000a6040 0xc0000a6040 0xc0000a6040]
&{asong2020 18}
&{asong2020 18}
&{asong2020 18}
```
* 在for range中，变量v是用来保存迭代切片所得的值，因为v只被声明了一次，每次迭代的值都是赋值给v，该变量的内存地址始终未变，这样讲他的地址追加到新的切片中，该切片保存的都是同一个地址
* 变量v的地址也并不是指向原来切片u[2]的，因我在使用range迭代的时候，变量v的数据是切片的拷贝数据，所以直接copy了结构体数据。
#### 解决一：
```package main

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
n := make([]*user,0,len(u))
for _,v := range u{
o := v // 引入了中间变量
n = append(n, &o)
}
fmt.Println(n)
for _,v := range n{
fmt.Println(v)
}
}
```
#### 解决二
```
......略
for k,_ := range u{
  n = append(n, &u[k])
 }
......略
```

### 2. 迭代修改变量问题
```
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
//
[{asong 23} {song 19} {asong2020 18}]

```
#### 解决
```
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
 for k,v := range u{
  if v.age != 18{
   u[k].age = 18
  }
 }
 fmt.Println(u)
}
```
### 3. 是否会造成死循环
```
func main() {
 v := []int{1, 2, 3}
 for i := range v {
  v = append(v, i)
 }
}//no
```

### map 的delete&add&for range
```
delete
func main()  {
 d := map[string]string{
  "asong": "帅",
  "song": "太帅了",
 }
 for k := range d{
  if k == "asong"{
   delete(d,k)
  }
 }
 fmt.Println(d)
}
 
// 运行结果:
map[song:太帅了]
add
func main()  {
 var addTomap = func() {
  var t = map[string]string{
   "asong": "太帅",
   "song": "好帅",
   "asong1": "非常帅",
  }
  for k := range t {
   t["song2020"] = "真帅"
   fmt.Printf("%s%s ", k, t[k])
  }
 }
 for i := 0; i < 10; i++ {
  addTomap()
  fmt.Println()
 }
}
// map内部实现是一个链式hash表，为了保证无顺序，初始化时会随机一个遍历开始的位置，所以新增的元素被遍历到就变的不确定了，同样删除也是一个道理，但是删除元素后边就不会出现，所以一定不会被遍历到。
```
