package Disks

import (
	"time"

	"fmt"

	//"encoding/binary"

	//para crear archivo binarop
	"log"
	"os"

	//"unsafe"
	"bytes"
	"encoding/binary"

	///leer pantalla
	"bufio"
	"strings"

	"strconv" // int to string

	"os/exec"

	"../FormatP"
)

///Master Boot Record *  175 Bytes
type MBR struct {
	Mbr_tamanio uint64 // 8 bytes
	//Mbr_fecha_creacion time.Time //24 bytes
	Mbr_fecha_creacion_s [19]byte // 19 bytes
	Mbr_disk_signature   uint64   // 8 bytes

	///140
	/*Mbr_part1 Particion // 35 bytes
	Mbr_part2 Particion // 35 bytes
	Mbr_part3 Particion // 35 bytes
	Mbr_part4 Particion // 35 bytes*/
	Mbr_Particiones [4]Particion
}

///Extended Boot Record
//42 bytes
type EBR struct {
	Part_status byte     // 1 bytes
	Part_fit    byte     // 1 bytes
	Part_start  uint64   // 8 bytes
	Part_size   uint64   // 8 bytes
	Part_next   int64    // 8 bytes
	Part_name   [16]byte // 16 bytes
	//Part_next uint64
}

/*listado de monta*/
var monta []Montaciones
var let_mon []Letras

/*Struct para Montar particiones*/
type Montaciones struct {
	M_path    string
	M_name    [16]byte
	M_id      string
	M_tipo    byte
	M_P_Parti Particion
	M_P_Log   EBR

	M_letra byte
	M_num   int64
}

type Letras struct {
	L_path  string
	L_letra byte
}

func DeleteDisk(path_all string) {

	confirmacion := false
	for confirmacion == false {

		fmt.Println("Esta Seguro de Eliminar Disco?")
		fmt.Println("S/N")

		reader := bufio.NewReader(os.Stdin)
		lin_coman, _ := reader.ReadString('\n')
		comando := strings.ReplaceAll(lin_coman, "\r", "")
		comando = strings.ReplaceAll(comando, "\n", "")

		comando = strings.ToLower(comando)
		//fmt.Println("-",comando,"+")

		if comando == "s" {
			confirmacion = true
		} else if comando == "n" {
			confirmacion = true
			return
		} else {
			confirmacion = false
		}
	}

	//if (comando == "s") {
	err := os.Remove(path_all)
	if err != nil {
		fmt.Printf("Error eliminando Disco: %v\n", err)
	} else {
		fmt.Println("Disco Eliminado correctamente")
	}
	//}
}

