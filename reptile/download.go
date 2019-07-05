package reptile

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

const (
	LOGINURL    = "https://yiqixie.com/e/useradmin/emlogin"
	FOLDERNAME  = "./download"
	BASEDOWNURL = "https://yiqixie.com/d/home/"
)

type FileInfo struct {
	FileId   string // 文件ID
	FileType string // 文件类型
	FilePath string // 文件路径
	FileName string // 文件名称(图片用)
	FileAddr string // 文件地址(图片用)
}

func GetDownloadList(c *chromedp.CDP, ctxt context.Context, targeturl string, fatherPath string) []FileInfo {
	// 生成文件夹
	CreateFolder(fatherPath)
	// run task list
	var res string
	err := c.Run(ctxt, DownloadPage(&res, targeturl))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(res)
	return AnalyseContent(res, c, ctxt, fatherPath)
}

func DownloadPage(res *string, targetUrl string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(targetUrl),
		chromedp.WaitVisible(`.i-o-sr-ho`),
		chromedp.OuterHTML(`.office-infinite-list-items-wrapper-outter`, res),
		chromedp.Sleep(2 * time.Second),
	}
}

func AnalyseContent(res string, c *chromedp.CDP, ctxt context.Context, fatherPath string) []FileInfo {

	var fileArray []FileInfo
	var downloadArray map[int]FileInfo = make(map[int]FileInfo)
	// add query
	var docio io.Reader = strings.NewReader(res)
	doc, err := goquery.NewDocumentFromReader(docio)
	if err != nil {
		log.Fatal(err)
	}
	// Find the review items
	doc.Find(".i-o-ho-cz-nc-td-rib").Each(func(i int, s *goquery.Selection) {
		fmt.Printf("Review %d: %s -\n", i, s.Text())
		var file FileInfo
		file.FileId = s.Text()
		file.FileType = "0"
		downloadArray[i] = file

		fileArray = append(fileArray, file)
		fmt.Println("target file:", file.FileId)
	})

	// Find folder name
	doc.Find(".cz-io").Each(func(i int, s *goquery.Selection) {
		fmt.Printf("Review %d: %s -\n", i, s.Text())
		var file FileInfo
		file.FileId = downloadArray[i].FileId
		file.FilePath = fatherPath + "/" + s.Text()
		file.FileName = s.Text()
		downloadArray[i] = file

		fileArray[i] = file

		fmt.Println("file path:", fileArray[i].FilePath)
	})

	// Find the review items
	doc.Find(".i-o-ho-cz-nc-td-xhb").Each(func(i int, s *goquery.Selection) {
		fmt.Printf("Review %d: %s -\n", i, s.Text())
		var file FileInfo
		file = downloadArray[i]
		//file.FilePath = downloadArray[i].FilePath
		file.FileType = s.Text()
		downloadArray[i] = file

		fileArray[i] = file

		// 对文件夹再次分析
		if file.FileType == "1" {

			// 创建文件夹
			CreateFolder(file.FilePath)

			fmt.Println("folder,start recursion...")
			newdownAddr := BASEDOWNURL + file.FileId
			subfilelist := GetDownloadList(c, ctxt, newdownAddr, file.FilePath)
			if len(subfilelist) > 0 {
				for _, val := range subfilelist {
					fileArray = append(fileArray, val)
				}
			}
			fmt.Println("recursion done!")
		} else if file.FileType == "4" {
			// 非文件夹,要重置路径
			//fileArray[i].FilePath = fatherPath

			picUrl := "https://yiqixie.com/d/home/" + fileArray[i].FileId
			picAddr := GetPictureSrc(picUrl, c, ctxt)
			fmt.Println("图片地址:", picAddr)
			fileArray[i].FileAddr = picAddr
		} else {
			// 非文件夹,要重置路径
			fileArray[i].FilePath = fatherPath
		}

	})

	//	fmt.Println(downloadArray)
	fmt.Println("current folder done!")

	return fileArray
}

