package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const size = 8

type cell interface {
	String() string
}

type ocell struct {
	count int
}

func (oc ocell) String() string {
	return strconv.Itoa(oc.count)
}

type xcell struct {
	x      int
	y      int
	ocells []*ocell
}

func (xc *xcell) increment() {
	for _, c := range xc.ocells {
		if c != nil {
			c.count++
		}
	}
}

func (xc xcell) String() string {
	return "X"
}

type field struct {
	cells  [size][size]cell
	xcells []*xcell
}

func (f *field) linkCells() {
	for _, xc := range f.xcells {
		x := xc.x
		y := xc.y

		xy := [][]int{
			{x - 1, y},
			{x + 1, y},
			{x, y - 1},
			{x, y + 1},
			{x - 1, y - 1},
			{x + 1, y - 1},
			{x + 1, y + 1},
			{x - 1, y + 1},
		}

		for i := 0; i < len(xy); i++ {
			r := xy[i]

			x, y := r[0], r[1]

			if x >= 0 && x < size && y >= 0 && y < size {
				c := f.cells[x][y]
				if oc, ok := c.(*ocell); ok {
					xc.ocells = append(xc.ocells, oc)
				}
			}
		}

	}
}

func (f *field) calc() {
	for _, xc := range f.xcells {
		xc.increment()
	}
}

func (f field) String() string {
	lines := make([]string, 0, size)

	for i := 0; i < size; i++ {
		parts := make([]string, 0, size)

		for j := 0; j < size; j++ {
			parts = append(parts, f.cells[i][j].String())
		}

		line := strings.Join(parts, " ")
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func buildField(r *bufio.Reader) (f field) {
	var (
		err  error
		line string
		i    int
	)

	for err != io.EOF {
		line, err = r.ReadString('\n')

		var j int
		for _, c := range line {
			switch c {
			case 'X':
				xc := &xcell{x: i, y: j}
				f.cells[i][j] = xc
				f.xcells = append(f.xcells, xc)
				j++

			case 'O':
				f.cells[i][j] = &ocell{}
				j++
			}

		}

		i++
	}

	return f
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Provide input file path")
	}

	input := os.Args[1]

	file, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := bufio.NewReader(file)

	f := buildField(r)
	f.linkCells()
	f.calc()

	fmt.Println(f)
}
