package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Traveler struct {
	id int
	x  int
	y  int
}

type Request struct {
	travelerId int
	x, y       int
	response   chan bool
}

type PrintRequest struct {
	travelerId   int
	fromX, fromY int
	toX, toY     int
}

type Cell struct {
	travelerId int
	requests   chan *Request
}

func main() {
	n := 10 // Number of rows
	m := 10 // Number of columns
	maxTravelers := 25
	numTravelers := 0

	printRequestsChan := make(chan *PrintRequest)
	printChan := make(chan bool)

	// Triggers print function
	go func() {
		for {
			time.Sleep(1 * time.Second)

			printChan <- true
		}
	}()

	// Print function logic
	go func() {
		printRequests := make([]*PrintRequest, 0)
		travelers := make([]*Traveler, 0)

		for {
			select {
			case <-printChan:
				{
					for _, printRequest := range printRequests {
						containsTraveler := false

						for _, traveler := range travelers {
							if traveler.id == printRequest.travelerId {
								traveler.x = printRequest.toX
								traveler.y = printRequest.toY

								containsTraveler = true
								break
							}
						}

						if !containsTraveler {
							travelers = append(travelers, &Traveler{
								id: printRequest.travelerId,
								x:  printRequest.toX,
								y:  printRequest.toY,
							})
						}
					}

					for _, printRequest := range printRequests {
						if printRequest.fromX > printRequest.toX || printRequest.fromY > printRequest.toY {
							tmpX := printRequest.fromX
							tmpY := printRequest.fromY

							printRequest.fromX = printRequest.toX
							printRequest.fromY = printRequest.toY

							printRequest.toX = tmpX
							printRequest.toY = tmpY
						}
					}

					for _, traveler := range travelers {
						fmt.Println(traveler)
					}

					for _, printRequest := range printRequests {
						fmt.Println(printRequest)
					}

					for j := 0; j < 3*m+2; j++ {
						fmt.Print("#")
					}

					fmt.Println()

					for j := 0; j < n; j++ {
						fmt.Print("#")

						for i := 0; i < m; i++ {
							foundMove := false
							foundTraveler := false
							foundTravelerId := 0

							for _, traveler := range travelers {
								if traveler.x == i && traveler.y == j {
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

							for _, printRequest := range printRequests {
								if printRequest.fromX > -1 && printRequest.fromX == i && printRequest.fromY == j {
									if printRequest.toX > printRequest.fromX {
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

						for i := 0; i < m; i++ {
							foundMove := false

							for _, printRequest := range printRequests {
								if printRequest.fromX > -1 && printRequest.fromX == i && printRequest.fromY == j {
									if printRequest.toY > printRequest.fromY {
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

					for j := 0; j < 3*m+2; j++ {
						fmt.Print("#")
					}

					fmt.Println()

					printRequests = printRequests[:0]
				}
			case printRequest := <-printRequestsChan:
				{
					printRequests = append(printRequests, printRequest)
				}
			}
		}
	}()

	// Create grid
	grid := make([][]*Cell, n)
	for i := range grid {
		grid[i] = make([]*Cell, m)

		for j := range grid[i] {
			grid[i][j] = &Cell{
				travelerId: -1,
				requests:   make(chan *Request),
			}
		}
	}

	// Launch cell goroutines
	for i := range grid {
		for j := range grid[i] {
			moveChan := make(chan bool)

			// Triggers move
			go func() {
				for {
					time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

					moveChan <- true
				}
			}()

			go func(i, j int) {
				responseChan := make(chan bool)

				for {
					select {
					case response := <-responseChan:
						{
							// Succesfully move traveler
							if response {
								grid[i][j].travelerId = -1
							}
						}
					case <-moveChan:
						{
							// Move traveler to adjacent field
							if grid[i][j].travelerId != -1 {
								directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
								move := directions[rand.Intn(4)]
								newRow, newCol := i+move[0], j+move[1]

								if newRow >= 0 && newRow < n && newCol >= 0 && newCol < m {
									grid[newRow][newCol].requests <- &Request{
										travelerId: grid[i][j].travelerId,
										x:          i,
										y:          j,
										response:   responseChan,
									}
								}
							}
						}
					case request := <-grid[i][j].requests:
						{
							// Handle move request
							if grid[i][j].travelerId == -1 {
								grid[i][j].travelerId = request.travelerId
								request.response <- true

								printRequestsChan <- &PrintRequest{
									travelerId: request.travelerId,
									fromX:      request.x,
									fromY:      request.y,
									toX:        i,
									toY:        j,
								}
							} else {
								request.response <- false
							}
						}
					}
				}
			}(i, j)
		}
	}

	// Create travelers
	for numTravelers < maxTravelers {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		x := rand.Intn(n)
		y := rand.Intn(m)

		responseChan := make(chan bool)

		grid[x][y].requests <- &Request{
			travelerId: numTravelers,
			x:          -1,
			y:          -1,
			response:   responseChan,
		}

		response := <-responseChan

		if response {
			numTravelers++
		}
	}

	select {}
}
