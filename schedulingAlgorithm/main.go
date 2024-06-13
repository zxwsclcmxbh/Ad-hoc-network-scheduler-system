package main

import (
	"fmt"
	"math/rand"
	"sort"
)

// Node 节点结构体
type Node struct {
	C     int   // 算力大小
	M     int   // 存储大小
	Delay []int // 通信时延
}

// Task 任务结构体
type Task struct {
	C int // 算力需求
	M int // 存储需求
}

// Population 个体编码为任务部署方案
type Population [][]int

// NewPopulation 创建种群
func NewPopulation(populationSize, tasksCount int) Population {
	pop := make(Population, populationSize)
	for i := 0; i < populationSize; i++ {
		individual := make([]int, tasksCount)
		for j := range individual {
			individual[j] = rand.Intn(5) // 任务随机部署到节点上
		}
		pop[i] = individual
	}
	return pop
}

// Fitness 计算个体的适应度
func Fitness(individual []int, nodes []Node, tasks []Task) int {
	totalDelay := 0
	for i, task := range tasks {
		node := nodes[individual[i]]
		totalDelay += node.Delay[0]*task.C + node.Delay[1]*task.M
	}
	return -totalDelay // 负值因为我们是最大化问题
}

// Crossover 交叉操作
func Crossover(individual1, individual2 []int) ([]int, []int) {
	crossoverPoint := rand.Intn(len(individual1))
	newIndividual1 := make([]int, len(individual1))
	newIndividual2 := make([]int, len(individual2))
	copy(newIndividual1, individual1)
	copy(newIndividual2, individual2)
	newIndividual1 = append(newIndividual1[:crossoverPoint], individual2[crossoverPoint:]...)
	newIndividual2 = append(newIndividual2[:crossoverPoint], individual1[crossoverPoint:]...)
	return newIndividual1, newIndividual2
}

// Mutate 变异操作
func Mutate(individual []int) []int {
	mutatePoint := rand.Intn(len(individual))
	newIndividual := make([]int, len(individual))
	copy(newIndividual, individual)
	newIndividual[mutatePoint] = rand.Intn(5)
	return newIndividual
}

// GeneticAlgorithm 遗传算法
func GeneticAlgorithm(population Population, nodes []Node, tasks []Task, generations int) ([]int, int) {
	for generation := 0; generation < generations; generation++ {
		sort.SliceStable(population, func(i, j int) bool {
			return Fitness(population[i], nodes, tasks) > Fitness(population[j], nodes, tasks)
		})

		newPopulation := make(Population, len(population))
		newPopulation[0] = population[0] // 保留最优个体
		for i := 1; i < len(population); i += 2 {
			individual1 := population[rand.Intn(len(population))]
			individual2 := population[rand.Intn(len(population))]
			newIndividual1, newIndividual2 := Crossover(individual1, individual2)
			newPopulation[i] = newIndividual1
			newPopulation[i+1] = newIndividual2
		}

		for i := 1; i < len(newPopulation); i++ {
			if rand.Float64() < 0.1 { // 10% 的概率进行变异
				newPopulation[i] = Mutate(newPopulation[i])
			}
		}

		population = newPopulation
	}

	bestIndividual := population[0]
	bestFitness := Fitness(bestIndividual, nodes, tasks)
	return bestIndividual, bestFitness
}

func main() {
	rand.Seed(42)

	// 节点信息
	nodes := []Node{
		{C: 10, M: 20, Delay: []int{0, 5, 10}},  // 节点1
		{C: 15, M: 25, Delay: []int{5, 0, 8}},   // 节点2
		{C: 20, M: 30, Delay: []int{10, 8, 0}},  // 节点3
		{C: 25, M: 35, Delay: []int{15, 10, 0}}, // 节点4
		{C: 30, M: 40, Delay: []int{20, 15, 0}}, // 节点5
	}

	// 任务信息
	tasks := []Task{
		{C: 8, M: 12},  // 任务1
		{C: 10, M: 15}, // 任务2
		{C: 12, M: 18}, // 任务3
	}

	populationSize := 100
	generations := 1000

	population := NewPopulation(populationSize, len(tasks))
	bestIndividual, bestFitness := GeneticAlgorithm(population, nodes, tasks, generations)

	fmt.Println("最优解:", bestIndividual)
	fmt.Println("最优适应度:", bestFitness)
}
