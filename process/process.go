package process

import (
	"fmt"
	"time"
)

type ProcessList struct {
	Processes       map[uint64]Process
	ContinueRunning map[uint64]bool
}

type Process struct {
	Id              uint64
	Value           uint64
	ContinueRunning bool
}

func (lp *ProcessList) AddProcess(p Process) {
	lp.Processes[p.Id] = p
	lp.ContinueRunning[p.Id] = true
	go lp.StartProcess(p.Id)
}

func (lp *ProcessList) UpdateProcess(p *Process) {
	lp.Processes[p.Id] = *p
}

func (lp *ProcessList) StartConcurrentProcesses() {
	for _, p := range lp.Processes {
		go lp.StartProcess(p.Id)
	}
}

func (lp *ProcessList) StopRunningProcess(id uint64) {
	if lp.ProcessIdExists(id) {
		lp.ContinueRunning[id] = false
		fmt.Printf("Proceso #%d fue detenido\n", id)
	} else {
		fmt.Printf("Proceso #%d no encontrado\n", id)
	}
}

func (lp *ProcessList) DeleteProcess(id uint64) {
	delete(lp.Processes, id)
	delete(lp.ContinueRunning, id)
}

func (lp *ProcessList) StartProcess(pId uint64) {
	p := lp.GetProcess(pId)
	for {
		if lp.ContinueRunning[p.Id] {
			p.Value++ //Arreglar problema con el contador cuando sale del server.
			fmt.Printf("ID Proc: %d Contador: %d\n", p.Id, p.Value)
			time.Sleep(time.Millisecond * 500)
		} else {
			lp.UpdateProcess(p)
			break
		}
	}
}

func (p *Process) StartProcess() {
	for {
		if p.ContinueRunning {
			p.Value++
			fmt.Printf("ID Proc: %d Contador: %d\n", p.Id, p.Value)
			time.Sleep(time.Millisecond * 500)
		} else {
			break
		}
	}
}

func (lp *ProcessList) GetProcess(pId uint64) *Process {
	for k, p := range lp.Processes {
		if k == pId {
			return &p
		}
	}
	return nil
}

func (lp *ProcessList) ProcessIdExists(id uint64) bool {
	_, found := lp.Processes[id]
	if found {
		return true
	}
	return false
}
