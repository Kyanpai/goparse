package xml

import (
	"bytes"
	xmlencode "encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Parse returns a string representing the content of the element(s) found at the path position
//
// If more than one element is found, Parse will join them with a space separator
func Parse(path string, data []byte) string {
	var n Node

	stack := strings.Split(path, "/")

	buf := bytes.NewBuffer(data)

	err := xmlencode.NewDecoder(buf).Decode(&n)
	if err != nil {
		panic(err)
	}
	values := make([]string, 0)
	xmlreader([]Node{n}, &stack, &values)
	return strings.Join(values, " ")
}

// xmlreader is a recursive function that will unstack an array of elements
// until reaching the desired node
func xmlreader(nodes []Node, stack *[]string, values *[]string) {
	for _, n := range nodes {
		if len(*stack) > 0 {
			if !trimStackIndex(n, stack) && len(*stack) > 1 && n.XMLName.Local == (*stack)[0] {
				*stack = (*stack)[1:]
			}
			if v := getValues(n, stack); v != nil {
				*values = append(*values, v...)
				return
			}
			if v := getValuesWithAttributes(n, stack); v != nil {
				*values = append(*values, v...)
				return
			}
			xmlreader(n.Nodes, stack, values)
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
func getValues(n Node, stack *[]string) []string {
	if len(*stack) == 1 && n.XMLName.Local == (*stack)[0] {
		*stack = []string{}
		return n.values()
	}
	return nil
}

// getValuesWithAttributes is like getValues but matches the node's attributes
func getValuesWithAttributes(n Node, stack *[]string) []string {
	if len(*stack) == 1 {
		re := regexp.MustCompile(`(.*)@([a-zA-Z]+)`)
		if re.MatchString((*stack)[0]) {
			submatch := re.FindStringSubmatch((*stack)[0])
			if submatch[1] == n.XMLName.Local {
				*stack = []string{}
				return n.valuesFromAttributes(submatch[2])
			}
		}
	}
	return nil
}
