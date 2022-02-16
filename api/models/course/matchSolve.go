package course

import (
	"Course-Selection-Scheduling/utils"
	"gopkg.in/eapache/queue.v1"
	"math"
)

//点
type point struct {
	ID  string //原始ID
	typ int    //类型（0:教师｜1:课程）
}

//离散化 映射解决教师ID与课程ID相同的情况
var (
	pointCnt   int           //点数
	pointID    map[point]int //映射后的ID
	originalID []point       //原始ID
)

//边
type edge struct {
	v, w, next int //出点、权值、邻接表下一项
}

//图
var graph []edge

//邻接表添加边
func add(u, v, w int) {
	graph = append(graph, edge{v, w, head[u]})
	head[u] = len(graph) - 1
}

//离散化
func (scheduleInfo ScheduleCourseRequest) Discretize() int {
	//初始化
	pointCnt = 1 //留出0
	pointID = map[point]int{}
	originalID = make([]point, 1)
	originalID[0] = point{ID: "", typ: -1}
	request := scheduleInfo.TeacherCourseRelationShip
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
	return pointCnt
}

/**************************Dinic算法**************************/
//dinic辅助信息
var (
	startID, endID, eSize int   //源点 汇点 中间边的数量
	dep, cur, head        []int //深度 当前弧标记 邻接表头
)

//dinic初始化
func initDinic(request map[string][]string) {
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
			add(uID, vID, 1)
			add(vID, uID, 0)
		}
	}
	eSize = len(graph)

	for id, pointInfo := range originalID {
		if pointInfo.typ == 0 {
			add(startID, id, 1)
			add(id, startID, 0)
		} else if pointInfo.typ == 1 {
			add(id, endID, 1)
			add(endID, id, 0)
		}
	}
}

//dinic构造分层
func dinicBfs() bool {
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

//dinic找增广
func dinicDfs(u, limit int) int {
	if u == endID {
		return limit
	}
	flow := 0
	for x := cur[u]; x != 0 && flow < limit; x = graph[x].next {
		cur[u] = x //当前弧优化
		v := graph[x].v
		w := graph[x].w
		if w != 0 && dep[v] == dep[u]+1 {
			k := dinicDfs(v, utils.Min(limit-flow, w)) //增广流量
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

//dinic运行
func dinicRun() int {
	maxflow := 0
	for dinicBfs() {
		for flow := dinicDfs(startID, math.MaxInt); flow != 0; flow = dinicDfs(startID, math.MaxInt) {
			maxflow += flow
		}
	}
	return maxflow
}

//dinic入口
func (scheduleInfo ScheduleCourseRequest) Dinic() map[string]string {
	request := scheduleInfo.TeacherCourseRelationShip
	initDinic(request)
	respond := make(map[string]string, dinicRun())
	for x := 2; x < eSize; x += 2 {
		if graph[x].w == 0 {
			respond[originalID[graph[x^1].v].ID] = originalID[graph[x].v].ID
		}
	}
	return respond
}

/**************************匈牙利算法**************************/

var (
	match []int  //匹配项
	vis   []bool //访问标记
)

func initHun(request map[string][]string) {
	graph = make([]edge, 2) //前2个不使用
	head = make([]int, pointCnt+2)
	match = make([]int, pointCnt+2)
	vis = make([]bool, pointCnt+2)
	for u, courses := range request {
		uID := pointID[point{u, 0}]
		for _, v := range courses {
			vID := pointID[point{v, 1}]
			add(uID, vID, 1)
		}
	}
}

func hunDfs(u int) bool {
	for x := head[u]; x != 0; x = graph[x].next {
		v := graph[x].v
		if vis[v] {
			continue
		}
		vis[v] = true
		if match[v] == 0 || hunDfs(match[v]) {
			match[v] = u
			return true
		}
	}
	return false
}

func hunRun() int {
	res := 0
	for id, pointInfo := range originalID {
		if pointInfo.typ == 0 {
			for pos, _ := range vis {
				vis[pos] = false
			}
			if hunDfs(id) {
				res++
			}
		}
	}
	return res
}

func (scheduleInfo ScheduleCourseRequest) Hungarian() map[string]string {
	request := scheduleInfo.TeacherCourseRelationShip
	initHun(request)
	respond := make(map[string]string, hunRun())
	for id, pointInfo := range originalID {
		if pointInfo.typ == 1 && match[id] != 0 {
			respond[originalID[match[id]].ID] = originalID[id].ID
		}
	}
	return respond
}
