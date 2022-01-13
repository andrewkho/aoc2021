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

type Vector2D [2]int
type Matrix2D [2][2]int

func (a Matrix2D) Dot(x Vector2D) Vector2D{
	b := Vector2D{}
	for i := 0; i<2; i++ {
		for j := 0; j<2; j++ {
			b[i] += a[i][j]*x[j]
		}
	}
	return b
}

func (v Vector2D) Add(o Vector2D) Vector2D {
	z := Vector2D{}
	for i := 0; i<2; i++ {
		z[i] = v[i] + o[i]
	}
	return z
}

func (v Vector2D) Sub(o Vector2D) Vector2D {
	z := Vector2D{}
	for i := 0; i<2; i++ {
		z[i] = v[i] - o[i]
	}
	return z
}

func (a Vector2D) Dot(b Vector2D) int {
	c := 0
	for i := 0; i<2; i++ {
		c += a[i]*b[i]
	}
	return c
}

type Vector3D [3]int
type Matrix3D [3][3]int

func (a Matrix3D) Dot(x Vector3D) Vector3D{
	b := Vector3D{}
	for i := 0; i<3; i++ {
		for j := 0; j<3; j++ {
			b[i] += a[i][j]*x[j]
		}
	}
	
	return b
}

func (v Vector3D) Add(o Vector3D) Vector3D {
	z := Vector3D{}
	for i := 0; i<3; i++ {
		z[i] = v[i] + o[i]
	}
	return z
}

func (v Vector3D) Sub(o Vector3D) Vector3D {
	z := Vector3D{}
	for i := 0; i<3; i++ {
		z[i] = v[i] - o[i]
	}
	return z
}

func (a Vector3D) Dot(b Vector3D) int {
	c := 0
	for i := 0; i<3; i++ {
		c += a[i]*b[i]
	}
	return c
}

func (v Vector3D) Cross(o Vector3D) Vector3D {
	z := Vector3D{}
	for i := 0; i<3; i++ {
		j, k := (i+1)%3, (i+2)%3
		z[i] = v[j]*o[k] - v[k]*o[j]
	}
	return z
}

type Vector4D [4]int
type Matrix4D [4][4]int

func (a Matrix4D) Dot(x Vector4D) Vector4D {
	b := Vector4D{}
	for i := 0; i<4; i++ {
		for j := 0; j<4; j++ {
			b[i] += a[i][j]*x[j]
		}
	}
	
	return b
}

func (v Vector4D) Add(o Vector4D) Vector4D {
	z := Vector4D{}
	for i := 0; i<4; i++ {
		z[i] = v[i] + o[i]
	}
	return z
}

func (v Vector4D) Sub(o Vector4D) Vector4D {
	z := Vector4D{
		v[0]-o[0],
		v[1]-o[1],
		v[2]-o[2],
		v[3]-o[3],
	}
	return z
}

func (a Vector4D) Dot(b Vector4D) int {
	c := 0
	for i := 0; i<4; i++ {
		c += a[i]*b[i]
	}
	return c
}