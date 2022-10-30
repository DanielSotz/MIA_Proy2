package FormatP

import (
	"time"

	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	//para crear archivo binarop
	"log"
	"os"

	"strconv" // int to string

	"os/exec"
)

var report_tree_complet string
var report_directorio string
var report_archiv_incar string
var report_camino_file string

////SUPER BLOQUE
type SuperB struct {
	//nombre del disco duro
	Sb_nombre_hd [20]byte
	//cantidad de estructuras
	Sb_arbol_v_cant            int64
	Sb_detalle_directorio_cant int64
	Sb_inodo_cant              int64
	Sb_bloques_cant            int64

	Sb_arbol_virtual_free      int64
	Sb_detalle_directorio_free int64
	Sb_inodos_free             int64
	Sb_bloques_free            int64
	///fecha
	Sb_data_creacion       [19]byte
	Sb_data_ultimo_montaje [19]byte

	Sb_montajes_conta int64
	///apuntadores al inicio
	Sb_ap_bitmap_arbol_directorio_ini   int64
	Sb_ap_arbol_directorio_ini          int64
	Sb_ap_bitmap_detalle_directorio_ini int64
	Sb_ap_detalle_directorio_ini        int64
	Sb_ap_bitmap_tabla_inodo_ini        int64
	Sb_ap_tabla_inodo_ini               int64
	Sb_ap_bitmap_bloques_ini            int64
	Sb_ap_bloques_ini                   int64
	Sb_ap_log_ini                       int64
	////tamaños de las esctucturas
	Sb_size_struct_arbol_directorio   int64
	Sb_size_struct_detalle_directorio int64
	Sb_size_struct_inodo              int64
	Sb_size_struct_bloque             int64
	///primer bit en el bitmap
	Sb_first_free_bit_arbol_directorio   int64
	Sb_first_free_bit_detalle_directorio int64
	Sb_first_free_bit_tabla_inodo        int64
	Sb_first_free_bit_bloques            int64

	////mi carnet
	Sb_magic_num int64
}

////ARBOL VIRTUAL DE DIRECTORIO
type AVD struct {
	Avd_fecha_creacion    [19]byte
	Avd_nombre_directorio [20]byte

	///array de 6 subdireccionrios
	Avd_ap_array_subdirectorios [6]int64

	///apuntados a los detalles de directorios
	Avd_ap_detalle_directorio int64

	//apuntador indirecto
	Avd_ap_arbol_virtual_directorio int64

	///datos de propietario
	Avd_proper int64
	Avd_gid    int64
	Avd_perm   int64
}

////DETALLEDE DIRECTORIO
type DetDirectorio struct {
	Dd_array_files [5]DetDirINFO
	//apuntador indirecto
	Dd_ap_detalle_directorio_indirec int64
}

/*infotmacion al detalle de directorio*/
type DetDirINFO struct {
	Dd_file_nombre   [20]byte
	Dd_file_ap_inodo int64

	Dd_file_date_creacion     [19]byte
	Dd_file_date_modificacion [19]byte
}

////TABLA DE I-NODO
type Inodo struct {
	I_count_inodo             int64
	I_size_archivo            int64
	I_count_bloques_asignados int64

	///array de 4 bloques
	I_array_bloques [4]int64

	//apuntador indirecto
	I_ap_indirecto int64

	///datos de propietario
	I_id_proper int64
	I_gid       int64
	I_perm      int64
}

//BLOQUE DE DATOS
type BloqueDatos struct {
	Db_id   int64 //////
	Db_data [25]byte
}

//LOG/BITACORA
type LogBitacora struct {
	Log_tipo_operación [20]byte
	Log_tipo           int64
	Log_nombre         [20]byte
	Log_contenido      [50]byte
	Log_fecha          [19]byte
}

type Letra struct {
	Let byte
}

/*func main() {
	Format()
}*/

func Format(inicia_part uint64, part_size uint64, path string) {

	//var inicio_particion int64 = 0
	//var Tam_Particion int64 = 20 * 1024 * 1024

	var inicio_particion int64
	var Tam_Particion int64

	inicio_particion = int64(inicia_part)
	Tam_Particion = int64(part_size)

	//fmt.Println("\ninicio_particion",inicio_particion)
	//fmt.Println("Tam_Particion",Tam_Particion)

	/*tomando los tamaños de las escrituras*/
	var size_avd int64 = int64(binary.Size(AVD{}))
	var size_det_directorio int64 = int64(binary.Size(DetDirectorio{}))
	var size_inodo int64 = int64(binary.Size(Inodo{}))
	var size_bloque int64 = int64(binary.Size(BloqueDatos{}))
	var size_bitacora int64 = int64(binary.Size(LogBitacora{}))
	var size_SB int64 = int64(binary.Size(SuperB{}))

	var NumEstructuras int64 = (Tam_Particion - (2 * size_SB)) / (27 + size_avd + size_det_directorio + (5*size_inodo + (20 * size_bloque) + size_bitacora))

	//fmt.Println("\nNumEstructuras",NumEstructuras)

	cant_AVD := NumEstructuras
	cant_DetDir := NumEstructuras
	cant_Inodos := 5 * NumEstructuras
	cant_Bloques := 20 * NumEstructuras // 4*cant_Inodos

	//cant_Bita := NumEstructuras

	Ini_bitmapAVD := inicio_particion + size_SB
	Ini_AVD := Ini_bitmapAVD + cant_AVD
	Ini_bitmapDD := Ini_AVD + (size_avd * cant_AVD)
	Ini_DD := Ini_bitmapDD + cant_DetDir
	Ini_bitmapInodo := Ini_DD + (size_det_directorio * cant_DetDir)
	Ini_Inodos := Ini_bitmapInodo + cant_Inodos
	Ini_bitmapBloque := Ini_Inodos + (size_inodo * cant_Inodos)
	Ini_Bloque := Ini_bitmapBloque + cant_Bloques
	Ini_Bitacora := Ini_Bloque + (size_bloque * cant_Bloques)
	//Ini_CopSB :=  Ini_Bitacora + size_bitacora*cant_Bita

	sB := SuperB{}
	//nombre del disco duro
	name_d := strings.Split(path, "/")
	name_disk_str := name_d[len(name_d)-1]
	copy(sB.Sb_nombre_hd[:], name_disk_str)
	//cantidad de estructuras
	sB.Sb_arbol_v_cant = cant_AVD
	sB.Sb_detalle_directorio_cant = cant_DetDir
	sB.Sb_inodo_cant = cant_Inodos
	sB.Sb_bloques_cant = cant_Bloques

	sB.Sb_arbol_virtual_free = cant_AVD
	sB.Sb_detalle_directorio_free = cant_DetDir
	sB.Sb_inodos_free = cant_Inodos
	sB.Sb_bloques_free = cant_Bloques

	///fecha
	fecha_string := time.Now().Format("2006-01-02 15:04:05")
	copy(sB.Sb_data_creacion[:], fecha_string)
	sB.Sb_data_ultimo_montaje = sB.Sb_data_creacion

	sB.Sb_montajes_conta = 0
	///apuntadores al inicio
	/*Ini_CopSB :=  Ini_Bitacora + size_bitacora*cant_Bita
	fin_part :=  Ini_CopSB + size_SB*/
	sB.Sb_ap_bitmap_arbol_directorio_ini = Ini_bitmapAVD
	sB.Sb_ap_arbol_directorio_ini = Ini_AVD
	sB.Sb_ap_bitmap_detalle_directorio_ini = Ini_bitmapDD
	sB.Sb_ap_detalle_directorio_ini = Ini_DD
	sB.Sb_ap_bitmap_tabla_inodo_ini = Ini_bitmapInodo
	sB.Sb_ap_tabla_inodo_ini = Ini_Inodos
	sB.Sb_ap_bitmap_bloques_ini = Ini_bitmapBloque
	sB.Sb_ap_bloques_ini = Ini_Bloque
	sB.Sb_ap_log_ini = Ini_Bitacora
	////tamaños de las esctucturas
	/*var size_bitacora int64 = int64(binary.Size(LogBitacora{}))
	var size_SB int64 = int64(binary.Size(SuperB{}))*/
	sB.Sb_size_struct_arbol_directorio = size_avd
	sB.Sb_size_struct_detalle_directorio = size_det_directorio
	sB.Sb_size_struct_inodo = size_inodo
	sB.Sb_size_struct_bloque = size_bloque
	///primer bit en el bitmap
	sB.Sb_first_free_bit_arbol_directorio = 1
	sB.Sb_first_free_bit_detalle_directorio = 1
	sB.Sb_first_free_bit_tabla_inodo = 1
	sB.Sb_first_free_bit_bloques = 1

	////mi carnet
	sB.Sb_magic_num = 201430496
	/*Guardando SuperBoot en disco*/
	/*abrimos el archivo con todos los permisos*/
	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(inicio_particion, 0)
	disk := &sB

	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, disk)
	msg_exito := "Se Formateo Partición " + /*string(name_byte[:])  + */ " con Exito"
	writeBytesInDisk(file, binario.Bytes(), msg_exito)

	/*setenado 0 al mapa de bits*/

	// Bitmap de AVDs
	var cero byte = '0'
	c0 := &cero
	/*var uno byte = '1';
	u1 := &uno*/
	file.Seek(Ini_bitmapAVD, 0)

	for i := 0; i < int(cant_AVD); i++ {
		/*if (i == 2 || i == 5 || i == 27 ) {
			var binario_1 bytes.Buffer
			binary.Write(&binario_1, binary.BigEndian, u1)
			writeBytesDisk(file, binario_1.Bytes())
			//fmt.Println("iiii", i)
		} else  {*/
		var binario_0 bytes.Buffer
		binary.Write(&binario_0, binary.BigEndian, c0)
		writeBytesDisk(file, binario_0.Bytes())
		//}
	}

	/*bitma detalle de directorios*/
	file.Seek(Ini_bitmapDD, 0)
	for i := 0; i < int(sB.Sb_detalle_directorio_cant); i++ {
		var binario_0 bytes.Buffer
		binary.Write(&binario_0, binary.BigEndian, c0)
		writeBytesDisk(file, binario_0.Bytes())
	}

	/*bitma tabla de inodos*/
	file.Seek(Ini_bitmapInodo, 0)
	for i := 0; i < int(cant_Inodos); i++ {
		var binario_0 bytes.Buffer
		binary.Write(&binario_0, binary.BigEndian, c0)
		writeBytesDisk(file, binario_0.Bytes())
	}

	/*bitma de  bloques*/
	file.Seek(Ini_bitmapBloque, 0)
	for i := 0; i < int(cant_Bloques); i++ {
		var binario_0 bytes.Buffer
		binary.Write(&binario_0, binary.BigEndian, c0)
		writeBytesDisk(file, binario_0.Bytes())
	}

	///por defecto creo indice 1
	/////buscando pos en bitmap arbol
	pos_bm_arbol := Pos_bitmap_AVD(file, sB)
	//fmt.Println("pos_bm_arbol", pos_bm_arbol)

	path_archivo := "/"
	CreateCarpeta(path, path_archivo, Ini_bitmapAVD, Ini_AVD, int64(pos_bm_arbol), size_avd)

	///*escribiendo detalle de Directorios
	/*path_archivo = "/users.txt"
	os_bm_detalle := Pos_bitmap_Detalles(file, sB)
	fmt.Println("os_bm_detalle", os_bm_detalle)
	CreateDetalle(file, path_archivo, Ini_AVD) */

	/////////////////////////////////////////////////////////////////////////////
	fmt.Println("Formateo de bitma terminado")

	/*archivo a crear*/
	/////////////////////////archivo_crear := "/home/users.txt"
	archivo_crear := "/users.txt"
	Contenido_ar := "1,G,root\n1,U,root,root,123\n"
	fmt.Println("Contenido_ar", Contenido_ar)
	var_directorios := strings.Split(archivo_crear, "/")
	if len(var_directorios[0]) == 0 {
		var_directorios[0] = "/"
	}
	//fmt.Println("var_directorios", var_directorios)
	Creando_Archivos(var_directorios, file, Ini_AVD, sB, Contenido_ar, 'S', 'A')

	//////////////////
	//RecorroDirectorio(path, inicia_part)
}

func NewFile(inicia_part uint64, part_size uint64, path string, size_fil int64, p byte, cont string, dire_file string) {

	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(int64(inicia_part), 0)

	sB := SuperB{}
	var size int = int(binary.Size(sB))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &sB)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario", err)
	}
	//fmt.Println(sB)

	/*archivo a crear*/
	archivo_crear := dire_file
	Contenido_ar := cont

	if (len(cont) > int(size_fil)) && int(size_fil) == 0 {
		size_fil = int64(len(cont))
	}

	var let byte = 65
	if len(cont) < int(size_fil) {

		for i := len(cont); i < int(size_fil); i++ {
			Contenido_ar = Contenido_ar + string(let)
			let++
			if let == 91 {
				let = 65
			}
		}
	} else if len(cont) > 0 {

		if len(cont) > int(size_fil) {
			Contenido_ar = cont[0:int(size_fil)]
		}
	}

	var_directorios := strings.Split(archivo_crear, "/")
	if len(var_directorios[0]) == 0 {
		var_directorios[0] = "/"
	}
	fmt.Println("var_directorios", var_directorios, len(var_directorios))
	fmt.Println("size_fil", size_fil)
	fmt.Println("Contenido_ar", Contenido_ar)

	var tipo_mk byte = 'A'
	Ini_AVD := sB.Sb_ap_arbol_directorio_ini
	Creando_Archivos(var_directorios, file, Ini_AVD, sB, Contenido_ar, p, tipo_mk)

	//////////RecorroDirectorio(path, inicia_part)
}

/*editando archivo*/
func Edit_File(inicia_part uint64, part_size uint64, path string, size_fil int64, cont string, dire_file string) {

	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(int64(inicia_part), 0)

	sB := SuperB{}
	var size int = int(binary.Size(sB))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &sB)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario", err)
	}
	//fmt.Println(sB)

	/*archivo a crear*/
	archivo_crear := dire_file
	Contenido_ar := cont

	if (len(cont) > int(size_fil)) && int(size_fil) == 0 {
		size_fil = int64(len(cont))
	}

	var let byte = 65
	if len(cont) < int(size_fil) {

		for i := len(cont); i < int(size_fil); i++ {
			Contenido_ar = Contenido_ar + string(let)
			let++
			if let == 91 {
				let = 65
			}
		}
	} else if len(cont) > 0 {

		if len(cont) > int(size_fil) {
			Contenido_ar = cont[0:int(size_fil)]
		}
	}

	var_directorios := strings.Split(archivo_crear, "/")
	if len(var_directorios[0]) == 0 {
		var_directorios[0] = "/"
	}
	//fmt.Println("var_directorios", var_directorios, len(var_directorios))
	//fmt.Println("size_fil", size_fil)
	//fmt.Println("Contenido_ar", Contenido_ar)

	var tipo_accion string
	var tipo_mk byte = 'A'
	Ini_AVD := sB.Sb_ap_arbol_directorio_ini

	tipo_accion = "MO_ARCV"
	Leyendo_Recorrido(file, Ini_AVD, var_directorios, 0, sB, tipo_mk, tipo_accion, "", Contenido_ar)

	/*Creando_Archivos(var_directorios, file, Ini_AVD, sB, Contenido_ar, p, tipo_mk)*/

	//////////RecorroDirectorio(path, inicia_part)
}
func NewDirectorio(inicia_part uint64, part_size uint64, path string, p byte, dire_file string) {

	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(int64(inicia_part), 0)

	sB := SuperB{}
	var size int = int(binary.Size(sB))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &sB)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario", err)
	}
	//fmt.Println(sB)

	/*carpeta a crear*/
	carpeta_crear := dire_file
	Contenido_ar := ""

	var_directorios := strings.Split(carpeta_crear, "/")
	if len(var_directorios[0]) == 0 {
		var_directorios[0] = "/"
	}
	fmt.Println("var_directorios", var_directorios, len(var_directorios))
	//fmt.Println("Contenido_ar", Contenido_ar)

	var tipo_mk byte = 'C'
	Ini_AVD := sB.Sb_ap_arbol_directorio_ini
	Creando_Archivos(var_directorios, file, Ini_AVD, sB, Contenido_ar, p, tipo_mk)

	/////////////RecorroDirectorio(path, inicia_part)
}

func Creando_Archivos(var_directorios []string, file *os.File, Ini_AVD int64, sB SuperB, Contenido_ar string, p byte, tipo_mk byte) {

	///ar_tipo: P = crea todoas las carpetas
	///ar_tipo: N = si no existe carpeta padre, error

	for i := 0; i < len(var_directorios); i++ {
		fmt.Println(i, var_directorios[i])
		/*voy recorriendo el directorio*/
	}

	Recorriendo_directorio(file, Ini_AVD, var_directorios, 0, sB, Contenido_ar, p, tipo_mk)
}

func Existe_Carpeta(file *os.File, pos_carpeta_bus int64, carpeta_a_buscar string) bool {

	//fmt.Println(i,"--------------car actual--",var_directorios[i])
	//fmt.Println(i,"--------------------------ini avd",sB.Sb_ap_arbol_directorio_ini,"size" ,sB.Sb_size_struct_arbol_directorio)
	//fmt.Println(i,"--------------------------pos_actual AVD" ,pos_AVD)

	file.Seek(int64(pos_carpeta_bus), 0)

	Avds := AVD{}
	var size int = int(binary.Size(Avds))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &Avds)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario", err)
	}
	/////////////fmt.Printf("vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv   Avd_nombre_directorio   %s\n", Avds.Avd_nombre_directorio)

	var name_actual [20]byte
	copy(name_actual[:], carpeta_a_buscar)

	//fmt.Println("name_actual", name_actual)
	if name_actual == Avds.Avd_nombre_directorio {
		///////////////fmt.Println("---Avds", "son igualitos, exite")
		return true
	}
	return false

}

