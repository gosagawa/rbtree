package rbtree

import (
	"fmt"
)

// K:キーの型, V:値の型
func NewRBMAP() *RBMAP {
	return &RBMAP{}
}

type Color int

///////////////////////////////////////////////////////////////////////////
// 共通定義
///////////////////////////////////////////////////////////////////////////

// R:赤, B:黒, Error:debug 用
const (
	ColorR Color = iota
	ColorB
	ColorError
)

type Node struct { // ノードの型
	color Color // そのノードの色
	key   int   // そのノードのキー
	value int   // そのノードの値
	lst   *Node // 左部分木
	rst   *Node // 右部分木
}

func NewNode(color Color, key int, value int) *Node {
	return &Node{
		color: color,
		key:   key,
		value: value,
	}
}

type RBMAP struct {
	root   *Node // 赤黒木の根
	change bool  // 修正が必要かを示すフラグ(true:必要, false:不要)
	lmax   int   // 左部分木のキーの最大値
	value  int   // lmax に対応する値
}

// ノード n が赤かチェックする
func (n *Node) isR() bool {
	return n != nil && n.color == ColorR
}

// ノード n が黒かチェックする
func (n *Node) isB() bool {
	return n != nil && n.color == ColorB
}

// ２分探索木 v の左回転。回転した木を返す
func rotateL(v *Node) *Node {
	u := v.rst
	t2 := u.lst
	u.lst = v
	v.rst = t2
	return u
}

// ２分探索木 u の右回転。回転した木を返す
func rotateR(u *Node) *Node {
	v := u.lst
	t2 := v.rst
	v.rst = u
	u.lst = t2
	return v
}

// ２分探索木 t の二重回転(左回転 -> 右回転)。回転した木を返す
func rotateLR(t *Node) *Node {
	t.lst = rotateL(t.lst)
	return rotateR(t)
}

// ２分探索木 t の二重回転(右回転 -> 左回転)。回転した木を返す
func rotateRL(t *Node) *Node {
	t.rst = rotateR(t.rst)
	return rotateL(t)
}

///////////////////////////////////////////////////////////////////////////
// insert(挿入)
///////////////////////////////////////////////////////////////////////////

// エントリー(key, x のペア)を挿入する
func (m *RBMAP) Insert(key int, x int) {
	m.root = m.insertSub(m.root, key, x)
	m.root.color = ColorB
}

func (m *RBMAP) insertSub(t *Node, key int, x int) *Node {
	if t == nil {
		m.change = true
		return NewNode(ColorR, key, x)
	}
	cmp := 0
	if key > t.key {
		cmp = 1
	}
	if key < t.key {
		cmp = -1
	}
	if cmp < 0 {
		t.lst = m.insertSub(t.lst, key, x)
		return m.balance(t)
	} else if cmp > 0 {
		t.rst = m.insertSub(t.rst, key, x)
		return m.balance(t)
	}
	m.change = false
	t.value = x
	return t
}

// エントリー挿入に伴う赤黒木の修正(パターンマッチ)
func (m *RBMAP) balance(t *Node) *Node {
	if !m.change {
		return t
	} else if !t.isB() {
		return t // 根が黒でないなら何もしない
	} else if t.lst.isR() && t.lst.lst.isR() {
		t = rotateR(t)
		t.lst.color = ColorB
	} else if t.lst.isR() && t.lst.rst.isR() {
		t = rotateLR(t)
		t.lst.color = ColorB
	} else if t.rst.isR() && t.rst.lst.isR() {
		t = rotateRL(t)
		t.rst.color = ColorB
	} else if t.rst.isR() && t.rst.rst.isR() {
		t = rotateL(t)
		t.rst.color = ColorB
	} else {
		m.change = false
	}
	return t
}

///////////////////////////////////////////////////////////////////////////
// delete(削除)
///////////////////////////////////////////////////////////////////////////

// key で指すエントリー(ノード)を削除する
func (m *RBMAP) Delete(key int) {
	m.root = m.deleteSub(m.root, key)
	if m.root != nil {
		m.root.color = ColorB
	}
}

