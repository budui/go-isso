package util

type Bloomfilter struct {
	Array    []byte
	elements int
	k        int
	m        int
}

func (bf *Bloomfilter) Add(key string) {

}

func (bf *Bloomfilter) getProbes(key string) {

}

func GenBloomfilterfunc(p string) []byte {
	array := make([]byte, 256)
	return array
}
