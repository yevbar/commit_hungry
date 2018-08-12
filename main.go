package main

import (
	"os"
	"os/exec"
	"fmt"
	"strings"
	"time"
	"github.com/MarinX/keylogger"
)

func IsDesired(Input string) bool {
	DesiredInputs := "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-=[]\\;',./`"
	InputLength := len(Input)
	//fmt.Println(Input)
	if (InputLength != 1) {
		return false
	}
	return strings.ContainsAny(DesiredInputs, Input)
}

func FormatTime(t time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func Commit(Input string, date time.Time) {
	Formatted := FormatTime(date)

	os.Setenv("GIT_AUTHOR_DATE", Formatted)
	os.Setenv("GIT_COMMIT_DATE", Formatted)

	//fmt.Println(Formatted)

	msg := fmt.Sprintf("Typed " + Input)
	//fmt.Println(msg)
	exec.Command("git", "commit", "--allow-empty", "-m " + msg).Start()
	//exec.Command("git", "push", "origin", "master").Run()
	//fmt.Println("Pushed")
}

func main() {
	if (len(os.Args) != 2) {
		fmt.Println("Usage ./main filename")
		return
	}

	cmd := exec.Command("vim", os.Args[1])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Start()

	if err != nil {
		panic(err)
	}


	devs, err := keylogger.NewDevices()
	if err != nil {
		panic(err)
	}

	for _, val := range devs {
		fmt.Println("Id->", val.Id, "Device->", val.Name)
	}

	//keyboard device file, on your system it will be diffrent!
	rd := keylogger.NewKeyLogger(devs[3])

	in, err := rd.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	//Our_PID := cmd.Process.Pid

	var LastPressed string = ""
	TheTime := time.Now()
	TheTime = TheTime.AddDate(-1,0,-5)

	for i := range in {
		/*_ , err := os.FindProcess(int(Our_PID))
		if err != nil {
			exec.Command("git", "push", "origin", "master").Run()
			return
		}*/
		if i.Type == keylogger.EV_KEY {
			var CurrentlyPressed = i.KeyString()
			if (CurrentlyPressed != "") {
				if (CurrentlyPressed != LastPressed) {
					//fmt.Println("{" + CurrentlyPressed + "}")
					if (IsDesired(CurrentlyPressed)) {
						//CurrentString += CurrentlyPressed
						//fmt.Println(CurrentString)
						//exec.Command("curl", "127.0.0.1:5000/" + CurrentlyPressed).Run()
						Commit(CurrentlyPressed, TheTime)
						TheTime = TheTime.AddDate(0,0,1)
						if (time.Now().Before(TheTime)) {
							TheTime = time.Now()
							TheTime = TheTime.AddDate(-1,0,-5)
						}
						LastPressed = CurrentlyPressed
					} else if (CurrentlyPressed == "SPACE") {
						//CurrentString += " "
						//fmt.Println(CurrentString)
						//exec.Command("curl", "127.0.0.1:5000/space").Run()
						Commit("SPACE", TheTime)
						TheTime = TheTime.AddDate(0,0,1)
						if (time.Now().Before(TheTime)) {
							TheTime = time.Now()
							TheTime = TheTime.AddDate(-1,0,-5)
						}
						LastPressed = CurrentlyPressed
					} else if (CurrentlyPressed == "BS") {
						//CurrentString = CurrentString[:len(CurrentString)-1]
						//fmt.Println(CurrentString)
						//exec.Command("curl", "127.0.0.1:5000/backspace").Run()
						Commit("BACKSPACE", TheTime)
						TheTime = TheTime.AddDate(0,0,1)
						if (time.Now().Before(TheTime)) {
							TheTime = time.Now()
							TheTime = TheTime.AddDate(-1,0,-5)
						}
						LastPressed = CurrentlyPressed
					}
				} else {
					LastPressed = ""
				}
			}
		}

	}
}
