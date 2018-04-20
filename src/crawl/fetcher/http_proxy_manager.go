package fetcher

import (
    "errors"
    "fmt"
    "math/rand"
    "net/url"
    "strconv"
    "galaxy_walker/internal/gcodebase/file"
    LOG "galaxy_walker/internal/gcodebase/log"
)

/*
Load proxy from file.
Selection:
1. round robin
2. random
3. score??
*/

type SelectMode int

const (
    PROXY_SELECT_RR     SelectMode = 1 << iota
    PROXY_SELECT_RANDOM
)

type ProxyManager struct {
    deads  []string // "host:port"
    alives []string
    mode   SelectMode
    index  int
    file.ConfigLoader
}

func (p *ProxyManager) loadConf() {
    result := p.LoadConfigWithTwoField("ProxyConf", *CONF.Crawler.ProxyConfFile, ":")
    p.deads = nil
    p.alives = nil
    for k, v := range result {
        port, err := strconv.Atoi(v)
        if err != nil {
            LOG.Errorf("Load Config Atoi Error, %s %s:%s", *CONF.Crawler.ProxyConfFile, k, v)
            continue
        }
        p.alives = append(p.alives, fmt.Sprintf("%s:%d", k, port))
        LOG.VLog(3).Debugf("Load fetch proxy %s : %d", k, port)
    }
}
func (p *ProxyManager) MarkDeadProxy(_url *url.URL) {
    alive := []string{}
    deadurl := _url.Host
    for _, c := range p.alives {
        if c != deadurl {
            alive = append(alive, c)
        } else {
            p.deads = append(p.deads, c)
        }
    }
    p.alives = alive
}
func (p *ProxyManager) GetProxyUrl() (*url.URL, error) {
    p.loadConf() // reload conf
    switch p.mode {
    case PROXY_SELECT_RR:
        return p.rrProxyUrl()
    case PROXY_SELECT_RANDOM:
        return p.randomProxyUrl()
    }
    return p.randomProxyUrl()
}
func (p *ProxyManager) randomProxyUrl() (*url.URL, error) {
    if len(p.alives) == 0 {
        return nil, errors.New("No Alive Proxy")
    }
    id := rand.Intn(len(p.alives))
    rawUrl := fmt.Sprintf("http://%s", p.alives[id])
    return url.Parse(rawUrl)
}
func (p *ProxyManager) rrProxyUrl() (*url.URL, error) {
    if len(p.alives) == 0 {
        return nil, errors.New("No Alive Proxy")
    }
    id := p.index % len(p.alives)
    p.index = (p.index + 1) % len(p.alives)
    rawUrl := fmt.Sprintf("http://%s", p.alives[id])
    return url.Parse(rawUrl)
}
func NewProxyManager(mode SelectMode) *ProxyManager {
    p := &ProxyManager{
        mode:  mode,
        index: 0,
    }
    p.loadConf()
    return p
}
