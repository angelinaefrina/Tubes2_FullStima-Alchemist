package main

type stack []string

func (s *stack) isEmpty() bool {
	return len(*s) == 0
}

func (s *stack) push(data string) {
	*s = append((*s), data)
}

func (s *stack) pop() {
	if s.isEmpty() {
		return
	} else {
		index := len(*s) - 1
		// // Print popped value
		// pop := (*s)[index]
		// fmt.Printf("%s popped\n", pop)
		*s = (*s)[:index]
	}
}

// func main() {
// 	s := stack{}
// 	s.push("Aku")
// 	s.push("Adalah")
// 	s.push("Seorang")
// 	s.push("Kapiten")
// 	s.push("Hadokk")

// 	fmt.Println(s)

// 	s.pop()
// 	s.pop()
// 	s.pop()

// 	fmt.Println(s)

// 	s.pop()
// 	s.pop()
// 	fmt.Println(s)
// }
