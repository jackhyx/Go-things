package 分布式相关

import (
	"fmt"
	"hash/crc32"
	"sort"
	"sync"
)

type ConsistenceHash struct {
	nodesMap        map[uint32]string // hash slot和虚拟node的映射关系
	nodesSlots      slots             // 虚拟node所有hash slot组成的切片
	NumVirtualNodes int               // 为每台机器在hash圆环上创建多少个虚拟Node
	mu              sync.RWMutex
}

// 使用sort.Sort函数，传入的参数需要实现的接口
type slots []uint32

func (s slots) Len() int {
	return len(s)
}

func (s slots) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s slots) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// 由于hash圆环上有2^32个hash slot，使用uint32类型来标识hash slot。为了能够使用sort.Sort函数排序对虚拟Node的hash slot进行排序，自定义一个[]uint32的类型实现sort.Sort函数约定的参数接口。
// 通过crc32函数计算散列值
func (h *ConsistenceHash) hash(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

// 集群中增加机器

func (h *ConsistenceHash) AddNode(addr string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	// 根据定义的数量生成虚拟Node
	// addr加上不同的后缀计算散列值得到每个虚拟Node的hash slot
	// 同一个机器的所有hash slot最终都指向同一个ip/port
	for i := 0; i < h.NumVirtualNodes; i++ {
		slot := h.hash(fmt.Sprintf("%s%d", addr, i))
		h.nodesMap[slot] = addr
	}
	h.sortNodesSlots()
}

// 所有虚拟Node映射到的hash slot排序后保存到切片
func (h *ConsistenceHash) sortNodesSlots() {
	slots := h.nodesSlots[:]
	for slot := range h.nodesMap {
		slots = append(slots, slot)
	}
	sort.Sort(slots)
	h.nodesSlots = slots
}

// 从集群中摘除机器

func (h *ConsistenceHash) DeleteNode(addr string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	// 删除所有的虚拟节点
	for i := 0; i < h.NumVirtualNodes; i++ {
		slot := h.hash(fmt.Sprintf("%s%d", addr, i))
		delete(h.nodesMap, slot)
	}
	h.sortNodesSlots()
}
