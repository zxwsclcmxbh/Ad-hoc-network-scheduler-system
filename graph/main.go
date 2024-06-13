package main

import (
	"fmt"
	"os"
	"os/exec"
)

// Graph 图结构
type Graph struct {
	adjacency map[string][]string
}

// NewGraph 创建新的图
func NewGraph() *Graph {
	return &Graph{
		adjacency: make(map[string][]string),
	}
}

// AddEdge 添加边
func (g *Graph) AddEdge(u, v string) {
	g.adjacency[u] = append(g.adjacency[u], v)
}

// PrintDot 打印 DOT 格式
func (g *Graph) PrintDot() {
	fmt.Println("graph G {")
	for node, neighbors := range g.adjacency {
		for _, neighbor := range neighbors {
			fmt.Printf("  %s -- %s;\n", node, neighbor)
		}
	}
	fmt.Println("}")
}

func main() {
	graph := NewGraph()

	// 添加节点的连接关系
	graph.AddEdge("EdgeNode1", "EdgeNode2")
	graph.AddEdge("EdgeNode2", "EdgeNode3")
	graph.AddEdge("EdgeNode3", "EdgeNode4")
	// graph.AddEdge("EdgeNode4", "EdgeNode1")
	graph.AddEdge("EdgeNode1", "EdgeNode3")
	// graph.AddEdge("EdgeNode2", "EdgeNode4")

	// 打印 DOT 格式
	graph.PrintDot()

	// 将 DOT 格式保存到文件
	file, err := os.Create("graph.dot")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	file.WriteString("graph G {\n")
	for node, neighbors := range graph.adjacency {
		for _, neighbor := range neighbors {
			file.WriteString(fmt.Sprintf("  %s -- %s;\n", node, neighbor))
		}
	}
	file.WriteString("}\n")

	// 使用 Graphviz 生成图形
	cmd := exec.Command("dot", "-Tpng", "graph.dot", "-o", "graph.png")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error generating graph:", err)
		return
	}
	fmt.Println("Graph generated as graph.png")
}
