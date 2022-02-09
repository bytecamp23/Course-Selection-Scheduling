//基于dinic算法
package course

import (
	"Course-Selection-Scheduling/utils"
	"bufio"
	"fmt"
	"gopkg.in/eapache/queue.v1"
	"math"
	"os"
)

type edge struct {
	u, v int
}
type graph struct {
	v, w, next int
}

var (
	st, ed, m, n, tot int
	dep, cur, head    []int
	edges             []edge
	e                 []graph
)

func add(u, v int) {
	e = append(e, graph{v, 1, head[u]})
	head[u] = len(e) - 1
	e = append(e, graph{u, 0, head[v]})
	head[v] = len(e) - 1
}
func input() {
	reader := bufio.NewReader(os.Stdin)
	reader = bufio.NewReaderSize(reader, 1e7)
	fmt.Fscanf(reader, "%d %d\n", &m, &n)
	var u, v int
	for {
		fmt.Fscanf(reader, "%d %d\n", &u, &v)
		if u == -1 && v == -1 {
			break
		}
		edges = append(edges, edge{u, v})
	}
}

//构造分层
func bfs() bool {
	for pos, _ := range dep {
		dep[pos] = 0
	}
	dep[st] = 1
	que := queue.New()
	que.Add(st)
	cur[st] = head[st]
	for que.Length() > 0 {
		u := que.Remove().(int)
		for x := head[u]; x != 0; x = e[x].next {
			v := e[x].v
			w := e[x].w
			if w != 0 && (dep[v] == 0) {
				cur[v] = head[v] //复原当前弧
				dep[v] = dep[u] + 1
				que.Add(v)
				if v == ed {
					return true
				}
			}
		}
	}
	return false
}

//找增广
func dfs(u, limit int) int {
	if u == ed {
		return limit
	}
	flow := 0
	for x := cur[u]; x != 0 && flow < limit; x = e[x].next {
		cur[u] = x //当前弧优化
		v := e[x].v
		w := e[x].w
		if w != 0 && dep[v] == dep[u]+1 {
			k := dfs(v, utils.Min(limit-flow, w)) //增广流量
			if k != 0 {
				dep[v] = 0 //剪枝，去掉增广完的点
			}
			e[x].w -= k
			e[x^1].w += k
			flow += k
		}
	}
	return flow //增广流量
}

func dinic() int {
	maxflow := 0
	for bfs() {
		for flow := dfs(st, math.MaxInt); flow != 0; flow = dfs(st, math.MaxInt) {
			maxflow += flow
		}
	}
	return maxflow
}

/*func main() {
	input()
	//st:=0 ed:=n+1 teacher:=[1,m],course:=[m+1,n]
	dep = make([]int, n+2)
	cur = make([]int, n+2)
	head = make([]int, n+2)
	e = make([]graph, 2) //前2个不使用
	st = 0
	ed = n + 1
	for v := 1; v <= m; v++ {
		add(st, v)
	}
	for u := m + 1; u <= n; u++ {
		add(u, ed)
	}
	for _, edge := range edges {
		add(edge.u, edge.v)
	}
	fmt.Println(dinic())
	for x := 2; x < len(e); x += 2 {
		if e[x].v >= m+1 && e[x].v <= n && e[x].w == 0 {
			fmt.Printf("%d %d\n", e[x^1].v, e[x].v)
		}
	}
}*/
