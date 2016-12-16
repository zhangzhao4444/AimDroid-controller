package android

import (
	"bufio"
	"log"
	"monidroid/util"
	"os"
	"path"
	"strings"
)

var adb string = "adb"

//Launch an application
func LaunchApp(pck, act string) error {
	cmd := adb + " shell am start -n " + pck + "/" + act
	_, err := util.ExeCmd(cmd)
	if err != nil {
		log.Println(err)
	}
	return err
}

//kill ape
//func KillProApp(pck string) {
//	cmd := adb + " shell ps | grep " + pck
//	out, err := util.ExeCmd(cmd)
//	if err != nil || len(out) <= 0 {
//		log.Println(pck, "is not running.")
//		return
//	}

//	iterms := strings.Fields(out)
//	if len(iterms) >= 9 {
//		pid := iterms[1]
//		cmd = adb + " shell su -c kill " + pid
//		_, err = util.ExeCmd(cmd)

//		if err != nil {
//			log.Println("Cannot kill", pck)
//		}
//	}

//}

//Kill an application
func KillApp(pck string) error {
	cmd := adb + " shell su -c am force-stop " + pck
	_, err := util.ExeCmd(cmd)
	return err
}

func ClearApp(pck string) error {
	cmd := adb + " shell pm clear " + pck
	_, err := util.ExeCmd(cmd)
	return err
}

func InitADB(sdk string) {
	adb = path.Join(sdk, "platform-tools/adb")
}

//start logcat
func StartLogcat() (*bufio.Reader, error) {

	_, err := util.ExeCmd(adb + " logcat -c")
	if err != nil {
		return nil, err
	}

	cmd := util.CreateCmd(adb + " logcat Monitor_Log:V art:I *:S")

	// Create stdout, stderr streams of type io.Reader
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	// Start command
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	read := bufio.NewReader(stdout)
	return read, nil
}

//push file to the device
func PushFile(src, dst string) error {
	if _, err := os.Stat(src); err == nil {
		cmd := adb + " push " + src + " " + dst
		_, err = util.ExeCmd(cmd)
		return err
	} else {
		return err
	}
}

//remove file from the device
func RemoveFile(dst string) error {
	cmd := adb + " shell rm " + dst
	_, err := util.ExeCmd(cmd)
	return err
}

func StartMonkey(pkg string) (string, error) {
	cmd := adb + " shell monkey --throttle 700 -p " + pkg + " -v 2000"
	//cmd := GetADBPath(sdk) + " shell monkey --pct-touch 80 --pct-trackball 20 --throttle 300 --uiautomator -v 1000"
	//cmd := GetADBPath(sdk) + " shell monkey --throttle 300 --uiautomator-dfs -v 100"

	out, err := util.ExeCmd(cmd)
	return out, err
	//time.Sleep(time.Millisecond * 10000)
	//return "", nil
}

//start ape
func StartApe(port string) (*bufio.Reader, error) {

	content := adb + " shell ape --ignore-crashes --ignore-timeouts --ignore-native-crashes --port " + port
	cmd := util.CreateCmd(content)

	// Create stdout, stderr streams of type io.Reader
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	// Start command
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	read := bufio.NewReader(stdout)
	return read, nil
}

//kill ape
func KillApe() {
	cmd := adb + " shell ps | grep com.android.commands.monkey"
	out, err := util.ExeCmd(cmd)
	if err != nil || len(out) <= 0 {
		log.Println("Ape is not running.")
		return
	}

	iterms := strings.Fields(out)
	if len(iterms) >= 9 {
		pid := iterms[1]
		cmd = adb + " shell su -c kill " + pid
		_, err = util.ExeCmd(cmd)

		if err != nil {
			log.Println("Cannot kill ape!", err)
		}
	}

}

//adb forward
func Forward(pcPort, mobilePort string) error {
	cmd := adb + " forward tcp:" + pcPort + " tcp:" + mobilePort
	_, err := util.ExeCmd(cmd)
	return err
}

//get current focused activity
func GetCurrentActivity(want string) string {
	name := ""
	cmd := adb + " shell dumpsys activity activities | grep mFocusedActivity"
	out, err := util.ExeCmd(cmd)
	if err != nil {
		log.Println("Cannot get current activity!", err)
		return name
	}
	iterms := strings.Split(out, " ")
	for i, iterm := range iterms {
		if iterm == "u0" {
			i++
			if i < len(iterms) {
				name = iterms[i]
			}
			break
		}
	}
	//log.Println(name)
	if len(name) > 0 {
		names := strings.Split(name, "/")
		if len(names) == 2 {
			name = names[1]
		}
	}
	if name == ".permission.ui.GrantPermissionsActivity" && len(want) > 0 {
		name = want
	}
	return name
}
