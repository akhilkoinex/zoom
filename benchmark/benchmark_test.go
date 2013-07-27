package benchmark

import (
	"github.com/stephenalexbrowne/zoom"
	"testing"
)

// not a real test. Just setup
func BenchmarkSetUp(b *testing.B) {
	setUp()
	b.SkipNow() // only need to do this once
}

func BenchmarkSave(b *testing.B) {

	// try once to make sure there's no errors
	pt := NewPerson("Alice", 25)
	err := zoom.Save(pt)
	if err != nil {
		b.Error(err)
	}

	// reset the timer
	b.ResetTimer()

	// run the actual test
	for i := 0; i < b.N; i++ {
		p := NewPerson("Bob", 25)
		zoom.Save(p)
	}

}

func BenchmarkFindById(b *testing.B) {

	// save a Person model
	p := NewPerson("Clarence", 25)
	err := zoom.Save(p)
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()

	// run the actual test
	for i := 0; i < b.N; i++ {
		zoom.FindById("person", p.Id)
	}
}

func BenchmarkDeleteById(b *testing.B) {

	// First, create and delete one person to
	// make sure there's no errors
	p := NewPerson("Dennis", 25)
	err := zoom.Save(p)
	if err != nil {
		b.Error(err)
	}
	err = zoom.DeleteById("person", p.Id)
	if err != nil {
		b.Error(err)
	}

	// save a shit ton of person models
	// we don't know how big b.N will be,
	// so better be safe
	NUM_IDS := 100000
	ids := make([]string, 0, NUM_IDS)
	for i := 0; i < NUM_IDS; i++ {
		p := NewPerson("Fred", 25)
		err := zoom.Save(p)
		if err != nil {
			b.Error(err)
		}
		ids = append(ids, p.Id)
	}
	b.ResetTimer()

	// run the actual test
	for i := 0; i < b.N; i++ {
		zoom.DeleteById("person", p.Id)
	}
}

// not an actual test
// delete any keys leftover from the previous test
func BenchmarkCleanUp(b *testing.B) {
	tearDown()
	b.SkipNow() // only need to do this once
}
