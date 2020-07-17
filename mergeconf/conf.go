package mergeconf

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

const (
	// Node shows an element is a node
	Node = iota
	// Leaf shows an element is a leaf
	Leaf
)

var (
	whiteSpaceChars = " \n\t"
)

type elem struct {
	kind     int
	name     string
	value    string
	children map[string]*elem
}

func newElem(kind int, name string) *elem {
	return &elem{kind, name, "", make(map[string]*elem)}
}

func (e *elem) setValue(value string) *elem {
	e.value = value
	return e
}

func (e *elem) addChild(name string, child *elem) *elem {
	e.children[name] = child
	return e
}

func (e *elem) findChild(name string) (ret *elem, ok bool) {
	ret, ok = e.children[name]
	return
}

func (e *elem) format(level int) string {
	buf := bytes.NewBuffer(nil)
	spaces := strings.Repeat("\t", level)

	if e.name != "" {
		buf.WriteString(fmt.Sprintf("%s<%s>\n", spaces, e.name))
	} else {
		level-- // for root elem
	}

	if len(e.children) > 0 {
		isFirst := true
		for _, v := range e.children {
			if v.kind == Leaf {
				if isFirst {
					isFirst = false
				} else {
					buf.WriteString("\n")
				}
				if v.value == "" {
					buf.WriteString(fmt.Sprintf("\t%s%s", spaces, v.name))
				} else {
					buf.WriteString(fmt.Sprintf("\t%s%s = %s", spaces, v.name, v.value))
				}
			}
		}
		for _, v := range e.children {
			if v.kind != Leaf {
				if isFirst {
					isFirst = false
				} else {
					buf.WriteString("\n")
				}
				buf.WriteString(v.format(level + 1))
			}
		}
	}
	if e.name != "" {
		buf.WriteString(fmt.Sprintf("\n%s</%s>", spaces, e.name))
	}
	return buf.String()
}

func (e *elem) String() string {
	if e.kind == Leaf {
		return fmt.Sprintf("%s = %s", e.name, e.value)
	}
	return e.format(0)
}

func initFromBytes(content []byte) (*elem, error) {
	root := &elem{children: make(map[string]*elem)}
	xmlDecoder := xml.NewDecoder(bytes.NewReader(content))
	var nodeStack []*elem
	nodeStack = append(nodeStack, root)
	for {
		currNode := nodeStack[len(nodeStack)-1]
		token, _ := xmlDecoder.Token()
		if token == nil {
			break
		}
		switch token.(type) {
		case xml.CharData:
			lineDecoder := bufio.NewScanner(bytes.NewReader(token.(xml.CharData)))
			lineDecoder.Split(bufio.ScanLines)
			for lineDecoder.Scan() {
				line := strings.Trim(lineDecoder.Text(), whiteSpaceChars)
				if len(line) > 0 && line[0] == '#' {
					continue
				}
				kv := strings.SplitN(line, "=", 2)
				k := strings.Trim(kv[0], whiteSpaceChars)
				v := ""
				if len(kv) == 2 {
					v = strings.Trim(kv[1], whiteSpaceChars)
				}
				if k == "" {
					continue
				}
				leaf := newElem(Leaf, k)
				leaf.setValue(v)
				currNode.addChild(k, leaf)
			}
		case xml.StartElement:
			nodeName := token.(xml.StartElement).Name.Local
			node, ok := currNode.findChild(nodeName)
			if !ok {
				node = newElem(Node, nodeName)
				currNode.addChild(nodeName, node)
			}
			nodeStack = append(nodeStack, node)
		case xml.EndElement:
			nodeName := token.(xml.EndElement).Name.Local
			if currNode.name != nodeName {
				return nil, fmt.Errorf("xml end not match :%s", nodeName)
			}
			nodeStack = nodeStack[:len(nodeStack)-1]
		}
	}
	return root, nil
}
