package main

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var db *gorm.DB

func init() {
	//open a db connection
	var err error
	db, err = gorm.Open("mysql", "root:@/misdb?charset=utf8&parseTime=True&loc=Local")


	if err != nil {
		panic("failed to connect database")
	}

	//Migrate the schema
	db.AutoMigrate(&materialModel{})
}

func main() {

	router := gin.Default()
	router.Use(cors.Default())

	v1 := router.Group("/api/v1/materials")
	{
		v1.POST("/", createMaterial)
		v1.GET("/", getAllMaterial)
		v1.GET("/:ID", getMaterialById)
		v1.PUT("/:ID", updateMaterial)
		v1.DELETE("/:ID", deleteMaterial)
	}
	router.Run()

}

type (
	// materialModel describes a materialModel type
	materialModel struct {
		gorm.Model
		Trademark     string `json:"Trademark"`
		IsBroken int    `json:"IsBroken"`
		Color string `json:"Color"`
		Date time.Time `json:"Date"`
		Description string `json:"Description"`
		InputBy string `json:"InputBy"`
		Name string `json:"Name"`
		Size int `json:"Size"`
		Type string `json:"Type"`
		Vendor string `json:"Vendor"`


	}

	// materialViewModel represents a formatted material
	materialViewModel struct {
		ID        uint   `json:"ID"`
		Trademark     string `json:"Trademark"`
		IsBroken bool   `json:"IsBroken"`
		Color string `json:"Color"`
		Date time.Time `json:"Date" time_format:"unix"`
		Description string `json:"Description"`
		InputBy string `json:"InputBy"`
		Name string `json:"Name"`
		Size int `json:"Size"`
		Type string `json:"Type"`
		Vendor string `json:"Vendor"`
	}
)

// createMaterial add a new material
func createMaterial(c *gin.Context) {

	isBroken, _ := strconv.Atoi(c.PostForm("IsBroken"))
	size, _ :=strconv.Atoi(c.PostForm("Size"))

	date, _ := time.Parse(time.RFC822Z, c.PostForm("Date"))
	material := materialModel{
		Trademark: c.PostForm("Trademark"),
		IsBroken: isBroken,
		Color:c.PostForm("Color"),
		Date: date,
		Description: c.PostForm("Description"),
		InputBy: c.PostForm("InputBy"),
		Name: c.PostForm("Name"),
		Type: c.PostForm("Type"),
		Vendor: c.PostForm("Vendor"),
		Size: size,

	}
	db.Save(&material)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Material item created successfully!", "resourceId": material.ID})
}

// getAllMaterial fetch all materials
func getAllMaterial(c *gin.Context) {

	var materials []materialModel
	var materialInfo []materialViewModel

	db.Find(&materials)

	if len(materials) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No material found!"})
		return
	}

	//transforms the material for building a good response
	for _, item := range materials {
		isBroken := false
		if item.IsBroken == 1 {
			isBroken = true
		} else {
			isBroken = false
		}
		materialInfo = append(materialInfo,
			materialViewModel{
				ID: item.ID,
				Trademark: item.Trademark,
				IsBroken: isBroken,
				Type: item.Type,
				Name: item.Name,
				Vendor: item.Vendor,
				InputBy: item.InputBy,
				Description: item.Description,
				Color: item.Color,
				Size: item.Size,
				Date: item.Date,
			})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": materialInfo})
}

// getMaterialById fetch a single material
func getMaterialById(c *gin.Context) {
	var material materialModel
	materialID := c.Param("ID")

	db.First(&material, materialID)

	if material.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No material found!"})
		return
	}

	isBroken := false
	if material.IsBroken == 1 {
		isBroken = true
	} else {
		isBroken = false
	}

	materialInfo := materialViewModel{
		ID: material.ID,
		Trademark: material.Trademark,
		IsBroken: isBroken,
		Date: material.Date,
		Size: material.Size,
		Color: material.Color,
		InputBy: material.InputBy,
		Description: material.Description,
		Name: material.Name,
		Vendor: material.Vendor,
		Type: material.Type,
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": materialInfo})
}



// updateMaterial update a material
func updateMaterial(c *gin.Context) {
	var material materialModel
	materialID := c.Param("ID")

	db.First(&material, materialID)

	if material.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No material found!"})
		return
	}

	db.Model(&material).Update(
		"Trademark", c.PostForm("Trademark"),
		"Vendor",c.PostForm("Vendor"),
		"Color",c.PostForm("Color"),
		"Color",c.PostForm("Type"),
		"Description",c.PostForm("Description"),
		"Name",c.PostForm("Name"),
		"isBroken",c.PostForm("isBroken"),
		"Size",c.PostForm("Size"),
		)
	isBroken, _ := strconv.Atoi(c.PostForm("IsBroken"))
	db.Model(&material).Update("IsBroken", isBroken)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Material updated successfully!"})
}

// deleteMaterial remove a material
func deleteMaterial(c *gin.Context) {
	var material materialModel
	materialID := c.Param("ID")

	db.First(&material, materialID)

	if material.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No material found!"})
		return
	}

	db.Delete(&material)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Material deleted successfully!"})
}