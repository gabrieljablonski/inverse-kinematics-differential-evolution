package roboticSystem

import (
	"arrays"
	"utils"
	"vectors"
)

type Link struct {
	DHParameters DHParameters
	ThetaSpace   utils.Range1D
}

type System struct {
	BasePosition vectors.Vector3D
	Links        []Link
}

func NewSystem(x, y, z float64) System {
	return System{
		BasePosition: vectors.NewVector3D(x, y, z),
		Links:        []Link{},
	}
}

func (s *System) Length() int {
	return len(s.Links)
}

func (s *System) AddLink(dh DHParameters, space utils.Range1D) {
	s.Links = append(s.Links, Link{
		DHParameters: dh,
		ThetaSpace:   space,
	})
}

func (s *System) SetTheta(link int, theta float64) {
	s.Links[link].DHParameters.Theta = theta
}

func (s *System) UpdateThetas(thetas *arrays.Array1D) {
	for i, theta := range thetas.Items() {
		s.SetTheta(i, theta)
	}
}

func (s *System) GetThetaValueSpace() []utils.Range1D {
	valueSpace := make([]utils.Range1D, s.Length())
	for i, link := range s.Links {
		valueSpace[i] = link.ThetaSpace
	}
	return valueSpace
}

func (s *System) DHParameters() []DHParameters {
	dh := make([]DHParameters, s.Length())
	for _, link := range s.Links {
		dh = append(dh, link.DHParameters)
	}
	return dh
}

// Calculates position of the junctions for each link
func (s *System) LinkPositions() []vectors.Vector3D {
	linkPositions := make([]vectors.Vector3D, len(s.Links)+1)
	linkPositions[0] = s.BasePosition
	transformMatrices := make([]*arrays.Array2D, s.Length())
	// first if T1
	transformMatrices[0] = s.Links[0].DHParameters.TransformationMatrix()
	for i := 1; i < len(s.Links); i++ {
		param := s.Links[i].DHParameters
		// i-th is T1*T2*...*Ti
		transformMatrices[i] = transformMatrices[i-1].Multiply(param.TransformationMatrix())
	}
	// link 1 -> T1*link0
	// link 2 -> T1*T2*link0
	// link i -> T1*T2*...*Ti*link0
	for i := 0; i < len(s.Links); i++ {
		newPosition := linkPositions[0].Transform(transformMatrices[i])
		linkPositions[i+1] = newPosition
	}
	return linkPositions
}

func (s *System) ManipulatorPosition() vectors.Vector3D {
	point := vectors.NewVector3D(0, 0, 0)
	transformationMatrix := ParametersToTransformationMatrix(s.DHParameters())
	return point.Transform(transformationMatrix)
}

//func SystemFromArray1D(array *arrays.Array1D) System {
//	if array.Length()%4 != 0 {
//		log.Fatalf("Invalid array (length not multiple of 4)")
//	}
//	dhParams := make([]DHParameters, array.Length()/4)
//	for i := 0; i < array.Length(); i += 4 {
//		dhParams[i/4] = DHParameters{
//			Theta: array.Get(i + 0),
//			D:     array.Get(i + 1),
//			R:     array.Get(i + 2),
//			Alpha: array.Get(i + 3),
//		}
//	}
//	return System{
//		BasePosition: vectors.NewVector3D(0, 0, 0),
//		Links: dhParams,
//	}
//}
