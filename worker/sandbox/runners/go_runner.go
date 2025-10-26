package runners

import (
	"io"
	"os/exec"

	"github.com/kash2104/taskflow/worker/sandbox"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RunGo(taskId primitive.ObjectID, code string) (string, string, error) {
	boxId:= ConvertNumberFromHex(taskId)
	
	// _, err := WriteCodeToSandbox(taskId,code,"go");

	// if err != nil{
	// 	return "","",err;
	// }

	cmd := exec.Command("go","run", "main_"+boxId+".go");
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