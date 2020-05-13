package advert

type item struct {
	name  string
	value int
}

type KeyValue []item

func (h KeyValue) Len() int {
	return len(h)
}

func (h KeyValue) Less(i, j int) bool {
	return h[i].value > h[j].value
}
func (h KeyValue) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *KeyValue) Push(x interface{}) {
	*h = append(*h, x.(item))
}
func (h *KeyValue) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
