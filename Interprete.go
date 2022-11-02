package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	/*locales*/
	"proyecto2/inicio/Disks"
	//"./FormatP"

	//
	"log"
)

func main() {
	Interprete()
}

func lexico(comando string) []string {

	var arr_comandos []string
	comando = comando + " "
	estado := 0
	var lexema string
	for i := 0; i < len(comando); i++ {

		//lexema := ""
		c := comando[i]
		//fmt.Println(i,"-", c)
		//////fmt.Printf("%c\n", c)

		switch estado {
		case 0:

			if c == '\r' {
				//fmt.Println("es retorno de carro")
				estado = 0

			} else if c == ' ' {
				estado = 0

			} else if c == '"' {

				estado = 2
				lexema += string(c)

				/*si es comentario*/
			} else if c == '#' {

				estado = 3
				//lexema += string(c)

			} else if c != ' ' {
				estado = 1
				//lexema += string(c)
				//fmt.Printf("%c\n", c)
				//fmt.Println(string(c))
				lexema += string(c)
				//lexema = append(lexema, string(c))

				//lexema = append(salexema, encodeToUtf8(c))
			}

		case 1:

			if c == '"' {
				arr_comandos = append(arr_comandos, lexema)
				estado = 0
				lexema = ""
				i = i - 1

				/*aqui sigue si es retorno de carro*/
			} else if c == '\r' {
				arr_comandos = append(arr_comandos, lexema)
				estado = 0
				lexema = ""
				i = i - 1

			} else if c != ' ' {
				estado = 1
				//fmt.Println(string(chars))
				//lexema = append(lexema, string(c))
				lexema += string(c)
				//fmt.Println(lexema)

				//lexema = append(lexema, encodeToUtf8(c))

			} else {

				//fmt.Println("lexema: ", lexema, " es: ", estado)
				arr_comandos = append(arr_comandos, lexema)
				estado = 0
				lexema = ""
			}

		case 2:

			if c != '"' {
				lexema += string(c)
				estado = 2

			} else {
				lexema += string(c)
				arr_comandos = append(arr_comandos, lexema)
				estado = 0
				lexema = ""

			}
		/*comentario*/
		case 3:

			/*default:
			fmt.Println("No soportado para exportacion")*/
			//fmt.Println(os)
		}

	}

	return arr_comandos

}

func Interprete_back() {

	fmt.Println("Ingrese comando")

	reader := bufio.NewReader(os.Stdin)
	lin_coman, _ := reader.ReadString('\n')
	comando := strings.ReplaceAll(lin_coman, "\n", "")
	fmt.Println(comando)

	//ReadComando(comando)
	arr_comandos := lexico(comando)

	for i := 0; i < len(arr_comandos); i++ {
		fmt.Print(i, ") ", arr_comandos[i], "-\n")
		//resultado++
	}
	fmt.Println("fin")

	/*leyendo lines de comandos
	no tomar en cuenta si es comentario, o vacio*/
	//fmt.Println("len comandos", len(arr_comandos))
	if len(arr_comandos) > 0 {
		ejecutandoLineComandos(arr_comandos)
	}

}

func Interprete() {

	for true {
		fmt.Println("--- MIA PROYECTO 2---")
		fmt.Println("--- Daniel Sotz A.---")
		fmt.Println("---   201430496   ---")
		fmt.Println("---Ingrese comando---")
		reader := bufio.NewReader(os.Stdin)
		lin_coman, _ := reader.ReadString('\n')
		comando := strings.ReplaceAll(lin_coman, "\n", "")
		//////////fmt.Println("read:",comando)

		var comando_all string

		/*******inicia si es multilinea****/
		com_multi := strings.ReplaceAll(comando, " ", "")
		com_multi = strings.ReplaceAll(com_multi, "\r", "")

		cant_mul := strings.Count(com_multi, "\\*")
		is_multilina := strings.HasSuffix(com_multi, "\\*")
		/////fmt.Println(cant_mul, is_multilina)

		if cant_mul == 1 && is_multilina == true {
			//fmt.Println("es multilinea")

			comando = strings.ReplaceAll(comando, "\r", "")
			comando = strings.ReplaceAll(comando, "\\*", "")
			comando_all = comando

		}

		for cant_mul == 1 && is_multilina == true {
			fmt.Println("2---Ingrese comando---")
			readerm := bufio.NewReader(os.Stdin)
			lin_coman_m, _ := readerm.ReadString('\n')
			comando_t := strings.ReplaceAll(lin_coman_m, "\n", "")
			//fmt.Println("2 read:",comando_t)

			/*ya se declaro arriba*/
			com_multi = strings.ReplaceAll(comando_t, " ", "")
			com_multi = strings.ReplaceAll(com_multi, "\r", "")

			cant_mul = strings.Count(com_multi, "\\*")
			is_multilina = strings.HasSuffix(com_multi, "\\*")
			///fmt.Println("--", cant_mul, is_multilina)

			if cant_mul == 1 && is_multilina == true {
				//fmt.Println("es multilinea")
				comando_t = strings.ReplaceAll(comando_t, "\r", "")
				comando_t = strings.ReplaceAll(comando_t, "\\*", "")
			}

			comando_all = comando_all + " " + comando_t
			//fmt.Println("****comando_all***:",comando_all)
			comando = comando_all
		}

		//fmt.Println("**fin fin fin fin fin***")
		//fmt.Println(comando)

		//return
		/*******finaliza si es multilinea****/

		//ReadComando(comando)
		arr_comandos := lexico(comando)

		for i := 0; i < len(arr_comandos); i++ {
			///fmt.Print(i,") ", arr_comandos[i], "-\n")
			//resultado++
		}
		/////fmt.Println("fin")

		/*leyendo lines de comandos
		no tomar en cuenta si es comentario, o vacio*/
		//fmt.Println("len comandos", len(arr_comandos))
		if len(arr_comandos) > 0 {
			ejecutandoLineComandos(arr_comandos)
		}

	}

}

