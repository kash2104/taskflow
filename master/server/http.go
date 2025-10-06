package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kash2104/taskflow/master/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


var TaskChannel chan db.UserRequest


func StartServer(taskChannel chan db.UserRequest){

	router := gin.New();

	PORT := os.Getenv("PORT")

	TaskChannel = taskChannel

	router.POST("/submit", HandleSubmission);
	router.GET("/findResult", HandleGetResult)
	router.GET("/checkPending", HandleCheckPending);
	// router.PATCH("/updateStatus", HandleUpdateStatus);

	router.Run(PORT)	
}


func HandleSubmission(c *gin.Context){

	var userRequest db.UserRequest

	// decoder := json.NewDecoder(c.Request.Body);
	// err := decoder.Decode(&userRequest);
	language := c.PostForm("language");
	code := c.PostForm("code");

	if language == "" || code == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Missing language or code in form data",
        })
        return
    }
	

	userRequest.Id = primitive.NewObjectID();
	userRequest.Status = "Pending"
	userRequest.Code = code;
	userRequest.Language = language;

	//adding to db
	newRequest,err := db.AddtoDB(userRequest);
	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":"false",
			"error":"error while adding to db",
			"actual_error":err,
		})
		return;
	}

	fmt.Println("added request to database")

	//adding to taskchannel
	TaskChannel <- newRequest
	fmt.Println("added request to task queue")


	c.JSON(http.StatusCreated,gin.H{
		"success":true,
		"message":"added new request to db",
		"data":newRequest,
	})

}

func HandleGetResult(c *gin.Context){
	var id primitive.ObjectID
	decoder := json.NewDecoder(c.Request.Body);
	err := decoder.Decode(&id);

	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"success":false,
			"error":"id not given",
		})
		return;
	}


	user,err:= db.GetResult(id);

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":false,
			"error":"error while fetching result from db",
			"actual_error":err,
		})
		return;
	}

	c.JSON(http.StatusOK, gin.H{
		"success":true,
		"data":user,
	})
}

func HandleCheckPending(c *gin.Context){
	var id primitive.ObjectID
	decoder := json.NewDecoder(c.Request.Body);
	err := decoder.Decode(&id);

	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"success":false,
			"error":"id not found",
		})
		return;
	}


	pending,err:= db.CheckPending(id);

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":false,
			"error":"error while checking status from db",
			"actual_error":err,
		})
		return;
	}

	c.JSON(http.StatusOK, gin.H{
		"success":true,
		"pending":pending,
	})
}

// func HandleUpdateStatus(c *gin.Context){
// 	var id primitive.ObjectID;

// 	decoder := json.NewDecoder(c.Request.Body);
// 	err := decoder.Decode(&id);

// 	if err != nil{
// 		c.JSON(http.StatusBadRequest,gin.H{
// 			"success":false,
// 			"error":"id not found",
// 		})
// 		return;
// 	}

// 	user, err := db.UpdateStatus(id);

// 	if err != nil{
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success":false,
// 			"error":"error while updating the status in db",
// 			"actual_error": err,
// 		})
// 		return;
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success":true,
// 		"data":user,
// 	})
// }