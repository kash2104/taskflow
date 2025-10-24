package runners

import (
	"fmt"
	"os"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


func ConvertNumberFromHex(taskId primitive.ObjectID) string {
	number, _ := strconv.ParseInt(taskId.Hex()[len(taskId.Hex())-2:], 16, 64)

	numberString := fmt.Sprintf("%d", number)
	return numberString
}


func WriteCodeToSandbox(taskId primitive.ObjectID, code string, language string) (string ,error) {
	boxId:= ConvertNumberFromHex(taskId)
	
	fileName := fmt.Sprintf("/var/local/lib/isolate/%s/box/main_%s.%s",boxId,boxId,language);

	err := os.WriteFile(fileName,[]byte(code),0644);
	if err != nil{
		return "",fmt.Errorf("failed to write code to sandbox file %v", err);
	}

	return fileName,nil
}

// func WriteCodeToSandbox(taskId primitive.ObjectID, code string, language string) (string, error) {
//     boxId := ConvertNumberFromHex(taskId)
//     boxPath := fmt.Sprintf("/var/local/lib/isolate/%s/box", boxId)

//     // ensure directory exists
//     if err := os.MkdirAll(boxPath, 0755); err != nil {
//         return "", fmt.Errorf("failed to create sandbox dir %s: %w", boxPath, err)
//     }

//     ext := language
//     if language == "cpp" {
//         ext = "cpp"
//     } else if language == "java" {
//         ext = "java"
//     }

//     fileName := fmt.Sprintf("main_%s.%s", boxId, ext)
//     fullPath := filepath.Join(boxPath, fileName)

//     if err := os.WriteFile(fullPath, []byte(code), 0644); err != nil {
//         return "", fmt.Errorf("failed to write code to sandbox file %s: %w", fullPath, err)
//     }

//     // debug log
//     fmt.Printf("Wrote code to sandbox: %s\n", fullPath)
//     return fullPath, nil
// }


// func GetExecutableToSandbox(taskId string, language string) error {
// 	boxId:= ConvertNumberFromHex(taskId)
// 	
// 	var cmd *exec.Cmd
// 	switch language {
// 	case "cpp":
// 		cmd = exec.Command("isolate", "--box-id="+boxId, "--run", "--", "g++", "main_"+boxId+".cpp", "-o", "main_"+boxId)
// 	case "java":
// 		cmd = exec.Command("isolate", "--box-id="+boxId, "--run", "--", "javac", "Main_"+boxId+".java")

// 	}
// 	if err := cmd.Run(); err != nil {
// 		return fmt.Errorf("failed to create a executable cpp code: %v", err)
// 	}
// 	return nil
// }