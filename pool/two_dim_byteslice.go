package pool

import (
    "sync"
)

var TwoDimBS *TwoDimByteSlicePool

type TwoDimByteSlicePool struct {
    pool *sync.Pool
}

// TwoDimByteSlice 二维递增数组
type TwoDimByteSlice struct {
    data []byte
    flat []sliceHeader
    dim  []sliceHeader
}

func (p *TwoDimByteSlicePool)Get() *TwoDimByteSlice {
    return p.pool.Get().(*TwoDimByteSlice)
}

func init() {
    TwoDimBS = &TwoDimByteSlicePool {
        pool : &sync.Pool{
            New: func() interface{} {
                return &TwoDimByteSlice{}
            },
        },
    }
}

func (b *TwoDimByteSlice)Release() {
    b.Reset()
    TwoDimBS.pool.Put(b)
}

func (b *TwoDimByteSlice)Reset() {
    b.data = b.data[:0]
    b.flat = b.flat[:0]
    b.dim  = b.dim[:0]
}

func (b *TwoDimByteSlice)Grow(dataCap, flatCap, dimCap int) {
    dataLen := len(b.data)
    dataExtend := dataCap - cap(b.data)
    if dataExtend > 0 {
        b.data = append(b.data, make([]byte, dataExtend)...)
        b.data = b.data[:dataLen]
    }

    flatLen := len(b.flat)
    flatExtend := flatCap * dimCap - cap(b.flat)
    if flatExtend > 0 {
        b.flat = append(b.flat, make([]sliceHeader, flatExtend)...)
        b.flat = b.flat[:flatLen]
    }

    dimLen := len(b.dim)
    dimExtend := dimCap - cap(b.dim)
    if dimExtend > 0 {
        b.dim = append(b.dim, make([]sliceHeader, dimExtend)...)
        b.dim = b.dim[:dimLen]
    }
}

func (b *TwoDimByteSlice)NewDim() int {
    b.dim = append(b.dim, sliceHeader{ len(b.flat), 0 })
    return b.Dim()
}

func (b *TwoDimByteSlice)Dim() int {
    return len(b.dim)
}

func (b *TwoDimByteSlice)Len(dim int) int {
    return b.dim[dim].len
}

// Index 不要修改返回的数组内容
func (b *TwoDimByteSlice)Index(dim, index int) []byte {
    flat := b.dim[dim].offset
    start := b.flat[flat+index].offset
    end   := start + b.flat[flat+index].len
    return b.data[start:end]
}

func (b *TwoDimByteSlice)Copy(dim, index int, buf []byte) []byte {
    flat := b.dim[dim].offset
    start := b.flat[flat+index].offset
    end   := start + b.flat[flat+index].len
    buf = append(buf, b.data[start:end]...)
    return buf
}

func (b *TwoDimByteSlice)Concat(bs ...[]byte) {
    offset := len(b.data)

    length := 0
    for _, s := range bs {
        length = length + len(s)
        b.data = append(b.data, s...)
    }
    b.flat = append(b.flat, sliceHeader{ offset, length })
    b.dim[len(b.dim)-1].len++
}

func (b *TwoDimByteSlice)Append(bs []byte) {
    length := len(b.data)
    b.data  = append(b.data, bs...)
    b.flat = append(b.flat, sliceHeader{ length, len(bs) })
    b.dim[len(b.dim)-1].len++
}

func (b *TwoDimByteSlice)ToByteSlice(dim int, bs [][]byte) [][]byte {
    start := b.dim[dim].offset
    end   := start + b.dim[dim].len
    for _, h := range b.flat[start:end] {
        bs = append(bs, b.data[h.offset:h.offset+h.len])
    }
    return bs
}

