package controller

import (
	"baseModule/models"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	ipfs "github.com/ipfs/go-ipfs-api"
)

const localPath = "./download"
const publicKey = ""

type AudioController struct {
}

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

func (ac *AudioController) GetTextByhash(c *gin.Context) {
	sh := ipfs.NewShell("localhost:5001")

	hash := c.Param("hash")

	// Agregar el texto a IPFS
	content, err := readFile(sh, hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Se obtuvo con exito", "texto": content})

}

func addFile(sh *ipfs.Shell, text string) (string, error) {
	return sh.Add(strings.NewReader(text))
}

func readFile(sh *ipfs.Shell, cid string) (*string, error) {
	reader, err := sh.Cat(fmt.Sprint("/ipfs/", cid))
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
