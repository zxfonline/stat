package stat

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/zxfonline/fileutil"
	"github.com/zxfonline/timefix"
)

var mystat *log.Logger
var statchan chan string
var statFile *os.File
var appName string

func init() {
	appName = strings.Replace(os.Args[0], "\\", "/", -1)
	_, name := path.Split(appName)
	names := strings.Split(name, ".")
	appName = names[0]

	fileName := "../log/" + appName + time.Now().Format("20060102") + ".stat"
	var err error
	statFile, err = fileutil.OpenFile(fileName, fileutil.DefaultFileFlag, fileutil.DefaultFileMode)
	if err != nil {
		log.Fatalln("open file error !")
		os.Exit(-1)
		return
	}
	statchan = make(chan string, 1000)
	mystat = log.New(statFile, "", log.Ldate|log.Ltime)

	go writeloop()
}

func Record(actcion string, v ...interface{}) {
	statchan <- fmt.Sprintf("[%15s]", strings.ToUpper(actcion)) + fmt.Sprint(v)
}

func writeloop() {
	pm := time.NewTimer(time.Duration(timefix.NextMidnight(time.Now(), 1).Unix()-time.Now().Unix()) * time.Second)
	for {
		select {
		case str := <-statchan:
			mystat.Println(str)
		case <-pm.C:
			// 关闭原来的文件
			statFile.Close()

			time.Sleep(time.Second * 1)

			fileName := "../log/" + appName + time.Now().Format("20060102") + ".stat"
			var err error
			statFile, err = fileutil.OpenFile(fileName, fileutil.DefaultFileFlag, fileutil.DefaultFileMode)
			if err != nil {
				log.Fatalln("open file error !")
				os.Exit(-1)
				return
			}

			mystat.SetOutput(statFile)

			pm.Reset(time.Second * 24 * 60 * 60)
		}
	}
}