func ExisteCarpeta(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func existeDisk(path_dis string) bool {
	if _, err := os.Stat(path_dis); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

/*escribiendo archivo, disco*/
func CreateDisk(size_d uint64, path_d string, nombre string, unidad_d byte) {

	/*calcualando tamaño de disco*/
	var size_disk uint64
	if unidad_d == 'k' {
		size_disk = size_d * 1024

	} else if unidad_d == 'm' {
		size_disk = size_d * 1024 * 1024

	}

	///////fmt.Println("**new disk: en class disk", name_d, size_disk, string(unidad_d))

	exis, _ := ExisteCarpeta(path_d)

	if exis == false {

		err := os.MkdirAll(path_d, 0777)

		if err != nil {
			panic(err)
		}

	}

	/*antes verifica si existe path*/
	path_completo := path_d + nombre
	//file, err := os.Create(name_d)
	file, err := os.Create(path_completo)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	///valor inicial-final
	var cfin int8 = 0
	fi := &cfin
	//file.Seek(1024, 0) // 1 kb
	file.Seek(int64(size_disk), 0) // tamaño del disco

	//Escribimos un 0 al final del archivo.
	var binario_fin bytes.Buffer
	binary.Write(&binario_fin, binary.BigEndian, fi)
	writeBytesDisk(file, binario_fin.Bytes())

	file.Seek(0, 0) // nos posicionamos en el inicio del archivo.

	//Master Boot Record *  175 Bytes
	mbr_disk := MBR{Mbr_tamanio: size_disk}
	fecha_string := time.Now().Format("2006-01-02 15:04:05")
	copy(mbr_disk.Mbr_fecha_creacion_s[:], fecha_string)
	mbr_disk.Mbr_disk_signature = 1

	/*valor en particiones, init*/
	par1 := Particion{}
	par1.Part_status = 'N'
	par2 := Particion{}
	par2.Part_status = 'N'
	par3 := Particion{}
	par3.Part_status = 'N'
	par4 := Particion{}
	par4.Part_status = 'N'

	//para prueba
	////ejemplo, prueba original
	/*par3.Part_status = 'S'
	par3.Part_start = 175
	par3.Part_size = 5 * 1024 * 1024
	copy(par3.Part_name[:], "par31")

	par1.Part_status = 'S'
	par1.Part_start = 8 * 1024 * 1024
	par1.Part_size = 2 * 1024 * 1024
	copy(par1.Part_name[:], "part13")*/

	//par3.Part_status = 'S'
	//par3.Part_start = 13 * 1024 * 1024
	//par3.Part_size = 5242880
	//copy(par3.Part_name[:], "par31")

	//par1.Part_status = 'S'
	//par1.Part_start = 10 * 1024 * 1024
	//par1.Part_size = 3145728
	//copy(par1.Part_name[:], "part13")

	/*mbr_disk.Mbr_part1 = par1
	mbr_disk.Mbr_part2 = par2
	mbr_disk.Mbr_part3 = par3
	mbr_disk.Mbr_part4 = par4*/
	mbr_disk.Mbr_Particiones[0] = par1
	mbr_disk.Mbr_Particiones[1] = par2
	mbr_disk.Mbr_Particiones[2] = par3
	mbr_disk.Mbr_Particiones[3] = par4

	disk_w := &mbr_disk

	///escribimos el struct
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, disk_w)
	writeBytesDisk(file, binario.Bytes())

	/*
		disk2 := mbr2{}
		disk2.N1 = 14
		disk2.N2 = 22
		disk2.Cadb = 'D'

		disk_w22 := &disk2

		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, disk_w22)
		writeBytesDisk(file, binario2.Bytes())*/

}

////PARTICIONES 35 bytes
type Particion struct {
	Part_status byte
	Part_type   byte
	Part_fit    byte
	Part_start  uint64 // 8 bytes
	Part_size   uint64 // 8 bytes
	Part_name   [16]byte
}

/*escribiendo particion en disco*/
func CreatePartD(size uint64, path string, name string, unit byte, btype byte, fit byte) {

	/*calcualando tamaño de disco*/
	var size_part uint64
	if unit == 'b' {
		size_part = size
	} else if unit == 'k' {
		size_part = size * 1024
	} else if unit == 'm' {
		size_part = size * 1024 * 1024

	}
	////fmt.Println("size_part:",size_part )

	/*verificando si existe disco*/
	exist_d := existeDisk(path)
	///////////fmt.Println("existe ar exist_d:",exist_d )
	if exist_d == false {
		fmt.Println("No existe Ruta del Disco:", path)
		return
	}

	readDisk_new(path, btype, fit, size_part, name)
	/*creando estrucura*/
}

/*ordenando llenos*/
func OrdenmMLLenos(libre []EspacioLibre) []EspacioLibre {

	//tmp := 0
	for x := 0; x < len(libre); x++ {
		for y := 0; y < len(libre); y++ {
			if libre[x].free_start < libre[y].free_start {
				tmp := libre[y]
				libre[y] = libre[x]
				libre[x] = tmp
			}
		}
	}

	return libre
}

func OrdenmMParticion(libre []Particion) []Particion {

	//tmp := 0
	for x := 0; x < len(libre); x++ {
		for y := 0; y < len(libre); y++ {
			if libre[x].Part_start < libre[y].Part_start {
				tmp := libre[y]
				libre[y] = libre[x]
				libre[x] = tmp
			}
		}
	}

	return libre
}

func OrdenmMLetras(let []Letras) []Letras {
	for x := 0; x < len(let); x++ {
		for y := 0; y < len(let); y++ {
			if let[x].L_letra < let[y].L_letra {
				tmp := let[y]
				let[y] = let[x]
				let[x] = tmp
			}
		}
	}
	return let
}

func GetIdMon(path string /*, name_byte [16]byte*/) int64 {
	var idvol int64
	idvol = 1
	for i := 0; i < len(monta); i++ {

		if monta[i].M_path == path /*&& monta[i].M_name == name_byte*/ {
			idvol++
		}
	}
	return idvol
}

func YaMontado(path string, name_byte [16]byte) bool {

	for i := 0; i < len(monta); i++ {
		if monta[i].M_path == path && monta[i].M_name == name_byte {
			return true
		}
	}
	return false
}

func ExisteIdMontado(idvol string) (bool, string) {
	for i := 0; i < len(monta); i++ {
		if monta[i].M_id == idvol {
			return true, monta[i].M_path
		}
	}
	return false, ""
}

func GetPart_montado(idvol string) (bool, string, uint64, uint64) {
	for i := 0; i < len(monta); i++ {
		if monta[i].M_id == idvol {

			var inicia_part uint64
			var part_size uint64
			if monta[i].M_tipo == 'p' {
				inicia_part = monta[i].M_P_Parti.Part_start
				part_size = monta[i].M_P_Parti.Part_size
			} else if monta[i].M_tipo == 'l' {
				inicia_part = monta[i].M_P_Log.Part_start
				part_size = monta[i].M_P_Log.Part_size
			}

			return true, monta[i].M_path, inicia_part, part_size
		}
	}
	return false, "", 0, 0
}

func FormatPart(type_s string, idvol string) {

	exis, path_disco, inicia_part, part_size := GetPart_montado(idvol)
	if exis == false {
		fmt.Println("No existe Partición Montada", idvol)
		return
	}

	/////////////fmt.Println("path_disco", path_disco, "inicia_part", inicia_part, "part_size", part_size  )
	FormatP.Format(inicia_part, part_size, path_disco)

}

/*creando archivo*/
func CrearFile(size int64, path string, idvol string, p byte, cont string) {

	exis, path_disco, inicia_part, part_size := GetPart_montado(idvol)
	if exis == false {
		fmt.Println("No existe Partición Montada", idvol)
		return
	}
	//fmt.Println("path_disco", path_disco, "inicia_part", inicia_part, "part_size", part_size  )
	FormatP.NewFile(inicia_part, part_size, path_disco, size, p, cont, path)
}

/*editando archivo*/
func EditFile(size int64, path string, idvol string, cont string) {

	exis, path_disco, inicia_part, part_size := GetPart_montado(idvol)
	if exis == false {
		fmt.Println("No existe Partición Montada", idvol)
		return
	}
	//fmt.Println("path_disco", path_disco, "inicia_part", inicia_part, "part_size", part_size  )
	FormatP.Edit_File(inicia_part, part_size, path_disco, size, cont, path)
}

/*creando carpeta*/
func CrearCarpeta(path string, idvol string, p byte) {

	exis, path_disco, inicia_part, part_size := GetPart_montado(idvol)
	if exis == false {
		fmt.Println("No existe Partición Montada", idvol)
		return
	}
	//fmt.Println("path_disco", path_disco, "inicia_part", inicia_part, "part_size", part_size  )
	FormatP.NewDirectorio(inicia_part, part_size, path_disco, p, path)
}

func Print_Contenido(arr_file_read []string, idvol string) {

	exis, path_disco, inicia_part, _ := GetPart_montado(idvol)
	if exis == false {
		fmt.Println("No existe Partición Montada", idvol)
		return
	}

	FormatP.Cat_Print(path_disco, inicia_part, arr_file_read)
	//Leyendo_Recorrido(file *os.File, pos_AVD int64, var_directorios []string, i int64, sB SuperB, Contenido_ar string, p byte, tipo_mk byte)
}

////cambiar nombre
func Cambiar_Name(path_ac string, new_name string, idvol string) {

	exis, path_disco, inicia_part, _ := GetPart_montado(idvol)
	if exis == false {
		fmt.Println("No existe Partición Montada", idvol)
		return
	}

	FormatP.Rem_car_file(path_disco, inicia_part, path_ac, new_name)
}
func ReportesGraf(path_save string, nombre_rep string, idvol string, ruta_ruta string) {

	//exis, path_disco := ExisteIdMontado(idvol)
	//exis, path_disco, inicia_part, part_size := GetPart_montado(idvol)
	exis, path_disco, inicia_part, _ := GetPart_montado(idvol)
	if exis == false {
		fmt.Println("No existe Partición Montada", idvol)
		return
	}

	/*creando carpetas*/
	name_imagen := strings.Split(path_save, "/")
	name_imagen_str := name_imagen[len(name_imagen)-1]
	//fmt.Println("nn", name_imagen)
	//fmt.Println("name_imagen_str", name_imagen_str)
	path_save_car := strings.ReplaceAll(path_save, name_imagen_str, "")
	//fmt.Println("path_save", path_save)
	//fmt.Println("path_save_car", path_save_car)

	exis_car, _ := ExisteCarpeta(path_save_car)
	//fmt.Println("*ExisteCarpeta", exis)
	if exis_car == false {

		err := os.MkdirAll(path_save_car, 0777)
		if err != nil {
			panic(err)
		}
	}

	if strings.ToLower(nombre_rep) == "mbr" {
		//GraficMBR(path_disco, path_save)
		GraficMBR_new(path_disco, path_save)

	} else if strings.ToLower(nombre_rep) == "disk" {
		GraficDISK(path_disco, path_save)

	} else if strings.ToLower(nombre_rep) == "sb" {
		FormatP.GraficSB(path_disco, path_save, inicia_part)

		///arbol de directorio
	} else if strings.ToLower(nombre_rep) == "bm_arbdir" {
		FormatP.GraficBM_Arbol(path_disco, path_save, inicia_part)

		////detalle direc
	} else if strings.ToLower(nombre_rep) == "bm_detdir" {
		FormatP.GraficBM_Rep(path_disco, path_save, inicia_part, "bm_detdir")

		////// inodos
	} else if strings.ToLower(nombre_rep) == "bm_inode" {
		FormatP.GraficBM_Rep(path_disco, path_save, inicia_part, "bm_inode")

		////// bloques
	} else if strings.ToLower(nombre_rep) == "bm_block" {
		FormatP.GraficBM_Rep(path_disco, path_save, inicia_part, "bm_block")

		//////reporte del arbol completo
	} else if strings.ToLower(nombre_rep) == "tree_complete" {
		FormatP.Gra_TreeComplete(path_disco, path_save, inicia_part, "tree_complete")

		//////reporte solo directorio
	} else if strings.ToLower(nombre_rep) == "directorio" {
		FormatP.Gra_Directorio(path_disco, path_save, inicia_part, "directorio")

		//////reporte arhivos del directorio
	} else if strings.ToLower(nombre_rep) == "file" {
		FormatP.Graf_Tree_directorio(path_disco, path_save, inicia_part, "tree_directorio", ruta_ruta)

		//////reporte camino del archivo
	} else if strings.ToLower(nombre_rep) == "tree" {
		FormatP.Graf_Tree_File(path_disco, path_save, inicia_part, "tree_file", ruta_ruta)

	} else {
		fmt.Println("Nombre de Reporte invalido")
	}

}

func NameByteToString(name_byte [16]byte) string {

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
func Graf_extendidas(file *os.File, m MBR, num_part_extend int) string {

	var graf_extendidas string

	fin_extendida := m.Mbr_Particiones[num_part_extend].Part_start + m.Mbr_Particiones[num_part_extend].Part_size
	////////////fmt.Println("fin_extendida:", fin_extendida)
	/*posicionandome en el primer ebr de la extendida*/
	file.Seek(0, 0)
	file.Seek(int64(m.Mbr_Particiones[num_part_extend].Part_start), 0)

	ebr_part := EBR{}
	size_ebr := int(binary.Size(ebr_part))

	data := readBytesDisk(file, size_ebr)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &ebr_part)
	if err != nil {
		log.Fatal("binary. se encontro error al leer EBR de la Extendida", err)
	}

	///*************************************************///
	///////////fmt.Println("*************************************")
	///fmt.Printf("--n %s\n", ebr_part.Part_name)
	///fmt.Println("ebr", ebr_part.Part_start - uint64(size_ebr), "-", ebr_part.Part_start)
	//fmt.Println(" ", ebr_part.Part_start, "-", ebr_part.Part_start +  ebr_part.Part_size, "|",  ebr_part.Part_size)
	/*libre en logicas*/
	var libre_log uint64
	if ebr_part.Part_next == -1 {
		if ebr_part.Part_status == 'N' {
			libre_log = fin_extendida - (m.Mbr_Particiones[num_part_extend].Part_start)
		} else {
			libre_log = fin_extendida - (ebr_part.Part_start + ebr_part.Part_size)
		}
	} else {

		if ebr_part.Part_status == 'N' {
			libre_log = uint64(ebr_part.Part_next) - (m.Mbr_Particiones[num_part_extend].Part_start)
		} else {
			libre_log = uint64(ebr_part.Part_next) - (ebr_part.Part_start + ebr_part.Part_size)
		}
	}
	/////fmt.Println("                    -libre: ", libre_log)

	if ebr_part.Part_status == 'N' {

		/*Verifico si hay espacio disponible*/
		graf_extendidas = graf_extendidas + "	<td>EBR</td>\n"

	} else if ebr_part.Part_status == 'S' /*&& ebr_part.Part_next != -1*/ {
		//fmt.Println(" insertar a mitad de primero ")

		graf_extendidas = graf_extendidas + "	<td>EBR</td>\n" +
			"	<td>" + NameByteToString(ebr_part.Part_name) + "</td>\n"
	}

	if libre_log > 0 {
		graf_extendidas = graf_extendidas + "	<td>Libre</td>\n"
	}

	ebr_anterior := ebr_part
	fmt.Printf("						--anterior %s\n", ebr_anterior.Part_name)
	tempo_sig := ebr_part.Part_next
	for tempo_sig != -1 {

		file.Seek(int64(tempo_sig), 0)
		ebr_tem := EBR{}
		size_tem := int(binary.Size(ebr_tem))
		data_tem := readBytesDisk(file, size_tem)
		buffer_tem := bytes.NewBuffer(data_tem)
		err = binary.Read(buffer_tem, binary.BigEndian, &ebr_tem)
		if err != nil {
			log.Fatal("binary. se encontro error al leer EBR de la Extendida", err)
			break
		}

		graf_extendidas = graf_extendidas + "	<td>EBR</td>\n" +
			"	<td>" + NameByteToString(ebr_tem.Part_name) + "</td>\n"

		//fmt.Printf("--n %s\n", ebr_tem.Part_name)
		//fmt.Println("ebr", ebr_tem.Part_start - uint64(size_tem), "-", ebr_tem.Part_start)
		//fmt.Println(" ", ebr_tem.Part_start, "-", ebr_tem.Part_start + ebr_tem.Part_size, "|",  ebr_tem.Part_size)
		/*libre en logicas*/
		var libre_log uint64
		if ebr_tem.Part_next == -1 {
			libre_log = fin_extendida - (ebr_tem.Part_start + ebr_tem.Part_size)
		} else {
			libre_log = uint64(ebr_tem.Part_next) - (ebr_tem.Part_start + ebr_tem.Part_size)
		}
		/////fmt.Println("                    libre: ", libre_log)

		if libre_log > 0 {
			graf_extendidas = graf_extendidas + "	<td>Libre</td>\n"
		}

		tempo_sig = ebr_tem.Part_next
		ebr_anterior = ebr_tem
		fmt.Printf("						--anterior %s\n", ebr_anterior.Part_name)

		/*******************fin read sig**************************/
	}

	//////////////////}

	return graf_extendidas
}

func EBR_extendidas(file *os.File, m MBR, num_part_extend int) string {

	//var graf_extendidas string
	var graf_extendidas_tab string

	/////////fin_extendida := m.Mbr_Particiones[num_part_extend].Part_start + m.Mbr_Particiones[num_part_extend].Part_size
	/////////fmt.Println("fin_extendida:", fin_extendida)
	/*posicionandome en el primer ebr de la extendida*/
	file.Seek(0, 0)
	file.Seek(int64(m.Mbr_Particiones[num_part_extend].Part_start), 0)

	ebr_part := EBR{}
	size_ebr := int(binary.Size(ebr_part))

	data := readBytesDisk(file, size_ebr)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &ebr_part)
	if err != nil {
		log.Fatal("binary. se encontro error al leer EBR de la Extendida", err)
	}

	///*************************************************///
	var index int = 0
	/////fmt.Println("*************************************")
	///fmt.Printf("--n %s\n", ebr_part.Part_name)
	///fmt.Println("ebr", ebr_part.Part_start - uint64(size_ebr), "-", ebr_part.Part_start)
	//fmt.Println(" ", ebr_part.Part_start, "-", ebr_part.Part_start +  ebr_part.Part_size, "|",  ebr_part.Part_size)

	idavl := "EX" + strconv.Itoa(int(index))

	graf_extendidas_tab = idavl + "[label=<\n" +
		"<TABLE ALIGN=\"LEFT\">\n" +
		"<tr>\n" +
		"	<TD>Nombre</TD>\n" +
		"	<TD>Valor</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Part_status</TD>\n" +
		"	<TD>" + string(ebr_part.Part_status) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Part_fit</TD>\n" +
		"	<TD>" + string(ebr_part.Part_fit) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Part_start</TD>\n" +
		"	<TD>" + strconv.FormatUint(ebr_part.Part_start, 10) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Part_size</TD>\n" +
		"	<TD>" + strconv.FormatUint(ebr_part.Part_size, 10) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Part_next</TD>\n" +
		"	<TD>" + strconv.Itoa(int(ebr_part.Part_next)) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Part_name</TD>\n" +
		"	<TD>" + NameByteToString(ebr_part.Part_name) + "</TD>\n" +
		"</tr>\n"

	graf_extendidas_tab = graf_extendidas_tab + "</TABLE>\n" +
		">];\n"

		/*libre en logicas*/
		/*var libre_log uint64
		if ebr_part.Part_next == -1 {
			if (ebr_part.Part_status == 'N') {
				libre_log = fin_extendida - ( m.Mbr_Particiones[num_part_extend].Part_start)
			} else {
				libre_log = fin_extendida - (ebr_part.Part_start +  ebr_part.Part_size)
			}
		} else {

			if (ebr_part.Part_status == 'N') {
				libre_log = uint64(ebr_part.Part_next) - ( m.Mbr_Particiones[num_part_extend].Part_start)
			} else {
				libre_log = uint64(ebr_part.Part_next) - (ebr_part.Part_start +  ebr_part.Part_size)
			}
		}
		fmt.Println("                    -libre: ", libre_log)*/

		/*if (ebr_part.Part_status == 'N') {

			graf_extendidas = graf_extendidas +"	<td>EBR</td>\n"

		} else if (ebr_part.Part_status == 'S' ) {

			graf_extendidas = graf_extendidas +"	<td>EBR</td>\n"+
			"	<td>"+ NameByteToString(ebr_part.Part_name) + "</td>\n"
		}*/

		/*if (libre_log > 0){
			graf_extendidas = graf_extendidas +"	<td>Libre</td>\n"
		}*/

	//ebr_anterior := ebr_part
	//fmt.Printf("						--anterior %s\n", ebr_anterior.Part_name)
	tempo_sig := ebr_part.Part_next
	for tempo_sig != -1 {

		file.Seek(int64(tempo_sig), 0)
		ebr_tem := EBR{}
		size_tem := int(binary.Size(ebr_tem))
		data_tem := readBytesDisk(file, size_tem)
		buffer_tem := bytes.NewBuffer(data_tem)
		err = binary.Read(buffer_tem, binary.BigEndian, &ebr_tem)
		if err != nil {
			log.Fatal("binary. se encontro error al leer EBR de la Extendida", err)
			break
		}

		/*graf_extendidas = graf_extendidas +"	<td>EBR</td>\n"+
		"	<td>"+ NameByteToString(ebr_tem.Part_name) + "</td>\n"*/
		index++
		//	fmt.Println("index", index)

		//fmt.Printf("--n %s\n", ebr_tem.Part_name)
		//fmt.Println("ebr", ebr_tem.Part_start - uint64(size_tem), "-", ebr_tem.Part_start)
		//fmt.Println(" ", ebr_tem.Part_start, "-", ebr_tem.Part_start + ebr_tem.Part_size, "|",  ebr_tem.Part_size)

		idavl = "EX" + strconv.Itoa(int(index))
		graf_extendidas_tab = graf_extendidas_tab + idavl + "[label=<\n" +
			"<TABLE ALIGN=\"LEFT\">\n" +
			"<tr>\n" +
			"	<TD>Nombre</TD>\n" +
			"	<TD>Valor</TD>\n" +
			"</tr>\n" +

			"<tr>\n" +
			"	<TD>Part_status</TD>\n" +
			"	<TD>" + string(ebr_tem.Part_status) + "</TD>\n" +
			"</tr>\n" +

			"<tr>\n" +
			"	<TD>Part_fit</TD>\n" +
			"	<TD>" + string(ebr_tem.Part_fit) + "</TD>\n" +
			"</tr>\n" +

			"<tr>\n" +
			"	<TD>Part_start</TD>\n" +
			"	<TD>" + strconv.FormatUint(ebr_tem.Part_start, 10) + "</TD>\n" +
			"</tr>\n" +

			"<tr>\n" +
			"	<TD>Part_size</TD>\n" +
			"	<TD>" + strconv.FormatUint(ebr_tem.Part_size, 10) + "</TD>\n" +
			"</tr>\n" +

			"<tr>\n" +
			"	<TD>Part_next</TD>\n" +
			"	<TD>" + strconv.Itoa(int(ebr_tem.Part_next)) + "</TD>\n" +
			"</tr>\n" +

			"<tr>\n" +
			"	<TD>Part_name</TD>\n" +
			"	<TD>" + NameByteToString(ebr_tem.Part_name) + "</TD>\n" +
			"</tr>\n"

		graf_extendidas_tab = graf_extendidas_tab + "</TABLE>\n" +
			">];\n"

		/*libre en logicas*/

		/*var libre_log uint64
		if ebr_tem.Part_next == -1 {
			libre_log = fin_extendida - (ebr_tem.Part_start +  ebr_tem.Part_size)
		} else {
			libre_log = uint64(ebr_tem.Part_next) - (ebr_tem.Part_start +  ebr_tem.Part_size)
		}
		fmt.Println("                    libre: ", libre_log)

		if (libre_log > 0){
			graf_extendidas = graf_extendidas +"	<td>Libre</td>\n"
		}*/

		tempo_sig = ebr_tem.Part_next
		//ebr_anterior = ebr_tem
		//fmt.Printf("						--anterior %s\n", ebr_anterior.Part_name)

		/*******************fin read sig**************************/
	}

	//////////////////}

	return graf_extendidas_tab
}
func GraficDISK(path_disco string, path_save string) {

	var graf_mbr string

	var graf_extendidas string
	file, err := os.OpenFile(path_disco, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	m := MBR{}
	var size int = int(binary.Size(m))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary. se encontro error al leer archivo binario", err)
	}

	////////////////////////////////////////////////////////////////////
	var inicial uint64
	inicial = uint64(size)
	var particiones_disco []EspacioLibre

	hayExtendida := false
	var ocupado_libre []EspacioLibre
	var num_part_extend int
	for i := 0; i < len(m.Mbr_Particiones); i++ {
		fmt.Printf("%s ", m.Mbr_Particiones[i].Part_name)
		fmt.Println("----part", i, string(m.Mbr_Particiones[i].Part_status), string(m.Mbr_Particiones[i].Part_type))

		/*ver rango de ocupado*/
		if m.Mbr_Particiones[i].Part_status == 'S' {

			ocp := EspacioLibre{}
			ocp.free_start = m.Mbr_Particiones[i].Part_start
			ocp.free_end = m.Mbr_Particiones[i].Part_start + m.Mbr_Particiones[i].Part_size
			ocp.free_size = m.Mbr_Particiones[i].Part_size
			ocupado_libre = append(ocupado_libre, ocp)

			p_oc := EspacioLibre{}
			p_oc.free_start = m.Mbr_Particiones[i].Part_start
			p_oc.free_end = m.Mbr_Particiones[i].Part_start + m.Mbr_Particiones[i].Part_size
			p_oc.free_size = m.Mbr_Particiones[i].Part_size
			p_oc.P_type = m.Mbr_Particiones[i].Part_type
			p_oc.P_name = m.Mbr_Particiones[i].Part_name
			particiones_disco = append(particiones_disco, p_oc)

		}
		/*si hay una extendida*/
		if m.Mbr_Particiones[i].Part_type == 'e' {
			hayExtendida = true
			num_part_extend = i
		}
	}

	libre := OrdenmMLLenos(ocupado_libre)
	//var espacios_libres []EspacioLibre
	if len(libre) == 0 {
		//fmt.Println("--             libre:", inicial, "-", m.Mbr_tamanio, "|", m.Mbr_tamanio - inicial)
		lib := EspacioLibre{}
		lib.free_start = inicial
		lib.free_end = m.Mbr_tamanio
		lib.free_size = m.Mbr_tamanio - inicial

		lib.P_type = 'f'
		copy(lib.P_name[:], "Libre")
		particiones_disco = append(particiones_disco, lib)
	}

	for i := 0; i < len(libre); i++ {
		//fmt.Println("***", libre[i].free_start, "-", libre[i].free_end /*,"|" , libre[i].free_size*/)
		if i == 0 {
			if inicial != libre[i].free_start {
				lib := EspacioLibre{}
				lib.free_start = inicial
				lib.free_end = libre[i].free_start
				lib.free_size = libre[i].free_start - inicial

				lib.P_type = 'f'
				copy(lib.P_name[:], "Libre")
				particiones_disco = append(particiones_disco, lib)
				//fmt.Println("0             libre:", inicial, "-", libre[i].free_start ,"|" ,  libre[i].free_start - inicial)
			}
		}
		if (i + 1) < len(libre) {
			if libre[i].free_end != libre[i+1].free_start {
				lib := EspacioLibre{}
				lib.free_start = libre[i].free_end
				lib.free_end = libre[i+1].free_start
				lib.free_size = libre[i+1].free_start - libre[i].free_end

				lib.P_type = 'f'
				copy(lib.P_name[:], "Libre")
				particiones_disco = append(particiones_disco, lib)
				//fmt.Println("1             libre:", libre[i].free_end, "-", libre[i+1].free_start ,"|" , libre[i+1].free_start - libre[i].free_end)
			}
		} else {

			if libre[i].free_end != m.Mbr_tamanio {
				lib := EspacioLibre{}
				lib.free_start = libre[i].free_end
				lib.free_end = m.Mbr_tamanio
				lib.free_size = m.Mbr_tamanio - libre[i].free_end

				lib.P_type = 'f'
				copy(lib.P_name[:], "Libre")
				particiones_disco = append(particiones_disco, lib)
				//fmt.Println("2              libre:", libre[i].free_end, "-", m.Mbr_tamanio ,"|" , m.Mbr_tamanio - libre[i].free_end)
			}
		}
	}

	particiones_disco = OrdenmMLLenos(particiones_disco)
	/////////////////////////////////////////////////////////////
	if hayExtendida == true {
		graf_extendidas = Graf_extendidas(file, m, num_part_extend)
	}

	/*graf_extendidas = "<TABLE ALIGN=\"LEFT\">\n"+
	"<tr>\n"+
		"	<td>Extendida</td>\n"+
	"</tr>\n"+

	"<tr>\n"+
	"	<td>"+

	"<TABLE ALIGN=\"LEFT\">\n"+
	"<tr>\n"+
		"	<td>EBR</td>\n"+
		"	<td>logica1</td>\n"+
	"</tr>\n"+
	"</TABLE>\n"+

	"</td>\n"+
	"</tr>\n"+

	"</TABLE>\n"*/

	graf_mbr = "digraph test {\n" +
		"graph [ratio=fill];\n" +
		"node [label=\"\\N\", fontsize=12, shape=plaintext];\n" +
		"graph [bb=\"2,2,362,164\"];\n" +
		"arset[label=<\n" +
		"<TABLE ALIGN=\"LEFT\">\n" +

		"<tr>\n" +
		"	<TD>MBR</TD>\n"

	var name_part string
	var tipo_par string
	for i := 0; i < len(particiones_disco); i++ {
		fmt.Printf("%s ", particiones_disco[i].P_name)
		fmt.Println(i, string(particiones_disco[i].P_type))

		if particiones_disco[i].P_type == 'p' {
			tipo_par = "(Primaria)"
		} else if particiones_disco[i].P_type == 'e' {
			tipo_par = "(Extendida)"
		} else if particiones_disco[i].P_type == 'f' {
			tipo_par = ""
		}

		name_part = ""
		for c := 0; c < len(particiones_disco[i].P_name); c++ {
			if particiones_disco[i].P_name[c] == 0 {
				name_part = name_part + " "
			} else {
				name_part = name_part + string(particiones_disco[i].P_name[c])
			}
		}

		if particiones_disco[i].P_type == 'e' {
			//graf_mbr = graf_mbr +"<TD>"+graf_extendidas +"</TD>\n"
			graf_mbr = graf_mbr + "<TD>" +
				//graf_extendidas +
				"<TABLE ALIGN=\"LEFT\">\n" +
				"<tr>\n" +
				"	<TD>" + tipo_par + name_part + "</TD>\n" +
				"</tr>\n" +

				"<tr>\n" +
				"	<td>" +

				"<TABLE ALIGN=\"LEFT\">\n" +
				"<tr>\n" +
				//	"	<td>EBR</td>\n"+
				//	"	<td>logica1</td>\n"+
				graf_extendidas +
				"</tr>\n" +
				"</TABLE>\n" +

				"</td>\n" +
				"</tr>\n" +

				"</TABLE>\n" +

				"</TD>\n"

		} else {
			graf_mbr = graf_mbr + "	<TD>" + tipo_par + name_part + "</TD>\n"
		}

		if particiones_disco[i].P_type == 'e' {
			//graf_mbr = graf_mbr /*+ "<tr>" */+  graf_extendidas /*+"</tr>\n"*/
		}

	}

	graf_mbr = graf_mbr + "</tr>\n"
	graf_mbr = graf_mbr + "</TABLE>\n" +
		">, ];\n" +
		"}"

	EscribirImagen(graf_mbr, path_save)
}

