package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
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
	//if strings.Contains(comando, "mostrar") {
	//	commandArray = append(commandArray, comando)
	//} else {
	//	commandArray = strings.Split(comando, " ")
	//}

	commandArray = strings.Split(comando, " ")

	// Ejecicion de comando leido
	ejecucion_comando(commandArray)
}

func ejecucion_comando(commandArray []string) {
	// Identificacion de comando y ejecucion
	data := strings.ToLower(commandArray[0])
	if data == "mkdisk" {
		crear_disco(commandArray)
	} else {
		fmt.Println("Comando ingresado no es valido")
	}
}

// crear_disco -tamaño=numero -dimensional=dimension/"dimension"
func crear_disco(commandArray []string) {
	tamano := 0
	dimensional := " "
	fit := " "
	path := ""
	tamano_archivo := 0
	limite := 0
	bloque := make([]byte, 1024)

	// Lectura de parametros del comando
	for i := 0; i < len(commandArray); i++ {
		data := strings.ToLower(commandArray[i])
		if strings.Contains(data, "-size=") {
			strtam := strings.Replace(data, "-size=", "", 1)
			strtam = strings.Replace(strtam, "\"", "", 2)
			strtam = strings.Replace(strtam, "\r", "", 1)
			tamano2, err := strconv.Atoi(strtam)
			tamano = tamano2
			if err != nil {
				msg_error(err)
			}
		} else if strings.Contains(data, "-unit=") {
			dimensional = strings.Replace(data, "-unit=", "", 1)
			dimensional = strings.Replace(dimensional, "\"", "", 2)
		} else if strings.Contains(data, "-fit=") {
			fit = strings.Replace(data, "-fit=", "", 1)
			fit = strings.Replace(fit, "\"", "", 2)
		} else if strings.Contains(data, "-path=") {
			path = strings.Replace(data, "-path=", "", 1)
			path = strings.Replace(path, "\"", "", 2)
		}
	}

	if (tamano > 0) && (path != "") {
		nmbr := MBR{}
		// Calculo de tamaño del archivo
		if strings.Contains(dimensional, "k") {
			tamano_archivo = tamano
			copy(nmbr.Tamano[:], strconv.Itoa(tamano_archivo*1024))
		} else if strings.Contains(dimensional, "m") || strings.Contains(dimensional, " ") {
			tamano_archivo = tamano * 1024
			copy(nmbr.Tamano[:], strconv.Itoa(tamano_archivo*1024))
		} else if strings.Contains(dimensional, "g") {
			tamano_archivo = tamano * 1024 * 1024
			copy(nmbr.Tamano[:], strconv.Itoa(tamano_archivo*1024))
		} else {
			fmt.Print("Error: Dimensional No Reconocida.")
		}

		// Calculo de FIT
		if strings.Contains(fit, "bf") {
			copy(nmbr.Fit[:], "BF")
		} else if strings.Contains(fit, "ff") || strings.Contains(fit, " ") {
			copy(nmbr.Fit[:], "FF")
		} else if strings.Contains(fit, "wf") {
			copy(nmbr.Fit[:], "WF")
		} else {
			fmt.Print("Error: Fit No Reconocido.")
		}

		// Preparacion del bloque a escribir en archivo
		for j := 0; j < 1024; j++ {
			bloque[j] = 0
		}

		//Creando Directorio
		directorio := ""
		carpetas := strings.Split(path, "/")

		for j := 0; j < len(carpetas)-1; j++ {
			directorio += carpetas[j] + "/"
		}

		directorio = strings.TrimRight(directorio, "/")
		crearDirectorioSiNoExiste(directorio)

		// Creacion, escritura y cierre de archivo
		disco, err := os.Create(path)
		if err != nil {
			msg_error(err)
		}
		for limite < tamano_archivo {
			_, err := disco.Write(bloque)
			if err != nil {
				msg_error(err)
			}
			limite++
		}
		disco.Close()

		//Firma Aleatoria
		rand.Seed(time.Now().UnixNano())
		copy(nmbr.Signature[:], strconv.Itoa(rand.Intn(1000)))

		//Fecha De Creación
		dt := time.Now()
		copy(nmbr.Fecha[:], dt.Format("01-02-2006 15:04:05"))

		//Particiones Con Status F (Inactiva)
		copy(nmbr.Part1.Status[:], "F")
		copy(nmbr.Part2.Status[:], "F")
		copy(nmbr.Part3.Status[:], "F")
		copy(nmbr.Part4.Status[:], "F")

		//Escribir MBR
		disco2, err := os.OpenFile(path, os.O_RDWR, 0660)
		if err != nil {
			msg_error(err)
		}

		// Conversion de struct a bytes
		mbrbyte := struct_to_bytes(nmbr)
		// Cambio de posicion de puntero dentro del archivo
		newpos, err := disco2.Seek(int64(0), os.SEEK_SET)
		if err != nil {
			msg_error(err)
		}
		// Escritura de struct en archivo binario
		_, err = disco2.WriteAt(mbrbyte, newpos)
		if err != nil {
			msg_error(err)
		}

		// Resumen de accion realizada
		fmt.Print("Creacion de Disco:")
		fmt.Print(" Tamaño: ")
		fmt.Print(tamano)
		fmt.Print(" Dimensional: ")
		fmt.Println(dimensional)

	} else {
		msg_parametrosObligatorios()
	}
}

func struct_to_bytes(p interface{}) []byte {
	// Codificacion de Struct a []Bytes
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil && err != io.EOF {
		msg_error(err)
	}
	return buf.Bytes()
}

func msg_error(err error) {
	fmt.Println("Error: ", err)
}

func msg_parametrosObligatorios() {
	fmt.Println("Error: Parametros Obligatorios No Definidos.")
}

func crearDirectorioSiNoExiste(directorio string) {
	if _, err := os.Stat(directorio); os.IsNotExist(err) {
		err = os.MkdirAll(directorio, 0755)
		if err != nil {
			// Aquí puedes manejar mejor el error, es un ejemplo
			panic(err)
		}
	}
}
