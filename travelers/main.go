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
	travelerType int // 0 - normal, 1 - wild
	travelerId   int
	x, y         int
	response     chan bool
}

type PrintRequest struct {
	travelerType int
	travelerId   int
	fromX, fromY int
	toX, toY     int
}

type KillRequest struct {
	travelerId   int
	travelerType int
}

type Cell struct {
	travelerType int
	travelerId   int
	requests     chan *Request
	killRequests chan *KillRequest
}

func main() {
	n := 10 // Number of rows
	m := 10 // Number of columns

	printRequestsChan := make(chan *PrintRequest)
	printKillChan := make(chan *KillRequest)
	printChan := make(chan bool)

	// Triggers print function
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)

			printChan <- true
		}
	}()

	// Print function logic
	go func() {
		printRequests := make([]*PrintRequest, 0)
		travelers := make([]*Traveler, 0)
		wildLocators := make([]*Traveler, 0)
		threats := make([]*Traveler, 0)

		for {
			select {
			case killRequest := <-printKillChan:
				{
					if killRequest.travelerType == 1 {
						for i := 0; i < len(wildLocators); i++ {
							if wildLocators[i].id == killRequest.travelerId {
								wildLocators = append(wildLocators[:i], wildLocators[i+1:]...)
								break
							}
						}
					} else if killRequest.travelerType == 2 {
						for i := 0; i < len(threats); i++ {
							if threats[i].id == killRequest.travelerId {
								threats = append(threats[:i], threats[i+1:]...)
								break
							}
						}
					}
				}
			case <-printChan:
				{
					for _, printRequest := range printRequests {
						if printRequest.travelerType == 0 {
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
						} else if printRequest.travelerType == 1 {
							containsTraveler := false

							for _, traveler := range wildLocators {
								if traveler.id == printRequest.travelerId {
									traveler.x = printRequest.toX
									traveler.y = printRequest.toY

									containsTraveler = true
									break
								}
							}

							if !containsTraveler {
								wildLocators = append(wildLocators, &Traveler{
									id: printRequest.travelerId,
									x:  printRequest.toX,
									y:  printRequest.toY,
								})
							}
						} else if printRequest.travelerType == 2 {
							threats = append(threats, &Traveler{
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

					for j := 0; j < 3*m+2; j++ {
						fmt.Print("#")
					}

					fmt.Println()

					for j := 0; j < n; j++ {
						fmt.Print("#")

						for i := 0; i < m; i++ {
							foundMove := false
							foundTraveler := false
							foundWildTraveler := false
							foundThreat := false
							foundTravelerId := 0

							for _, traveler := range travelers {
								if traveler.x == i && traveler.y == j {
									foundTraveler = true
									foundTravelerId = traveler.id
									break
								}
							}

							for _, traveler := range wildLocators {
								if traveler.x == i && traveler.y == j {
									foundWildTraveler = true
									break
								}
							}

							for _, traveler := range threats {
								if traveler.x == i && traveler.y == j {
									foundThreat = true
									break
								}
							}

							if foundTraveler {
								fmt.Print(foundTravelerId)

								if foundTravelerId < 10 {
									fmt.Print(" ")
								}
							} else if foundWildTraveler {
								fmt.Print("* ")
							} else if foundThreat {
								fmt.Print("# ")
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
				travelerId:   -1,
				requests:     make(chan *Request),
				killRequests: make(chan *KillRequest),
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
					case killRequest := <-grid[i][j].killRequests:
						{
							if grid[i][j].travelerId == killRequest.travelerId && grid[i][j].travelerType == killRequest.travelerType {
								fmt.Println(killRequest)
								grid[i][j].travelerId = -1

								printKillChan <- killRequest
							}
						}
					case <-moveChan:
						{
							// Move traveler to adjacent field
							if grid[i][j].travelerId != -1 && grid[i][j].travelerType == 0 {
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
								grid[i][j].travelerType = request.travelerType
								grid[i][j].travelerId = request.travelerId

								request.response <- true

								printRequestsChan <- &PrintRequest{
									travelerType: request.travelerType,
									travelerId:   request.travelerId,
									fromX:        request.x,
									fromY:        request.y,
									toX:          i,
									toY:          j,
								}
							} else if grid[i][j].travelerType == 1 && request.travelerType == 0 {
								directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

								for _, move := range directions {
									newRow, newCol := i+move[0], j+move[1]

									if newRow >= 0 && newRow < n && newCol >= 0 && newCol < m {
										responseChan := make(chan bool)

										grid[newRow][newCol].requests <- &Request{
											travelerType: 1,
											travelerId:   grid[i][j].travelerId,
											x:            i,
											y:            j,
											response:     responseChan,
										}

										response := <-responseChan

										if response {
											grid[i][j].travelerType = request.travelerType
											grid[i][j].travelerId = request.travelerId

											request.response <- true

											printRequestsChan <- &PrintRequest{
												travelerType: request.travelerType,
												travelerId:   request.travelerId,
												fromX:        request.x,
												fromY:        request.y,
												toX:          i,
												toY:          j,
											}

											break
										}
									}
								}
							} else if grid[i][j].travelerType == 2 && request.travelerType == 0 {
								fmt.Println("SHOW")

								grid[i][j].travelerId = -1

								printKillChan <- &KillRequest{
									travelerType: 2,
									travelerId:   grid[i][j].travelerId,
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
	go func() {
		maxTravelers := 25
		numTravelers := 0

		for numTravelers < maxTravelers {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

			x := rand.Intn(n)
			y := rand.Intn(m)

			responseChan := make(chan bool)

			grid[x][y].requests <- &Request{
				travelerType: 0,
				travelerId:   numTravelers,
				x:            -1,
				y:            -1,
				response:     responseChan,
			}

			response := <-responseChan

			if response {
				numTravelers++
			}
		}
	}()

	wildLocatorLifetime := 5000 * time.Millisecond

	// Create wild locator
	go func() {
		travelerId := 0

		for {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

			x := rand.Intn(n)
			y := rand.Intn(m)

			responseChan := make(chan bool)

			grid[x][y].requests <- &Request{
				travelerType: 1,
				travelerId:   travelerId,
				x:            -1,
				y:            -1,
				response:     responseChan,
			}

			response := <-responseChan

			if response {
				internalId := travelerId

				go func() {
					time.Sleep(wildLocatorLifetime)

					for i := range grid {
						for j := range grid[i] {
							grid[i][j].killRequests <- &KillRequest{
								travelerType: 1,
								travelerId:   internalId,
							}
						}
					}
				}()

				travelerId++
			}
		}
	}()

	threatLifetime := 2000 * time.Millisecond

	// Create threat
	go func() {
		travelerId := 0

		for {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

			x := rand.Intn(n)
			y := rand.Intn(m)

			responseChan := make(chan bool)

			grid[x][y].requests <- &Request{
				travelerType: 2,
				travelerId:   travelerId,
				x:            -1,
				y:            -1,
				response:     responseChan,
			}

			response := <-responseChan

			if response {
				internalId := travelerId

				go func() {
					time.Sleep(threatLifetime)

					grid[x][y].killRequests <- &KillRequest{
						travelerType: 2,
						travelerId:   internalId,
					}
				}()

				travelerId++
			}
		}
	}()

	select {}
}
