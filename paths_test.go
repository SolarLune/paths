package paths

import "testing"

var p *Path

func BenchmarkGetPathFromCells(b *testing.B) {
	for i := 0; i < b.N; i++ {
		firstMap := NewGrid(200, 200, 16, 16)
		p = firstMap.GetPathFromCells(firstMap.Get(0, 0), firstMap.Get(199, 199), false, false)
	}
}
