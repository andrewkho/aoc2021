package util

type Point2D struct {
	X int
	Y int
}

type Point3D struct {
	X int
	Y int
	Z int
}

type Array2D [][]int
type Array3D [][][]int

func New2DZeros(nrows int, ncols int) Array2D {
	a := make(Array2D, nrows)
	for i := 0; i<nrows; i++ {
		a[i] = make([]int, ncols)
	}

	return a
}

func New3DZeros(nx int, ny int, nz int) Array3D {
	a := make(Array3D, nx)
	for i := 0; i<nx; i++ {
		a[i] = make([][]int, ny)
		for j := 0; j<nx; j++ {
			a[i][j] = make([]int, nz)
		}
	}

	return a
}