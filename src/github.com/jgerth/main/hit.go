package main


type Hit struct {
	m   float64
	dot float64
}

type Hits []Hit

func (slice Hits) Len() int {
	return len(slice)
}

func (slice Hits) Less(i, j int) bool {
	return slice[i].m < slice[j].m
}

func (slice Hits) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