func ejecutandoLineComandos(lin_coman []string) {
	com_start := strings.ToLower(lin_coman[0])

	if com_start == "exec" {
		fmt.Println("ejecutar exec archivo")
		Exec(lin_coman)

	} else if com_start == "pause" {
		Pausa(lin_coman)
	} else if com_start == "mkdisk" {

		fmt.Println("para crear disco")
		MKdisk(lin_coman)

	} else if com_start == "rmdisk" {

		fmt.Println("Elimina disco")
		RKdisk(lin_coman)

	} else if com_start == "mount" {

		//fmt.Println("len(lin_coman)", len(lin_coman))
		if len(lin_coman) == 1 {
			fmt.Println("Listado de particiones")
			Mount_list()
		} else {
			fmt.Println("montar particion de disco")
			Mount(lin_coman)
		}

	} else if com_start == "unmount" {

		fmt.Println("quitar una particion montada")
		UnMount(lin_coman)

	} else if com_start == "fdisk" {

		fmt.Println("creando las particiones")
		Fdisk(lin_coman)

		/*comando para reportes*/
	} else if com_start == "rep" {

		fmt.Println("imprime mbr")
		RepGraf(lin_coman)

		/////////************cominza para sistema de archivos******************///
	} else if com_start == "mkfs" {

		fmt.Println("formateo mbr")
		MKFS(lin_coman)

		/*crear archivo*/
	} else if com_start == "mkfile" {

		fmt.Println("creando archivo")
		MKFILE(lin_coman)

		/*crear carpeta*/
	} else if com_start == "mkdir" {

		fmt.Println("creando carpeta")
		MKDIR(lin_coman)

		/*leer archivos*/
	} else if com_start == "cat" {

		fmt.Println("leyendo archivo")
		CAT(lin_coman)

		/*cambiar nombre*/
	} else if com_start == "ren" {

		fmt.Println("cambiar nombre")
		REM(lin_coman)

		/*edita un archivo*/
	} else if com_start == "edit" {

		fmt.Println("editando archivo")
		EDIT(lin_coman)

		/*} else if com_start == "format" {
		fmt.Println("formateo mbr")
		//FormatP.Format()*/

		/*para salir de consola*/
		//} else if com_start == "salir" {

	} else {
		fmt.Println("Error, comando no reconocido")
	}
}

// //Renombrar
func REM(comandos []string) {

	var path string
	var name string
	var id string

	en_path := false
	en_name := false
	en_id := false

	/*si todo esta correcto y yo que me alegro*/
	all_com_correc := true

	var des_err string = ""

	path = ""
	for i := 1; i < len(comandos); i++ {
		//fmt.Print(i,"-REM(", comandos[i], ")\n")

		var_comand := strings.Split(comandos[i], "->")
		//fmt.Println("0",  var_comand,":" , len(var_comand))

		if len(var_comand) == 2 {

			command := strings.ToLower(var_comand[0])

			if command == "-path" {

				/*si hay duplicados*/
				if en_path == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					///////////////////////////////
					path = var_comand[1]
					en_path, path, i = VerificarPath(path, i, comandos, var_comand)
					//////////////////////
				}

			} else if command == "-name" {

				/*si hay duplicados*/
				if en_name == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					///////////////////
					name = var_comand[1]
					en_name, name, i = VerificarCont(name, i, comandos, var_comand)
					//////////////////

				}

				////id obligatorio
			} else if command == "-id" {

				/*si hay duplicados*/
				if en_id == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					////////////////
					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 { ///new
						all_com_correc = false
						des_err = des_err + " -id "
					} else {
						//fmt.Println("var_comand[1]", var_comand[1])
						id = var_comand[1]
						//fmt.Println("unit", string(unit))
						en_id = true
						//fmt.Print("     id encontrado-", name, "-\n" )
					}
					////////////

				}

			} else {
				all_com_correc = false
				des_err = des_err + "(" + command + " No reconocido)"
				fmt.Println(command, "commando incorrecto*")
			}

		} else {
			des_err = des_err + "(" + var_comand[0] + " No reconocido)"
			all_com_correc = false
			fmt.Println("Error en REM, SOBRAN parametros o comando no reconocido")
		}
	}

	if en_path == false {
		des_err = des_err + " -path "
	}
	if en_name == false {
		des_err = des_err + " -name "
	}
	if en_id == false {
		des_err = des_err + " -id "
	}

	if en_path == false || en_name == false || en_id == false || all_com_correc == false {
		fmt.Println("Error en REM, en", des_err)
	} else {
		//fmt.Println("comando valido, montar")
		fmt.Println("path", path, "nombre", name, "id", id)
		Disks.Cambiar_Name(path, name, id)

	}

}

func CAT(comandos []string) {

	var id string
	var file string
	/*en_path := false*/

	en_id := false
	en_file := false

	var arr_file_read []string

	/*si todo esta correcto y yo que me alegro*/
	all_com_correc := true
	var des_err string = ""

	for i := 1; i < len(comandos); i++ {
		//fmt.Print(i,"-CAT(", comandos[i], ")\n")

		var_comand := strings.Split(comandos[i], "->")
		//fmt.Println("0",  var_comand,":" , len(var_comand))

		if len(var_comand) == 2 {

			command := strings.ToLower(var_comand[0])

			//if command == "-id" {
			if len(command) >= 5 && command[0:5] == "-file" {

				////////////////
				/*verificando si tiene valor del parametro*/
				/*if (len(var_comand[1]) == 0) {
					all_com_correc = false
					des_err =  des_err + " -file "
				} else {*/

				///////////////////
				file = var_comand[1]
				en_file, file, i = VerificarCont(file, i, comandos, var_comand)
				//////////////////
				arr_file_read = append(arr_file_read, file)

				//}
				////////////

				////id obligatorio
			} else if command == "-id" {

				/*si hay duplicados*/
				if en_id == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					////////////////
					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 { ///new
						all_com_correc = false
						des_err = des_err + " -id "
					} else {
						//fmt.Println("var_comand[1]", var_comand[1])
						id = var_comand[1]
						//fmt.Println("unit", string(unit))
						en_id = true
						//fmt.Print("     id encontrado-", name, "-\n" )
					}
					////////////

				}
			} else {
				all_com_correc = false
				des_err = des_err + "(" + command + " No reconocido)"
				fmt.Println(command, "commando incorrecto*")
			}

		} else {
			des_err = des_err + "(" + var_comand[0] + " No reconocido)"
			all_com_correc = false
			fmt.Println("Error en CAT, SOBRAN parametros o comando no reconocido")
		}
	}

	if en_file == false {
		des_err = des_err + " -file "
	}

	if en_id == false {
		des_err = des_err + " -id "
	}

	//fmt.Println("arr_file_read", arr_file_read)
	//fmt.Println("id", id)
	if en_id == false || en_file == false || all_com_correc == false {
		fmt.Println("Error en CAT, en", des_err)
	} else {
		Disks.Print_Contenido(arr_file_read, id)
	}

}

