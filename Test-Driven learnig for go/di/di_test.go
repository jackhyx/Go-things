package main

import (
	"bytes"
	"testing"
)

/*
func Greet(writer io.Writer, name string) {
	fmt.Fprintf(writer, "Hello, %s", name)
}

func main() {
	Greet(os.Stdout, "Elodie")
}
*/
func TestGreet(t *testing.T) {
	buffer := bytes.Buffer{}
	Greet(&buffer, "Chris")

	got := buffer.String()
	want := "Hello, Chris"

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

// bytes 包中的 buffer 类型实现了 Writer 接口。
// 因此，我们可以在测试中，用它来作为我们的 Writer，接着调用了 Greet 后，我们可以用它来检查写入了什么
