package main

import (
	"arrays"
	de "differentialEvolution"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	rs "roboticSystem"
	"runtime"
	"strings"
	"time"
	"utils"
	"vectors"
)

const (
	PopulationSize  = 15
	CrossoverRate   = 0.5
	WeightingFactor = 0.5
	MaxGenerations  = 2000
	TargetFitness   = 0.000
	StallPeriod     = 50  // in generations
	StallFactor     = 0.1 // 0~1
)

// https://en.wikipedia.org/wiki/Ackley_function
// example minimization function for agent of size 2
// optimal point is {20, 20}
//func Ackley(xs *arrays.Array1D) float64 {
//	x, y := xs.Get(0)-20, xs.Get(1)-20
//	return -20*math.Exp(-.2*math.Sqrt(.5*(x*x+y*y))) -
//		math.Exp(.5*(math.Cos(2*math.Pi*x)+math.Cos(2*math.Pi*y))) + math.E + 20
//}

func buildFitnessFunction(target vectors.Vector3D, baseSystem rs.System) de.FitnessFunction {
	return func(agent *arrays.Array1D) float64 {
		// agent is the theta parameters for each link, one after another
		// example for n links: `agent = [θ0 θ1 ... θn]`
		baseSystem.UpdateThetas(agent)
		return baseSystem.ManipulatorPosition().Distance(target)
	}
}

//func buildSearchSpace(baseSearchSpace []utils.Range1D, repetitions int) []utils.Range1D {
//	var searchSpace []utils.Range1D
//	for i := 0; i < repetitions; i++ {
//		searchSpace = append(searchSpace, baseSearchSpace...)
//	}
//	return searchSpace
//}

func main() {
	rand.Seed(time.Now().Unix())
	baseSystem := rs.NewSystem(0, 0, 0)
	baseSystem.AddLink(
		rs.DHParameters{
			D:     0.03,
			R:     0,
			Alpha: math.Pi / 2.0,
		},
		utils.Range1D{UpperBound: math.Pi})
	baseSystem.AddLink(
		rs.DHParameters{
			D:     0,
			R:     0.1,
			Alpha: 0,
		},
		utils.Range1D{UpperBound: math.Pi})
	baseSystem.AddLink(
		rs.DHParameters{
			D:     0,
			R:     0.1,
			Alpha: 0,
		},
		utils.Range1D{LowerBound: -math.Pi})
	baseSystem.AddLink(
		rs.DHParameters{
			D:     0,
			R:     0.18,
			Alpha: 0,
		},
		utils.Range1D{
			LowerBound: -math.Pi / 2.0,
			UpperBound: math.Pi / 2.0,
		})
	// Target should have a distance smaller than 0.5 from the base of the system
	// Maximum values:
	//     |x|: x0+0.38
	//     |y|: y0+0.38
	//     |z|: z0+0.41
	// The sum of the coordinates should probably not exceed 0.5
	//target := vectors.NewVector3D(.1, .1, .1)
	target := vectors.RandomVector3D(utils.Range1D{
		LowerBound: -.2,
		UpperBound:  .2,
	})

	evolver := de.NewEvolver(de.NewEvolverParams{
		AgentSize:       baseSystem.Length(),
		PopulationSize:  PopulationSize,
		CrossoverRate:   CrossoverRate,
		WeightingFactor: WeightingFactor,
		SearchSpace:     baseSystem.GetThetaValueSpace(),
		MaxGenerations:  MaxGenerations,
		TargetFitness:   TargetFitness,
		StallPeriod:     StallPeriod,
		StallFactor:     StallFactor,
		FitnessFunction: buildFitnessFunction(target, baseSystem),
	})
	evolver.InitializePopulation()
	var bestAgentLinkPositions [][]vectors.Vector3D
	for evolver.ShouldContinue() {
		err := evolver.Evolve()
		if err != nil {
			log.Fatal(err)
		}
		baseSystem.UpdateThetas(evolver.CurrentBestAgent)
		bestAgentLinkPositions = append(bestAgentLinkPositions, baseSystem.LinkPositions())
		log.Printf("---Generation %d---", evolver.CurrentGeneration)
		log.Printf("Best agent: %s", evolver.CurrentBestAgent)
		log.Printf("Position: %s", baseSystem.ManipulatorPosition())
		log.Printf("Fitness: %.3f", evolver.CurrentBestFitness)
	}
	output := make([]string, len(bestAgentLinkPositions)+1)
	output[0] = target.String()
	for i, generation := range bestAgentLinkPositions {
		line := make([]string, len(generation))
		for j, linkPosition := range generation {
			line[j] = linkPosition.String()
		}
		output[i+1] = strings.Join(line, "\t")
	}

	var filename string
	if len(os.Args) == 1 {
		filename = "example_output.txt"
	} else {
		filename = os.Args[1]
	}
	err := ioutil.WriteFile(filename, []byte(strings.Join(output, "\n")), 0644)
	if err != nil {
		log.Print(err)
	}
	var python string
	if runtime.GOOS == "windows" {
		python = "python"
	} else {
		python = "python3"
	}
	scriptName := "plot_link_generations.py"
	cmd := exec.Command(python, scriptName, filename)
	err = cmd.Run()

	if err != nil {
		log.Fatalf("%#v", err)
	}
	//delta := math.Pi/20.0
	//var positions []vectors.Vector3D
	//for t1 := searchSpace[0].LowerBound; t1 <= searchSpace[0].UpperBound; t1 += delta {
	//	for t2 := searchSpace[1].LowerBound; t2 <= searchSpace[1].UpperBound; t2 += delta {
	//		for t3 := searchSpace[2].LowerBound; t3 <= searchSpace[2].UpperBound; t3 += delta {
	//			for t4 := searchSpace[3].LowerBound; t4 <= searchSpace[3].UpperBound; t4 += delta {
	//				baseSystem.UpdateThetas(&arrays.Array1D{t1, t2, t3, t4})
	//				positions = append(positions, baseSystem.ManipulatorPosition())
	//			}
	//		}
	//	}
	//}
	//sSearchSpace := make([]string, len(positions))
	//for i, p := range positions {
	//	sSearchSpace[i] = p.String()
	//}
	//err := ioutil.WriteFile("search_space.txt", []byte(strings.Join(sSearchSpace, "\n")), 0644)
	//
	//if err != nil {
	//	log.Print(err)
	//}
}
