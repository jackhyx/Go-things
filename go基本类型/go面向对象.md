### Go 实现面向对象编程
封装--通过首字母大小写来控制
继承--组合
多态--接口
#### 封装
* 面向对象中的 “封装” 指的是可以隐藏对象的内部属性和实现细节，仅对外提供公开接口调用，这样子用户就不需要关注你内部是怎么实现的。
* 在 Go 语言中的属性访问权限，通过首字母大小写来控制：
* 首字母大写，代表是公共的、可被外部访问的。
* 首字母小写，代表是私有的，不可以被外部访问。
```
Go 语言的例子如下：

type Animal struct {
name string
}

func NewAnimal() *Animal {
return &Animal{}
}

func (p *Animal) SetName(name string) {
p.name = name
}

func (p *Animal) GetName() string {
return p.name
}
```
* 在上述例子中，我们声明了一个结构体 Animal，其属性 name 为小写。没法通过外部方法，在配套上存在 Setter 和 Getter 的方法，用于统一的访问和设置控制。
* 以此实现在 Go 语言中的基本封装。

### 继承
#### 面向对象中的 “继承” 指的是子类继承父类的特征和行为，使得子类对象（实例）具有父类的实例域和方法，或子类从父类继承方法，使得子类具有父类相同的行为。
* 在 Go 语言中，是没有类似 extends 关键字的这种继承的方式，在语言设计上采取的是组合的方式：
```
type Animal struct {
Name string
}

type Cat struct {
Animal
FeatureA string
}

type Dog struct {
Animal
FeatureB string
}
在上述例子中，我们声明了 Cat 和 Dog 结构体，其在内部匿名组合了 Animal 结构体。因此 Cat 和 Dog 的实例都可以调用 Animal 结构体的方法：

func main() {
p := NewAnimal()
p.SetName("煎鱼，记得点赞~")

dog := Dog{Animal: *p}
fmt.Println(dog.GetName())
}
同时 Cat 和 Dog 的实例可以拥有自己的方法：

func (dog *Dog) HelloWorld() {
fmt.Println("脑子进煎鱼了")
}

func (cat *Cat) HelloWorld() {
fmt.Println("煎鱼进脑子了")
}
```
* 上述例子能够正常包含调用 Animal 的相关属性和方法，也能够拥有自己的独立属性和方法，在 Go 语言中达到了类似继承的效果。

### 多态
#### 面向对象中的 “多态” 指的同一个行为具有多种不同表现形式或形态的能力，具体是指一个类实例（对象）的相同方法在不同情形有不同表现形式。
* 多态也使得不同内部结构的对象可以共享相同的外部接口，也就是都是一套外部模板，内部实际是什么，只要符合规格就可以。
```
在 Go 语言中，多态是通过接口来实现的：

type AnimalSounder interface {
MakeDNA()
}

func MakeSomeDNA(animalSounder AnimalSounder) {
animalSounder.MakeDNA()
}
在上述例子中，我们声明了一个接口类型 AnimalSounder，配套一个 MakeSomeDNA 方法，其接受 AnimalSounder 接口类型作为入参。

因此在 Go 语言中。只要配套的 Cat 和 Dog 的实例也实现了 MakeSomeDNA 方法，那么我们就可以认为他是 AnimalSounder 接口类型：

type AnimalSounder interface {
MakeDNA()
}

func MakeSomeDNA(animalSounder AnimalSounder) {
animalSounder.MakeDNA()
}

func (c *Cat) MakeDNA() {
fmt.Println("煎鱼是煎鱼")
}

func (c *Dog) MakeDNA() {
fmt.Println("煎鱼其实不是煎鱼")
}

func main() {
MakeSomeDNA(&Cat{})
MakeSomeDNA(&Dog{})
}
```
* 当 Cat 和 Dog 的实例实现了 AnimalSounder 接口类型的约束后，就意味着满足了条件，他们在 Go 语言中就是一个东西。能够作为入参传入 MakeSomeDNA 方法中，再根据不同的实例实现多态行为。

### 总结
通过今天这篇文章，我们基本了解了面向对象的定义和 Go 官方对面向对象这一件事的看法，同时针对面向对象的三大特性：“封装、继承、多态” 在 Go 语言中的实现方法就进行了一一讲解。

五大原则 “单一职责原则（SRP）、开放封闭原则（OCP）、里氏替换原则（LSP）、依赖倒置原则（DIP）、接口隔离原则（ISP）” 
