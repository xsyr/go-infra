package pool

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestTwoByteSlice(t *testing.T) {
    var bs TwoDimByteSlice
    testF := func(b *TwoDimByteSlice) {
        b.Grow(512, 10, 10)
        assert.Equal(t, 0  , b.Dim())
        assert.Equal(t, 0  , len(b.data))
        assert.Equal(t, 512, cap(b.data))

        assert.Equal(t, 0 , len(b.flat))
        assert.Equal(t, 10, cap(b.flat))

        assert.Equal(t, 0 , len(b.dim))
        assert.Equal(t, 10, cap(b.dim))

        {
            // dim 1
            assert.Equal(t, 1, b.NewDim())

            b.Append([]byte("a"))
            b.Append([]byte("b"))
            b.Append([]byte("c"))
            b.Concat([]byte("d"), []byte("e"))

            assert.Equal(t, 1  , b.Dim())
            assert.Equal(t, 4  , b.Len(0))
            assert.Equal(t, 5  , len(b.data))
            assert.Equal(t, 512, cap(b.data))

            assert.Equal(t, 4 , len(b.flat))
            assert.Equal(t, 10, cap(b.flat))
        }

        {
            // dim 2
            assert.Equal(t, 2, b.NewDim())

            b.Append([]byte("j"))
            b.Append([]byte("k"))

            assert.Equal(t, 2  , b.Dim())
            assert.Equal(t, 2  , b.Len(1))
            assert.Equal(t, 7  , len(b.data))
            assert.Equal(t, 512, cap(b.data))

            assert.Equal(t, 6 , len(b.flat))
            assert.Equal(t, 10, cap(b.flat))
        }

        {
            // growing
            b.Grow(1024, 20, 20)
            assert.Equal(t, 2   , b.Dim())
            assert.Equal(t, 7   , len(b.data))
            assert.Equal(t, 1024, cap(b.data))

            assert.Equal(t, 6 , len(b.flat))
            assert.Equal(t, 20, cap(b.flat))

            assert.Equal(t, 2 , len(b.dim))
            assert.Equal(t, 20, cap(b.dim))
        }

        var bs2 [50][]byte
        {
            // dim 1
            s := b.ToByteSlice(0, bs2[:0])
            assert.Equal(t, 4, len(s))
            assert.Equal(t, []byte("a") , s[0])
            assert.Equal(t, []byte("b") , s[1])
            assert.Equal(t, []byte("c") , s[2])
            assert.Equal(t, []byte("de"), s[3])
        }
        {
            // dim 2
            s := b.ToByteSlice(1, bs2[:0])
            assert.Equal(t, 2, len(s))
            assert.Equal(t, []byte("j") , s[0])
            assert.Equal(t, []byte("k") , s[1])
        }

        assert.Equal(t, []byte("a") , b.Index(0, 0))
        assert.Equal(t, []byte("b") , b.Index(0, 1))
        assert.Equal(t, []byte("c") , b.Index(0, 2))
        assert.Equal(t, []byte("de"), b.Index(0, 3))

        assert.Equal(t, []byte("j") , b.Index(1, 0))
        assert.Equal(t, []byte("k") , b.Index(1, 1))
    }
    testF(&bs)

    bs.Reset()
    assert.Equal(t, 0   , bs.Dim())
    assert.Equal(t, 0   , len(bs.data))
    assert.Equal(t, 1024, cap(bs.data))

    assert.Equal(t, 0  , len(bs.flat))
    assert.Equal(t, 20 , cap(bs.flat))

    assert.Equal(t, 0  , len(bs.dim))
    assert.Equal(t, 20 , cap(bs.dim))
}

