package main

import "fmt"

type queue []string

func (q *queue) isEmpty() bool {
	return len(*q) == 0
}
func (q *queue) enqueue(data string) {
	*q = append((*q), data)
}

func (q *queue) dequeue() {
	if q.isEmpty() {
		return
	} else {
		index := len(*q)
		// Print dequeued value
		dq := (*q)[0]
		fmt.Printf("%s dequeued\n", dq)
		*q = (*q)[1:index]
	}
}

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
