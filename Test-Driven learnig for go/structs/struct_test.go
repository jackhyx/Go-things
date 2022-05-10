package main

import (
	"math"
	"testing"
)

type Rectangle struct {
	Width  float64
	Height float64
} //定义我们自己的类型 Rectangle，它可以封装长方形的信息。

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

/*有些编程语言中我们可以这样做即重载

func Area(circle Circle) float64 { ... }
func Area(rectangle Rectangle) float64 { ... }
./shapes.go:20:32: Area redeclared in this block
解决：
不同的包可以有函数名相同的函数。所以我们可以在一个新的包里创建函数 Area(Circle)。但是感觉有点大才小用了
我们可以为新类型定义方法：方法和函数很相似但是方法是通过一个特定类型的实例调用的。函数可以随时被调用，比如 Area(rectangle)。不像方法需要在某个事物上调用。
*/
type Circle struct {
	Radius float64
}

/* 方法声明
当方法被这种类型的变量调用时，数据的引用通过变量 receiverName 获得。在其他许多编程语言中这些被隐藏起来并且通过 this 来获得接收者。
把类型的第一个字母作为接收者变量是 Go 语言的一个惯例。
r Rectangle
*/
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

type Triangle struct {
	Base   float64
	Height float64
}

// 我们可以通过下面的语法来访问一个 struct 中的域： myStruct.field
func (c Triangle) Area() float64 {
	return (c.Base * c.Height) * 0.5
}

type Shape interface {
	Area() float64
}

// 我们的辅助函数是怎样实现不需要关心参数是矩形，圆形还是三角形的。通过声明一个接口，辅助函数能从具体类型解耦而只关心方法本身需要做的工作。
// 接口在 Go 这种静态类型语言中是一种非常强有力的概念。因为接口可以让函数接受不同类型的参数并能创造类型安全且高解耦的代码
// 在 Go 语言中 interface resolution 是隐式的。如果传入的类型匹配接口需要的，则编译正确。
/*
func TestArea(t *testing.T) {

	checkArea := func(t *testing.T, shape Shape, want float64) {
		t.Helper() // 辅助函数
		got := shape.Area()
		if got != want {
			t.Errorf("got %.2f want %.2f", got, want) // f 对应 float64，.2 表示输出 2 位小数
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
		t.Run(tt.name, func(t *testing.T) { // 在每个用例中使用 t.Run，测试用例的错误输出中会包含用例的名字：
			got := tt.shape.Area()
			if got != tt.hasArea {
				t.Errorf("%#v got %.2f want %.2f", tt.shape, got, tt.hasArea)
			}
		})

	}

}
