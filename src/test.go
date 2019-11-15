package main

import (
	"log"
	"roboticSystem"
)

func main() {
	dh := []roboticSystem.DHParameters{
		{0.1,      0,      0,  1.571},
		{0.2,      0, 0.4318,      0},
		{0.3,   0.15, 0.0203, -1.571},
		{0, 0.4318,      0,  1.571},
		{0,      0,      0, -1.571},
		{0,      0,      0,      0},
	}
	log.Print(roboticSystem.ParametersToTransformationMatrix(dh))
}
