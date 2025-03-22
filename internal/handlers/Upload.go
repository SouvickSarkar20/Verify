package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	algo "verification/internal/algo"

	"github.com/gin-gonic/gin"
)

func UploadHandler(c *gin.Context) {
	fmt.Println("Upload handler triggered")

	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed", "err": err.Error()})
		return
	}

	// Ensure storage directory exists
	storageDir := "./storage"
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		err := os.MkdirAll(storageDir, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create storage directory"})
			return
		}
	}

	savePath := fmt.Sprintf("%s/%s", storageDir, file.Filename)
	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		fmt.Printf("Error saving file: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File save failed", "details": err.Error()})
		return
	}

	//generate the rsa keys here
	privateKey, err := algo.GenerateRSAKeys(2048)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Key generation failed", "details": err.Error()})
		return
	}

	//next we will save the public key -> further in the blockchain
	publicKeyPath := filepath.Join(storageDir, file.Filename+"_public_key.pem")
	err = algo.SavePublicKey(&privateKey.PublicKey, file.Filename+"_public_key.pem")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Public key save failed", "details": err.Error()})
		return
	}

	//hash the file
	hash, err := algo.HashFile(savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File hashing failed", "details": err.Error()})
		return
	}

	//sign the hash
	signature, err := algo.SignHash(privateKey, hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Hash signing failed", "details": err.Error()})
		return
	}

	//save the signature
	signaturePath := filepath.Join(storageDir, file.Filename+"_signature.txt")
	err = algo.SaveSignature(signature, signaturePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Signature save failed", "details": err.Error()})
		return
	}

	//respond to client
	c.JSON(http.StatusOK, gin.H{
		"message":        "File uploaded and signed successfully",
		"file":           file.Filename,
		"public_key":     publicKeyPath,
		"signature_file": signaturePath,
	})
}

//done with phase-1