func GraficDISK_back(path_disco string, path_save string) {

	var graf_mbr string

	var graf_extendidas string
	file, err := os.OpenFile(path_disco, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	m := MBR{}
	var size int = int(binary.Size(m))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary. se encontro error al leer archivo binario", err)
	}

	////////////////////////////////////////////////////////////////////
	var inicial uint64
	inicial = uint64(size)
	var particiones_disco []EspacioLibre

	//hayExtendida := false
	var ocupado_libre []EspacioLibre
	//var ext_i int
	for i := 0; i < len(m.Mbr_Particiones); i++ {
		fmt.Printf("%s ", m.Mbr_Particiones[i].Part_name)
		fmt.Println("----part", i, string(m.Mbr_Particiones[i].Part_status), string(m.Mbr_Particiones[i].Part_type))

		/*ver rango de ocupado*/
		if m.Mbr_Particiones[i].Part_status == 'S' {

			ocp := EspacioLibre{}
			ocp.free_start = m.Mbr_Particiones[i].Part_start
			ocp.free_end = m.Mbr_Particiones[i].Part_start + m.Mbr_Particiones[i].Part_size
			ocp.free_size = m.Mbr_Particiones[i].Part_size
			ocupado_libre = append(ocupado_libre, ocp)

			p_oc := EspacioLibre{}
			p_oc.free_start = m.Mbr_Particiones[i].Part_start
			p_oc.free_end = m.Mbr_Particiones[i].Part_start + m.Mbr_Particiones[i].Part_size
			p_oc.free_size = m.Mbr_Particiones[i].Part_size
			p_oc.P_type = m.Mbr_Particiones[i].Part_type
			p_oc.P_name = m.Mbr_Particiones[i].Part_name
			particiones_disco = append(particiones_disco, p_oc)

		}
		/*si hay una extendida*/
		/*if m.Mbr_Particiones[i].Part_type == 'e' {
			hayExtendida = true
			ext_i = i
		}*/
	}

	libre := OrdenmMLLenos(ocupado_libre)
	//var espacios_libres []EspacioLibre
	if len(libre) == 0 {
		//fmt.Println("--             libre:", inicial, "-", m.Mbr_tamanio, "|", m.Mbr_tamanio - inicial)
		lib := EspacioLibre{}
		lib.free_start = inicial
		lib.free_end = m.Mbr_tamanio
		lib.free_size = m.Mbr_tamanio - inicial

		lib.P_type = 'f'
		copy(lib.P_name[:], "Libre")
		particiones_disco = append(particiones_disco, lib)
	}

	for i := 0; i < len(libre); i++ {
		//fmt.Println("***", libre[i].free_start, "-", libre[i].free_end /*,"|" , libre[i].free_size*/)
		if i == 0 {
			if inicial != libre[i].free_start {
				lib := EspacioLibre{}
				lib.free_start = inicial
				lib.free_end = libre[i].free_start
				lib.free_size = libre[i].free_start - inicial

				lib.P_type = 'f'
				copy(lib.P_name[:], "Libre")
				particiones_disco = append(particiones_disco, lib)
				//fmt.Println("0             libre:", inicial, "-", libre[i].free_start ,"|" ,  libre[i].free_start - inicial)
			}
		}
		if (i + 1) < len(libre) {
			if libre[i].free_end != libre[i+1].free_start {
				lib := EspacioLibre{}
				lib.free_start = libre[i].free_end
				lib.free_end = libre[i+1].free_start
				lib.free_size = libre[i+1].free_start - libre[i].free_end

				lib.P_type = 'f'
				copy(lib.P_name[:], "Libre")
				particiones_disco = append(particiones_disco, lib)
				//fmt.Println("1             libre:", libre[i].free_end, "-", libre[i+1].free_start ,"|" , libre[i+1].free_start - libre[i].free_end)
			}
		} else {

			if libre[i].free_end != m.Mbr_tamanio {
				lib := EspacioLibre{}
				lib.free_start = libre[i].free_end
				lib.free_end = m.Mbr_tamanio
				lib.free_size = m.Mbr_tamanio - libre[i].free_end

				lib.P_type = 'f'
				copy(lib.P_name[:], "Libre")
				particiones_disco = append(particiones_disco, lib)
				//fmt.Println("2              libre:", libre[i].free_end, "-", m.Mbr_tamanio ,"|" , m.Mbr_tamanio - libre[i].free_end)
			}
		}
	}

	particiones_disco = OrdenmMLLenos(particiones_disco)
	/////////////////////////////////////////////////////////////

	/*graf_extendidas = "<TABLE ALIGN=\"LEFT\">\n"+
	"<tr>\n"+
		"	<td>Extendida</td>\n"+
	"</tr>\n"+

	"<tr>\n"+
		"	<td>EBR2</td>\n"+
		"	<td>parti 222</td>\n"+
	"</tr>\n"+

	"</TABLE>\n"*/

	graf_extendidas = "<TABLE ALIGN=\"LEFT\">\n" +
		"<tr>\n" +
		"	<td>Extendida</td>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<td>" +

		"<TABLE ALIGN=\"LEFT\">\n" +
		"<tr>\n" +
		"	<td>EBR</td>\n" +
		"	<td>logica1</td>\n" +
		"</tr>\n" +
		"</TABLE>\n" +

		"</td>\n" +
		"</tr>\n" +

		"</TABLE>\n"

	graf_mbr = "digraph test {\n" +
		"graph [ratio=fill];\n" +
		"node [label=\"\\N\", fontsize=12, shape=plaintext];\n" +
		"graph [bb=\"2,2,362,164\"];\n" +
		"arset[label=<\n" +
		"<TABLE ALIGN=\"LEFT\">\n" +

		"<tr>\n" +
		"	<TD>MBR</TD>\n"

	var name_part string
	var tipo_par string
	for i := 0; i < len(particiones_disco); i++ {
		fmt.Printf("%s ", particiones_disco[i].P_name)
		fmt.Println(i, string(particiones_disco[i].P_type))

		if particiones_disco[i].P_type == 'p' {
			tipo_par = "(Primaria)"
		} else if particiones_disco[i].P_type == 'e' {
			tipo_par = "(Extendida)"
		} else if particiones_disco[i].P_type == 'f' {
			tipo_par = ""
		}

		name_part = ""
		for c := 0; c < len(particiones_disco[i].P_name); c++ {
			if particiones_disco[i].P_name[c] == 0 {
				name_part = name_part + " "
			} else {
				name_part = name_part + string(particiones_disco[i].P_name[c])
			}
		}

		graf_mbr = graf_mbr + "	<TD>" + tipo_par + name_part + "</TD>\n"

		if particiones_disco[i].P_type == 'e' {
			graf_mbr = graf_mbr + "<TD>" + graf_extendidas + "</TD>\n"
		}

		if particiones_disco[i].P_type == 'e' {
			//graf_mbr = graf_mbr /*+ "<tr>" */+  graf_extendidas /*+"</tr>\n"*/
		}

	}

	graf_mbr = graf_mbr + "</tr>\n"
	graf_mbr = graf_mbr + "</TABLE>\n" +
		">, ];\n" +
		"}"

	EscribirImagen(graf_mbr, path_save)
}

