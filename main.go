package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/goccy/go-graphviz"
	"github.com/rs/cors"
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
	Status      [100]byte
	Type        [100]byte
	Fit         [100]byte
	Start       [100]byte
	Size        [100]byte
	Name        [100]byte
	Particiones [10]EBR
}

type EBR = struct {
	Status [100]byte
	Fit    [100]byte
	Start  [100]byte
	Size   [100]byte
	Name   [100]byte
	Next   [100]byte
}

type DiscoMontado = struct {
	path string
	name string
	id   string
	mbr  MBR
	num  int
}

type SuperBloque = struct {
	s_filesystem_type   [100]byte
	s_inodes_count      [100]byte
	s_blocks_count      [100]byte
	s_free_blocks_count [100]byte
	s_free_inodes_count [100]byte
	s_mtime             [100]byte
	s_mnt_count         [100]byte
	s_magic             [100]byte
	s_inode_size        [100]byte
	s_block_size        [100]byte
	s_firts_ino         [100]byte
	s_first_blo         [100]byte
	s_bm_inode_start    [100]byte
	s_bm_block_start    [100]byte
	s_inode_start       [100]byte
	s_block_start       [100]byte
}

type TablaInodos = struct {
	i_uid   [100]byte
	i_gid   [100]byte
	i_size  [100]byte
	i_atime [100]byte
	i_ctime [100]byte
	i_mtime [100]byte
	i_block [100]byte
	i_type  [100]byte
	i_perm  [100]byte
}

type BloqueArchivo = struct {
	b_content [100]byte
}

type BloqueCarpeta = struct {
	b_content [4]Content
}

type Content = struct {
	b_name  [100]byte
	b_inodo [100]byte
}

type cmdstruct struct {
	Cmd string `json:"cmd"`
}

type Usuario struct {
	usuario string
	contra  string
	id      string
	ruta    string
	nogrupo string
}

var discos = [20]DiscoMontado{} //Discos Montados
var cant = 1
var ultimoG = ""
var usuarioLogueado = Usuario{}

