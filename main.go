package main 

import (
	"fmt"
	"io/ioutil"
	"github.com/gin-gonic/gin"
)

func PostHomePage(c *gin.Context) {
	// var req createPerson
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	fmt.Println("error!!!!")
	// }

	// arg := createPerson{
	// 	name: req.name,
	// 	age: req.age,
	// }

	// fmt.Println("name",req.name, req.age)

	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println((string(value)))

	c.JSON(200, gin.H{
		"message": string(value),
	})

}


func createPerson(c *gin.Context) {
	name := c.Param("name")
	birthday := c.Param("birthday")

	c.JSON(200, gin.H{
		"name": name,
		"birthday": birthday,
	})
}

func getPersonById(c *gin.Context) {
	id := c.Param("id")

	c.JSON(200, gin.H{
		"id": id,
	})
}

func updatePersonBirthday(c *gin.Context) {
	id := c.Param("id")
	birthday := c.Param("birthday")

	c.JSON(200, gin.H{
		"id": id,
		"birthday": birthday,
	})
}

func deletePersonById(c *gin.Context) {
	id := c.Param("id")

	c.JSON(200, gin.H{
		"id to delete": id,
	})
}


//crud - create, read, update, delete


// func main() {
// 	r := gin.Default()
// 	// r.POST("/", PostHomePage)

// 	r.POST("/:name/:birthday", createPerson)
// 	r.GET("/:id", getPersonById)
// 	r.PUT("/:id/:birthday", updatePersonBirthday)
// 	r.DELETE("/:id", deletePersonById)

// 	r.Run()
// }