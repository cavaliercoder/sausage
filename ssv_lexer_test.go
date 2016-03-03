package main

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

var ssvLog = `
      07/Aug/2015 00:02:18  60334 10.241.144.12 TCP_MISS/200 2533 CONNECT webmail.det.wa.edu.au:443 - HIER_DIRECT/10.1.10.169 - "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; WOW64; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C)"
      07/Aug/2015 00:02:18      1 10.241.144.12 TCP_DENIED/302 359 CONNECT urs.microsoft.com:443 - HIER_NONE/- text/html "VCSoapClient"
      07/Aug/2015 00:02:18      0 10.241.144.12 TCP_DENIED/302 359 CONNECT urs.microsoft.com:443 - HIER_NONE/- text/html "VCSoapClient"
      07/Aug/2015 00:02:18     22 10.85.8.11 TCP_MISS/200 226 GET http://gsp1.apple.com/pep/gcc - TIMEOUT_HIER_DIRECT/167.30.50.135 text/html "GeoServices/982.64 CFNetwork/711.4.6 Darwin/14.0.0"
      07/Aug/2015 00:02:18      0 10.241.144.12 TCP_DENIED/302 359 CONNECT urs.microsoft.com:443 - HIER_NONE/- text/html "VCSoapClient"
      07/Aug/2015 00:02:18    427 10.155.40.254 TCP_MISS/200 19118 CONNECT login.det.wa.edu.au:443 - HIER_DIRECT/10.1.143.126 - "Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; WOW64; Trident/6.0)"
      07/Aug/2015 00:02:18      0 10.241.144.12 TCP_DENIED/302 359 CONNECT urs.microsoft.com:443 - HIER_NONE/- text/html "VCSoapClient"
      07/Aug/2015 00:02:18      0 10.241.144.12 TCP_DENIED/302 359 CONNECT urs.microsoft.com:443 - HIER_NONE/- text/html "VCSoapClient"
      07/Aug/2015 00:02:18      0 10.241.144.12 TCP_DENIED/302 359 CONNECT urs.microsoft.com:443 - HIER_NONE/- text/html "VCSoapClient"
      07/Aug/2015 00:02:18      0 10.241.144.12 TCP_DENIED/302 359 CONNECT urs.microsoft.com:443 - HIER_NONE/- text/html "VCSoapClient"
      07/Aug/2015 00:02:18      0 10.241.144.12 TCP_DENIED/302 359 CONNECT urs.microsoft.com:443 - HIER_NONE/- text/html "VCSoapClient"
      07/Aug/2015 00:02:18    843 10.209.241.62 TCP_MISS/200 423 GET http://v4.moatads.com/pixel.gif?e=17&i=MNINE1&bq=0&f=0&j=http%3A%2F%2Fwww.ninemsn.com.au&o=3&t=1438876935436&de=799181944690&m=0&ar=3853b57-clean&q=3&cb=0&cu=1438876935425&ll=16&ln=0&r=18.0.0&em=0&en=0&d=10812498%3A300x250%3A2000000000196132%3Aundefined&qs=5&bo=2000000000196132&bd=undefined&ac=1&it=500&cs=0 - HIER_DIRECT/54.152.243.169 image/gif "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.125 Safari/537.36"
      07/Aug/2015 00:02:18  30289 10.202.153.148 TCP_MISS/200 4872 CONNECT api.smoot.apple.com:443 - HIER_DIRECT/17.252.249.246 - "Parsec/1 (iPad2,4; iPhone OS 12D508) Spotlight/1.0"
      07/Aug/2015 00:02:18   5318 10.235.56.11 TCP_MISS/200 19299 CONNECT login.det.wa.edu.au:443 - HIER_DIRECT/10.1.143.126 - "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.125 Safari/537.36"
      07/Aug/2015 00:02:18    462 10.209.81.62 TCP_MISS/200 0 CONNECT push.collobos.com:443 - HIER_DIRECT/54.218.39.91 - "-"
      07/Aug/2015 00:02:18    494 10.100.72.57 TCP_MISS/200 4704 CONNECT configuration.apple.com:443 - HIER_DIRECT/23.63.48.176 - "Mail/53 CFNetwork/711.3.18 Darwin/14.0.0"
      07/Aug/2015 00:02:18 38892492 10.219.129.62 TCP_MISS/200 2550695 CONNECT s8.mylivechat.com:443 - HIER_DIRECT/74.86.208.242 - "Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; WOW64; Trident/6.0)"
      07/Aug/2015 00:02:18  25464 10.219.129.62 TCP_MISS/200 5331 CONNECT safebrowsing-cache.google.com:443 - HIER_DIRECT/218.100.43.222 - "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:39.0) Gecko/20100101 Firefox/39.0"
      07/Aug/2015 00:02:18  25845 10.219.129.62 TCP_MISS/200 1003 CONNECT safebrowsing.google.com:443 - HIER_DIRECT/218.100.43.249 - "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:39.0) Gecko/20100101 Firefox/39.0"
      07/Aug/2015 00:02:18 848458 10.219.129.62 TCP_MISS/200 20996 CONNECT feedbackws.icloud.com:443 - HIER_DIRECT/17.151.239.54 - "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.130 Safari/537.36"
      07/Aug/2015 00:02:18 861698 10.219.129.62 TCP_MISS/200 120050 CONNECT p10-calendarws.icloud.com:443 - HIER_DIRECT/17.151.224.15 - "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.130 Safari/537.36"
      07/Aug/2015 00:02:18  50617 10.219.129.62 TCP_MISS/200 6326 CONNECT classdojo.pubnub.com:443 - HIER_DIRECT/54.249.82.171 - "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.125 Safari/537.36"
      07/Aug/2015 00:02:18  31284 10.219.129.62 TCP_MISS/200 5581 CONNECT classdojo.pubnub.com:443 - HIER_DIRECT/54.249.82.171 - "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.130 Safari/537.36"
      07/Aug/2015 00:02:18  46061 10.219.129.62 TCP_MISS/200 6326 CONNECT classdojo.pubnub.com:443 - HIER_DIRECT/54.249.82.171 - "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.125 Safari/537.36"
      07/Aug/2015 00:02:19      0 10.1.76.250 TCP_MEM_HIT/200 375 HEAD http://www.microsoft.com/ - HIER_NONE/- text/html "-"
      07/Aug/2015 00:02:19      0 10.1.76.250 TCP_MEM_HIT/200 375 HEAD http://www.microsoft.com/ - HIER_NONE/- text/html "-"
      07/Aug/2015 00:02:19     64 10.1.76.250 TCP_MISS/302 665 GET http://www.google.com/ - TIMEOUT_HIER_DIRECT/218.100.43.241 text/html "-"
      07/Aug/2015 00:02:19     69 10.1.76.250 TCP_MISS/302 665 GET http://www.google.com/ - TIMEOUT_HIER_DIRECT/218.100.43.230 text/html "-"
      07/Aug/2015 00:02:19     68 10.1.76.251 TCP_MISS/302 665 GET http://www.google.com/ - TIMEOUT_HIER_DIRECT/218.100.43.245 text/html "-"
      07/Aug/2015 00:02:19     67 10.1.76.250 TCP_MISS/302 665 GET http://www.google.com/ - TIMEOUT_HIER_DIRECT/218.100.43.237 text/html "-"
      07/Aug/2015 00:02:19     16 10.209.241.62 TCP_MISS/204 353 GET http://b.scorecardresearch.com/b?c1=7&c2=8973917&c3=1&ns__t=1438876938568&ns_c=UTF-8&c8=ninemsn%20Homepage%20-%20News%2C%20Sport%2C%20Finance%2C%20Lifestyle%2C%20TV%2C%20Competitions%2C%20Horoscopes%2C%20Daily%20Quiz&c7=http%3A%2F%2Fwww.ninemsn.com.au%2F%3Frf%3Dtrue&c9=http%3A%2F%2Fwww.ninemsn.com.au%2F%3Frf%3Dtrue - TIMEOUT_HIER_DIRECT/167.30.50.134 - "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.125 Safari/537.36"
      07/Aug/2015 00:02:19    276 10.155.40.254 TCP_MISS/200 19142 CONNECT login.det.wa.edu.au:443 - HIER_DIRECT/10.1.143.126 - "Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; WOW64; Trident/6.0)"
      07/Aug/2015 00:02:19    301 10.209.241.62 TCP_MISS/200 423 GET http://v4.moatads.com/pixel.gif?e=34&ud=0&up=0&qa=1280&qb=1024&qc=0&qd=0&qf=1280&qe=899&qh=1280&qg=984&qi=1280&qj=984&qk=0&ql=%3B%5BpwxnRd%7Dt%3Aa%5DmJVOG)%2C~%405%2F%5BGI%3Fi6%5EB61%2F%3DSqcMr1%7B%2CW5.NO)(aWFGuW3%23(kQ%5DQ%3A__vQapTyKIw%40%40soz4%5EC%2CYRd%7Cw%24_%3Fj!LNSiyMm4MQ2%24(%5DmGUpR%5B.%25B%7BGs11_2CTpj%2F%24b%3Dh_G3%253(N%5BvPUDby7p1vU3ZtuCG&qm=-480&qn=6OZw%3DoHB%2CEF%3FKC1I%3Cq.bWoCSV2W0Su*TDXlCfX2iR2%25(GyHN%3DI(%2C%3Ba15lK1t!9ZpAH..4iwM%25z4mc4%7Di3MTg%26B%3BLm!__PyDN(%2BWx*h~%3F03*%5B)%2C2iVSWfV%7D%2F%2FRA7R.eJKx%7Ci6sGm!ryh%7Cek)3.%5BqC%7Dq%40Dgh%2C%7B%5BH%3BRy%5EQ%5E%5BhPSI.%24ki)sV~1HmDkx2EF6pJBPJ.(0E%3AUdBE)ea*X%3Dy%3E%5B%25B7k.%3ETy%25.8e%40GW*_)9L%2CzVx)rOS2z.%5BOCDTWRe%2Ba%2Fke%3BR30982iYBgDzb%23Ls1(u0EnUa%3Fwb%26k!C%24%26J%3BBcJVrwLy%3Aaq%24St%3Fxny%3Am%5EGbv5*7*7UO0%40M%7CQDt%3ExZq%224%7CQjw%60.%7Bi%3F%5DQZ%2CA2%2BNhloI%40s1%7CZ5*%3FVl%3Fe3%7CqL5%40J%3D%5BwPjrG%3D2fb%2CM%249!0t9.aS%3B4oD%7D%60%3Fjc!L2LmqMs%3Cex1bxNTK7%2BuCTpY%3CZ.e&qo=0&qp=10000&qq=000001100000&qr=0&qt=0&i=MNINE1&bq=0&f=0&j=http%3A%2F%2Fwww.ninemsn.com.au&o=3&t=1438876935449&de=370793092530&m=0&ar=3853b57-clean&q=4&cb=0&cu=1438876935440&ll=16&ln=0&r=18.0.0&em=0&en=0&d=11556764%3A1x1%3A80000000000018486%3Aundefined&qs=5&bo=80000000000018486&bd=undefined&ac=1&it=500&cs=0 - TIMEOUT_HIER_DIRECT/54.152.243.169 image/gif "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.125 Safari/537.36"
      07/Aug/2015 00:02:19  29694 10.217.161.62 TCP_MISS/200 5351 CONNECT safebrowsing-cache.google.com:443 - HIER_DIRECT/218.100.43.207 - "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:39.0) Gecko/20100101 Firefox/39.0"
      07/Aug/2015 00:02:19 3628593 10.217.161.62 TCP_MISS/200 14432 CONNECT 1.client-channel.google.com:443 - HIER_DIRECT/74.125.204.189 - "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:39.0) Gecko/20100101 Firefox/39.0"
      07/Aug/2015 00:02:19   2819 10.197.225.62 TCP_MISS/200 375 POST http://meeting03.prezi.com/ - HIER_DIRECT/54.237.62.179 text/plain "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:39.0) Gecko/20100101 Firefox/39.0"
      07/Aug/2015 00:02:19    566 10.197.225.62 TCP_MISS/200 287 POST http://meeting03.prezi.com/ - HIER_DIRECT/54.237.62.179 text/plain "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:39.0) Gecko/20100101 Firefox/39.0"
      07/Aug/2015 00:02:19    718 10.235.88.169 TCP_MISS/204 338 POST http://logs.spilgames.com/lg/pb/1/ut/ - HIER_DIRECT/212.72.60.214 - "Mozilla/5.0 (iPad; CPU OS 8_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12F69 Safari/600.1.4"
      07/Aug/2015 00:02:19  30137 10.217.161.62 TCP_MISS/200 1003 CONNECT safebrowsing.google.com:443 - HIER_DIRECT/218.100.43.236 - "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:39.0) Gecko/20100101 Firefox/39.0"
      07/Aug/2015 00:02:19    466 10.195.121.62 TCP_MISS/200 0 CONNECT push.collobos.com:443 - HIER_DIRECT/54.218.39.91 - "-"
      07/Aug/2015 00:02:19  57728 10.243.33.62 TCP_MISS/200 438 CONNECT clients4.google.com:443 - HIER_DIRECT/218.100.43.241 - "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.125 Safari/537.36"
      07/Aug/2015 00:02:19 146672 10.243.33.62 TCP_MISS/200 5768 CONNECT classdojo.pubnub.com:443 - HIER_DIRECT/54.249.82.171 - "Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; WOW64; Trident/6.0)"
      07/Aug/2015 00:02:19 1043595 10.243.33.62 TCP_MISS/200 4140 CONNECT 5.client-channel.google.com:443 - HIER_DIRECT/74.125.203.189 - "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.125 Safari/537.36"
      07/Aug/2015 00:02:19  36370 10.243.33.62 TCP_MISS_ABORTED/200 631 GET http://realtime.services.disqus.com/api/2/thread/1884102185?bust=460 - HIER_DIRECT/108.168.151.6 application/json "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.125 Safari/537.36"
      07/Aug/2015 00:02:19  69238 10.243.33.62 TCP_MISS/200 1462 CONNECT classdojo.zendesk.com:443 - HIER_DIRECT/192.161.147.1 - "Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; WOW64; Trident/6.0)"
      07/Aug/2015 00:02:19  68212 10.243.33.62 TCP_MISS/200 5768 CONNECT ps9.pubnub.com:443 - HIER_DIRECT/54.249.82.174 - "Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; WOW64; Trident/6.0)"
      07/Aug/2015 00:02:19  68193 10.243.33.62 TCP_MISS/200 5203 CONNECT dialog.filepicker.io:443 - HIER_DIRECT/54.241.21.84 - "Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; WOW64; Trident/6.0)"
      07/Aug/2015 00:02:19  68196 10.243.33.62 TCP_MISS/200 6141 CONNECT classdojo.pubnub.com:443 - HIER_DIRECT/54.249.82.175 - "Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; WOW64; Trident/6.0)"
      07/Aug/2015 00:02:19  68197 10.243.33.62 TCP_MISS/200 5203 CONNECT www.filepicker.io:443 - HIER_DIRECT/54.241.21.84 - "Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; WOW64; Trident/6.0)"
      07/Aug/2015 00:02:19     99 10.127.9.62 TCP_MISS/200 231 GET http://gsp1.apple.com/pep/gcc - TIMEOUT_HIER_DIRECT/167.30.50.135 text/html "GeoServices/982.64 CFNetwork/711.4.6 Darwin/14.0.0"
`

func TestSSVLexer(t *testing.T) {
	fields := []string{"", "", "size", "", "", "duration", "method", "url", "", "", "mime_type", "agent"}
	r := bufio.NewReader(strings.NewReader(ssvLog))
	l := NewSSVLexer(fields)

	for {
		// read line
		b, _, err := r.ReadLine()
		if err != nil {
			if err != io.EOF {
				t.Errorf(err.Error())
			}

			break
		}

		if len(b) > 0 {
			// lex it
			if m, err := l.Lex(string(b)); err == nil {
				// ensure each field is present
				for _, field := range fields {
					if field != "" {
						if _, ok := m[field]; !ok {
							t.Errorf("field '%s' is missing from lexer output", field)
						}
					}
				}
			} else {
				t.Errorf("%v", err)
			}
		}
	}
}
