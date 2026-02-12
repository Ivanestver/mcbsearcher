package data_structs

type PQItem struct {
	Dist float64
	Num  int
}

type PriorityQueue []PQItem

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Dist < pq[j].Dist
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	*pq = append(*pq, x.(PQItem))
}

func (pq *PriorityQueue) Pop() any {
	n := pq.Len()
	it := (*pq)[n-1]
	*pq = (*pq)[:n-1]
	return it
}
