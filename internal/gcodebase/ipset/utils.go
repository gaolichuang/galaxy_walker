package ipset

import (
	"strings"
	"net"
	"fmt"
	"errors"
        "sync"
)


func GetParentDomain(domain string) string {
        if domain == "." {
                return ""
        }
        idx := 0
        for idx < len(domain) && domain[idx] != '.' {
                idx += 1
        }
        // the last dot or not
        if idx < len(domain)-1 {
                idx += 1
        }
        return domain[idx:]
}
// 255.255.255.255 ==> 4294967295
func IP2Uint(ip net.IP) uint64 {
	// TODO. only ipv4 just....
	nip := ip.To4()
        if nip == nil {
                return 0
        }
	return uint64(nip[0])<<24 | uint64(nip[1])<<16 |
		uint64(nip[2])<<8 | uint64(nip[3])
}
func Uint2IP(ipnum uint64) net.IP {
        var p [net.IPv4len]byte
        p[3] = byte(ipnum & 0xff)
        ipnum = ipnum >> 8
        p[2] = byte(ipnum & 0xff)
        ipnum = ipnum >> 8
        p[1] = byte(ipnum & 0xff)
        ipnum = ipnum >> 8
        p[0] = byte(ipnum & 0xff)
        return net.IPv4(p[0], p[1], p[2], p[3])
}
// 1.1.1.1/8 ==> 00000001
func CIDR2BinarySeries(cidr string) (error, []byte) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return err, nil
	}
	ipi := IP2Uint(ip)
	mask, _ := ipnet.Mask.Size()
	binCIDR := [32]byte{0}
	idx := 31
	for ipi != 0 {
		binCIDR[idx] = byte(ipi & 1)
		ipi = ipi >> 1
		idx -= 1
	}
	return nil, binCIDR[:mask]
}
// 00000001 => 1.0.0.0/8
func BinarySeries2CIDR(cidr []byte) string {
	ipSer := [32]byte{0}
	for k, v := range cidr {
		ipSer[k] = v
	}
	ip := BinarySeries2IP(ipSer)
	return fmt.Sprintf("%s/%d", ip.String(), len(cidr))
}

