package data_structs

import "cycles/types"

type Graph [][]*types.Edge

func (graph *Graph) IsConnected(i, j int) bool {
	return (*graph)[i][j] != nil
}

func (graph *Graph) GetOnlyNumbers() [][]int {
	onlyNumbers := make([][]int, len(*graph))
	for i := 0; i < len(onlyNumbers); i++ {
		onlyNumbers[i] = make([]int, len((*graph)[i]))
		for j := 0; j < len(onlyNumbers[i]); j++ {
			if (*graph)[i][j] != nil {
				onlyNumbers[i][j] = 1
			}
		}
	}
	return onlyNumbers
}