func main() {
	//analizar()
	fmt.Println("MIA - T4, API Rest GO")
	mux := http.NewServeMux()

	mux.HandleFunc("/ejecutar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var Content cmdstruct
		respuesta := "Conectado"
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &Content)

		respuesta = analizar(Content.Cmd)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "` + respuesta + `" }`))
	})

	fmt.Println("Server ON in port 5000")
	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":5000", handler))
}

func analizar(cmd string) string {

	instrucciones := strings.Split(cmd, "\n")
	salida := ""
	//finalizar := false
	//reader := bufio.NewReader(os.Stdin)
	////  Ciclo para lectura de multiples comandos
	//for !finalizar {
	//	fmt.Print("[MIA]@Proyecto2:~$  ")
	//	comando, _ := reader.ReadString('\n')
	//	if strings.Contains(comando, "exit") {
	//		finalizar = true
	//	} else {
	//		if comando != "" && comando != "exit\n" {
	//			//  Separacion de comando y parametros
	//			split_comando(comando)
	//		}
	//	}
	//}

	for i := 0; i < len(instrucciones); i++ {
		comando := instrucciones[i]
		if comando != "" && comando != "exit\n" {
			//  Separacion de comando y parametros
			salida += split_comando(comando) + "\\n"
		}
	}

	fmt.Println(salida)

	return salida
}

func split_comando(comando string) string {
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
	return ejecucion_comando(commandArray)
}

func ejecucion_comando(commandArray []string) string {
	// Identificacion de comando y ejecucion
	data := strings.ToLower(commandArray[0])
	if data == "mkdisk" {
		return crear_disco(commandArray)
	} else if data == "rmdisk" {
		return eliminar_disco(commandArray)
	} else if data == "fdisk" {
		return crear_particion(commandArray)
	} else if data == "mount" {
		return montar_disco(commandArray)
	} else if data == "mkfs" {
		return procesarMkfs(commandArray)
	} else if data == "login" {
		return procesarLogin(commandArray)
	} else if data == "logout" {
		return procesarLogout(commandArray)
	} else if data == "mkgrp" {
		return crearGrupo(commandArray)
	} else if data == "rmgrp" {
		return eliminarGrupo(commandArray)
	} else if data == "mkuser" {
		return crearUsuario(commandArray)
	} else if data == "rmuser" {
		return eliminarUsuario(commandArray)
	} else if data == "pause" {
		return "Pausa ..."
	} else if data == "rep" {
		return procesarRep(commandArray)
	} else {
		fmt.Println("Comando ingresado no es valido")
		return "Comando ingresado no es valido"
	}
}

func eliminarUsuario(commandArray []string) string {
	usuario := ""
	// Lectura de parametros del comando
	for i := 0; i < len(commandArray); i++ {
		data := strings.ToLower(commandArray[i])
		if strings.Contains(data, "-usuario=") {
			usuario = strings.Replace(data, "-usuario=", "", 1)
		}
	}

	if usuario != "" {
		if usuarioLogueado.usuario != "" {
			if usuarioLogueado.usuario == "root" && usuarioLogueado.nogrupo == "1" {
				content, err := ioutil.ReadFile(usuarioLogueado.ruta)

				if err != nil {
					log.Fatal(err)
				}

				nuevoContent := eliminarUs(usuario, string(content))
				b := []byte(nuevoContent)
				err = ioutil.WriteFile(usuarioLogueado.ruta, b, 0644)
				if err != nil {
					log.Fatal(err)
				}

				return "> Usuario " + usuario + " eliminado."
			} else {
				return "Solo el usuario root puede eliminar grupos."
			}
		} else {
			return "No hay usuario en sesión para eliminar grupo."
		}
	} else {
		return msg_parametrosObligatorios()
	}

}

func eliminarUs(usuario string, texto string) string {
	usuariosArray := strings.Split(texto, "\n")

	nuevoT := ""

	for i := 0; i < len(usuariosArray); i++ {
		usuarioLeido := strings.Replace(usuariosArray[i], " ", "", 10)
		credencialesLeidas := strings.Split(usuarioLeido, ",")

		aux := usuariosArray[i] + "\n"

		if credencialesLeidas[1] == "U" {
			if credencialesLeidas[3] == usuario {
				aux = "0, U, " + credencialesLeidas[2] + ", " + credencialesLeidas[3] + ", " + credencialesLeidas[4] + ", " + "\n"
			}
		}
		nuevoT += aux
	}
	nuevoT = strings.TrimRight(nuevoT, "\n")
	return nuevoT
}

func crearUsuario(commandArray []string) string {
	usuario := ""
	contra := ""
	grp := ""
	// Lectura de parametros del comando
	for i := 0; i < len(commandArray); i++ {
		data := strings.ToLower(commandArray[i])
		if strings.Contains(data, "-usuario=") {
			usuario = strings.Replace(data, "-usuario=", "", 1)
		} else if strings.Contains(data, "-pwd=") {
			contra = strings.Replace(data, "-pwd=", "", 1)
		} else if strings.Contains(data, "-grp=") {
			grp = strings.Replace(data, "-grp=", "", 1)
		}
	}

	if usuario != "" && contra != "" && grp != "" {
		if usuarioLogueado.usuario != "" {
			if usuarioLogueado.usuario == "root" && usuarioLogueado.nogrupo == "1" {
				content, err := ioutil.ReadFile(usuarioLogueado.ruta)

				if err != nil {
					log.Fatal(err)
				}

				if verificarGrupo(grp, string(content)) {

					nuevoContent := agregarUs(usuario, contra, grp, string(content))
					b := []byte(nuevoContent)
					err = ioutil.WriteFile(usuarioLogueado.ruta, b, 0644)
					if err != nil {
						log.Fatal(err)
					}

					return "> Usuario " + usuario + " creado en " + grp + "."
				} else {
					return "Grupo " + grp + " no existe."
				}
			} else {
				return "Solo el usuario root puede eliminar grupos."
			}
		} else {
			return "No hay usuario en sesión para eliminar grupo."
		}
	} else {
		return msg_parametrosObligatorios()
	}
}

func agregarUs(usuario string, contra string, grupo string, texto string) string {
	usuariosArray := strings.Split(texto, "\n")

	nuevoT := ""

	for i := 0; i < len(usuariosArray); i++ {
		usuarioLeido := strings.Replace(usuariosArray[i], " ", "", 10)
		credencialesLeidas := strings.Split(usuarioLeido, ",")

		aux := usuariosArray[i] + "\n"

		if credencialesLeidas[1] == "G" {
			if credencialesLeidas[2] == grupo {
				aux += credencialesLeidas[0] + ", U, " + credencialesLeidas[2] + ", " + usuario + ", " + contra + "\n"
			}
		}
		nuevoT += aux
	}
	nuevoT = strings.TrimRight(nuevoT, "\n")
	return nuevoT
}

func eliminarGrupo(commandArray []string) string {
	name := ""
	// Lectura de parametros del comando
	for i := 0; i < len(commandArray); i++ {
		data := strings.ToLower(commandArray[i])
		if strings.Contains(data, "-name=") {
			name = strings.Replace(data, "-name=", "", 1)
		}
	}

	if name != "" {
		if usuarioLogueado.usuario != "" {
			if usuarioLogueado.usuario == "root" && usuarioLogueado.nogrupo == "1" {
				content, err := ioutil.ReadFile(usuarioLogueado.ruta)

				if err != nil {
					log.Fatal(err)
				}

				if verificarGrupo(name, string(content)) {

					nuevoContent := eliminarG(name, string(content))
					b := []byte(nuevoContent)
					err = ioutil.WriteFile(usuarioLogueado.ruta, b, 0644)
					if err != nil {
						log.Fatal(err)
					}

					return "> Grupo " + name + " eliminado."
				} else {
					return "Grupo " + name + " no existe."
				}
			} else {
				return "Solo el usuario root puede eliminar grupos."
			}
		} else {
			return "No hay usuario en sesión para eliminar grupo."
		}
	} else {
		return msg_parametrosObligatorios()
	}
}

func eliminarG(nombre string, texto string) string {
	usuariosArray := strings.Split(texto, "\n")

	nuevoT := ""

	for i := 0; i < len(usuariosArray); i++ {
		usuarioLeido := strings.Replace(usuariosArray[i], " ", "", 10)
		credencialesLeidas := strings.Split(usuarioLeido, ",")

		aux := usuariosArray[i] + "\n"

		if credencialesLeidas[1] == "G" {
			if credencialesLeidas[2] == nombre {
				aux = "0, " + credencialesLeidas[1] + ", " + credencialesLeidas[2] + "\n"
			}
		}
		nuevoT += aux
	}
	nuevoT = strings.TrimRight(nuevoT, "\n")
	return nuevoT
}

func crearGrupo(commandArray []string) string {
	name := ""
	// Lectura de parametros del comando
	for i := 0; i < len(commandArray); i++ {
		data := strings.ToLower(commandArray[i])
		if strings.Contains(data, "-name=") {
			name = strings.Replace(data, "-name=", "", 1)
		}
	}

	if name != "" {

		if usuarioLogueado.usuario != "" {
			if usuarioLogueado.usuario == "root" && usuarioLogueado.nogrupo == "1" {
				content, err := ioutil.ReadFile(usuarioLogueado.ruta)

				if err != nil {
					log.Fatal(err)
				}

				if !verificarGrupo(name, string(content)) {
					ultimo, err := strconv.Atoi(ultimoG)
					if err != nil {
						msg_error(err)
					}
					ultimo++

					nuevoG := strconv.Itoa(ultimo)

					nuevoContent := string(content) + "\n" + nuevoG + ", G, " + name

					b := []byte(nuevoContent)
					err = ioutil.WriteFile(usuarioLogueado.ruta, b, 0644)
					if err != nil {
						log.Fatal(err)
					}

					return "> Grupo " + name + " creado."
				} else {
					return "Grupo " + name + " ya existe."
				}

			} else {
				return "Solo el usuario root puede crear grupos."
			}
		} else {
			return "No hay usuario en sesión para crear grupo."
		}

	} else {
		return msg_parametrosObligatorios()
	}
}

func verificarGrupo(name string, texto string) bool {
	usuariosArray := strings.Split(texto, "\n")

	for i := 0; i < len(usuariosArray); i++ {
		usuarioLeido := strings.Replace(usuariosArray[i], " ", "", 10)
		credencialesLeidas := strings.Split(usuarioLeido, ",")

		if credencialesLeidas[1] == "G" {
			if credencialesLeidas[2] == name && credencialesLeidas[0] != "0" {
				return true
			}
			if credencialesLeidas[0] != "0" {
				ultimoG = credencialesLeidas[0]
			}
		}
	}
	return false
}

func procesarLogout(commandArray []string) string {
	if usuarioLogueado.usuario != "" {
		aux := usuarioLogueado.usuario
		usuarioLogueado.usuario = ""
		usuarioLogueado.contra = ""
		usuarioLogueado.ruta = ""
		usuarioLogueado.id = ""
		usuarioLogueado.nogrupo = ""
		return "> Cerrando sesión de " + aux
	} else {
		return "No hay sesión activa que cerrar."
	}
}

func procesarLogin(commandArray []string) string {
	usuario := ""
	contra := ""
	id := ""
	// Lectura de parametros del comando
	for i := 0; i < len(commandArray); i++ {
		data := strings.ToLower(commandArray[i])
		if strings.Contains(data, "-id=") {
			id = strings.Replace(data, "-id=", "", 1)
		} else if strings.Contains(data, "-password=") {
			contra = strings.Replace(data, "-password=", "", 1)
		} else if strings.Contains(data, "-usuario=") {
			usuario = strings.Replace(data, "-usuario=", "", 1)
		}
	}

	if usuario != "" && contra != "" && id != "" {
		if usuarioLogueado.usuario != "" {
			return "Ya hay un usuario logueado."
		}

		for i := 0; i < len(discos); i++ {
			if discos[i].id == id {

				//Creando Directorio
				directorio := ""
				carpetas := strings.Split(discos[i].path, "/")

				for j := 0; j < len(carpetas)-1; j++ {
					directorio += carpetas[j] + "/"
				}

				directorio += "users.txt"

				content, err := ioutil.ReadFile(directorio)

				if err != nil {
					log.Fatal(err)
				}

				if verificarLogin(usuario, contra, string(content)) {
					usuarioLogueado.id = id
					usuarioLogueado.ruta = directorio
					return "> " + usuario + " logueado en" + id + "."
				} else {
					return "> Credenciales incorrectas para login en " + id + "."
				}
			}
		}

		return "No se encontro id para Login."
	} else {
		return msg_parametrosObligatorios()
	}
}

func verificarLogin(usuario string, contra string, texto string) bool {
	usuariosArray := strings.Split(texto, "\n")

	for i := 0; i < len(usuariosArray); i++ {
		usuarioLeido := strings.Replace(usuariosArray[i], " ", "", 10)
		credencialesLeidas := strings.Split(usuarioLeido, ",")

		if credencialesLeidas[1] == "U" {
			if credencialesLeidas[3] == usuario && credencialesLeidas[4] == contra {
				usuarioLogueado.usuario = usuario
				usuarioLogueado.contra = contra
				usuarioLogueado.nogrupo = credencialesLeidas[0]
				return true
			}
		}
	}
	return false
}

func procesarMkfs(commandArray []string) string {
	id := ""
	tipo := ""

	// Lectura de parametros del comando
	for i := 0; i < len(commandArray); i++ {
		data := strings.ToLower(commandArray[i])
		if strings.Contains(data, "-id=") {
			id = strings.Replace(data, "-id=", "", 1)
		} else if strings.Contains(data, "-type=") {
			tipo = strings.Replace(data, "-type=", "", 1)
		}
	}

	if id != "" {
		if tipo == "full" || tipo == "" {

			for i := 0; i < len(discos); i++ {
				if discos[i].id == id {

					//Creando Directorio
					directorio := ""
					carpetas := strings.Split(discos[i].path, "/")

					for j := 0; j < len(carpetas)-1; j++ {
						directorio += carpetas[j] + "/"
					}

					directorio += "users.txt"

					b := []byte("1, G, root \n1, U, root, root, 123")
					err := ioutil.WriteFile(directorio, b, 0644)
					if err != nil {
						log.Fatal(err)
					}

					return "> " + id + " Formateado con ext2."
				}
			}

			return "No se encontro id para MKFS."
		} else {
			return "Tipo de Formateo no reconocido."
		}
	} else {
		return msg_parametrosObligatorios()
	}
}

//Reportes
func procesarRep(commandArray []string) string {
	path := ""
	name := ""
	id := ""
	ruta := ""

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
		} else if strings.Contains(data, "-name=") {
			name = strings.Replace(data, "-name=", "", 1)
			name = strings.Replace(name, "\"", "", 2)
		} else if strings.Contains(data, "-id=") {
			id = strings.Replace(data, "-id=", "", 1)
			id = strings.Replace(id, "\"", "", 2)
		} else if strings.Contains(data, "-ruta=") {
			ruta = strings.Replace(data, "-ruta=", "", 1)
			ruta = strings.Replace(ruta, "\"", "", 2)
		}
	}

	if name != "" && path != "" && id != "" {
		if name == "disk" {
			return repDisk(path, id)
		} else if name == "tree" {
			return "Reporte Tree Creado"
		} else if name == "file" {
			return "Reporte File Creado"
		} else if name == "sb" {
			return "Reporte SB Creado"
		}
		return "Reporte no encontrado."
	} else {
		return msg_parametrosObligatorios()
	}
}

func repDisk(path string, id string) string {
	for i := 0; i < len(discos); i++ {
		if discos[i].id == id {
			mbrleido := leerMBR(discos[i].path)
			particionesaux := [4]Partition{}
			particionesaux[0] = mbrleido.Part1
			particionesaux[1] = mbrleido.Part2
			particionesaux[2] = mbrleido.Part3
			particionesaux[3] = mbrleido.Part4

			g := graphviz.New()
			graph, err := g.Graph()
			if err != nil {
				log.Fatal(err)
			}

			// create your graph
			n, err := graph.CreateNode("n")
			if err != nil {
				log.Fatal(err)
			}
			n.SetShape("record")

			text := "MBR"

			inicio := int(unsafe.Sizeof(mbrleido))
			for i := 0; i < 4; i++ {
				if CToGoString(particionesaux[i].Status) == "V" {

					if CToGoString(particionesaux[i].Type) == "P" {
						t, err := strconv.Atoi(CToGoString(particionesaux[i].Size))
						if err != nil {
							msg_error(err)
						}
						tt, err := strconv.Atoi(CToGoString(mbrleido.Tamano))
						if err != nil {
							msg_error(err)
						}
						percent := float64(t) / float64(tt)
						percent = percent * 100
						text += " | Primaria " + fmt.Sprintf("%.2f", percent) + "%"
					} else if CToGoString(particionesaux[i].Type) == "E" {
						text += " | {{Extendida}|{"

						for j := 0; j < 10; j++ {
							if CToGoString(particionesaux[i].Particiones[j].Status) == "V" {
								t, err := strconv.Atoi(CToGoString(particionesaux[i].Particiones[j].Size))
								if err != nil {
									msg_error(err)
								}
								tt, err := strconv.Atoi(CToGoString(mbrleido.Tamano))
								if err != nil {
									msg_error(err)
								}
								percent := float64(t) / float64(tt)
								percent = percent * 100
								text += "EBR | Lógica " + fmt.Sprintf("%.2f", percent) + "%|"
							}

						}
						text += "}}"
					}

					st, err := strconv.Atoi(CToGoString(particionesaux[i].Start))
					if err != nil {
						msg_error(err)
					}
					tam, err := strconv.Atoi(CToGoString(particionesaux[i].Size))
					if err != nil {
						msg_error(err)
					}
					inicio = st + tam + 1
				} else if i == 3 {
					t, err := strconv.Atoi(CToGoString(mbrleido.Tamano))
					if err != nil {
						msg_error(err)
					}
					t = t - int(inicio)
					tt, err := strconv.Atoi(CToGoString(mbrleido.Tamano))
					if err != nil {
						msg_error(err)
					}
					percent := float64(t) / float64(tt)
					percent = percent * 100

					text += " | Libre " + fmt.Sprintf("%.2f", percent) + "%"
				}
			}

			n.SetLabel(text)

			// 1. write encoded PNG data to buffer
			var buf bytes.Buffer
			if err := g.Render(graph, graphviz.PNG, &buf); err != nil {
				log.Fatal(err)
			}

			//Creando Directorio
			directorio := ""
			carpetas := strings.Split(path, "/")

			for j := 0; j < len(carpetas)-1; j++ {
				directorio += carpetas[j] + "/"
			}

			directorio = strings.TrimRight(directorio, "/")
			crearDirectorioSiNoExiste(directorio)

			// 3. write to file directly
			if err := g.RenderFilename(graph, graphviz.PNG, path); err != nil {
				log.Fatal(err)
			}
			return "Reporte Disk Creado"
		}
	}
	return "No se encontro disco montado para Reporte Disk"
}

//Montar Disco
func montar_disco(commandArray []string) string {
	path := ""
	name := ""
	salida := ""
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
		} else if strings.Contains(data, "-name=") {
			name = strings.Replace(data, "-name=", "", 1)
			name = strings.Replace(name, "\"", "", 2)
		}
	}

	if path != "" && name != "" {
		salida += "\\nDiscos Montados: \\n"
		mbrleido := leerMBR(path)
		nuevoDisco := DiscoMontado{}

		nuevoDisco.mbr = mbrleido
		nuevoDisco.name = name
		nuevoDisco.path = path

		cont := 0
		encontrado := false
		numaux := -1
		for i := 0; i < len(discos); i++ {
			if discos[i].path == "" {
				letra := obtenerLetra(cont)
				if numaux >= 0 {
					nuevoDisco.num = numaux
				} else {
					nuevoDisco.num = cant
				}
				nuevoDisco.id = "06" + strconv.Itoa(nuevoDisco.num) + letra
				//fmt.Println(nuevoDisco.id)
				discos[i] = nuevoDisco
				break
			} else if discos[i].path == path {
				cont++
				numaux = discos[i].num
				encontrado = true
				salida += "    > " + discos[i].id + "\\n"
			} else {
				salida += "    > " + discos[i].id + "\\n"
			}
		}

		if !encontrado {
			cant++
		}
		return salida + "Disco " + nuevoDisco.id + " montado.\\n"
	} else {
		return msg_parametrosObligatorios()
	}
}

func obtenerLetra(num int) string {
	if num == 0 {
		return "a"
	} else if num == 1 {
		return "b"
	} else if num == 2 {
		return "c"
	} else if num == 3 {
		return "d"
	} else if num == 4 {
		return "e"
	} else if num == 5 {
		return "f"
	} else if num == 6 {
		return "g"
	} else if num == 7 {
		return "h"
	} else if num == 8 {
		return "i"
	} else if num == 9 {
		return "j"
	} else {
		return "k"
	}
}

//Crear Particion
func crear_particion(commandArray []string) string {
	tamano := 0
	dimensional := " "
	path := ""
	tipo := " "
	fit := " "
	name := ""
	salida := ""

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
		} else {
			salida += "Parametro no identificado.\\n"
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
			tamano = tamano * 1024
			copy(nuevaP.Size[:], strconv.Itoa(tamano))
		} else if strings.Contains(dimensional, "m") {
			tamano = tamano * 1024 * 1024
			copy(nuevaP.Size[:], strconv.Itoa(tamano))
		} else {
			fmt.Print("Error: Dimensional No Reconocida.")
			return salida + "Error -> Dimensional No Reconocida."
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
			return salida + "Error -> Fit No Reconocido."
		}

		// Tipo de Particion
		if strings.Contains(tipo, "e") {
			copy(nuevaP.Type[:], "E")
		} else if strings.Contains(tipo, "p") || strings.Contains(tipo, " ") {
			copy(nuevaP.Type[:], "P")
		} else if strings.Contains(tipo, "l") {
			copy(nuevaP.Type[:], "L")
		} else {
			fmt.Print("Error: Tipo De Particion No Reconocido.")
			return "Error -> Tipo De Particion No Reconocido."
		}

		// Nombre de la particion
		copy(nuevaP.Name[:], name)

		if CToGoString(nuevaP.Type) == "P" { //**************** PRIMARIA *******************
			creada := false
			errorp := false
			// Calculo del tamano de struct en bytes
			ejm2 := struct_to_bytes(mbrleido)
			start := len(ejm2)

			fin, err := strconv.Atoi(CToGoString(mbrleido.Tamano))
			if err != nil {
				msg_error(err)
			}

			//Ver si part1 esta libre
			if CToGoString(mbrleido.Part1.Status) == "F" && !creada && !errorp {
				//Verificar que si haya espacio
				if start+tamano < fin {
					start = start + 1
					copy(nuevaP.Start[:], strconv.Itoa(start))
					mbrleido.Part1 = nuevaP
					creada = true
					fmt.Println("Particion Primaria Creada")
				} else {
					fmt.Println("Error -> No hay espacio suficiente para la particion")
					errorp = true
					return "Error -> No hay espacio suficiente para la particion"
				}
			} else if CToGoString(mbrleido.Part1.Status) == "V" {
				tamanopart, err := strconv.Atoi(CToGoString(mbrleido.Part1.Size))
				if err != nil {
					msg_error(err)
				}
				start = start + tamanopart

				if name == CToGoString(mbrleido.Part1.Name) {
					fmt.Println("Error -> El nombre no puede repetise dentro de las particiones")
					errorp = true
					return "Error -> El nombre no puede repetise dentro de las particiones"
				}
			}

			//Ver si part2 esta libre
			if CToGoString(mbrleido.Part2.Status) == "F" && !creada && !errorp {
				//Verificar que si haya espacio
				if start+tamano < fin {
					start = start + 1
					copy(nuevaP.Start[:], strconv.Itoa(start))
					mbrleido.Part2 = nuevaP
					creada = true
					fmt.Println("Particion Primaria Creada")
				} else {
					fmt.Println("Error -> No hay espacio suficiente para la particion")
					errorp = true
					return "Error -> No hay espacio suficiente para la particion"
				}
			} else if CToGoString(mbrleido.Part2.Status) == "V" {
				tamanopart, err := strconv.Atoi(CToGoString(mbrleido.Part2.Size))
				if err != nil {
					msg_error(err)
				}
				start = start + tamanopart

				if name == CToGoString(mbrleido.Part2.Name) {
					fmt.Println("Error -> El nombre no puede repetise dentro de las particiones")
					errorp = true
					return "Error -> El nombre no puede repetise dentro de las particiones"
				}
			}

			//Ver si part3 esta libre
			if CToGoString(mbrleido.Part3.Status) == "F" && !creada && !errorp {
				//Verificar que si haya espacio
				if start+tamano < fin {
					start = start + 1
					copy(nuevaP.Start[:], strconv.Itoa(start))
					mbrleido.Part3 = nuevaP
					creada = true
					fmt.Println("Particion Primaria Creada")
				} else {
					fmt.Println("Error -> No hay espacio suficiente para la particion")
					errorp = true
					return "Error -> No hay espacio suficiente para la particion"
				}
			} else if CToGoString(mbrleido.Part3.Status) == "V" {
				tamanopart, err := strconv.Atoi(CToGoString(mbrleido.Part3.Size))
				if err != nil {
					msg_error(err)
				}
				start = start + tamanopart

				if name == CToGoString(mbrleido.Part3.Name) {
					fmt.Println("Error -> El nombre no puede repetise dentro de las particiones")
					errorp = true
					return "Error -> El nombre no puede repetise dentro de las particiones"
				}
			}

			//Ver si part4 esta libre
			if CToGoString(mbrleido.Part4.Status) == "F" && !creada && !errorp {
				//Verificar que si haya espacio
				if start+tamano < fin {
					start = start + 1
					copy(nuevaP.Start[:], strconv.Itoa(start))
					mbrleido.Part4 = nuevaP
					creada = true
					fmt.Println("Particion Primaria Creada")
				} else {
					fmt.Println("Error -> No hay espacio suficiente para la particion")
					errorp = true
					return "Error -> No hay espacio suficiente para la particion"
				}
			} else if CToGoString(mbrleido.Part4.Status) == "V" {
				tamanopart, err := strconv.Atoi(CToGoString(mbrleido.Part4.Size))
				if err != nil {
					msg_error(err)
				}
				start = start + tamanopart

				if name == CToGoString(mbrleido.Part4.Name) {
					fmt.Println("Error -> El nombre no puede repetise dentro de las particiones")
					errorp = true
					return "Error -> El nombre no puede repetise dentro de las particiones"
				}
			}

			if !creada && !errorp {
				fmt.Println("Error: Ya hay un máximo de 4 particiones activas.")
				return "Error: Ya hay un máximo de 4 particiones activas."
			}
		} else if CToGoString(nuevaP.Type) == "E" { //**************** EXTENDIDA ***************

			if CToGoString(mbrleido.Part1.Type) == "E" || CToGoString(mbrleido.Part2.Type) == "E" || CToGoString(mbrleido.Part3.Type) == "E" || CToGoString(mbrleido.Part4.Type) == "E" {
				fmt.Println("Error: Ya existe un máximo de 1 partición extendida.")
				return "Error -> Ya existe un máximo de 1 partición extendida."
			} else {
				creada := false
				errorp := false
				// Calculo del tamano de struct en bytes
				ejm2 := struct_to_bytes(mbrleido)
				start := len(ejm2)

				fin, err := strconv.Atoi(CToGoString(mbrleido.Tamano))
				if err != nil {
					msg_error(err)
				}

				//Ver si part1 esta libre
				if CToGoString(mbrleido.Part1.Status) == "F" && !creada && !errorp {
					//Verificar que si haya espacio
					if start+tamano < fin {
						start = start + 1
						copy(nuevaP.Start[:], strconv.Itoa(start))
						mbrleido.Part1 = nuevaP
						creada = true
						fmt.Println("Particion Extendida Creada")
					} else {
						fmt.Println("Error -> No hay espacio suficiente para la particion")
						errorp = true
						return "Error -> No hay espacio suficiente para la particion"
					}
				} else if CToGoString(mbrleido.Part1.Status) == "V" {
					tamanopart, err := strconv.Atoi(CToGoString(mbrleido.Part1.Size))
					if err != nil {
						msg_error(err)
					}
					start = start + tamanopart

					if name == CToGoString(mbrleido.Part1.Name) {
						fmt.Println("Error -> El nombre no puede repetise dentro de las particiones")
						errorp = true
						return "Error -> El nombre no puede repetise dentro de las particiones"
					}
				}

				//Ver si part2 esta libre
				if CToGoString(mbrleido.Part2.Status) == "F" && !creada && !errorp {
					//Verificar que si haya espacio
					if start+tamano < fin {
						start = start + 1
						copy(nuevaP.Start[:], strconv.Itoa(start))
						mbrleido.Part2 = nuevaP
						creada = true
						fmt.Println("Particion Extendida Creada")
					} else {
						fmt.Println("Error -> No hay espacio suficiente para la particion")
						errorp = true
						return "Error -> No hay espacio suficiente para la particion"
					}
				} else if CToGoString(mbrleido.Part2.Status) == "V" {
					tamanopart, err := strconv.Atoi(CToGoString(mbrleido.Part2.Size))
					if err != nil {
						msg_error(err)
					}
					start = start + tamanopart

					if name == CToGoString(mbrleido.Part2.Name) {
						fmt.Println("Error -> El nombre no puede repetise dentro de las particiones")
						errorp = true
						return "Error -> El nombre no puede repetise dentro de las particiones"
					}
				}

				//Ver si part3 esta libre
				if CToGoString(mbrleido.Part3.Status) == "F" && !creada && !errorp {
					//Verificar que si haya espacio
					if start+tamano < fin {
						start = start + 1
						copy(nuevaP.Start[:], strconv.Itoa(start))
						mbrleido.Part3 = nuevaP
						creada = true
						fmt.Println("Particion Extendida Creada")
					} else {
						fmt.Println("Error -> No hay espacio suficiente para la particion")
						errorp = true
						return "Error -> No hay espacio suficiente para la particion"
					}
				} else if CToGoString(mbrleido.Part3.Status) == "V" {
					tamanopart, err := strconv.Atoi(CToGoString(mbrleido.Part3.Size))
					if err != nil {
						msg_error(err)
					}
					start = start + tamanopart

					if name == CToGoString(mbrleido.Part3.Name) {
						fmt.Println("Error -> El nombre no puede repetise dentro de las particiones")
						errorp = true
						return "Error -> El nombre no puede repetise dentro de las particiones"
					}
				}

				//Ver si part4 esta libre
				if CToGoString(mbrleido.Part4.Status) == "F" && !creada && !errorp {
					//Verificar que si haya espacio
					if start+tamano < fin {
						start = start + 1
						copy(nuevaP.Start[:], strconv.Itoa(start))
						mbrleido.Part4 = nuevaP
						creada = true
						fmt.Println("Particion Extendida Creada")
					} else {
						fmt.Println("Error -> No hay espacio suficiente para la particion")
						errorp = true
						return "Error -> No hay espacio suficiente para la particion"
					}
				} else if CToGoString(mbrleido.Part4.Status) == "V" {
					tamanopart, err := strconv.Atoi(CToGoString(mbrleido.Part4.Size))
					if err != nil {
						msg_error(err)
					}
					start = start + tamanopart

					if name == CToGoString(mbrleido.Part4.Name) {
						fmt.Println("Error -> El nombre no puede repetise dentro de las particiones")
						errorp = true
						return "Error -> El nombre no puede repetise dentro de las particiones"
					}
				}

				if !creada && !errorp {
					fmt.Println("Error: Ya hay un máximo de 4 particiones activas.")
					return "Error -> Ya hay un máximo de 4 particiones activas."
				}
			}
		} else if CToGoString(nuevaP.Type) == "L" {
			nuevoebr := EBR{}
			nuevoebr.Status = nuevaP.Status
			nuevoebr.Fit = nuevaP.Fit
			nuevoebr.Size = nuevaP.Size
			copy(nuevoebr.Next[:], "-1")
			nuevoebr.Name = nuevaP.Name

			if CToGoString(mbrleido.Part1.Type) == "E" {
				tamanoParticion, err := strconv.Atoi(CToGoString(mbrleido.Part1.Size))
				if err != nil {
					msg_error(err)
				}
				for i := 0; i < len(mbrleido.Part1.Particiones); i++ {
					if CToGoString(mbrleido.Part1.Particiones[i].Status) != "V" {
						if i > 0 {
							ultimo := i - 1
							ultimostart, err := strconv.Atoi(CToGoString(mbrleido.Part1.Particiones[ultimo].Start))
							if err != nil {
								msg_error(err)
							}
							ultimotamano, err := strconv.Atoi(CToGoString(mbrleido.Part1.Particiones[ultimo].Size))
							if err != nil {
								msg_error(err)
							}

							if ultimostart+ultimotamano+tamano <= tamanoParticion {
								copy(nuevoebr.Start[:], strconv.Itoa(ultimostart+ultimotamano+1))
								copy(mbrleido.Part1.Particiones[ultimo].Next[:], strconv.Itoa(ultimostart+ultimotamano+1))
								mbrleido.Part1.Particiones[i] = nuevoebr
								fmt.Println("Partición Logica Creada")
							} else {
								fmt.Println("Error: No hay almacenamiento para nueva partición lógica")
								return "Error -> No hay almacenamiento para nueva partición lógica"
							}

						} else {
							copy(nuevoebr.Start[:], "0")
							if tamano <= tamanoParticion {
								mbrleido.Part1.Particiones[i] = nuevoebr
								fmt.Println("Partición Logica Creada")
							} else {
								fmt.Println("Error: No hay almacenamiento para nueva partición lógica")
								return "Error -> No hay almacenamiento para nueva partición lógica"
							}
						}
						break
					} else if name == CToGoString(mbrleido.Part1.Particiones[i].Name) {
						fmt.Println("Error -> El nombre no puede repetise dentro de las particiones")
						return "Error -> El nombre no puede repetise dentro de las particiones"
					}
				}
			} else if CToGoString(mbrleido.Part2.Type) == "E" {
				tamanoParticion, err := strconv.Atoi(CToGoString(mbrleido.Part2.Size))
				if err != nil {
					msg_error(err)
				}
				for i := 0; i < len(mbrleido.Part2.Particiones); i++ {
					if CToGoString(mbrleido.Part2.Particiones[i].Status) != "V" {
						if i > 0 {
							ultimo := i - 1
							ultimostart, err := strconv.Atoi(CToGoString(mbrleido.Part2.Particiones[ultimo].Start))
							if err != nil {
								msg_error(err)
							}
							ultimotamano, err := strconv.Atoi(CToGoString(mbrleido.Part2.Particiones[ultimo].Size))
							if err != nil {
								msg_error(err)
							}

							if ultimostart+ultimotamano+tamano <= tamanoParticion {
								copy(nuevoebr.Start[:], strconv.Itoa(ultimostart+ultimotamano+1))
								copy(mbrleido.Part2.Particiones[ultimo].Next[:], strconv.Itoa(ultimostart+ultimotamano+1))
								mbrleido.Part2.Particiones[i] = nuevoebr
								fmt.Println("Partición Logica Creada")
							} else {
								fmt.Println("Error: No hay almacenamiento para nueva partición lógica")
								return "Error -> No hay almacenamiento para nueva partición lógica"
							}

						} else {
							copy(nuevoebr.Start[:], "0")
							if tamano <= tamanoParticion {
								mbrleido.Part2.Particiones[i] = nuevoebr
								fmt.Println("Partición Logica Creada")
							} else {
								fmt.Println("Error: No hay almacenamiento para nueva partición lógica")
								return "Error -> No hay almacenamiento para nueva partición lógica"
							}
						}
						break
					} else if name == CToGoString(mbrleido.Part2.Particiones[i].Name) {
						fmt.Println("Error -> El nombre no puede repetise dentro de las particiones")
						return "Error -> El nombre no puede repetise dentro de las particiones"
					}
				}
			} else if CToGoString(mbrleido.Part3.Type) == "E" {
				tamanoParticion, err := strconv.Atoi(CToGoString(mbrleido.Part3.Size))
				if err != nil {
					msg_error(err)
				}
				for i := 0; i < len(mbrleido.Part3.Particiones); i++ {
					if CToGoString(mbrleido.Part3.Particiones[i].Status) != "V" {
						if i > 0 {
							ultimo := i - 1
							ultimostart, err := strconv.Atoi(CToGoString(mbrleido.Part3.Particiones[ultimo].Start))
							if err != nil {
								msg_error(err)
							}
							ultimotamano, err := strconv.Atoi(CToGoString(mbrleido.Part3.Particiones[ultimo].Size))
							if err != nil {
								msg_error(err)
							}

							if ultimostart+ultimotamano+tamano <= tamanoParticion {
								copy(nuevoebr.Start[:], strconv.Itoa(ultimostart+ultimotamano+1))
								copy(mbrleido.Part3.Particiones[ultimo].Next[:], strconv.Itoa(ultimostart+ultimotamano+1))
								mbrleido.Part3.Particiones[i] = nuevoebr
								fmt.Println("Partición Logica Creada")
							} else {
								fmt.Println("Error: No hay almacenamiento para nueva partición lógica")
								return "Error -> No hay almacenamiento para nueva partición lógica"
							}

						} else {
							copy(nuevoebr.Start[:], "0")
							if tamano <= tamanoParticion {
								mbrleido.Part3.Particiones[i] = nuevoebr
								fmt.Println("Partición Logica Creada")
							} else {
								fmt.Println("Error: No hay almacenamiento para nueva partición lógica")
								return "Error -> No hay almacenamiento para nueva partición lógica"
							}
						}
						break
					} else if name == CToGoString(mbrleido.Part3.Particiones[i].Name) {
						fmt.Println("Error -> El nombre no puede repetise dentro de las particiones")
						return "Error -> El nombre no puede repetise dentro de las particiones"
					}
				}
			} else if CToGoString(mbrleido.Part4.Type) == "E" {
				tamanoParticion, err := strconv.Atoi(CToGoString(mbrleido.Part4.Size))
				if err != nil {
					msg_error(err)
				}
				for i := 0; i < len(mbrleido.Part4.Particiones); i++ {
					if CToGoString(mbrleido.Part4.Particiones[i].Status) != "V" {
						if i > 0 {
							ultimo := i - 1
							ultimostart, err := strconv.Atoi(CToGoString(mbrleido.Part4.Particiones[ultimo].Start))
							if err != nil {
								msg_error(err)
							}
							ultimotamano, err := strconv.Atoi(CToGoString(mbrleido.Part4.Particiones[ultimo].Size))
							if err != nil {
								msg_error(err)
							}

							if ultimostart+ultimotamano+tamano <= tamanoParticion {
								copy(nuevoebr.Start[:], strconv.Itoa(ultimostart+ultimotamano+1))
								copy(mbrleido.Part4.Particiones[ultimo].Next[:], strconv.Itoa(ultimostart+ultimotamano+1))
								mbrleido.Part4.Particiones[i] = nuevoebr
								fmt.Println("Partición Logica Creada")
							} else {
								fmt.Println("Error: No hay almacenamiento para nueva partición lógica")
								return "Error-> No hay almacenamiento para nueva partición lógica"
							}

						} else {
							copy(nuevoebr.Start[:], "0")
							if tamano <= tamanoParticion {
								mbrleido.Part4.Particiones[i] = nuevoebr
								fmt.Println("Partición Logica Creada")
							} else {
								fmt.Println("Error -> No hay almacenamiento para nueva partición lógica")
								return "Error -> No hay almacenamiento para nueva partición lógica"
							}
						}
						break
					} else if name == CToGoString(mbrleido.Part4.Particiones[i].Name) {
						fmt.Println("Error -> El nombre no puede repetise dentro de las particiones")
						return "Error -> El nombre no puede repetise dentro de las particiones"
					}
				}
			} else {
				fmt.Print("Error: No existe partición extendida.")
				return "Error -> No existe partición extendida."
			}

		} else {
			fmt.Println("Error: Tipo de particion no existe -> ")
			tt := CToGoString(nuevaP.Type)
			fmt.Println(tt)
			return "Error -> Tipo de particion no existe -> " + tt
		}
		escribirMBR(mbrleido, path)
		if tipo == " " {
			return "> Partición p creada."
		}
		return "> Partición " + tipo + " creada."
	} else {
		return msg_parametrosObligatorios()
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

	//fmt.Print("Fecha: ")
	//fmt.Println(string(mbrleido.Fecha[:]))
	//fmt.Print("Firma: ")
	//fmt.Println(string(mbrleido.Signature[:]))
	//fmt.Print("Tamaño: ")
	//fmt.Println(string(mbrleido.Tamano[:]))
	//fmt.Print("Fit: ")
	//fmt.Println(string(mbrleido.Fit[:]))

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
func eliminar_disco(commandArray []string) string {
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
			return "<Error> eliminando archivo:"
		} else {
			fmt.Println("Eliminado correctamente")
			return "> Disco eliminado correctamente."
		}
	} else {
		return msg_parametrosObligatorios()
	}
}

// crear_disco -tamaño=numero -dimensional=dimension/"dimension"
func crear_disco(commandArray []string) string {
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
			return "Error: Dimensional No Reconocida."
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
			return "Error: Fit No Reconocido."
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
		return "> Se creo el disco correctamente."
	} else {
		return msg_parametrosObligatorios()
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

func msg_parametrosObligatorios() string {
	fmt.Println("Error: Parametros Obligatorios No Definidos.")
	return "Error: Parametros Obligatorios No Definidos."
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

//Byte[] a string puro
func CToGoString(c [100]byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}