// creando carpetas
func MKDIR(comandos []string) {
	var path string
	var id string
	var p byte

	en_path := false
	en_id := false
	en_p := false

	/*si todo esta correcto y yo que me alegro*/
	all_com_correc := true

	var des_err string = ""

	path = ""
	for i := 1; i < len(comandos); i++ {
		//fmt.Print(i,"-M(", comandos[i], ")\n")

		var_comand := strings.Split(comandos[i], "=")
		//fmt.Println("0",  var_comand,":" , len(var_comand))

		if len(var_comand) == 2 {

			command := strings.ToLower(var_comand[0])
			/**OBLIGATORIO**/
			if command == "-path" {

				/*si hay duplicados*/
				if en_path == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					///////////////////
					path = var_comand[1]
					en_path, path, i = VerificarPath(path, i, comandos, var_comand)
					//////////////////

				}

				////id obligatorio
			} else if command == "-id" {

				/*si hay duplicados*/
				if en_id == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					////////////////
					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 { ///new
						all_com_correc = false
						des_err = des_err + " -id "
					} else {
						//fmt.Println("var_comand[1]", var_comand[1])
						id = var_comand[1]
						//fmt.Println("unit", string(unit))
						en_id = true
						//fmt.Print("     id encontrado-", name, "-\n" )
					}
					////////////

				}

			} else {
				all_com_correc = false
				des_err = des_err + "(" + command + " No reconocido)"
				//fmt.Println(command, "commando incorrecto*")
			}

		} else if len(var_comand) == 1 {

			command := strings.ToLower(var_comand[0])
			if command == "-p" {

				/*si hay duplicados*/
				if en_p == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					////////////////
					//fmt.Println("var_comand[0]", var_comand[0])
					p = 'S'
					en_p = true

				}
			} else {
				all_com_correc = false
				des_err = des_err + "(" + command + " No reconocido)"
				//fmt.Println(command, "commando incorrecto*")

			}

		} else {
			des_err = des_err + "(" + var_comand[0] + " No reconocido)"
			all_com_correc = false
			fmt.Println("Error en MKDIR, SOBRAN parametros o comando no reconocido")
		}
	}

	if en_p == false {
		p = 'N'
	}

	if en_path == false {
		des_err = des_err + " -path "
	}

	if en_id == false {
		des_err = des_err + " -id "
	}

	if en_id == false || en_path == false || all_com_correc == false {
		fmt.Println("Error en MKDIR, en", des_err)
	} else {
		fmt.Println("path", path, "id", id, "p", string(p))
		Disks.CrearCarpeta(path, id, p)
	}
}

// /editar archivo
func EDIT(comandos []string) {
	var size int64
	var path string
	var id string
	var cont string

	en_size := false
	en_path := false
	en_id := false
	en_cont := false

	/*si todo esta correcto y yo que me alegro*/
	all_com_correc := true

	var des_err string = ""

	path = ""
	for i := 1; i < len(comandos); i++ {
		//fmt.Print(i,"-E(", comandos[i], ")\n")

		var_comand := strings.Split(comandos[i], "->")
		//fmt.Println("0",  var_comand,":" , len(var_comand))

		if len(var_comand) == 2 {

			command := strings.ToLower(var_comand[0])
			/**OBLIGATORIO**/
			if command == "-size" {

				/*si hay duplicados*/
				if en_size == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {
					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 { ///new
					} else {
						ss, _ := strconv.ParseInt(var_comand[1], 10, 64)
						size = ss

						if size > 0 {
							en_size = true
							//fmt.Println("        size encontrado ", size)
						} else {
							des_err = des_err + "(size, debe de ser > 0 ) "
						}
					}
				}

			} else if command == "-path" {

				/*si hay duplicados*/
				if en_path == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					///////////////////
					path = var_comand[1]
					en_path, path, i = VerificarPath(path, i, comandos, var_comand)
					//////////////////

				}

			} else if command == "-cont" {

				/*si hay duplicados*/
				if en_cont == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					///////////////////
					cont = var_comand[1]
					en_cont, cont, i = VerificarCont(cont, i, comandos, var_comand)
					//////////////////

				}

				////id obligatorio
			} else if command == "-id" {

				/*si hay duplicados*/
				if en_id == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					////////////////
					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 { ///new
						all_com_correc = false
						des_err = des_err + " -id "
					} else {
						//fmt.Println("var_comand[1]", var_comand[1])
						id = var_comand[1]
						//fmt.Println("unit", string(unit))
						en_id = true
						//fmt.Print("     id encontrado-", name, "-\n" )
					}
					////////////

				}

			} else {
				all_com_correc = false
				des_err = des_err + "(" + command + " No reconocido)"
				fmt.Println(command, "commando incorrecto*")
			}

		} else {
			des_err = des_err + "(" + var_comand[0] + " No reconocido)"
			all_com_correc = false
			fmt.Println("Error en EDIT, SOBRAN parametros o comando no reconocido")
		}
	}

	if en_size == false {

		size = 0
	}
	if en_cont == false {

		cont = ""
	}
	if en_path == false {
		des_err = des_err + " -path "
	}

	if en_id == false {
		des_err = des_err + " -id "
	}

	if en_id == false || en_path == false || all_com_correc == false {
		fmt.Println("Error en EDIT, en", des_err)
	} else {
		//fmt.Println("comando valido")
		fmt.Println("size", size, "path", path, "id", id, "cont", cont)
		Disks.EditFile(size, path, id, cont)
	}
}

func MKFILE(comandos []string) {
	var size int64
	var path string
	var id string
	var p byte
	var cont string

	en_size := false
	en_path := false
	en_id := false
	en_p := false
	en_cont := false

	/*si todo esta correcto y yo que me alegro*/
	all_com_correc := true

	var des_err string = ""

	path = ""
	for i := 1; i < len(comandos); i++ {
		//fmt.Print(i,"-M(", comandos[i], ")\n")

		var_comand := strings.Split(comandos[i], "=")
		//fmt.Println("0",  var_comand,":" , len(var_comand))

		if len(var_comand) == 2 {

			command := strings.ToLower(var_comand[0])
			/**OBLIGATORIO**/
			if command == "-size" {

				/*si hay duplicados*/
				if en_size == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {
					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 { ///new
					} else {
						ss, _ := strconv.ParseInt(var_comand[1], 10, 64)
						size = ss

						if size > 0 {
							en_size = true
							//fmt.Println("        size encontrado ", size)
						} else {
							des_err = des_err + "(size, debe de ser > 0 ) "
						}
					}
				}

			} else if command == "-path" {

				/*si hay duplicados*/
				if en_path == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					///////////////////
					path = var_comand[1]
					en_path, path, i = VerificarPath(path, i, comandos, var_comand)
					//////////////////

				}

			} else if command == "-cont" {

				/*si hay duplicados*/
				if en_cont == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					///////////////////
					cont = var_comand[1]
					en_cont, cont, i = VerificarCont(cont, i, comandos, var_comand)
					//////////////////

				}

				////id obligatorio
			} else if command == "-id" {

				/*si hay duplicados*/
				if en_id == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					////////////////
					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 { ///new
						all_com_correc = false
						des_err = des_err + " -id "
					} else {
						//fmt.Println("var_comand[1]", var_comand[1])
						id = var_comand[1]
						//fmt.Println("unit", string(unit))
						en_id = true
						//fmt.Print("     id encontrado-", name, "-\n" )
					}
					////////////

				}

			} else {
				all_com_correc = false
				des_err = des_err + "(" + command + " No reconocido)"
				fmt.Println(command, "commando incorrecto*")
			}

		} else if len(var_comand) == 1 {

			command := strings.ToLower(var_comand[0])
			if command == "-r" {

				/*si hay duplicados*/
				if en_p == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					////////////////
					//fmt.Println("var_comand[0]", var_comand[0])
					p = 'S'
					en_p = true

				}
			} else {
				all_com_correc = false
				des_err = des_err + "(" + command + " No reconocido)"
				fmt.Println(command, "commando incorrecto*")

			}

		} else {
			des_err = des_err + "(" + var_comand[0] + " No reconocido)"
			all_com_correc = false
			fmt.Println("Error en MKFILE, SOBRAN parametros o comando no reconocido")
		}
	}

	if en_p == false {
		p = 'N'
	}
	if en_size == false {

		size = 0
	}
	if en_cont == false {

		cont = ""
	}
	if en_path == false {
		des_err = des_err + " -path "
	}

	if en_id == false {
		des_err = des_err + " -id "
	}

	if en_path == false || all_com_correc == false {
		fmt.Println("Error en MKFILE, en", des_err)
	} else {
		//fmt.Println("comando valido")
		fmt.Println("size", size, "path", path, "id", id, "p", string(p), "cont", cont)
		Disks.CrearFile(size, path, id, p, cont)
	}
}

