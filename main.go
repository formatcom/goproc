package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"log"
	"strconv"
)

type Process struct {
	Pid int
	Children []int
	Cmd string
	Path string
}

// https://github.com/c9s/goprocinfo/blob/master/linux/process_cmdline.go
func ReadProcessCmdline(path string) (string, error) {

	b, err := ioutil.ReadFile(path)

	if err != nil {
		return "", err
	}

	l := len(b) - 1 // Define limit before last byte ('\0')
	z := byte(0)    // '\0' or null byte
	s := byte(0x20) // space byte
	c := 0          // cursor of useful bytes

	for i := 0; i < l; i++ {

		// Check if next byte is not a '\0' byte.
		if b[i+1] != z {

			// Offset must match a '\0' byte.
			c = i + 2

			// If current byte is '\0', replace it with a space byte.
			if b[i] == z {
				b[i] = s
			}
		}
	}

	x := strings.TrimSpace(string(b[0:c]))

	return x, nil
}

func main() {

	root := "/proc"

	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	process := make([]Process, 0, 5)

	// recorremos todos los archivos en proc
	for _, file := range files {
		name := file.Name()


		// se verifica que sea un directorio y que sea un numero
		pid, err := strconv.Atoi(name)
		if err != nil || !file.IsDir() {
			break
		}

		p := Process{Pid: pid}
		p.Path = fmt.Sprintf("%s/%d", root, pid)

		cmdline, err := ReadProcessCmdline(fmt.Sprintf("%s/cmdline", p.Path))
		if err != nil {
			log.Fatal(err)
		}
		cmd := strings.Split(cmdline, " ")[0]

		p.Cmd = cmd

		task := fmt.Sprintf("%s/task", p.Path)
		children, err := ioutil.ReadDir(task)
		if err != nil {
			log.Fatal(err)
		}

		for _, child := range children {
			name := child.Name()
			thread, err := strconv.Atoi(name)
			if err != nil || !child.IsDir() {
				break
			}
			f, err := ioutil.ReadFile(fmt.Sprintf(
				"%s/%d/children", task, thread))

			if err != nil {
				log.Fatal(err)
			}

			children := strings.Split(string(f), " ")
			for _, child := range children {
				pid, err := strconv.Atoi(child)
				if err == nil {
					p.Children = append(p.Children, pid)
				}
			}
		}

		process = append(process, p)
		fmt.Println(p)
	}
}