func Recorriendo_directorio(file *os.File, pos_AVD int64, var_directorios []string, i int64, sB SuperB, Contenido_ar string, p byte, tipo_mk byte) {

	//fmt.Println(i, "--------------car actual--", var_directorios[i])
	//fmt.Println(i, "--------------------------ini avd", sB.Sb_ap_arbol_directorio_ini, "size", sB.Sb_size_struct_arbol_directorio)
	//fmt.Println(i, "--------------------------pos_actual AVD", pos_AVD)

	file.Seek(int64(pos_AVD), 0)

	Avds := AVD{}
	var size int = int(binary.Size(Avds))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &Avds)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario", err)
	}
	/////////////////////////////////////fmt.Println("Avds", Avds)

	var name_actual [20]byte
	copy(name_actual[:], var_directorios[i])
	//fmt.Printf("				name_actual   %s\n", name_actual)
	//fmt.Printf("				Avd_nombre_directorio   %s\n", Avds.Avd_nombre_directorio)

	////si es igual, me quedo ahi, si no busco en el indiercto
	if name_actual == Avds.Avd_nombre_directorio {
		//fmt.Println("Avds", "son igualitos, exite")
		/*verificando la carpeta siguiente*/
		i++
		/*verificando si el id no soprepase*/
		if int(i) >= len(var_directorios) {
			return
		}
		//fmt.Println("			siguinte", var_directorios[i])
		/*verigico si es carpeta o archivo*/

		var tipo_insert byte

		esfinal := false
		/*archiv := strings.Split(var_directorios[i], ".")
		fmt.Println("			archiv", archiv)
		if archiv[len(archiv)-1] == "txt" {
			fmt.Println("*************************************** ES ARCHIVO")
			tipo_insert = 'A'
		} else {
			tipo_insert = 'C'
		}*/
		tipo_insert = 'C'

		if int(i) == (len(var_directorios) - 1) {
			//fmt.Println("*************************************** FINAL FINAL FINAL FINAL", var_directorios[i])
			tipo_insert = tipo_mk
			esfinal = true
		}
		//fmt.Println("***************************************esfinal*", esfinal)

		//fmt.Println("***********************************************", string(tipo_insert))

		//////******************INICIO**CREACION*DE*CARPETAS*****************************************************////
		if tipo_insert == 'C' {
			/*recoriiendo los subdirectorios para encontar la carpeta siguiente*/
			encotrado := false
			var pos_carpeta_bus int64
			for j := 0; j < len(Avds.Avd_ap_array_subdirectorios); j++ {
				//fmt.Println(j,"MMMMMMMMMMMMMMMMMMMMM", Avds.Avd_ap_array_subdirectorios[j])
				/*Aqui verifico si existe la carpeta xisuitente*/
				if Avds.Avd_ap_array_subdirectorios[j] != -1 {
					pos_carpeta_bus = sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * Avds.Avd_ap_array_subdirectorios[j])
					encotrado = Existe_Carpeta(file, pos_carpeta_bus, var_directorios[i])

					if encotrado == true {
						break
					}
				}
			}

			//fmt.Println("encotrado", encotrado)
			/*no existe, entonces creo*/
			if encotrado == false {
				//busco la primera pos de los apuntadores de avd
				encotrado_dirlibre := false
				var index int
				for i := 0; i < len(Avds.Avd_ap_array_subdirectorios); i++ {
					if Avds.Avd_ap_array_subdirectorios[i] == -1 {
						index = i
						encotrado_dirlibre = true
						break
					}
				}
				//fmt.Println("----------------------------Avds index", index)
				//fmt.Println("----------------------------Avds encotrado_dirlibre", encotrado_dirlibre)

				if encotrado_dirlibre == false {
					//fmt.Println("0000000000000 no se encontro libre normal")

					//fmt.Println("Avds.Avd_ap_arbol_virtual_directorio", Avds.Avd_ap_arbol_virtual_directorio)
					if Avds.Avd_ap_arbol_virtual_directorio == -1 {
						/*creo indirecto y luego creo la nueva carpeta*/

						/*verificando si crea padres anteriores*/
						if esfinal == false && p == 'N' {

							fmt.Println("Err, No existe Carpeta Padre (", var_directorios[i], ")")
							return
						}

						/*creo indirecto*/
						////////////////////////////////////
						/*buscando la pos segun el bitmap*/
						index_bm_arbol := Pos_bitmap_AVD(file, sB)
						cord_avd := sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * int64(index_bm_arbol))
						//fmt.Println("Avds indirecto index_bm_arbol", index_bm_arbol)
						//fmt.Println("Avds indirecto cord_avd", cord_avd)

						/*creando nueva carpetas*/
						Avds_new_indi := AVD{}
						//NameByteToString(Avds.Avd_nombre_directorio)
						Avds_new_indi = CreateNewCarpeta_indirec(file, string(Avds.Avd_nombre_directorio[:]), Avds_new_indi, cord_avd, int64(index_bm_arbol), sB)
						/*actualizando AVDS*/
						Avds.Avd_ap_arbol_virtual_directorio = int64(index_bm_arbol)
						//fmt.Println("++++++++ Avds.Avd_ap_arbol_virtual_directorio", Avds.Avd_ap_arbol_virtual_directorio)
						//fmt.Println("uuuuuuuuuuuuuu Avds_new",Avds_new_indi)
						//fmt.Println("uuuuuuuuuuuuuu Avds",Avds)
						UpdateAVD(file, Avds, pos_AVD)
						////////////////////////////////////

						//////////////////////////////////////creo la carpeta nueva
						file.Seek(int64(cord_avd), 0)

						Avds_indi := AVD{}
						var size_indi int = int(binary.Size(Avds_indi))
						data_indi := readBytesDisk(file, size_indi)
						buffer_indi := bytes.NewBuffer(data_indi)

						err := binary.Read(buffer_indi, binary.BigEndian, &Avds_indi)
						if err != nil {
							fmt.Println("binary. se encontro error al leer archivo binario", err)
						}

						index_bm_arbol_in := Pos_bitmap_AVD(file, sB)
						cord_avd_in := sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * int64(index_bm_arbol_in))
						//fmt.Println("Avds indi_new index_bm_arbol_in", index_bm_arbol_in)
						//fmt.Println("Avds indi_new cord_avd_in", cord_avd_in)

						///*creando nueva carpetas
						Avds_new_in_indirec := AVD{}
						CreateNewCarpeta(file, var_directorios[i], Avds_new_in_indirec, cord_avd_in, int64(index_bm_arbol_in), sB)
						///*actualizando AVDS
						Avds_new_indi.Avd_ap_array_subdirectorios[0] = int64(index_bm_arbol_in)
						UpdateAVD(file, Avds_new_indi, cord_avd)
						/////////////////////////////////

						/*recursivamente a la carpeta*/
						Recorriendo_directorio(file, cord_avd_in, var_directorios, i, sB, Contenido_ar, p, tipo_mk)

					} else {
						///*SI YA EXISTE APUNTADOR INDURECTO, ENTONCES EM POSICIONO EN ELLA*/
						//fmt.Println("SI YA EXISTE APUNTADOR INDURECTO, ENTONCES EM POSICIONO EN ELLA")
						cord_avd_ind := sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * int64(Avds.Avd_ap_arbol_virtual_directorio))
						i = i - 1
						Recorriendo_directorio(file, cord_avd_ind, var_directorios, i, sB, Contenido_ar, p, tipo_mk)
					}

				} else {

					///////////////////////////////////////////
					////fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@22")

					/*verificando si crea padres anteriores*/
					//fmt.Println("esfinal", esfinal, "p", string(p))
					if esfinal == false && p == 'N' {

						fmt.Println("2 Err, No existe Carpeta Padre (", var_directorios[i], ")")
						return
					}
					/*buscando la pos segun el bitmap*/
					index_bm_arbol := Pos_bitmap_AVD(file, sB)
					cord_avd := sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * int64(index_bm_arbol))
					//fmt.Println("Avds index_bm_arbol", index_bm_arbol)
					//fmt.Println("Avds cord_avd", cord_avd)

					/*creando en index*/
					Avds.Avd_ap_array_subdirectorios[index] = int64(index_bm_arbol)
					/*creando nueva carpetas*/
					Avds_new := AVD{}
					CreateNewCarpeta(file, var_directorios[i], Avds_new, cord_avd, int64(index_bm_arbol), sB)
					/*actualizando AVDS*/
					Avds.Avd_ap_array_subdirectorios[index] = int64(index_bm_arbol)
					//fmt.Println("++++++++ Avds.Avd_ap_arbol_virtual_directorio", Avds.Avd_ap_arbol_virtual_directorio)
					UpdateAVD(file, Avds, pos_AVD)

					/*recursivamente a la carpeta*/
					Recorriendo_directorio(file, cord_avd, var_directorios, i, sB, Contenido_ar, p, tipo_mk)
					/////////////////////////////////
				}

			} else {
				//fmt.Println("Avds", "LA CARPETA YA EXISTE, ME VOY EN ELLA")
				/*recursivamente a la carpeta*/
				Recorriendo_directorio(file, pos_carpeta_bus, var_directorios, i, sB, Contenido_ar, p, tipo_mk)

			}
			//////******************FIN**CREACION*DE*CARPETAS*****************************************************////

			//////******************INICIO**CREACION*DE*ARCHIVOS*****************************************************////
		} else if tipo_insert == 'A' {
			/*buscando apuntador de archivos*/
			var detDir DetDirectorio
			var index_array_DD int64
			var index_array_inodo int64

			var pos_detdir int64 /////por ahora toma el detdir actual
			if Avds.Avd_ap_detalle_directorio == -1 {
				//fmt.Println("DetDir", "no exite DEtalle de directorio, entonces lo CREO")

				/*buscando la pos segun en el bitmap*/
				index_bm_DD := Pos_bitmap_Detalles(file, sB)
				cord_DD := sB.Sb_ap_detalle_directorio_ini + (sB.Sb_size_struct_detalle_directorio * int64(index_bm_DD))
				//fmt.Println("DD index_bm_DD", index_bm_DD)
				//fmt.Println("DD cord_DD", cord_DD)

				pos_detdir = cord_DD ////////-----------------------------

				/////////////////////////////
				detDir_new := DetDirectorio{}
				index_array_DD = 0
				NewDetalleDirectorio(file, var_directorios[i], &detDir_new, cord_DD, int64(index_bm_DD), sB, index_array_DD)
				/*reviso primer index libre*/
				/*como es nuevo, entonces hago el primero del inodo*/
				/*actualizando AVDS*/
				Avds.Avd_ap_detalle_directorio = int64(index_bm_DD)
				UpdateAVD(file, Avds, pos_AVD)
				//fmt.Println("Avds.Avd_ap_detalle_directorio", Avds.Avd_ap_detalle_directorio)

				/*como  es nuevo, entonces uso el primer inodo*/
				detDir = detDir_new ///////////
				index_array_inodo = 0
				//fmt.Printf("+++++++++++++++detDir.Dd_array_files[index_array_DD].Dd_file_nombre  %s\n", detDir.Dd_array_files[index_array_DD].Dd_file_nombre)
				//fmt.Printf("+++++++++++++++detDir_new.Dd_array_files[index_array_DD].Dd_file_nombre  %s\n", detDir_new.Dd_array_files[index_array_DD].Dd_file_nombre)

				/*primer inodo**/
				Creando_First_Inodo(file, var_directorios, i, sB, Contenido_ar, index_array_DD, index_array_inodo, detDir, pos_detdir)
			} else {
				//fmt.Println("IIIIIIII si existe", "directorio reviso ahi")
				/*si ya existe detalle, entonces me busco la posicion del
				detalle y voy insertando los inodos*/

				/*me posiciono en Det
				luego busco espacio disponible*/
				cord_DD_exis := sB.Sb_ap_detalle_directorio_ini + (sB.Sb_size_struct_detalle_directorio * Avds.Avd_ap_detalle_directorio)
				//fmt.Println("existente DD cord_DD_exis", cord_DD_exis)
				pos_detdir = cord_DD_exis ////////-----------------------------

				file.Seek(int64(cord_DD_exis), 0)

				detDir_exis := DetDirectorio{}
				var size int = int(binary.Size(detDir_exis))
				data := readBytesDisk(file, size)
				buffer := bytes.NewBuffer(data)

				err := binary.Read(buffer, binary.BigEndian, &detDir_exis)
				if err != nil {
					fmt.Println("binary. se encontro error al leer archivo binario", err)
				}

				detDir = detDir_exis
				index_array_inodo = 0

				//detDir_exis, pos_detdir, index_array_DD = BuscandoDetalleDisponible(file, detDir_exis,sB, Contenido_ar, var_directorios[i], cord_DD_exis, var_directorios, i)
				BuscandoDetalleDisponible(file, detDir_exis, sB, Contenido_ar, var_directorios[i], cord_DD_exis, var_directorios, i)
				/*primer inodo**/
				//Creando_First_Inodo(file, var_directorios, i, sB, Contenido_ar, index_array_DD, index_array_inodo, detDir, pos_detdir)

			}

		}
		//////******************FIN**CREACION*DE*ARCHIVOS*****************************************************////

		////Recorriendo_directorio(file, Ini_AVD, var_directorios, i)
	} else {

		//fmt.Println("Avds", "no exite")
		//fmt.Println("Avds.Avd_ap_detalle_directorio", Avds.Avd_ap_detalle_directorio)
		/*no existe*/
	}
}

func Rename_CarIndirectos(file *os.File, pos_AVD int64, sB SuperB, new_name string) {

	file.Seek(int64(pos_AVD), 0)

	Avds := AVD{}
	var size int = int(binary.Size(Avds))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &Avds)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario", err)
	}

	copy(Avds.Avd_nombre_directorio[:], new_name)
	//update avd
	UpdateAVD(file, Avds, pos_AVD)

	if Avds.Avd_ap_arbol_virtual_directorio != -1 {

		///*SI YA EXISTE APUNTADOR INDURECTO, ENTONCES EM POSICIONO EN ELLA*/
		//fmt.Println("SI YA EXISTE APUNTADOR INDURECTO, ENTONCES EM POSICIONO EN ELLA")
		cord_avd_ind := sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * int64(Avds.Avd_ap_arbol_virtual_directorio))
		//i = i-1
		Rename_CarIndirectos(file, cord_avd_ind, sB, new_name)
	}
}

var posicion_de_carpeta_rep int64

