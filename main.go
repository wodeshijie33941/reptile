// Command click is a chromedp example demonstrating how to use a selector to
// click on an element.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"chromedp/reptile"

	"github.com/chromedp/chromedp"
)

const (

	// your target web addr
	HOMEADDR   = `https://yiqixie.com/d/home/fcACKS7AzHdvN0htrsNcqhJMr`

	// dowload dir
	FOLDERNAME = "/download"
)

//var downloadArray map[int]FileInfo = make(map[int]FileInfo)

func main() {

	var err error
	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt) //chromedp.WithLog(log.Printf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("journey is start soon...")
	fmt.Println("simulate login...")
	// run task list
	var res string
	err = c.Run(ctxt, reptile.Onlogin(&res))
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Print(res)
	fmt.Println("login done!")

	//AnalyseContent(res)
	fmt.Println("analyse html data...")

	curPathStr := getCurrentDirectory()
	fmt.Println("application path:", curPathStr)
	curPathStr += FOLDERNAME
	fmt.Println("download path:", curPathStr)
	//return
	filelist := reptile.GetDownloadList(c, ctxt, HOMEADDR, curPathStr)
	fmt.Println("analyse done!")
	reptile.Download(filelist, curPathStr)
	// add end

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n\n\n\n\n\ntask done,bye~~~")
	time.Sleep(8 * time.Second)
}

/*
Get current dir
*/
func getCurrentDirectory() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	/*if err != nil {
		beego.Debug(err)
	}*/
	return strings.Replace(dir, "\\", "/", -1)
}
