package main

type Hit struct {
	m   float64
	dot float64
}

type Hits struct {
	hits  []Hit
	index int
}

//type Hits []Hit

func (h *Hits) Len() int {
	return h.index
}

func (h *Hits) Less(i, j int) bool {
	return h.hits[i].m < h.hits[j].m
}

func (h *Hits) Swap(i, j int) {
	h.hits[i], h.hits[j] = h.hits[j], h.hits[i]
}

func (h *Hits) Add(hit Hit) {
	h.hits[h.index] = hit
	h.index++
}

func (h *Hits) Get(index int) *Hit {
	return &h.hits[index]
}
