package cmap

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"unsafe"
)

// 一个map有多个段用来并发，每个段会上锁，每个段中包含多个bucket，每个bucket是一个Pair的链表

// 表示并发安全的字典的接口
type ConcurrentMap interface {
	// 返回并发量
	Concurrency() int
	// element不能为nil，是否新增了值，若键已存在则替换
	Put(key string, element interface{}) (bool, error)
	// 若返回nil，说明指定键不存在
	Get(key string) interface{}
	// 删除指定键，若返回true说明键已存在且已删除，否则说明键不存在
	Delete(key string) bool
	// 返回当前字典中键-元素对的数量
	Len() uint64
}

// 实现
type myConcurrentMap struct {
	// 大的方向的并发
	concurrency int
	segments    []Segment
	total       uint64
}

func (cmap *myConcurrentMap) Concurrency() int {
	return cmap.concurrency
}

func (cmap *myConcurrentMap) Len() uint64 {
	return cmap.total
}

const MAX_CONCURRENCY = 16
const DEFAULT_BUCKET_NUMBER = 16

// pairRedistributor可以为nil
func NewConcurrentMap(concurrency int, pairRedistributor PairRedistributor) (ConcurrentMap, error) {
	if concurrency <= 0 {
		return nil, newIllegalParameterError("concurrency is too small")
	}

	if concurrency > MAX_CONCURRENCY {
		return nil, newIllegalParameterError("concurrency is too large")
	}

	cmap := &myConcurrentMap{}
	cmap.concurrency = concurrency
	cmap.segments = make([]Segment, concurrency)
	for i := 0; i < concurrency; i++ {
		cmap.segments[i] = newSegment(DEFAULT_BUCKET_NUMBER, pairRedistributor)
	}
	return cmap, nil
}

func newIllegalParameterError(s string) error {
	return nil
}

func (cmap *myConcurrentMap) Put(key string, element interface{}) (bool, error) {
	p, err := newPair(key, element)
	if err != nil {
		return false, err
	}
	s := cmap.findSegment(p.Hash())
	ok, err := s.Put(p)
	if ok {
		atomic.AddUint64(&cmap.total, 1)
	}
	return ok, err
}

// 表示并发安全的键-元素对的接口
type Pair interface {
	// 单链键-元素对接口
	linkedPair
	// 返回键的值
	Key() string
	// 返回键的散列值
	Hash() uint64
	// 返回元素的值
	Element() interface{}
	SetElement(element interface{}) error
	// 生成一个当前键-元素对的副本并返回
	Copy() Pair
	String() string
}

// 表示单向链接的键-元素对的接口
type linkedPair interface {
	// 获得下一个键-元素对，若返回nil说明当前已在单链表的末尾
	Next() Pair
	// 设定下一个键-元素对
	SetNext(nextPair Pair) error
}

// 表示键-元素对的类型
type pair struct {
	key string
	// 表示键的散列值
	hash    uint64
	element unsafe.Pointer
	next    unsafe.Pointer
}

func (p pair) Next() Pair {
	panic("implement me")
}

func (p pair) SetNext(nextPair Pair) error {
	panic("implement me")
}

func (p pair) Key() string {
	panic("implement me")
}

func (p pair) Hash() uint64 {
	panic("implement me")
}

func (p pair) Element() interface{} {
	return atomic.LoadPointer(&p.element)
}

func (p pair) SetElement(element interface{}) error {
	atomic.StorePointer(&p.element, unsafe.Pointer(&element))
	return nil
}

func (p pair) Copy() Pair {
	i, err := newPair(p.key, *(atomic.LoadPointer(&p.element)))
	if err != nil {
		fmt.Printf("Copy err:%v\n", err)
		return nil
	}

	return i
}

func (p pair) String() string {
	return fmt.Sprintf("pair[key:%v, element:%v]\n", p.key, *(atomic.LoadPointer(&p.element)))
}

func newPair(key string, element interface{}) (Pair, error) {
	p := &pair{
		key:  key,
		hash: hash(key),
	}
	if element == nil {
		return nil, newIllegalParameterError("element is nil")
	}
	p.element = unsafe.Pointer(&element)
	return p, nil
}

func hash(s string) uint64 {
	return -1
}