func GraficMBR(path_disco string, path_save string) {

	var graf_mbr string
	file, err := os.OpenFile(path_disco, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	m := MBR{}
	var size int = int(binary.Size(m))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary. se encontro error al leer archivo binario", err)
	}
	///////////fmt.Println(m)
	graf_mbr = "digraph test {\n" +
		"graph [ratio=fill];\n" +
		"node [label=\"\\N\", fontsize=15, shape=plaintext];\n" +
		//"graph [bb=\"0,0,352,154\"];\n"+
		"arset[label=<\n" +
		"<TABLE ALIGN=\"LEFT\">\n" +
		"<tr>\n" +
		"	<TD>Nombre</TD>\n" +
		"	<TD>Valor</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Mbr_tamanio</TD>\n" +
		"	<TD>" + strconv.FormatUint(m.Mbr_tamanio, 10) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Mbr_fecha_creacion_s</TD>\n" +
		"	<TD>" + string(m.Mbr_fecha_creacion_s[:]) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Mbr_disk_signature</TD>\n" +
		"	<TD>" + strconv.FormatUint(m.Mbr_disk_signature, 10) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Particiones</TD>\n" +
		"</tr>\n"

	for i := 0; i < len(m.Mbr_Particiones); i++ {
		//fmt.Print( "type-",m.Mbr_Particiones[i].Part_type, "-\n")

		graf_mbr = graf_mbr + "<tr>\n" +
			"	<td>Part_status" + strconv.Itoa(i) + "</td>\n" +
			"	<td> " + string(m.Mbr_Particiones[i].Part_status) + "</td>\n" +
			"</tr>\n"

		if m.Mbr_Particiones[i].Part_status == 'S' {

			//name_part := string(m.Mbr_Particiones[i].Part_name[:])
			//name_part := bytes.NewBuffer(m.Mbr_Particiones[i].Part_name).String()
			var name_part string
			name_part = ""
			for c := 0; c < len(m.Mbr_Particiones[i].Part_name); c++ {

				if m.Mbr_Particiones[i].Part_name[c] == 0 {
					name_part = name_part + " "
				} else {
					name_part = name_part + string(m.Mbr_Particiones[i].Part_name[c])
				}
			}

			fmt.Println(name_part)
			graf_mbr = graf_mbr + "<tr>\n" +
				"	<td>Part_type" + strconv.Itoa(i) + "</td>\n" +
				"	<td> " + string(m.Mbr_Particiones[i].Part_type) + "</td>\n" +
				"</tr>\n" +

				"<tr>\n" +
				"	<td>Part_fit" + strconv.Itoa(i) + "</td>\n" +
				"	<td>" + string(m.Mbr_Particiones[i].Part_fit) + "</td>\n" +
				"</tr>\n" +

				"<tr>\n" +
				"	<td>Part_start" + strconv.Itoa(i) + "</td>\n" +
				"	<td>" + strconv.FormatUint(m.Mbr_Particiones[i].Part_start, 10) + "</td>\n" +
				"</tr>\n" +

				"<tr>\n" +
				"	<td>Part_size" + strconv.Itoa(i) + "</td>\n" +
				"	<td>" + strconv.FormatUint(m.Mbr_Particiones[i].Part_size, 10) + "</td>\n" +
				"</tr>\n" +

				"<tr>\n" +
				"	<td>Part_name" + strconv.Itoa(i) + "</td>\n" +
				"	<td>" + name_part + "</td>\n" +
				"</tr>\n"

		} else {

			graf_mbr = graf_mbr + "<tr>\n" +
				"	<td>Part_type" + strconv.Itoa(i) + "</td>\n" +
				"	<td></td>\n" +
				"</tr>\n" +

				"<tr>\n" +
				"	<td>Part_fit" + strconv.Itoa(i) + "</td>\n" +
				"	<td></td>\n" +
				"</tr>\n" +

				"<tr>\n" +
				"	<td>Part_start" + strconv.Itoa(i) + "</td>\n" +
				"	<td></td>\n" +
				"</tr>\n" +

				"<tr>\n" +
				"	<td>Part_size" + strconv.Itoa(i) + "</td>\n" +
				"	<td></td>\n" +
				"</tr>\n" +

				"<tr>\n" +
				"	<td>Part_name" + strconv.Itoa(i) + "</td>\n" +
				"	<td></td>\n" +
				"</tr>\n"

		}
	}

	graf_mbr = graf_mbr + "<tr>\n" +
		"	<TD></TD>\n" +
		"</tr>\n"

	graf_mbr = graf_mbr + "</TABLE>\n" +
		">, ];\n" +
		"}"

	EscribirImagen(graf_mbr, path_save)

}

