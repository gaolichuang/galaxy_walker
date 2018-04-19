/*
http://www.cnblogs.com/golove/p/3269099.html
正则如何设置flag???

基本页面解析部分
1.fetch
2.format unicode
3.slim 获取有效部分 html parser
4.regex提取有效部分
*/
package test

import (
        "galaxy_walker/internal/gcodebase/http_lib"
        LOG "galaxy_walker/internal/gcodebase/log"
        "galaxy_walker/internal/gcodebase/file"
        "regexp"
        "fmt"
        "strings"
        "path/filepath"
        "strconv"
        "time"
        "math/rand"
        "os"
)

const (
        kpagenum = 70
        projectPath = "/Users/gaolichuang/workspace/go/src/galaxy_walker/logs/"
)

var (
        listpagePath = filepath.Join(projectPath, "page")
        contentpagePath = filepath.Join(projectPath, "content")
        resultfile = filepath.Join(projectPath,"result.csv")
)

//var output = "/Users/gaolichuang/workspace/go/src/galaxy_walker/logs/2.html"
//var url = "http://bj.58.com/lipinxianhua/33236178538168x.shtml"
func fetchUrl(url string, output string) error {
        content, err := http_lib.GetUrl("GET", url, "")
        if err != nil {
                LOG.Errorf("fetch err %v", err)
                return err
        }
        return file.WriteStringToFile([]byte(content), output)
}
/*
func gencounterurl(url string) (error,string) {
        fmtstr := "http://jst1.58.com/counter?infoid=%s&userid=&uname=&sid=0&lid=0&px=0&cfpath="
        r, _ := regexp.Compile("/([0-9]+)x.(shtml)")
        id := ""
        strs := r.FindStringSubmatch(url)
        if len(strs) == 2 {
                id = strs[1]
                return nil,fmt.Sprintf(fmtstr,id)
        }
        // http://jst1.58.com/counter?infoid=33307266507976&userid=&uname=&sid=0&lid=0&px=0&cfpath=
        return fmt.Errorf("get id fail %s",url),""
}
*/
func parseHtml1(content string) (err error, cateValue, subcateValue, regionvalue, contectvalue, addressvalue string) {
        // category
        r, err := regexp.Compile("suUl.*基础信息列表")
        if err != nil {
                fmt.Println(err)
                return
        }

        strs := r.FindStringSubmatch(content)
        // 设置flag的模式。i m s U
        rr, err := regexp.Compile("(?U)li>.*</li")
        rrstr := rr.FindAllStringSubmatch(strs[0], -1)

        cate := "类别"
        subcate := "小类"
        region := "服务区域"
        contect := "联系人"
        address := "商家地址"

        //rrr,err := regexp.Compile("[^0-9A-Za-z%;,()!-:\"_>< =/]*")
        rrr, err := regexp.Compile(`(?U)([\p{Han}]+.*[\p{Han}]+.*)<`)
        for _, v := range rrstr {
                for _, vv := range v {
                        //fmt.Println(vv)
                        kk := rrr.FindAllString(vv, -1)
                        rrrstr := make([]string, 0)
                        for _, k := range kk {
                                if len(k) > 0 {
                                        rrrstr = append(rrrstr, k)
                                }
                        }
                        //fmt.Println(rrrstr)
                        if len(rrrstr) > 1 {
                                key := strings.Replace(rrrstr[0], "<", "", -1)
                                value := strings.Join(rrrstr[1:], ",")
                                value = strings.Replace(value, "<", "", -1)
                                //fmt.Println(key,value)
                                if strings.HasPrefix(key, cate) {
                                        cateValue = value
                                        cateValue = strings.Replace(cateValue,",", " ",-1)
                                } else if strings.HasPrefix(key, subcate) {
                                        subcateValue = value
                                        subcateValue = strings.Replace(subcateValue,",", " ",-1)
                                } else if strings.HasPrefix(key, region) {
                                        regionvalue = value
                                        regionvalue = strings.Replace(regionvalue,",", " ",-1)
                                } else if strings.HasPrefix(key, contect) {
                                        contectvalue = strings.Replace(rrrstr[1], "<", "", -1)
                                        contectvalue = strings.Replace(contectvalue,",", " ",-1)
                                } else if strings.HasPrefix(key, address) {
                                        addressvalue = value
                                        addressvalue = strings.Replace(addressvalue,",", " ",-1)
                                }

                        }
                }
        }
        return
}
func parseHtml2(content string) (err error, title,rqvalue, hyvalue, fwvalue, gwvalue string) {
        // category
        r, err := regexp.Compile("userinfo.*进入官网<")
        if err != nil {
                fmt.Println(err)
                return
        }

        strs := r.FindStringSubmatch(content)
        if len(strs) < 1 {
                err = fmt.Errorf("no title.")
                return
        }
        // title
        rr, err := regexp.Compile(`<h2>[\p{Han}]+.*[\p{Han}]+</h2>`)
        title = rr.FindString(strs[0])
        title = strings.TrimPrefix(title, "<h2>")
        title = strings.TrimSuffix(title, "</h2>")

        rrr, err := regexp.Compile(`<li>人气<em>([\d.]+)&#[\d]+;</em></li>`)
        rq := rrr.FindStringSubmatch(strs[0])
        if len(rq) > 1 {
                rqvalue = rq[1]
        }
        rrr, err = regexp.Compile(`<li>活跃<em>([\d.]+)&#[\d]+;</em></li>`)
        rq = rrr.FindStringSubmatch(strs[0])
        if len(rq) > 1 {
                hyvalue = rq[1]
        }
        rrr, err = regexp.Compile(`<li>服务<em>([\d.]+)&#[\d]+;</em></li>`)
        rq = rrr.FindStringSubmatch(strs[0])
        if len(rq) > 1 {
                fwvalue = rq[1]
        }

        rrr, err = regexp.Compile(`(?U)进入官网.*href="(.*)".*进入官网`)
        rq = rrr.FindStringSubmatch(strs[0])
        if len(rq) > 1 {
                gwvalue = rq[1]
        }
        return
}
func parseHtml3(content string) (err error,detail string){
        // category
        r, err := regexp.Compile("详细描述文字.*详细描述文字")
        if err != nil {
                fmt.Println(err)
                return
        }

        strs := r.FindStringSubmatch(content)
        fmt.Println(strs)
        return
}

