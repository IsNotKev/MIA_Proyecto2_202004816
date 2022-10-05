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
	} else if data == "rmdisk" {
		eliminar_disco(commandArray)
	} else if data == "fdisk" {
		crear_particion(commandArray)
	} else {
		fmt.Println("Comando ingresado no es valido")
	}
}

//Crear Particion
func crear_particion(commandArray []string) {
	tamano := 0
	dimensional := " "
	path := ""
	tipo := ""
	fit := " "
	name := ""

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
		} else if strings.Contains(data, "-type=") {
			tipo = strings.Replace(data, "-type=", "", 1)
			tipo = strings.Replace(tipo, "\"", "", 2)
		} else if strings.Contains(data, "-name=") {
			name = strings.Replace(data, "-name=", "", 1)
			name = strings.Replace(name, "\"", "", 2)
		} else if strings.Contains(data, "-path=\"") {
			ultimo := data[len(data)-1:]
			path = data
			indice := i + 1
			for ultimo != "\"" {
				path += " " + strings.ToLower(commandArray[indice])
				ultimo = path[len(path)-1:]
				indice++
			}
			i = indice - 1
			path = strings.Replace(path, "-path=", "", 1)
			path = strings.Replace(path, "\"", "", 2)
		} else if strings.Contains(data, "-path=") {
			path = strings.Replace(data, "-path=", "", 1)
		}
	}

	if (tamano > 0) && (path != "") && (name != "") {
		mbrleido := leerMBR(path)

		nuevaP := Partition{}
		copy(nuevaP.Status[:], "V")

		// Calculo de tamaño de la Particion
		if strings.Contains(dimensional, "b") {
			copy(nuevaP.Size[:], strconv.Itoa(tamano))
		} else if strings.Contains(dimensional, "k") || strings.Contains(dimensional, " ") {
			copy(nuevaP.Size[:], strconv.Itoa(tamano*1024))
		} else if strings.Contains(dimensional, "m") {
			copy(nuevaP.Size[:], strconv.Itoa(tamano*1024*1024))
		} else {
			fmt.Print("Error: Dimensional No Reconocida.")
		}

		// Calculo de FIT
		if strings.Contains(fit, "bf") {
			copy(nuevaP.Fit[:], "BF")
		} else if strings.Contains(fit, "ff") {
			copy(nuevaP.Fit[:], "FF")
		} else if strings.Contains(fit, "wf") || strings.Contains(fit, " ") {
			copy(nuevaP.Fit[:], "WF")
		} else {
			fmt.Print("Error: Fit No Reconocido.")
		}

		// Tipo de Particion
		if strings.Contains(tipo, "e") {
			copy(nuevaP.Type[:], "E")
		} else if strings.Contains(tipo, "p") || strings.Contains(fit, " ") {
			copy(nuevaP.Type[:], "P")
		} else if strings.Contains(fit, "l") {
			copy(nuevaP.Type[:], "L")
		} else {
			fmt.Print("Error: Tipo De Particion No Reconocido.")
		}

		// Nombre de la particion
		copy(nuevaP.Name[:], name)

		if string(nuevaP.Type[:]) == "P" {
			if string(mbrleido.Part1.Status[:]) == "F" {
				mbrleido.Part1 = nuevaP
			} else if string(mbrleido.Part2.Status[:]) == "F" {
				mbrleido.Part2 = nuevaP
			} else if string(mbrleido.Part3.Status[:]) == "F" {
				mbrleido.Part3 = nuevaP
			} else if string(mbrleido.Part4.Status[:]) == "F" {
				mbrleido.Part4 = nuevaP
			} else {
				fmt.Print("Error: Ya hay un máximo de 4 particiones activas.")
			}
		} else if string(nuevaP.Type[:]) == "E" {

		} else if string(nuevaP.Type[:]) == "L" {

		}

	} else {
		msg_parametrosObligatorios()
	}
}

func leerMBR(ruta string) MBR {
	mbr_empty := MBR{}

	// Apertura de archivo
	disco, err := os.OpenFile(ruta, os.O_RDWR, 0660)
	if err != nil {
		msg_error(err)
	}
	// Calculo del tamano de struct en bytes
	ejm2 := struct_to_bytes(mbr_empty)
	sstruct := len(ejm2)

	lectura := make([]byte, sstruct)
	_, err = disco.ReadAt(lectura, int64(0))
	if err != nil && err != io.EOF {
		msg_error(err)
	}

	mbrleido := bytes_to_struct(lectura)

	fmt.Print("Fecha: ")
	fmt.Println(string(mbrleido.Fecha[:]))
	fmt.Print("Firma: ")
	fmt.Println(string(mbrleido.Signature[:]))
	fmt.Print("Tamaño: ")
	fmt.Println(string(mbrleido.Tamano[:]))
	fmt.Print("Fit: ")
	fmt.Println(string(mbrleido.Fit[:]))

	return mbrleido
}

func escribirMBR(mbr MBR, ruta string) {
	//Escribir MBR
	disco2, err := os.OpenFile(ruta, os.O_RDWR, 0660)
	if err != nil {
		msg_error(err)
	}

	// Conversion de struct a bytes
	mbrbyte := struct_to_bytes(mbr)
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
	disco2.Close()
}

//Eliminar Disco
func eliminar_disco(commandArray []string) {
	path := ""
	// Lectura de parametros del comando
	for i := 0; i < len(commandArray); i++ {
		data := strings.ToLower(commandArray[i])
		if strings.Contains(data, "-path=\"") {
			ultimo := data[len(data)-1:]
			path = data
			indice := i + 1
			for ultimo != "\"" {
				path += " " + strings.ToLower(commandArray[indice])
				ultimo = path[len(path)-1:]
				indice++
			}
			i = indice - 1
			path = strings.Replace(path, "-path=", "", 1)
			path = strings.Replace(path, "\"", "", 2)
		} else if strings.Contains(data, "-path=") {
			path = strings.Replace(data, "-path=", "", 1)
		}
	}
	if path != "" {
		err := os.Remove(path)
		if err != nil {
			fmt.Printf("Error eliminando archivo: %v\n", err)
		} else {
			fmt.Println("Eliminado correctamente")
		}
	} else {
		msg_parametrosObligatorios()
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
		} else if strings.Contains(data, "-path=\"") {
			ultimo := data[len(data)-1:]
			path = data
			indice := i + 1
			for ultimo != "\"" {
				path += " " + strings.ToLower(commandArray[indice])
				ultimo = path[len(path)-1:]
				indice++
			}
			i = indice - 1
			path = strings.Replace(path, "-path=", "", 1)
			path = strings.Replace(path, "\"", "", 2)
		} else if strings.Contains(data, "-path=") {
			path = strings.Replace(data, "-path=", "", 1)
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
		copy(nmbr.Fecha[:], dt.Format("02-01-2006 15:04:05"))

		//Particiones Con Status F (Inactiva)
		copy(nmbr.Part1.Status[:], "F")
		copy(nmbr.Part2.Status[:], "F")
		copy(nmbr.Part3.Status[:], "F")
		copy(nmbr.Part4.Status[:], "F")

		escribirMBR(nmbr, path)

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

func bytes_to_struct(s []byte) MBR {
	// Decodificacion de [] Bytes a Struct ejemplo
	p := MBR{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil && err != io.EOF {
		msg_error(err)
	}
	return p
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
