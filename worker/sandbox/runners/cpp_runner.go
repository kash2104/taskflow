package runners

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/kash2104/taskflow/worker/sandbox"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func CompileCpp(taskId primitive.ObjectID, code string) error{
	boxId:= ConvertNumberFromHex(taskId)
	
	// _, err := WriteCodeToSandbox(taskId,code,"cpp");

	// if err != nil{
	// 	return err;
	// }

	// compile code
	// compile := exec.Command("isolate","--box-id="+boxId,"--run","--", "g++",basePath+boxId+"/box/main_"+boxId+".cpp", "-o", "main_"+boxId);

	compile := exec.Command("g++", "main_" + boxId + ".cpp", "-o", "main_" + boxId);
	compile.Dir = basePath + boxId + "/box"

	if err := compile.Run(); err != nil{
		return fmt.Errorf("compilation error!! %v", err);
	}
	return nil;
}

func RunCpp(taskId primitive.ObjectID) (string, string, error) {
	boxId:= ConvertNumberFromHex(taskId)
	

	//run code
	cmd := exec.Command("isolate", "--box-id="+boxId,"--time=4", "--run", "./main_"+boxId);
	stdoutpipe,_ := cmd.StdoutPipe()
	stderrpipe,_ := cmd.StderrPipe()

	cmd.Start()

	stdout, _ := io.ReadAll(stdoutpipe);
	stderr,_ := io.ReadAll(stderrpipe);

	err := cmd.Wait();

	defer sandbox.CleanupSandbox(taskId);

	return string(stdout), string(stderr), err;

}