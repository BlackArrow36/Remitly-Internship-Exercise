package controller

import (
	"net/http"
	"path/filepath"
	"strings"
	"swift-api/importer"

	"github.com/gin-gonic/gin"
)

func ImportSWIFTData(c *gin.Context) {
	fileName := c.Query("file")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'file' query parameter",
		})
		return
	}
	if !strings.HasSuffix(fileName, ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Only .csv files are allowed",
		})
		return
	}

	filePath := filepath.Join("import", fileName)

	parsedData, err := importer.ParseCSV(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to parse CSV",
			"details": err.Error(),
		})
		return
	}

	err = importer.InsertSwiftCodes(parsedData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to insert SWIFT codes",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SWIFT codes imported successfully!"})
}
