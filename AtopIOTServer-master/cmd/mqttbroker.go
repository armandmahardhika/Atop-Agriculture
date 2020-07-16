package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// func main() {
//     os := runtime.GOOS
//     switch os {
//     case "windows":
//         fmt.Println("Windows")
//     case "darwin":
//         fmt.Println("MAC operating system")
//     case "linux":
//         fmt.Println("Linux")
//     default:
//         fmt.Printf("%s.\n", os)
//     }
// }

func GetMosquittoInfo() (bson.M, error) {

	lsof := exec.Command("lsof", "-nP", "-i")
	grepMosquitto := exec.Command("grep", "mosquitto")
	grepListen := exec.Command("grep", "LISTEN")
	lsofOut, _ := lsof.StdoutPipe()
	grepMosquittoOut, _ := grepMosquitto.StdoutPipe()
	err := lsof.Start()
	if err != nil {
		return bson.M{}, err
	}
	grepMosquitto.Stdin = lsofOut
	grepListen.Stdin = grepMosquittoOut

	err = grepMosquitto.Start()
	if err != nil {
		return bson.M{}, err
	}

	out, err := grepListen.Output()
	str := string(out)
	m := strings.Split(str, "\n")
	fmt.Println(m)

	return bson.M{}, nil
}
