package copypasta

import (
	. "fmt"
	"index/suffixarray"
	"reflect"
	"unsafe"
)

func hashCollection() {
	const prime uint64 = 1e8 + 7

	var maxLen int
	powP := make([]uint64, maxLen+1)
	powP[0] = 1
	for i := 1; i <= maxLen; i++ {
		powP[i] = powP[i-1] * prime
	}

	hashVal := func(s string) (val uint64) {
		for i, c := range s {
			val += uint64(c) * powP[i]
		}
		return
	}

	_ = hashVal
}

func stringCollection() {
	// code from my answer at https://www.zhihu.com/question/21923021/answer/37475572
	calcMaxMatchLengths := func(pattern string) []int {
		n := len(pattern)
		maxMatchLengths := make([]int, n)
		maxLength := 0
		for i := 1; i < n; i++ {
			c := pattern[i]
			for maxLength > 0 && pattern[maxLength] != c {
				maxLength = maxMatchLengths[maxLength-1]
			}
			if pattern[maxLength] == c {
				maxLength++
			}
			maxMatchLengths[i] = maxLength
		}
		return maxMatchLengths
	}
	// search pattern from text, return all start positions
	kmpSearch := func(text, pattern string) (positions []int) {
		maxMatchLengths := calcMaxMatchLengths(pattern)
		lenP := len(pattern)
		count := 0
		for i := range text {
			c := text[i]
			for count > 0 && pattern[count] != c {
				count = maxMatchLengths[count-1]
			}
			if pattern[count] == c {
				count++
			}
			if count == lenP {
				positions = append(positions, i-lenP+1)
				count = maxMatchLengths[count-1]
			}
		}
		return
	}

	calcMinPeriod := func(pattern string) int {
		maxMatchLengths := calcMaxMatchLengths(pattern)
		n := len(pattern)
		if val := maxMatchLengths[n-1]; val > 0 {
			if n%(n-val) == 0 {
				return n / (n - val)
			}
		}
		return 1 // or -1
	}

	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	max := func(a, b int) int {
		if a >= b {
			return a
		}
		return b
	}
	smallestRepresentation := func(s string) string {
		n := len(s)
		s += s
		i := 0
		for j := 1; j < n; {
			k := 0
			for ; k < n && s[i+k] == s[j+k]; k++ {
			}
			if k >= n {
				break
			}
			if s[i+k] < s[j+k] {
				j += k + 1
			} else {
				i, j = j, max(j, i+k)+1
			}
		}
		return s[i : i+n]
	}

	// https://blog.csdn.net/synapse7/article/details/18908413
	// http://manacher-viz.s3-website-us-east-1.amazonaws.com/#/
	manacher := func(origin string) int {
		n := len(origin)
		s := make([]byte, 2*n+3)
		s[0] = '^'
		for i, c := range origin {
			s[i<<1|1] = '#'
			s[i<<1+2] = byte(c)
		}
		s[n<<1|1] = '#'
		s[n<<1+2] = '$'
		maxLen := make([]int, 2*n+3)
		var ans, mid, right int
		for i := 1; i < 2*n+2; i++ {
			if i < right {
				maxLen[i] = min(maxLen[mid<<1-i], right-i)
			} else {
				maxLen[i] = 1
			}
			// 取 min 的原因：记点 i 关于 mid 的对称点为 i'，
			// 若以 i' 为中心的回文串范围超过了以 mid 为中心的回文串的范围
			// (此时有 i + len[(mid<<1)-i] >= right，注意 len 是包括中心的半长度)
			// 则 len[i] 应取 right - i (总不能超过边界吧)
			for s[i+maxLen[i]] == s[i-maxLen[i]] {
				maxLen[i]++
			}
			ans = max(ans, maxLen[i])
			if right < i+maxLen[i] {
				mid = i
				right = i + maxLen[i]
			}
		}
		return ans - 1
	}

	// 判断源串中的某一子串 [l...r] 是否为回文串
	//manacherQuery := func(l, r int) bool {
	//return maxLen[l+r+2] >= r-l+1
	//}

	_ = []interface{}{kmpSearch, calcMinPeriod, smallestRepresentation, manacher}
}

func suffixArray() {
	// lcp[i] = lcp(s[sa[i]:], s[sa[i+1]:])
	calcLCP := func(s []byte, sa []int) (lcp []int) {
		n := len(s)
		rank := make([]int, n+1)
		for i := range rank {
			rank[sa[i]] = i
		}
		lcp = make([]int, n, n+1)
		h := 0
		for i := range lcp {
			j := sa[rank[i]-1]
			if h > 0 {
				h--
			}
			for ; j+h < n && i+h < n && s[j+h] == s[i+h]; h++ {
			}
			lcp[rank[i]-1] = h
		}
		return
	}

	var s []byte
	index := suffixarray.New(s)
	sa := *(*[]int)(unsafe.Pointer(reflect.ValueOf(index).Elem().FieldByName("sa").UnsafeAddr()))
	sa = append([]int{len(s)}, sa...) // 方便定义 lcp
	lcp := calcLCP(s, sa)
	lcp = append(lcp, 0)

	// debug
	for i := range sa {
		if lcp[i] == 0 {
			Println("  " + string(s[sa[i]:]))
		} else {
			Println(lcp[i], string(s[sa[i]:]))
		}
	}

	// TODO: []int 的后缀数组
}

//

type trieNode struct {
	childIdx [26]int
	dupCnt   int // 重复插入计数
	val      int // 节点附加信息
}
type trie struct {
	nodes []*trieNode
}

func newTrie() *trie {
	return &trie{
		nodes: []*trieNode{{}}, // init with root
	}
}

// insert `s` into trie and add `val` at leaf
func (t *trie) insert(s string, val int) {
	o := t.nodes[0]
	for _, c := range s {
		c -= 'a'
		if o.childIdx[c] == 0 {
			o.childIdx[c] = len(t.nodes)
			t.nodes = append(t.nodes, &trieNode{})
		}
		o = t.nodes[o.childIdx[c]]
	}
	o.dupCnt++
	o.val = val
}

func (t *trie) contain(s string) bool {
	o := t.nodes[0]
	for _, c := range s {
		idx := o.childIdx[c-'a']
		if o.childIdx[idx] == 0 {
			return false
		}
		o = t.nodes[idx]
	}
	return true
}

// childIdx 长度为 2，且 trie 上所有字符串长度与 bits 一致
func (t *trie) maxXor(bits []byte) {
	o := t.nodes[0]
	for i, b := range bits {
		if o.childIdx[b^1] > 0 {
			bits[i] = 1
			b ^= 1
		} else { // o.childIdx[b] > 0
			bits[i] = 0
		}
		o = t.nodes[o.childIdx[b]]
	}
	return
}