// 1.2.3.4 = > 00000001 00000010 00000011 00000100
func IP2BinarySeries(ip net.IP) [32]byte {
	ipi := IP2Uint(ip)
	binIP := [32]byte{0}
	idx := 31
	for ipi != 0 {
		binIP[idx] = byte(ipi & 1)
		ipi = ipi >> 1
		idx -= 1
	}
	return binIP
}
func BinarySeries2IP(ipi [32]byte) net.IP {
	i := 0
	a, b, c, d := byte(0), byte(0), byte(0), byte(0)
	for {
		a += ipi[i]
		i += 1
		if i == 8 {
			break
		}
		a = a << 1
	}
	for {
		b += ipi[i]
		i += 1
		if i == 16 {
			break
		}
		b = b << 1
	}
	for {
		c += ipi[i]
		i += 1
		if i == 24 {
			break
		}
		c = c << 1
	}
	for {
		d += ipi[i]
		i += 1
		if i == 32 {
			break
		}
		d = d << 1
	}
	return net.IPv4(a, b, c, d)
}
func GenCIDRFromIP(ip1, ip2 net.IP) (cidrs []string, err error) {
	cidrs = make([]string, 0)
	ipBeg := IP2BinarySeries(ip1)
	ipEnd := IP2BinarySeries(ip2)
	larger, _ := LargerWithIpByte(ipBeg, ipEnd)
	if larger == 0 {
		cidrs = append(cidrs, fmt.Sprintf("%s/32", ip1.String()))
		return cidrs, nil
	}else if larger == 1 {
		return cidrs, errors.New("first argument is larger than second")
	}
	newIpBeg := ipBeg
	for newIpBeg != ipEnd {
		i := 31
		for ; i > -1; i-- {
			if newIpBeg[i] == 1 {
				break
			}
		}//i代表最开头的连续0的位数， 例如010010000 返回4
		finished := false
		for {
			j := i + 1
			tmp := newIpBeg
			for ; j < 32; j++ {
				tmp[j] = 1
			}
			larger, _ := LargerWithIpByte(tmp, ipEnd)
			if larger == -1 {
				break
			}else if larger == 1 {
				i += 1
				continue
			}else {
				finished = true
				break
			}
		}
		cidrs = append(cidrs, fmt.Sprintf("%s/%d", BinarySeries2IP(newIpBeg).String(), i+1))
		oldIpBeg := newIpBeg
		tmp := [32]byte{}
		tmp[i] = 1
		newIpBeg, err = IPByteSum(oldIpBeg, tmp)
		if err != nil {
			return cidrs, nil
		}
		if finished {
			break
		}
	}
	return cidrs, nil
}
func IPByteSum(ip1byte, ip2byte [32]byte) ([32]byte, error) {
	output := [32]byte{}
	carry := false
	for i := len(output) - 1; i > -1 ; i-- {
		placeValue := ip1byte[i] + ip2byte[i]
		if carry {
			placeValue += 1
		}
		if placeValue%2 == 0 {
			output[i] = 0
		}else{
			output[i] = 1
		}
		if placeValue > 1{
			carry = true
		}else {
			carry = false
		}
	}
	if carry == true {
		return output, errors.New(fmt.Sprint(output))
	}
	return output, nil
}
func LargerWithIpByte(ip1, ip2 [32]byte) (int, int) {
        for index, _ := range ip1 {
                if ip1[index] > ip2[index] {
                        return 1, index
                }else if ip1[index] < ip2[index] {
                        return -1, index
                }else {
                        continue
                }
        }
        return 0, -1
}
func GetLevelByViewName(viewname string) int {
        /*
                default
                china-_-_-ALL
                china-hebei-shijiazhuang-_
                china-rbejing-hebei-dianxin
        */
        strs := strings.Split(viewname,"-")
        if len(strs) == 1 {
                return 0
        }
        if strs[0] != "-" && strs[1] != "-" && strs[2] != "-" {
                return 4
        } else if strs[0] != "-" && strs[1] != "-" && strs[2] == "-" {
                if strings.HasPrefix(strs[1],"r") {
                        return 2
                } else {
                        return 3
                }
        } else if strs[0] != "-" && strs[1] == "-" && strs[2] == "-" {
                return 1
        }
        return 5
}
//////////////////////////////////
//////////////////////////////////
func IsReservedIP(addr net.IP) bool {
        start := []uint64{
                0, // 0.0.0.0
                167772160,//10.0.0.0
                1681915904,//100.64.0.0
                2130706432,//127.0.0.0
                2851995648,//169.254.0.0
                2886729728,//172.16.0.0
                3221225472,//192.0.0.0
                3221225984,//192.0.2.0
                3227017984,//192.88.99.0
                3232235520,//192.168.0.0
                3323068416,//198.18.0.0
                3325256704,//198.51.100.0
                3405803776,//203.0.113.0
                3758096384,//224.0.0.0
        }
        end := []uint64{
                16777215,//0.255.255.255
                184549375,//10.255.255.255
                1686110207,//100.127.255.255
                2147483647,//127.255.255.255
                2852061183,//169.254.255.255
                2887778303,//172.31.255.255
                3221225727,//192.0.0.255
                3221226239,//192.0.2.255
                3227018239,//192.88.99.255
                3232301055,//192.168.255.255
                3323199487,//198.19.255.255
                3325256959,//198.51.100.255
                3405804031,//203.0.113.255
                4294967295,//255.255.255.255
        }
        addr = addr.To4()
        if addr == nil {
                return false
        }
        firstByte := []byte(addr)[0]
        if firstByte == 0 || firstByte == 10 || firstByte == 100 || firstByte == 127 || firstByte == 169 || firstByte == 172 || firstByte == 192 || firstByte == 198 || firstByte == 203 || firstByte >= 224 {
                err, ok := judgeIPByNum(addr, start, end)
                if err != nil {
                        return true
                }
                return ok
        }
        return false
}
func IsPrivateIPByScope(addr net.IP) bool {
        cidrs := []string {
                "0.0.0.0/8",
                "10.0.0.0/8",
                "100.64.0.0/10",
                "127.0.0.0/8",
                "169.254.0.0/16",
                "172.16.0.0/12",
                "192.0.0.0/24",
                "192.0.2.0/24",
                "192.88.99.0/24",
                "192.168.0.0/16",
                "198.18.0.0/15",
                "198.51.100.0/24",
                "203.0.113.0/24",
                "224.0.0.0/4",
                "240.0.0.0/4",
                "255.255.255.255/32",
        }
        scope := make([]*net.IPNet,0)
        for _,cidr := range cidrs {
                _,ipn,_ := net.ParseCIDR(cidr)
                scope = append(scope,ipn)
        }

        err,ok := judgeIPInScope(addr,scope)
        if err != nil {
                return true
        }
        return ok
}
func judgeIPInScope(addr net.IP, scope []*net.IPNet) (error,bool) {
        if addr == nil {
                return errors.New("No Ip"),false
        }
        for _,ipnet := range scope {
                if ipnet.Contains(addr) {
                        return nil,true
                }
        }
        return nil,false
}

