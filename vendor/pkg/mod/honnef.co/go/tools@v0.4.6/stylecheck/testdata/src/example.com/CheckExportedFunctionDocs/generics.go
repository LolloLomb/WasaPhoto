//go:build go1.18

package pkg

// Whatever //@ diag(`comment on exported function`)
func TPFoo[T any]() {}

// Whatever //@ diag(`comment on exported function`)
func TPBar[T1, T2 any]() {}

// TPBaz is amazing
func TPBaz[T any]() {}

type TPT[T any] struct{}

// Foo is amazing
func (TPT[T]) Foo() {}

// Whatever //@ diag(`comment on exported method`)
func (TPT[T]) Bar() {}

type TPT2[T1, T2 any] struct{}

// Foo is amazing
func (TPT2[T1, T2]) Foo() {}

// Whatever //@ diag(`comment on exported method`)
func (*TPT2[T1, T2]) Bar() {}
