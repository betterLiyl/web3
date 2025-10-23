// ### 任务11：高性能数据结构
// **目标**：实现零分配的数据结构
// **描述**：实现一系列高性能、内存安全的数据结构

// **包含结构**：
// - 无锁队列（Lock-free Queue）
// - 跳表（Skip List）
// - 布隆过滤器（Bloom Filter）
// - 一致性哈希（Consistent Hash）
// - 时间轮（Timing Wheel）
// - LRU缓存（LRU Cache）

struct Queue<T> {
    data: Vec<T>,
}

impl<T> Queue<T> {
    pub fn new() -> Self {
        Self { data: Vec::new() }
    }

    pub fn enqueue(&mut self, item: T) {
        self.data.push(item);
    }

    pub fn dequeue(&mut self) -> Option<T> {
        self.data.pop()
    }

}

struct SkipList<T: Ord> {
    data: Vec<T>,
}

impl<T: Ord> SkipList<T> {
    pub fn new() -> Self {
        Self { data: Vec::new() }
    }

    pub fn insert(&mut self, item: T) {
        self.data.push(item);
    }
    
    pub fn search(&self, item: &T) -> bool {
        self.data.contains(item)
    }


}
use std::hash::{Hash, Hasher};
struct BloomFilter {
    
    data: Vec<bool>,
    size: usize,
}

impl BloomFilter {
    pub fn new() -> Self {
        Self { data: Vec::new(), size: 1000 }
    }
    pub fn hash(&self, item: &impl Hash) -> usize {
        let mut hasher = std::collections::hash_map::DefaultHasher::new();
        item.hash(&mut hasher);
        hasher.finish() as usize
    }
    
    pub fn insert(&mut self, item: &impl Hash) {
        let index = self.hash(item) % self.size;
        self.data[index] = true;
    }
    
    pub fn search(&self, item: &impl Hash) -> bool {
        let index = self.hash(item) % self.size;
        self.data[index]
    }
    pub fn init(&mut self, size: usize) {
        self.data = vec![false; size];
    }

}

fn main() {
    let mut bloom_filter = BloomFilter::new();
    bloom_filter.init(1000);
    bloom_filter.insert(&"hello");
    println!("{}", bloom_filter.search(&"hello"));
    println!("{}", bloom_filter.search(&"world"));
}