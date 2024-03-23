package network

import (
	"bufio"
	"fmt"
	"github.com/samber/lo"
	"strings"
)

const End = "ZZZ"

type Node struct {
	label string
	left  string
	right string
}

type Network struct {
	directions []byte
	nodes      []*Node
	lookup     map[string]*Node
}

func NewNetwork(s *bufio.Scanner) (*Network, error) {
	getDirections := func(dirStr string) []byte {
		return []byte(dirStr)
	}
	getNodes := func(lines []string) []*Node {
		return lo.Map(lines, func(s string, i int) *Node {
			pieces := strings.Split(s, "=")
			label := strings.Trim(pieces[0], " ")
			next := strings.Trim(pieces[1], " ()")
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

func (n *Network) StepsToFinish() int64 {
	steps := int64(0)
	currNode := n.nodes[0]
	exitFound := false
	fmt.Println("directions: %s", n.directions)
	fmt.Println("")
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
			if steps%1e6 == 0 {
				fmt.Printf("currLabel: %s, steps: %d\n", currNode.label, steps)
			}

			if currNode.label == "ZZZ" {
				exitFound = true
				break
			}
		}
	}
	return steps
}
