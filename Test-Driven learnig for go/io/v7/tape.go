package main

import (
	"os"
)

// os.File 文件有一个 truncate 函数，可以让我们有效地清空文件。我们应该能够调用它来得到我们想要的功能
type tape struct {
	file *os.File // file io.ReadWriteSeeker
}

func (t *tape) Write(p []byte) (n int, err error) {
	t.file.Truncate(0)
	t.file.Seek(0, 0)
	return t.file.Write(p)
}

/*
在我们处理文件的过程中有一些非常天真的行为，这可能会在以后产生非常严重的错误。
当我们 Recordwin 时，我们返回到文件的开头，然后写入新的数据，但是如果新的数据比之前的数据要小怎么办?
在我们目前的情况下，这是不可能的。我们从不编辑或删除得分，因此数据只会变得更大，但是这样的代码是不负责任的，出现删除场景的结果是不可想象的。
但是我们要怎么测试这种问题呢？我们需要做的是首先重构我们的代码，这样就可以将我们所编写的数据和正在写入的分开。然后我们可以分别测试它是否以我们期望的方式运行。
我们将创建一个新类型来封装我们的「当写入时，从头部开始」功能。
我把它叫做 Tape。创建一个包含以下内容的新文件



*/
