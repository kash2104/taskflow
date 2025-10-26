package runners

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/kash2104/taskflow/worker/sandbox"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func CompileJava(taskId primitive.ObjectID, code string) error{
	boxId:= ConvertNumberFromHex(taskId)
	
	_, err := WriteCodeToSandbox(taskId,code,"cpp");

	if err != nil{
		return err;
	}

	// compile code
	compile := exec.Command("javac","main_"+boxId+".java");
	compile.Dir = basePath + boxId + "/box"

	if err := compile.Run(); err != nil{
		return fmt.Errorf("compilation error!! %v", err);
	}
	return nil;
}


func RunJava(taskId primitive.ObjectID, code string) (string, string, error) {
	boxId:= ConvertNumberFromHex(taskId)
	
	// _, err := WriteCodeToSandbox(taskId,code,"java");

	// if err != nil{
	// 	return "","",err;
	// }


	//run code
	//for java, the class name should be Main_<boxId> inorder to run it
	cmd := exec.Command("java", "Main");
	cmd.Dir = basePath + boxId + "/box"
	stdoutpipe,_ := cmd.StdoutPipe()
	stderrpipe,_ := cmd.StderrPipe()

	cmd.Start()

	stdout, _ := io.ReadAll(stdoutpipe);
	stderr,_ := io.ReadAll(stderrpipe);

	err := cmd.Wait();

	defer sandbox.CleanupSandbox(taskId);

	return string(stdout), string(stderr), err;

}