// 下载文件
func Download(filelist []FileInfo, folder string) {
	fmt.Println("start download...")

	var client http.Client
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client.Jar = jar

	v := url.Values{}
	v.Set("email", "wodeshijie3394@126.com")
	v.Set("password", "1qaz1QAZ")
	//v.Set("redirectTo", "/?jumpToDoclist=true")
	u := ioutil.NopCloser(strings.NewReader(v.Encode()))
	r, err := client.Post(LOGINURL, "application/x-www-form-urlencoded", u)
	if err != nil {
		fmt.Println("http post err:", err)
	}
	defer r.Body.Close()
	fmt.Println(r.StatusCode)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("http body err:", err)
	}
	fmt.Println(string(body))
	r.Body.Close()

	fmt.Println("下载流程已经开启!")

	var baseUrl string = "https://yiqixie.com/d/export/"
	var exelUrl string = "https://yiqixie.com/s/export/"
	var pptUrl string = "https://yiqixie.com/p/export/"
	var formatStr string = "?format="
	for _, value := range filelist { //downloadArray
		if value.FileId == "" || value.FileId == " " {
			continue
		}
		if value.FileType == "2" {
			if value.FileId == "" || value.FileId == " " {
				continue
			}
			temUrl := baseUrl + value.FileId + formatStr + "docx"
			fmt.Println("set up download add:", temUrl)
			DownloadFile(&client, temUrl, value.FilePath)
		} else if value.FileType == "4" {
			DownloadImg(value.FilePath, value.FileAddr)
		} else if value.FileType == "5" {
			temUrl := exelUrl + value.FileId + formatStr + "xlsx"
			fmt.Println("set up download add:", temUrl)
			DownloadFile(&client, temUrl, value.FilePath)
		} else if value.FileType == "7" {
			temUrl := pptUrl + value.FileId + formatStr + "pptx"
			fmt.Println("set up download add:", temUrl)
			DownloadFile(&client, temUrl, value.FilePath)
		}
	}
}

func DownloadFile(client *http.Client, downloadUrl string, folderPath string) {

	fmt.Println("downloading file...")
	resp, err := client.Get(downloadUrl)

	if err != nil {
		// handle error
		fmt.Println("下载文件出错:", err, "该文件可能已经被删除!")
		return
	}

	defer resp.Body.Close()

	//if err != nil {
	//	fmt.Println("下载出错:", err)
	//}

	//fmt.Println(resp.Header.Get("Content-Disposition"))
	downloadFileName := GetFileName(resp.Header.Get("Content-Disposition"))
	downloadFileName = downloadFileName[1 : len(downloadFileName)-1]
	fmt.Println("analyze name:", downloadFileName)

	downloadFileName = folderPath + "/" + downloadFileName
	downLoadFile, err := os.Create(downloadFileName)
	if err != nil {
		panic(err)
	}
	defer downLoadFile.Close()
	fileSize, _ := io.Copy(downLoadFile, resp.Body)
	fmt.Println("file size:", fileSize)
	fmt.Println("current file download done!")
}

func GetFileName(content string) string {
	contentArray := strings.Split(content, ";")
	//fmt.Println(contentArray)
	for _, v := range contentArray {
		targetArray := strings.Split(v, "=")
		fmt.Println(targetArray)
		if targetArray[0] == " filename" {
			return targetArray[1]
		}
	}

	return ""
}

// 获取图片src地址
func GetPictureSrc(url string, c *chromedp.CDP, ctxt context.Context) string {
	var res string
	err := c.Run(ctxt, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`.i-ofb-sr`),
		chromedp.OuterHTML(`.e7-m`, &res),
		chromedp.Sleep(2 * time.Second),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(res)

	var docio io.Reader = strings.NewReader(res)
	doc, err := goquery.NewDocumentFromReader(docio)
	if err != nil {
		log.Fatal(err)
	}

	var src string
	doc.Find("img").Each(func(_ int, tag *goquery.Selection) {

		if tag.HasClass("e7-an") {
			src, _ = tag.Attr("src")
		} else {
			src, _ = tag.Attr("data-actualsrc")
		}
	})
	return src
}

// 下载图片
func DownloadImg(imageName, imageUrl string) {
	/*path := strings.Split(url, "/")
	var name string
	if len(path) > 1 {
		name = path[len(path)-1]
	}
	fmt.Println(name)*/
	out, _ := os.Create(imageName)
	defer out.Close()
	resp, _ := http.Get(imageUrl)
	defer resp.Body.Close()
	pix, _ := ioutil.ReadAll(resp.Body)
	io.Copy(out, bytes.NewReader(pix))

	fmt.Println("image download done!!")
}
