package rbtree_test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/gosagawa/rbtree"
)

func TestMain(t *testing.T) {

	rand.Seed(time.Now().UnixNano())

	n := 30
	m := rbtree.NewRBMAP()
	keys := make([]int, n)
	for i := 0; i < n; i++ {
		keys[i] = i
	}
	rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
	for i := 0; i < n; i++ {
		m.Insert(i, keys[i])
	}

	for i := 0; i < 5; i++ {
		val := m.Lookup(i)
		fmt.Printf("m.lookup(%2d) == %2d\n", i, val)
	}
	fmt.Printf("size: %v\n", m.Size())
	fmt.Printf("keys: %v\n", m.Keys())
	fmt.Printf("%v\n", m)

	N := 1000000
	m.Clear()
	answer := make(map[int]int)

	insertErrors := 0
	deleteErrors := 0
	for i := 0; i < N; i++ {
		key := rand.Intn(math.MaxInt64)
		m.Insert(key, i)
		answer[key] = i
	}
	for key := range answer {
		x := m.Lookup(key)
		y := answer[key]
		if x != y {
			insertErrors++
		}
	}
	answerKeys := make([]int, len(answer))
	i := 0
	for key := range answer {
		answerKeys[i] = key
		i++
	}

	half := len(answer) / 2
	for _, key := range answerKeys {
		if half == 0 {
			break
		}
		m.Delete(key)
		half--
	}
	half = len(answer) / 2
	for _, key := range answerKeys {
		if half == 0 {
			break
		}
		if m.Member(key) {
			deleteErrors++
		}
		half--
	}
	if !m.Balanced() {
		t.Error("バランス: NG")
	}
	if !m.BstValid() {
		t.Error("２分探索木: NG")
	}
	if !m.ColorValid() {
		t.Error("配色: NG")
	}
	if insertErrors != 0 {
		t.Error("挿入: NG")
	}
	if deleteErrors != 0 {
		t.Error("削除: NG")
	}
	for _, key := range m.Keys() {
		m.Delete(key)
	}
	if !m.IsEmpty() {
		t.Error("全削除: NG")
	}
}
