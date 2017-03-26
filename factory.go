package main

import (
	"sem"
	"time"
	"fmt"
)

const (
	build_t = 200 * time.Millisecond
	nrobots = 3
	ncables = 5
)

type ProdCons struct {
	mat_id	int
	prodsem *sem.Sem
	conssem *sem.Sem
}

type Robot struct {
	id		 int
	cables	 [ncables]*int
	pantalla *int
	carcasa	 *int
	placa	 *int
}

var (
	cableprod	 = ProdCons{0, sem.NewSem(0), sem.NewSem(1)}
	pantallaprod = ProdCons{0, sem.NewSem(0), sem.NewSem(1)}
	carcasaprod	 = ProdCons{0, sem.NewSem(0), sem.NewSem(1)}
	placaprod	 = ProdCons{0, sem.NewSem(0), sem.NewSem(1)}

)

func produce(prod *ProdCons) {
	counter := 0
	for {
		prod.conssem.Down()
		prod.mat_id = counter
		counter++
		prod.prodsem.Up()
	}
}

func getMat(cons *ProdCons) int {
	cons.prodsem.Down()
	mat := cons.mat_id
	cons.conssem.Up()
	return mat
}

func launchRobot(robot *Robot) {
	var cables [ncables]int
	for {
		for !robot.check_mats() {
			for i := range robot.cables {
				if robot.cables[i] == nil {
					cables[i] = getMat(&cableprod)
					robot.cables[i] = &cables[i] 
				}
			}

			if robot.pantalla == nil {
				pantalla := getMat(&pantallaprod)
				robot.pantalla = &pantalla
			}

			if robot.carcasa == nil {
				carcasa := getMat(&carcasaprod)
				robot.carcasa = &carcasa
			}

			if robot.placa == nil {
				placa := getMat(&placaprod)
				robot.placa = &placa
			}
		}
		robot.build()
		robot.reset()
	}
}

func (robot *Robot) toString(status string) {
	fmt.Printf("robot %d, cables %d %d %d %d %d, pantalla %d, carcasa %d, placa %d %s\n",
		robot.id, *robot.cables[0], *robot.cables[1], *robot.cables[2], *robot.cables[3],
		*robot.cables[4], *robot.pantalla, *robot.carcasa, *robot.placa, status)
}

func (robot *Robot) build() {
	robot.toString("Comenzando")
	time.Sleep(build_t)
	robot.toString("Terminado")
}

func (robot *Robot) check_mats() bool {
	ready := true

	for cable := range robot.cables {
		if robot.cables[cable] == nil {
			ready = false
		}
	}
	if robot.pantalla == nil {
		ready = false
	}
	if robot.carcasa == nil {
		ready = false
	}
	if robot.placa == nil {
		ready = false
	}

	return ready
}

func (robot *Robot) reset() {
	for cable := range robot.cables {
		robot.cables[cable] = nil
	}
	robot.pantalla = nil
	robot.carcasa = nil
	robot.placa = nil
}

func main() {
	go produce(&cableprod)
	go produce(&pantallaprod)
	go produce(&carcasaprod)
	go produce(&placaprod)

	for i := 0; i < nrobots-1; i++ {
		robot := new(Robot)
		robot.id = i
		go launchRobot(robot)
	}

	robot := new(Robot)
	robot.id = nrobots - 1
	launchRobot(robot)
}