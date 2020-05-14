package bloomfilter

import (
	"fmt"
	"testing"
)

func TestBloomfilter_Contains(t *testing.T) {
	t.Run("same", func(t *testing.T) {
		bf := New()
		c1 := []byte("c1")
		c2 := []byte("c2")
		bf.Add(c1)
		bf.Add(c2)

		if got := bf.Contains(c1); !got {
			t.Errorf("Bloomfilter.Contains() = %v, want %v", got, true)
		}
		if got := bf.Contains(c2); !got {
			t.Errorf("Bloomfilter.Contains() = %v, want %v", got, true)
		}
		if got := bf.Contains([]byte("c3")); got {
			t.Errorf("Bloomfilter.Contains() = %v, want %v", got, true)
		}
	})
	t.Run("ip", func(t *testing.T) {
		bf := New()
		bf.Add([]byte("127.0.0.1"))

		for i := 0; i < 256; i++ {
			if got := bf.Contains([]byte(fmt.Sprintf("1.2.%d.4", i))); got {
				t.Errorf("Bloomfilter.Contains() = %v, want %v", got, false)
			}
		}

		for i := 0; i < 256; i++ {
			bf.Add([]byte(fmt.Sprintf("1.2.%d.4", i)))
		}

		if got := bf.Contains([]byte("127.0.0.1")); !got {
			t.Errorf("Bloomfilter.Contains() = %v, want %v", got, true)
		}
		if got := bf.Contains([]byte("1.3.3.4")); got {
			t.Errorf("Bloomfilter.Contains() = %v, want %v", got, false)
		}

		nbf := RecoverFrom(bf.Buffer(), bf.Len())

		for i := 0; i < 256; i++ {
			if got := nbf.Contains([]byte(fmt.Sprintf("1.2.%d.4", i))); !got {
				t.Errorf("Bloomfilter.Contains() = %v, want %v", got, true)
			}
		}
	})
}

func TestBloomfilter_Len(t *testing.T) {
	t.Run("zerot", func(t *testing.T) {
		bf := New()
		if got := bf.Len(); got != 0 {
			t.Errorf("Bloomfilter.Len() = %v, want %v", got, 0)
		}
	})
	t.Run("recover_add", func(t *testing.T) {
		bf := RecoverFrom([256]byte{}, 11)
		if got := bf.Len(); got != 11 {
			t.Errorf("Bloomfilter.Len() = %v, want %v", got, 11)
		}
		bf.Add([]byte("a"))
		if got := bf.Len(); got != 12 {
			t.Errorf("Bloomfilter.Len() = %v, want %v", got, 12)
		}
	})
}

func Test_sumshift(t *testing.T) {
	var dt1, dt2 [32]byte
	dt1[30] = 128
	dt2[30] = 128

	t.Run(">8", func(t *testing.T) {
		sumshift(&dt1, 11)
		if dt1[30+(11/8)] != 128>>(11%8) {
			t.Errorf("right shitf %d bits failed", 11)
		}
	})
	t.Run("<8", func(t *testing.T) {
		sumshift(&dt2, 3)
		if dt2[30] != 128>>3 {
			t.Errorf("right shift %d bits failed", 11)
		}
	})
}

func BenchmarkBloomfilter_Add(b *testing.B) {
	bf := New()
	for i := 0; i < b.N; i++ {
		bf.Add([]byte(fmt.Sprintf("%d", i)))
	}
}