///////////////////////////////recorriendo leyendo inicio
//Leyendo_Recorrido(file, pos_AVD, var_directorios, 0, sB, tipo_mk)
func Leyendo_Recorrido(file *os.File, pos_AVD int64, var_directorios []string, i int64, sB SuperB, tipo_mk byte, tipo_accion string, new_name string, Contenido_ar string) {

	////////////fmt.Println(i,"--------------car actual--",var_directorios[i])
	//fmt.Println(i,"--------------------------ini avd",sB.Sb_ap_arbol_directorio_ini,"size" ,sB.Sb_size_struct_arbol_directorio)
	//fmt.Println(i,"--------------------------pos_actual AVD" ,pos_AVD)

	var graf_tempo string = ""
	var conexiones string = ""

	file.Seek(int64(pos_AVD), 0)

	Avds := AVD{}
	var size int = int(binary.Size(Avds))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &Avds)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario", err)
	}
	/////////////////////////////////////fmt.Println("Avds", Avds)

	var name_actual [20]byte
	copy(name_actual[:], var_directorios[i])
	//////////////fmt.Printf("				name_actual   %s\n", name_actual)
	//////////////fmt.Printf("				Avd_nombre_directorio   %s\n", Avds.Avd_nombre_directorio)

	////si es igual, me quedo ahi, si no busco en el indiercto
	if name_actual == Avds.Avd_nombre_directorio {
		/////////////////fmt.Println("Avds", "son igualitos, exite")
		encontrado_tipo_art = true

		if tipo_accion == "POS_CAR" {
			posicion_de_carpeta_rep = pos_AVD
		}

		/*verificando la carpeta siguiente*/
		i++
		/*verificando si el id no soprepase*/
		if int(i) >= len(var_directorios) {

			if tipo_accion == "MOD_NAME" {
				copy(Avds.Avd_nombre_directorio[:], new_name)
				//update avd
				UpdateAVD(file, Avds, pos_AVD)

				if Avds.Avd_ap_arbol_virtual_directorio != -1 {
					cord_avd_ind := sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * int64(Avds.Avd_ap_arbol_virtual_directorio))
					Rename_CarIndirectos(file, cord_avd_ind, sB, new_name)
				}

				fmt.Println("Se actualizó nombre con exito")

			}

			return
		}
		//////////////////fmt.Println("			siguinte", var_directorios[i])
		/*verigico si es carpeta o archivo*/

		var tipo_insert byte

		////////esfinal := false

		tipo_insert = 'C'

		if int(i) == (len(var_directorios) - 1) {
			/////////////fmt.Println("*************************************** FINAL FINAL FINAL FINAL", var_directorios[i])
			tipo_insert = tipo_mk
			////////////esfinal = true
		}
		////////////////fmt.Println("***************************************esfinal*", esfinal)

		////////////fmt.Println("***********************************************", string(tipo_insert))

		//////******************INICIO**CREACION*DE*CARPETAS*****************************************************////
		if tipo_insert == 'C' {
			/*recoriiendo los subdirectorios para encontar la carpeta siguiente*/
			encotrado := false
			var pos_carpeta_bus int64
			var index_carcar int64
			for j := 0; j < len(Avds.Avd_ap_array_subdirectorios); j++ {
				/*Aqui verifico si existe la carpeta xisuitente*/
				if Avds.Avd_ap_array_subdirectorios[j] != -1 {
					pos_carpeta_bus = sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * Avds.Avd_ap_array_subdirectorios[j])
					encotrado = Existe_Carpeta(file, pos_carpeta_bus, var_directorios[i])

					if encotrado == true {
						index_carcar = int64(j)
						break
					}
				}
			}

			//	fmt.Println("encotrado", encotrado)

			/*no existe, entonces creo*/
			if encotrado == false {

				//fmt.Println("Avds.Avd_ap_arbol_virtual_directorio", Avds.Avd_ap_arbol_virtual_directorio)
				if Avds.Avd_ap_arbol_virtual_directorio == -1 {
					/*creo indirecto y luego creo la nueva carpeta*/

					/*verificando si crea padres anteriores*/
					//fmt.Println("esfinal", esfinal )
					//if esfinal == false  {

					fmt.Println("Err, No existe Carepta (", var_directorios[i], ")")
					encontrado_tipo_art = false
					return
					//}

				} else {
					///*SI YA EXISTE APUNTADOR INDURECTO, ENTONCES EM POSICIONO EN ELLA*/
					//fmt.Println("SI YA EXISTE APUNTADOR INDURECTO, ENTONCES EM POSICIONO EN ELLA")
					cord_avd_ind := sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * int64(Avds.Avd_ap_arbol_virtual_directorio))
					i = i - 1

					if tipo_accion == "REP_TF" {

						//idavl := "AVD"+ strconv.Itoa(int(avd_index))
						idavl := "AVD" + strconv.Itoa(int(pos_AVD))

						graf_tempo =
							idavl + "[label=<\n" +
								"<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n"

						report_camino_file = report_camino_file + graf_tempo

						graf_tempo = "<TR port='0'>\n" +
							"<TD COLSPAN='2'><font color='black'>" + NameByteToString(Avds.Avd_nombre_directorio) + "</font></TD>\n" +
							"</TR>\n" +

							"<TR>\n" +
							"<TD BGCOLOR='#99ccff'>Avd_fecha_creacion</TD><TD>" + string(Avds.Avd_fecha_creacion[:]) + "</TD>\n" +
							"</TR>\n"
						report_camino_file = report_camino_file + graf_tempo

						graf_tempo =
							"<TR>\n" +
								"<TD BGCOLOR='#99ccff'>aptr_ind</TD><TD port='8'>" + strconv.Itoa(int(Avds.Avd_ap_arbol_virtual_directorio)) + "</TD>\n" +
								"</TR>\n"

						report_camino_file = report_camino_file + graf_tempo
						conexiones = conexiones + idavl + ":8->" + "AVD" + strconv.Itoa(int(cord_avd_ind)) + "\n"

						graf_tempo = "</TABLE>\n" +
							">];\n" +
							conexiones
						report_camino_file = report_camino_file + graf_tempo

					}

					Leyendo_Recorrido(file, cord_avd_ind, var_directorios, i, sB, tipo_mk, tipo_accion, new_name, Contenido_ar)
				}

				///////////
			} else {

				if tipo_accion == "REP_TF" {

					//idavl := "AVD"+ strconv.Itoa(int(avd_index))
					idavl := "AVD" + strconv.Itoa(int(pos_AVD))

					graf_tempo =
						idavl + "[label=<\n" +
							"<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n"

					report_camino_file = report_camino_file + graf_tempo

					graf_tempo = "<TR port='0'>\n" +
						"<TD COLSPAN='2'><font color='black'>" + NameByteToString(Avds.Avd_nombre_directorio) + "</font></TD>\n" +
						"</TR>\n" +

						"<TR>\n" +
						"<TD BGCOLOR='#99ccff'>Avd_fecha_creacion</TD><TD>" + string(Avds.Avd_fecha_creacion[:]) + "</TD>\n" +
						"</TR>\n"
					report_camino_file = report_camino_file + graf_tempo

					//pos_carpeta_bus := 	sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * int64( Avds.Avd_ap_array_subdirectorios[int(index_carcar)] ) )
					graf_tempo =
						"<TR>\n" +
							"<TD BGCOLOR='#99ccff'>aptr" + strconv.Itoa(int(index_carcar)+1) + "</TD><TD port='" + strconv.Itoa(int(index_carcar)+1) + "'>" + strconv.Itoa(int(Avds.Avd_ap_array_subdirectorios[int(index_carcar)])) + "</TD>\n" +
							"</TR>\n"

					report_camino_file = report_camino_file + graf_tempo
					conexiones = idavl + ":" + strconv.Itoa(int(index_carcar)+1) + "->" + "AVD" + strconv.Itoa(int(pos_carpeta_bus)) + "\n"

					graf_tempo = "</TABLE>\n" +
						">];\n" +
						conexiones
					report_camino_file = report_camino_file + graf_tempo

				}

				//fmt.Println("Avds", "LA CARPETA YA EXISTE, ME VOY EN ELLA")
				/*recursivamente a la carpeta*/
				Leyendo_Recorrido(file, pos_carpeta_bus, var_directorios, i, sB, tipo_mk, tipo_accion, new_name, Contenido_ar)

			}
			//////******************FIN**CREACION*DE*CARPETAS*****************************************************////

			//////******************INICIO**CREACION*DE*ARCHIVOS*****************************************************////
		} else if tipo_insert == 'A' {
			///*buscando apuntador de archivos

			if Avds.Avd_ap_detalle_directorio == -1 {
				//fmt.Println("DetDir", "no exite DEtalle de directorio, entonces lo CREO")

				//////////////////////////fmt.Println("esfinal", esfinal )
				//if esfinal == false  {

				fmt.Println("Err, No existe Carepta (", var_directorios[i], ")")
				encontrado_tipo_art = false
				return
				//}

			} else {

				///fmt.Println("IIIIIIII si existe", "directorio reviso ahi")

				///*si ya existe detalle, entonces me busco la posicion del
				//detalle y voy insertando los inodos

				///*me posiciono en Det
				//luego busco espacio disponible
				cord_DD_exis := sB.Sb_ap_detalle_directorio_ini + (sB.Sb_size_struct_detalle_directorio * Avds.Avd_ap_detalle_directorio)
				//fmt.Println("existente DD cord_DD_exis", cord_DD_exis)
				//////pos_detdir = cord_DD_exis

				file.Seek(int64(cord_DD_exis), 0)

				detDir_exis := DetDirectorio{}
				var size int = int(binary.Size(detDir_exis))
				data := readBytesDisk(file, size)
				buffer := bytes.NewBuffer(data)

				err := binary.Read(buffer, binary.BigEndian, &detDir_exis)
				if err != nil {
					fmt.Println("binary. se encontro error al leer archivo binario", err)
				}
				/////fmt.Println("existente detDir_exis", detDir_exis)

				if tipo_accion == "REP_TF" {

					//idavl := "AVD"+ strconv.Itoa(int(avd_index))
					idavl := "AVD" + strconv.Itoa(int(pos_AVD))

					graf_tempo =
						idavl + "[label=<\n" +
							"<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n"

					report_camino_file = report_camino_file + graf_tempo

					graf_tempo = "<TR port='0'>\n" +
						"<TD COLSPAN='2'><font color='black'>" + NameByteToString(Avds.Avd_nombre_directorio) + "</font></TD>\n" +
						"</TR>\n" +

						"<TR>\n" +
						"<TD BGCOLOR='#99ccff'>Avd_fecha_creacion</TD><TD>" + string(Avds.Avd_fecha_creacion[:]) + "</TD>\n" +
						"</TR>\n"
					report_camino_file = report_camino_file + graf_tempo

					//pos_carpeta_bus := 	sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * int64( Avds.Avd_ap_array_subdirectorios[int(index_carcar)] ) )
					if Avds.Avd_ap_detalle_directorio != -1 {
						graf_tempo =
							"<TR>\n" +
								"<TD BGCOLOR='#99ccff'>Detalle D</TD><TD port='7'>" + strconv.Itoa(int(Avds.Avd_ap_detalle_directorio)) + "</TD>\n" +
								"</TR>\n"

						report_camino_file = report_camino_file + graf_tempo
						conexiones = idavl + ":7->" + "DETD" + strconv.Itoa(int(cord_DD_exis)) + "\n"
					}

					graf_tempo = "</TABLE>\n" +
						">];\n" +
						conexiones
					report_camino_file = report_camino_file + graf_tempo

				}

				/*BUSCANDO EL ARCHIVO EN LOS DETALLES*/
				//BuscandoDetalleDisponible(file, detDir_exis,sB, Contenido_ar, var_directorios[i], cord_DD_exis, var_directorios, i)
				BuscandoDetalle_Leer(file, detDir_exis, sB, var_directorios[i], cord_DD_exis, var_directorios, i, tipo_accion, new_name, Contenido_ar)

			}

		}
		//////******************FIN**CREACION*DE*ARCHIVOS*****************************************************////

	} else {

		//fmt.Println("Avds", "no exite")
		//fmt.Println("Avds.Avd_ap_detalle_directorio", Avds.Avd_ap_detalle_directorio)
		/*no existe*/
	}
}

//////////////////////////////////fin recorriendo leyendo

func BuscandoDetalleDisponible(file *os.File, detDir_exis DetDirectorio, sB SuperB, Contenido_ar string, name_archiv string, pos_DD int64, var_directorios []string, i int64) /*(DetDirectorio, int64, int64)*/ {

	encotrado_idet := false
	var index_array_DD int64

	file.Seek(int64(pos_DD), 0)

	detDir_exis = DetDirectorio{}
	var size int = int(binary.Size(detDir_exis))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &detDir_exis)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario", err)
	}

	for j := 0; j < len(detDir_exis.Dd_array_files); j++ {
		//////////////////////////fmt.Println(j,"AAAAAAAAAAAAAAAAAA", detDir_exis.Dd_array_files[j].Dd_file_ap_inodo )
		//////////////////fmt.Printf("	Detdir %s\n", detDir_exis.Dd_array_files[j].Dd_file_nombre)
		/*Aqui verifico si hay espacio en detalle de directorio*/
		if detDir_exis.Dd_array_files[j].Dd_file_ap_inodo == -1 {
			/*pos_carpeta_bus = sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * Avds.Avd_ap_array_subdirectorios[j] )
			encotrado = Existe_Carpeta(file, pos_carpeta_bus, var_directorios[i])*/

			//if encotrado_idet == true {
			index_array_DD = int64(j)
			encotrado_idet = true
			break
		}
	}

	//fmt.Println("encotrado_idet", encotrado_idet )
	//fmt.Println("index_array_DD", index_array_DD )

	if encotrado_idet == false {
		//////////////////////fmt.Println("3333333333333333333333 no se encontro libre normal encotrado_idet", encotrado_idet)
		////fmt.Println("33333333333333333333333 creo apuntador indirecto", encotrado_idet)

		/////////////////fmt.Println("detDir_exis.Dd_ap_detalle_directorio_indirec", detDir_exis.Dd_ap_detalle_directorio_indirec)
		if detDir_exis.Dd_ap_detalle_directorio_indirec == -1 {
			//new det directorio
			///////////////////////////////////////////////////////////////////
			/*buscando la pos segun en el bitmap*/
			index_bm_DD := Pos_bitmap_Detalles(file, sB)
			cord_DD := sB.Sb_ap_detalle_directorio_ini + (sB.Sb_size_struct_detalle_directorio * int64(index_bm_DD))
			//fmt.Println("indirec DD index_bm_DD", index_bm_DD)
			//fmt.Println("indirec DD cord_DD", cord_DD)

			//pos_detdir = cord_DD////////-----------------------------

			/////////////////////////////
			detDir_new := DetDirectorio{}
			index_array_DD = 0
			NewDetalleDirectorio(file, name_archiv, &detDir_new, cord_DD, int64(index_bm_DD), sB, index_array_DD)
			/*reviso primer index libre*/
			/*como es nuevo, entonces hago el primero del inodo*/
			/*actualizando Det indirecto*/

			detDir_exis.Dd_ap_detalle_directorio_indirec = int64(index_bm_DD)
			UpdateDD(file, detDir_exis, pos_DD)
			///////////////////////fmt.Println("detDir_exis.Dd_ap_detalle_directorio_indirec ", detDir_exis.Dd_ap_detalle_directorio_indirec )
			////////////////////////////////////////////////////////////////////

			//detDir = detDir_exis
			var index_array_inodo int64 = 0
			//detDir_exis, pos_detdir, index_array_DD = BuscandoDetalleDisponible(file, detDir_exis,sB, Contenido_ar, var_directorios[i], cord_DD_exis)
			Creando_First_Inodo(file, var_directorios, i, sB, Contenido_ar, index_array_DD, index_array_inodo, detDir_new, cord_DD)
			//return detDir_new, cord_DD, index_array_DD

		} else {
			/*SI YA EXISTE APUNTADOR INDURECTO, ENTONCES EM POSICIONO EN ELLA*/
			///fmt.Println("SI YA EXISTE APUNTADOR INDURECTO, ENTONCES EM POSICIONO EN ELLA")
			cord_DD_sig := sB.Sb_ap_detalle_directorio_ini + (sB.Sb_size_struct_detalle_directorio * detDir_exis.Dd_ap_detalle_directorio_indirec)
			BuscandoDetalleDisponible(file, detDir_exis, sB, Contenido_ar, name_archiv, cord_DD_sig, var_directorios, i)
			//fmt.Println("			DESPUES DE SER RECURSIVA")
		}

	} else {

		/*actualizando datos de archivo*/
		copy(detDir_exis.Dd_array_files[index_array_DD].Dd_file_nombre[:], name_archiv)
		fecha_string := time.Now().Format("2006-01-02 15:04:05")
		copy(detDir_exis.Dd_array_files[index_array_DD].Dd_file_date_creacion[:], fecha_string)
		detDir_exis.Dd_array_files[index_array_DD].Dd_file_date_modificacion = detDir_exis.Dd_array_files[index_array_DD].Dd_file_date_creacion

		//fmt.Println("			RETORNONU FILA DISPO", index_array_DD)

		//detDir = detDir_exis
		var index_array_inodo int64 = 0
		//detDir_exis, pos_detdir, index_array_DD = BuscandoDetalleDisponible(file, detDir_exis,sB, Contenido_ar, var_directorios[i], cord_DD_exis)
		Creando_First_Inodo(file, var_directorios, i, sB, Contenido_ar, index_array_DD, index_array_inodo, detDir_exis, pos_DD)
		//return detDir_exis, pos_DD, index_array_DD
	}

}

/*BUSCANDO DETALLE  ARCHIVO*/
func BuscandoDetalle_Leer(file *os.File, detDir_exis DetDirectorio, sB SuperB, name_archiv string, pos_DD int64, var_directorios []string, i int64, tipo_accion string, new_name string, Contenido_ar string) /*(DetDirectorio, int64, int64)*/ {

	encotrado_idet := false
	var index_array_DD int64
	var graf_tempo string = ""
	var conexiones string = ""

	file.Seek(int64(pos_DD), 0)

	detDir_exis = DetDirectorio{}
	var size int = int(binary.Size(detDir_exis))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &detDir_exis)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario", err)
	}

	for j := 0; j < len(detDir_exis.Dd_array_files); j++ {
		//fmt.Println(j,"AAAAAAAAAAAAAAAAAA", detDir_exis.Dd_array_files[j].Dd_file_ap_inodo )
		//fmt.Printf("	Detdir %s\n", detDir_exis.Dd_array_files[j].Dd_file_nombre)

		/*verifico si existe*/
		var name_actual [20]byte
		copy(name_actual[:], name_archiv)

		//fmt.Println("name_actual", name_actual)
		if detDir_exis.Dd_array_files[j].Dd_file_ap_inodo != -1 && name_actual == detDir_exis.Dd_array_files[j].Dd_file_nombre {
			//fmt.Println("---Avds", "son igualitos, exite")
			index_array_DD = int64(j)
			encotrado_idet = true
			break
		}
	}

	//fmt.Println("encotrado_idet", encotrado_idet )
	//fmt.Println("index_array_DD", index_array_DD )

	if encotrado_idet == false {
		//fmt.Println("3333333333333333333333 no se encontro libre normal encotrado_idet", encotrado_idet)

		//fmt.Println("detDir_exis.Dd_ap_detalle_directorio_indirec", detDir_exis.Dd_ap_detalle_directorio_indirec)
		if detDir_exis.Dd_ap_detalle_directorio_indirec == -1 {

			fmt.Println("Err, No existe Archivo (", var_directorios[i], ")")
			encontrado_tipo_art = false
			return

		} else {
			/*SI YA EXISTE APUNTADOR INDURECTO, ENTONCES EM POSICIONO EN ELLA*/
			cord_DD_sig := sB.Sb_ap_detalle_directorio_ini + (sB.Sb_size_struct_detalle_directorio * detDir_exis.Dd_ap_detalle_directorio_indirec)

			id_dd := (pos_DD - sB.Sb_ap_detalle_directorio_ini) / sB.Sb_size_struct_detalle_directorio

			if tipo_accion == "REP_TF" {

				iddet := "DETD" + strconv.Itoa(int(pos_DD))

				graf_tempo =
					iddet + "[label=<\n" +
						"<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n" +
						"<TR>\n" +
						"<TD BGCOLOR='#a3d977'>Detalle D</TD><TD BGCOLOR='#a3d977'>" + strconv.Itoa(int(id_dd)) + "</TD>\n" +
						"</TR>\n"

				report_camino_file = report_camino_file + graf_tempo

				graf_tempo =
					"<TR>\n" +
						"<TD BGCOLOR='#a3d977'>aptr_ind</TD><TD port='6'>" + strconv.Itoa(int(detDir_exis.Dd_ap_detalle_directorio_indirec)) + "</TD>\n" +
						"</TR>\n"

				report_camino_file = report_camino_file + graf_tempo

				conexiones = iddet + ":6->" + "DETD" + strconv.Itoa(int(cord_DD_sig)) + "\n"

				graf_tempo = "</TABLE>\n" +
					">];\n" +
					conexiones
				report_camino_file = report_camino_file + graf_tempo

			}

			BuscandoDetalle_Leer(file, detDir_exis, sB, name_archiv, cord_DD_sig, var_directorios, i, tipo_accion, new_name, Contenido_ar)
			//fmt.Println("			DESPUES DE SER RECURSIVA")
		}

	} else {

		//fmt.Println("-----------------",  detDir_exis.Dd_array_files[index_array_DD].Dd_file_nombre, "-----------------")

		if tipo_accion == "REP_TF" {

			id_dd := (pos_DD - sB.Sb_ap_detalle_directorio_ini) / sB.Sb_size_struct_detalle_directorio
			iddet := "DETD" + strconv.Itoa(int(pos_DD))
			pos_inod_sig := sB.Sb_ap_tabla_inodo_ini + (sB.Sb_size_struct_inodo * detDir_exis.Dd_array_files[index_array_DD].Dd_file_ap_inodo)
			graf_tempo =
				iddet + "[label=<\n" +
					"<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n" +
					"<TR>\n" +
					"<TD BGCOLOR='#a3d977'>Detalle D</TD><TD BGCOLOR='#a3d977'>" + strconv.Itoa(int(id_dd)) + "</TD>\n" +
					"</TR>\n"

			report_camino_file = report_camino_file + graf_tempo

			graf_tempo =
				"<TR>\n" +
					"<TD BGCOLOR='#a3d977'>" + NameByteToString(detDir_exis.Dd_array_files[index_array_DD].Dd_file_nombre) + "</TD><TD port='" + strconv.Itoa(int(index_array_DD)+1) + "'>" + strconv.Itoa(int(detDir_exis.Dd_array_files[index_array_DD].Dd_file_ap_inodo)) + "</TD>\n" +
					"</TR>\n"

			report_camino_file = report_camino_file + graf_tempo

			conexiones = iddet + ":" + strconv.Itoa(int(index_array_DD)+1) + "->" + "INOD" + strconv.Itoa(int(pos_inod_sig)) + "\n"

			graf_tempo = "</TABLE>\n" +
				">];\n" +
				conexiones
			report_camino_file = report_camino_file + graf_tempo

			Inodos_deArchivo(file, pos_inod_sig, sB, 0)

		}

		////////para imprimir bloque
		if tipo_accion == "PRI_AR" {

			fmt.Printf("-----------------%s-----------------\n", detDir_exis.Dd_array_files[index_array_DD].Dd_file_nombre)
			fmt.Printf("#")
			pos_inod_sig := sB.Sb_ap_tabla_inodo_ini + (sB.Sb_size_struct_inodo * detDir_exis.Dd_array_files[index_array_DD].Dd_file_ap_inodo)
			Print_Inodos(file, pos_inod_sig, sB)
			fmt.Println("")
		}
		encontrado_tipo_art = true

		////////modificar nombre de archivo
		if tipo_accion == "MOD_NAME" {
			//detDir_exis.Dd_array_files[index_array_DD].Dd_file_nombre = new_name
			copy(detDir_exis.Dd_array_files[index_array_DD].Dd_file_nombre[:], new_name)
			//update detdir
			UpdateDD(file, detDir_exis, pos_DD)
			fmt.Println("Se actualizó nombre con exito")

		}

		////////editando contenido
		if tipo_accion == "MO_ARCV" {
			/*borrando inodos y bloques para actulizar nuevo*/

			///name

			pos_inod_sig := sB.Sb_ap_tabla_inodo_ini + (sB.Sb_size_struct_inodo * detDir_exis.Dd_array_files[index_array_DD].Dd_file_ap_inodo)
			Borrando_Inodos(file, pos_inod_sig, sB)

			//UpdateInode(file, inod, pos_Inodo)
			Update_bitmapInodo(file, detDir_exis.Dd_array_files[index_array_DD].Dd_file_ap_inodo, sB, '0')
			/*borrando nodo*/

			//detDir_exis.Dd_array_files[index_array_DD].Dd_file_ap_inodo = -1
			//UpdateDD(file, detDir_exis, pos_DD)

			/*actualizando datos de archivo*/
			//copy(detDir_exis.Dd_array_files[index_array_DD].Dd_file_nombre[:], name_archiv)
			//fecha_actulizacion:= time.Now().Format("2006-01-02 15:04:05")
			////copy(detDir_exis.Dd_array_files[index_array_DD].Dd_file_date_creacion[:], fecha_actulizacion)
			//copy(detDir_exis.Dd_array_files[index_array_DD].Dd_file_date_modificacion[:], fecha_actulizacion)

			var index_array_inodo int64 = 0
			Creando_First_Inodo(file, var_directorios, i, sB, Contenido_ar, index_array_DD, index_array_inodo, detDir_exis, pos_DD)

		}

		/////////ME VOY A BUSCAR INODOS

	}

}

