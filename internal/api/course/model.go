package course

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

//离散化
func discretize(request map[string][]string) {
	//初始化
	pointCnt = 1 //留出0
	pointID = map[point]int{}
	originalID = make([]point, 1)
	originalID[0] = point{ID: "", typ: -1}

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
}

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
