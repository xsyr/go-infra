package pool

import (
    "io"
    "sync"
)

var BS *ByteSlicePool

type sliceHeader struct {
    offset, len int
}

type ByteSlice struct {
    data []byte
    elems []sliceHeader
}

func (b *ByteSlice)Grow(dataCap, elemCap int) {
    dataLen := len(b.data)
    dataExtend := dataCap - cap(b.data)
    if dataExtend > 0 {
        b.data = append(b.data, make([]byte, dataExtend)...)
        b.data = b.data[:dataLen]
    }

    elemLen := len(b.elems)
    elemExtend := elemCap - cap(b.elems)
    if elemExtend > 0 {
        b.elems = append(b.elems, make([]sliceHeader, elemExtend)...)
        b.elems = b.elems[:elemLen]
    }
}

func (b *ByteSlice)Len() int {
    return len(b.elems)
}

func (b *ByteSlice)Index(index int) []byte {
    start := b.elems[index].offset
    end := start + b.elems[index].len
    return b.data[start:end]
}

func (b *ByteSlice)CopyTo(index int, buf []byte) []byte {
    start := b.elems[index].offset
    end := start + b.elems[index].len
    buf = append(buf, b.data[start:end]...)
    return buf
}

func (b *ByteSlice)Reset() {
    b.data  = b.data[:0]
    b.elems = b.elems[:0]
}

func (b *ByteSlice)Release() {
    b.Reset()
    BS.pool.Put(b)
}

func (b *ByteSlice)AppendConcat(bs ...[]byte) *ByteSlice{
    used := len(b.data)

    length := 0
    for _, s := range bs {
        length = length + len(s)
        b.data = append(b.data, s...)
    }
    b.elems = append(b.elems, sliceHeader{ used, length })
    return b
}

func (b *ByteSlice)AppendFromReaderN(rd io.Reader, expect int) (n int, err error) {
    used := len(b.data)
    b.data = append(b.data, make([]byte, expect)...)

    n, err = io.ReadFull(rd, b.data[used:])
    if err != nil {
        b.data = b.data[:used]
        return n, err
    }

    b.elems = append(b.elems, sliceHeader{ used, expect })
    return n, err
}

func (b *ByteSlice)Append(bs ...[]byte) *ByteSlice{
    for _, bts := range bs {
        length := len(b.data)
        b.data  = append(b.data, bts...)
        b.elems = append(b.elems, sliceHeader{ length, len(bts) })
    }
    return b
}

func (b *ByteSlice)ToByteSlice(bs [][]byte) [][]byte {
    for _, h := range b.elems {
        bs = append(bs, b.data[h.offset:h.offset+h.len])
    }
    return bs
}

type ByteSlicePool struct {
    pool *sync.Pool
}

func (p *ByteSlicePool)Get() *ByteSlice {
    return p.pool.Get().(*ByteSlice)
}

func init() {
    BS = &ByteSlicePool {
        pool : &sync.Pool{
            New: func() interface{} {
                return &ByteSlice{}
            },
        },
    }
}
