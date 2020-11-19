package main

import (
	"encoding/gob"
	"fmt"
	"math"
	"net"
	"time"

	"./process"
	"./useful"
)

type Server struct {
	Clients map[int64]bool
	Lp      process.ProcessList
}

func (s *Server) servidor() {
	server, err := net.Listen("tcp", useful.PORT)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for {
		client, err := server.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		go s.handleClient(client)
	}
}

func (s *Server) handleClient(client net.Conn) {
	var c useful.Client
	err := gob.NewDecoder(client).Decode(&c)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Client #", c.ClientId, "Received")
		if s.clientIdExists(c.ClientId) {
			fmt.Println("Removing client #", c.ClientId)
			s.Lp.AddProcess(c.Process)
			s.removeClientId(c.ClientId)
			fmt.Println("Client #", c.ClientId, "was removed succesfully from known clients")
		} else {
			p := s.GetTopProcess()
			if p != nil {
				s.Lp.DeleteProcess(p.Id)
				c = useful.Client{ClientId: int64(len(s.Clients)), Process: *p} //Creating client first to take into account the actual length of the ClientsList
				s.addNewClient()
				fmt.Println("Added client #", c.ClientId)
				err = gob.NewEncoder(client).Encode(c)
				if err != nil {
					fmt.Println(err)
					return
				}
			} else {
				fmt.Println("No more Processes available on server")
			}

		}
	}

}

func (s *Server) GetTopProcess() *process.Process {
	/*rand.Seed(time.Now().Unix())
	keys := make([]uint64, 0, len(s.Lp.Processes))
	for k, _ := range s.Lp.Processes {
		keys = append(keys, k)
	}
	temp := rand.Intn(len(keys))
	rIndx := uint64(temp) Initially I did this with a random seed, then I watched the video of what was required*/
	keys := make([]uint64, 0, len(s.Lp.Processes))
	for k, _ := range s.Lp.Processes {
		keys = append(keys, k)
	}
	lowestKey := getMinVal(keys)
	if lowestKey != math.MaxUint64 {
		s.Lp.StopRunningProcess(lowestKey)
		time.Sleep(time.Millisecond * 500) //Wait needed to allow update on the Process
		return s.Lp.GetProcess(lowestKey)
	}
	return nil
}

func getMinVal(s []uint64) uint64 {
	if len(s) > 0 {
		min := s[0] //base value
		for _, v := range s {
			if v < min {
				min = v
			}
		}
		return min
	}
	return math.MaxUint64
}

func (s *Server) clientIdExists(clientId int64) bool {
	if clientId != -1 {
		_, found := s.Clients[clientId]
		if found {
			return true
		}
		return false
	}
	return false
}

func (s *Server) removeClientId(clientId int64) {
	delete(s.Clients, clientId)
}

func (s *Server) addNewClient() {
	s.Clients[int64(len(s.Clients))] = false
}

func createProcessList(size uint64) *process.ProcessList {
	lp := process.ProcessList{Processes: map[uint64]process.Process{}, ContinueRunning: map[uint64]bool{}}
	initialCounterVal := uint64(0)
	for id := uint64(0); id < size; id++ {
		lp.Processes[id] = process.Process{Id: id, Value: initialCounterVal, ContinueRunning: false} //The process continue running value is for when it goes to the client
		lp.ContinueRunning[id] = true
	}
	return &lp
}

func main() {
	size := uint64(5)
	lp := createProcessList(size)
	s := Server{Clients: map[int64]bool{}, Lp: *lp}
	s.Lp.StartConcurrentProcesses()
	go s.servidor()

	fmt.Print("Press Enter to finalize...")
	fmt.Scanln()
}