func MKFS(comandos []string) {

	var type_s string
	var unit byte
	var id string
	var add int64

	en_type := false
	en_unit := false
	en_id := false
	en_add := false

	/*si todo esta correcto y yo que me alegro*/
	all_com_correc := true

	var des_err string = ""

	for i := 1; i < len(comandos); i++ {
		//fmt.Print(i,"-MKFS(", comandos[i], ")\n")

		var_comand := strings.Split(comandos[i], "=")
		//fmt.Println("0",  var_comand,":" , len(var_comand))

		if len(var_comand) == 2 {

			command := strings.ToLower(var_comand[0])

			/*OPCIONAL*/
			if command == "-type" {

				/*si hay duplicados*/
				if en_type == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					////////////////
					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 {
						all_com_correc = false
						des_err = des_err + " -type "
					} else {
						type_s = strings.ToLower(var_comand[1])

						if type_s == "fast" || type_s == "full" {
							en_type = true
						} else {
							all_com_correc = false
							des_err = des_err + "( type debe ser fast ó full) "
						}
					}
					////////////

				}

				/*OPCIONAL*/
			} else if command == "-unit" {

				/*si hay duplicados*/
				if en_unit == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 || len(var_comand[1]) > 1 { ///new
						all_com_correc = false
						des_err = des_err + " -unit "
					} else {
						//fmt.Println("var_comand[1]", var_comand[1])
						//unit =   var_comand[1][0]
						unidad_s := strings.ToLower(var_comand[1])
						unit = unidad_s[0]

						//fmt.Println("unit", string(unit))
						if unit == 'b' || unit == 'k' || unit == 'm' {
							en_unit = true
							/////////fmt.Println("        unit encontrado ", unit)
						} else {
							all_com_correc = false
							des_err = des_err + "( unit debe ser B, K ó M) "
						}
					}

				}

			} else if command == "-add" {

				/*si hay duplicados*/
				if en_add == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {
					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 { ///new
					} else {

						ss, _ := strconv.ParseInt(var_comand[1], 10, 64)
						//ss, _ := strconv.Parseint(var_comand[1], 0, 64)
						add = ss
						en_add = true
						fmt.Println("        add encontrado ", add)

					}
				}

				////id obligatorio
			} else if command == "-id" {

				/*si hay duplicados*/
				if en_id == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					////////////////
					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 { ///new
						all_com_correc = false
						des_err = des_err + " -id "
					} else {
						//fmt.Println("var_comand[1]", var_comand[1])
						id = var_comand[1]
						//fmt.Println("unit", string(unit))
						en_id = true
						//fmt.Print("     id encontrado-", name, "-\n" )
					}
					////////////

				}

			} else {
				all_com_correc = false
				des_err = des_err + "(" + command + " No reconocido)"
				fmt.Println(command, "commando incorrecto*")
			}

		} else {
			des_err = des_err + "(" + var_comand[0] + " No reconocido)"
			all_com_correc = false
			fmt.Println("Error en REP, SOBRAN parametros o comando no reconocido")
		}
	}

	if en_type == false {
		type_s = "full"
	}
	if en_unit == false {
		unit = 'k'
	}
	if en_id == false {
		des_err = des_err + " -id "
	}

	//if (en_path == false || en_name == false || en_id == false || all_com_correc == false) {
	if en_id == false || all_com_correc == false {
		fmt.Println("Error en MKFS, en", des_err)
	} else {

		//fmt.Println("comando valido, formatear")
		//fmt.Println("id", id, "type_s", type_s, "unit", string(unit), "add", add)
		Disks.FormatPart(type_s, id)

	}

}

func RepGraf(comandos []string) {

	var path string
	var name string
	var id string
	var ruta string

	en_path := false
	en_name := false
	en_id := false
	en_ruta := false

	/*si todo esta correcto y yo que me alegro*/
	all_com_correc := true

	var des_err string = ""

	path = ""
	for i := 1; i < len(comandos); i++ {
		//fmt.Print(i,"-GRA(", comandos[i], ")\n")

		var_comand := strings.Split(comandos[i], "=")
		//fmt.Println("0",  var_comand,":" , len(var_comand))

		if len(var_comand) == 2 {

			command := strings.ToLower(var_comand[0])

			if command == "-path" {

				/*si hay duplicados*/
				if en_path == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					///////////////////////////////
					path = var_comand[1]
					en_path, path, i = VerificarPath(path, i, comandos, var_comand)
					//////////////////////
				}

			} else if command == "-name" {

				/*si hay duplicados*/
				if en_name == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					////////////////
					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 { ///new
						all_com_correc = false
						des_err = des_err + " -name "
					} else {
						//fmt.Println("var_comand[1]", var_comand[1])
						name = var_comand[1]
						//fmt.Println("unit", string(unit))
						en_name = true
						fmt.Print("     name encontrado-", name, "-\n")
					}
					////////////

				}

			} else if command == "-ruta" {

				/*si hay duplicados*/
				if en_ruta == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					///////////////////
					ruta = var_comand[1]
					en_ruta, ruta, i = VerificarCont(ruta, i, comandos, var_comand)
					//////////////////

				}

				////id obligatorio
			} else if command == "-id" {

				/*si hay duplicados*/
				if en_id == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					////////////////
					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 { ///new
						all_com_correc = false
						des_err = des_err + " -id "
					} else {
						//fmt.Println("var_comand[1]", var_comand[1])
						id = var_comand[1]
						//fmt.Println("unit", string(unit))
						en_id = true
						//fmt.Print("     id encontrado-", name, "-\n" )
					}
					////////////

				}

			} else {
				all_com_correc = false
				des_err = des_err + "(" + command + " No reconocido)"
				fmt.Println(command, "commando incorrecto*")
			}

		} else {
			des_err = des_err + "(" + var_comand[0] + " No reconocido)"
			all_com_correc = false
			fmt.Println("Error en REP, SOBRAN parametros o comando no reconocido")
		}
	}

	if en_path == false {
		des_err = des_err + " -path "
	}
	if en_name == false {
		des_err = des_err + " -name "
	}
	if en_id == false {
		des_err = des_err + " -id "
	}

	if en_path == false || en_name == false || en_id == false || all_com_correc == false {
		fmt.Println("Error en REP, en", des_err)
	} else {

		/*creando disco con los parametros*/
		//fmt.Println("path", path, "nombre", name, "id", id, "ruta", ruta)
		Disks.ReportesGraf(path, name, id, ruta)

	}

}