func RecorroDirectorio(path_disco string, inicia_part uint64) {

	file, err := os.OpenFile(path_disco, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(int64(inicia_part), 0)

	sB := SuperB{}
	var size int = int(binary.Size(sB))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &sB)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario", err)
	}

	report_tree_complet = "digraph arbolito_com {\n" +
		"node [shape=plaintext]\n" +
		"rankdir=LR;\n"

	//report_tree_complet = ""

	report_directorio = report_tree_complet

	RecorroCarpetas(file, sB.Sb_ap_arbol_directorio_ini, sB, 0)

}

func CreandoInodosconContenido(file *os.File, sB SuperB, Contenido_ar string, Inodo_new Inodo, pos_inod int64, aptr_b int64) {

	tempo_conte := Contenido_ar
	size_ar := int64(len(tempo_conte))
	var contenid_restante string
	var size_restante int64
	//fmt.Println("                                           Inodo_new.I_size_archivo", Inodo_new.I_size_archivo)
	//fmt.Println("++++++++++++++++++++++++++size_ar", size_ar)
	Inodo_new.I_count_bloques_asignados = aptr_b + 1
	if size_ar > 25 {
		tempo_conte = Contenido_ar[0:25]
		contenid_restante = Contenido_ar[25:]
		size_restante = int64(len(contenid_restante))

	}
	//fmt.Println("------CREATE--------------**tempo_conte", tempo_conte)
	//fmt.Println("--------------------size_restante", size_restante)
	//fmt.Println("--------------------contenid_restante", contenid_restante)

	/*creando inodos xd*/

	/***************INICIO*********creando bloque de comentarios*********************************/
	/*buscando la pos segun en el bitmap*/
	index_bm_bloque := Pos_bitmap_Bloques(file, sB)
	cord_bloque := sB.Sb_ap_bloques_ini + (sB.Sb_size_struct_bloque * int64(index_bm_bloque))
	//fmt.Println("inodo index_bm_bloque", index_bm_bloque)
	//fmt.Println("inodo cord_bloque", cord_bloque)

	var contenid_bloque [25]byte
	copy(contenid_bloque[:], tempo_conte)
	/*new bloque*/
	Bloque_new := BloqueDatos{}
	NewBloque(file, Bloque_new, cord_bloque, int64(index_bm_bloque), sB, contenid_bloque)

	//detDir.Dd_array_files[index_array_DD].Dd_file_ap_inodo = int64(index_bm_inodo)
	Inodo_new.I_array_bloques[aptr_b] = int64(index_bm_bloque)
	UpdateInode(file, Inodo_new, pos_inod)

	/***************FIN*********creando bloques de comentarios*********************************/

	/*si tiene mas, creo siguientes nodos*/
	if size_restante > 0 {
		aptr_b++
		//fmt.Println("							aptr_b", aptr_b)
		/*si es menor que 3 creo bloque, si es creo apuntador indirecto de inodo*/
		if aptr_b < 4 {
			CreandoInodosconContenido(file, sB, contenid_restante, Inodo_new, pos_inod, aptr_b)
		} else {
			/*creando nuevo inodo*/
			//Inodo_new.I_ap_indirecto
			if Inodo_new.I_ap_indirecto == -1 {
				CreandoInodoIndirecto(file, sB, contenid_restante, Inodo_new, pos_inod, aptr_b)
			}
		}

	}

}

////////creando primer INODO
func Creando_First_Inodo(file *os.File, var_directorios []string, i int64, sB SuperB, Contenido_ar string, index_array_DD int64, index_array_inodo int64, detDir DetDirectorio, pos_detdir int64) {

	/***************INICIO*********creando i-nodos*********************************/
	///////////fmt.Println("////////////////// index_array_DD", index_array_DD)
	/*buscando la pos segun en el bitmap*/
	index_bm_inodo := Pos_bitmap_Inodos(file, sB)
	cord_inodo := sB.Sb_ap_tabla_inodo_ini + (sB.Sb_size_struct_inodo * int64(index_bm_inodo))
	//fmt.Println("inodo index_bm_inodo", index_bm_inodo)
	//fmt.Println("inodo cord_inodo", cord_inodo)

	/*Contenido_ar = "123456789d123456789v123456789t123456789c123456789c123456789S123456789S123456789o123456789n123456789C"+
	"123456789d123456789v123456789t123456789c123456789c123456789S123456789S123456789o123456789n123456789C" +
	"123456789d"*/

	/*new inodo*/
	Inodo_new := Inodo{}
	NewInodo(file, var_directorios[i], &Inodo_new, cord_inodo, int64(index_bm_inodo), sB, index_array_inodo, Contenido_ar)

	detDir.Dd_array_files[index_array_DD].Dd_file_ap_inodo = int64(index_bm_inodo)
	//detDir.Dd_array_files[index_array_DD].Dd_file_ap_inodo = 123
	//fmt.Printf("detDir_new.Dd_array_files[index_array_DD].Dd_file_nombre  %s\n", detDir.Dd_array_files[index_array_DD].Dd_file_nombre)
	//fmt.Println("pos_detdir", pos_detdir)
	UpdateDD(file, detDir, pos_detdir)

	/*creando bloque en nodos*/
	//CreandoInodosconContenido(file, sB, Contenido_ar, Inodo_new, cord_inodo, 0)
	CreandoInodosconContenido(file, sB, Contenido_ar, Inodo_new, cord_inodo, 0)
	//CreandoInodosconContenido("123456789d123456789v123456789t123456789c123456789s123456789T")
	/***************FIN*********creando i-nodos*********************************/

}

////////creando new Inode de Apuntador Indirecto
func CreandoInodoIndirecto(file *os.File, sB SuperB, Contenido_ar string, Inodo_new Inodo, pos_inod int64, aptr_b int64) {

	/***************INICIO*********creando i-nodos*********************************/
	//////fmt.Println("							Creando inodo apuntador indirecto")
	/*buscando la pos segun en el bitmap*/
	index_bm_inodo := Pos_bitmap_Inodos(file, sB)
	cord_inodo := sB.Sb_ap_tabla_inodo_ini + (sB.Sb_size_struct_inodo * int64(index_bm_inodo))
	///////////////fmt.Println("i indirec index_bm_inodo", index_bm_inodo)
	///////////fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx   i indirec cord_inodo", cord_inodo)

	/*new inodo*/
	Inodo_aptr := Inodo{}
	NewInodo(file, " ", &Inodo_aptr, cord_inodo, int64(index_bm_inodo), sB, 0, Contenido_ar)

	Inodo_new.I_ap_indirecto = int64(index_bm_inodo)
	pos_inod_pad := sB.Sb_ap_tabla_inodo_ini + (sB.Sb_size_struct_inodo * Inodo_new.I_count_inodo)
	UpdateInode(file, Inodo_new, pos_inod_pad)
	//UpdateInode(file, Inodo_new, pos_inod)

	/*creando bloque en nodos*/
	//CreandoInodosconContenido(file, sB, Contenido_ar, Inodo_aptr, cord_inodo, 0)
	//Contenido_ar = "123456789d123456789v123456789t123456789c123456789c123456789S123456789S123456789o123456789n123456789C123456789M"
	CreandoInodosconContenido(file, sB, Contenido_ar, Inodo_aptr, cord_inodo, 0)
	//CreandoInodosconContenido("123456789d123456789v123456789t123456789c123456789s123456789T")

	/***************FIN*********creando i-nodos*********************************/
}

func CreandoInodosconContenido_back(Contenido_ar string) {

	/*tempo_conte := Contenido_ar
	size_ar := int64(len(tempo_conte))
	var contenid_restante string
	var size_restante int64
	fmt.Println("++++++++++++++++++++++++++size_ar", size_ar)
	if size_ar > 25 {
		tempo_conte = Contenido_ar[0:25]
		contenid_restante = Contenido_ar[25:len(Contenido_ar)]
		size_restante = int64(len(contenid_restante))

	}
	fmt.Println("------CREATE--------------**tempo_conte", tempo_conte)
	fmt.Println("--------------------size_restante", size_restante)
	fmt.Println("--------------------contenid_restante", contenid_restante)

	if size_restante > 0 {
		//CreandoInodosconContenido(contenid_restante)
	}*/

}

func RecorroCarpetas(file *os.File, pos_carpeta int64, sB SuperB, avd_index int64) {

	//////recorro carpetas
	var graf_tempo string = ""
	file.Seek(int64(pos_carpeta), 0)

	avd_dir := AVD{}
	//var size int  = int(binary.Size(sB))
	data := readBytesDisk(file, int(sB.Sb_size_struct_arbol_directorio))
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &avd_dir)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario AVD", err)
	}

	//(pos_carpeta - sB.Sb_ap_arbol_directorio_ini)/ sB.Sb_size_struct_arbol_directorio  :=  avd_dir.Avd_ap_array_subdirectorios[i]
	var conexiones string = ""
	var conexiones_direct string = ""
	idavl := "AVD" + strconv.Itoa(int(avd_index))

	graf_tempo =
		idavl + "[label=<\n" +
			"<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n"

	report_tree_complet = report_tree_complet + graf_tempo
	/*reporte directorio*/
	report_directorio = report_directorio + graf_tempo

	graf_tempo = "<TR port='0'>\n" +
		"<TD COLSPAN='2'><font color='black'>" + NameByteToString(avd_dir.Avd_nombre_directorio) + "</font></TD>\n" +
		"</TR>\n" +

		"<TR>\n" +
		"<TD BGCOLOR='#99ccff'>Avd_fecha_creacion</TD><TD>" + string(avd_dir.Avd_fecha_creacion[:]) + "</TD>\n" +
		"</TR>\n"

	report_tree_complet = report_tree_complet + graf_tempo
	/*reporte directorio*/
	report_directorio = report_directorio + graf_tempo

	for i := 0; i < int(len(avd_dir.Avd_ap_array_subdirectorios)); i++ {

		if avd_dir.Avd_ap_array_subdirectorios[i] != -1 {
			//pos_carpeta_sig := sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * avd_dir.Avd_ap_array_subdirectorios[i] )
			//RecorroCarpetas(file, pos_carpeta_sig, sB, avd_dir.Avd_ap_array_subdirectorios[i])

			graf_tempo =
				"<TR>\n" +
					"<TD BGCOLOR='#99ccff'>aptr" + strconv.Itoa(i+1) + "</TD><TD port='" + strconv.Itoa(i+1) + "'>" + strconv.Itoa(int(avd_dir.Avd_ap_array_subdirectorios[i])) + "</TD>\n" +
					"</TR>\n"

			report_tree_complet = report_tree_complet + graf_tempo
			conexiones = conexiones + idavl + ":" + strconv.Itoa(i+1) + "->" + "AVD" + strconv.Itoa(int(avd_dir.Avd_ap_array_subdirectorios[i])) + "\n"

			/*reporte directorio*/
			report_directorio = report_directorio + graf_tempo
			conexiones_direct = conexiones_direct + idavl + ":" + strconv.Itoa(i+1) + "->" + "AVD" + strconv.Itoa(int(avd_dir.Avd_ap_array_subdirectorios[i])) + "\n"
		} else {
			graf_tempo =
				"<TR>\n" +
					"<TD BGCOLOR='#99ccff'>aptr" + strconv.Itoa(i+1) + "</TD><TD port='" + strconv.Itoa(i+1) + "'> </TD>\n" +
					"</TR>\n"

			report_tree_complet = report_tree_complet + graf_tempo
			/*reporte directorio*/
			report_directorio = report_directorio + graf_tempo
		}
	}

	/*para deatalle*/
	if avd_dir.Avd_ap_detalle_directorio != -1 {
		//pos_archiv_sig := sB.Sb_ap_detalle_directorio_ini+ (sB.Sb_size_struct_detalle_directorio* avd_dir.Avd_ap_detalle_directorio )
		//RecorroDetalleDir(file, pos_archiv_sig, sB)
		graf_tempo =
			"<TR>\n" +
				"<TD BGCOLOR='#99ccff'>Detalle D</TD><TD port='7'>" + strconv.Itoa(int(avd_dir.Avd_ap_detalle_directorio)) + "</TD>\n" +
				"</TR>\n"

		report_tree_complet = report_tree_complet + graf_tempo
		conexiones = conexiones + idavl + ":7->" + "DETD" + strconv.Itoa(int(avd_dir.Avd_ap_detalle_directorio)) + "\n"
	} else {
		graf_tempo =
			"<TR>\n" +
				"<TD BGCOLOR='#99ccff'>Detalle D</TD><TD port='7'> </TD>\n" +
				"</TR>\n"

		report_tree_complet = report_tree_complet + graf_tempo
	}

	/*para recorrer apuntador indirecto*/
	if avd_dir.Avd_ap_arbol_virtual_directorio != -1 {
		//pos_carpeta_sig := sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * avd_dir.Avd_ap_array_subdirectorios[i] )
		//RecorroCarpetas(file, pos_carpeta_sig, sB, avd_dir.Avd_ap_array_subdirectorios[i])
		graf_tempo =
			"<TR>\n" +
				"<TD BGCOLOR='#99ccff'>aptr_ind</TD><TD port='8'>" + strconv.Itoa(int(avd_dir.Avd_ap_arbol_virtual_directorio)) + "</TD>\n" +
				"</TR>\n"

		report_tree_complet = report_tree_complet + graf_tempo
		conexiones = conexiones + idavl + ":8->" + "AVD" + strconv.Itoa(int(avd_dir.Avd_ap_arbol_virtual_directorio)) + "\n"

		/*reporte directorio*/
		report_directorio = report_directorio + graf_tempo
		conexiones_direct = conexiones_direct + idavl + ":8->" + "AVD" + strconv.Itoa(int(avd_dir.Avd_ap_arbol_virtual_directorio)) + "\n"

	} else {
		graf_tempo =
			"<TR>\n" +
				"<TD BGCOLOR='#99ccff'>aptr_ind</TD><TD port='8'> </TD>\n" +
				"</TR>\n"
		report_tree_complet = report_tree_complet + graf_tempo

		/*reporte directorio*/
		report_directorio = report_directorio + graf_tempo
	}

	graf_tempo = "</TABLE>\n" +
		">];\n" +
		conexiones

	report_tree_complet = report_tree_complet + graf_tempo

	/*reporte directorio*/
	report_directorio = report_directorio + "</TABLE>\n" +
		">];\n" +
		conexiones_direct

	//AVD1:7 -> DD1

	fmt.Printf("car %s\n", avd_dir.Avd_nombre_directorio)
	for i := 0; i < int(len(avd_dir.Avd_ap_array_subdirectorios)); i++ {

		//fmt.Println(i,"----", avd_dir.Avd_ap_array_subdirectorios[i])
		if avd_dir.Avd_ap_array_subdirectorios[i] != -1 {
			pos_carpeta_sig := sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * avd_dir.Avd_ap_array_subdirectorios[i])
			RecorroCarpetas(file, pos_carpeta_sig, sB, avd_dir.Avd_ap_array_subdirectorios[i])
		}
	}

	/*para recorrer apuntador indirecto*/
	if avd_dir.Avd_ap_arbol_virtual_directorio != -1 {
		pos_carpeta_indi := sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * avd_dir.Avd_ap_arbol_virtual_directorio)
		RecorroCarpetas(file, pos_carpeta_indi, sB, avd_dir.Avd_ap_arbol_virtual_directorio)
	}

	/*para recorrer archivos*/
	if avd_dir.Avd_ap_detalle_directorio != -1 {
		pos_archiv_sig := sB.Sb_ap_detalle_directorio_ini + (sB.Sb_size_struct_detalle_directorio * avd_dir.Avd_ap_detalle_directorio)
		RecorroDetalleDir(file, pos_archiv_sig, sB, avd_dir.Avd_ap_detalle_directorio)
	}

}

