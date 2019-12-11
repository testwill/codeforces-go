package copypasta

import (
	"math/bits"
	"sort"
)

// https://www.csie.ntu.edu.tw/~hsinmu/courses/_media/dsa_13spring/horowitz_306_311_biconnected.pdf
// low(u) is the lowest dfn that we can reach from u using a path of descendants followed by at most one back edge

// namespace
type tree struct{}

// 树的重心 https://oi-wiki.org/graph/tree-centroid/
// 应用：求树上距离不超过 upperDis 的点对数 http://poj.org/problem?id=1741
func (*tree) numPairsWithDistanceLimit(n int, upperDis int64) int64 {
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	type neighbor struct {
		to     int
		weight int64
	}
	g := make([][]neighbor, n+1)
	for i := 0; i < n-1; i++ {
		var v, w int
		var weight int64
		//v, w, weight := read(), read(), read()
		g[v] = append(g[v], neighbor{w, weight})
		g[w] = append(g[w], neighbor{v, weight})
	}
	usedCentroid := make([]bool, n+1)

	subtreeSize := make([]int, n+1)
	var calcSubtreeSize func(v, fa int) int
	calcSubtreeSize = func(v, fa int) int {
		sz := 1
		for _, e := range g[v] {
			if w := e.to; w != fa && !usedCentroid[w] {
				sz += calcSubtreeSize(w, v)
			}
		}
		subtreeSize[v] = sz
		return sz
	}

	var findCentroid func(v, fa, compSize int) (int, int)
	findCentroid = func(v, fa, compSize int) (minSize, ct int) {
		minSize = int(1e9)
		//ct = -1
		maxSubSize := 0
		sizeV := 1 // 除去了 usedCentroid 子树的剩余大小
		for _, e := range g[v] {
			if w := e.to; w != fa && !usedCentroid[w] {
				if minSizeW, ctW := findCentroid(w, v, compSize); minSizeW < minSize {
					minSize = minSizeW
					ct = ctW
				}
				maxSubSize = max(maxSubSize, subtreeSize[w])
				sizeV += subtreeSize[w]
			}
		}
		maxSubSize = max(maxSubSize, compSize-sizeV)
		if maxSubSize < minSize {
			minSize = maxSubSize
			ct = v
		}
		return
	}

	var disToCentroid []int64
	var calcDisToCentroid func(v, fa int, d int64)
	calcDisToCentroid = func(v, fa int, d int64) {
		disToCentroid = append(disToCentroid, d)
		for _, e := range g[v] {
			if w := e.to; w != fa && !usedCentroid[w] {
				calcDisToCentroid(w, v, d+e.weight)
			}
		}
	}

	countPairs := func(ds []int64) int64 {
		cnt := int64(0)
		//sort.Ints(ds)
		sort.Slice(ds, func(i, j int) bool { return ds[i] < ds[j] })
		j := len(ds)
		for i, di := range ds {
			for ; j > 0 && di+ds[j-1] > upperDis; j-- {
			}
			cnt += int64(j)
			if j > i {
				cnt--
			}
		}
		return cnt >> 1
	}

	var f func(int) int64
	f = func(v int) (ans int64) {
		calcSubtreeSize(v, -1)
		_, ct := findCentroid(v, -1, subtreeSize[v])
		usedCentroid[ct] = true
		// 统计按 ct 分割后的子树中的点对数
		for _, e := range g[ct] {
			if !usedCentroid[e.to] {
				ans += f(e.to)
			}
		}
		// 统计经过 ct 的点对数
		// 0 是方便统计包含 ct 的部分
		ds := []int64{0}
		for _, e := range g[ct] {
			if !usedCentroid[e.to] {
				disToCentroid = []int64{}
				calcDisToCentroid(e.to, ct, e.weight)
				ans -= countPairs(disToCentroid)
				ds = append(ds, disToCentroid...)
			}
		}
		ans += countPairs(ds)
		usedCentroid[ct] = false
		return
	}
	return f(0)
}