func Pausa(comandos []string) {
	//fmt.Println("len(comandos)", len(comandos))
	if len(comandos) == 1 {

		fmt.Println("Presione tecla ENTER para Continuar...")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')

	} else {
		fmt.Println("Error en pause, parametros incorrectos")

	}
}

func UnMount(comandos []string) {

	/*var path string
	var name string
	en_path := false
	en_name := false*/
	var arr_desmont []string

	/*si todo esta correcto y yo que me alegro*/
	all_com_correc := true
	var des_err string = ""

	for i := 1; i < len(comandos); i++ {
		//fmt.Print(i,"-UNMO(", comandos[i], ")\n")

		var_comand := strings.Split(comandos[i], "=")
		//fmt.Println("0",  var_comand,":" , len(var_comand))

		if len(var_comand) == 2 {

			command := strings.ToLower(var_comand[0])

			//if command == "-id" {
			if len(command) >= 3 && command[0:3] == "-id" {

				////////////////
				/*verificando si tiene valor del parametro*/
				if len(var_comand[1]) == 0 { ///new
					all_com_correc = false
					des_err = des_err + " -id "
				} else {
					arr_desmont = append(arr_desmont, var_comand[1])
					//fmt.Println("var_comand[1]", var_comand[1])
					//name = var_comand[1]
					//fmt.Print("     name encontrado-", name, "-\n" )
				}
				////////////

			} else {
				all_com_correc = false
				des_err = des_err + "(" + command + " No reconocido)"
				fmt.Println(command, "commando incorrecto*")
			}

		} else {
			des_err = des_err + "(" + var_comand[0] + " No reconocido)"
			all_com_correc = false
			fmt.Println("Error en UNMOUNT, SOBRAN parametros o comando no reconocido")
		}
	}

	fmt.Println("arr_desmont", arr_desmont)
	if all_com_correc == false {
		fmt.Println("Error en UNMOUNT, en", des_err)
	} else {
		Disks.DesMontar(arr_desmont)
	}

}

/*ini mount particion*/
func Mount(comandos []string) {
	var path string
	var name string

	en_path := false
	en_name := false

	/*si todo esta correcto y yo que me alegro*/
	all_com_correc := true

	var des_err string = ""

	path = ""
	for i := 1; i < len(comandos); i++ {
		//fmt.Print(i,"-MO(", comandos[i], ")\n")

		var_comand := strings.Split(comandos[i], "=")
		//fmt.Println("0",  var_comand,":" , len(var_comand))

		if len(var_comand) == 2 {

			command := strings.ToLower(var_comand[0])

			if command == "-path" {

				/*si hay duplicados*/
				if en_path == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					///////////////////////////////
					path = var_comand[1]
					en_path, path, i = VerificarPath(path, i, comandos, var_comand)
					//////////////////////
				}

			} else if command == "-name" {

				/*si hay duplicados*/
				if en_name == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					////////////////
					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 { ///new
						all_com_correc = false
						des_err = des_err + " -name "
					} else {
						//fmt.Println("var_comand[1]", var_comand[1])
						name = var_comand[1]
						//fmt.Println("unit", string(unit))
						en_name = true
						fmt.Print("     name encontrado-", name, "-\n")
					}
					////////////

				}

			} else {
				all_com_correc = false
				des_err = des_err + "(" + command + " No reconocido)"
				fmt.Println(command, "commando incorrecto*")
			}

		} else {
			des_err = des_err + "(" + var_comand[0] + " No reconocido)"
			all_com_correc = false
			fmt.Println("Error en MOUNT, SOBRAN parametros o comando no reconocido")
		}
	}

	if en_path == false {
		des_err = des_err + " -path "
	}
	if en_name == false {
		des_err = des_err + " -name "
	}

	if en_path == false || en_name == false || all_com_correc == false {
		fmt.Println("Error en MOUNT, en", des_err)
	} else {
		//fmt.Println("comando valido, montar")
		/*creando disco con los parametros*/
		fmt.Println("path", path, "name", name)
		Disks.Mount_Par(path, name)

	}
}

/*fin mount particion*/

func Mount_list() {
	Disks.Mount_list()
}

/*eliminado disco*/
func RKdisk(comandos []string) {
	var path string
	//var name string

	en_path := false
	//en_name := false

	/*si todo esta correcto y yo que me alegro*/
	all_com_correc := true

	var des_err string = ""

	path = ""
	for i := 1; i < len(comandos); i++ {
		//fmt.Print(i,"-R(", comandos[i], ")\n")

		var_comand := strings.Split(comandos[i], "=")
		//fmt.Println("0",  var_comand,":" , len(var_comand))

		if len(var_comand) == 2 {

			command := strings.ToLower(var_comand[0])

			if command == "-path" {

				/*si hay duplicados*/
				if en_path == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					////////////////////
					path = var_comand[1]
					en_path, path, i = VerificarPath(path, i, comandos, var_comand)
					////////////
				}
			} else {
				all_com_correc = false
				des_err = des_err + "(" + command + " No reconocido)"
				fmt.Println(command, "commando incorrecto*")
			}

		} else {
			des_err = des_err + "(" + var_comand[0] + " No reconocido)"
			all_com_correc = false
			fmt.Println("Error en EXEC, SOBRAN parametros o comando no reconocido")
		}
	}

	if en_path == false {
		des_err = des_err + " -path "
	} /*else {
		////*verificando la extension
		name := strings.Split(path, "/")
		name_imp := name[len(name)-1]
		//fmt.Println("2 name_imp",  name_imp )

		name_ext := strings.Split(name_imp, ".")
		///*****************verificar despues si hay nombre antes de la extencin
		if ( len(name_ext) == 2 && name_ext[1] == "mia") {
			////fmt.Println("******ext correcto")
		} else {
			all_com_correc = false
			des_err =  des_err +  " extension invalida "
			///fmt.Println("******ext NO CORECTO, o nombre incorrecto")
		}
	}*/

	if en_path == false || all_com_correc == false {
		fmt.Println("Error en RMDISK, en", des_err)
	} else {
		fmt.Println("comando valido")
		/*creando disco con los parametros*/

		/*Eliminando Disco*/
		Disks.DeleteDisk(path)

		//////
	}

}

func VerificarPath(path string, i int, comandos []string, var_comand []string) (bool, string, int) {

	////////////////////
	en_path := false

	/*verificando path*/
	if len(path) == 0 {
		//fmt.Println("len(path)", len(path), "i", i+1, len(comandos) )
		/*verificando si el que sigue es el path*/
		if i+1 < len(comandos) {
			if comandos[i+1][0] == '"' {
				i = i + 1

				path = comandos[i]
				path = strings.ReplaceAll(path, "\"", "")
				en_path = true
				//fmt.Println("*       path encontrado -", path,"-", len(var_comand))
			}
		}

	} else {
		path = var_comand[1]
		en_path = true
		//fmt.Println("        path encontrado -", path,"-", len(var_comand))
	}

	return en_path, path, i
}

