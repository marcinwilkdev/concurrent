package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Traveler struct {
	id int
	row, col int
}

type MoveToPrint struct {
	row, col, newRow, newCol int
}

func main() {
	rand.Seed(time.Now().UnixNano())

	n := 10 // Number of rows
	m := 10 // Number of columns
	numTravelers := 25

	travelers := make([]Traveler, 0, numTravelers)
	movesToPrint := []MoveToPrint{}

	var mutex sync.Mutex

	// Function to print the grid
	printGrid := func() {
		mutex.Lock()
		defer mutex.Unlock()

		for j := 0; j < 3 * m + 2; j++ {
			fmt.Print("#")
		}

		fmt.Println()

		for i := 0; i < n; i++ {
			fmt.Print("#")

			for j := 0; j < m; j++ {
				foundMove := false
				foundTraveler := false
				foundTravelerId := 0

				for _, traveler := range travelers {
					if traveler.row == i && traveler.col == j {
						foundTraveler = true
						foundTravelerId = traveler.id
						break
					}
				}

				if foundTraveler {
					fmt.Print(foundTravelerId)

					if foundTravelerId < 10 {
						fmt.Print(" ")
					}
				} else {
					fmt.Print("  ")
				}

				for _, moveToPrint := range movesToPrint {
					if moveToPrint.row == i && moveToPrint.col == j {
						if moveToPrint.newCol > moveToPrint.col {
							fmt.Print("-")
							foundMove = true
							break
						}
					}
				}

				if !foundMove {
					fmt.Print(" ")
				}
			}

			fmt.Println("#")
			fmt.Print("#")

			for j := 0; j < m; j++ {
				foundMove := false

				for _, moveToPrint := range movesToPrint {
					if moveToPrint.row == i && moveToPrint.col == j {
						if moveToPrint.newRow > moveToPrint.row {
							fmt.Print("|  ")

							foundMove = true
							break
						}
					}
				}

				if !foundMove {
					fmt.Print("   ")
				}
			}

			fmt.Println("#")
		}

		for j := 0; j < 3 * m + 2; j++ {
			fmt.Print("#")
		}

		fmt.Println()

		movesToPrint = []MoveToPrint{}
	}

	// Start a goroutine to periodically print the grid
	go func() {
		for {
			printGrid()
			time.Sleep(1 * time.Second) // Adjust the interval as needed
		}
	}()

	// Function to simulate traveler movement
	moveTraveler := func(traveler *Traveler) {
		for {
			mutex.Lock()
			// Generate a random direction
			directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
			move := directions[rand.Intn(4)]
			newRow, newCol := traveler.row+move[0], traveler.col+move[1]

			// Check if the new position is within the grid and not occupied
			if newRow >= 0 && newRow < n && newCol >= 0 && newCol < m {
				// Check if the new position is unoccupied
				isOccupied := false

				for _, otherTraveler := range travelers {
					if otherTraveler.row == newRow && otherTraveler.col == newCol {
						isOccupied = true
						break
					}
				}

				if !isOccupied {
					if traveler.row < newRow || traveler.col < newCol {
						moveToPrint := MoveToPrint{
							row: traveler.row,
							col: traveler.col,
							newRow: newRow,
							newCol: newCol,
						}

						movesToPrint = append(movesToPrint, moveToPrint)
					} else {
						moveToPrint := MoveToPrint{
							row: newRow,
							col: newCol,
							newRow: traveler.row,
							newCol: traveler.col,
						}

						movesToPrint = append(movesToPrint, moveToPrint)
					}

					// Move the traveler to the new position
					traveler.row, traveler.col = newRow, newCol
				}
			}

			mutex.Unlock()

			// Sleep for a random duration before the next move
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}

	// Create and start multiple travelers
	travelerId := 0

	for len(travelers) < numTravelers {
		mutex.Lock()

		row := rand.Intn(n)
		col := rand.Intn(m)

		isOccupied := false

		for _, otherTraveler := range travelers {
			if otherTraveler.row == row && otherTraveler.col == col {
				isOccupied = true
				mutex.Unlock()
				break
			}
		}

		if !isOccupied {
			traveler := Traveler{id: travelerId, row: rand.Intn(n), col: rand.Intn(m)}
			travelers = append(travelers, traveler)
			go moveTraveler(&travelers[len(travelers) - 1])

			travelerId += 1

			mutex.Unlock()

			// Sleep for a random duration before the next traveler
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}

	// Run the program indefinitely
	select {}
}
