package differentialEvolution

import (
	"arrays"
	"fmt"
	"log"
	"math"
	"math/rand"
	"utils"
)

// evolution optimizes fitness to 0
type FitnessFunction func(agent *arrays.Array1D) float64

type Evolver struct {
	// init factors
	AgentSize       int
	PopulationSize  int
	CrossoverRate   float64
	WeightingFactor float64
	SearchSpace     []utils.Range1D
	// termination criteria
	MaxGenerations int
	TargetFitness  float64
	StallPeriod    int
	StallFactor    float64
	stallCount     int

	CurrentGeneration  int
	CurrentBestFitness float64
	CurrentBestAgent   *arrays.Array1D
	Population         *arrays.Array2D
	FitnessFunction    FitnessFunction
}

type NewEvolverParams struct {
	AgentSize       int
	PopulationSize  int
	CrossoverRate   float64
	WeightingFactor float64
	SearchSpace     []utils.Range1D
	MaxGenerations  int
	TargetFitness   float64
	StallPeriod     int
	StallFactor     float64
	FitnessFunction FitnessFunction
}

func NewEvolver(p NewEvolverParams) Evolver {
	if p.AgentSize != len(p.SearchSpace) {
		if len(p.SearchSpace) != 1 {
			log.Fatalf("Invalid search space")
		}
		for i := 0; i < p.AgentSize-1; i++ {
			// repeat the search space given for all agent features
			p.SearchSpace = append(p.SearchSpace, p.SearchSpace[0])
		}
	}
	return Evolver{
		AgentSize:          p.AgentSize,
		PopulationSize:     p.PopulationSize,
		CrossoverRate:      p.CrossoverRate,
		WeightingFactor:    p.WeightingFactor,
		SearchSpace:        p.SearchSpace,
		MaxGenerations:     p.MaxGenerations,
		TargetFitness:      p.TargetFitness,
		StallPeriod:        p.StallPeriod,
		StallFactor:        p.StallFactor,
		CurrentGeneration:  0,
		CurrentBestFitness: math.Inf(1),
		CurrentBestAgent:   nil,
		Population:         nil,
		FitnessFunction:    p.FitnessFunction,
	}
}

func (e *Evolver) InitializePopulation() {
	e.Population = &arrays.Array2D{}
	for i := 0; i < e.PopulationSize; i++ {
		agent := &arrays.Array1D{}
		for j := 0; j < e.AgentSize; j++ {
			agent.Append(utils.RandomInRange(e.SearchSpace[j]))
		}
		e.Population.Append(agent)
	}
}

// Default termination criterion is `TargetFitness` equals to 0
// If `MaxGenerations` is set to a number greater than 0, it is checked first
func (e *Evolver) ShouldContinue() bool {
	if e.MaxGenerations > 0 && e.CurrentGeneration == e.MaxGenerations {
		log.Print("Max generations reached.")
		return false
	}
	if e.TargetFitness >= 0 && e.CurrentBestFitness <= e.TargetFitness {
		log.Print("Target fitness reached.")
		return false
	}
	// if fitness varies in a proportion less than `StallFactor` to the last fitness
	// for `StallPeriod` times, evolution halts
	if e.StallPeriod > 0 && e.stallCount >= e.StallPeriod {
		log.Print("Evolution stalled.")
		return false
	}
	return true
}

func (e *Evolver) pickRandomAgents(exclude int) (arrays.Array1D, arrays.Array1D, arrays.Array1D) {
	howMany := 3
	agentNumbers := utils.PickRandom(e.PopulationSize, howMany, exclude)
	agents := arrays.Array2D{}
	for _, n := range *agentNumbers {
		agent := e.Population.GetRow(n).Copy()
		agents.Append(agent)
	}
	return *agents[0], *agents[1], *agents[2]
}

func (e *Evolver) mutate(referenceAgentNumber int) *arrays.Array1D {
	r1, r2, r3 := e.pickRandomAgents(referenceAgentNumber)
	// Compute v = r1 + F*(r3 - r2)
	// in which
	//     `v` is the mutated value
	//     `ri` is the i-th random agent from the population
	//     `F` is the weighting factor
	mutated := r1.Add((r3.Subtract(&r2)).MultiplyByConstant(e.WeightingFactor))
	// ensure the mutated values are within the search space
	for i, value := range mutated.Items() {
		mutated.Set(i, utils.ConstrainValue(value, e.SearchSpace[i]))
	}
	return mutated
}

func (e *Evolver) crossover(referenceAgent, mutatedAgent *arrays.Array1D) *arrays.Array1D {
	crossed := referenceAgent.Copy()
	randomIndex := rand.Intn(e.AgentSize) // random index so at least one feature gets crossed
	for i := range referenceAgent.Items() {
		ri := rand.Float64()
		if ri <= e.CrossoverRate || i == randomIndex {
			crossed.Set(i, mutatedAgent.Get(i))
		} // else keep reference agent value
	}
	return crossed
}

func (e *Evolver) mutateAndCrossover(referenceAgentNumber int) *arrays.Array1D {
	referenceAgent := e.Population.GetRow(referenceAgentNumber).Copy()
	mutated := e.mutate(referenceAgentNumber)
	crossed := e.crossover(referenceAgent, mutated)
	return crossed
}

func (e *Evolver) tryReplaceAgent(referenceAgentNumber int) (newAgent *arrays.Array1D, fitness float64) {
	referenceAgent := e.Population.GetRow(referenceAgentNumber).Copy()
	crossed := e.mutateAndCrossover(referenceAgentNumber)

	referenceFitness := e.FitnessFunction(referenceAgent)
	crossedFitness := e.FitnessFunction(crossed)
	if crossedFitness <= referenceFitness {
		return crossed, crossedFitness
	}
	return referenceAgent, referenceFitness
}

func (e *Evolver) Evolve() error {
	if e.Population == nil {
		return fmt.Errorf("population not initialized")
	}
	newPopulation := make(arrays.Array2D, e.PopulationSize)
	lastBestFitness := e.CurrentBestFitness
	for i := range e.Population.Items() {
		newAgent, fitness := e.tryReplaceAgent(i)
		newPopulation.SetRow(i, *newAgent)
		if fitness <= e.CurrentBestFitness {
			e.CurrentBestFitness = fitness
			e.CurrentBestAgent = newAgent
		}
	}
	e.Population = &newPopulation

	fitnessImprovementRatio := (lastBestFitness-e.CurrentBestFitness)/lastBestFitness
	if fitnessImprovementRatio <= e.StallFactor {
		e.stallCount++
	} else {
		e.stallCount = 0
	}
	e.CurrentGeneration++
	return nil
}
