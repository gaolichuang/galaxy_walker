package test

import (
        "bytes"
        "errors"
        "fmt"
        "golang.org/x/net/html"
        "io"
        "strings"
)

func getBody(doc *html.Node) (*html.Node, error) {
        var b *html.Node
        var f func(*html.Node)
        f = func(n *html.Node) {
                fmt.Printf("Attr:%v,data:%s,atom:%s,ns:%s,type:%d\n",n.Attr,n.Data,n.DataAtom,n.Namespace,n.Type)
                fmt.Println("XXXXXXXXXXXXx",n.Data)
                fmt.Println(renderNode(n))
                fmt.Println("=============",n.Data)
                if n.Type == html.ElementNode && n.Data == "body" {
                        b = n
                }
                for c := n.FirstChild; c != nil; c = c.NextSibling {
                        f(c)
                }

        }
        f(doc)
        if b != nil {
                return b, nil
        }
        return nil, errors.New("Missing <body> in the node tree")
}

func renderNode(n *html.Node) string {
        var buf bytes.Buffer
        w := io.Writer(&buf)
        html.Render(w, n)
        return buf.String()
}

func main() {
        doc, _ := html.Parse(strings.NewReader(htm))
        bn, err := getBody(doc)
        if err != nil {
                return
        }
        body := renderNode(bn)
        fmt.Println(body)
}

const htm = `<!DOCTYPE html>
<html>
<head>
    <title color=1></title>
</head>
<body>
    body content
    <p>more content</p>
</body>
</html>`