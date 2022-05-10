## 总结
### 在本节中，我们介绍了很多内容。我们为一个字典应用创建了完整的 CRUD API。在整个过程中，我们学会了如何：
* 创建 map

* 在 map 中搜索值
```go

const (
ErrNotFound         = DictionaryErr("could not find the word you were looking for")
ErrWordExists       = DictionaryErr("cannot add word because it already exists")
ErrWordDoesNotExist = DictionaryErr("cannot update word because it does not exist")
)

type DictionaryErr string

func (e DictionaryErr) Error() string {
return string(e)
}

```
* 向 map 添加新值
```go

func (d Dictionary) Add(word, definition string) error {
_, err := d.Search(word)

switch err {
case ErrNotFound:
d[word] = definition
case nil:
return ErrWordExists
default:
return err
}

return nil
}

```
* 更新 map 中的值
```go

func (d Dictionary) Update(word, definition string) error {
_, err := d.Search(word)

switch err {
case ErrNotFound:
return ErrWordDoesNotExist
case nil:
d[word] = definition
default:
return err
}

return nil
}

```

* 从 map 中删除值
```go

func (d Dictionary) Delete(word string) {
delete(d, word)
}

```

* 了解更多错误相关的知识
* 如何创建常量类型的错误

```go

const (
ErrNotFound         = DictionaryErr("could not find the word you were looking for")
ErrWordExists       = DictionaryErr("cannot add word because it already exists")
ErrWordDoesNotExist = DictionaryErr("cannot update word because it does not exist")
)

type DictionaryErr string

func (e DictionaryErr) Error() string {
return string(e)
}

```

* 对错误进行封装