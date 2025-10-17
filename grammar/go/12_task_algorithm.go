package main

import (
	"fmt"
	"strings"
)

// ### 任务9：数据结构库
// **目标**：实现常用数据结构
// **描述**：实现一个数据结构库，包含栈、队列、链表、树、图等

// **包含结构**：
// - 栈（Stack）
// - 队列（Queue）
// - 链表（LinkedList）
// - 二叉搜索树（BST）
// - 哈希表（HashTable）
// - 图（Graph）

type Stack struct {
	items []interface{}
}
func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
}
func (s *Stack) Pop() interface{} {
	if s.IsEmpty() {
		return nil
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}
func (s *Stack) Peek() interface{} {
	if s.IsEmpty() {
		return nil
	}
	return s.items[len(s.items)-1]
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}
func (s *Stack) Size() int {
	return len(s.items)
}
func (s *Stack) Clear() {
	s.items = []interface{}{}
}


type Queue struct {
	items []interface{}
}
func (q *Queue) Enqueue(item interface{}) {
	q.items = append(q.items, item)
}
func (q *Queue) Dequeue() interface{} {
	if q.IsEmpty() {
		return nil
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}
func (q *Queue) Peek() interface{} {
	if q.IsEmpty() {
		return nil
	}
	return q.items[0]
}
func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}
func (q *Queue) Size() int {
	return len(q.items)
}
func (q *Queue) Clear() {
	q.items = []interface{}{}
}

type LinkedList struct {
	head *Node
	tail *Node
}
type Node struct {
	value interface{}
	next  *Node
	prev  *Node
}
func (ll *LinkedList) Add(value interface{}) {
	node := &Node{value: value}
	if ll.head == nil {
		ll.head = node
		ll.tail = node
	} else {
		ll.tail.next = node
		node.prev = ll.tail
		ll.tail = node
	}
}
func (ll *LinkedList) Remove(value interface{}) {
	if ll.head == nil {
		return
	}
	if ll.head.value == value {
		ll.head = ll.head.next
		if ll.head != nil {
			ll.head.prev = nil
		} else {
			ll.tail = nil
		}
		return
	}
	if ll.tail.value == value {
		ll.tail = ll.tail.prev
		if ll.tail != nil {
			ll.tail.next = nil
		} else {
			ll.head = nil
		}
		return
	}
	for node := ll.head.next; node != nil; node = node.next {
		if node.value == value {
			node.prev.next = node.next
			if node.next != nil {
				node.next.prev = node.prev
			}
			return
		}
	}
}

func (ll *LinkedList) Contains(value interface{}) bool {
	for node := ll.head; node != nil; node = node.next {
		if node.value == value {
			return true
		}
	}
	return false
}
func (ll *LinkedList) Size() int{
	if ll.head == nil {
		return 0;
	}
	size := 0;
	for node := ll.head; node != nil; node = node.next {
		size ++
	}
	return size;
}

// 泛型版本的BST，支持任何可比较的类型
type BST[T comparable] struct {
	root *TreeNode[T]
	compare func(a, b T) int // 自定义比较函数
}

type TreeNode[T comparable] struct {
	value T
	left  *TreeNode[T]
	right *TreeNode[T]
}

// 泛型BST的方法
func NewBST[T comparable](compare func(a, b T) int) *BST[T] {
	return &BST[T]{compare: compare}
}

func (bst *BST[T]) Insert(value T) {
	if bst.root == nil {
		bst.root = &TreeNode[T]{value: value}
	} else {
		bst.root.insert(value, bst.compare)
	}
}

func (node *TreeNode[T]) insert(value T, compare func(a, b T) int) {
	cmp := compare(value, node.value)
	if cmp < 0 {
		if node.left == nil {
			node.left = &TreeNode[T]{value: value}
		} else {
			node.left.insert(value, compare)
		}
	} else if cmp > 0 {
		if node.right == nil {
			node.right = &TreeNode[T]{value: value}
		} else {
			node.right.insert(value, compare)
		}
	}
	// cmp == 0 表示相等，不插入重复值
}

func (bst *BST[T]) Search(value T) bool {
	return bst.root.search(value, bst.compare)
}

func (node *TreeNode[T]) search(value T, compare func(a, b T) int) bool {
	if node == nil {
		return false
	}
	cmp := compare(value, node.value)
	if cmp == 0 {
		return true
	} else if cmp < 0 {
		return node.left.search(value, compare)
	} else {
		return node.right.search(value, compare)
	}
}

// Iterate 方法用于泛型BST
func (bst *BST[T]) Iterate() []T {
	var values []T
	bst.root.iterate(&values)
	return values
}

func (node *TreeNode[T]) iterate(values *[]T) {
	if node == nil {
		return
	}
	node.left.iterate(values)
	*values = append(*values, node.value)
	node.right.iterate(values)
}
type HashTable[T comparable,V any] struct {
	entries map[T]V
}
func NewHashTable[T comparable,V any]() *HashTable[T,V] {
	return &HashTable[T,V]{entries: make(map[T]V)}
}
func (ht *HashTable[T,V]) Set(key T, value V) {
	ht.entries[key] = value
}
func (ht *HashTable[T,V]) Get(key T) (V, bool) {
	value, ok := ht.entries[key]
	return value, ok
}
func (ht *HashTable[T,V]) Delete(key T) {
	delete(ht.entries, key)
}
func (ht *HashTable[T,V]) Size() int {
	return len(ht.entries)
}
func (ht *HashTable[T,V]) Clear() {
	ht.entries = make(map[T]V)
}
func (ht *HashTable[T,V]) Keys() []T {
	keys := make([]T, 0, len(ht.entries))
	for key := range ht.entries {
		keys = append(keys, key)
	}
	return keys
}

func (ht *HashTable[T,V]) Values() []V {
	values := make([]V, 0, len(ht.entries))
	for _, value := range ht.entries {
		values = append(values, value)
	}
	return values
}
func (ht *HashTable[T,V]) Contain(key T) bool{
   contain := false
   for k := range ht.entries{
	if key == k{
		contain = true
		break
	}
   }
   return contain
}
type Graph[T comparable] struct {
	vertices map[T]*Vertex[T]
}
type Vertex[T comparable] struct {
	value T
	edges []*Edge[T]
}
type Edge[T comparable] struct {
	to     *Vertex[T]
	weight int
}
func (g *Graph[T]) AddVertex(value T) {
	if _, ok := g.vertices[value]; !ok {
		g.vertices[value] = &Vertex[T]{value: value}
	}
}
func (g *Graph[T]) AddEdge(from, to T, weight int) {
	fromVertex, ok := g.vertices[from]
	if !ok {
		g.AddVertex(from)
		fromVertex = g.vertices[from]
	}
	toVertex, ok := g.vertices[to]
	if !ok {
		g.AddVertex(to)
		toVertex = g.vertices[to]
	}
	fromVertex.edges = append(fromVertex.edges, &Edge[T]{to: toVertex, weight: weight})
}


func NewGraph[T comparable]() *Graph[T] {
	return &Graph[T]{vertices: make(map[T]*Vertex[T])}
}

func (g *Graph[T]) DFS(start T, visited map[T]bool, result *[]T) {
	if visited == nil {
		visited = make(map[T]bool)
	}
	
	visited[start] = true
	*result = append(*result, start)
	
	if vertex, ok := g.vertices[start]; ok {
		for _, edge := range vertex.edges {
			if !visited[edge.to.value] {
				g.DFS(edge.to.value, visited, result)
			}
		}
	}
}

func (g *Graph[T]) BFS(start T) []T {
	var result []T
	visited := make(map[T]bool)
	queue := []T{start}
	
	visited[start] = true
	
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)
		
		if vertex, ok := g.vertices[current]; ok {
			for _, edge := range vertex.edges {
				if !visited[edge.to.value] {
					visited[edge.to.value] = true
					queue = append(queue, edge.to.value)
				}
			}
		}
	}
	
	return result
}


