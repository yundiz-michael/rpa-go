package common

type hash struct {
	str string
}

func NewHash() *hash {
	return &hash{}
}

func (h *hash) Seed(str string) *hash {
	h.str = str
	return h
}

func (h *hash) DeHash() int {
	hash := 5381
	length := len(h.str)
	seed := 5
	for i := 0; i < length; i++ {
		hash += (hash << seed) + int(h.str[i])
	}
	return hash & 0x7FFFFFFF
}
