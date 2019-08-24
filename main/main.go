package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net"
	"strings"
)

type PlayerPacket struct {
	id int32
	pack interface{}
}

// Wait for packet from reader and then read
func waitForPacket(reader *bufio.Reader) (pack interface{}, retErr error) {
	defer func() {
		if r := recover(); r != nil {
			pack = nil
			retErr = errors.New(fmt.Sprint(r))
		}
	}()

	retErr = nil

	line, _, err := reader.ReadLine()
	if err != nil {
		panic(fmt.Sprintf("Read Error: %s", err.Error()))
	}

	packMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(line), &packMap); err != nil {
		panic(fmt.Sprintf("JSON Parse Error: %s", err.Error()))
	}

	pack, err = ParsePacket(packMap)
	if err != nil {
		panic(fmt.Sprintf("PlayerPacket Parse Error: %s", err.Error()))
	}

	return
}

// Handle new connection
// Marcus-chans everywhere
func handlePlayer(conn net.Conn, outCh chan PlayerPacket, id int32, distInCh chan float64) {
	reader := bufio.NewReader(conn)

	fmt.Printf("New connection: %s | UserID: %d\n", conn.LocalAddr().String(), id)

	for {
		pack, err := waitForPacket(reader)
		if err != nil {
			if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host.") {
				fmt.Printf("Player (ID: %d) forcibly quit\n", id)
				outCh <- PlayerPacket{id, Quit{}}
				return
			} else {
				fmt.Printf("[ERROR (ID: %d)]: %s\n", id, err.Error())
			}
		} else {
			outCh <- PlayerPacket{id, pack}

			switch pack.(type) {
			// Wait for system to send back distance info
			case Location:
				distObj := struct {
					Dist float64
				} { <- distInCh }

				outStr, err := json.Marshal(distObj)
				if err != nil {
					fmt.Printf("[ERROR (ID: %d)]: JSON OUT ERROR: %s\n", id, err.Error())
				} else {
					outStr = append(outStr, '\n')

					if _, err := conn.Write(outStr); err != nil {
						fmt.Printf("[ERROR (ID: %d)]: SOCK WRITE ERROR: %s\n", id, err.Error())
					} else {
						fmt.Printf("(ID: %d): DIST OUT: Wrote %s\n", id, outStr)
					}
				}
			case Quit:
				return
			}
		}
	}
}

// Variables to keep track of players
var playerChs = make(map[int32]chan float64) // Channels to send dist data (by ID)
var locs = make(map[int32]map[int32]Location) // List of locations by team
var team = make(map[int32]int32) // Map of teams (by ID)

func handlePackets(ch chan PlayerPacket) {
	for playerPack := range ch {
		fmt.Printf("Received packet from ID: %d, PACK: %+v\n", playerPack.id, playerPack.pack)

		switch pack := playerPack.pack.(type) {
		case Location:
			if _, ok := team[playerPack.id]; !ok {
				fmt.Printf("[ERROR]: Invalid ID %d (for set team packet)\n", playerPack.id)
			} else {
				// Get min dist
				dist := math.MaxFloat64
				playerTeam := team[playerPack.id]

				for k, mp := range locs {
					if k != playerTeam {
						for _, loc := range mp {
							dist = math.Min(dist, pack.Dist(&loc))
						}
					}
				}

				// Set loc
				if locs[playerTeam] == nil { // Make map if it is nil
					locs[playerTeam] = make(map[int32]Location)
				}
				locs[playerTeam][playerPack.id] = pack

				// Send back dist
				playerChs[playerPack.id] <- dist
			}
		case SetTeam:
			if _, ok := team[playerPack.id]; !ok {
				fmt.Printf("[ERROR]: Invalid ID %d (for set team packet)\n", playerPack.id)
			} else {
				fmt.Printf("Player (ID: %d) has set team to %d\n", playerPack.id, pack.team)

				team[playerPack.id] = pack.team
			}
		case Quit:
			fmt.Printf("Player (ID: %d) has disconnected\n", playerPack.id)

			delete(locs[team[playerPack.id]], playerPack.id)
			delete(team, playerPack.id)
			delete(playerChs, playerPack.id)
		default:
			fmt.Printf("[ERROR]: Unrecognized packet type %T\n", playerPack)
		}
	}
}

// Address stuff for connecting
const ADDR = "0.0.0.0:80"

func main() {
	server, err := net.Listen("tcp", ADDR)
	packetCh := make(chan PlayerPacket)

	// Init event handler
	go handlePackets(packetCh)

	if err != nil {
		fmt.Printf("[ERROR]: Listen Error: %s\n", err.Error())
	} else {
		fmt.Printf("Started listening on %s\n", ADDR)

		for {
			conn, err := server.Accept()

			if err != nil {
				fmt.Printf("[ERROR]: ConnError: Conn: %s, Error: %s\n", conn.LocalAddr().String(), err.Error())
			} else {
				id := MakeId()
				playerChs[id] = make(chan float64)
				team[id] = 0

				go handlePlayer(conn, packetCh, id, playerChs[id])
			}
		}
	}
}