func (g *Graph[T]) GetNeighbors(vertex T) []T {
	var neighbors []T
	if v, ok := g.vertices[vertex]; ok {
		for _, edge := range v.edges {
			neighbors = append(neighbors, edge.to.value)
		}
	}
	return neighbors
}

func (g *Graph[T]) HasPath(from, to T) bool {
	visited := make(map[T]bool)
	return g.hasPathDFS(from, to, visited)
}

func (g *Graph[T]) hasPathDFS(current, target T, visited map[T]bool) bool {
	if current == target {
		return true
	}
	
	visited[current] = true
	
	if vertex, ok := g.vertices[current]; ok {
		for _, edge := range vertex.edges {
			if !visited[edge.to.value] {
				if g.hasPathDFS(edge.to.value, target, visited) {
					return true
				}
			}
		}
	}
	
	return false
}


func main(){
	fmt.Println("=== 泛型BST演示开始 ===")
	bst := NewBST(func(a, b string) int { return strings.Compare(a, b) })
	bst.Insert("banana")
	bst.Insert("apple")
	bst.Insert("orange")
	bst.Insert("peach")
	bst.Insert("grape")
	bst.Insert("watermelon")
	bst.Insert("mango")
	fmt.Println("BST中序遍历结果:", bst.Iterate())
	
	fmt.Println("\n=== 图数据结构演示开始 ===")
	// 创建一个字符串类型的图
	graph := NewGraph[string]()
	
	// 添加顶点和边
	graph.AddVertex("A")
	graph.AddVertex("B")
	graph.AddVertex("C")
	graph.AddVertex("D")
	graph.AddVertex("E")
	
	// 添加边 (有向图)
	graph.AddEdge("A", "B", 1)
	graph.AddEdge("A", "C", 1)
	graph.AddEdge("B", "D", 1)
	graph.AddEdge("C", "D", 1)
	graph.AddEdge("D", "E", 1)
	graph.AddEdge("B", "E", 1)
	
	// 深度优先搜索
	var dfsResult []string
	graph.DFS("A", nil, &dfsResult)
	fmt.Println("从A开始的DFS遍历:", dfsResult)
	
	// 广度优先搜索
	bfsResult := graph.BFS("A")
	fmt.Println("从A开始的BFS遍历:", bfsResult)
	
	// 获取邻居
	neighbors := graph.GetNeighbors("A")
	fmt.Println("A的邻居:", neighbors)
	
	// 检查路径
	hasPath := graph.HasPath("A", "E")
	fmt.Println("A到E是否有路径:", hasPath)
	
	hasPath2 := graph.HasPath("E", "A")
	fmt.Println("E到A是否有路径:", hasPath2)
	
	fmt.Println("\n=== 字符串匹配算法演示 ===")
	text := "ABABCABABA"
	pattern := "ABABA"
	
	// KMP算法
	kmpResult := KMPSearch(text, pattern)
	fmt.Printf("KMP算法: 在文本'%s'中查找模式'%s'，位置: %d\n", text, pattern, kmpResult)
	
	// Boyer-Moore算法
	bmResult := BoyerMooreSearch(text, pattern)
	fmt.Printf("Boyer-Moore算法: 在文本'%s'中查找模式'%s'，位置: %d\n", text, pattern, bmResult)
	
	// 测试不存在的模式
	notFoundPattern := "XYZ"
	kmpNotFound := KMPSearch(text, notFoundPattern)
	bmNotFound := BoyerMooreSearch(text, notFoundPattern)
	fmt.Printf("查找不存在的模式'%s': KMP=%d, Boyer-Moore=%d\n", notFoundPattern, kmpNotFound, bmNotFound)
}

