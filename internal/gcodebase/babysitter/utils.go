package babysitter

import (
        "strings"
        "regexp"
)

func ParseStatusiHtml(html string) map[string]string {
        html = strings.Replace(html,"\n","",-1)
        pairs := make(map[string]string)
        r, _ := regexp.Compile("<key>.*?</key>:<value>.*?</value>")
        for _,v := range r.FindAllString(html,-1) {
                fs := strings.Split(v,"</key>:<value>")
                if len(fs) != 2 {
                        continue
                }
                fs[0] = strings.TrimLeft(fs[0],"<key>")
                fs[1] = strings.TrimRight(fs[1],"</value>")
                pairs[fs[0]] = fs[1]
        }
        // get healthy
        rh,_ :=regexp.Compile("<h1>Healthy</h1><div>.*?</div>")
        healthy := rh.FindString(html)
        healthy = strings.Replace(healthy,"<h1>Healthy</h1><div>","",-1)
        healthy = strings.TrimRight(healthy,"</div>")
        pairs["Healthy"] = healthy

        rk, _ := regexp.Compile("br>.*?:.*?<")
        for _,v := range rk.FindAllString(html,-1) {
                v = strings.TrimLeft(v,"br>")
                v = strings.TrimRight(v,"<")
                fs := strings.Split(v,":")
                if len(fs) != 2 {
                        continue
                }
                pairs[fs[0]] = fs[1]
        }
        return pairs
}
