package predictor

type Predictor struct {
	// history of measurements for a particular value
	history []int
	// telescoping list of deltas between one measurement and the next,
	// calculated until the deltas on a particular level converge to zero.
	deltas [][]int
}

func NewPredictor() *Predictor {
	return &Predictor{
		history: []int{},
		deltas:  [][]int{},
	}
}

func (dc *Predictor) Next() int {
	return 0
}