func GraficMBR_new(path_disco string, path_save string) {

	var graf_mbr string
	file, err := os.OpenFile(path_disco, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	m := MBR{}
	var size int = int(binary.Size(m))
	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary. se encontro error al leer archivo binario", err)
	}

	var graf_extend string = ""
	graf_extend = "digraph mbr_todo {\n" +
		"node [shape=plaintext]\n" +
		"rankdir=LR;\n"

	///////////fmt.Println(m)
	/*graf_mbr = "digraph test {\n" +
	"graph [ratio=fill];\n"+
	"node [label=\"\\N\", fontsize=15, shape=plaintext];\n"*/
	graf_mbr = "arset[label=<\n" +
		"<TABLE ALIGN=\"LEFT\">\n" +
		"<tr>\n" +
		"	<TD>Nombre</TD>\n" +
		"	<TD>Valor</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Mbr_tamanio</TD>\n" +
		"	<TD>" + strconv.FormatUint(m.Mbr_tamanio, 10) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Mbr_fecha_creacion_s</TD>\n" +
		"	<TD>" + string(m.Mbr_fecha_creacion_s[:]) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Mbr_disk_signature</TD>\n" +
		"	<TD>" + strconv.FormatUint(m.Mbr_disk_signature, 10) + "</TD>\n" +
		"</tr>\n" +

		"<tr>\n" +
		"	<TD>Particiones</TD>\n" +
		"</tr>\n"

	for i := 0; i < len(m.Mbr_Particiones); i++ {
		//fmt.Print( "type-",m.Mbr_Particiones[i].Part_type, "-\n")

		graf_mbr = graf_mbr + "<tr>\n" +
			"	<td>Part_status" + strconv.Itoa(i) + "</td>\n" +
			"	<td> " + string(m.Mbr_Particiones[i].Part_status) + "</td>\n" +
			"</tr>\n"

		if m.Mbr_Particiones[i].Part_status == 'S' {

			//name_part := string(m.Mbr_Particiones[i].Part_name[:])
			//name_part := bytes.NewBuffer(m.Mbr_Particiones[i].Part_name).String()
			var name_part string
			name_part = ""
			for c := 0; c < len(m.Mbr_Particiones[i].Part_name); c++ {

				if m.Mbr_Particiones[i].Part_name[c] == 0 {
					name_part = name_part + " "
				} else {
					name_part = name_part + string(m.Mbr_Particiones[i].Part_name[c])
				}
			}

			fmt.Println(name_part)
			graf_mbr = graf_mbr + "<tr>\n" +
				"	<td>Part_type" + strconv.Itoa(i) + "</td>\n" +
				"	<td> " + string(m.Mbr_Particiones[i].Part_type) + "</td>\n" +
				"</tr>\n" +

				"<tr>\n" +
				"	<td>Part_fit" + strconv.Itoa(i) + "</td>\n" +
				"	<td>" + string(m.Mbr_Particiones[i].Part_fit) + "</td>\n" +
				"</tr>\n" +

				"<tr>\n" +
				"	<td>Part_start" + strconv.Itoa(i) + "</td>\n" +
				"	<td>" + strconv.FormatUint(m.Mbr_Particiones[i].Part_start, 10) + "</td>\n" +
				"</tr>\n" +

				"<tr>\n" +
				"	<td>Part_size" + strconv.Itoa(i) + "</td>\n" +
				"	<td>" + strconv.FormatUint(m.Mbr_Particiones[i].Part_size, 10) + "</td>\n" +
				"</tr>\n" +

				"<tr>\n" +
				"	<td>Part_name" + strconv.Itoa(i) + "</td>\n" +
				"	<td>" + name_part + "</td>\n" +
				"</tr>\n"

		} else {

			graf_mbr = graf_mbr + "<tr>\n" +
				"	<td>Part_type" + strconv.Itoa(i) + "</td>\n" +
				"	<td></td>\n" +
				"</tr>\n" +

				"<tr>\n" +
				"	<td>Part_fit" + strconv.Itoa(i) + "</td>\n" +
				"	<td></td>\n" +
				"</tr>\n" +

				"<tr>\n" +
				"	<td>Part_start" + strconv.Itoa(i) + "</td>\n" +
				"	<td></td>\n" +
				"</tr>\n" +

				"<tr>\n" +
				"	<td>Part_size" + strconv.Itoa(i) + "</td>\n" +
				"	<td></td>\n" +
				"</tr>\n" +

				"<tr>\n" +
				"	<td>Part_name" + strconv.Itoa(i) + "</td>\n" +
				"	<td></td>\n" +
				"</tr>\n"

		}
	}

	graf_mbr = graf_mbr + "<tr>\n" +
		"	<TD></TD>\n" +
		"</tr>\n"

	graf_mbr = graf_mbr + "</TABLE>\n" +
		">, ];\n"

	graf_extend = graf_extend + graf_mbr

	hayExtendida := false
	var num_part_extend int
	for i := 0; i < len(m.Mbr_Particiones); i++ {
		/*si hay una extendida*/
		if m.Mbr_Particiones[i].Part_type == 'e' {
			hayExtendida = true
			num_part_extend = i
		}
	}

	var graf_ex_tab string = ""
	if hayExtendida == true {
		graf_ex_tab = EBR_extendidas(file, m, num_part_extend)
	}

	graf_extend = graf_extend + graf_ex_tab + "}\n"
	EscribirImagen(graf_extend, path_save)

}
func EscribirImagen(image string, path_save string) {

	//name_com := "/home/rafaelc/Dis/" + "mbr_gra.txt"

	//name_txt:= name_imagen + ".txt"
	//path_sav := path_save + name_txt

	//fmt.Println("name_com :",name_com )

	file, err := os.Create(path_save + ".txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(image)
	if err != nil {
		log.Fatal("err ", err)
	}

	ExeCommand(path_save)

}

func ExeCommand(path_save string) {
	//cmd := exec.Command("/home/rafaelc/Dis/","dot -Tpng mbr_gra.txt -o mbr_gra.png")
	//cmd.Stdin = strings.NewReader("some input")

	//cmd := exec.Command("dot", "-Tpng", "mbr_gra.txt", "-o", "mbr_gra.png")
	cmd := exec.Command("dot", "-Tpng", path_save+".txt", "-o", path_save)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf( out)

	fmt.Printf("%q", out.String())
	/*cmd = exec.Command(path_save)
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}*/

}

/*leyendo archivo, disco*/
func readDisk_new(path string, btype byte, fit byte, size_part uint64, name string) {
	/*abrimos el archivo*/
	///file, err := os.Open(path)
	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	m := MBR{}
	//var size int = int(unsafe.Sizeof(m))
	//fmt.Println("sizeof :", size )

	var size int = int(binary.Size(m))
	///////var part_ini uint64
	fmt.Println("binary :", size)

	/*175*/
	//////part_ini = uint64(size) + 1
	//////fmt.Println("part_ini :", part_ini )

	data := readBytesDisk(file, size)
	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary. se encontro error al leer archivo binario", err)
	}

	//fmt.Println(m)
	/*verificando particiones libres del mbr*/
	fmt.Println("tam_disk", m.Mbr_tamanio)

	/*convertiendo en bytes el nombre*/
	var name_byte [16]byte
	copy(name_byte[:], name)
	//fmt.Println("name_byte", name_byte)

	hayExtendida := false
	var libre []EspacioLibre
	var ext_i int
	for i := 0; i < len(m.Mbr_Particiones); i++ {
		//fmt.Printf( "%s ", m.Mbr_Particiones[i].Part_name)
		//fmt.Println("part", i, string(m.Mbr_Particiones[i].Part_status ), string(m.Mbr_Particiones[i].Part_type))

		/*ver rango de ocupado*/
		if m.Mbr_Particiones[i].Part_status == 'S' {
			//discovacio = false
			lib := EspacioLibre{}
			lib.free_start = m.Mbr_Particiones[i].Part_start
			lib.free_end = m.Mbr_Particiones[i].Part_start + m.Mbr_Particiones[i].Part_size
			lib.free_size = m.Mbr_Particiones[i].Part_size
			libre = append(libre, lib)

			/*existe nombre, verfico que tiene que ser igual, mayus, minus*/

			/*verifico NOMBRE PRIMARIAS O EXTERNA*/
			if m.Mbr_Particiones[i].Part_name == name_byte {
				//fmt.Println("El nombre del Disco", name_byte, "ya Existe")
				fmt.Printf("Ya Existe una Particion con el mismo nombre en Disco | %s \n", name_byte)
				return
			}

		}
		/*si hay una extendida*/
		if m.Mbr_Particiones[i].Part_type == 'e' {
			hayExtendida = true
			ext_i = i
		}
	}

	//////
	//return

	/***inicio de validacion de tipo**/
	if hayExtendida == true && btype == 'e' {
		fmt.Println("No puede crear una partición EXTENDIDA, ya existe uno")
		return
	}
	if hayExtendida == false && btype == 'l' {
		fmt.Println("No puede crear una partición LOGICA, porque no existe Extendida")
		return
	}

	if btype != 'l' && m.Mbr_Particiones[0].Part_status == 'S' && m.Mbr_Particiones[1].Part_status == 'S' && m.Mbr_Particiones[2].Part_status == 'S' && m.Mbr_Particiones[3].Part_status == 'S' {
		fmt.Println("Las 4 Particiones ya estan ocupadas")
		return
	}

	/*verifico NOMBRE en LOGICAS*/
	if hayExtendida == true {

		/*verificando si exite particion logica con el mismo nombre*/
		exislog := exist_partlogicas(file, name, m, ext_i)
		if exislog == true {
			fmt.Printf("Ya Existe una Particion LOGICA con el mismo nombre en Disco | %s \n", name_byte)
			return
		}
		//fmt.Println("exislog:", exislog)

	}

	/***fin de validacion de tipo**/

	/*****************************INICIO*SI*ES*TIPO**PRIMARIO*O*EXTENDIDA**************************************************/
	if btype == 'p' || btype == 'e' {
		/*verificando espacio libre y usados*/
		//var parcerc uint64
		var inicial uint64
		inicial = uint64(size) //+ 1
		////fmt.Println(inicial)
		for i := 0; i < len(m.Mbr_Particiones); i++ {
			fmt.Printf("%s ", m.Mbr_Particiones[i].Part_name)
			fmt.Println(i, string(m.Mbr_Particiones[i].Part_status), string(m.Mbr_Particiones[i].Part_type), m.Mbr_Particiones[i].Part_start, "-", m.Mbr_Particiones[i].Part_start+m.Mbr_Particiones[i].Part_size, "|", m.Mbr_Particiones[i].Part_size)
			///obtengo el valor menor

		}

		/////fmt.Println("            incia libre:", ini_libre)
		/////fmt.Println("------------")

		fmt.Println(" llen:", len(libre))

		for i := 0; i < len(libre); i++ {
			fmt.Println(libre[i].free_start, "-", libre[i].free_end, "|", libre[i].free_size)
		}

		libre = OrdenmMLLenos(libre)
		var espacios_libres []EspacioLibre
		if len(libre) == 0 {
			fmt.Println("--             libre:", inicial, "-", m.Mbr_tamanio, "|", m.Mbr_tamanio-inicial)

			lib := EspacioLibre{}
			lib.free_start = inicial
			lib.free_end = m.Mbr_tamanio
			lib.free_size = m.Mbr_tamanio - inicial
			espacios_libres = append(espacios_libres, lib)
		}

		for i := 0; i < len(libre); i++ {
			fmt.Println("***", libre[i].free_start, "-", libre[i].free_end /*,"|" , libre[i].free_size*/)

			if i == 0 {
				if inicial != libre[i].free_start {

					lib := EspacioLibre{}
					lib.free_start = inicial
					lib.free_end = libre[i].free_start
					lib.free_size = libre[i].free_start - inicial
					espacios_libres = append(espacios_libres, lib)
					fmt.Println("0             libre:", inicial, "-", libre[i].free_start, "|", libre[i].free_start-inicial)
				}
			}
			if (i + 1) < len(libre) {
				if libre[i].free_end != libre[i+1].free_start {

					lib := EspacioLibre{}
					lib.free_start = libre[i].free_end
					lib.free_end = libre[i+1].free_start
					lib.free_size = libre[i+1].free_start - libre[i].free_end
					espacios_libres = append(espacios_libres, lib)

					fmt.Println("1             libre:", libre[i].free_end, "-", libre[i+1].free_start, "|", libre[i+1].free_start-libre[i].free_end)
				}
			} else {

				if libre[i].free_end != m.Mbr_tamanio {

					lib := EspacioLibre{}
					lib.free_start = libre[i].free_end
					lib.free_end = m.Mbr_tamanio
					lib.free_size = m.Mbr_tamanio - libre[i].free_end
					espacios_libres = append(espacios_libres, lib)
					fmt.Println("2              libre:", libre[i].free_end, "-", m.Mbr_tamanio, "|", m.Mbr_tamanio-libre[i].free_end)
				}
			}

		}

		fmt.Println("\n espacios_libres:", len(espacios_libres))
		fmt.Println("insert:", size_part)
		var start_new_part uint64
		part_encontrada := false
		for i := 0; i < len(espacios_libres); i++ {
			fmt.Println(i, "--", espacios_libres[i].free_start, "-", espacios_libres[i].free_end, "|", espacios_libres[i].free_size)
			/*verificando si la particion caba en el primer espacio*/
			if size_part <= espacios_libres[i].free_size {
				fmt.Println("----Se encontro particion en:", i)
				start_new_part = espacios_libres[i].free_start
				part_encontrada = true
				break
			}
		}

		if part_encontrada == false {
			fmt.Println("No se encontró sufienete espacio para la particion")
			return
		}

		var num_part int
		for i := 0; i < len(m.Mbr_Particiones); i++ {
			if m.Mbr_Particiones[i].Part_status == 'N' {
				fmt.Println("part seleccionado:", i)
				num_part = i
				break
			}
		}

		/*actulizando datos de MBR*/
		m.Mbr_Particiones[num_part].Part_status = 'S'
		m.Mbr_Particiones[num_part].Part_type = btype
		m.Mbr_Particiones[num_part].Part_fit = fit
		m.Mbr_Particiones[num_part].Part_start = start_new_part
		m.Mbr_Particiones[num_part].Part_size = size_part
		m.Mbr_Particiones[num_part].Part_name = name_byte
		//copy(m.Mbr_Particiones[num_part].Part_name[:], name)

		file.Seek(0, 0)
		disk_w22 := &m

		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, disk_w22)
		writeBytesDisk(file, binario2.Bytes())

		if btype == 'e' {
			/*creando EBR*/
			part_ex := EBR{}

			part_ex.Part_status = 'N'
			part_ex.Part_fit = fit
			//part_ex.Part_start uint64
			//part_ex.Part_size uint64
			part_ex.Part_next = -1
			//part_ex.Part_name [16]byte
			file.Seek(int64(start_new_part), 0)

			/*escribiendo ebr inicial en disco, particion externa*/
			disk_ebr_ext := &part_ex
			var binario_ebr bytes.Buffer
			binary.Write(&binario_ebr, binary.BigEndian, disk_ebr_ext)
			writeBytesDisk(file, binario_ebr.Bytes())

		}
	}

	/*****************************FIN*SI*ES*TIPO**PRIMARIO*O*EXTENDIDA**************************************************/
	/***inicio para logico***/
	if btype == 'l' {

		var num_part_extend int
		for i := 0; i < len(m.Mbr_Particiones); i++ {
			if m.Mbr_Particiones[i].Part_type == 'e' {
				num_part_extend = i
				////fmt.Println("part de extendida:", num_part_extend)
				break
			}
		}

		/*la siguiente particion*/
		fmt.Println("Part_status:", string(m.Mbr_Particiones[num_part_extend].Part_status))
		fmt.Println("Part_start:", m.Mbr_Particiones[num_part_extend].Part_start)
		fmt.Println("Part_size:", m.Mbr_Particiones[num_part_extend].Part_size)

		fin_extendida := m.Mbr_Particiones[num_part_extend].Part_start + m.Mbr_Particiones[num_part_extend].Part_size
		fmt.Println("fin_extendida:", fin_extendida)
		/*posicionandome en el primer ebr de la extendida*/
		file.Seek(0, 0)
		file.Seek(int64(m.Mbr_Particiones[num_part_extend].Part_start), 0)

		ebr_part := EBR{}
		size_ebr := int(binary.Size(ebr_part))
		//////fmt.Println("size_ebr:" , size_ebr)

		data = readBytesDisk(file, size_ebr)
		//Convierte la data en un buffer,necesario para
		//decodificar binario
		buffer = bytes.NewBuffer(data)

		//Decodificamos y guardamos en la variable ebr_part
		err = binary.Read(buffer, binary.BigEndian, &ebr_part)
		if err != nil {
			log.Fatal("binary. se encontro error al leer EBR de la Extendida", err)
		}
		fmt.Println(ebr_part)
		fmt.Printf("i-------nsert name %s\n", name_byte)

		fmt.Println("0ebr_part.Part_status", string(ebr_part.Part_status))
		fmt.Println("ebr_part.Part_next", ebr_part.Part_next)
		fmt.Println("ebr_part.Part_start", ebr_part.Part_start)
		fmt.Println("ebr_part.Part_size", ebr_part.Part_size)
		fmt.Printf("ebr_part.Part_name %s\n", ebr_part.Part_name)

		///*************************************************///
		fmt.Println("*************************************")
		fmt.Printf("--n %s\n", ebr_part.Part_name)
		fmt.Println("ebr", ebr_part.Part_start-uint64(size_ebr), "-", ebr_part.Part_start)
		fmt.Println(" ", ebr_part.Part_start, "-", ebr_part.Part_start+ebr_part.Part_size, "|", ebr_part.Part_size)
		/*libre en logicas*/
		var libre_log uint64
		if ebr_part.Part_next == -1 {
			if ebr_part.Part_status == 'N' {
				libre_log = fin_extendida - (m.Mbr_Particiones[num_part_extend].Part_start)
			} else {
				libre_log = fin_extendida - (ebr_part.Part_start + ebr_part.Part_size)
			}
		} else {

			if ebr_part.Part_status == 'N' {
				libre_log = uint64(ebr_part.Part_next) - (m.Mbr_Particiones[num_part_extend].Part_start)
			} else {
				libre_log = uint64(ebr_part.Part_next) - (ebr_part.Part_start + ebr_part.Part_size)
			}
		}
		fmt.Println("                    -libre: ", libre_log)
		//fmt.Println("               *ebr_part.Part_next: ", ebr_part.Part_next)
		//fmt.Println("                    -ebr_part.Part_status: ", string(ebr_part.Part_status))
		/*ingreso la nueva particion logica en el ebr de la extendida*/
		/*simulando una lista enlazada, el status N se le toma como cabecera nula*/
		if ebr_part.Part_status == 'N' {

			/*Verifico si hay espacio disponible*/
			iniciar_particion := uint64(m.Mbr_Particiones[num_part_extend].Part_start)
			start := m.Mbr_Particiones[num_part_extend].Part_start + uint64(size_ebr)
			//libre_ex := fin_extendida - start

			///fmt.Println("                    libre: ", libre_log)
			//fmt.Println("libre_ex primero", libre_ex )
			//////////fmt.Println("size_part primero", size_part)
			//ingresar := size_part + uint64(size_ebr) ****
			ingresar := size_part
			if ingresar < libre_log {
				//NewPartEBR(file *os.File, iniciar_particion uint64, status byte, fit byte, start uint64, size_part uint64, siguiente int64, name_byte [16]byte )
				NewPartEBR(file, iniciar_particion, 'S', fit, start, size_part, ebr_part.Part_next, name_byte)
				return

			} else /*if (ingresar > libre_log) */ {
				////fmt.Println("No se encontró sufienete espacio para la particion, Espacio Libre:", libre_log, "bytes")
				////fmt.Println("Espacio de Particion + EBR:", ingresar, "bytes")
				/////////return
			}
			//Fdisk -sizE->17 -path->"/home/rafaelc/Dis/Disco4.dsk"  -name->Parlog6 -unit->m -type->l return
			//fmt.Println("espacio encotrado")
			//return

		} else if ebr_part.Part_status == 'S' && ebr_part.Part_next != -1 { ////////////////////////else {
			//fmt.Println(" insertar a mitad de primero ")

			ingresar := size_part + uint64(size_ebr)
			if ingresar <= libre_log && ebr_part.Part_next != -1 {
				fmt.Println("11 mitad se puede ingresar en espacio libre")
				/*insertando particion en medio*/
				iniciar_par_bmr := ebr_part.Part_start + ebr_part.Part_size
				start := iniciar_par_bmr + uint64(size_ebr)
				////NewPartEBR(file *os.File, iniciar_par_bmr uint64, status byte, fit byte, start uint64, size_part uint64, siguiente int64, name_byte [16]byte )
				fmt.Println("mitad iniciar_par_bmr", iniciar_par_bmr)
				fmt.Println("mitad start", start)

				NewPartEBR(file, iniciar_par_bmr, 'S', fit, start, size_part, ebr_part.Part_next, name_byte)

				/**********actuliza valor anteior************/
				//fmt.Printf("ebr_tem.Part_name:%s\n", ebr_tem.Part_name)
				ebr_part.Part_next = int64(iniciar_par_bmr)
				pos_ante_ebr := ebr_part.Part_start - uint64(size_ebr)
				//fmt.Println("pos_ante_ebr Primer EBR:", pos_ante_ebr)

				file.Seek(int64(pos_ante_ebr), 0)

				ebr_ext_ante := &ebr_part
				var binario_ebr_an bytes.Buffer
				binary.Write(&binario_ebr_an, binary.BigEndian, ebr_ext_ante)
				writeBytesDisk(file, binario_ebr_an.Bytes())
				/*******din actuluza atruas*******/

				return /// retorno porque ya encontrño espacio, ya no necsita recorrer el resto
			}

		}

		///////fmt.Println("tempo_sig", tempo_sig)
		/*recorriendo como una lista enlazada*/
		ebr_anterior := ebr_part
		fmt.Printf("						--anterior %s\n", ebr_anterior.Part_name)
		tempo_sig := ebr_part.Part_next
		for tempo_sig != -1 {

			///////fmt.Println("entra en tempo_sig", tempo_sig)

			file.Seek(int64(tempo_sig), 0)
			/*leyendo siguiente ebr*/
			/*******************ini read sig**************************/
			ebr_tem := EBR{}
			size_tem := int(binary.Size(ebr_tem))
			data_tem := readBytesDisk(file, size_tem)
			buffer_tem := bytes.NewBuffer(data_tem)
			err = binary.Read(buffer_tem, binary.BigEndian, &ebr_tem)
			if err != nil {
				log.Fatal("binary. se encontro error al leer EBR de la Extendida", err)
				break
			}

			fmt.Printf("--n %s\n", ebr_tem.Part_name)
			fmt.Println("ebr", ebr_tem.Part_start-uint64(size_tem), "-", ebr_tem.Part_start)
			fmt.Println(" ", ebr_tem.Part_start, "-", ebr_tem.Part_start+ebr_tem.Part_size, "|", ebr_tem.Part_size)
			/*libre en logicas*/
			var libre_log uint64
			if ebr_tem.Part_next == -1 {
				libre_log = fin_extendida - (ebr_tem.Part_start + ebr_tem.Part_size)
			} else {
				libre_log = uint64(ebr_tem.Part_next) - (ebr_tem.Part_start + ebr_tem.Part_size)
			}
			fmt.Println("                    libre: ", libre_log)
			//fmt.Println("               ebr_tem.Part_next: ", ebr_tem.Part_next)

			ingresar := size_part + uint64(size_ebr)
			if ingresar <= libre_log && ebr_tem.Part_next != -1 {
				fmt.Println("11 se puede ingresar en espacio libre")
				/*insertando particion en medio*/
				iniciar_par_bmr := ebr_tem.Part_start + ebr_tem.Part_size
				start := iniciar_par_bmr + uint64(size_ebr)
				////NewPartEBR(file *os.File, iniciar_par_bmr uint64, status byte, fit byte, start uint64, size_part uint64, siguiente int64, name_byte [16]byte )
				fmt.Println("mitad iniciar_par_bmr", iniciar_par_bmr)
				fmt.Println("mitad start", start)

				NewPartEBR(file, iniciar_par_bmr, 'S', fit, start, size_part, ebr_tem.Part_next, name_byte)

				/**********actuliza valor anteior************/
				//fmt.Printf("ebr_tem.Part_name:%s\n", ebr_tem.Part_name)
				ebr_tem.Part_next = int64(iniciar_par_bmr)
				pos_ante_ebr := ebr_tem.Part_start - uint64(size_ebr)
				//fmt.Println("pos_ante_ebr Primer EBR:", pos_ante_ebr)

				file.Seek(int64(pos_ante_ebr), 0)

				ebr_ext_ante := &ebr_tem
				var binario_ebr_an bytes.Buffer
				binary.Write(&binario_ebr_an, binary.BigEndian, ebr_ext_ante)
				writeBytesDisk(file, binario_ebr_an.Bytes())
				/*******din actuluza atruas*******/

				return /// retorno porque ya encontrño espacio, ya no necsita recorrer el resto
			}

			//tempo = tempo.siguiente
			tempo_sig = ebr_tem.Part_next
			///asigno mbr anterior
			ebr_anterior = ebr_tem
			fmt.Printf("						--anterior %s\n", ebr_anterior.Part_name)

			/*******************fin read sig**************************/
		}

		/*Verifico si hay espacio disponible*/
		fmt.Println("stempo_sig", tempo_sig)

		iniciar_par_bmr := ebr_anterior.Part_start + ebr_anterior.Part_size
		start := iniciar_par_bmr + uint64(size_ebr)
		fmt.Println("\nstart ultimo ebr", iniciar_par_bmr, "(+42) start ultimo", start)

		fmt.Println("size_part ultimo", size_part)

		libre_log = fin_extendida - (ebr_anterior.Part_start + ebr_anterior.Part_size)
		fmt.Println("xxx libre_log", libre_log)
		ingresar := size_part + uint64(size_ebr)
		if ingresar > libre_log {
			fmt.Println("No se encontró sufienete espacio para la particion, Espacio Libre:", libre_log, "bytes")
			fmt.Println("Espacio de Particion + EBR:", ingresar, "bytes")
			return
		}
		fmt.Println("*******222 se puede ingresar en espacio libre")
		//return

		////NewPartEBR(file *os.File, iniciar_par_bmr uint64, status byte, fit byte, start uint64, size_part uint64, siguiente int64, name_byte [16]byte )
		NewPartEBR(file, iniciar_par_bmr, 'S', fit, start, size_part, -1, name_byte)

		/**********actuliza valor anteior************/
		/////fmt.Println("Primer EBR:", m.Mbr_Particiones[num_part_extend].Part_start)
		ebr_anterior.Part_next = int64(iniciar_par_bmr)
		pos_ante_ebr := ebr_anterior.Part_start - uint64(size_ebr)
		fmt.Println("pos_ante_ebr Primer EBR:", pos_ante_ebr)
		file.Seek(int64(pos_ante_ebr), 0)

		ebr_ext_ante := &ebr_anterior
		var binario_ebr_an bytes.Buffer
		binary.Write(&binario_ebr_an, binary.BigEndian, ebr_ext_ante)
		writeBytesDisk(file, binario_ebr_an.Bytes())
		/*******din actuluza atruas*******/
		//////////////////}

	}
	/***fin para logico***/
}

