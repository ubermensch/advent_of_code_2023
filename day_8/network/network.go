package network

import "bufio"

type Node struct {
	label string
	left  *Node
	right *Node
}

type Network struct {
	directions []byte
	nodes      []*Node
}

func NewNetwork(s *bufio.Scanner) (*Network, error) {
	return nil, nil
}

func (n *Network) StepsToFinish() int {
	return 0
}
