//基于dinic算法

package course

import (
	"Course-Selection-Scheduling/utils"
	"gopkg.in/eapache/queue.v1"
	"math"
)

type edge struct {
	v, w, next int //出点、权值、邻接表下一项
}
type point struct {
	ID  string //原始ID
	typ int    //类型（0:教师｜1:课程）
}

//图的一些信息
var (
	startID, endID, eSize int    //源点 汇点 中间边的数量
	dep, cur, head        []int  //深度 当前弧标记 邻接表头
	graph                 []edge //图
)

//映射解决教师ID与课程ID相同的情况
var (
	pointCnt   int               //点数
	pointID    = map[point]int{} //映射后的ID
	originalID []point           //原始ID
)

func add(u, v int) {
	graph = append(graph, edge{v, 1, head[u]})
	head[u] = len(graph) - 1
	graph = append(graph, edge{u, 0, head[v]})
	head[v] = len(graph) - 1
}

//构造分层
func bfs() bool {
	for pos, _ := range dep {
		dep[pos] = 0
	}
	dep[startID] = 1
	que := queue.New()
	que.Add(startID)
	cur[startID] = head[startID]
	for que.Length() > 0 {
		u := que.Remove().(int)
		for x := head[u]; x != 0; x = graph[x].next {
			v := graph[x].v
			w := graph[x].w
			if w != 0 && (dep[v] == 0) {
				cur[v] = head[v] //复原当前弧
				dep[v] = dep[u] + 1
				que.Add(v)
				if v == endID {
					return true
				}
			}
		}
	}
	return false
}

//找增广
func dfs(u, limit int) int {
	if u == endID {
		return limit
	}
	flow := 0
	for x := cur[u]; x != 0 && flow < limit; x = graph[x].next {
		cur[u] = x //当前弧优化
		v := graph[x].v
		w := graph[x].w
		if w != 0 && dep[v] == dep[u]+1 {
			k := dfs(v, utils.Min(limit-flow, w)) //增广流量
			if k != 0 {
				dep[v] = 0 //剪枝，去掉增广完的点
			}
			graph[x].w -= k
			graph[x^1].w += k
			flow += k
		}
	}
	return flow //增广流量
}

func dinicRun() int {
	maxflow := 0
	for bfs() {
		for flow := dfs(startID, math.MaxInt); flow != 0; flow = dfs(startID, math.MaxInt) {
			maxflow += flow
		}
	}
	return maxflow
}

func initDinic(request map[string][]string) {
	//离散化
	for u, courses := range request {
		pointU := point{u, 0}
		if pointID[pointU] == 0 {
			originalID = append(originalID, pointU)
			pointID[pointU] = pointCnt
			pointCnt++
		}
		for _, v := range courses {
			pointV := point{v, 1}
			if pointID[pointV] == 0 {
				originalID = append(originalID, pointV)
				pointID[pointV] = pointCnt
				pointCnt++
			}
		}
	}

	graph = make([]edge, 2) //前2个不使用
	dep = make([]int, pointCnt+2)
	cur = make([]int, pointCnt+2)
	head = make([]int, pointCnt+2)
	startID = pointCnt
	endID = pointCnt + 1

	for u, courses := range request {
		uID := pointID[point{u, 0}]
		for _, v := range courses {
			vID := pointID[point{v, 1}]
			add(uID, vID)
		}
	}
	eSize = len(graph)

	for id, pointInfo := range originalID {
		if pointInfo.typ == 0 {
			add(startID, id)
		} else {
			add(id, endID)
		}
	}
}

func dinic(request map[string][]string) map[string]string {
	initDinic(request)
	respond := make(map[string]string, dinicRun())
	//log.Printf("排课完成 总共%d节课\n", run())
	for x := 2; x < eSize; x += 2 {
		if graph[x].w == 0 {
			respond[originalID[graph[x^1].v].ID] = originalID[graph[x].v].ID
			//log.Printf("教师号：%s课程号：%s\n", originalID[graph[x^1].v].ID, originalID[graph[x].v].ID)
		}
	}
	return respond
}
