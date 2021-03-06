package stat

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/zxfonline/fileutil"
	"github.com/zxfonline/golog"
	"github.com/zxfonline/timefix"
)

var mystat *log.Logger
var statchan chan string
var statFile *os.File
var appName string

func init() {
	appName = path.Clean(os.Args[0])
	_, appName = filepath.Split(appName)
	names := strings.Split(appName, ".")
	appName = names[0]

	var err error
	statFile, err = fileutil.OpenFile("../log/"+appName+time.Now().Format("20060102")+".stat", fileutil.DefaultFileFlag, fileutil.DefaultFileMode)
	if err != nil {
		log.Fatalln("open file error !")
		os.Exit(-1)
		return
	}
	statchan = make(chan string, 10000)
	mystat = log.New(statFile, "", log.Ldate|log.Ltime)

	go writeloop()
}

func Record(actcion string, v ...interface{}) {
	statchan <- fmt.Sprintf("[%15s]", strings.ToUpper(actcion)) + fmt.Sprint(v)
	if wait := len(statchan); wait > cap(statchan)/10*5 && wait%100 == 0 {
		golog.Warnf("state record taskchan process,waitchan:%d/%d.", wait, cap(statchan))
	}
}

func writeloop() {
	pm := time.NewTimer(time.Duration(timefix.NextMidnight(time.Now(), 1).Unix()-time.Now().Unix()) * time.Second)
	for {
		select {
		case str := <-statchan:
			mystat.Println(str)
			select {
			case <-pm.C:
				if statFile1, err := fileutil.OpenFile("../log/"+appName+time.Now().Format("20060102")+".stat", fileutil.DefaultFileFlag, fileutil.DefaultFileMode); err != nil {
					log.Printf("open file err:%v\n", err)
				} else {
					statFile.Close()
					mystat.SetOutput(statFile1)
					statFile = statFile1
					pm.Reset(time.Second * 24 * 60 * 60)
				}
			default:
			}
		}
	}
}
