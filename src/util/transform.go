package util

type Rect struct {
	X float32
	Y float32
	W float32
	H float32
}

func (r1 *Rect) CenterAndFit(r2 *Rect) Rect {
	newRect := r1
	aspect := r1.W / r1.H
	if r1.W > r2.W {
		newRect.W = r2.W
		newRect.H = newRect.W / aspect
	}
	if r1.H > r2.H {
		newRect.H = r2.H
		newRect.W = newRect.H * aspect
	}
	newRect.X = (r2.W - newRect.W) / 2
	newRect.Y = (r2.H - newRect.H) / 2
	return *newRect
}

func NewRect(x float32, y float32, w float32, h float32) Rect {
	return Rect{x, y, w, h}
}
