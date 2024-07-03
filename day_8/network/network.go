package network

import (
	"bufio"
	"github.com/samber/lo"
	"strings"
)

// greatest common divisor (GCD) via Euclidean algorithm
func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func lcm(a, b int, integers ...int) int {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
}

const Start = "AAA"
const End = "ZZZ"

type Node struct {
	label string
	left  string
	right string
}

func (n *Node) EndsWith(b byte) bool {
	return []byte(n.label)[len(n.label)-1] == b
}

type Network struct {
	directions []rune
	nodes      []*Node
	lookup     map[string]*Node
}

func NewNetwork(s *bufio.Scanner) (*Network, error) {
	getDirections := func(dirStr string) []rune {
		return []rune(dirStr)
	}
	getNodes := func(lines []string) []*Node {
		return lo.Map(lines, func(s string, i int) *Node {
			pieces := strings.Split(s, "=")
			label := strings.Trim(pieces[0], " ")
			next := strings.Trim(pieces[1], " (4)")
			leftAndRight := strings.Split(next, ",")
			left, right := strings.Trim(leftAndRight[0], " "), strings.Trim(leftAndRight[1], " ")
			return &Node{
				label: label,
				left:  left,
				right: right,
			}
		})
	}
	getLookup := func(nodes []*Node) map[string]*Node {
		lookup := make(map[string]*Node)
		for _, node := range nodes {
			lookup[node.label] = node
		}
		return lookup
	}

	var lines []string
	for s.Scan() {
		line := s.Text()
		lines = append(lines, line)
	}
	directions := getDirections(lines[0])
	nodes := getNodes(lines[2:])
	lookup := getLookup(nodes)

	return &Network{
		directions: directions,
		nodes:      nodes,
		lookup:     lookup,
	}, nil

}

func (n *Network) StepsToFinish() int {
	steps := 0
	currNode := n.lookup[Start]
	exitFound := false

	for !exitFound {
		for _, curr := range n.directions {
			switch curr {
			case 'L':
				currNode = n.lookup[currNode.left]
			case 'R':
				currNode = n.lookup[currNode.right]
			default:
				panic("direction not left or right")
			}
			steps += 1

			if currNode.label == End {
				exitFound = true
				break
			}
		}
	}
	return steps
}

func (n *Network) LcmStepsToFinish() int {
	startingNodes := lo.Filter(n.nodes, func(n *Node, i int) bool {
		return []byte(n.label)[2] == 'A'
	})

	nextNode := func(curr *Node, direction rune) *Node {
		switch direction {
		case 'L':
			return n.lookup[curr.left]
		case 'R':
			return n.lookup[curr.right]
		default:
			panic("direction not left or right")
		}
	}

	stepsToFinish := func(node *Node) int {
		var currNode = node
		steps := 0
		dirI := 0
		var dir rune
		for {
			dir = n.directions[dirI]
			steps += 1
			currNode = nextNode(currNode, dir)
			if currNode.EndsWith('Z') {
				break
			}
			if dirI == len(n.directions)-1 {
				dirI = 0
			} else {
				dirI += 1
			}
		}

		return steps
	}

	steps := []int{}
	for _, curr := range startingNodes {
		steps = append(steps, stepsToFinish(curr))
	}

	return lcm(steps[0], steps[1], steps[2:]...)
}
