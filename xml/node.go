package xml

import (
	xmlencode "encoding/xml"
)

type Node struct {
	XMLName xmlencode.Name
	Content []byte           `xml:",innerxml"`
	Nodes   []Node           `xml:",any"`
	Attrs   []xmlencode.Attr `xml:"-"`
}

func (n *Node) UnmarshalXML(d *xmlencode.Decoder, start xmlencode.StartElement) error {
	n.Attrs = start.Attr
	type node Node

	return d.DecodeElement((*node)(n), &start)
}

func (n Node) values() []string {
	values := make([]string, 0)

	if len(n.Nodes) == 0 {
		values = append(values, string(n.Content))
	}

	for _, subnode := range n.Nodes {
		if len(subnode.Nodes) == 0 {
			values = append(values, string(subnode.Content))
		}
	}

	return values
}

func (n Node) valuesFromAttributes(att string) []string {
	values := make([]string, 0)

	if len(n.Nodes) == 0 {
		values = append(values, n.findAttributeValue(att))
	}

	for _, subnode := range n.Nodes {
		if len(subnode.Nodes) == 0 {
			values = append(values, n.findAttributeValue(att))
		}
	}

	return values
}

func (n Node) findAttributeValue(att string) string {
	for _, a := range n.Attrs {
		if a.Name.Local == att {
			return a.Value
		}
	}
	return ""
}
