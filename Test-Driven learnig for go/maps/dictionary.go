package maps

type Dictionary map[string]string

const (
	ErrNotFound         = DictionaryErr("could not find the word you were looking for")
	ErrWordExists       = DictionaryErr("cannot add word because it already exists")
	ErrWordDoesNotExist = DictionaryErr("cannot update word because it does not exist")
)

type DictionaryErr string

func (e DictionaryErr) Error() string {
	return string(e)
}

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

func (d Dictionary) Search(word string) (string, error) {
	definition, ok := d[word]
	if !ok {
		return "", ErrNotFound
	}
	//为了使测试通过，我们使用了一个 map 查找的有趣特性。它可以返回两个值。第二个值是一个布尔值，表示是否成功找到 key。
	//此特性允许我们区分单词不存在还是未定义。

	return definition, nil
}
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

func (d Dictionary) Delete(word string) {
	delete(d, word)
}

/*
Go 的 map 有一个内置函数 delete。它需要两个参数。第一个是这个 map，第二个是要删除的键。
delete 函数不返回任何内容，我们基于相同的概念构建 Delete 方法。由于删除一个不存在的值是没有影响的，与我们的 Update 和 Create 方法不同，我们不需要用错误复杂化 API。
*/
/*  向 map 添加元素也类似于数组。你只需指定键并给它赋一个值。
Map 有一个有趣的特性，不使用指针传递你就可以修改它们。这是因为 map 是引用类型。这意味着它拥有对底层数据结构的引用，就像指针一样。它底层的数据结构是 hash table 或 hash map，你可以在这里阅读有关 hash tables 的更多信息。
Map 作为引用类型是非常好的，因为无论 map 有多大，都只会有一个副本。
引用类型引入了 maps 可以是 nil 值。如果你尝试使用一个 nil 的 map，你会得到一个 nil 指针异常，这将导致程序终止运行。
由于 nil 指针异常，你永远不应该初始化一个空的 map 变量：
var m map[string]string
相反，你可以像我们上面那样初始化空 map，或使用 make 关键字创建 map：


dictionary = map[string]string{}

// OR

dictionary = make(map[string]string)

这两种方法都可以创建一个空的 hash map 并指向 dictionary。这确保永远不会获得 nil 指针异常。
*/
