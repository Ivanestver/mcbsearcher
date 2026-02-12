package types

type SupportVector []uint64

func (s SupportVector) GetScalarMultiplication(other SupportVector) uint64 {
	if len(s) != len(other) {
		panic("Support vectors must be of the same size")
	}
	var res uint64 = 0
	for i := range s {
		res += s[i] * other[i]
	}
	return res
}

func (s SupportVector) XOR(other SupportVector) SupportVector {
	res := make(SupportVector, len(s))
	for i := range s {
		res[i] = s[i] ^ other[i]
	}
	return res
}
