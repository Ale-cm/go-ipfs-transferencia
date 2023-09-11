package controller

import (
	"api-audio-go/models"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"io"

	ipfs "github.com/ipfs/go-ipfs-api"
)

type AudioController struct {
}

const localPath = "./download"

const publicKey = ""

type Respon struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    interface{}
}

func (e *AudioController) handleSucces(c *gin.Context, data interface{}) {
	var returnData = Respon{
		Status:  "0000",
		Message: "Success",
		Data:    data,
	}
	c.JSON(http.StatusOK, returnData)
}

func (e *AudioController) handleError(c *gin.Context, message string) {
	var returnData = Respon{
		Status:  "501",
		Message: message,
	}
	c.JSON(http.StatusBadRequest, returnData)
}

func (ac *AudioController) HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World"})
}

func (ac *AudioController) UploadText(c *gin.Context) {
	// Obtiene el archivo de audio del formulario con el nombre "audio"
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(c.Request.Body)
	file, err := c.FormFile("text")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Crea un directorio de almacenamiento si no existe
	storagePath := "./uploads/text" // Cambia esto a la ubicación deseada
	if err := os.MkdirAll(storagePath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Guarda el archivo de audio en el directorio de almacenamiento
	filePath := storagePath + "/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// En este punto, el archivo de audio ha sido guardado con éxito
	// Puedes realizar otras operaciones aquí, como procesar el audio o registrar su ubicación

	c.JSON(http.StatusOK, gin.H{"message": "Archivo de audio guardado exitosamente", "filePath": filePath})
}

func (ac *AudioController) UploadTextB64(c *gin.Context) {
	var text = models.File{}
	err := c.Bind(&text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(text.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	decodedString := string(decodedBytes)

	sh := ipfs.NewShell("localhost:5001")

	// Agregar el texto a IPFS
	hash, err := addFile(sh, decodedString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf(hash)

	c.JSON(http.StatusOK, gin.H{"message": "Archivo de audio guardado exitosamente", "texto": decodedString, "hash": hash})

}

func addFile(sh *ipfs.Shell, text string) (string, error) {
	return sh.Add(strings.NewReader(text))
}

func readFile(sh *ipfs.Shell, cid string) (*string, error) {
	reader, err := sh.Cat(fmt.Sprint("/ipfs/%s", cid))
	if err != nil {
		return nil, fmt.Errorf("error reading the file %s", err.Error())
	}
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading bytes: %s", err.Error())
	}
	text := string(bytes)

	return &text, nil
}

func downloadFile(sh *ipfs.Shell, cid string) error {
	return sh.Get(cid, localPath)
}