func Recorro_Carpetas_TreeDir(file *os.File, pos_carpeta int64, sB SuperB, avd_index int64) {

	//////recorro carpetas
	var graf_tempo string = ""
	file.Seek(int64(pos_carpeta), 0)

	avd_dir := AVD{}
	//var size int  = int(binary.Size(sB))
	data := readBytesDisk(file, int(sB.Sb_size_struct_arbol_directorio))
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &avd_dir)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario AVD", err)
	}

	//(pos_carpeta - sB.Sb_ap_arbol_directorio_ini)/ sB.Sb_size_struct_arbol_directorio  :=  avd_dir.Avd_ap_array_subdirectorios[i]
	var conexiones string = ""
	idavl := "AVD" + strconv.Itoa(int(avd_index))

	graf_tempo =
		idavl + "[label=<\n" +
			"<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n"

	report_archiv_incar = report_archiv_incar + graf_tempo

	graf_tempo = "<TR port='0'>\n" +
		"<TD COLSPAN='2'><font color='black'>" + NameByteToString(avd_dir.Avd_nombre_directorio) + "</font></TD>\n" +
		"</TR>\n" +

		"<TR>\n" +
		"<TD BGCOLOR='#99ccff'>Avd_fecha_creacion</TD><TD>" + string(avd_dir.Avd_fecha_creacion[:]) + "</TD>\n" +
		"</TR>\n"

	report_archiv_incar = report_archiv_incar + graf_tempo

	/*for i := 0; i < int(len(avd_dir.Avd_ap_array_subdirectorios)); i++ {

		if avd_dir.Avd_ap_array_subdirectorios[i] != -1 {
			//pos_carpeta_sig := sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * avd_dir.Avd_ap_array_subdirectorios[i] )
			//RecorroCarpetas(file, pos_carpeta_sig, sB, avd_dir.Avd_ap_array_subdirectorios[i])

			graf_tempo =
			"<TR>\n"+
			"<TD BGCOLOR='#99ccff'>aptr"+strconv.Itoa(i+1)+"</TD><TD port='"+strconv.Itoa(i+1)+"'>"+strconv.Itoa(int(avd_dir.Avd_ap_array_subdirectorios[i]))+"</TD>\n"+
			"</TR>\n"


			report_tree_complet =  report_tree_complet + graf_tempo
			conexiones = conexiones + idavl+":" +strconv.Itoa(i+1)+"->" +  "AVD"+ strconv.Itoa(int(avd_dir.Avd_ap_array_subdirectorios[i])) +"\n"

			///*reporte directorio
			report_directorio = report_directorio + graf_tempo
			conexiones_direct = conexiones_direct + idavl+":" +strconv.Itoa(i+1)+"->" +  "AVD"+ strconv.Itoa(int(avd_dir.Avd_ap_array_subdirectorios[i])) +"\n"
		} else {
			graf_tempo =
			"<TR>\n"+
			"<TD BGCOLOR='#99ccff'>aptr"+strconv.Itoa(i+1)+"</TD><TD port='"+strconv.Itoa(i+1)+"'> </TD>\n"+
			"</TR>\n"

			report_tree_complet =  report_tree_complet + graf_tempo
			///*reporte directorio
			report_directorio = report_directorio + graf_tempo
		}
	}*/

	/*para deatalle*/
	if avd_dir.Avd_ap_detalle_directorio != -1 {
		//pos_archiv_sig := sB.Sb_ap_detalle_directorio_ini+ (sB.Sb_size_struct_detalle_directorio* avd_dir.Avd_ap_detalle_directorio )
		//RecorroDetalleDir(file, pos_archiv_sig, sB)
		graf_tempo =
			"<TR>\n" +
				"<TD BGCOLOR='#99ccff'>Detalle D</TD><TD port='7'>" + strconv.Itoa(int(avd_dir.Avd_ap_detalle_directorio)) + "</TD>\n" +
				"</TR>\n"

		report_archiv_incar = report_archiv_incar + graf_tempo
		conexiones = conexiones + idavl + ":7->" + "DETD" + strconv.Itoa(int(avd_dir.Avd_ap_detalle_directorio)) + "\n"
	} else {
		graf_tempo =
			"<TR>\n" +
				"<TD BGCOLOR='#99ccff'>Detalle D</TD><TD port='7'> </TD>\n" +
				"</TR>\n"

		report_archiv_incar = report_archiv_incar + graf_tempo
	}

	/*para recorrer apuntador indirecto*/
	/*if avd_dir.Avd_ap_arbol_virtual_directorio != -1 {
		//pos_carpeta_sig := sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * avd_dir.Avd_ap_array_subdirectorios[i] )
		//RecorroCarpetas(file, pos_carpeta_sig, sB, avd_dir.Avd_ap_array_subdirectorios[i])
		graf_tempo =
		"<TR>\n"+
		"<TD BGCOLOR='#99ccff'>aptr_ind</TD><TD port='8'>"+strconv.Itoa(int(avd_dir.Avd_ap_arbol_virtual_directorio))+"</TD>\n"+
		"</TR>\n"

		report_tree_complet =  report_tree_complet + graf_tempo
		conexiones = conexiones + idavl+":8->" +  "AVD"+ strconv.Itoa(int(avd_dir.Avd_ap_arbol_virtual_directorio  )) +"\n"

		///*reporte directorio
		report_directorio = report_directorio + graf_tempo
		conexiones_direct = conexiones_direct + idavl+":8->" +  "AVD"+ strconv.Itoa(int(avd_dir.Avd_ap_arbol_virtual_directorio  )) +"\n"

	} else {
		graf_tempo =
		"<TR>\n"+
		"<TD BGCOLOR='#99ccff'>aptr_ind</TD><TD port='8'> </TD>\n"+
		"</TR>\n"
		report_tree_complet =  report_tree_complet + graf_tempo

		///*reporte directorio
		report_directorio = report_directorio + graf_tempo
	}
	*/

	graf_tempo = "</TABLE>\n" +
		">];\n" +
		conexiones

	report_archiv_incar = report_archiv_incar + graf_tempo

	//AVD1:7 -> DD1

	/*fmt.Printf("car %s\n", avd_dir.Avd_nombre_directorio)
	for i := 0; i < int(len(avd_dir.Avd_ap_array_subdirectorios)); i++ {
		//fmt.Println(i,"----", avd_dir.Avd_ap_array_subdirectorios[i])
		if avd_dir.Avd_ap_array_subdirectorios[i] != -1 {
			pos_carpeta_sig := sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * avd_dir.Avd_ap_array_subdirectorios[i] )
			RecorroCarpetas(file, pos_carpeta_sig, sB, avd_dir.Avd_ap_array_subdirectorios[i])
		}
	}*/

	/*para recorrer apuntador indirecto*/
	/*if avd_dir.Avd_ap_arbol_virtual_directorio != -1 {
		pos_carpeta_indi := sB.Sb_ap_arbol_directorio_ini + (sB.Sb_size_struct_arbol_directorio * avd_dir.Avd_ap_arbol_virtual_directorio )
		RecorroCarpetas(file, pos_carpeta_indi, sB, avd_dir.Avd_ap_arbol_virtual_directorio )
	}*/

	/*para recorrer archivos*/
	if avd_dir.Avd_ap_detalle_directorio != -1 {
		pos_archiv_sig := sB.Sb_ap_detalle_directorio_ini + (sB.Sb_size_struct_detalle_directorio * avd_dir.Avd_ap_detalle_directorio)
		RecorroDetalleDir(file, pos_archiv_sig, sB, avd_dir.Avd_ap_detalle_directorio)
	}

}

func RecorroDetalleDir(file *os.File, pos_DetDir int64, sB SuperB, index_d int64) {

	//////recorro carpetas
	var graf_tempo string = ""
	file.Seek(int64(pos_DetDir), 0)

	detdir := DetDirectorio{}
	//var size int  = int(binary.Size(sB))
	data := readBytesDisk(file, int(sB.Sb_size_struct_detalle_directorio))
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &detdir)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario DD", err)
	}

	var conexiones string = ""
	var conexiones_det_indi string = ""
	iddet := "DETD" + strconv.Itoa(int(index_d))

	graf_tempo =
		iddet + "[label=<\n" +
			"<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n" +
			"<TR>\n" +
			"<TD BGCOLOR='#a3d977'>Detalle D</TD><TD BGCOLOR='#a3d977'>" + strconv.Itoa(int(index_d)) + "</TD>\n" +
			"</TR>\n"

	report_tree_complet = report_tree_complet + graf_tempo
	report_archiv_incar = report_archiv_incar + graf_tempo

	for i := 0; i < int(len(detdir.Dd_array_files)); i++ {

		if detdir.Dd_array_files[i].Dd_file_ap_inodo != -1 {
			//fmt.Printf("**************************	Detdir %s\n", detdir.Dd_array_files[i].Dd_file_nombre)
			//fmt.Println("	detdir.Dd_array_files[",i,"].Dd_file_ap_inodo", detdir.Dd_array_files[i].Dd_file_ap_inodo)

			//pos_inod_sig := sB.Sb_ap_tabla_inodo_ini+ (sB.Sb_size_struct_inodo* detdir.Dd_array_files[i].Dd_file_ap_inodo )
			//RecorroInodos(file, pos_inod_sig, sB)

			graf_tempo =
				"<TR>\n" +
					"<TD BGCOLOR='#a3d977'>" + NameByteToString(detdir.Dd_array_files[i].Dd_file_nombre) + "</TD><TD port='" + strconv.Itoa(i+1) + "'>" + strconv.Itoa(int(detdir.Dd_array_files[i].Dd_file_ap_inodo)) + "</TD>\n" +
					"</TR>\n"

			report_tree_complet = report_tree_complet + graf_tempo
			report_archiv_incar = report_archiv_incar + graf_tempo

			conexiones = conexiones + iddet + ":" + strconv.Itoa(i+1) + "->" + "INOD" + strconv.Itoa(int(detdir.Dd_array_files[i].Dd_file_ap_inodo)) + "\n"
		} else {
			graf_tempo =
				"<TR>\n" +
					"<TD BGCOLOR='#a3d977'> </TD><TD port='" + strconv.Itoa(i+1) + "'> </TD>\n" +
					"</TR>\n"

			report_tree_complet = report_tree_complet + graf_tempo
			report_archiv_incar = report_archiv_incar + graf_tempo
		}
	}

	/*recorro indirecto*/
	if detdir.Dd_ap_detalle_directorio_indirec != -1 {
		//fmt.Println("+++++++++++++++ detdir.Dd_ap_detalle_directorio_indirec", detdir.Dd_ap_detalle_directorio_indirec)
		//pos_det_sig := sB.Sb_ap_detalle_directorio_ini+ (sB.Sb_size_struct_detalle_directorio * detdir.Dd_ap_detalle_directorio_indirec )
		//RecorroDetalleDir(file, pos_det_sig, sB, detdir.Dd_ap_detalle_directorio_indirec)

		graf_tempo =
			"<TR>\n" +
				"<TD BGCOLOR='#a3d977'>aptr_ind</TD><TD port='6'>" + strconv.Itoa(int(detdir.Dd_ap_detalle_directorio_indirec)) + "</TD>\n" +
				"</TR>\n"

		report_tree_complet = report_tree_complet + graf_tempo
		report_archiv_incar = report_archiv_incar + graf_tempo

		conexiones = conexiones + iddet + ":6->" + "DETD" + strconv.Itoa(int(detdir.Dd_ap_detalle_directorio_indirec)) + "\n"
		conexiones_det_indi = conexiones_det_indi + iddet + ":6->" + "DETD" + strconv.Itoa(int(detdir.Dd_ap_detalle_directorio_indirec)) + "\n"

	} else {
		graf_tempo =
			"<TR>\n" +
				"<TD BGCOLOR='#a3d977'>aptr_ind</TD><TD port='6'> </TD>\n" +
				"</TR>\n"

		report_tree_complet = report_tree_complet + graf_tempo
		report_archiv_incar = report_archiv_incar + graf_tempo
	}

	report_tree_complet = report_tree_complet + "</TABLE>\n" +
		">];\n" +
		conexiones

	report_archiv_incar = report_archiv_incar + "</TABLE>\n" +
		">];\n" +
		conexiones_det_indi

	for i := 0; i < int(len(detdir.Dd_array_files)); i++ {

		if detdir.Dd_array_files[i].Dd_file_ap_inodo != -1 {
			//fmt.Printf("**************************	Detdir %s\n", detdir.Dd_array_files[i].Dd_file_nombre)
			//fmt.Println("	detdir.Dd_array_files[", i, "].Dd_file_ap_inodo", detdir.Dd_array_files[i].Dd_file_ap_inodo)

			pos_inod_sig := sB.Sb_ap_tabla_inodo_ini + (sB.Sb_size_struct_inodo * detdir.Dd_array_files[i].Dd_file_ap_inodo)
			RecorroInodos(file, pos_inod_sig, sB, detdir.Dd_array_files[i].Dd_file_ap_inodo)

		}
	}

	/*recorro indirecto*/
	if detdir.Dd_ap_detalle_directorio_indirec != -1 {
		//fmt.Println("+++++++++++++++ detdir.Dd_ap_detalle_directorio_indirec", detdir.Dd_ap_detalle_directorio_indirec)
		pos_det_sig := sB.Sb_ap_detalle_directorio_ini + (sB.Sb_size_struct_detalle_directorio * detdir.Dd_ap_detalle_directorio_indirec)
		RecorroDetalleDir(file, pos_det_sig, sB, detdir.Dd_ap_detalle_directorio_indirec)

	}

}

//////////////recorriendo nodo
func RecorroInodos(file *os.File, pos_Inodo int64, sB SuperB, index_ino int64) {

	//////recorro carpetas
	file.Seek(int64(pos_Inodo), 0)

	inod := Inodo{}
	//var size int  = int(binary.Size(sB))
	data := readBytesDisk(file, int(sB.Sb_size_struct_inodo))
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &inod)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario DD", err)
	}

	//fmt.Println("inod.I_count_inodo", inod.I_count_inodo)
	//fmt.Println("inod.I_size_archivo", inod.I_size_archivo)
	//fmt.Println("	inod.I_count_bloques_asignados", inod.I_count_bloques_asignados)

	var conexiones string = ""
	idinod := "INOD" + strconv.Itoa(int(index_ino))

	report_tree_complet = report_tree_complet +
		idinod + "[label=<\n" +
		"<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n" +
		"<TR>\n" +
		"<TD BGCOLOR='#ffc374'>i-nodo</TD><TD BGCOLOR='#ffc374'>" + strconv.Itoa(int(inod.I_count_inodo)) + "</TD>\n" +
		"</TR>\n" +
		"<TR>\n" +
		"<TD BGCOLOR='#ffc374'>Size</TD><TD>" + strconv.Itoa(int(inod.I_size_archivo)) + "</TD>\n" +
		"</TR>\n" +
		"<TR>\n" +
		"<TD BGCOLOR='#ffc374'>Bloques</TD><TD>" + strconv.Itoa(int(inod.I_count_bloques_asignados)) + "</TD>\n" +
		"</TR>\n"

	for i := 0; i < int(len(inod.I_array_bloques)); i++ {

		if inod.I_array_bloques[i] != -1 {
			//fmt.Println("inod.I_array_bloques[",i,"]", inod.I_array_bloques[i])
			report_tree_complet = report_tree_complet +
				"<TR>\n" +
				"<TD BGCOLOR='#ffc374'>aptr" + strconv.Itoa(i+1) + "</TD><TD port='" + strconv.Itoa(i+1) + "'>" + strconv.Itoa(int(inod.I_array_bloques[i])) + "</TD>\n" +
				"</TR>\n"

			conexiones = conexiones + idinod + ":" + strconv.Itoa(i+1) + "->" + "BLO" + strconv.Itoa(int(inod.I_array_bloques[i])) + "\n"
		} else {
			report_tree_complet = report_tree_complet +
				"<TR>\n" +
				"<TD BGCOLOR='#ffc374'>aptr" + strconv.Itoa(i+1) + "</TD><TD port='" + strconv.Itoa(i+1) + "'> </TD>\n" +
				"</TR>\n"
		}
	}

	/*apuntador indirecto*/
	if inod.I_ap_indirecto != -1 {
		report_tree_complet = report_tree_complet +
			"<TR>\n" +
			"<TD BGCOLOR='#ffc374'>aptr_ind</TD><TD port='5'>" + strconv.Itoa(int(inod.I_ap_indirecto)) + "</TD>\n" +
			"</TR>\n"

		conexiones = conexiones + idinod + ":5->" + "INOD" + strconv.Itoa(int(inod.I_ap_indirecto)) + "\n"
	} else {
		report_tree_complet = report_tree_complet +
			"<TR>\n" +
			"<TD BGCOLOR='#ffc374'>aptr_ind</TD><TD port='5'> </TD>\n" +
			"</TR>\n"
	}

	report_tree_complet = report_tree_complet + "</TABLE>\n" +
		">];\n" +
		conexiones

	for i := 0; i < int(len(inod.I_array_bloques)); i++ {

		if inod.I_array_bloques[i] != -1 {
			fmt.Println("inod.I_array_bloques[", i, "]", inod.I_array_bloques[i])

			pos_bloq_sig := sB.Sb_ap_bloques_ini + (sB.Sb_size_struct_bloque * inod.I_array_bloques[i])
			RecorroBloques(file, pos_bloq_sig, sB, inod.I_array_bloques[i])
		}
	}

	/*apuntador indirecto*/
	if inod.I_ap_indirecto != -1 {
		fmt.Println("--inod.I_ap_indirecto", inod.I_ap_indirecto)
		pos_inod_in_sig := sB.Sb_ap_tabla_inodo_ini + (sB.Sb_size_struct_inodo * inod.I_ap_indirecto)
		RecorroInodos(file, pos_inod_in_sig, sB, inod.I_ap_indirecto)
	}
}

