package bloomfilter

import (
	"crypto/sha256"
	"encoding/binary"
)

// Bloomfilter is a space-efficient probabilistic data structure.
// False-positive rate:
// 	* 1e-05 for  <80 elements
// 	* 1e-04 for <105 elements
//	* 1e-03 for <142 elements
// Uses a 256 byte array (2048 bits) and 11 hash functions. 256 byte because
// of space efficiency (array is saved for each comment) and 11 hash functions
// because of best overall false-positive rate in that range.
type Bloomfilter struct {
	buffer [256]byte
	num  int
	k    int
	m    uint16
}

// New return a empty Bloomfilter.
func New() *Bloomfilter {
	return &Bloomfilter{
		buffer: [256]byte{},
		num:  0,
		k:    11,
		m:    256 * 8,
	}
}

// RecoverFrom recover a Bloomfilter instance use preovide `data` and `num` of elements
func RecoverFrom(data [256]byte, num int) *Bloomfilter {
	return &Bloomfilter{
		buffer: data,
		num:  num,
		k:    11,
		m:    uint16(len(data) * 8),
	}
}

// Len return the amount of added elements
func (bf *Bloomfilter) Len() int {
	return bf.num
}

// Buffer return the buffer of Bloomfilter
func (bf *Bloomfilter) Buffer() [256]byte {
	return bf.buffer
}

// Add add a element into Bloomfilter
func (bf *Bloomfilter) Add(e []byte) {
	sum := sha256.Sum256(e)
	for i := 0; i < bf.k; i++ {
		p := binary.BigEndian.Uint16(sum[30:])
		n := p & (bf.m - 1)
		bf.buffer[n/8] |= 1 << (n % 8)
		sumshift(&sum, bf.k)
	}
	bf.num++
}

// Contains check a element whether in Bloomfilter
func (bf *Bloomfilter) Contains(e []byte) bool {
	sum := sha256.Sum256(e)
	for i := 0; i < bf.k; i++ {
		p := binary.BigEndian.Uint16(sum[sha256.Size-2:])
		n := p & (bf.m - 1)
		if bf.buffer[n/8]&(1<<(n%8)) == 0 {
			return false
		}
		sumshift(&sum, bf.k)
	}
	return true
}

func sumshift(data *[32]byte, bits int) {
	n := 32
	r8, r := bits/8, bits%8

	var shifted [32]byte
	for i := n - 1; i-r8-1 >= 0; i-- {
		shifted[i] = (*data)[i-r8] >> r
		shifted[i] |= (*data)[i-r8-1] << (8 - r)
	}
	*data = shifted
}
