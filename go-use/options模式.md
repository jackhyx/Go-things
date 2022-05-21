### 函数式选项
* go的函数式选项模式[Functional Options Pattern in Go]
* 前言： go语言开发遇到的许多问题之一是尝试将一个函数的参数设置为可选. 这是一个非常常见的用例, 类似python等语言中的函数默认参数，有些对象应该使用一些基本的默认设置来开箱即用, 你偶尔可能需要提供一些更详细的配置.
#### 实例一
```
 */
package main
// 函数式选项模式（灵活使用默认值，又不影响对元素的修改）

//1、定义一个Options函数类型
type Options func(*Option)
//2、利用闭包为每个字段编写一个设置值的With函数：
func WithA(a,b string)Options  {
	return func(o *Option) {
		o.A = a
		o.B = b
	}
}

func WithC(c int)Options  {
	return func(o *Option) {
		o.C = c
	}
}
// 3、我们定义一个默认的Option
var defaultO  = &Option{ // 指针消耗少
	A: "A",
	B: "B",
	C: 100,
}
// 创造新的构造函数
func NewOptions(opt ...Options) (o *Option ){

	o = defaultO
	for _,O :=range opt{
		O(o)
	}

	return o
}
func main() {
	o := NewOption("a","b",1)
	fmt.Println("o>>:",o)
	//
	o1 := NewOptions()
	fmt.Println("o1:>>",o1)
	o2 := NewOptions(WithA("AA","BB"))
	fmt.Println("o1:>>",o2)

}
// 输出
o>>: &{a b 1}
o1:>> &{A B 100}
o1:>> &{AA BB 100

```
* 总结： 函数式选项模式的本质是利用go对闭包的支持，实现了函数默认值，并且以后再要为Option添加新的字段也不会影响之前的代码

#### 实例二
* 函数式选项（Functional Options）: Golang中实现简洁API的一种方式
* 在使用NewXXX函数构建struct的时候，struct中的属性并不都是必须的，这些非必须属性，在构建struct的过程中可以通过函数式选项的方式，实现更加简洁的API
* 假设需要实现一个协程池GPool，其中必备的属性有协程数量size，还有可选项：是否异步async，错误处理errorHandler，最大缓存任务数maxTaskNums，那么struct的设计应该如下
```
package pool

// Option 定义函数式选项
type Option func(options *Options)

// GPool 协程池
type GPool struct {
size    int64 // 协程数量
options *Options
}

type ErrorHandler func(err error)

// Options 将非必须的选项都放到这里
type Options struct {
async    bool         // 是否支持异步提交任务
handler  ErrorHandler // 任务执行出错时，回调该函数
maxTasks int64        // 协程池所接受的最大缓存任务数
}

// NewGPool 新建协程池
func NewGPool(size int64, opts ...Option) *GPool {
options := loadOpts(opts)
return &GPool{
size:    size,
options: options,
}
}

func loadOpts(opts []Option) *Options {
options := &Options{}
for _, opt := range opts {
opt(options)
}
return options
}

func WithAsync(async bool) Option {
return func(options *Options) {
options.async = async
}
}

func WithErrorHandler(handler ErrorHandler) Option {
return func(options *Options) {
options.handler = handler
}
}

func WithMaxTasks(maxTasks int64) Option {
return func(options *Options) {
options.maxTasks = maxTasks
}
}

// 如果需要创建一个协程池，协程数量为100，只需要这样写
p := pool.NewGPool(100)

// 如果需要创建一个协程池，协程数量为100并支持异步提交，只需要这样写
p := pool.NewGPool(100, pool.WithAsync(true))

// 如果需要穿件一个协程池，协程数量为100、支持异步提交，并且回调自定义错误处理，只需要这样写
p := pool.NewGPool(100,
pool.WithAsync(true),
pool.WithErrorHandler(func(err error) {
// 处理任务执行过程中发生的error
}),
)


// 如果不使用函数式选项:第一种，直接构建struct，但是需要填写非常非常多的属性，对调用者并不友好
func NewGPool(size int64, async bool, handler ErrorHandler, maxTasks int64) *GPool {
return &GPool{
size:    size,
options: &Options{
async:    async,
handler:  handler,
maxTasks: maxTasks,
},
}
}

// 当struct中的属性变得越来越多时候，这长长的函数签名，对于调用者而言，简直是噩梦般的存在
```
* 第二种，使用建造者模式
```
func (builder *GPoolBuilder) Builder(size int64) *GPoolBuilder {
return &GPoolBuilder{p: &GPool{
size: size,
options: &Options{},
}}
}

func (builder *GPoolBuilder) WithAsync(async bool) *GPoolBuilder {
builder.p.options.async = async
return builder
}

func (builder *GPoolBuilder) Build() *GPool {
return builder.p
}

// 调用者使用经构建者模式封装后的API，还是非常舒服的
builder := GPoolBuilder{}
p := builder.Builder(100).WithAsync(true).Build()
```

