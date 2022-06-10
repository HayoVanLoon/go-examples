package main

import (
	"math"
)

type Group[K comparable, V any] interface {
	Key() K
	Values() []V
}

type group[K comparable, V any] struct {
	key    K
	values []V
}

func (g *group[K, V]) Key() K {
	return g.key
}

func (g *group[K, V]) Values() []V {
	return g.values
}

func DoMap[K comparable, V any](fn func(V) K, xs []V) []Group[K, V] {
	var groups []Group[K, V]
	m := make(map[K]int)
	for _, x := range xs {
		k := fn(x)
		idx, ok := m[k]
		if !ok {
			idx = len(groups)
			m[k] = idx
			groups = append(groups, &group[K, V]{key: k})
		}
		groups[idx].(*group[K, V]).values = append(groups[idx].(*group[K, V]).values, x)
	}
	return groups
}

func Merge[K comparable, V any](groups ...Group[K, V]) []Group[K, V] {
	var out []Group[K, V]
	m := make(map[K]int)
	for _, g := range groups {
		idx, ok := m[g.Key()]
		if ok {
			old := len(out[idx].Values())
			ng := &group[K, V]{g.Key(), make([]V, old+len(g.Values()))}
			copy(ng.values, out[idx].Values())
			copy(ng.values[old:], g.Values())
			out[idx] = ng
			continue
		}
		m[g.Key()] = len(out)
		out = append(out, g)
	}
	return out
}

type ShuffleGroup[K comparable, V any] struct {
	Group[K, V]
	Id int
}

func Shuffle[K comparable, V any](groups []Group[K, V], num int) []ShuffleGroup[K, V] {
	if num <= 0 {
		panic("num <= 0")
	}
	var keys []K
	var counts []int
	for _, g := range groups {
		keys = append(keys, g.Key())
		counts = append(counts, len(g.Values()))
	}

	m2 := make([]ShuffleGroup[K, V], len(groups))
	parts := partition(counts, num)
	for w := num - 1; w >= 0; w -= 1 {
		end := len(keys)
		if w < num-1 {
			end = parts[w+1]
		}
		for i := parts[w]; i < end; i += 1 {
			m2[i] = ShuffleGroup[K, V]{Id: w, Group: groups[i]}
		}
	}
	return m2
}

func partition(xs []int, num int) []int {
	type cell struct{ loss, idx int }
	m := make([][]cell, len(xs))
	for i := range xs {
		m[i] = make([]cell, num)
	}
	m[0][0].loss = xs[0]
	for i := 1; i < len(xs); i += 1 {
		m[i][0].loss = m[i-1][0].loss + xs[i]
	}

	for w := 1; w < num; w += 1 {
		for i := 1; i < len(xs); i += 1 {
			m[i][w].loss = math.MaxInt
			acc := 0
			for k := i; k >= w; k -= 1 {
				acc += xs[k]
				max := m[k-1][w-1].loss
				if max < acc {
					max = acc
				}
				if max < m[i][w].loss {
					m[i][w] = cell{max, k}
				}
			}
		}
	}

	positions := make([]int, num)
	row, col := len(xs)-1, num-1
	for {
		positions[col] = m[row][col].idx
		if row == 0 || col == 0 {
			break
		}
		row = positions[col] - 1
		col -= 1
	}
	return positions
}

func Reduce[K comparable, V any](fn func([]V) V, groups []Group[K, V]) map[K]V {
	m2 := make(map[K]V)
	for _, g := range groups {
		m2[g.Key()] = fn(g.Values())
	}
	return m2
}
