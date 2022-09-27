package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type MBR = struct {
	Tamano    [100]byte
	Fecha     [100]byte
	Signature [100]byte
	Fit       [100]byte
	Part1     Partition
	Part2     Partition
	Part3     Partition
	Part4     Partition
}

type Partition = struct {
	Status [100]byte
	Type   [100]byte
	Fit    [100]byte
	Start  [100]byte
	Size   [100]byte
	Name   [100]byte
}

type EBR = struct {
	Status [100]byte
	Fit    [100]byte
	Start  [100]byte
	Size   [100]byte
	Name   [100]byte
	Next   [100]byte
}

func main() {
	analizar()
}

func analizar() {
	finalizar := false
	reader := bufio.NewReader(os.Stdin)
	//  Ciclo para lectura de multiples comandos
	for !finalizar {
		fmt.Print("[MIA]@Proyecto2:~$  ")
		comando, _ := reader.ReadString('\n')
		if strings.Contains(comando, "exit") {
			finalizar = true
		} else {
			if comando != "" && comando != "exit\n" {
				//  Separacion de comando y parametros
				split_comando(comando)
			}
		}
	}
}

func split_comando(comando string) {
	var commandArray []string
	// Eliminacion de saltos de linea
	comando = strings.Replace(comando, "\n", "", 1)
	comando = strings.Replace(comando, "\r", "", 1)
	// Guardado de parametros
	if strings.Contains(comando, "mostrar") {
		commandArray = append(commandArray, comando)
	} else {
		commandArray = strings.Split(comando, " ")
	}
	// Ejecicion de comando leido
	//ejecucion_comando(commandArray)
}