* 但是，却要额外维护一份属于builder的代码，虽然使用简洁，但是具备一定的维护成本！！
* 总的来看，函数式选项还是最佳的选择方案，开发者通过它能够构建简洁，友好的AP


#### 实例三
前置条件：我们从一个比较常见的问题入手，当我们有一个结构体，我们想到得到一个该结构体的变量，一个常见的方式就是通过工厂模式创建一个new函数，然后传入相应的参数。那么就有了下面的代码：
```
type Diss struct {
	Topic string
	Person string
	Time int
}

// 优点: 直观简单
// 缺点: 当我们需要修改Diss结构体里面的内容，我们同时也需要修改这个传统的new函数
func traditionalNew(topic, person string, time int) Diss {
	return Diss{
		Topic:  topic,
		Person: person,
		Time:   time,
	}
}
```
* 假设我们不想传入时间参数，这个时候该怎么去解决呢？一种方式是直接通过构造结构体来实现：

```diss := Diss{"some topic", "someone"}


这里我们并不传入Time字段的值，直接使用他的默认值
但是如果我们使用上面的传统方式的new的时候就无法成功满足这个要求，我们只能通过新建另一个传统的new2来实现。这样就会比较麻烦。
```
* 这里我们就引入option编程模式：具体来说要知道golang提供的基本的编程方式：可变参数，函数式编程。

* 可变参数示例
```
func variedFunc(arg ...int){
    //这里的arg就是可变参数，我们可以对他进行遍历
}
```
* 函数式编程：
```
type func01 func(int)

func funcationalFunc() func01 {
	return func(int){

	}
}
```
* 有了这些基础，我们来构想一个模型出来
```
type Option func(d *Diss)

// 我们现在只需要传入Option变量就可以了，所以我们需要构造函数

func optionNew(option ...Option) *Diss{
	diss := &Diss{}
	for _, o := range option {
		o(diss)
	}
	return diss
}

func WithTopic(topic string) Option {
	return func(d *Diss) {
		d.Topic = topic
	}
}

func WithPerson(person string) Option {
	return func(d *Diss) {
		d.Person = person
	}
}

func WithTime(time int) Option {
	return func(d *Diss) {
		d.Time = time
	}
}

func main() {
	dissByOption := optionNew(WithTopic("something bad"), WithTime(2), WithPerson("funk"))
	fmt.Println(dissByOption)
	dissRaw := traditionalNew("something bad too","jackal",3)
	fmt.Println(dissRaw)
}
```
* 这里可以把option看成对diss对象的装饰器。然后我们通过函数返回的方式生成这些装饰器，当然你也可以使用其他方式实现类似的装饰器的效果。这里还有一个优势就是，我们传参的时候可以无序地传进去，而在传统的构造函数里面只能通过参数指定的方式传参。
* 关于option的用法在很多源码库里面都是用到了，学习它对于读源码能力有很大的提高，对写代码能力也有不小的促进。