// https://oi-wiki.org/graph/lca/#rmq
func (*tree) lca(n, root int) {
	g := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var v, w int
		//v, w := read()-1, read()-1
		g[v] = append(g[v], w)
		g[w] = append(g[w], v)
	}

	vs := make([]int, 0, 2*n-1) // 欧拉序列
	pos := make([]int, n)       // pos[v] 表示 v 在 vs 中第一次出现的位置编号
	depths := make([]int, 0, 2*n-1)
	var dfs func(v, fa, d int)
	dfs = func(v, fa, d int) {
		pos[v] = len(vs)
		vs = append(vs, v)
		depths = append(depths, d)
		for _, w := range g[v] {
			if w != fa {
				dfs(w, v, d+1)
				vs = append(vs, v)
				depths = append(depths, d)
			}
		}
	}
	dfs(root, -1, 0)

	type pair struct{ v, i int }
	var st [][20]pair
	stInit := func(a []int) {
		n := len(a)
		st = make([][20]pair, n)
		for i := range st {
			st[i][0] = pair{a[i], i}
		}
		for j := uint(1); 1<<j <= n; j++ {
			for i := 0; i+(1<<j)-1 < n; i++ {
				st0, st1 := st[i][j-1], st[i+(1<<(j-1))][j-1]
				if st0.v < st1.v {
					st[i][j] = st0
				} else {
					st[i][j] = st1
				}
			}
		}
	}
	stQuery := func(l, r int) int { // [l,r] 注意 l r 是从 0 开始算的
		k := uint(bits.Len(uint(r-l+1)) - 1)
		st0, st1 := st[l][k], st[r-(1<<k)+1][k]
		if st0.v < st1.v {
			return st0.i
		}
		return st1.i
	}
	// 注意下标的换算，输出 LCA 的话要 +1
	calcLCA := func(v, w int) int {
		pv, pw := pos[v], pos[w]
		if pv > pw {
			pv, pw = pw, pv
		}
		return vs[stQuery(pv, pw)]
	}

	stInit(depths)
	// ...

	_ = calcLCA
}

// https://en.wikipedia.org/wiki/Heavy_path_decomposition
// https://oi-wiki.org/graph/hld/
func (*tree) hld(n, root int) {
	read := func() int { return 0 }

	// TODO: 处理边权的情况
	g := make([][]int, n)
	for i := 0; i < n-1; i++ {
		v, w := read()-1, read()-1
		g[v] = append(g[v], w)
		g[w] = append(g[w], v)
	}
	// 点权
	vals := make([]int64, n)
	for i := range vals {
		vals[i] = int64(read())
	}

	// 重儿子，父节点，深度，子树大小，所处重链顶点（深度最小），DFS 序（作为线段树中的编号，从 1 开始）
	type node struct{ hson, fa, depth, size, top, dfn int }
	nodes := make([]node, n)
	//idv := make([]int, n+1) // idv[nodes[v].dfn] == v

	var build func(v, fa, d int) *node
	build = func(v, fa, d int) *node {
		nodes[v] = node{hson: -1, fa: fa, depth: d, size: 1}
		o := &nodes[v]
		for _, w := range g[v] {
			if w == fa {
				continue
			}
			son := build(w, v, d+1)
			o.size += son.size
			if o.hson == -1 || son.size > nodes[o.hson].size {
				o.hson = w
			}
		}
		return o
	}
	build(root, -1, 0)

	dfn := 0
	var decomposition func(v, fa, top int)
	decomposition = func(v, fa, top int) {
		o := &nodes[v]
		o.top = top
		dfn++
		o.dfn = dfn
		//idv[dfn] = v
		if o.hson != -1 {
			// 优先遍历重儿子，保证在同一条重链上的点的 DFS 序是连续的
			decomposition(o.hson, v, top)
			for _, w := range g[v] {
				if w != fa && w != o.hson {
					decomposition(w, v, w)
				}
			}
		}
	}
	decomposition(root, -1, root)

	t := make(lazySegmentTree, 4*n)
	// 按照 DFS 序初始化线段树
	dfnVals := make([]int64, n)
	for i, v := range vals {
		dfnVals[nodes[i].dfn-1] = v
	}
	t.init(dfnVals)

	doPath := func(v, w int, do func(l, r int)) {
		ov, ow := nodes[v], nodes[w]
		for ; ov.top != ow.top; ov, ow = nodes[v], nodes[w] {
			topv, topw := nodes[ov.top], nodes[ow.top]
			// v 所处的重链顶点必须比 w 的深
			if topv.depth < topw.depth {
				v, w = w, v
				ov, ow = ow, ov
				topv, topw = topw, topv
			}
			do(topv.dfn, ov.dfn)
			// TODO: 边权下，处理轻边的情况
			v = topv.fa
		}
		if ov.depth > ow.depth {
			//v, w = w, v
			ov, ow = ow, ov
		}
		do(ov.dfn, ow.dfn)
		// TODO: 边权下，处理轻边的情况
	}
	updatePath := func(v, w int, add int64) {
		doPath(v, w, func(l, r int) { t.update(l, r, add) })
	}
	queryPath := func(v, w int) (sum int64) {
		doPath(v, w, func(l, r int) { sum += t.query(l, r) }) // TODO % mod
		return
	}
	updateSubtree := func(v int, add int64) {
		o := nodes[v]
		t.update(o.dfn, o.dfn+o.size-1, add)
	}
	querySubtree := func(v int) (sum int64) {
		o := nodes[v]
		return t.query(o.dfn, o.dfn+o.size-1)
	}

	_ = []interface{}{updatePath, queryPath, updateSubtree, querySubtree}
}

// TODO LCT
