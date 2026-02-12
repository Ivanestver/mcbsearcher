package data_structs

type Stack[T any] []T

func (stack *Stack[T]) Push(p T) {
	*stack = append(*stack, p)
}

func (stack *Stack[T]) Pop() T {
	last := stack.last()
	p := (*stack)[last]
	*stack = (*stack)[:last]
	return p
}

func (stack *Stack[T]) Pick() T {
	return (*stack)[stack.last()]
}

func (stack *Stack[T]) IsEmpty() bool {
	return len(*stack) == 0
}

func (stack *Stack[T]) last() int {
	return len(*stack) - 1
}
