package vortice

import (
	"testing"
	"vortice/object"
)

// dummy types for Register1..6 (arguments)
type ra struct{}
type rb struct{}
type rc struct{}
type rd struct{}
type re struct{}
type rf struct{}

// distinct return types to satisfy unique registration
type rt0 struct{ a int }
type rt1 struct{ a int }
type rt2 struct{ a int }
type rt3 struct{ a int }
type rt4 struct{ a int }
type rt5 struct{ a int }
type rt6 struct{ a int }

func f0() *rt0                                      { return &rt0{a: 0} }
func f1(a *ra) *rt1                                 { return &rt1{a: 1} }
func f2(a *ra, b *rb) *rt2                          { return &rt2{a: 2} }
func f3(a *ra, b *rb, c *rc) *rt3                   { return &rt3{a: 3} }
func f4(a *ra, b *rb, c *rc, d *rd) *rt4            { return &rt4{a: 4} }
func f5(a *ra, b *rb, c *rc, d *rd, e *re) *rt5     { return &rt5{a: 5} }
func f6(a *ra, b *rb, c *rc, d *rd, e *re, f *rf) *rt6 { return &rt6{a: 6} }

func TestRegister0To6_NoPanic(t *testing.T) {
	Register0(f0, WithDesc("d"))
	Register1[*rt1, *ra](f1, WithSingleton())
	Register2[*rt2, *ra, *rb](f2, WithPrototype())
	Register3[*rt3, *ra, *rb, *rc](f3)
	Register4[*rt4, *ra, *rb, *rc, *rd](f4)
	Register5[*rt5, *ra, *rb, *rc, *rd, *re](f5)
	Register6[*rt6, *ra, *rb, *rc, *rd, *re, *rf](f6)
}

type rtOpt struct{}

func TestRegister_WithOption_SetTag(t *testing.T) {
	called := false
	op := func(p *object.Property) { called = true; p.SetTag("k", "v") }
	Register0(func() *rtOpt { return &rtOpt{} }, op)
	if !called {
		t.Fatalf("option not called")
	}
}

// Duplicate registration should panic (unique=true in register)
func TestRegister0_DisallowDuplicate(t *testing.T) {
	type rtDup0 struct{}
	f := func() *rtDup0 { return &rtDup0{} }
	Register0(f)
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on duplicate Register0")
		}
	}()
	Register0(f)
}


