package listedit_test

import (
	. "gopkg.in/check.v1"

	"testing"

	"github.com/canonical/go-algo/listedit"
)

type editTest struct {
	a, b string
	f    listedit.CostFunc
	eq   listedit.EqualFunc
	r    int64
}

type multiEditTest struct {
	a  []string
	b  string
	f  listedit.CostFunc
	eq listedit.EqualFunc
	r  int64
}

func uniqueCost(ar, br any) listedit.Cost {
	return listedit.Cost{SwapAB: 1, DeleteA: 3, InsertB: 5}
}

var editTests = []editTest{
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 0, a: "abc", b: "abc"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 1, a: "abc", b: "abd"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 1, a: "abc", b: "adc"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 1, a: "abc", b: "dbc"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 2, a: "abc", b: "add"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 2, a: "abc", b: "ddc"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 2, a: "abc", b: "dbd"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 3, a: "abc", b: "ddd"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 3, a: "abc", b: "ab"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 3, a: "abc", b: "bc"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 3, a: "abc", b: "ac"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 6, a: "abc", b: "a"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 6, a: "abc", b: "b"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 6, a: "abc", b: "c"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 9, a: "abc", b: ""},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 5, a: "abc", b: "abcd"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 5, a: "abc", b: "dabc"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 10, a: "abc", b: "adbdc"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 10, a: "abc", b: "dabcd"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 40, a: "abc", b: "ddaddbddcdd"},
	{f: listedit.StandardCost, eq: listedit.DefaultEqual, r: 3, a: "abcdefg", b: "axcdfgh"},
	{f: listedit.StandardCost, eq: listedit.DefaultEqual, r: 3, a: "aaaabcaaa", b: "aaaaafa"},
	{f: listedit.StandardCost, eq: listedit.DefaultEqual, r: 5, a: "", b: "abcde"},
	{f: listedit.StandardCost, eq: listedit.DefaultEqual, r: 5, a: "abcde", b: ""},
}

func (s *S) TestEdit(c *C) {
	for _, test := range editTests {
		c.Logf("Test: %v", test)
		f := test.f
		if f == nil {
			f = listedit.StandardCost
		}
		eq := test.eq
		if eq == nil {
			eq = listedit.DefaultEqual
		}
		alist := splitString(test.a)
		blist := splitString(test.b)
		r, ops := listedit.Edit(alist, blist, f, eq)
		c.Assert(r, Equals, test.r)
		c.Assert(applyOps(alist, ops), DeepEquals, blist)
	}
}

var multiEditTests = []multiEditTest{
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 0, a: []string{"abc"}, b: "abc"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 0, a: []string{"ab", "c"}, b: "abc"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 0, a: []string{"a", "b", "c"}, b: "abc"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 0, a: []string{"a", "bc"}, b: "abc"},

	{f: uniqueCost, eq: listedit.DefaultEqual, r: 0, a: []string{"abc", ""}, b: "abc"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 0, a: []string{"ab", "", "c"}, b: "abc"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 0, a: []string{"ab", "c", ""}, b: "abc"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 0, a: []string{"a", "", "b", "", "c"}, b: "abc"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 0, a: []string{"a", "", "bc"}, b: "abc"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 0, a: []string{"a", "bc", ""}, b: "abc"},
	{f: uniqueCost, eq: listedit.DefaultEqual, r: 0, a: []string{"", "abc"}, b: "abc"},

	{f: listedit.StandardCost, eq: listedit.DefaultEqual, r: 3, a: []string{"abcdefg"}, b: "axcdfgh"},
	{f: listedit.StandardCost, eq: listedit.DefaultEqual, r: 3, a: []string{"a", "bcdefg"}, b: "axcdfgh"},
	{f: listedit.StandardCost, eq: listedit.DefaultEqual, r: 3, a: []string{"ab", "cdefg"}, b: "axcdfgh"},
	{f: listedit.StandardCost, eq: listedit.DefaultEqual, r: 3, a: []string{"a", "b", "cdefg"}, b: "axcdfgh"},
	{f: listedit.StandardCost, eq: listedit.DefaultEqual, r: 3, a: []string{"a", "bc", "defg"}, b: "axcdfgh"},
	{f: listedit.StandardCost, eq: listedit.DefaultEqual, r: 3, a: []string{"a", "bc", "d", "efg"}, b: "axcdfgh"},
}

func (s *S) TestMultiEdit(c *C) {
	for _, test := range multiEditTests {
		c.Logf("Test: %v", test)
		f := test.f
		if f == nil {
			f = listedit.StandardCost
		}
		eq := test.eq
		if eq == nil {
			eq = listedit.DefaultEqual
		}
		var alist [][]any
		for _, astr := range test.a {
			alist = append(alist, splitString(astr))
		}
		blist := splitString(test.b)
		r, multiOps := listedit.MultiEdit(alist, blist, f, eq)
		c.Assert(r, Equals, test.r)
		c.Assert(merge(applyMultiOps(alist, multiOps)), DeepEquals, blist)
	}
}

func splitString(s string) []any {
	r := make([]any, len(s))
	for i, c := range s {
		r[i] = string(c)
	}
	return r
}

func applyOps(a []any, ops []listedit.Op) []any {
	// Simulate applying the operations to array a
	for _, op := range ops {
		switch op.Kind {
		case "match":
			continue
		case "replace":
			a[op.IdxA] = op.ValueB
		case "delete":
			a = append(a[:op.IdxA], a[op.IdxA+1:]...)
		case "insert":
			a = append(a, "")                // extend slice with dummy element
			copy(a[op.IdxA+1:], a[op.IdxA:]) // shift elements after insertion index to the right
			a[op.IdxA] = op.ValueB
		}
	}
	return a
}

func applyMultiOps(a [][]any, ops []listedit.MultiOp) [][]any {
	// Simulate applying the operations to multi-array a
	for _, op := range ops {
		switch op.Kind {
		case "match":
			continue
		case "replace":
			a[op.SubA][op.SubIdxA] = op.ValueB
		case "delete":
			subA := a[op.SubA]
			subA = append(subA[:op.SubIdxA], subA[op.SubIdxA+1:]...)
			a[op.SubA] = subA
		case "insert":
			subA := a[op.SubA]
			subA = append(subA, "")                      // extend slice with dummy element
			copy(subA[op.SubIdxA+1:], subA[op.SubIdxA:]) // shift elements after insertion index to the right
			subA[op.SubIdxA] = op.ValueB
			a[op.SubA] = subA
		}
	}
	return a
}

func merge(a [][]any) []any {
	merged := []any{}
	for _, subA := range a {
		merged = append(merged, subA...)
	}
	return merged
}

func BenchmarkDistance(b *testing.B) {
	one := splitString("abdefghijklmnopqrstuvwxyz")
	two := splitString("a.d.f.h.j.l.n.p.r.t.v.x.z")
	for i := 0; i < b.N; i++ {
		listedit.Edit(one, two, listedit.StandardCost, listedit.DefaultEqual)
	}
}
