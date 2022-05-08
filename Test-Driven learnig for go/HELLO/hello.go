package main

import "fmt"

const englishHelloPrefix = "Hello, "
const spanish = "Spanish"
const frenchHelloPrefix = "Bonjour, "
const spanishHelloPrefix = "Hola, "
const french = "French"

func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}
	return greetingPrefix(language) + name
}
func greetingPrefix(language string) (prefix string) {

	switch language {
	case french:
		prefix = frenchHelloPrefix
	case spanish:
		prefix = spanishHelloPrefix
	default:
		prefix = englishHelloPrefix
	}

	return
}

func main() {
	fmt.Println(Hello("Elodie", french))
}

/*  常量应该可以提高应用程序的性能，它避免了每次使用 Hello 时创建 "Hello, " 字符串实例。
显然，对于这个例子来说，性能提升是微不足道的！
但是创建常量的价值是可以快速理解值的含义，有时还可以帮助提高性能
当你有很多 if 语句检查一个特定的值时，通常使用 switch 语句来代替。
如果我们希望稍后添加更多的语言支持，我们可以使用 switch 来重构代码，使代码更易于阅读和扩展

一些新的概念：
在我们的函数签名中，我们使用了 命名返回值（prefix string）。
这将在你的函数中创建一个名为 prefix 的变量。
	它将被分配「零」值。这取决于类型，例如 int 是 0，对于字符串它是 ""。
		你只需调用 return 而不是 return prefix 即可返回所设置的值。
	这将显示在 Go Doc 中，所以它使你的代码更加清晰。
如果没有其他 case 语句匹配，将会执行 default 分支。
函数名称以小写字母开头。在 Go 中，公共函数以大写字母开始，私有函数以小写字母开头。
我们不希望我们算法的内部结构暴露给外部，所以我们将这个功能私有化。
*/
