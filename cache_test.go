package cjsonsource

import "testing"

func TestMergeMap(t *testing.T) {
	a := map[string]string{"a": "a", "b": "b", "c": "c"}
	b := map[string]string{"b": "b", "c": "c", "d": "d"}
	x := map[string]string{"x": "x", "y": "y", "z": "z"}
	MergeMap(a, b)
	t.Log(a)
	MergeMap(b, x)
	t.Log(b)
}