func VerificarCont(cont string, i int, comandos []string, var_comand []string) (bool, string, int) {

	////////////////////
	en_cont := false

	/*verificando path*/
	if len(cont) == 0 {
		//fmt.Println("len(path)", len(path), "i", i+1, len(comandos) )
		/*verificando si el que sigue es el path*/
		if i+1 < len(comandos) {
			if comandos[i+1][0] == '"' {
				i = i + 1

				cont = comandos[i]
				cont = strings.ReplaceAll(cont, "\"", "")
				en_cont = true
				//fmt.Println("*       path encontrado -", path,"-", len(var_comand))
			}
		}

	} else {
		cont = var_comand[1]
		en_cont = true
		//fmt.Println("        path encontrado -", path,"-", len(var_comand))
	}

	return en_cont, cont, i
}

func VerificarNameCom(name string, i int, comandos []string, var_comand []string) (bool, string, int) {

	////////////////////
	en_name := false

	/*verificando name*/
	if len(name) == 0 {
		//fmt.Println("len(name)", len(name), "i", i+1, len(comandos) )
		/*verificando si el que sigue es el name*/
		if i+1 < len(comandos) {
			if comandos[i+1][0] == '"' {
				i = i + 1

				name = comandos[i]
				name = strings.ReplaceAll(name, "\"", "")
				en_name = true
				//fmt.Println("*       name encontrado -", name,"-", len(var_comand))
			}
		}

	} else {
		name = var_comand[1]
		en_name = true
		//fmt.Print("     name encontrado-", name, "-\n" )
	}

	return en_name, name, i
}

