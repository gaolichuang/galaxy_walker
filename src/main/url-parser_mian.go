package main

import (
        "net/url"
        "fmt"
        "net"
)

func main() {
        u,err := url.ParseRequestURI("http://www.baidu.com:99")
        //u,err := url.ParseRequestURI("http://www.baidu.com:99/aa/1.txt?a=1&b=2#1")
        if err != nil {
                fmt.Println(err)
                return
        }
        fmt.Println("ForceQuery:",u.ForceQuery)
        fmt.Println("Fragment:",u.Fragment)
        fmt.Println("Host:",u.Host)
        h,p,err := net.SplitHostPort(u.Host)
        fmt.Println(h,p,err)
        u.Host = "xx.xxxx"
        fmt.Println("Opaque:",u.Opaque)
        fmt.Println("Path",u.Path)
        fmt.Println("RawPath",u.RawPath)
        fmt.Println("RawQuery",u.RawQuery)
        fmt.Println("Scheme:",u.Scheme)
        fmt.Println("User:",u.User)
        fmt.Println("EscapedPath:",u.EscapedPath())
        fmt.Println("IsAbs:",u.IsAbs())
        fmt.Println("RequestURI:",u.RequestURI())
        fmt.Println("String:",u.String())
}