func NewPartEBR(file *os.File, iniciar_particion uint64, status byte, fit byte, start uint64, size_part uint64, siguiente int64, name_byte [16]byte) {

	part_ex := EBR{}

	part_ex.Part_status = status
	part_ex.Part_fit = fit
	part_ex.Part_start = start
	part_ex.Part_size = size_part
	part_ex.Part_next = siguiente
	part_ex.Part_name = name_byte

	/*posicionandome en el primer ebr de la extendida*/
	//file.Seek(0, 0)
	file.Seek(int64(iniciar_particion), 0)

	/*escribiendo ebr inicial en disco CON DATOS (CABECERA), particion externa*/
	ebr_ext_n := &part_ex
	var binario_ebr_n bytes.Buffer
	binary.Write(&binario_ebr_n, binary.BigEndian, ebr_ext_n)

	msg_exito := "Se creo partición " + string(name_byte[:]) + " de " + strconv.FormatUint(size_part, 10) + " bytes"
	writeBytesPart(file, binario_ebr_n.Bytes(), msg_exito)

}

//file, pos_ebr_ant, ebr_anterior,  deletp
func ActualizarEBR_par(file *os.File, iniciar_ebr uint64, ebr_anterior EBR, deletp string, name_byte [16]byte) {

	/*posicionandome en el primer ebr de la extendida*/
	file.Seek(int64(iniciar_ebr), 0)

	/*escribiendo ebr inicial en disco CON DATOS (CABECERA), particion externa*/
	ebr_ext_n := &ebr_anterior
	var binario_ebr_n bytes.Buffer
	binary.Write(&binario_ebr_n, binary.BigEndian, ebr_ext_n)

	msg_exito := "Se Eliminó Partición " + string(name_byte[:]) + " con Exito"
	//str2 := string(name_byte[:])
	writeBytesPart(file, binario_ebr_n.Bytes(), msg_exito)

}

/*modificando tamaño de la partiicion*/
func Mount_Par(path string, name string) {

	exist_d := existeDisk(path)
	//fmt.Println("existe ar exist_d:",exist_d )
	if exist_d == false {
		fmt.Println("No existe Ruta del Disco:", path)
		return
	}

	Montando(path, name)
}

func Letra_Montada(path string) (byte, bool) {

	var let_ultimo byte
	if len(let_mon) == 0 {
		let_ultimo = 96
	}

	for i := 0; i < len(let_mon); i++ {
		let_ultimo = let_mon[i].L_letra
		if let_mon[i].L_path == path {
			return let_mon[i].L_letra, true
		}
	}
	return let_ultimo, false

}

func DesMontar(arr_desmont []string) {

	for i := 0; i < len(arr_desmont); i++ {
		//fmt.Print("id->", monta[i].M_id, " -path->", monta[i].M_path)
		//fmt.Printf(" -name->%s\n", monta[i].M_name )
		encontrado := Set_Desmont(arr_desmont[i])
		if encontrado == false {
			fmt.Println("No se encontró particion montada:", arr_desmont[i])
		} else {

		}

	}
}

func Set_Desmont(idmont string) bool {

	for i := 0; i < len(monta); i++ {
		if monta[i].M_id == idmont {

			fmt.Println("Particion desmontada:", idmont)
			monta = append(monta[:i], monta[i+1:]...)

			return true
		}
	}
	return false
}

func Mount_list() {

	//id->96a1 -path->"/home/Disco1.dsk" -name->"Part1"
	for i := 0; i < len(monta); i++ {
		fmt.Print("id=", monta[i].M_id, " -path=", monta[i].M_path)
		fmt.Printf(" -name=%s\n", monta[i].M_name)
	}
}

func Montando(path string, name string) {
	/*no se puede montar una particion extendida*/

	/*abrimos el archivo*/
	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	m := MBR{}
	var size int = int(binary.Size(m))

	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary. se encontro error al leer archivo binario", err)
	}

	/*convertiendo en bytes el nombre*/
	var name_byte [16]byte
	copy(name_byte[:], name)

	yaEstaMontado := YaMontado(path, name_byte)
	if yaEstaMontado == true {
		fmt.Println("Ya esta montada esa particion")
		return
	}

	hayExtendida := false
	nameinPart := false
	nameinLogicas := false
	var num_part_extend int
	for i := 0; i < len(m.Mbr_Particiones); i++ {
		fmt.Printf("%s ", m.Mbr_Particiones[i].Part_name)
		fmt.Println("part", i, string(m.Mbr_Particiones[i].Part_status), string(m.Mbr_Particiones[i].Part_type))

		if m.Mbr_Particiones[i].Part_status == 'S' {

			/*verifico NOMBRE PRIMARIAS O EXTERNA*/
			if m.Mbr_Particiones[i].Part_name == name_byte {
				nameinPart = true

				/*M_path string
				M_name [16]byte
				M_id string
				M_tipo byte
				M_P_Parti Particion
				M_P_Log EBR*/

				let, exis := Letra_Montada(path)
				if exis == true {
				} else {
					let = let + 1
					new_let := Letras{}
					new_let.L_path = path
					new_let.L_letra = let
					let_mon = append(let_mon, new_let)

					let_mon = OrdenmMLetras(let_mon)
				}

				idv := GetIdMon(path)
				name_mon := "96" + string(let) + strconv.FormatInt(idv, 10)

				/*fmt.Println("  idv", idv)
				fmt.Println("   let", string(let), let)
				fmt.Println("   exis", exis)
				fmt.Println("***   name_mon", name_mon)*/

				//monta []Montaciones
				mon_p := Montaciones{}
				mon_p.M_path = path
				mon_p.M_name = name_byte
				mon_p.M_tipo = m.Mbr_Particiones[i].Part_type
				mon_p.M_id = name_mon
				mon_p.M_num = idv
				mon_p.M_letra = let

				if m.Mbr_Particiones[i].Part_type == 'e' || m.Mbr_Particiones[i].Part_type == 'p' {

					mon_p.M_P_Parti = m.Mbr_Particiones[i]
				}
				monta = append(monta, mon_p)
				fmt.Println("Se montó partición id: ", name_mon)
				return
			}

			/*si hay una extendida*/
			if m.Mbr_Particiones[i].Part_type == 'e' {
				hayExtendida = true
				num_part_extend = i
			}
		}

	}
	/*verifico NOMBRE en normales*/
	if hayExtendida == false && nameinPart == false {
		fmt.Printf("No existe una particion con el  nombre de %s\n", name_byte)
		return
	}

	////fmt.Println("num_part_extend", num_part_extend)
	//fmt.Println("mon_p", monta)

	/***********inicio logico**************************/
	/*posicionandome en el primer ebr de la extendida*/
	file.Seek(0, 0)
	file.Seek(int64(m.Mbr_Particiones[num_part_extend].Part_start), 0)

	ebr_part := EBR{}
	size_ebr := int(binary.Size(ebr_part))

	data = readBytesDisk(file, size_ebr)
	buffer = bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable ebr_part
	err = binary.Read(buffer, binary.BigEndian, &ebr_part)
	if err != nil {
		log.Fatal("binary. se encontro error al leer EBR de la Extendida", err)
	}

	///*************************************************///
	fmt.Println("*************************************")
	/*fmt.Printf("--n %s\n", ebr_part.Part_name)
	fmt.Println("ebr", ebr_part.Part_start - uint64(size_ebr), "-", ebr_part.Part_start)
	fmt.Println(" ", ebr_part.Part_start, "-", ebr_part.Part_start +  ebr_part.Part_size, "|",  ebr_part.Part_size)
	*/
	if ebr_part.Part_name == name_byte {
		////fmt.Printf("--ENCONTRADO EN EL PRIMERO %s\n", ebr_part.Part_name)
		nameinLogicas = true

		if ebr_part.Part_status == 'S' {

			let, exis := Letra_Montada(path)
			if exis == true {
			} else {
				let = let + 1
				new_let := Letras{}
				new_let.L_path = path
				new_let.L_letra = let
				let_mon = append(let_mon, new_let)

				let_mon = OrdenmMLetras(let_mon)
			}

			idv := GetIdMon(path)
			name_mon := "96" + string(let) + strconv.FormatInt(idv, 10)

			/*fmt.Println("  idv", idv)
			fmt.Println("   let", string(let), let)
			fmt.Println("   exis", exis)
			fmt.Println("***   name_mon", name_mon)*/

			//monta []Montaciones
			mon_p := Montaciones{}
			mon_p.M_path = path
			mon_p.M_name = name_byte
			mon_p.M_tipo = 'l'
			mon_p.M_id = name_mon
			mon_p.M_num = idv
			mon_p.M_letra = let
			mon_p.M_P_Log = ebr_part

			monta = append(monta, mon_p)
			fmt.Println("Se montó partición id: ", name_mon)
			return
		}
	}

	/*recorriendo como una lista enlazada*/
	/////fmt.Printf("						--anterior %s\n", ebr_anterior.Part_name)
	tempo_sig := ebr_part.Part_next
	for tempo_sig != -1 {

		file.Seek(int64(tempo_sig), 0)
		/*leyendo siguiente ebr*/
		/*******************ini read sig**************************/
		ebr_tem := EBR{}
		size_tem := int(binary.Size(ebr_tem))
		data_tem := readBytesDisk(file, size_tem)
		buffer_tem := bytes.NewBuffer(data_tem)
		err = binary.Read(buffer_tem, binary.BigEndian, &ebr_tem)
		if err != nil {
			log.Fatal("binary. se encontro error al leer EBR de la Extendida", err)
			break
		}

		/*fmt.Printf("--n %s\n", ebr_tem.Part_name)
		fmt.Println("ebr", ebr_tem.Part_start - uint64(size_tem), "-", ebr_tem.Part_start)
		fmt.Println(" ", ebr_tem.Part_start, "-", ebr_tem.Part_start + ebr_tem.Part_size, "|",  ebr_tem.Part_size)
		*/
		if ebr_tem.Part_name == name_byte {

			nameinLogicas = true
			///fmt.Printf("------PARTICION ENCOTRADA %s\n", ebr_tem.Part_name)

			let, exis := Letra_Montada(path)
			if exis == true {
			} else {
				let = let + 1
				new_let := Letras{}
				new_let.L_path = path
				new_let.L_letra = let
				let_mon = append(let_mon, new_let)

				let_mon = OrdenmMLetras(let_mon)
			}

			idv := GetIdMon(path)
			name_mon := "96" + string(let) + strconv.FormatInt(idv, 10)

			/*fmt.Println("  idv", idv)
			fmt.Println("   let", string(let), let)
			fmt.Println("   exis", exis)
			fmt.Println("***   name_mon", name_mon)*/

			//monta []Montaciones
			mon_p := Montaciones{}
			mon_p.M_path = path
			mon_p.M_name = name_byte
			mon_p.M_tipo = 'l'
			mon_p.M_id = name_mon
			mon_p.M_num = idv
			mon_p.M_letra = let
			mon_p.M_P_Log = ebr_tem

			monta = append(monta, mon_p)
			fmt.Println("Se montó partición id: ", name_mon)
			return

		}

		//tempo = tempo.siguiente
		tempo_sig = ebr_tem.Part_next

		/*******************fin read sig**************************/
	}
	///**********************///fin de logicas

	//fmt.Println("nameinLogicas", nameinLogicas)
	/*verifico NOMBRE en LOGICAS*/
	if nameinLogicas == false {
		fmt.Printf("No existe una particion PRIMARIA/EXTEND/LOGICA con el  nombre de %s\n", name_byte)
		return
	}

}

/*modificando tamaño de la partiicion*/
func AddPartD(path string, name string, add int64, unit byte) {

	exist_d := existeDisk(path)
	//fmt.Println("existe ar exist_d:",exist_d )
	if exist_d == false {
		fmt.Println("No existe Ruta del Disco:", path)
		return
	}

	/*calcualando tamaño de disco*/
	var add_part int64
	if unit == 'b' {
		add_part = add
	} else if unit == 'k' {
		add_part = add * 1024
	} else if unit == 'm' {
		add_part = add * 1024 * 1024

	}

	if add_part > 0 {
		fmt.Println("agregar")
		Add_part(path, name, add_part, "+")
	} else if add_part < 0 {
		fmt.Println("quitar")
		add_part = add_part * -1
		Add_part(path, name, add_part, "-")
	}

}

