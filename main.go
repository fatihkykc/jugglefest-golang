package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type vector3 struct {
	h int
	e int
	p int
}

type juggler struct {
	name     string
	vec      vector3
	prefList []string
	counter  int
	dots     map[string]int
}

type jugglers []*juggler

func (j *juggler) getNextPref() string {
	dest := j.prefList[j.counter]
	j.counter++
	return dest
}

func addJuggler(j *juggler) {
	dest := j.getNextPref()
	fmt.Println(circuitlist[dest])
	j.dots = make(map[string]int)
	if j.counter >= 6 {
		leftjugglerlist = append(leftjugglerlist, j)
		return
	}
	if circuitlist[dest] != nil {
		j.dots[dest] = dotProduct(j, circuitlist[dest])
		r, e := circuitlist[dest].push(j)
		if e == true && r == nil {
			//fmt.Println("got in")
			return
		}
		if reflect.TypeOf(r).String() == reflect.TypeOf(j).String() && e == true {
			//fmt.Println("got in but kicked someone")
			addJuggler(r)
		}
		if e == false && r == nil {
			//fmt.Println("doesnt have enough points")
			addJuggler(j)
		}
	}

}

func addLeftJuggler(j *juggler) {
	var circuits circuits
	for dest := range circuits {
		if len(circuits[dest].assignedJugglers) < circuits[dest].maxJugglers {
			j.dots[dest] = dotProduct(j, circuits[dest])
			r, e := circuits[dest].push(j)
			if e == true && r == nil {
				return
			}
			if reflect.TypeOf(r).String() == reflect.TypeOf(j).String() && e == true {
				addJuggler(r)
			}
			if e == false && r == nil {
				addJuggler(j)
			}
		}
	}
}

type circuit struct {
	name             string
	vec              vector3
	assignedJugglers []*juggler
	maxJugglers      int
}

func (c *circuit) push(j *juggler) (*juggler, bool) {
	if len(c.assignedJugglers) >= c.maxJugglers {
		sort.Slice(c.assignedJugglers, func(i, j int) bool {
			return c.assignedJugglers[i].dots[c.name] < c.assignedJugglers[j].dots[c.name]
		})
		minJuggler := c.assignedJugglers[len(c.assignedJugglers)-1]
		if j.dots[c.name] > minJuggler.dots[c.name] {
			loser := minJuggler
			c.assignedJugglers = c.assignedJugglers[:len(c.assignedJugglers)-1+copy(c.assignedJugglers[len(c.assignedJugglers)-1:],
				c.assignedJugglers[len(c.assignedJugglers)-1+1:])]
			c.assignedJugglers = append(c.assignedJugglers, j)
			return loser, true
		}
		if j.dots[c.name] <= minJuggler.dots[c.name] {
			return nil, false
		}
	} else {
		c.assignedJugglers = append(c.assignedJugglers, j)
		return nil, true
	}
	panic("should never happen")
}

type circuits map[string]*circuit

func dotProduct(j *juggler, c *circuit) int {
	return int(j.vec.h)*int(c.vec.h) + int(j.vec.e)*int(c.vec.e) + int(j.vec.p)*int(c.vec.p)
}

type leftJugglers []*juggler

var circuitlist circuits
var jugglerlist jugglers
var leftjugglerlist leftJugglers

func main() {
	///Make a method that can fill the data for the circuit object.
	file, err := os.Open("/home/fatihkykc/go/src/goJuggleFest/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	circuitlist = make(map[string]*circuit)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.HasPrefix(s, "C") {
			//var obj circuit
			name := strings.Split(s, " ")[1]
			h := strings.Split(s, " ")[2]
			h = strings.Split(h, ":")[1]
			e := strings.Split(s, " ")[3]
			e = strings.Split(e, ":")[1]
			p := strings.Split(s, " ")[4]
			p = strings.Split(p, ":")[1]
			H, err := strconv.Atoi(h)
			E, err := strconv.Atoi(e)
			P, err := strconv.Atoi(p)
			obj := new(circuit)
			obj.name = name
			obj.vec = vector3{H, E, P}
			obj.maxJugglers = 6
			obj.assignedJugglers = nil
			//circuits = map[string]*circuit{}
			circuitlist[name] = obj
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println(circuits[name])
		}

		if strings.HasPrefix(s, "J") {
			name := strings.Split(s, " ")[1]
			h := strings.Split(s, " ")[2]
			h = strings.Split(h, ":")[1]
			e := strings.Split(s, " ")[3]
			e = strings.Split(e, ":")[1]
			p := strings.Split(s, " ")[4]
			p = strings.Split(p, ":")[1]
			H, err := strconv.Atoi(h)
			E, err := strconv.Atoi(e)
			P, err := strconv.Atoi(p)
			prefList := strings.Split(s, " ")[5]
			obj := new(juggler)
			obj.name = name
			obj.vec = vector3{H, E, P}
			var prefs []string
			preflim := 10
			for i := 0; i < preflim; i++ {
				prefs = append(prefs, strings.Split(prefList, ",")[i])
			}
			obj.prefList = prefs
			jugglerlist = append(jugglerlist, obj)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for _, juggler := range jugglerlist {
		addJuggler(juggler)
		//fmt.Println(juggler.prefList[0])
	}

	// for _, juggler := range leftJugglers {
	// 	addJuggler(juggler)
	// }
	for i := range jugglerlist {
		fmt.Println(jugglerlist[i])
	}
	//fmt.Println(circuitlist["C1998"])
}