//////////////inidos_tree file
func Inodos_deArchivo(file *os.File, pos_Inodo int64, sB SuperB, index_ino int64) {

	//////recorro carpetas
	file.Seek(int64(pos_Inodo), 0)

	inod := Inodo{}
	//var size int  = int(binary.Size(sB))
	data := readBytesDisk(file, int(sB.Sb_size_struct_inodo))
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &inod)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario DD", err)
	}

	/*fmt.Println("inod.I_count_inodo", inod.I_count_inodo)
	fmt.Println("inod.I_size_archivo", inod.I_size_archivo)
	fmt.Println("	inod.I_count_bloques_asignados", inod.I_count_bloques_asignados)*/

	var conexiones string = ""
	idinod := "INOD" + strconv.Itoa(int(pos_Inodo))

	report_camino_file = report_camino_file +
		idinod + "[label=<\n" +
		"<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n" +
		"<TR>\n" +
		"<TD BGCOLOR='#ffc374'>i-nodo</TD><TD BGCOLOR='#ffc374'>" + strconv.Itoa(int(inod.I_count_inodo)) + "</TD>\n" +
		"</TR>\n" +
		"<TR>\n" +
		"<TD BGCOLOR='#ffc374'>Size</TD><TD>" + strconv.Itoa(int(inod.I_size_archivo)) + "</TD>\n" +
		"</TR>\n" +
		"<TR>\n" +
		"<TD BGCOLOR='#ffc374'>Bloques</TD><TD>" + strconv.Itoa(int(inod.I_count_bloques_asignados)) + "</TD>\n" +
		"</TR>\n"

	for i := 0; i < int(len(inod.I_array_bloques)); i++ {

		if inod.I_array_bloques[i] != -1 {
			//fmt.Println("inod.I_array_bloques[",i,"]", inod.I_array_bloques[i])
			report_camino_file = report_camino_file +
				"<TR>\n" +
				"<TD BGCOLOR='#ffc374'>aptr" + strconv.Itoa(i+1) + "</TD><TD port='" + strconv.Itoa(i+1) + "'>" + strconv.Itoa(int(inod.I_array_bloques[i])) + "</TD>\n" +
				"</TR>\n"

			conexiones = conexiones + idinod + ":" + strconv.Itoa(i+1) + "->" + "BLO" + strconv.Itoa(int(inod.I_array_bloques[i])) + "\n"
		} else {
			report_camino_file = report_camino_file +
				"<TR>\n" +
				"<TD BGCOLOR='#ffc374'>aptr" + strconv.Itoa(i+1) + "</TD><TD port='" + strconv.Itoa(i+1) + "'> </TD>\n" +
				"</TR>\n"
		}
	}

	/*apuntador indirecto*/
	if inod.I_ap_indirecto != -1 {

		pos_inod_in_sig := sB.Sb_ap_tabla_inodo_ini + (sB.Sb_size_struct_inodo * inod.I_ap_indirecto)

		report_camino_file = report_camino_file +
			"<TR>\n" +
			"<TD BGCOLOR='#ffc374'>aptr_ind</TD><TD port='5'>" + strconv.Itoa(int(inod.I_ap_indirecto)) + "</TD>\n" +
			"</TR>\n"

		conexiones = conexiones + idinod + ":5->" + "INOD" + strconv.Itoa(int(pos_inod_in_sig)) + "\n"
	} else {
		report_camino_file = report_camino_file +
			"<TR>\n" +
			"<TD BGCOLOR='#ffc374'>aptr_ind</TD><TD port='5'> </TD>\n" +
			"</TR>\n"
	}

	report_camino_file = report_camino_file + "</TABLE>\n" +
		">];\n" +
		conexiones

	for i := 0; i < int(len(inod.I_array_bloques)); i++ {

		if inod.I_array_bloques[i] != -1 {

			pos_bloq_sig := sB.Sb_ap_bloques_ini + (sB.Sb_size_struct_bloque * inod.I_array_bloques[i])
			Camino_ar_Bloques(file, pos_bloq_sig, sB, inod.I_array_bloques[i])
		}
	}

	/*apuntador indirecto*/
	if inod.I_ap_indirecto != -1 {
		pos_inod_in_sig := sB.Sb_ap_tabla_inodo_ini + (sB.Sb_size_struct_inodo * inod.I_ap_indirecto)
		Inodos_deArchivo(file, pos_inod_in_sig, sB, inod.I_ap_indirecto)
	}
}

//////////////imprimiendo nodos nodo
func Print_Inodos(file *os.File, pos_Inodo int64, sB SuperB) {

	//////recorro carpetas
	file.Seek(int64(pos_Inodo), 0)

	inod := Inodo{}
	//var size int  = int(binary.Size(sB))
	data := readBytesDisk(file, int(sB.Sb_size_struct_inodo))
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &inod)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario DD", err)
	}

	//fmt.Println("inod.I_count_inodo", inod.I_count_inodo)
	//fmt.Println("inod.I_size_archivo", inod.I_size_archivo)
	//fmt.Println("	inod.I_count_bloques_asignados", inod.I_count_bloques_asignados)

	for i := 0; i < int(len(inod.I_array_bloques)); i++ {

		if inod.I_array_bloques[i] != -1 {
			//fmt.Println("inod.I_array_bloques[",i,"]", inod.I_array_bloques[i])
			pos_bloq_sig := sB.Sb_ap_bloques_ini + (sB.Sb_size_struct_bloque * inod.I_array_bloques[i])
			Print_Bloques(file, pos_bloq_sig, sB)
		}
	}

	/*apuntador indirecto*/
	if inod.I_ap_indirecto != -1 {
		//fmt.Println("--inod.I_ap_indirecto", inod.I_ap_indirecto)
		pos_inod_in_sig := sB.Sb_ap_tabla_inodo_ini + (sB.Sb_size_struct_inodo * inod.I_ap_indirecto)
		Print_Inodos(file, pos_inod_in_sig, sB)
	}

}

//////////////borrando inodos
func Borrando_Inodos(file *os.File, pos_Inodo int64, sB SuperB) {

	//////recorro carpetas
	file.Seek(int64(pos_Inodo), 0)

	inod := Inodo{}
	//var size int  = int(binary.Size(sB))
	data := readBytesDisk(file, int(sB.Sb_size_struct_inodo))
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &inod)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario DD", err)
	}

	//fmt.Println("inod.I_count_inodo", inod.I_count_inodo)
	//fmt.Println("inod.I_size_archivo", inod.I_size_archivo)
	//fmt.Println("	inod.I_count_bloques_asignados", inod.I_count_bloques_asignados)

	for i := 0; i < int(len(inod.I_array_bloques)); i++ {

		if inod.I_array_bloques[i] != -1 {
			//fmt.Println("inod.I_array_bloques[",i,"]", inod.I_array_bloques[i])

			/////pos_bloq_sig := sB.Sb_ap_bloques_ini+ (sB.Sb_size_struct_bloque * inod.I_array_bloques[i] )

			/*borrando bitmap bloque*/
			BorrandoBloque_bitmap(file, inod.I_array_bloques[i], sB, '0')
			//borrando inodo
			inod.I_array_bloques[i] = -1
			UpdateInode(file, inod, pos_Inodo)

			////////////Print_Bloques(file, pos_bloq_sig, sB)
		}
	}

	/*apuntador indirecto*/
	if inod.I_ap_indirecto != -1 {
		//fmt.Println("--inod.I_ap_indirecto", inod.I_ap_indirecto)
		pos_inod_in_sig := sB.Sb_ap_tabla_inodo_ini + (sB.Sb_size_struct_inodo * inod.I_ap_indirecto)

		Update_bitmapInodo(file, inod.I_ap_indirecto, sB, '0')
		inod.I_ap_indirecto = -1
		UpdateInode(file, inod, pos_Inodo)
		Borrando_Inodos(file, pos_inod_in_sig, sB)
	}

}

//////////////imprimiendo bloques
func Print_Bloques(file *os.File, pos_Bloque int64, sB SuperB) {

	file.Seek(int64(pos_Bloque), 0)

	bloq := BloqueDatos{}
	//var size int  = int(binary.Size(sB))
	data := readBytesDisk(file, int(sB.Sb_size_struct_bloque))
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &bloq)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario DD", err)
	}

	fmt.Printf("%s", bloq.Db_data)

}

//////////////borrando bloques
func Borrando_Bloques(file *os.File, pos_Bloque int64, sB SuperB) {

	file.Seek(int64(pos_Bloque), 0)

	bloq := BloqueDatos{}
	//var size int  = int(binary.Size(sB))
	data := readBytesDisk(file, int(sB.Sb_size_struct_bloque))
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &bloq)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario DD", err)
	}

	fmt.Printf("%s", bloq.Db_data)

}

//////////////recorriendo bloques
func RecorroBloques(file *os.File, pos_Bloque int64, sB SuperB, index_bloq int64) {

	//////recorro carpetas
	file.Seek(int64(pos_Bloque), 0)

	bloq := BloqueDatos{}
	//var size int  = int(binary.Size(sB))
	data := readBytesDisk(file, int(sB.Sb_size_struct_bloque))
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &bloq)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario DD", err)
	}

	//fmt.Printf("		bloq.Db_data  %s\n", bloq.Db_data)

	idbloq := "BLO" + strconv.Itoa(int(index_bloq))

	report_tree_complet = report_tree_complet +
		idbloq + "[label=<\n" +
		"<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n" +

		"<TR>\n" +
		"<TD BGCOLOR='#ff8f80'>bloque</TD><TD BGCOLOR='#ff8f80'>" + strconv.Itoa(int(index_bloq)) + "</TD>\n" +
		"</TR>\n" +

		"<TR>\n" +
		"<TD COLSPAN='2'>" + BloqueByteToString(bloq.Db_data) + "</TD>\n" +
		"</TR>\n"

	report_tree_complet = report_tree_complet + "</TABLE>\n" +
		">];\n"
}

//////////////rbloque camino del archivo
func Camino_ar_Bloques(file *os.File, pos_Bloque int64, sB SuperB, index_bloq int64) {

	//////recorro carpetas
	file.Seek(int64(pos_Bloque), 0)

	bloq := BloqueDatos{}
	//var size int  = int(binary.Size(sB))
	data := readBytesDisk(file, int(sB.Sb_size_struct_bloque))
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &bloq)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario DD", err)
	}

	//fmt.Printf("		bloq.Db_data  %s\n", bloq.Db_data)

	idbloq := "BLO" + strconv.Itoa(int(index_bloq))

	report_camino_file = report_camino_file +
		idbloq + "[label=<\n" +
		"<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n" +

		"<TR>\n" +
		"<TD BGCOLOR='#ff8f80'>bloque</TD><TD BGCOLOR='#ff8f80'>" + strconv.Itoa(int(index_bloq)) + "</TD>\n" +
		"</TR>\n" +

		"<TR>\n" +
		"<TD COLSPAN='2'>" + BloqueByteToString(bloq.Db_data) + "</TD>\n" +
		"</TR>\n"

	report_camino_file = report_camino_file + "</TABLE>\n" +
		">];\n"

}

//////Bloques
func NewBloque(file *os.File, bloque_new BloqueDatos, posicion_bloq int64, index_bm_bloq int64, sB SuperB, Contenido_ar [25]byte) {

	//bloque_new.bloque_new = -1
	bloque_new.Db_data = Contenido_ar

	file.Seek(posicion_bloq, 0)
	bloq := &bloque_new

	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, bloq)
	msg_exito := "Se Creó bloque " + /*string(name_byte[:])  + */ " con Exito"
	writeBytesInDisk(file, binario.Bytes(), msg_exito)

	/////////////////////////////escribiendo en bitmap
	var uno byte = '1'
	epos_bitmap_dd := sB.Sb_ap_bitmap_bloques_ini + index_bm_bloq
	file.Seek(epos_bitmap_dd, 0)
	bitmal := &uno

	var binario_bmap bytes.Buffer
	binary.Write(&binario_bmap, binary.BigEndian, bitmal)
	msg_exito = "Se Creó Actualizo bitmap bloq" + /*string(name_byte[:])  + */ " con Exito"
	writeBytesInDisk(file, binario_bmap.Bytes(), msg_exito)

}

//////borrando blo0que
func BorrandoBloque_bitmap(file *os.File, index_bm_bloq int64, sB SuperB, val_bit byte) {

	/////////////////////////////escribiendo en bitmap
	//var uno byte = '1'
	var uno byte = val_bit
	epos_bitmap_dd := sB.Sb_ap_bitmap_bloques_ini + index_bm_bloq
	file.Seek(epos_bitmap_dd, 0)
	bitmal := &uno

	var binario_bmap bytes.Buffer
	binary.Write(&binario_bmap, binary.BigEndian, bitmal)
	msg_exito := ""
	writeBytesInDisk(file, binario_bmap.Bytes(), msg_exito)

}

//////i-nodos
//NewInodo(    file,      		" ", &Inodo_aptr, 			cord_inodo, 	int64(index_bm_inodo),	 sB, 		index_array_inodo, Contenido_ar)
func NewInodo(file *os.File, name string, Inodo_new *Inodo, posicion_inod int64, index_bm_inod int64, sB SuperB, i_det int64, Contenido_ar string) {

	////si es nuevo
	/// i_det = 0
	size_ar := int64(len(Contenido_ar))
	//fmt.Println("size_ar", size_ar)
	Inodo_new.I_count_inodo = index_bm_inod
	Inodo_new.I_size_archivo = size_ar
	Inodo_new.I_count_bloques_asignados = 0

	Inodo_new.I_array_bloques[0] = -1
	Inodo_new.I_array_bloques[1] = -1
	Inodo_new.I_array_bloques[2] = -1
	Inodo_new.I_array_bloques[3] = -1

	Inodo_new.I_ap_indirecto = -1

	//Inodo_new.I_id_proper
	//Inodo_new.I_gid
	//Inodo_new.I_perm

	////////

	file.Seek(posicion_inod, 0)
	carpet := &Inodo_new

	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, carpet)
	msg_exito := "Se Creó inod " + /*string(name_byte[:])  + */ " con Exito"
	writeBytesInDisk(file, binario.Bytes(), msg_exito)

	/////////////////////////////escribiendo en bitmap
	var uno byte = '1'
	epos_bitmap_dd := sB.Sb_ap_bitmap_tabla_inodo_ini + index_bm_inod
	file.Seek(epos_bitmap_dd, 0)
	bitmal := &uno

	var binario_bmap bytes.Buffer
	binary.Write(&binario_bmap, binary.BigEndian, bitmal)
	msg_exito = "Se Creó Actualizo bitmap ino" + /*string(name_byte[:])  + */ " con Exito"
	writeBytesInDisk(file, binario_bmap.Bytes(), msg_exito)

}

func Update_bitmapInodo(file *os.File, index_bm_inod int64, sB SuperB, val_bit byte) {

	/////////////////////////////escribiendo en bitmap
	//var uno byte = '1'

	var uno byte = val_bit

	epos_bitmap_dd := sB.Sb_ap_bitmap_tabla_inodo_ini + index_bm_inod
	file.Seek(epos_bitmap_dd, 0)
	bitmal := &uno

	var binario_bmap bytes.Buffer
	binary.Write(&binario_bmap, binary.BigEndian, bitmal)
	//msg_exito = "Se Creó Actualizo bitmap ino" + /*string(name_byte[:])  + */ " con Exito"
	msg_exito := ""
	writeBytesInDisk(file, binario_bmap.Bytes(), msg_exito)

}