func Add_part(path string, name string, add_part int64, signo string) {

	fmt.Println("add_part", add_part, "signo", signo)
	size_part := uint64(add_part)
	/*abrimos el archivo*/
	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	m := MBR{}
	var size int = int(binary.Size(m))

	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary. se encontro error al leer archivo binario", err)
	}

	fmt.Println("-----------tam_disk", m.Mbr_tamanio)

	/*convertiendo en bytes el nombre*/
	var name_byte [16]byte
	copy(name_byte[:], name)

	var orden_particion []Particion

	hayExtendida := false
	nameinPart := false
	nameinLogicas := false
	var num_part_extend int
	for i := 0; i < len(m.Mbr_Particiones); i++ {
		fmt.Printf("%s ", m.Mbr_Particiones[i].Part_name)
		fmt.Println("part", i, string(m.Mbr_Particiones[i].Part_status), string(m.Mbr_Particiones[i].Part_type))

		if m.Mbr_Particiones[i].Part_status == 'S' {
			lib := Particion{}
			lib.Part_status = m.Mbr_Particiones[i].Part_status
			lib.Part_type = m.Mbr_Particiones[i].Part_type
			lib.Part_fit = m.Mbr_Particiones[i].Part_fit
			lib.Part_start = m.Mbr_Particiones[i].Part_start
			lib.Part_size = m.Mbr_Particiones[i].Part_size
			lib.Part_name = m.Mbr_Particiones[i].Part_name
			orden_particion = append(orden_particion, lib)
		}

		/*ver rango de ocupado*/
		if m.Mbr_Particiones[i].Part_status == 'S' {

			/*verifico NOMBRE PRIMARIAS O EXTERNA*/
			if m.Mbr_Particiones[i].Part_name == name_byte {
				nameinPart = true
			}

		}
		/*si hay una extendida*/
		if m.Mbr_Particiones[i].Part_type == 'e' {
			hayExtendida = true
			num_part_extend = i
		}
	}
	////
	/*verifico NOMBRE en normales*/
	if hayExtendida == false && nameinPart == false {
		fmt.Printf("No existe una particion con el  nombre de %s\n", name_byte)
		return
	}

	orden_particion = OrdenmMParticion(orden_particion)
	for i := 0; i < len(orden_particion); i++ {
		fmt.Printf("000 %s ", orden_particion[i].Part_name)
		fmt.Println(i, string(orden_particion[i].Part_status), string(orden_particion[i].Part_type), orden_particion[i].Part_start, "-", orden_particion[i].Part_start+orden_particion[i].Part_size, "|", orden_particion[i].Part_size)

	}
	/****VERIFICANOD SI EXISTE PARTICION NORMAL*****/
	//return
	for i := 0; i < len(m.Mbr_Particiones); i++ {
		fmt.Printf("%s ", m.Mbr_Particiones[i].Part_name)
		fmt.Println(i, string(m.Mbr_Particiones[i].Part_status), string(m.Mbr_Particiones[i].Part_type), m.Mbr_Particiones[i].Part_start, "-", m.Mbr_Particiones[i].Part_start+m.Mbr_Particiones[i].Part_size, "|", m.Mbr_Particiones[i].Part_size)

		if m.Mbr_Particiones[i].Part_name == name_byte {
			fmt.Println("--**encontrado*")
			ingresar := size_part

			start := 0
			var new_size uint64
			var msg_exito string

			var ini_free uint64
			var fin_free uint64
			var size_free uint64

			if signo == "+" {

				/*si hay espacio libre*/
				for pa := 0; pa < len(orden_particion); pa++ {

					if orden_particion[pa].Part_name == name_byte {
						ini_free = orden_particion[pa].Part_start + orden_particion[pa].Part_size

						if (pa + 1) < len(orden_particion) {
							fin_free = orden_particion[pa+1].Part_start
						} else {
							fin_free = m.Mbr_tamanio
						}

						size_free = fin_free - ini_free
						fmt.Println(ini_free, "-", fin_free, "/", size_free)
						break
					}
				}

				if ingresar <= size_free {

					fmt.Println("11 se puede agregar en espacio libre")
					new_size = m.Mbr_Particiones[i].Part_size + ingresar
					fmt.Println("mitad new_size", new_size)
					fmt.Println("mitad iniciar_bmr", start)
					m.Mbr_Particiones[i].Part_size = new_size

					msg_exito = "Se Agrego a la Particion " + string(name_byte[:]) + " con Exito"

				} else {

					fmt.Println("No hay suficiente espacio, Espacio Libre:", size_free, "bytes")
					fmt.Println("Espacio a Agregar:", ingresar, "bytes")
					return
				}
				////////////////
			} else if signo == "-" {

				/*si sobra espacio para quitar*/
				if ingresar < m.Mbr_Particiones[i].Part_size {

					fmt.Println("p2 se puede quitar en espacio en particion")
					new_size = m.Mbr_Particiones[i].Part_size - ingresar
					fmt.Println("mitad new_size", new_size)
					fmt.Println("mitad iniciar_bmr", start)
					m.Mbr_Particiones[i].Part_size = new_size

					msg_exito = "Se Quitó a la Particion " + string(name_byte[:]) + " con Exito"

				} else {

					fmt.Println("No hay suficiente espacio para quitar Espacio:", m.Mbr_Particiones[i].Part_size, "bytes")
					fmt.Println("Espacio a quitar:", ingresar, "bytes")
					return
				}

			}

			/**********actuliza valor anteior************/
			file.Seek(int64(start), 0)
			ebr_ext_ante := &m
			var binario_ebr_an bytes.Buffer
			binary.Write(&binario_ebr_an, binary.BigEndian, ebr_ext_ante)
			writeBytesPart(file, binario_ebr_an.Bytes(), msg_exito)
			/*******din actuluza atruas*******/
			return

			//return
		}
	}
	/***********FIN**************************/

	/***********inicio logico**************************/

	fin_extendida := m.Mbr_Particiones[num_part_extend].Part_start + m.Mbr_Particiones[num_part_extend].Part_size
	fmt.Println("fin_extendida:", fin_extendida)
	/*posicionandome en el primer ebr de la extendida*/
	file.Seek(0, 0)
	file.Seek(int64(m.Mbr_Particiones[num_part_extend].Part_start), 0)

	ebr_part := EBR{}
	size_ebr := int(binary.Size(ebr_part))

	data = readBytesDisk(file, size_ebr)
	buffer = bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable ebr_part
	err = binary.Read(buffer, binary.BigEndian, &ebr_part)
	if err != nil {
		log.Fatal("binary. se encontro error al leer EBR de la Extendida", err)
	}
	fmt.Println(ebr_part)
	fmt.Printf("i-------nsert name %s\n", name_byte)

	fmt.Println("0ebr_part.Part_status", string(ebr_part.Part_status))
	fmt.Println("ebr_part.Part_next", ebr_part.Part_next)
	fmt.Println("ebr_part.Part_start", ebr_part.Part_start)
	fmt.Println("ebr_part.Part_size", ebr_part.Part_size)
	fmt.Printf("ebr_part.Part_name %s\n", ebr_part.Part_name)

	///*************************************************///
	fmt.Println("*************************************")
	fmt.Printf("--n %s\n", ebr_part.Part_name)
	fmt.Println("ebr", ebr_part.Part_start-uint64(size_ebr), "-", ebr_part.Part_start)
	fmt.Println(" ", ebr_part.Part_start, "-", ebr_part.Part_start+ebr_part.Part_size, "|", ebr_part.Part_size)
	/*libre en logicas*/
	var libre_log uint64
	if ebr_part.Part_next == -1 {
		if ebr_part.Part_status == 'N' {
			libre_log = fin_extendida - (m.Mbr_Particiones[num_part_extend].Part_start)
		} else {
			libre_log = fin_extendida - (ebr_part.Part_start + ebr_part.Part_size)
		}
	} else { ////////////////

		if ebr_part.Part_status == 'N' {
			libre_log = uint64(ebr_part.Part_next) - (m.Mbr_Particiones[num_part_extend].Part_start)
		} else { /////////////
			libre_log = uint64(ebr_part.Part_next) - (ebr_part.Part_start + ebr_part.Part_size)
		}
	}
	fmt.Println("                    -libre: ", libre_log)

	if ebr_part.Part_name == name_byte {
		fmt.Printf("--ENCONTRADO EN EL PRIMERO %s\n", ebr_part.Part_name)
		nameinLogicas = true

		ingresar := size_part

		if ebr_part.Part_status == 'S' /*&& ebr_part.Part_next != -1*/ {

			start := ebr_part.Part_start - uint64(size_ebr)
			var new_size uint64
			var msg_exito string
			if signo == "+" {
				/*si hay espacio libre*/
				if ingresar <= libre_log /*&& ebr_tem.Part_next != -1*/ {

					fmt.Println("001 se puede agregar en espacio libre")
					new_size = ebr_part.Part_size + ingresar
					fmt.Println("mitad new_size", new_size)
					fmt.Println("mitad iniciar_bmr", start)
					ebr_part.Part_size = new_size

					msg_exito = "Se Agrego a la Particion " + string(name_byte[:]) + " con Exito"

				} else {

					fmt.Println("No hay suficiente espacio, Espacio Libre:", libre_log, "bytes")
					fmt.Println("Espacio a Agregar:", ingresar, "bytes")
					return
				}
			} else if signo == "-" {

				/*si sobra espacio para quitar*/
				if ingresar < ebr_part.Part_size {

					fmt.Println("002 se puede quitar en espacio en particion")
					new_size = ebr_part.Part_size - ingresar
					fmt.Println("mitad new_size", new_size)
					fmt.Println("mitad iniciar_bmr", start)
					ebr_part.Part_size = new_size

					msg_exito = "Se Quitó a la Particion " + string(name_byte[:]) + " con Exito"

				} else {

					fmt.Println("No hay suficiente espacio para quitar Espacio:", ebr_part.Part_size, "bytes")
					fmt.Println("Espacio a quitar:", ingresar, "bytes")
					return
				}

			}

			/**********actuliza valor anteior************/
			file.Seek(int64(start), 0)
			ebr_ext_ante := &ebr_part
			var binario_ebr_an bytes.Buffer
			binary.Write(&binario_ebr_an, binary.BigEndian, ebr_ext_ante)

			writeBytesPart(file, binario_ebr_an.Bytes(), msg_exito)
			/*******din actuluza atruas*******/
			return

		}
	}

	///////fmt.Println("tempo_sig", tempo_sig)
	/*recorriendo como una lista enlazada*/
	ebr_anterior := ebr_part
	fmt.Printf("						--anterior %s\n", ebr_anterior.Part_name)
	tempo_sig := ebr_part.Part_next
	for tempo_sig != -1 {

		///////fmt.Println("entra en tempo_sig", tempo_sig)

		file.Seek(int64(tempo_sig), 0)
		/*leyendo siguiente ebr*/
		/*******************ini read sig**************************/
		ebr_tem := EBR{}
		size_tem := int(binary.Size(ebr_tem))
		data_tem := readBytesDisk(file, size_tem)
		buffer_tem := bytes.NewBuffer(data_tem)
		err = binary.Read(buffer_tem, binary.BigEndian, &ebr_tem)
		if err != nil {
			log.Fatal("binary. se encontro error al leer EBR de la Extendida", err)
			break
		}

		fmt.Printf("--n %s\n", ebr_tem.Part_name)
		fmt.Println("ebr", ebr_tem.Part_start-uint64(size_tem), "-", ebr_tem.Part_start)
		fmt.Println(" ", ebr_tem.Part_start, "-", ebr_tem.Part_start+ebr_tem.Part_size, "|", ebr_tem.Part_size)
		/*libre en logicas*/
		var libre_log uint64
		if ebr_tem.Part_next == -1 {
			libre_log = fin_extendida - (ebr_tem.Part_start + ebr_tem.Part_size)
		} else {
			libre_log = uint64(ebr_tem.Part_next) - (ebr_tem.Part_start + ebr_tem.Part_size)
		}
		fmt.Println("                    libre: ", libre_log)

		ingresar := size_part
		if ebr_tem.Part_name == name_byte {
			nameinLogicas = true
			fmt.Printf("------PARTICION ENCOTRADA %s\n", ebr_tem.Part_name)
			start := ebr_tem.Part_start - uint64(size_tem)
			var new_size uint64
			var msg_exito string
			if signo == "+" {
				/*si hay espacio libre*/
				if ingresar <= libre_log /*&& ebr_tem.Part_next != -1*/ {

					fmt.Println("11 se puede agregar en espacio libre")
					new_size = ebr_tem.Part_size + ingresar
					fmt.Println("mitad new_size", new_size)
					fmt.Println("mitad iniciar_bmr", start)
					ebr_tem.Part_size = new_size

					msg_exito = "Se Agrego a la Particion " + string(name_byte[:]) + " con Exito"

				} else {

					fmt.Println("No hay suficiente espacio, Espacio Libre:", libre_log, "bytes")
					fmt.Println("Espacio a Agregar:", ingresar, "bytes")
					return
				}
			} else if signo == "-" {

				/*si sobra espacio para quitar*/
				if ingresar < ebr_tem.Part_size {

					fmt.Println("12 se puede quitar en espacio en particion")
					new_size = ebr_tem.Part_size - ingresar
					fmt.Println("mitad new_size", new_size)
					fmt.Println("mitad iniciar_bmr", start)
					ebr_tem.Part_size = new_size

					msg_exito = "Se Quitó a la Particion " + string(name_byte[:]) + " con Exito"

				} else {

					fmt.Println("No hay suficiente espacio para quitar Espacio:", ebr_tem.Part_size, "bytes")
					fmt.Println("Espacio a quitar:", ingresar, "bytes")
					return
				}

			}

			/**********actuliza valor anteior************/
			file.Seek(int64(start), 0)
			ebr_ext_ante := &ebr_tem
			var binario_ebr_an bytes.Buffer
			binary.Write(&binario_ebr_an, binary.BigEndian, ebr_ext_ante)

			writeBytesPart(file, binario_ebr_an.Bytes(), msg_exito)
			/*******din actuluza atruas*******/
			return

		}

		//tempo = tempo.siguiente
		tempo_sig = ebr_tem.Part_next
		///asigno mbr anterior
		ebr_anterior = ebr_tem
		fmt.Printf("						--anterior %s\n", ebr_anterior.Part_name)

		/*******************fin read sig**************************/
	}
	///**********************///fin de logicas

	//fmt.Println("nameinLogicas", nameinLogicas)
	/*verifico NOMBRE en LOGICAS*/
	if nameinLogicas == false {
		fmt.Printf("No existe una particion PRIMARIA/EXTEND/LOGICA con el  nombre de %s\n", name_byte)
		return
	}

}

/*eliminando particion en disco*/
func DeletePartD(path string, name string, deletp string) {

	exist_d := existeDisk(path)
	//fmt.Println("existe ar exist_d:",exist_d )
	if exist_d == false {
		fmt.Println("No existe Ruta del Disco:", path)
		return
	}

	read_n_delete(path, name, deletp)
}

