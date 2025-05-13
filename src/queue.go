package main

// import "sync"

// import "fmt"

type queue []string

func (q *queue) isEmpty() bool {
	return len(*q) == 0
}
func (q *queue) enqueue(data string) {
	*q = append((*q), data)
}

func (q *queue) dequeue()(string, bool) {
	if q.isEmpty() {
		return "", false
	} else {
		// index := len(*q)
		// Print dequeued value
		dq := (*q)[0]
		// fmt.Printf("%s dequeued\n", dq)
		*q = (*q)[1:]
		return dq, true
	}
}

// type queueMulti struct {
// 	items []string
// 	mu sync.Mutex
// }

// func (q *queueMulti) enqueue(item string) {
// 	q.mu.Lock()
// 	defer q.mu.Unlock()
// 	q.items = append(q.items, item)
// }

// func (q *queueMulti) dequeue() (string, bool) {
// 	q.mu.Lock()
// 	defer q.mu.Unlock()
// 	if len(q.items) == 0 {
// 		return "", false
// 	}
// 	item := q.items[0]
// 	q.items = q.items[1:]
// 	return item, true
// }

// func (q *queueMulti) isEmpty() bool {
// 	q.mu.Lock()
// 	defer q.mu.Unlock()
// 	return len(q.items) == 0
// }
// func main() {
// 	q := queue{}
// 	q.enqueue("Ana")
// 	q.enqueue("Boni")
// 	q.enqueue("Cika")
// 	q.enqueue("Delirius")
// 	q.enqueue("Eka")

// 	fmt.Println(q)

// 	q.dequeue()
// 	q.dequeue()
// 	q.dequeue()

// 	fmt.Println(q)

// 	q.dequeue()
// 	q.dequeue()
// 	fmt.Println(q)
// }

// func main() {
// 	q := queue{}
// 	q.enqueue("Ana")
// 	q.enqueue("Boni")
// 	q.enqueue("Cika")
// 	q.enqueue("Delirius")
// 	q.enqueue("Eka")

// 	fmt.Println(q)

// 	q.dequeue()
// 	q.dequeue()
// 	q.dequeue()

// 	fmt.Println(q)

// 	q.dequeue()
// 	q.dequeue()
// 	fmt.Println(q)
// }