func Exec(comandos []string) {
	var path string
	//var name string

	en_path := false
	//en_name := false

	/*si todo esta correcto y yo que me alegro*/
	all_com_correc := true

	var des_err string = ""

	path = ""
	for i := 1; i < len(comandos); i++ {
		//fmt.Print(i,"-E(", comandos[i], ")\n")

		var_comand := strings.Split(comandos[i], "=")
		//fmt.Println("0",  var_comand,":" , len(var_comand))

		if len(var_comand) == 2 {

			command := strings.ToLower(var_comand[0])

			if command == "-path" {

				/*si hay duplicados*/
				if en_path == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					////////////////////
					path = var_comand[1]

					/*verificando path*/
					if len(path) == 0 {
						//fmt.Println("len(path)", len(path), "i", i+1, len(comandos) )
						/*verificando si el que sigue es el path*/
						if i+1 < len(comandos) {
							if comandos[i+1][0] == '"' {
								i = i + 1

								path = comandos[i]
								path = strings.ReplaceAll(path, "\"", "")
								en_path = true
								//fmt.Println("*       path encontrado -", path,"-", len(var_comand))
							}
						}

					} else {
						path = var_comand[1]
						en_path = true
						//fmt.Println("        path encontrado -", path,"-", len(var_comand))
					}
					///////////////////////////
				}
			} else {
				all_com_correc = false
				des_err = des_err + "(" + command + " No reconocido)"
				fmt.Println(command, "commando incorrecto*")
			}

		} else {
			des_err = des_err + "(" + var_comand[0] + " No reconocido)"
			all_com_correc = false
			fmt.Println("Error en EXEC, SOBRAN parametros o comando no reconocido")
		}
	}

	if en_path == false {
		des_err = des_err + " -path "
	} else {
		/*verificando la extension*/
		name := strings.Split(path, "/")
		name_imp := name[len(name)-1]
		//fmt.Println("2 name_imp",  name_imp )

		name_ext := strings.Split(name_imp, ".")
		/*****************verificar despues si hay nombre antes de la extencin*/
		if len(name_ext) == 2 && name_ext[1] == "script" {
			////fmt.Println("******ext correcto")
		} else {
			all_com_correc = false
			des_err = des_err + " extension invalida "
			///fmt.Println("******ext NO CORECTO, o nombre incorrecto")
		}
	}

	if en_path == false || all_com_correc == false {
		fmt.Println("Error en EXEC, en", des_err)
	} else {
		fmt.Println("comando valido")

		fmt.Println("path", path /*, "name", name,*/)

		/*leyendo archivo .script*/
		leyendo_mia(path)
	}

}
func leyendo_mia(path_mia string) {

	//abriendo archivo
	file, err := os.Open(path_mia)

	if err != nil {
		log.Fatalf("Error abriendo archivo: %s\n", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var comando_all string

	/*leyendo linea por linea*/
	cant_mul := 0
	is_multilina := false
	com_multi := ""

	comando_all = ""
	for scanner.Scan() {

		//fmt.Println(scanner.Bytes())
		comando := scanner.Text()
		fmt.Println("|", comando)

		/*******inicia si es multilinea****/
		com_multi = strings.ReplaceAll(comando, " ", "")
		//com_multi = strings.ReplaceAll(com_multi, "\r", "")

		cant_mul = strings.Count(com_multi, "\\*")
		is_multilina = strings.HasSuffix(com_multi, "\\*")
		/////////fmt.Println(cant_mul, is_multilina)

		if cant_mul == 1 && is_multilina == true {
			////fmt.Println("es multilinea")

			//comando = strings.ReplaceAll(comando, "\r", "")
			comando = strings.ReplaceAll(comando, "\\*", "")
			comando_all = comando_all + comando
			///////fmt.Println("----------------",comando_all)

		} else {
			comando_all = comando_all + comando
			//////fmt.Println("+++++++++++++++++++",comando_all)
			///comando_all = ""

			///////////
			arr_comandos := lexico(comando_all)
			/*for i := 0; i < len(arr_comandos) ; i++ {
				fmt.Print("              ", i,") ", arr_comandos[i], "-\n")
			}*/
			if len(arr_comandos) > 0 {
				ejecutandoLineComandos(arr_comandos)
			}

			comando_all = ""

		}

	}

}

func Fdisk(comandos []string) {
	var size uint64 //
	var path string //******/
	var name string
	var unit byte //

	var btype byte //
	var fit byte   //

	var delete string
	var add int64

	en_size := false //
	en_path := false
	en_name := false
	en_unit := false //

	en_type := false //
	en_fit := false  //

	en_delete := false //

	en_add := false

	/*si todo esta correcto y yo que me alegro*/
	all_com_correc := true

	var des_err string = ""
	var des_err_del string = ""
	var des_err_add string = ""

	path = ""
	for i := 1; i < len(comandos); i++ {
		//fmt.Print(i,"-FF(", comandos[i], ")\n")

		var_comand := strings.Split(comandos[i], "=")
		//fmt.Println("0",  var_comand,":" , len(var_comand))

		if len(var_comand) == 2 {

			command := strings.ToLower(var_comand[0])

			if command == "-size" {

				/*si hay duplicados*/
				if en_size == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {
					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 { ///new
					} else {
						ss, _ := strconv.ParseUint(var_comand[1], 0, 64)
						size = ss

						if size > 0 {
							en_size = true
							//fmt.Println("        size encontrado ", size)
						} else {
							des_err = des_err + "(size, debe de ser > 0 ) "
						}

						/*fmt.Println("", var_comand[0],"")
						fmt.Println("*", len(var_comand[1]),"-")

						for y := 0; y < len(var_comand[1]) ; y++ {
							//fmt.Printf("%v",y ," %c", var_comand[1][y],"-\n")
							fmt.Println("-", var_comand[1][y],"-")

						}*/
					}
				}

			} else if command == "-add" {

				/*si hay duplicados*/
				if en_add == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {
					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 { ///new
					} else {

						ss, _ := strconv.ParseInt(var_comand[1], 10, 64)
						//ss, _ := strconv.Parseint(var_comand[1], 0, 64)
						add = ss
						en_add = true
						//fmt.Println("        add encontrado ", add)

					}
				}

				/*OBLIGA*/
			} else if command == "-path" {

				/*si hay duplicados*/
				if en_path == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					///////////////////
					path = var_comand[1]
					en_path, path, i = VerificarPath(path, i, comandos, var_comand)
					//////////////////
					//////////////////////

				}

				/*OPCIONAL*/
			} else if command == "-type" {

				/*si hay duplicados*/
				if en_type == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					/*verificando si tiene valor del parametro*/
					//fmt.Println("len(var_comand[1])", len(var_comand[1]))
					if len(var_comand[1]) == 0 || len(var_comand[1]) > 1 { ///new
						all_com_correc = false
						des_err = des_err + " -type"
					} else {
						//fmt.Println("var_comand[1]", var_comand[1])
						//unit =   var_comand[1][0]
						type_s := strings.ToLower(var_comand[1])
						btype = type_s[0]

						//fmt.Println("unit", string(unit))
						if btype == 'p' || btype == 'e' || btype == 'l' {
							en_type = true
							//fmt.Println("        type encontrado ", btype)
						} else {
							all_com_correc = false
							des_err = des_err + "( type debe ser P, E ó L) "
						}
					}

				}

				/*OPCIONAL*/
			} else if command == "-fit" {

				/*si hay duplicados*/
				if en_fit == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					/*verificando si tiene valor del parametro*/
					//fmt.Println("len(var_comand[1])", len(var_comand[1]))
					if len(var_comand[1]) == 0 || len(var_comand[1]) > 2 { ///new
						all_com_correc = false
						des_err = des_err + " -fit"
					} else {
						//fmt.Println("var_comand[1]", var_comand[1])
						//unit =   var_comand[1][0]
						fit_s := strings.ToLower(var_comand[1])
						////btype =  type_s[0]

						//fmt.Println("fit_s", string(fit_s))
						if fit_s == "bf" || fit_s == "ff" || fit_s == "wf" {
							en_fit = true
							fit = fit_s[0]
							//fmt.Println("        fit encontrado ", fit_s)
						} else {
							all_com_correc = false
							des_err = des_err + "( fit debe ser BF, FF ó WF) "
						}
					}

				}

			} else if command == "-name" {

				/*si hay duplicados*/
				if en_name == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					////////////////////
					name = var_comand[1]
					en_name, name, i = VerificarNameCom(name, i, comandos, var_comand)

					/*verificando si tiene valor del parametro*/
					/*if (len(var_comand[1]) == 0) { ///new
						all_com_correc = false
						des_err =  des_err + " -name "
					} else {
						//fmt.Println("var_comand[1]", var_comand[1])
						name = var_comand[1]
						//fmt.Println("unit", string(unit))
						en_name = true
						fmt.Print("     name encontrado-", name, "-\n" )
					}*/
					////////////////

				}
				/*OPCIONAL*/
			} else if command == "-unit" {

				/*si hay duplicados*/
				if en_unit == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 || len(var_comand[1]) > 1 { ///new
						all_com_correc = false
						des_err = des_err + " -unit "
					} else {
						//fmt.Println("var_comand[1]", var_comand[1])
						//unit =   var_comand[1][0]
						unidad_s := strings.ToLower(var_comand[1])
						unit = unidad_s[0]

						//fmt.Println("unit", string(unit))
						if unit == 'b' || unit == 'k' || unit == 'm' {
							en_unit = true
							/////////fmt.Println("        unit encontrado ", unit)
						} else {
							all_com_correc = false
							des_err = des_err + "( unit debe ser B, K ó M) "
						}
					}

				}

				/*OPCIONAL*/
			} else if command == "-delete" {

				/*si hay duplicados*/
				if en_delete == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					////////////////
					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 { ///new
						all_com_correc = false
						des_err = des_err + " -delete "
					} else {
						//fmt.Println("var_comand[1]", var_comand[1])
						delete = strings.ToLower(var_comand[1])

						if delete == "fast" || delete == "full" {
							en_delete = true
						} else {
							all_com_correc = false
							des_err = des_err + "( delete debe ser fast ó full) "
						}

						//fmt.Print("     delete encontrado-", delete, "-\n" )
					}
					////////////

				}

			} else {
				all_com_correc = false
				des_err = des_err + "(" + command + " No reconocido)"
				fmt.Println(command, "commando incorrecto*")
			}

		} else {
			des_err = des_err + "(" + var_comand[0] + " No reconocido)"
			all_com_correc = false
			fmt.Println("Error en FDISK, SOBRAN parametros o comando no reconocido")
		}
	}

	/*si no se ingreso K*/
	if en_unit == false {
		unit = 'k'
	} else {
		des_err_del = des_err_del + " -unit "
	}
	/*si no se ingreso P*/
	if en_type == false {
		btype = 'p'
	} else {
		des_err_del = des_err_del + " -type "
		des_err_add = des_err_add + " -type "
	}
	/*si no se ingreso WF = w*/
	if en_fit == false {
		fit = 'w'
	} else {
		des_err_del = des_err_del + " -fit "
		des_err_add = des_err_add + " -type "
	}

	if en_size == false && en_delete == false && en_add == false {

		des_err = des_err + " -size "
	}
	if en_path == false {
		des_err = des_err + " -path "
	}
	if en_name == false {
		des_err = des_err + " -name "
	}
	/////err de delete
	if en_size == true {
		des_err_del = des_err_del + " -size "
		des_err_add = des_err_add + " -size "
	}

	if en_delete == true {
		des_err_add = des_err_add + " -delete "
	}

	if en_add == true {
		des_err_del = des_err_del + " -delete "
	}

	/*si es delete, name y path*/
	if en_delete == true {

		if en_path == true && en_name == true && en_unit == false && en_type == false && en_fit == false && en_size == false && en_add == false && all_com_correc == true {
			fmt.Println("----comando eliminado")
			fmt.Println("path", path, "name", name, "delete", delete)

			Disks.DeletePartD(path, name, delete)
		} else {
			fmt.Println("Error en FDISK delete, ", des_err)
			fmt.Println("Parametros invalidos en delete", des_err_del)
		}

	} else if en_add == true {

		if en_path == true && en_name == true /*&& en_unit == true*/ && en_type == false && en_fit == false && en_size == false && all_com_correc == true {
			fmt.Println("----comando add")
			fmt.Println("path", path, "name", name, "add", add, "unit", string(unit))

			Disks.AddPartD(path, name, add, unit)
		} else {
			fmt.Println("Error en FDISK add, ", des_err)
			fmt.Println("Parametros invalidos en add", des_err_add)
		}

	} else {

		if en_size == false || en_path == false || en_name == false || all_com_correc == false {
			fmt.Println("Error en FDISK, ", des_err)
		} else {
			fmt.Println("comando valido")
			/*creando disco con los parametros*/
			//fmt.Println("size", size, "path", path, "name", name, "unit", string(unit))
			//fmt.Println("-size", size, "path", path, "name", name, "unit", string(unit), "type", string(btype), "fit", string(fit))

			/*se comento solo par pruebas*/
			////creando el disco
			Disks.CreatePartD(size, path, name, unit, btype, fit)

		}
	}

	///fmt.Println("fin despues de part id")

}

