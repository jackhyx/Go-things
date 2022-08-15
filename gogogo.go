package main

import (
	"errors"
	"fmt"
)

func test() error {
	num1 := 10
	num2 := 0
	if num2 == 0 {
		return errors.New("除数为零")
	}
	res := num1 / num2
	fmt.Println("计算结果为=", res)
	return nil

}

func main() {
	//无论err是否为空,都进行程序的终止

	panic() //内置函数2,参数使一个interface接口
}