//////detalle de directorio
func NewDetalleDirectorio(file *os.File, name string, detDir_new *DetDirectorio, posicion_DD int64, index_bm_DD int64, sB SuperB, i_det int64) {

	////si es nuevo
	/// i_det = 0

	detDir_new.Dd_array_files[0].Dd_file_ap_inodo = -1
	detDir_new.Dd_array_files[1].Dd_file_ap_inodo = -1
	detDir_new.Dd_array_files[2].Dd_file_ap_inodo = -1
	detDir_new.Dd_array_files[3].Dd_file_ap_inodo = -1
	detDir_new.Dd_array_files[4].Dd_file_ap_inodo = -1

	copy(detDir_new.Dd_array_files[i_det].Dd_file_nombre[:], name)
	fecha_string := time.Now().Format("2006-01-02 15:04:05")
	copy(detDir_new.Dd_array_files[i_det].Dd_file_date_creacion[:], fecha_string)

	detDir_new.Dd_array_files[i_det].Dd_file_ap_inodo = -1
	detDir_new.Dd_array_files[i_det].Dd_file_date_modificacion = detDir_new.Dd_array_files[i_det].Dd_file_date_creacion

	//apuntador indirecto
	detDir_new.Dd_ap_detalle_directorio_indirec = -1

	file.Seek(posicion_DD, 0)
	carpet := &detDir_new

	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, carpet)
	msg_exito := "Se Creó DD " + /*string(name_byte[:])  + */ " con Exito"
	writeBytesInDisk(file, binario.Bytes(), msg_exito)

	/////////////////////////////escribiendo en bitmap
	var uno byte = '1'
	epos_bitmap_dd := sB.Sb_ap_bitmap_detalle_directorio_ini + index_bm_DD
	file.Seek(epos_bitmap_dd, 0)
	bitmal := &uno

	var binario_bmap bytes.Buffer
	binary.Write(&binario_bmap, binary.BigEndian, bitmal)
	msg_exito = "Se Creó Actualizo bitmap DD" + /*string(name_byte[:])  + */ " con Exito"
	writeBytesInDisk(file, binario_bmap.Bytes(), msg_exito)

}

func UpdateDD(file *os.File, Detdir DetDirectorio, pos_DD int64) {
	file.Seek(pos_DD, 0)
	detdi := &Detdir

	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, detdi)
	msg_exito := "Se Actulizó datos del DD" + /*string(name_byte[:])  + */ " con Exito"
	writeBytesInDisk(file, binario.Bytes(), msg_exito)

}

func UpdateInode(file *os.File, Inod Inodo, pos_inod int64) {
	file.Seek(pos_inod, 0)
	bloqcont := &Inod

	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, bloqcont)
	msg_exito := "Se Actulizó datos del Bloquecont" + /*string(name_byte[:])  + */ " con Exito"
	writeBytesInDisk(file, binario.Bytes(), msg_exito)
}

func UpdateAVD(file *os.File, Avds AVD, pos_AVD int64) {

	file.Seek(pos_AVD, 0)
	carpet := &Avds

	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, carpet)
	msg_exito := "Se Actulizó datos del AVD" + /*string(name_byte[:])  + */ " con Exito"
	writeBytesInDisk(file, binario.Bytes(), msg_exito)

}

func CreateNewCarpeta_indirec(file *os.File, name string, Avds AVD, posicion_avd int64, index_bm_arbol int64, sB SuperB) AVD {

	//Avds := AVD{}
	fecha_string := time.Now().Format("2006-01-02 15:04:05")
	copy(Avds.Avd_fecha_creacion[:], fecha_string)
	copy(Avds.Avd_nombre_directorio[:], name)
	///array de 6 subdireccionrios
	Avds.Avd_ap_array_subdirectorios[0] = -1
	Avds.Avd_ap_array_subdirectorios[1] = -1
	Avds.Avd_ap_array_subdirectorios[2] = -1
	Avds.Avd_ap_array_subdirectorios[3] = -1
	Avds.Avd_ap_array_subdirectorios[4] = -1
	Avds.Avd_ap_array_subdirectorios[5] = -1
	///apuntados a los detalles de directorios
	Avds.Avd_ap_detalle_directorio = -1
	//apuntador indirecto
	Avds.Avd_ap_arbol_virtual_directorio = -1
	///datos de propietario
	/*Avd_proper int64
	Avd_gid int64
	Avd_perm int64*/
	///////////////////////////////empezar_escribir := Ini_AVD ///+ (index * size_avd)
	/*escribiendo carpeta, directorio*/
	//////////WriteCarpeta(path, Avds, empezar_escribir, Ini_bitmapAVD, index)

	file.Seek(posicion_avd, 0)
	carpet := &Avds

	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, carpet)
	msg_exito := "Se Creó Carpeta " + /*string(name_byte[:])  + */ " con Exito"
	writeBytesInDisk(file, binario.Bytes(), msg_exito)

	/////////////////////////////escribiendo en bitmap
	var uno byte = '1'
	epos_bitmap_avl := sB.Sb_ap_bitmap_arbol_directorio_ini + index_bm_arbol
	file.Seek(epos_bitmap_avl, 0)
	bitmal := &uno

	var binario_bmap bytes.Buffer
	binary.Write(&binario_bmap, binary.BigEndian, bitmal)
	msg_exito = "Se Creó Actualizo bitmap " + /*string(name_byte[:])  + */ " con Exito"
	writeBytesInDisk(file, binario_bmap.Bytes(), msg_exito)

	//fmt.Println("++++++++ Avds",Avds)

	return Avds

}

func CreateNewCarpeta(file *os.File, name string, Avds AVD, posicion_avd int64, index_bm_arbol int64, sB SuperB) {

	//Avds := AVD{}
	fecha_string := time.Now().Format("2006-01-02 15:04:05")
	copy(Avds.Avd_fecha_creacion[:], fecha_string)
	copy(Avds.Avd_nombre_directorio[:], name)
	///array de 6 subdireccionrios
	Avds.Avd_ap_array_subdirectorios[0] = -1
	Avds.Avd_ap_array_subdirectorios[1] = -1
	Avds.Avd_ap_array_subdirectorios[2] = -1
	Avds.Avd_ap_array_subdirectorios[3] = -1
	Avds.Avd_ap_array_subdirectorios[4] = -1
	Avds.Avd_ap_array_subdirectorios[5] = -1
	///apuntados a los detalles de directorios
	Avds.Avd_ap_detalle_directorio = -1
	//apuntador indirecto
	Avds.Avd_ap_arbol_virtual_directorio = -1
	///datos de propietario
	/*Avd_proper int64
	Avd_gid int64
	Avd_perm int64*/
	///////////////////////////////empezar_escribir := Ini_AVD ///+ (index * size_avd)
	/*escribiendo carpeta, directorio*/
	//////////WriteCarpeta(path, Avds, empezar_escribir, Ini_bitmapAVD, index)

	file.Seek(posicion_avd, 0)
	carpet := &Avds

	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, carpet)
	msg_exito := "Se Creó Carpeta " + /*string(name_byte[:])  + */ " con Exito"
	writeBytesInDisk(file, binario.Bytes(), msg_exito)

	/////////////////////////////escribiendo en bitmap
	var uno byte = '1'
	epos_bitmap_avl := sB.Sb_ap_bitmap_arbol_directorio_ini + index_bm_arbol
	file.Seek(epos_bitmap_avl, 0)
	bitmal := &uno

	var binario_bmap bytes.Buffer
	binary.Write(&binario_bmap, binary.BigEndian, bitmal)
	msg_exito = "Se Creó Actualizo bitmap " + /*string(name_byte[:])  + */ " con Exito"
	writeBytesInDisk(file, binario_bmap.Bytes(), msg_exito)

	//fmt.Println("++++++++ Avds", Avds)

}

/////////pos arbol/////////

func Pos_bitmap_AVD(file *os.File, sB SuperB) int {

	file.Seek(sB.Sb_ap_bitmap_arbol_directorio_ini, 0)
	//fmt.Println("sB.Sb_ap_bitmap_arbol_directorio_ini", sB.Sb_ap_bitmap_arbol_directorio_ini)

	var cer byte
	var size_bit = int(binary.Size(cer))

	for i := 0; i < int(sB.Sb_arbol_v_cant); i++ {
		/*pos := sB.Sb_ap_bitmap_arbol_directorio_ini + (int64(i) * int64(size_bit))
		fmt.Println("pos", pos)
		file.Seek(pos,0);*/
		data_bit := readBytesDisk(file, size_bit)
		buffer_b := bytes.NewBuffer(data_bit)

		//Decodificamos y guardamos en la variable m
		err2 := binary.Read(buffer_b, binary.BigEndian, &cer)
		if err2 != nil {
			fmt.Println("binary. se encontro error al leer archivo binario", err2)
		}
		if cer == '0' {
			return i
		}

	}
	return -1

}

/*pos detalles*/
func Pos_bitmap_Detalles(file *os.File, sB SuperB) int {

	file.Seek(sB.Sb_ap_bitmap_detalle_directorio_ini, 0)
	fmt.Println("sB.Sb_ap_bitmap_detalle_directorio_ini", sB.Sb_ap_bitmap_detalle_directorio_ini)

	var cer byte
	var size_bit = int(binary.Size(cer))

	for i := 0; i < int(sB.Sb_detalle_directorio_cant); i++ {

		data_bit := readBytesDisk(file, size_bit)
		buffer_b := bytes.NewBuffer(data_bit)

		err2 := binary.Read(buffer_b, binary.BigEndian, &cer)
		if err2 != nil {
			fmt.Println("binary. se encontro error al leer archivo binario", err2)
		}
		if cer == '0' {
			return i
		}

	}
	return -1

}

/*pos inodo*/
func Pos_bitmap_Inodos(file *os.File, sB SuperB) int {

	file.Seek(sB.Sb_ap_bitmap_tabla_inodo_ini, 0)
	fmt.Println("sB.Sb_ap_bitmap_tabla_inodo_ini", sB.Sb_ap_bitmap_tabla_inodo_ini)

	var cer byte
	var size_bit = int(binary.Size(cer))
	//Sb_detalle_directorio_cant
	for i := 0; i < int(sB.Sb_inodo_cant); i++ {

		data_bit := readBytesDisk(file, size_bit)
		buffer_b := bytes.NewBuffer(data_bit)

		err2 := binary.Read(buffer_b, binary.BigEndian, &cer)
		if err2 != nil {
			fmt.Println("binary. se encontro error al leer archivo binario", err2)
		}
		if cer == '0' {
			return i
		}
	}
	return -1
}

/*pos Bloques*/
func Pos_bitmap_Bloques(file *os.File, sB SuperB) int {

	file.Seek(sB.Sb_ap_bitmap_bloques_ini, 0)
	fmt.Println("sB.Sb_ap_bitmap_bloques_ini", sB.Sb_ap_bitmap_bloques_ini)

	var cer byte
	var size_bit = int(binary.Size(cer))
	//Sb_detalle_directorio_cant
	for i := 0; i < int(sB.Sb_bloques_cant); i++ {

		data_bit := readBytesDisk(file, size_bit)
		buffer_b := bytes.NewBuffer(data_bit)

		err2 := binary.Read(buffer_b, binary.BigEndian, &cer)
		if err2 != nil {
			fmt.Println("binary. se encontro error al leer archivo binario bloq", err2)
		}
		if cer == '0' {
			return i
		}
	}
	return -1
}

//CreateDetalle(file *os.File, path_archivo, Ini_AVD)
func CreateDetalle(file *os.File, path_archivo string, Ini_AVD int64) {

	file.Seek(int64(Ini_AVD), 0)

	Avds := AVD{}
	var size int = int(binary.Size(Avds))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &Avds)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario", err)
	}

	//fmt.Println("Avds", Avds)

	////array de directorio
	var_directorios := strings.Split(path_archivo, "/")
	////////////////carpeta root
	if len(var_directorios[0]) == 0 {
		var_directorios[0] = "/"
	}
	fmt.Println("var_directorios", var_directorios, len(var_directorios))

	var name_compar [20]byte
	copy(name_compar[:], var_directorios[0])
	fmt.Printf("name_compar   %s\n", name_compar)
	fmt.Printf("Avd_nombre_directorio   %s\n", Avds.Avd_nombre_directorio)

	/////////////////creando archivo
	/*si es igual, entonces busco insertar en detalle si no es*/
	if name_compar == Avds.Avd_nombre_directorio {
		fmt.Println("Avds", "son igualitos")
		/*busco si tiene detalle de articulo*/

		/*no existe, creo*/
		if Avds.Avd_ap_detalle_directorio == -1 {

			/*DetDirec := DetDirectorio{}

			Detinfo := DetDirINFO{}

			type DetDirINFO struct {
				Dd_file_nombre [20]byte
				Dd_file_ap_inodo int64

				Dd_file_date_creacion [19]byte
				Dd_file_date_modificacion  [19]byte
			}

			DetDirec.Dd_array_files
			//apuntador indirecto
			DetDirec.Dd_ap_detalle_directorio_indirec = -1*/

			/*exite, entonces*/
		} else {

		}
	}

	/*fecha_string := time.Now().Format("2006-01-02 15:04:05")
	copy(Avds.Avd_fecha_creacion[:], fecha_string)
	copy(Avds.Avd_nombre_directorio[:], name)

	///array de 6 subdireccionrios
	Avds.Avd_ap_array_subdirectorios[0] = -1
	Avds.Avd_ap_array_subdirectorios[1] = -1
	Avds.Avd_ap_array_subdirectorios[2] = -1
	Avds.Avd_ap_array_subdirectorios[3] = -1
	Avds.Avd_ap_array_subdirectorios[4] = -1
	Avds.Avd_ap_array_subdirectorios[5] = -1

	///apuntados a los detalles de directorios
	Avds.Avd_ap_detalle_directorio = -1

	//apuntador indirecto
	Avds.Avd_ap_arbol_virtual_directorio  = -1

	///datos de propietario
	//Avd_proper int64
	//Avd_gid int64
	//Avd_perm int64
	empezar_escribir := Ini_AVD ///+ (index * size_avd)*/

}

func CreateCarpeta(path string, name string, Ini_bitmapAVD int64, Ini_AVD int64, index int64, size_avd int64) {

	Avds := AVD{}

	fecha_string := time.Now().Format("2006-01-02 15:04:05")
	copy(Avds.Avd_fecha_creacion[:], fecha_string)
	copy(Avds.Avd_nombre_directorio[:], name)

	///array de 6 subdireccionrios
	Avds.Avd_ap_array_subdirectorios[0] = -1
	Avds.Avd_ap_array_subdirectorios[1] = -1
	Avds.Avd_ap_array_subdirectorios[2] = -1
	Avds.Avd_ap_array_subdirectorios[3] = -1
	Avds.Avd_ap_array_subdirectorios[4] = -1
	Avds.Avd_ap_array_subdirectorios[5] = -1

	///apuntados a los detalles de directorios
	Avds.Avd_ap_detalle_directorio = -1

	//apuntador indirecto
	Avds.Avd_ap_arbol_virtual_directorio = -1

	///datos de propietario
	/*Avd_proper int64
	Avd_gid int64
	Avd_perm int64*/
	empezar_escribir := Ini_AVD ///+ (index * size_avd)

	/*escribiendo carpeta, directorio*/
	WriteCarpeta(path, Avds, empezar_escribir, Ini_bitmapAVD, index)

}

func WriteCarpeta(path string, Avds AVD, inicio_struct int64, Ini_bitmapAVD int64, index int64) {

	/*escribir carpeta*/
	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(inicio_struct, 0)
	carpet := &Avds

	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, carpet)
	msg_exito := "Se Creó Carpeta " + /*string(name_byte[:])  + */ " con Exito"
	writeBytesInDisk(file, binario.Bytes(), msg_exito)

	WriteBitmap_arbol(file, Ini_bitmapAVD, index)
}

func WriteBitmap_arbol(file *os.File, Ini_bitmapAVD int64, index int64) {

	var uno byte = '1'
	//var size_bit = int(binary.Size(uno))

	escribir_en := Ini_bitmapAVD + index
	file.Seek(escribir_en, 0)
	bitmal := &uno

	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, bitmal)
	msg_exito := "Se Creó Actualizo bitmap " + /*string(name_byte[:])  + */ " con Exito"
	writeBytesInDisk(file, binario.Bytes(), msg_exito)
}

///////////////

/*leyendo archivo, disco*/
/*func readDisk(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	return file
}*/

func writeBytesDisk(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal("err ", err)
	} /*else {
		fmt.Println(msg_exito)
	}*/
}

func writeBytesInDisk(file *os.File, bytes []byte, msg_exito string) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal("err ", err)
	} else {
		//fmt.Println(msg_exito)
	}
}

//*********************reportes*************************************///
func NameByteToString(name_byte [20]byte) string {

	name_part := ""
	for c := 0; c < len(name_byte); c++ {
		if name_byte[c] == 0 {
			name_part = name_part + " "
		} else {
			name_part = name_part + string(name_byte[c])
		}
	}
	return name_part
}

///contenido bloque
func BloqueByteToString(name_byte [25]byte) string {
	name_part := ""
	for c := 0; c < len(name_byte); c++ {
		if name_byte[c] == 0 {
			name_part = name_part + " "
		} else {
			name_part = name_part + string(name_byte[c])
		}
	}
	return name_part
}

func GraficBM_Arbol(path_disco string, path_save string, inicia_part uint64) {

	var rep_bm_arbol string
	file, err := os.OpenFile(path_disco, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(int64(inicia_part), 0)

	sB := SuperB{}
	var size int = int(binary.Size(sB))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &sB)
	if err != nil {
		log.Fatal("binary. se encontro error al leer archivo binario", err)
	}

	////fmt.Println(sB)

	/*recorriendo el bitmap*/
	file.Seek(sB.Sb_ap_bitmap_arbol_directorio_ini, 0)
	fmt.Println("sB.Sb_ap_bitmap_arbol_directorio_ini", sB.Sb_ap_bitmap_arbol_directorio_ini)

	///cer := Letra{}
	var cer byte
	var size_bit = int(binary.Size(cer))
	fmt.Println("size_bit", size_bit)
	//cer := SuperB{}
	//pos := sB.Sb_ap_bitmap_arbol_directorio_ini
	//for i := 0; i < 30; i++ {
	rep_bm_arbol = ""
	col := 0
	for i := 0; i < int(sB.Sb_arbol_v_cant); i++ {

		/*pos := sB.Sb_ap_bitmap_arbol_directorio_ini + (int64(i) * int64(size_bit))
		fmt.Println("pos", pos)
		file.Seek(pos,0);*/

		col++
		data_bit := readBytesDisk(file, size_bit)
		buffer_b := bytes.NewBuffer(data_bit)

		//Decodificamos y guardamos en la variable m
		err2 := binary.Read(buffer_b, binary.BigEndian, &cer)
		if err2 != nil {
			fmt.Println("binary. se encontro error al leer archivo binario", err2)
		}
		rep_bm_arbol = rep_bm_arbol + string(cer) + "|"
		if col == 20 {
			rep_bm_arbol = rep_bm_arbol + "\n"
			col = 0
		}

		//fmt.Println(i , string(cer))
		//pos = pos +  int64(size_bit)

	}

	BM_arbol(rep_bm_arbol, path_save)

	//fmt.Println( cer.Sb_nombre_hd)

	/*btimap := make([]byte, sB.Sb_arbol_v_cant )
	cont := make([]byte, sB.Sb_arbol_v_cant )

	_, err = file.Read(cont)
	if err != nil {
		fmt.Println(err)
	}
	buffer_bit := bytes.NewBuffer(cont)
	err = binary.Read(buffer_bit, binary.BigEndian, &btimap)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario", err)
	}
	//fmt.Println("btimap", btimap)

	for i := 0; i < 25; i++ {
		fmt.Println(i, string(btimap[i]))
	}*/

}