/*leyendo archivo, y delete particion*/
func read_n_delete(path string, name string, deletp string) {
	/*abrimos el archivo*/
	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	m := MBR{}
	var size int = int(binary.Size(m))
	//fmt.Println("binary :",size )

	data := readBytesDisk(file, size)
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary. se encontro error al leer archivo binario", err)
	}

	///fmt.Println(m)

	/*convertiendo en bytes el nombre*/
	var name_byte [16]byte
	copy(name_byte[:], name)

	hayExtendida := false
	nameinPart := false
	nameinLogicas := false
	var num_part_extend int
	for i := 0; i < len(m.Mbr_Particiones); i++ {
		fmt.Printf("%s ", m.Mbr_Particiones[i].Part_name)
		fmt.Println("part", i, string(m.Mbr_Particiones[i].Part_status), string(m.Mbr_Particiones[i].Part_type))

		/*ver rango de ocupado*/
		if m.Mbr_Particiones[i].Part_status == 'S' {

			/*verifico NOMBRE PRIMARIAS O EXTERNA*/
			if m.Mbr_Particiones[i].Part_name == name_byte {
				nameinPart = true
			}

		}
		/*si hay una extendida*/
		if m.Mbr_Particiones[i].Part_type == 'e' {
			hayExtendida = true
			num_part_extend = i
		}
	}
	////
	/*verifico NOMBRE en normales*/
	if hayExtendida == false && nameinPart == false {
		fmt.Printf("No existe una particion con el  nombre de %s\n", name_byte)
		return
	}

	/**********valida si desea eliminar S/N***************************/
	confirmacion := false
	for confirmacion == false {
		fmt.Printf("Esta Seguro de Partición %s\n", name_byte)
		fmt.Println("S/N")

		reader := bufio.NewReader(os.Stdin)
		lin_coman, _ := reader.ReadString('\n')
		comando := strings.ReplaceAll(lin_coman, "\r", "")
		comando = strings.ReplaceAll(comando, "\n", "")

		comando = strings.ToLower(comando)
		///fmt.Println("-",comando,"+")

		if comando == "s" {
			confirmacion = true
			//fmt.Println("Partición Eliminada correctamente")
		} else if comando == "n" {
			confirmacion = true
			return
		} else {
			confirmacion = false
		}
	}

	//fmt.Println("Partición Eliminada correctamente")
	//return
	/*************************************/
	/****VERIFICANOD SI EXISTE PARTICION NORMAL*****/
	for i := 0; i < len(m.Mbr_Particiones); i++ {
		fmt.Printf("%s ", m.Mbr_Particiones[i].Part_name)
		fmt.Println(i, string(m.Mbr_Particiones[i].Part_status), string(m.Mbr_Particiones[i].Part_type), m.Mbr_Particiones[i].Part_start, "-", m.Mbr_Particiones[i].Part_start+m.Mbr_Particiones[i].Part_size, "|", m.Mbr_Particiones[i].Part_size)

		if m.Mbr_Particiones[i].Part_name == name_byte {
			fmt.Println("*encontrado*")

			par_star := m.Mbr_Particiones[i].Part_start
			par_size := m.Mbr_Particiones[i].Part_size

			m.Mbr_Particiones[i].Part_status = 'N'
			m.Mbr_Particiones[i].Part_type = ' '
			m.Mbr_Particiones[i].Part_fit = ' '
			m.Mbr_Particiones[i].Part_start = 0
			m.Mbr_Particiones[i].Part_size = 0
			//m.Mbr_Particiones[i].Part_name = name_byte
			copy(m.Mbr_Particiones[i].Part_name[:], " ")

			file.Seek(0, 0)
			disk_w22 := &m

			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, disk_w22)

			msg_exito := "Se Eliminó Partición " + string(name_byte[:]) + " con Exito"
			writeBytesPart(file, binario2.Bytes(), msg_exito)

			if deletp == "full" {
				FormatFull(file, par_star, par_size)
			}

			return
		}
	}
	/***********FIN**************************/

	/***********inicio logico**************************/

	file.Seek(0, 0)
	file.Seek(int64(m.Mbr_Particiones[num_part_extend].Part_start), 0)

	ebr_part := EBR{}
	size_ebr := int(binary.Size(ebr_part))

	data = readBytesDisk(file, size_ebr)
	buffer = bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable ebr_part
	err = binary.Read(buffer, binary.BigEndian, &ebr_part)
	if err != nil {
		log.Fatal("binary. se encontro error al leer EBR de la Extendida", err)
	}
	fmt.Println(ebr_part)

	fmt.Println("*************************************")
	fmt.Printf("--n %s\n", ebr_part.Part_name)
	fmt.Println("ebr", ebr_part.Part_start-uint64(size_ebr), "-", ebr_part.Part_start)
	fmt.Println(" ", ebr_part.Part_start, "-", ebr_part.Part_start+ebr_part.Part_size, "|", ebr_part.Part_size)
	fmt.Println(" Part_next", ebr_part.Part_next)

	if ebr_part.Part_name == name_byte {
		fmt.Printf("--ENCONTRADO EN EL PRIMERO %s\n", ebr_part.Part_name)
		nameinLogicas = true

		/*elimina particion encotrada*/
		par_star := ebr_part.Part_start
		par_size := ebr_part.Part_size
		ebr_part.Part_status = 'N'
		/*actuliza en mbr ante*/
		pos_ebr_ant := ebr_part.Part_start - uint64(size_ebr)
		/////////////ebr_part.Part_start = 0
		ebr_part.Part_size = 0
		copy(ebr_part.Part_name[:], "                ")
		ActualizarEBR_par(file, pos_ebr_ant, ebr_part, deletp, ebr_part.Part_name)
		if deletp == "full" {
			FormatFull(file, par_star, par_size)
		}

	}
	/*recorriendo como una lista enlazada*/
	ebr_anterior := ebr_part
	fmt.Printf("						--anterior %s\n", ebr_anterior.Part_name)
	tempo_sig := ebr_part.Part_next
	fmt.Println("--tempo_sig", tempo_sig)
	for tempo_sig != -1 {

		file.Seek(int64(tempo_sig), 0)
		/*leyendo siguiente ebr*/
		/*******************ini read sig**************************/
		ebr_tem := EBR{}
		size_tem := int(binary.Size(ebr_tem))
		data_tem := readBytesDisk(file, size_tem)
		buffer_tem := bytes.NewBuffer(data_tem)
		err = binary.Read(buffer_tem, binary.BigEndian, &ebr_tem)
		if err != nil {
			log.Fatal("binary. se encontro error al leer EBR de la Extendida", err)
			break
		}

		fmt.Printf("--n %s\n", ebr_tem.Part_name)
		fmt.Println("ebr", ebr_tem.Part_start-uint64(size_tem), "-", ebr_tem.Part_start)
		fmt.Println(" ", ebr_tem.Part_start, "-", ebr_tem.Part_start+ebr_tem.Part_size, "|", ebr_tem.Part_size)
		fmt.Println(" Part_next", ebr_tem.Part_next)
		if ebr_tem.Part_name == name_byte {
			nameinLogicas = true
			/*elimina particion encotrada*/
			fmt.Printf("---PARTICION ANTERIOR %s\n", ebr_anterior.Part_name)
			fmt.Println("+++++++++++ANTERIOR start", ebr_anterior.Part_start)
			fmt.Println("+++++++++++ANTERIOR siguiente", ebr_anterior.Part_next)
			fmt.Printf("---PARTICION ENCOTRADA %s\n", ebr_tem.Part_name)
			fmt.Println("+++++++++++PARTICION siguiente", ebr_tem.Part_next)
			fmt.Println("+++++++++++PARTICION start", ebr_tem.Part_start)

			//return
			/*actulizo apuntadores de anterior*/
			ebr_anterior.Part_next = ebr_tem.Part_next
			/*actuliza en mbr ante*/
			pos_ebr_ant := ebr_anterior.Part_start - uint64(size_tem)
			ActualizarEBR_par(file, pos_ebr_ant, ebr_anterior, deletp, ebr_tem.Part_name)
			if deletp == "full" {
				FormatFull(file, ebr_tem.Part_start, ebr_tem.Part_size)
			}
			tempo_sig = -1
		}
		//tempo = tempo.siguiente
		tempo_sig = ebr_tem.Part_next
		///asigno mbr anterior
		ebr_anterior = ebr_tem
		fmt.Printf("						--anterior %s\n", ebr_anterior.Part_name)

		/*******************fin read sig**************************/
	}
	///**********************///fin de logicas

	//fmt.Println("nameinLogicas", nameinLogicas)
	/*verifico NOMBRE en LOGICAS*/
	if nameinLogicas == false {
		fmt.Printf("No existe una particion PRIMARIA/EXTEND/LOGICA con el  nombre de %s\n", name_byte)
		return
	}
}

func FormatFull(file *os.File, start uint64, tam uint64) {
	var cero int8 = 0
	cer := &cero

	fmt.Println("tstart", start)
	file.Seek(int64(start), 0)

	limit := int(tam)
	fmt.Println("tam uint64", tam)
	fmt.Println("tam int", limit)

	fin := start + uint64(limit)
	fmt.Println("---", fin)

	for i := 0; i < limit; i++ {
		/*escribiendo */
		var binario bytes.Buffer
		binary.Write(&binario, binary.BigEndian, cer)
		////msg_exito := "Se Eliminó Partición " + string(name_byte[:])  + " con Exito"
		//writeBytesPart(file, binario.Bytes(), "")
		writeBytesDisk(file, binario.Bytes())
		//proces := (i/int(tam)*100)
		////fmt.Println(i, int(tam))
	}

	fmt.Println("Formateo full terminado")

}

type EspacioLibre struct {
	free_start uint64 // 8 bytes
	free_end   uint64 // 8 bytes
	free_size  uint64 // 8 bytes

	P_name [16]byte
	P_type byte
}

func exist_partlogicas(file *os.File, name string, m MBR, num_part_extend int) bool {

	/*convertiendo en bytes el nombre*/
	var name_byte [16]byte
	copy(name_byte[:], name)

	nameinLogicas := false

	file.Seek(int64(m.Mbr_Particiones[num_part_extend].Part_start), 0)

	ebr_part := EBR{}
	size_ebr := int(binary.Size(ebr_part))

	data := readBytesDisk(file, size_ebr)
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable ebr_part
	err := binary.Read(buffer, binary.BigEndian, &ebr_part)
	if err != nil {
		log.Fatal("binary. se encontro error al leer EBR de la Extendida", err)
	}
	fmt.Println(ebr_part)

	///fmt.Println("*************************************")
	///fmt.Printf("--n %s\n", ebr_part.Part_name)
	///fmt.Println("ebr", ebr_part.Part_start - uint64(size_ebr), "-", ebr_part.Part_start)
	//fmt.Println(" ", ebr_part.Part_start, "-", ebr_part.Part_start +  ebr_part.Part_size, "|",  ebr_part.Part_size)

	if ebr_part.Part_name == name_byte {
		nameinLogicas = true
	}
	/*recorriendo como una lista enlazada*/
	/////ebr_anterior := ebr_part
	////fmt.Printf("						--anterior %s\n", ebr_anterior.Part_name)
	tempo_sig := ebr_part.Part_next
	for tempo_sig != -1 {

		file.Seek(int64(tempo_sig), 0)
		/*leyendo siguiente ebr*/
		/*******************ini read sig**************************/
		ebr_tem := EBR{}
		size_tem := int(binary.Size(ebr_tem))
		data_tem := readBytesDisk(file, size_tem)
		buffer_tem := bytes.NewBuffer(data_tem)
		err = binary.Read(buffer_tem, binary.BigEndian, &ebr_tem)
		if err != nil {
			log.Fatal("binary. se encontro error al leer EBR de la Extendida", err)
			break
		}

		////fmt.Printf("--n %s\n", ebr_tem.Part_name)
		////fmt.Println("ebr", ebr_tem.Part_start - uint64(size_tem), "-", ebr_tem.Part_start)
		///fmt.Println(" ", ebr_tem.Part_start, "-", ebr_tem.Part_start +  ebr_tem.Part_size, "|",  ebr_tem.Part_size)

		if ebr_tem.Part_name == name_byte {
			nameinLogicas = true
			/*elimina particion encotrada*/
		}
		//tempo = tempo.siguiente
		tempo_sig = ebr_tem.Part_next
		///asigno mbr anterior
		/////ebr_anterior = ebr_tem
		////fmt.Printf("						--anterior %s\n", ebr_anterior.Part_name)

		/*******************fin read sig**************************/
	}
	//////fin recorriendo

	return nameinLogicas
}

func readBytesDisk(file *os.File, number int) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func writeBytesDisk(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal("err ", err)
	}

}

func writeBytesPart(file *os.File, bytes []byte, msg_exito string) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal("err ", err)
	} else {
		fmt.Println(msg_exito)
	}

}

/*
func main() {
	//var t time.Time
	//t := time.Now()
	//fmt.Println(t)

	mbr_n := MBR{Mbr_tamanio: 1}
	//mbr_n.Mbr_fecha_creacion = time.Now()
	mbr_n.Mbr_disk_signature = 3
	fecha_string := time.Now().Format("2006-01-02 15:04:05")
	copy(mbr_n.Mbr_fecha_creacion_s[:], fecha_string)

	//fmt.Println("fecha_string:",  fecha_string)
	fmt.Printf("Mbr_fecha_creacion_s: %s\n" ,mbr_n.Mbr_fecha_creacion_s)


	//fmt.Println("time.Now()---:", unsafe.Sizeof(time.Now()))
	//fmt.Println("bina:", binary.Size(time.Now()))

	fmt.Println("*****bunary mbr_n.Mbr_fecha_creacion_s:", binary.Size(mbr_n.Mbr_fecha_creacion_s))

	par1 := Particion{}




	par1.Part_status = 'S'
	par1.Part_type = 'P'
	par1.Part_fit = 'B'
	par1.Part_start = 250
	par1.Part_size = 100
	copy(par1.Part_name[:], "Parti_1")

	////fmt.Println("tam par1 ---:", int(unsafe.Sizeof(par1)))
	//fmt.Println("tam par1 bina:", binary.Size(par1))
	//fmt.Println("par1.Part_status---:", unsafe.Sizeof(par1.Part_status))
	//fmt.Println("par1.Part_type---:", unsafe.Sizeof(par1.Part_type))
	//fmt.Println("par1.Part_fit---:", unsafe.Sizeof(par1.Part_fit))
	//fmt.Println("par1.Part_start---:", unsafe.Sizeof(par1.Part_start))
	//fmt.Println("par1.Part_size---:", unsafe.Sizeof(par1.Part_size))
	//fmt.Println("par1.Part_name---:", unsafe.Sizeof(par1.Part_name))

	mbr_n.Mbr_part1 = par1


	fmt.Println("mbr_n :", unsafe.Sizeof(mbr_n))
	fmt.Println("mbr_n bina:", binary.Size(mbr_n))

	fmt.Println(mbr_n)

	fmt.Println("Mbr_tamanio: ", mbr_n.Mbr_tamanio)
	fmt.Println("Mbr_fecha_creacion_s: ", mbr_n.Mbr_fecha_creacion_s)
	//fmt.Println("fecha format: ", mbr_n.Mbr_fecha_creacion.Format("2006-01-02 15:04:05"))

	fmt.Println("Mbr_part1: ", mbr_n.Mbr_part1)

	fmt.Printf("Part_name: %s", mbr_n.Mbr_part1.Part_name)

	fmt.Println("\nMbr_part2: ", mbr_n.Mbr_part2)


}*/