func main() {
        // fetchListPage()
        //parseAndFetchAllPages()
        parseAllContent()
}
//////////////////
//////////////////

func parseAllContent() {
        outputFile,err := os.OpenFile(resultfile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
        if err != nil {
                LOG.Errorf("open file %s err %v",resultfile,err)
                return
        }
        titles := make([]string,0)
        titles = append(titles,"翻页")
        titles = append(titles,"内容页ID")
        titles = append(titles,"名称")
        titles = append(titles,"人气")
        titles = append(titles,"活跃")
        titles = append(titles,"服务")
        titles = append(titles,"官网")
        titles = append(titles,"分类")
        titles = append(titles,"子类")
        titles = append(titles,"区域")
        titles = append(titles,"联系人")
        titles = append(titles,"公司地址")
        outputFile.WriteString(strings.Join(titles,",") + "\n")
        defer outputFile.Close()
        for i:=1;i <= kpagenum;i++ {
                parseContent(strconv.Itoa(i),outputFile)
        }
}
func parseContent(pageid string,output *os.File) {
        fpath := filepath.Join(contentpagePath,pageid)
        fnames := file.ListFiles(fpath,"x")
        for _,fname := range fnames {
                c,_ := file.ReadFileToString(filepath.Join(fpath,fname))
                content := string(c)
                content = strings.Replace(content,"\n","",-1)
                content = strings.Replace(content,"\t","",-1)
                content = strings.Replace(content," ","",-1)
                content = strings.Replace(content,"&nbsp;","",-1)
                err,cateValue, subcateValue, regionvalue, contectvalue, addressvalue := parseHtml1(content)
                if err != nil {
                        LOG.Errorf("parseHtml1 Error %s,%s,%v",pageid,fname,err)
                        continue
                }
                err,title,rqvalue, hyvalue, fwvalue, gwvalue := parseHtml2(content)
                if err != nil {
                        LOG.Errorf("parseHtml2 Error %s,%s,%v",pageid,fname,err)
                        continue
                }
                fs := make([]string,0)
                fs = append(fs,pageid)
                fs = append(fs,fname)
                fs = append(fs,title)
                fs = append(fs,rqvalue)
                fs = append(fs,hyvalue)
                fs = append(fs,fwvalue)
                fs = append(fs,gwvalue)
                fs = append(fs,cateValue)
                fs = append(fs,subcateValue)
                fs = append(fs,regionvalue)
                fs = append(fs,contectvalue)
                fs = append(fs,addressvalue)

                output.WriteString(strings.Join(fs,",") + "\n")
        }
}

func fetchListPage() {
        /*
        output : /Users/gaolichuang/workspace/go/src/galaxy_walker/logs/page
        */
        firstPage := "http://bj.58.com/lipinxianhua/?key=%E7%BB%BF%E6%A4%8D%E7%A7%9F%E8%B5%81&cmcskey=%E7%BB%BF%E6%A4%8D%E7%A7%9F%E8%B5%81"
        pageprefix := "http://bj.58.com/lipinxianhua/pn"
        pagesubfix := "/?key=%E7%BB%BF%E6%A4%8D%E7%A7%9F%E8%B5%81&cmcskey=%E7%BB%BF%E6%A4%8D%E7%A7%9F%E8%B5%81"
        pagenum := kpagenum
        pagepath := listpagePath
        file.MkDirAll(pagepath)
        fetchUrl(firstPage, filepath.Join(pagepath, "1"))
        for i := 2; i <= pagenum; i++ {
                url := fmt.Sprintf("%s%d%s", pageprefix, i, pagesubfix)
                fetchUrl(url, filepath.Join(pagepath, strconv.Itoa(i)))
                fmt.Printf("%d fetch %s", i, url)
                time.Sleep(time.Second * time.Duration(rand.Int()/5))
        }
}
func parseAndFetchAllPages() {
        for i := 1; i <= kpagenum; i++ {
                filename := filepath.Join(listpagePath, strconv.Itoa(i))
                ids, pageid := parsePageId(filename)
                fmt.Println(pageid,ids)
                fetchContent(ids, pageid)
        }
}
func parsePageId(filename string) ([]string, string) {
        c, _ := file.ReadFileToString(filename)
        content := string(c)
        content = strings.Replace(content, "\n", "", -1)
        content = strings.Replace(content, "\t", "", -1)
        content = strings.Replace(content, " ", "", -1)
        content = strings.Replace(content, "&nbsp;", "", -1)
        idss := make([]string, 0)
        r, err := regexp.Compile("dataParam.*dispCateName")
        if err != nil {
                fmt.Println(err)
                return idss, filepath.Base(filename)
        }
        str := r.FindString(content)
        rr, _ := regexp.Compile(`,[\d]{14}_`)
        ids := rr.FindAllStringSubmatch(str, -1)
        idset := make(map[string]bool)
        for _, id := range ids {
                nid := strings.Replace(id[0], "_", "x", -1)
                nid=strings.TrimPrefix(nid,",")
                idset[nid] = true
        }

        for k, _ := range idset {
                idss = append(idss, k)
        }
        return idss, filepath.Base(filename)
}
func fetchContent(ids []string, pageid string) {
        path := filepath.Join(contentpagePath, pageid)
        file.MkDirAll(path)
        fmt.Println("XXXXXXXXXX", pageid)
        for _, id := range ids {
                url := fmt.Sprintf("http://bj.58.com/lipinxianhua/%s.shtml", id)
                fmt.Println(url)
                outputfilename := filepath.Join(path,id)
                err := fetchUrl(url,outputfilename)
                if err != nil {
                        LOG.Errorf("%s err %v",id,err)
                } else {
                        time.Sleep(time.Second)
                }
        }
}