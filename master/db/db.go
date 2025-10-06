package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRequest struct{
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Code string `json:"code" bson:"code"`
	Language string `json:"language" bson:"language"`
	Status string `json:"status" bson:"status"`
	Result string `json:"result" bson:"result"`
}


var collection *mongo.Collection

func Connect(url string) error{

	clientOptions := options.Client().ApplyURI(url)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil{
		log.Panicf("Failed to connect to Mongodb %v", err);;
	}

	err = client.Ping(context.Background(),nil);
	if err != nil{
		log.Panicf("Failed to establish connection with mongo %v", err);
	}

	db := client.Database("taskflow")

	collection = db.Collection("user_requests")

	fmt.Println("Connected to MongoDB!")

	return nil
}

func AddtoDB(req UserRequest) (UserRequest, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second);
	defer cancel();

	// if req.Code == "" || req.Language == ""{
	// 	return req, fmt.Errorf("code and language fields cannot be empty");
	// }

	// req.Status = "Pending";
	// req.Result = "";

	_, err := collection.InsertOne(ctx, req);

	if err != nil{
		return req, fmt.Errorf("failed to insert document: %v", err);
	}

	return req,nil;
}


func GetResult(id primitive.ObjectID) (UserRequest, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second);
	defer cancel();

	if id.IsZero(){
		return UserRequest{}, fmt.Errorf("id cannot be empty");
	}

	var result UserRequest;

	err := collection.FindOne(ctx, bson.M{"_id":id}).Decode(&result);

	if err != nil{
		if err == mongo.ErrNoDocuments{
			return result, fmt.Errorf("no document found with the given id");
		}
		return result, fmt.Errorf("failed to find document: %v", err);
	}

	return result, nil;
}

func CheckPending(taskId primitive.ObjectID) (bool, error){

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second);
	defer cancel();

	if taskId.IsZero(){
		return false, fmt.Errorf("id cannot be empty");
	}

	var result UserRequest;

	err := collection.FindOne(ctx, bson.M{"_id":taskId}).Decode(&result);

	if err != nil{
		if err == mongo.ErrNoDocuments{
			return false, fmt.Errorf("no document found with the given id");
		}
		return false, fmt.Errorf("failed to find document: %v", err);
	}

	if result.Status == "Pending"{
		return false, nil;
	}

	return true, nil;

}

func UpdateStatus(id primitive.ObjectID) (UserRequest, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second);
	defer cancel();

	if id.IsZero(){
		return UserRequest{}, fmt.Errorf("id cannot be empty");
	}

	var result UserRequest;

	filter := bson.M{"_id":id, "status":"Pending"};
	err := collection.FindOne(ctx,filter).Decode(&result);

	if err != nil{
		if err == mongo.ErrNoDocuments{
			return result, fmt.Errorf("no document found with the given id or the task is not pending");
		}
		return result, fmt.Errorf("failed to find document: %v", err);
	}

	
	update := bson.M{
		"$set": bson.M{
			"status": "Completed",
		},
	}

	_, err = collection.UpdateOne(ctx,filter,update);
	if err != nil{
		return result, fmt.Errorf("failed to update document: %v", err);
	}

	result.Status = "Completed";
	
	return result, nil;
}