// 根据给定参数寻找并返回对应散列段，用高位取模
func (cmap *myConcurrentMap) findSegment(keyHash uint64) Segment {
	if cmap.concurrency == 1 {
		return cmap.segments[0]
	}

	var keyHash32 uint32
	if keyHash > math.MaxUint32 {
		keyHash32 = uint32(keyHash >> 32)
	} else {
		keyHash32 = uint32(keyHash)
	}
	return cmap.segments[int(keyHash32>>16)%(cmap.concurrency-1)]
}

func (cmap *myConcurrentMap) Get(key string) interface{} {
	keyHash := hash(key)
	s := cmap.findSegment(keyHash)
	pair := s.GetWithHash(key, keyHash)
	if pair == nil {
		return nil
	}
	return pair.Element()
}

func (cmap *myConcurrentMap) Delete(key string) bool {
	s := cmap.findSegment(hash(key))
	if s.Delete(key) {
		atomic.AddUint64(&cmap.total, ^uint64(0))
		return true
	}
	return false
}

// 表示并发安全的散列段的接口
type Segment interface {
	// 返回值表示是否新增了键-元素对
	Put(p Pair) (bool, error)
	Get(key string) Pair
	// keyHash是基于参数key计算得出的散列值，避免重复计算hash
	GetWithHash(key string, keyHash uint64) Pair
	// 若返回true表明已删除，否则说明未找到
	Delete(key string) bool
	// 获取当前段的尺寸(包含的散列桶的数量)
	Size() uint64
}

// 表示并发安全的散列段的类型
type segment struct {
	// 表示散列桶，尽量均衡，便于查找，免去链表O(n)
	buckets    []Bucket
	bucketsLen int
	pairTotal  uint64
	// 表示键-元素对的再分布器，用于把散列段中所有键-元素对均匀地封不到所有散列桶中
	pairRedistributor PairRedistributor
	lock              sync.Mutex
}

func (s *segment) Size() uint64 {
	panic("implement me")
}

const DEFAULT_BUCKET_LOAD_FACTOR = 0.75

func newSegment(bucketNumber int, pairRedistributor PairRedistributor) Segment {
	if bucketNumber <= 0 {
		bucketNumber = DEFAULT_BUCKET_NUMBER
	}
	if pairRedistributor == nil {
		pairRedistributor = newDefaultPairRedistributor(DEFAULT_BUCKET_LOAD_FACTOR, bucketNumber)
	}
	buckets := make([]Bucket, bucketNumber)
	for i := 0; i < bucketNumber; i++ {
		buckets[i] = newBucket()
	}
	return &segment{
		buckets:           buckets,
		bucketsLen:        bucketNumber,
		pairRedistributor: pairRedistributor,
	}
}

func newDefaultPairRedistributor(f float64, i int) PairRedistributor {
	return nil
}

type BucketStatus int64

// 表示针对键-元素对的再分布器，当散列段内的键-元素对分布不均时进行重新分布
type PairRedistributor interface {
	// 根据键-元素对总数和散列桶总数计算并更新阈值
	UpdateThreshold(pairTotal uint64, bucketNumber int)
	CheckBucketStatus(pairTotal uint64, bucketSize uint64) (bucketStatus BucketStatus)
	// 实施键-元素对的再分布
	Redistribute(bucketStatus BucketStatus, buckets []Bucket) (newBuckets []Bucket, changed bool)
}

func (s *segment) Put(p Pair) (bool, error) {
	s.lock.Lock()
	b := s.buckets[int(p.Hash()%uint64(s.bucketsLen))]
	ok, err := b.Put(p, nil)
	if ok {
		newTotal := atomic.AddUint64(&s.pairTotal, 1)
		s.redistribute(newTotal, b.Size())
	}
	s.lock.Unlock()
	return ok, err
}

func (s *segment) Get(key string) Pair {
	return s.GetWithHash(key, hash(key))
}

func (s *segment) GetWithHash(key string, keyHash uint64) Pair {
	s.lock.Lock()
	b := s.buckets[int(keyHash%uint64(s.bucketsLen))]
	s.lock.Unlock()
	return b.Get(key)
}

func (s *segment) Delete(key string) bool {
	s.lock.Lock()
	b := s.buckets[int(hash(key)%uint64(s.bucketsLen))]
	ok := b.Delete(key, nil)
	if ok {
		newTotal := atomic.AddUint64(&s.pairTotal, ^uint64(0))
		s.redistribute(newTotal, b.Size())
	}
	s.lock.Unlock()
	return ok
}

