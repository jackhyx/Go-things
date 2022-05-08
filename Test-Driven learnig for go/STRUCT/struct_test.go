package main

import (
	"math"
	"testing"
)

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

type Triangle struct {
	Base   float64
	Height float64
}

func (c Triangle) Area() float64 {
	return (c.Base * c.Height) * 0.5
}

type Shape interface {
	Area() float64
}

/*
func TestArea(t *testing.T) {

	checkArea := func(t *testing.T, shape Shape, want float64) {
		t.Helper()
		got := shape.Area()
		if got != want {
			t.Errorf("got %.2f want %.2f", got, want)
		}
	}

	t.Run("rectangles", func(t *testing.T) {
		rectangle := Rectangle{12, 6}
		checkArea(t, rectangle, 72.0)
	})

	t.Run("circles", func(t *testing.T) {
		circle := Circle{10}
		checkArea(t, circle, 314.1592653589793)
	})

}
*/
func TestArea(t *testing.T) {

	areaTests := []struct {
		name    string
		shape   Shape
		hasArea float64
	}{
		{name: "Rectangle", shape: Rectangle{Width: 12, Height: 6}, hasArea: 72.0},
		{name: "Circle", shape: Circle{Radius: 10}, hasArea: 314.1592653589793},
		{name: "Triangle", shape: Triangle{Base: 12, Height: 6}, hasArea: 36.0},
	}

	for _, tt := range areaTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			got := tt.shape.Area()
			if got != tt.hasArea {
				t.Errorf("%#v got %.2f want %.2f", tt.shape, got, tt.hasArea)
			}
		})

	}

}

/*
声明方法的语法跟函数差不多，因为他们本身就很相似。唯一的不同是方法接收者的语法 func(receiverName ReceiverType) MethodName(args)。
当方法被这种类型的变量调用时，数据的引用通过变量 receiverName 获得。在其他许多编程语言中这些被隐藏起来并且通过 this 来获得接收者。
把类型的第一个字母作为接收者变量是 Go 语言的一个惯例

这种定义 interface 的方式与大部分其他编程语言不同。通常接口定义需要这样的代码 My type Foo implements interface Bar。
但是在我们的例子里，
Rectangle 有一个返回值类型为 float64 的方法 Area，所以它满足接口 Shape
Circle 有一个返回值类型为 float64 的方法 Area，所以它满足接口 Shape
string 没有这种方法，所以它不满足这个接口
等等
在 Go 语言中 interface resolution 是隐式的。如果传入的类型匹配接口需要的，则编译正确

这些数字代表什么并不一目了然，我们应该让我们的测试函数更容易理解。
到目前为止我们仅仅学到一种创建结构体 MyStruct{val1, val2} 的方法，但是我们可以选择命名这些域

在每个用例中使用 t.Run，测试用例的错误输出中会包含用例的名字：

总结
声明结构体以创建我们自己的类型，让我们把数据集合在一起并达到简化代码的目地
声明接口，这样我们可以定义适合不同参数类型的函数（参数多态）
在自己的数据类型中添加方法以实现接口
列表驱动测试让断言更清晰，这样可以使测试文件更易于扩展和维护
*/
