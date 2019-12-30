package pool

import (
	"bytes"
	"sync"
)

type Buffer struct {
	Buf  *bytes.Buffer
	pool *BufferPool
}

func (b *Buffer) Free() {
	b.pool.Recycle(b)
}

type BufferPool struct {
	lock sync.Mutex
	pool *sync.Pool

	// 是否回收, true为不再回收
	remove bool

	// 回收时是否清空,true为清空
	clear bool
}

func (bp *BufferPool) Init(clear bool) *BufferPool {
	bp.clear = clear
	return bp

}

// 从池中获取一个buffer, reset如果为true会将buffer置空
// 如果buffer按元素索引操作, 使用false
func (bp *BufferPool) Get() (*Buffer) {
	buf := bp.pool.Get().(*Buffer)
	buf.pool = bp
	return buf
}

// 回收buffer, reset如果为true会将buffer置空
func (bp *BufferPool) Recycle(buf *Buffer) {
	if bp.clear {
		buf.Buf.Reset()
	}
	if bp.remove {
		return
	}
	bp.pool.Put(buf)
}

// 清理buffer策略
func (bp *BufferPool) monitor() {

	//bp.remove
}

// []byte类型缓冲池
func NewBytesPool(ln, cap int) *BufferPool {
	if ln > cap {
		cap = ln
	}
	return &BufferPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return &Buffer{Buf: bytes.NewBuffer(make([]byte, ln, cap))}
			},
		},
	}
}