// 检查给定参数并设置阈值和计数，在必要时重新分配所有散列桶中的所有键-元素对，必须在互斥锁的保护下调用此方法
func (s *segment) redistribute(pairTotal uint64, bucketSize uint64) (err error) {
	defer func() {
		if p := recover(); p != nil {
			if pErr, ok := p.(error); ok {
				err = newPairRedistributorError(pErr.Error())
			} else {
				err = newPairRedistributorError(fmt.Sprintf("%s", p))
			}
		}
	}()
	s.pairRedistributor.UpdateThreshold(pairTotal, s.bucketsLen)
	bucketStatus := s.pairRedistributor.CheckBucketStatus(pairTotal, bucketSize)
	newBuckets, changed := s.pairRedistributor.Redistribute(bucketStatus, s.buckets)
	if changed {
		s.buckets = newBuckets
		s.bucketsLen = len(s.buckets)
	}
	return nil
}

func newPairRedistributorError(s string) error {
	return nil
}

// 表示并发安全的散列桶接口
type Bucket interface {
	// 返回值表示是否新增了键-元素对，若在调用此方法前已锁定lock，则不要传入lock，否则需要传入对应的lock
	Put(p Pair, lock sync.Locker) (bool, error)
	Get(key string) Pair
	GetFirstPair() Pair
	// 若在调用此方法前已锁定lock，则不要传入lock，否则需要传入对应的lock
	Delete(key string, lock sync.Locker) bool
	// 若在调用此方法前已锁定lock，则不要传入lock，否则需要传入对应的lock
	Clear(lock sync.Locker)
	// 当前散列桶的尺寸
	Size() uint64
	String() string
}

// 表示并发安全的散列桶的类型
type bucket struct {
	// 存储键-元素对列表的表头
	firstValue atomic.Value
	size       uint64
}

func (b *bucket) Size() uint64 {
	panic("implement me")
}

func (b *bucket) String() string {
	panic("implement me")
}

func newBucket() Bucket {
	b := &bucket{}
	b.firstValue.Store(placeholder)
	return b
}

// 占位符
// 由于原子值不能存储(Store)nil，所以当散列桶空时用此符占位
var placeholder Pair = &pair{}

func (b *bucket) GetFirstPair() Pair {
	if v := b.firstValue.Load(); v == nil {
		return nil
	} else if p, ok := v.(Pair); !ok || p == placeholder {
		return nil
	} else {
		return p
	}
}

func (b *bucket) Put(p Pair, lock sync.Locker) (bool, error) {
	if p == nil {
		return false, newIllegalParameterError("pair is nil")
	}

	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	firstPair := b.GetFirstPair()
	if firstPair == nil {
		b.firstValue.Store(p)
		atomic.AddUint64(&b.size, 1)
		return true, nil
	}
	var target Pair
	key := p.Key()
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			target = v
			break
		}
	}
	if target != nil {
		target.SetElement(p.Element())
		return false, nil
	}
	p.SetNext(firstPair)
	b.firstValue.Store(p)
	atomic.AddUint64(&b.size, 1)
	return true, nil
}

func (b *bucket) Delete(key string, lock sync.Locker) bool {
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	firstPair := b.GetFirstPair()
	if firstPair == nil {
		return false
	}

	var prevPairs []Pair
	var target Pair
	var breakpoint Pair
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			target = v
			breakpoint = v.Next()
			break
		}
		prevPairs = append(prevPairs, v)
	}
	if target == nil {
		return false
	}
	newFirstPair := breakpoint
	for i := len(prevPairs) - 1; i >= 0; i-- {
		pairCopy := prevPairs[i].Copy()
		pairCopy.SetNext(newFirstPair)
		newFirstPair = pairCopy
	}
	if newFirstPair != nil {
		b.firstValue.Store(newFirstPair)
	} else {
		b.firstValue.Store(placeholder)
	}
	atomic.AddUint64(&b.size, ^uint64(0))
	return true
}

func (b *bucket) Get(key string) Pair {
	firstPair := b.GetFirstPair()
	if firstPair == nil {
		return nil
	}
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			return v
		}
	}
	return nil
}

func (b *bucket) Clear(lock sync.Locker) {
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}

	atomic.StoreUint64(&b.size, 0)
	b.firstValue.Store(placeholder)
}
