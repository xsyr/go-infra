package pool

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestByteSlice(t *testing.T) {
    var bs ByteSlice
    testF := func(b *ByteSlice) {
        b.Grow(512, 10)
        assert.Equal(t, 0  , b.Len())
        assert.Equal(t, 0  , len(b.data))
        assert.Equal(t, 512, cap(b.data))

        assert.Equal(t, 0 , len(b.elems))
        assert.Equal(t, 10, cap(b.elems))

        b.Append([]byte("a"), []byte("b"), []byte("c"))
        b.Concat([]byte("d"), []byte("e"))

        assert.Equal(t, 4  , b.Len())
        assert.Equal(t, 5  , len(b.data))
        assert.Equal(t, 512, cap(b.data))

        assert.Equal(t, 4 , len(b.elems))
        assert.Equal(t, 10, cap(b.elems))

        {
            // growing
            b.Grow(1024, 20)
            assert.Equal(t, 4  , b.Len())
            assert.Equal(t, 5  , len(b.data))
            assert.Equal(t, 1024, cap(b.data))

            assert.Equal(t, 4 , len(b.elems))
            assert.Equal(t, 20, cap(b.elems))
        }

        var bs2 [50][]byte
        s := b.ToByteSlice(bs2[:0])
        assert.Equal(t, 4, len(s))
        assert.Equal(t, []byte("a") , s[0])
        assert.Equal(t, []byte("b") , s[1])
        assert.Equal(t, []byte("c") , s[2])
        assert.Equal(t, []byte("de"), s[3])
    }
    testF(&bs)

    bs.Reset()
    assert.Equal(t, 0   , bs.Len())
    assert.Equal(t, 0   , len(bs.data))
    assert.Equal(t, 1024, cap(bs.data))

    assert.Equal(t, 0  , len(bs.elems))
    assert.Equal(t, 20 , cap(bs.elems))
}

func TestByteSliceWithGrow(t *testing.T) {
    kvs := BS.Get()
    kvs.Append(
        []byte("k1base")        , []byte("k1base-value")        ,
        []byte("k1authors.base"), []byte("k1authors.base-value"),
        []byte("k1authors.imgs"), []byte("k1authors.imgs-value"),

        []byte("k2base")        , []byte("k2base-value")        ,
        []byte("k2authors.base"), []byte("k2authors.base-value"),
        []byte("k2authors.imgs"), []byte("k2authors.imgs-value"),
    )

    assert.Equal(t, 12, kvs.Len())
    assert.Equal(t, []byte("k1base")      , kvs.Index(0))
    assert.Equal(t, []byte("k1base-value"), kvs.Index(1))

    assert.Equal(t, []byte("k1authors.base")      , kvs.Index(2))
    assert.Equal(t, []byte("k1authors.base-value"), kvs.Index(3))

    assert.Equal(t, []byte("k1authors.imgs")      , kvs.Index(4))
    assert.Equal(t, []byte("k1authors.imgs-value"), kvs.Index(5))
}