func (m *RBMAP) deleteSub(t *Node, key int) *Node {
	if t == nil {
		m.change = false
		return nil
	}
	cmp := 0
	if key > t.key {
		cmp = 1
	}
	if key < t.key {
		cmp = -1
	}
	if cmp < 0 {
		t.lst = m.deleteSub(t.lst, key)
		return m.balanceL(t)
	} else if cmp > 0 {
		t.rst = m.deleteSub(t.rst, key)
		return m.balanceR(t)
	} else {
		if t.lst == nil {
			switch t.color {
			case ColorR:
				m.change = false
				break
			case ColorB:
				m.change = true
				break
			}
			return t.rst // 右部分木を昇格させる
		} else {
			t.lst = m.deleteMax(t.lst) // 左部分木の最大値を削除する
			t.key = m.lmax             // 左部分木の削除した最大値で置き換える
			t.value = m.value
			return m.balanceL(t)
		}
	}
}

// 部分木 t の最大値のノードを削除する
// 戻り値は削除により修正された部分木
// 削除した最大値を lmax に保存する
func (m *RBMAP) deleteMax(t *Node) *Node {
	if t.rst != nil {
		t.rst = m.deleteMax(t.rst)
		return m.balanceR(t)
	} else {
		m.lmax = t.key // 部分木のキーの最大値を保存
		m.value = t.value
		switch t.color {
		case ColorR:
			m.change = false
			break
		case ColorB:
			m.change = true
			break
		}
		return t.lst // 左部分木を昇格させる
	}
}

// 左部分木のノード削除に伴う赤黒木の修正(パターンマッチ)
// 戻り値は修正された木
func (m *RBMAP) balanceL(t *Node) *Node {
	if !m.change {
		return t // 修正なしと赤ノード削除の場合はここ

	} else if t.rst.isB() && t.rst.lst.isR() {
		rb := t.color
		t = rotateRL(t)
		t.color = rb
		t.lst.color = ColorB
		m.change = false
	} else if t.rst.isB() && t.rst.rst.isR() {
		rb := t.color
		t = rotateL(t)
		t.color = rb
		t.lst.color = ColorB
		t.rst.color = ColorB
		m.change = false
	} else if t.rst.isB() {
		rb := t.color
		t.color = ColorB
		t.rst.color = ColorR
		m.change = (rb == ColorB)
	} else if t.rst.isR() {
		t = rotateL(t)
		t.color = ColorB
		t.lst.color = ColorR
		t.lst = m.balanceL(t.lst)
		m.change = false
	} else { // 黒ノード削除の場合、ここはありえない
		panic("(L) This program is buggy")
	}
	return t
}

// 右部分木のノード削除に伴う赤黒木の修正(パターンマッチ)
// 戻り値は修正された木
func (m *RBMAP) balanceR(t *Node) *Node {
	if !m.change {
		return t // 修正なしと赤ノード削除の場合はここ
	} else if t.lst.isB() && t.lst.rst.isR() {
		rb := t.color
		t = rotateLR(t)
		t.color = rb
		t.rst.color = ColorB
		m.change = false
	} else if t.lst.isB() && t.lst.lst.isR() {
		rb := t.color
		t = rotateR(t)
		t.color = rb
		t.lst.color = ColorB
		t.rst.color = ColorB
		m.change = false
	} else if t.lst.isB() {
		rb := t.color
		t.color = ColorB
		t.lst.color = ColorR
		m.change = (rb == ColorB)
	} else if t.lst.isR() {
		t = rotateR(t)
		t.color = ColorB
		t.rst.color = ColorR
		t.rst = m.balanceR(t.rst)
		m.change = false
	} else { // 黒ノード削除の場合、ここはありえない
		panic("(R) This program is buggy")
	}
	return t
}

///////////////////////////////////////////////////////////////////////////
// member(検索)等
///////////////////////////////////////////////////////////////////////////

// キーの検索。ヒットすれば true、しなければ false
func (m *RBMAP) Member(key int) bool {
	t := m.root
	for t != nil {
		cmp := 0
		if key > t.key {
			cmp = 1
		}
		if key < t.key {
			cmp = -1
		}
		if cmp < 0 {
			t = t.lst
		} else if cmp > 0 {
			t = t.rst
		} else {
			return true
		}
	}
	return false
}

// キーから値を得る。キーがヒットしない場合は nil を返す
func (m *RBMAP) Lookup(key int) int {
	t := m.root
	for t != nil {
		cmp := 0
		if key > t.key {
			cmp = 1
		}
		if key < t.key {
			cmp = -1
		}
		if cmp < 0 {
			t = t.lst
		} else if cmp > 0 {
			t = t.rst
		} else {
			return t.value
		}
	}
	return 0
}

