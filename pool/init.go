package pool

import (
    "sync"
)

func init() {
    BS = &ByteSlicePool {
        pool : &sync.Pool{
            New: func() interface{} {
                return &ByteSlice{}
            },
        },
    }

    TwoDimBS = &TwoDimByteSlicePool {
        pool : &sync.Pool{
            New: func() interface{} {
                return &TwoDimByteSlice{}
            },
        },
    }

    for i:= 0; i < 2048; i++ {
        bs := new(ByteSlice)
        bs.Grow(4096, 64)
        BS.pool.Put(bs)

        bs2 := new(TwoDimByteSlice)
        bs2.Grow(4096, 20, 20)
        TwoDimBS.pool.Put(bs2)
    }
}
