package dag

type GraphList struct {
	adjList map[int][]int
}

func NewGraphList() *GraphList {
	return &GraphList{
		adjList: make(map[int][]int),
	}
}

func (g *GraphList) AddVertex(vertex int) {
	if _, ok := g.adjList[vertex]; !ok {
		g.adjList[vertex] = []int{}
	}
}

func (g *GraphList) AddEdge(from, to int) {
	g.AddVertex(from)
	g.AddVertex(to)
	g.adjList[from] = append(g.adjList[from], to)
}

func (g *GraphList) TopologicalSort() ([]int, bool) {
	var sortedVertexes []int
	var queue []int

	//Key is Vertex and value is the number of a Vertexes Pointing to that Vertex
	indegrees := make(map[int]int)

	//Filling Indegree
	for vertex, neighbours := range g.adjList {
		if _, ok := indegrees[vertex]; !ok {
			indegrees[vertex] = 0
		}
		for _, neighbourVertex := range neighbours {
			if _, ok := indegrees[neighbourVertex]; !ok {
				indegrees[neighbourVertex] = 1
			} else {
				indegrees[neighbourVertex]++
			}
		}
	}

	//Adding Vertexes with 0 indegree to queue i.e. no other vertex pointing towards it
	for attributeId, indegree := range indegrees {
		if indegree == 0 {
			queue = append(queue, attributeId)
		}
	}

	for len(queue) > 0 {
		top := queue[0]
		queue = queue[1:]
		sortedVertexes = append(sortedVertexes, top)

		for _, neighbourVertex := range g.adjList[top] {
			indegrees[neighbourVertex]--
			if indegrees[neighbourVertex] == 0 {
				queue = append(queue, neighbourVertex)
			}
		}
	}
	if len(sortedVertexes) != len(indegrees) {
		return []int{}, true
	}
	return sortedVertexes, false
}