func IsPrivateIP(addr net.IP) bool {
        start := []uint64{
                167772160,//10.0.0.0
                1681915904,//100.64.0.0
                2130706432,//127.0.0.0
                2886729728,//172.16.0.0
                3221225472,//192.0.0.0
                3232235520,//192.168.0.0
                3323068416,//198.18.0.0
        }
        end := []uint64{
                184549375,//10.255.255.255
                1686110207,//100.127.255.255
                2147483647,//127.255.255.255
                2887778303,//172.31.255.255
                3221225727,//192.0.0.255
                3232301055,//192.168.255.255
                3323199487,//198.19.255.255
        }
        addr = addr.To4()
        if addr == nil {
                return false
        }

        firstByte := []byte(addr)[0]
        if firstByte == 10 || firstByte == 100 || firstByte == 127 || firstByte == 172 || firstByte == 192 || firstByte == 198 {
                err, ok := judgeIPByNum(addr, start, end)
                if err != nil {
                        return true
                }
                return ok
        }
        return false
}
func judgeIPByNum(addr net.IP, start,end []uint64) (error,bool) {
        if len(start) != len(end) {
                return errors.New("Error Scope"),false
        }
        addrNum := IP2Uint(addr)
        for i := 0;i < len(start); i++ {
                if addrNum >= start[i] && addrNum <= end[i] {
                        return nil,true
                }
        }
        return nil,false
}
func ParseViewName(name string) (string,string,string,string) {
        fs := strings.Split(name,"-")
        if len(fs) != 4 {
                return "","","",""
        }
        return fs[0],fs[1],fs[2],fs[3]
}
const (
        KViewUnKnown      = 0
        KViewTypeCountry  = 1
        KViewTypeRegion   = 2
        KViewTypeProvince = 3
        KViewTypeCity     = 4
        KViewTypeOverSea  = 5
        KViewTypeCustom  = 6
        KViewTypeCustomParent  = 7
)
func ParseViewType(name string) int {
        g,p,c,i := ParseViewName(name)
        if strings.HasPrefix(g,"custom") {
                if i == "all" {
                        return KViewTypeCustomParent
                }
                return KViewTypeCustom
        }
        if g != "china" && g != "" {
                return KViewTypeOverSea
        }
        if p == "_" {
                if c == "_" {
                        return KViewTypeCountry
                }
        } else if c == "_" {
                if p != "rall" {
                        if strings.HasPrefix(p, "r") {
                                return KViewTypeRegion
                        } else {
                                return KViewTypeProvince
                        }
                }
        } else {
                return KViewTypeCity
        }
        return KViewUnKnown
}
////////////////////////////////////
////////////////////////////////////
/*
func ExternalIP() (string, string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), iface.HardwareAddr.String(), nil
		}
	}
	return "", "", errors.New("are you connected to the network?")
}
*/
func getLocalIPInterval(name string) net.IP {
        // reference http://studygolang.com/articles/01202
        ifaces, err := net.Interfaces()
        if err != nil {
                return nil
        }
        for _, iface := range ifaces {
                if iface.Flags&net.FlagUp == 0 {
                        continue // interface down
                }
                if iface.Flags&net.FlagLoopback != 0 {
                        continue // loopback interface
                }
                addrs,err := iface.Addrs()
                if err != nil {
                        continue
                }
                if name != "" && iface.Name != name {
                        continue
                }
                for _,addr := range addrs {
                        // check the address type and if it is not a loopback the display it
                        if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
                                if ipnet.IP.To4() != nil {
                                        return ipnet.IP
                                }
                        }
                }
        }
        return nil
}

func GetAllLocalIP() []net.IP {
        // reference http://studygolang.com/articles/01202
        ifaces, err := net.Interfaces()
        if err != nil {
                return nil
        }
        ips := make([]net.IP,0)
        reservedIPs := make([]net.IP,0)
        for _, iface := range ifaces {
                if iface.Flags&net.FlagUp == 0 {
                        continue // interface down
                }
                /*
                if iface.Flags&net.FlagLoopback != 0 {
                        continue // loopback interface
                }
                */
                addrs,err := iface.Addrs()
                if err != nil {
                        continue
                }
                for _,addr := range addrs {
                        // check the address type and if it is not a loopback the display it
                        if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
                                addr := ipnet.IP.To4()
                                if addr != nil {
                                        if IsReservedIP(addr) {
                                                reservedIPs = append(reservedIPs,addr)
                                        } else {
                                                ips = append(ips,ipnet.IP)
                                        }
                                }
                        }
                }
        }
        /*
        如果内外网ip都存在,则只返回外网ip
        但如果都是内网ip,则返回内网ip回去
        */
        if len(ips) <= 0 {
                return reservedIPs
        }
        return ips
}
func ParserViewName(name string) (string,string,string,string) {
        fs := strings.Split(name,"-")
        if len(fs) != 4 {
                return "","","",""
        }
        return fs[0],fs[1],fs[2],fs[3]
}

///////////// Singleton  ////////////////////////
var _localIP_instance net.IP = nil
var _localIP_init_ctx sync.Once

func GetLocalIP(nic string) net.IP {
        _localIP_init_ctx.Do(func() {
                _localIP_instance = getLocalIPInterval(nic)
        })
        return _localIP_instance
}
var _localIPStr_instance string
var _localIPStr_init_ctx sync.Once

func GetLocalIPStr(nic string) string {
        _localIPStr_init_ctx.Do(func() {
                _localIPStr_instance = getLocalIPInterval(nic).String()
        })
        return _localIPStr_instance
}