func MKdisk(comandos []string) {
	var size uint64
	var path string
	var pathcompleta string
	var prub []string
	var name string
	var unit byte

	en_fit := false //

	en_size := false
	en_path := false
	//en_name := false
	en_unit := false

	/*si todo esta correcto y yo que me alegro*/
	all_com_correc := true

	var des_err string = ""

	path = ""
	for i := 1; i < len(comandos); i++ {
		//fmt.Print(i,"-M(", comandos[i], ")\n")

		var_comand := strings.Split(comandos[i], "=")
		//fmt.Println("0",  var_comand,":" , len(var_comand))

		if len(var_comand) == 2 {

			command := strings.ToLower(var_comand[0])
			/**OBLIGATORIO**/
			if command == "-size" {

				/*si hay duplicados*/
				if en_size == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {
					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 { ///new
					} else {
						ss, _ := strconv.ParseUint(var_comand[1], 0, 64)
						size = ss

						if size > 0 {
							en_size = true
							//fmt.Println("        size encontrado ", size)
						} else {
							des_err = des_err + "(size, debe de ser > 0 ) "
						}

						/*fmt.Println("", var_comand[0],"")
						fmt.Println("*", len(var_comand[1]),"-")

						for y := 0; y < len(var_comand[1]) ; y++ {
							//fmt.Printf("%v",y ," %c", var_comand[1][y],"-\n")
							fmt.Println("-", var_comand[1][y],"-")

						}*/
					}
				}

			} else if command == "-path" {

				/*si hay duplicados*/
				if en_path == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					///////////////////
					pathcompleta = var_comand[1]
					en_path, pathcompleta, i = VerificarPath(pathcompleta, i, comandos, var_comand)

					prub = strings.Split(pathcompleta, "/")

					tam := len(prub)

					name = prub[tam-1]
					path = ""

					for i := 1; i < (tam - 1); i++ {

						path = path + "/" + prub[i]

					}
					path = path + "/"

				}

			} else if command == "-fit" {

				/*si hay duplicados*/
				if en_fit == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					/*verificando si tiene valor del parametro*/
					//fmt.Println("len(var_comand[1])", len(var_comand[1]))
					if len(var_comand[1]) == 0 || len(var_comand[1]) > 2 { ///new
						all_com_correc = false
						des_err = des_err + " -fit"
					} else {

						fit_s := strings.ToLower(var_comand[1])

						if fit_s == "bf" || fit_s == "ff" || fit_s == "wf" {
							en_fit = true
							//fit = fit_s[0]
							//fmt.Println("        fit encontrado ", fit_s)
						} else {
							all_com_correc = false
							des_err = des_err + "( fit debe ser BF, FF ó WF) "
						}
					}

				}

			} else if command == "-unit" {

				/*si hay duplicados*/
				if en_unit == true {
					all_com_correc = false
					des_err = des_err + "(" + command + " Duplicado) "
				} else {

					/*verificando si tiene valor del parametro*/
					if len(var_comand[1]) == 0 || len(var_comand[1]) > 1 { ///new
						all_com_correc = false
						des_err = des_err + " -unit "
					} else {
						//fmt.Println("var_comand[1]", var_comand[1])
						//unit =   var_comand[1][0]
						unidad_s := strings.ToLower(var_comand[1])
						unit = unidad_s[0]

						//fmt.Println("unit", string(unit))
						if unit == 'k' || unit == 'm' {
							en_unit = true
							//fmt.Println("        unit encontrado ", unit)
						} else {
							all_com_correc = false
							des_err = des_err + "( unit debe ser K ó M) "
						}
					}

				}

			} else {
				all_com_correc = false
				des_err = des_err + "(" + command + " No reconocido)"
				//fmt.Println(command, "commando incorrecto*")
			}

		} else {
			des_err = des_err + "(" + var_comand[0] + " No reconocido)"
			all_com_correc = false
			fmt.Println("Error en MKDISK, SOBRAN parametros o comando no reconocido")
		}
	}

	if en_unit == false {
		unit = 'm'
	}

	if en_size == false {

		des_err = des_err + " -size "
	}
	if en_path == false {
		des_err = des_err + " -path "
	}

	if en_size == false || en_path == false || all_com_correc == false {
		fmt.Println("Error en MKDISK, en", des_err)
	} else {
		fmt.Println("comando valido")
		/*creando disco con los parametros*/
		//fmt.Println("size", size, "path", path, "name", name, "unit", string(unit))

		////creando el disco
		Disks.CreateDisk(size, path, name, unit)
		//fmt.Println("new_disk", new_disk)
		//new_disk.CreateDisk(size, path, name, unit)

	}
}

func is_nameCorrect(name string) bool {
	letras_may := "ABCDEFGHIJKLMNÑOPQRSTUVWXYZ_"
	letras_min := "abcdefghijklmnñopqrstuvwxyz"
	numbers := "0123456789"

	caractes_cor := letras_may + letras_min + numbers

	if len(name) > 0 {
		for i := 0; i < len(name); i++ {
			correcto := strings.Contains(caractes_cor, string(name[i]))
			if correcto == false {
				return false
			}
		}
	} else {
		return false
	}
	return true

}

func tam_Correcto(tam int, tam_correcto int) bool {

	if tam == tam_correcto {
		return true
	} else {
		return false
	}
}

func ReadComando(comando string) {
	var ArrCommand []string
	ArrCommand = strings.Split(comando, " ")
	fmt.Println(ArrCommand)

	/*for i := 0; i < len(ArrCommand) ; i++ {
		fmt.Print(i,"-", ArrCommand[i], "\n")
		//resultado++
	}*/

	ArrCommand = ArrasinVaciones(ArrCommand)
	//fmt.Println(ArrCommand)
	//ejecutarComando(commandArray) //Ejecutamos el comando.

	for i := 0; i < len(ArrCommand); i++ {
		fmt.Print(i, "-", ArrCommand[i], "\n")
		//resultado++
	}
	fmt.Println("fin")
}

func ArrasinVaciones(command []string) []string {

	var arr_sinespacios []string
	for i := 0; i < len(command); i++ {
		if len(strings.ReplaceAll(command[i], " ", "")) != 0 {

			//fmt.Print(i,"000", command[i], "\n")
			arr_sinespacios = append(arr_sinespacios, strings.ReplaceAll(command[i], " ", ""))
		}
	}
	return arr_sinespacios
}