// マップが空なら true、空でないなら false
func (m *RBMAP) IsEmpty() bool {
	return m.root == nil
}

// マップを空にする
func (m *RBMAP) Clear() {
	m.root = nil
}

// キーのリスト
func (m *RBMAP) Keys() []int {
	al := []int{}
	al = m.keysSub(m.root, al)
	return al
}

// 値のリスト
func (m *RBMAP) Values() []int {
	al := []int{}
	al = m.valuesSub(m.root, al)
	return al
}

// マップのサイズ
func (m *RBMAP) Size() int {
	return len(m.Keys())
}

// キーの最小値
func (m *RBMAP) Min() int {
	t := m.root
	if t == nil {
		return 0
	}
	for t.lst != nil {
		t = t.lst
	}
	return t.key
}

// キーの最大値
func (m *RBMAP) Max() int {
	t := m.root
	if t == nil {
		return 0
	}
	for t.rst != nil {
		t = t.rst
	}
	return t.key
}

func (m *RBMAP) keysSub(t *Node, al []int) []int {
	if t != nil {
		al = m.keysSub(t.lst, al)
		al = append(al, t.key)
		al = m.keysSub(t.rst, al)
	}
	return al
}

func (m *RBMAP) valuesSub(t *Node, al []int) []int {
	if t != nil {
		al = m.valuesSub(t.lst, al)
		al = append(al, t.value)
		al = m.valuesSub(t.rst, al)
	}
	return al
}

///////////////////////////////////////////////////////////////////////////
// debug 用ルーチン
///////////////////////////////////////////////////////////////////////////

// 赤黒木をグラフ文字列に変換する
func (m *RBMAP) String() string {
	return m.toGraph("", "", m.root)
}

// 赤黒木のバランスが取れているか確認する
func (m *RBMAP) Balanced() bool {
	return m.blackHeight(m.root) >= 0
}

// 赤黒木の配色が正しいか確認する
func (m *RBMAP) ColorValid() bool {
	return m.colorChain(m.root) == ColorB
}

// ２分探索木になっているか確認する
func (m *RBMAP) BstValid() bool {
	return m.bstValidSub(m.root)
}

func (m *RBMAP) toGraph(head string, bar string, t *Node) string {
	graph := ""
	if t != nil {
		graph += m.toGraph(head+"　　", "／", t.rst)
		node := ""
		switch t.color {
		case ColorR:
			node = "R"
		case ColorB:
			node = "B"
		}
		node += fmt.Sprintf(":%v:%v", t.key, t.value)
		graph += fmt.Sprintf("%v%v%v\n", head, bar, node)
		graph += m.toGraph(head+"　　", "＼", t.lst)
	}
	return graph
}

func (m *RBMAP) blackHeight(t *Node) int {
	if t == nil {
		return 0
	}
	nl := m.blackHeight(t.lst)
	nr := m.blackHeight(t.rst)
	if nl < 0 || nr < 0 || nl != nr {
		return -1
	}
	if t.color == ColorB {
		return nl + 1
	}
	return nl
}

func (m *RBMAP) colorChain(t *Node) Color {
	if t == nil {
		return ColorB
	}
	p := t.color
	cl := m.colorChain(t.lst)
	cr := m.colorChain(t.rst)
	if cl == ColorError || cr == ColorError {
		return ColorError
	}
	if p == ColorB {
		return p
	}
	if p == ColorR && cl == ColorB && cr == ColorB {
		return p
	}
	return ColorError
}

func (m *RBMAP) bstValidSub(t *Node) bool {
	if t == nil {
		return true
	}
	flag := m.small(t.key, t.lst) && m.large(t.key, t.rst)
	return flag && m.bstValidSub(t.lst) && m.bstValidSub(t.rst)
}

func (m *RBMAP) small(key int, t *Node) bool {
	if t == nil {
		return true
	}
	flag := key > t.key
	return flag && m.small(key, t.lst) && m.small(key, t.rst)
}

func (m *RBMAP) large(key int, t *Node) bool {
	if t == nil {
		return true
	}
	flag := key < t.key
	return flag && m.large(key, t.lst) && m.large(key, t.rst)
}