///completo
func Gra_TreeComplete(path_disco string, path_save string, inicia_part uint64, tipo_rep string) {

	/*file, err := os.OpenFile(path_disco, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(int64(inicia_part) , 0)

	sB := SuperB{}
	var size int  = int(binary.Size(sB))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	err = binary.Read(buffer, binary.BigEndian, &sB)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario", err)
	}*/

	///////////////////////////////////////
	RecorroDirectorio(path_disco, inicia_part)
	report_tree_complet = report_tree_complet + "}\n"
	EscribirImagen(report_tree_complet, path_save)
}

///solo carpetas
func Gra_Directorio(path_disco string, path_save string, inicia_part uint64, tipo_rep string) {
	///////////////////////////////////////
	RecorroDirectorio(path_disco, inicia_part)
	report_directorio = report_directorio + "}\n"
	EscribirImagen(report_directorio, path_save)
}

///archivos de la carpeta
func Graf_Tree_directorio(path_disco string, path_save string, inicia_part uint64, tipo_rep string, ruta_ruta string) {
	///////////////////////////////////////
	file, err := os.OpenFile(path_disco, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(int64(inicia_part), 0)

	sB := SuperB{}
	var size int = int(binary.Size(sB))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &sB)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario", err)
	}
	Ini_AVD := sB.Sb_ap_arbol_directorio_ini

	carpeta_crear := ruta_ruta
	var_directorios := strings.Split(carpeta_crear, "/")
	if len(var_directorios[0]) == 0 {
		var_directorios[0] = "/"
	}
	fmt.Println("var_directorios", var_directorios, len(var_directorios))

	/////////////////////////

	report_archiv_incar = ""
	report_archiv_incar = "digraph arbolito_com {\n" +
		"node [shape=plaintext]\n" +
		"rankdir=LR;\n"

	var tipo_mk byte = 'C'
	var tipo_accion string = "POS_CAR"

	Leyendo_Recorrido(file, Ini_AVD, var_directorios, 0, sB, tipo_mk, tipo_accion, "", "")

	fmt.Println("posicion_de_carpeta_rep", posicion_de_carpeta_rep)
	fmt.Println("encontrado_tipo_art", encontrado_tipo_art)

	if encontrado_tipo_art == true {
		//Leyendo_Recorrido(file, Ini_AVD, var_directorios, 0, sB, tipo_mk, tipo_accion, new_name)
		//reporte de archivos de directorio

		Recorro_Carpetas_TreeDir(file, posicion_de_carpeta_rep, sB, 0)
		report_archiv_incar = report_archiv_incar + "}\n"
		EscribirImagen(report_archiv_incar, path_save)
	}

}

///el camino del archivo
func Graf_Tree_File(path_disco string, path_save string, inicia_part uint64, tipo_rep string, ruta_ruta string) {
	///////////////////////////////////////
	file, err := os.OpenFile(path_disco, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(int64(inicia_part), 0)

	sB := SuperB{}
	var size int = int(binary.Size(sB))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &sB)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario", err)
	}
	Ini_AVD := sB.Sb_ap_arbol_directorio_ini

	carpeta_crear := ruta_ruta
	var_directorios := strings.Split(carpeta_crear, "/")
	if len(var_directorios[0]) == 0 {
		var_directorios[0] = "/"
	}
	fmt.Println("var_directorios", var_directorios, len(var_directorios))

	/////////////////////////

	report_camino_file = ""
	report_camino_file = "digraph camino_fil {\n" +
		"node [shape=plaintext]\n" +
		"rankdir=LR;\n"

	var tipo_mk byte = 'A'
	var tipo_accion string = "REP_TF"

	Leyendo_Recorrido(file, Ini_AVD, var_directorios, 0, sB, tipo_mk, tipo_accion, "", "")

	//fmt.Println("posicion_de_carpeta_rep", posicion_de_carpeta_rep)
	//fmt.Println("encontrado_tipo_art", encontrado_tipo_art)

	if encontrado_tipo_art == true {
		//Leyendo_Recorrido(file, Ini_AVD, var_directorios, 0, sB, tipo_mk, tipo_accion, new_name)
		//reporte de archivos de directorio

		//Recorro_Carpetas_TreeDir(file, posicion_de_carpeta_rep, sB, 0)
		report_camino_file = report_camino_file + "}\n"
		EscribirImagen(report_camino_file, path_save)
	}

}

/*cat*/
func Cat_Print(path_disco string, inicia_part uint64, arr_file_read []string) {

	file, err := os.OpenFile(path_disco, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(int64(inicia_part), 0)

	sB := SuperB{}
	var size int = int(binary.Size(sB))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &sB)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario", err)
	}
	//fmt.Println(sB)

	/**
	Leer Contenido, CAt = A
	*/
	var tipo_mk byte = 'A'
	Ini_AVD := sB.Sb_ap_arbol_directorio_ini

	//Creando_Archivos(var_directorios, file, Ini_AVD, sB, p, tipo_mk)
	//Recorriendo_directorio(file, Ini_AVD, var_directorios, 0, sB, Contenido_ar, p, tipo_mk)

	//Leyendo_Recorrido(file, Ini_AVD, var_directorios, 0, sB, tipo_mk)
	var tipo_accion string = "PRI_AR"

	for i := 0; i < len(arr_file_read); i++ {

		carpeta_crear := arr_file_read[i]
		var_directorios := strings.Split(carpeta_crear, "/")
		if len(var_directorios[0]) == 0 {
			var_directorios[0] = "/"
		}
		//fmt.Println("var_directorios", var_directorios, len(var_directorios))
		Leyendo_Recorrido(file, Ini_AVD, var_directorios, 0, sB, tipo_mk, tipo_accion, "", "")
	}

}

/*renomnbar*/
var encontrado_tipo_art bool = false

func Rem_car_file(path_disco string, inicia_part uint64, path_file string, new_name string) {

	file, err := os.OpenFile(path_disco, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(int64(inicia_part), 0)

	sB := SuperB{}
	var size int = int(binary.Size(sB))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &sB)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario", err)
	}
	//fmt.Println(sB)

	var tipo_accion string = "TIPO_AR"

	/**
	Leer Contenido, CAt = A
	*/

	Ini_AVD := sB.Sb_ap_arbol_directorio_ini

	carpeta_crear := path_file
	var_directorios := strings.Split(carpeta_crear, "/")
	if len(var_directorios[0]) == 0 {
		var_directorios[0] = "/"
	}
	fmt.Println("var_directorios", var_directorios, len(var_directorios))

	//fmt.Println("tipo_mk", tipo_mk, Ini_AVD)

	encontrado_tipo_art = false
	/*buscandi tipo de archivo*/
	var tipo_mk byte = 'A'
	Leyendo_Recorrido(file, Ini_AVD, var_directorios, 0, sB, tipo_mk, tipo_accion, "", "")
	//fmt.Println("111 para archivo, encontrado_tipo_art", encontrado_tipo_art)

	//return
	if encontrado_tipo_art == false {
		tipo_mk = 'C'
		Leyendo_Recorrido(file, Ini_AVD, var_directorios, 0, sB, tipo_mk, tipo_accion, "", "")
		//fmt.Println("2222 para carpeta, encontrado_tipo_art", encontrado_tipo_art)
	}

	/*si existen, y es de algun tipo entonces ya cambio nombre*/
	if encontrado_tipo_art == true {
		tipo_accion = "MOD_NAME"
		Leyendo_Recorrido(file, Ini_AVD, var_directorios, 0, sB, tipo_mk, tipo_accion, new_name, "")
	}

}
func GraficBM_Rep(path_disco string, path_save string, inicia_part uint64, tipo_rep string) {

	var rep_bm_detalle string
	file, err := os.OpenFile(path_disco, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(int64(inicia_part), 0)

	sB := SuperB{}
	var size int = int(binary.Size(sB))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	err = binary.Read(buffer, binary.BigEndian, &sB)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario", err)
	}

	///////////////////////////////////////
	if tipo_rep == "bm_detdir" {

		file.Seek(sB.Sb_ap_bitmap_detalle_directorio_ini, 0)
		fmt.Println("sB.Sb_ap_bitmap_detalle_directorio_ini", sB.Sb_ap_bitmap_detalle_directorio_ini)

		var cer byte
		var size_bit = int(binary.Size(cer))
		fmt.Println("size_bit", size_bit)

		rep_bm_detalle = ""
		col := 0
		for i := 0; i < int(sB.Sb_detalle_directorio_cant); i++ {

			/*pos := sB.Sb_ap_bitmap_arbol_directorio_ini + (int64(i) * int64(size_bit))
			fmt.Println("pos", pos)
			file.Seek(pos,0);*/

			col++
			data_bit := readBytesDisk(file, size_bit)
			buffer_b := bytes.NewBuffer(data_bit)

			//Decodificamos y guardamos en la variable m
			err2 := binary.Read(buffer_b, binary.BigEndian, &cer)
			if err2 != nil {
				fmt.Println("binary. se encontro error al leer archivo binario", err2)
			}
			rep_bm_detalle = rep_bm_detalle + string(cer) + "|"
			if col == 20 {
				rep_bm_detalle = rep_bm_detalle + "\n"
				col = 0
			}
		}

		///////////////////inodos////////////////////
	} else if tipo_rep == "bm_inode" {

		file.Seek(sB.Sb_ap_bitmap_tabla_inodo_ini, 0)
		fmt.Println("sB.Sb_ap_bitmap_tabla_inodo_ini", sB.Sb_ap_bitmap_tabla_inodo_ini)

		var cer byte
		var size_bit = int(binary.Size(cer))
		fmt.Println("size_bit", size_bit)

		rep_bm_detalle = ""
		col := 0
		for i := 0; i < int(sB.Sb_inodo_cant); i++ {

			/*pos := sB.Sb_ap_bitmap_arbol_directorio_ini + (int64(i) * int64(size_bit))
			fmt.Println("pos", pos)
			file.Seek(pos,0);*/

			col++
			data_bit := readBytesDisk(file, size_bit)
			buffer_b := bytes.NewBuffer(data_bit)

			//Decodificamos y guardamos en la variable m
			err2 := binary.Read(buffer_b, binary.BigEndian, &cer)
			if err2 != nil {
				fmt.Println("binary. se encontro error al leer archivo binario", err2)
			}
			rep_bm_detalle = rep_bm_detalle + string(cer) + "|"
			if col == 20 {
				rep_bm_detalle = rep_bm_detalle + "\n"
				col = 0
			}
		}
		///////////////////bloques////////////////////
	} else if tipo_rep == "bm_block" {

		file.Seek(sB.Sb_ap_bitmap_bloques_ini, 0)
		fmt.Println("sB.Sb_ap_bitmap_bloques_ini", sB.Sb_ap_bitmap_bloques_ini)

		var cer byte
		var size_bit = int(binary.Size(cer))
		fmt.Println("size_bit", size_bit)

		rep_bm_detalle = ""
		col := 0
		for i := 0; i < int(sB.Sb_bloques_cant); i++ {

			/*pos := sB.Sb_ap_bitmap_arbol_directorio_ini + (int64(i) * int64(size_bit))
			fmt.Println("pos", pos)
			file.Seek(pos,0);*/

			col++
			data_bit := readBytesDisk(file, size_bit)
			buffer_b := bytes.NewBuffer(data_bit)

			//Decodificamos y guardamos en la variable m
			err2 := binary.Read(buffer_b, binary.BigEndian, &cer)
			if err2 != nil {
				fmt.Println("binary. se encontro error al leer archivo binario", err2)
			}
			rep_bm_detalle = rep_bm_detalle + string(cer) + "|"
			if col == 20 {
				rep_bm_detalle = rep_bm_detalle + "\n"
				col = 0
			}
		}
	}

	///////////////////////////////////////

	BM_arbol(rep_bm_detalle, path_save)
}

func GraficSB(path_disco string, path_save string, inicia_part uint64) {

	var graf_sb string
	file, err := os.OpenFile(path_disco, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(int64(inicia_part), 0)

	sB := SuperB{}
	var size int = int(binary.Size(sB))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &sB)
	if err != nil {
		fmt.Println("binary. se encontro error al leer archivo binario", err)
	}

	///fmt.Println(sB)
	/////////////
	graf_sb = "digraph test {\n" +
		"graph [ratio=fill];\n" +
		"node [label=\"\\N\", fontsize=15, shape=plaintext];\n" +
		//"graph [bb=\"0,0,352,154\"];\n"+
		"arset[label=<\n" +
		"<TABLE ALIGN=\"LEFT\">\n" +
		"<tr>\n" +
		"	<TD>Nombre</TD>\n" +
		"	<TD>Valor</TD>\n" +
		"</tr>\n" +

		//nombre del disco duro
		"<tr>\n" +
		"	<TD>Sb_nombre_hd</TD>\n" +
		"	<TD>" + NameByteToString(sB.Sb_nombre_hd) + "</TD>\n" +
		"</tr>\n" +

		//cantidad de estructuras
		"<tr>\n" +
		"	<TD>Sb_arbol_virtual_count</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_arbol_v_cant)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_detalle_directorio_count</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_detalle_directorio_cant)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_inodo_count</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_inodo_cant)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_bloques_count</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_bloques_cant)) + "</TD>\n" +
		"</tr>\n" +

		//libres
		"<tr>\n" +
		"	<TD>Sb_arbol_virtual_free</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_arbol_virtual_free)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_detalle_directorio_free</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_detalle_directorio_free)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_inodos_free</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_inodos_free)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_bloques_free</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_bloques_free)) + "</TD>\n" +
		"</tr>\n" +

		///fecha
		"<tr>\n" +
		"	<TD>Sb_data_creacion</TD>\n" +
		"	<TD>" + string(sB.Sb_data_creacion[:]) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_data_ultimo_montaje</TD>\n" +
		"	<TD>" + string(sB.Sb_data_ultimo_montaje[:]) + "</TD>\n" +
		"</tr>\n" +

		//Sb_montajes_conta
		"<tr>\n" +
		"	<TD>Sb_montajes_count</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_montajes_conta)) + "</TD>\n" +
		"</tr>\n" +

		///apuntadores al inicio
		"<tr>\n" +
		"	<TD>Sb_ap_bitmap_arbol_directorio_ini</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_ap_bitmap_arbol_directorio_ini)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_ap_arbol_directorio_ini</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_ap_arbol_directorio_ini)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_ap_bitmap_detalle_directorio_ini</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_ap_bitmap_detalle_directorio_ini)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_ap_detalle_directorio_ini</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_ap_detalle_directorio_ini)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_ap_bitmap_tabla_inodo_ini</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_ap_bitmap_tabla_inodo_ini)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_ap_bitmap_tabla_inodo_ini</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_ap_bitmap_tabla_inodo_ini)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_ap_tabla_inodo_ini</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_ap_tabla_inodo_ini)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_ap_bitmap_bloques_ini</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_ap_bitmap_bloques_ini)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_ap_bloques_ini</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_ap_bloques_ini)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_ap_log_ini</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_ap_log_ini)) + "</TD>\n" +
		"</tr>\n" +

		////tamaños de las esctucturas
		"<tr>\n" +
		"	<TD>Sb_size_struct_arbol_directorio</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_size_struct_arbol_directorio)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_size_struct_detalle_directorio</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_size_struct_detalle_directorio)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_size_struct_inodo</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_size_struct_inodo)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_size_struct_bloque</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_size_struct_bloque)) + "</TD>\n" +
		"</tr>\n" +

		///primer bit en el bitmap
		"<tr>\n" +
		"	<TD>Sb_first_free_bit_arbol_directorio</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_first_free_bit_arbol_directorio)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_first_free_bit_detalle_directorio</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_first_free_bit_detalle_directorio)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_first_free_bit_tabla_inodo</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_first_free_bit_tabla_inodo)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Sb_first_free_bit_bloques</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_first_free_bit_bloques)) + "</TD>\n" +
		"</tr>\n" +

		////mi carnet
		"<tr>\n" +
		"	<TD>Sb_magic_num</TD>\n" +
		"	<TD>" + strconv.Itoa(int(sB.Sb_magic_num)) + "</TD>\n" +
		"</tr>\n"

	graf_sb = graf_sb + "</TABLE>\n" +
		">, ];\n" +
		"}"

	EscribirImagen(graf_sb, path_save)

}

/*escribir reporte de bitmaps*/
func BM_arbol(txt_ar string, path_save string) {

	//file, err := os.Create(path_save + ".txt")
	file, err := os.Create(path_save)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	_, err = file.WriteString(txt_ar)
	if err != nil {
		fmt.Println("err ", err)
	} else {
		fmt.Println("Reporte Bm creado")
	}
}

func EscribirImagen(image string, path_save string) {

	file, err := os.Create(path_save + ".txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	_, err = file.WriteString(image)
	if err != nil {
		fmt.Println("err ", err)
	}

	ExeCommand(path_save)

}

func ExeCommand(path_save string) {
	//cmd := exec.Command("/home/rafaelc/Dis/","dot -Tpng mbr_gra.txt -o mbr_gra.png")
	cmd := exec.Command("dot", "-Tpng", path_save+".txt", "-o", path_save)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%q", out.String())
}

func readBytesDisk(file *os.File, number int) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		fmt.Println(err)
	}
	return bytes
}
