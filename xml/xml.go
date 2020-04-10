package xml

import (
	"bytes"
	xmlencode "encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Parse returns a slice of strings representing the content of the element(s)
// found at the path position
func Parse(path string, data []byte) []string {
	var n Node

	stack := strings.Split(path, "/")

	buf := bytes.NewBuffer(data)

	err := xmlencode.NewDecoder(buf).Decode(&n)
	if err != nil {
		panic(err)
	}
	values := make([]string, 0)
	xmlReader([]Node{n}, &stack, &values)
	return values
}

// ParseRecursive returns a slice of strings representing the content of the
// element(s) found at the path position recursively
func ParseRecursive(path string, data []byte) []string {
	var n Node

	stack := strings.Split(path, "/")

	buf := bytes.NewBuffer(data)

	err := xmlencode.NewDecoder(buf).Decode(&n)
	if err != nil {
		panic(err)
	}
	values := make([]string, 0)
	xmlReaderRecursive([]Node{n}, &stack, &values)
	return values
}

// xmlReaderRecursive is a recursive function that will unstack an array of elements
// until reaching the desired node
func xmlReaderRecursive(nodes []Node, stack *[]string, values *[]string) {
	for _, n := range nodes {
		if len(*stack) > 0 {
			if !trimStackIndex(n, stack) && len(*stack) > 1 && n.XMLName.Local == (*stack)[0] {
				*stack = (*stack)[1:]
			}

			if v := getValuesWithAttributes(n, stack, func(n Node, att string) []string {
				v := n.valuesFromAttributes(att)
				n.valuesRecursive(&v)
				return v
			}); v != nil {
				*values = append(*values, v...)
			}

			if v := getValues(n, stack, func(n Node) []string {
				v := []string{}
				n.valuesRecursive(&v)
				return v
			}); v != nil {
				*values = append(*values, v...)
			}
			xmlReaderRecursive(n.Nodes, stack, values)
		}
	}
}

// xmlReader is a recursive function that will unstack an array of elements
// until reaching the desired node
func xmlReader(nodes []Node, stack *[]string, values *[]string) {
	for _, n := range nodes {
		if len(*stack) > 0 {
			if !trimStackIndex(n, stack) && len(*stack) > 1 && n.XMLName.Local == (*stack)[0] {
				*stack = (*stack)[1:]
			}
			if v := getValues(n, stack, func(n Node) []string { return n.values() }); v != nil {
				*values = append(*values, v...)
				return
			}
			if v := getValuesWithAttributes(n, stack, func(n Node, att string) []string { return n.valuesFromAttributes(att) }); v != nil {
				*values = append(*values, v...)
				return
			}
			xmlReader(n.Nodes, stack, values)
		}
	}
}

func trimStackIndex(n Node, stack *[]string) bool {
	re := regexp.MustCompile(`(.*)\[([0-9])\](.*)`)
	if re.MatchString((*stack)[0]) {
		submatches := re.FindStringSubmatch((*stack)[0])
		pre := submatches[1]
		index, err := strconv.Atoi(submatches[2])
		post := submatches[3]
		if err != nil {
			panic(err)
		}
		if n.XMLName.Local == pre && index == 0 {
			if len(*stack) == 1 {
				(*stack)[0] = fmt.Sprintf("%s%s", pre, post)
			} else {
				*stack = (*stack)[1:]
			}
			return true
		} else if n.XMLName.Local == pre && index > 0 {
			(*stack)[0] = fmt.Sprintf("%s[%d]%s", pre, index-1, post)
			return true
		}
	}
	return false
}

// getValues returns a slice of string holding the individual values inside the node
func getValues(n Node, stack *[]string, f func(Node) []string) []string {
	if len(*stack) == 1 && n.XMLName.Local == (*stack)[0] {
		*stack = []string{}
		return f(n)
	}
	return nil
}

// getValuesWithAttributes is like getValues but matches the node's attributes
func getValuesWithAttributes(n Node, stack *[]string, f func(Node, string) []string) []string {
	if len(*stack) == 1 {
		re := regexp.MustCompile(`(.*)@([a-zA-Z]+)`)
		if re.MatchString((*stack)[0]) {
			submatch := re.FindStringSubmatch((*stack)[0])
			if submatch[1] == n.XMLName.Local {
				*stack = []string{}
				return f(n, submatch[2])
			}
		}
	}
	return nil
}
