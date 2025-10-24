package sandbox

import (
	"fmt"
	"os/exec"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertNumberFromHex(taskId primitive.ObjectID) string {
	number, _ := strconv.ParseInt(taskId.Hex()[len(taskId.Hex())-2:], 16, 64)

	numberString := fmt.Sprintf("%d", number)
	return numberString
}

func InitSandbox(taskId primitive.ObjectID) error {
	boxId := ConvertNumberFromHex(taskId);
	cmd := exec.Command("isolate", "--init", "--box-id="+boxId)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to init sandbox %v", err)
	}
	fmt.Printf("Sandbox initialised!!")
	return nil
}

func CleanupSandbox(taskId primitive.ObjectID) error {
	boxId := ConvertNumberFromHex(taskId)

	cmd := exec.Command("isolate", "--cleanup", "--box-id="+boxId)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clean sandbox %v", err)
	}
	fmt.Printf("Sandbox cleaned!!")
	return nil
}