// ### 任务10：排序和搜索算法
// **目标**：实现经典算法
// **描述**：实现各种排序和搜索算法，并进行性能比较

// **包含算法**：
// - 排序：冒泡、选择、插入、快排、归并、堆排序
// - 搜索：二分搜索、深度优先、广度优先
// - 字符串：KMP、Boyer-Moore

func bubbleSort(arr []int) {
    n := len(arr)
    for i := 0; i < n-1; i++ {
        for j := 0; j < n-i-1; j++ {
            if arr[j] > arr[j+1] {
                arr[j], arr[j+1] = arr[j+1], arr[j]
            }
        }
    }
}

func selectionSort(arr []int) {
    n := len(arr)
    for i := 0; i < n-1; i++ {
        minIdx := i
        for j := i+1; j < n; j++ {
            if arr[j] < arr[minIdx] {
                minIdx = j
            }
        }
        arr[i], arr[minIdx] = arr[minIdx], arr[i]
    }
}

func insertionSort(arr []int) {
    n := len(arr)
    for i := 1; i < n; i++ {
        key := arr[i]
        j := i - 1
        for j >= 0 && arr[j] > key {
            arr[j+1] = arr[j]
            j--
        }
        arr[j+1] = key
    }
}
func quickSort(arr []int) {
    if len(arr) <= 1 {
        return
 	}
	
	pivot := arr[len(arr)/2]
	left, right := 0, len(arr)-1
	
	for left <= right {
		for arr[left] < pivot {
			left++
		}
		for arr[right] > pivot {
			right--
		}
		if left <= right {
			arr[left], arr[right] = arr[right], arr[left]
			left++
			right--
		}
	}
	
	quickSort(arr[:right+1])
	quickSort(arr[left:])
}

func mergeSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    
    mid := len(arr) / 2
    left := mergeSort(arr[:mid])
    right := mergeSort(arr[mid:])
    
    return merge(left, right)
}

func merge(left, right []int) []int {
    result := make([]int, 0, len(left)+len(right))
    i, j := 0, 0
    
    for i < len(left) && j < len(right) {
        if left[i] <= right[j] {
            result = append(result, left[i])
            i++
        } else {
            result = append(result, right[j])
            j++
        }
    }
    
    result = append(result, left[i:]...)
    result = append(result, right[j:]...)
    
    return result
}

func heapSort(arr []int) {
    n := len(arr)
    
    // 构建最大堆
    for i := n/2 - 1; i >= 0; i-- {
        heapify(arr, n, i)
    }
    
    // 一个一个提取元素
    for i := n - 1; i > 0; i-- {
        arr[0], arr[i] = arr[i], arr[0] // 交换
        heapify(arr, i, 0)
    }
}

func heapify(arr []int, n, i int) {
    largest := i
    left := 2*i + 1
    right := 2*i + 2
    
    if left < n && arr[left] > arr[largest] {
        largest = left
    }
    
    if right < n && arr[right] > arr[largest] {
        largest = right
    }
    
    if largest != i {
        arr[i], arr[largest] = arr[largest], arr[i]
        heapify(arr, n, largest)
    }
}

func binarySearch(arr []int, target int) int {
    left, right := 0, len(arr)-1
    
    for left <= right {
        mid := left + (right-left)/2
        if arr[mid] == target {
            return mid
        } else if arr[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    
    return -1
}


func KMPSearch(text, pattern string) int {
    n := len(text)
    m := len(pattern)
    
    if m == 0 {
        return 0
    }
    
    lps := make([]int, m)
    computeLPSArray(pattern, m, lps)
    
    i, j := 0, 0
    for i < n {
        if pattern[j] == text[i] {
            i++
            j++
        }
        
        if j == m {
            return i - j
        } else if i < n && pattern[j] != text[i] {
            if j != 0 {
                j = lps[j-1]
            } else {
                i++
            }
        }
    }
    
    return -1
}

func computeLPSArray(pattern string, m int, lps []int) {
    length := 0
    lps[0] = 0
    i := 1
    
    for i < m {
        if pattern[i] == pattern[length] {
            length++
            lps[i] = length
            i++
        } else {
            if length != 0 {
                length = lps[length-1]
            } else {
                lps[i] = 0
                i++
            }
        }
    }
}

// Boyer-Moore字符串匹配算法
func BoyerMooreSearch(text, pattern string) int {
    n := len(text)
    m := len(pattern)
    
    if m == 0 {
        return 0
    }
    if m > n {
        return -1
    }
    
    // 构建坏字符表
    badChar := make(map[byte]int)
    for i := 0; i < m; i++ {
        badChar[pattern[i]] = i
    }
    
    // 从右向左匹配
    shift := 0
    for shift <= n-m {
        j := m - 1
        
        // 从模式的右端开始匹配
        for j >= 0 && pattern[j] == text[shift+j] {
            j--
        }
        
        // 如果完全匹配
        if j < 0 {
            return shift
        }
        
        // 计算坏字符规则的移动距离
        badCharShift := 1
        if badCharPos, exists := badChar[text[shift+j]]; exists {
            badCharShift = maxInt(1, j-badCharPos)
        } else {
            badCharShift = j + 1
        }
        
        shift += badCharShift
    }
    
    return -1
}

// 辅助函数：返回两个整数中的最大值
func maxInt(a, b int) int {
    if a > b {
        return a
    }
    return b
}

// 辅助函数：返回两个整数中的最小值
func minInt(a, b int) int {
    if a < b {
        return a
    }
    return b
}

