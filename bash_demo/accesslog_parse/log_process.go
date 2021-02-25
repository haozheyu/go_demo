package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	client "github.com/influxdata/influxdb-client-go/v2"
)

type Reader interface {
	Read(rc chan []byte)
}
type Writer interface {
	Writer(wc chan *Message)
}
type Message struct {
	TimeLocal                    time.Time
	BytesSent                    int
	Path, Method, Scheme, Status string
	UpstreamTime, RequestTime    float64
}
type LogProcess struct {
	rc    chan []byte
	wc    chan *Message
	read  Reader
	write Writer
}
type ReadFromFile struct { // 读取文件的路径
	FilePath string
}
type WriteToInfluxDB struct { // influx data source
	InfluxDBsn string
	OrgName    string
	BucketName string
	Token      string
}
type SystemInfo struct {
	HandleLine   int     `json:"handle_line"`
	Tps          float64 `json:"tps"`
	ReadChanLen  int     `json:"read_chan_len"`
	WriteChanLen int     `json:"write_chan_len"`
	RunTime      string  `json:"run_time"`
	ErrNum       int     `json:"err_num"`
}

const (
	TypeHandleLine = 0
	TypeErrNum = 1
)
var (
	TypeMonitorChan = make(chan int ,200)
)
type Monitor struct {
	startTime time.Time
	data SystemInfo
	tpsSlice []int
}
func (r *ReadFromFile) Read(rc chan []byte) { // 读取模块
	//1 打开文件
	var (
		f   = &os.File{}
		err error
		rd  = &bufio.Reader{}
		rl  []byte
	)
	if f, err = os.Open(r.FilePath); err != nil {
		panic(fmt.Sprintf("open file fail:%s", err))
	}
	f.Seek(0, 2) // 文件的Unix-like移动到末尾

	//2 从末尾开始逐行读取
	rd = bufio.NewReader(f)
	for {
		if rl, err = rd.ReadBytes('\n'); err == io.EOF {



			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("file readBytes err:%s", err))
		}
		//3 将读取到的一行发给chan
		TypeMonitorChan <- TypeHandleLine
		rc <- rl[:len(rl)-1]
	}
}
func (w *WriteToInfluxDB) Writer(wc chan *Message) { // 写入模块
	for v := range wc {
		// 初始化influxdb client database(数据库) measurement(表) points(一行数据|tags索引属性|fields记录值|time时间戳)
		// 写入 curl -i -XPOST 'http://localhost:8086/write?db=mydb' --data-binary 'cpu_usage, host=server01,region=us-west value=0.64 1434055562000000000'
		// 读取 curl -i -XPOST 'http://localhost:8086/query?pretty=true' --data-urlencode "db=mydb" "q=SELECT \"value"\FROM\"cpu_usage\"WHERE\"region"\='us-west'"
		var (
			i_client     client.Client
			i_api        api.WriteAPIBlocking
			i_writePoint *write.Point
			err          error
		)
		i_client = client.NewClient(w.InfluxDBsn, w.Token) //创建一个client
		defer i_client.Close()
		i_api = i_client.WriteAPIBlocking(w.OrgName, w.BucketName) //使用阻塞写客户端对所需桶进行写操作
		i_writePoint = client.NewPoint("access_log",               //构造一行数据
			map[string]string{"Path": v.Path, "Method": v.Method, "Scheme": v.Scheme, "Status": v.Status},
			map[string]interface{}{"UpstreamTime": v.UpstreamTime,
				"RequestTime": v.RequestTime,
				"BytesSent":   v.BytesSent,
			},
			time.Now(),
		)
		if err = i_api.WritePoint(context.Background(), i_writePoint); err != nil {
			log.Println("influxdb writer fail:", err)
		} //数据写入
		/**
		i_writePoint = client.NewPointWithMeasurement("access_log").  //第二种写入方法
			AddTag("unit", "temperature").
			AddField("avg", 23.2).
			AddField("max", 45.0).
			SetTime(time.Now())
		i_api.WritePoint(context.Background(), i_writePoint)

		line := fmt.Sprintf("stat,unit=temperature avg=%f,max=%f", 23.5, 45.0)
		writeAPI.WriteRecord(context.Background(), line) //第三种写入方法


		**/
		fmt.Println(v)
	}
}
func (l *LogProcess) Process() { //解析模块
	/**
	172.0.0.12 - - [18/Feb/2021:16:37:37 +0800] http "GET /admin/WebEditor/.hg/requires HTTP/1.1" 200 2133 "-" "KeepAliveClient" "-" 1.005 1.854
	([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)
	**/
	var (
		read_byte      []byte
		reg            *regexp.Regexp
		reg_match      []string
		reg_str        = `([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`
		mes            = &Message{}
		t_location     = &time.Location{}
		t_timestamp    = time.Time{}
		err            error
		b_send         int
		b_method       []string
		b_url          *url.URL
		b_upstreamtime float64
		b_requestime   float64
	)
	reg = regexp.MustCompile(reg_str)
	t_location, _ = time.LoadLocation("Asia/Shanghai")
	for read_byte = range l.rc {
		reg_match = reg.FindStringSubmatch(string(read_byte))
		if len(reg_match) != 14 {
			log.Println("find log filed reg_match fail:", string(read_byte))
			TypeMonitorChan <- TypeErrNum
			continue
		}
		if t_timestamp, err = time.ParseInLocation("02/Jan/2006:15:04:05 -0700", reg_match[4], t_location); err != nil {
			TypeMonitorChan <- TypeErrNum
			log.Println("parse timestamp location fail: ", err.Error(), reg_match[4])
			continue
		}
		mes.TimeLocal = t_timestamp
		b_send, err = strconv.Atoi(reg_match[8])
		mes.BytesSent = b_send
		b_method = strings.Split(reg_match[6], " ")
		mes.Method = b_method[0]
		if b_url, err = url.Parse(b_method[1]); err != nil {
			TypeMonitorChan <- TypeErrNum
			log.Println("url parse file:", err)
			continue
		}
		mes.Path = b_url.Path
		mes.Scheme = reg_match[5]
		mes.Status = reg_match[7]
		if b_upstreamtime, err = strconv.ParseFloat(reg_match[12], 64); err != nil {
			TypeMonitorChan <- TypeErrNum
			log.Println("upstreamtime parse file:", err)
			continue
		}
		if b_requestime, err = strconv.ParseFloat(reg_match[13], 64); err != nil {
			TypeMonitorChan <- TypeErrNum
			log.Println("request time parse file:", err)
			continue
		}
		mes.UpstreamTime = b_upstreamtime
		mes.RequestTime = b_requestime

		l.wc <- mes
	}
}
func (m *Monitor) start(lp *LogProcess){ //应用监控处理
	var(
		req_data []byte
		err error
		tps_ticker *time.Ticker
	)
	go func() {
		for n:=range TypeMonitorChan{
			switch n {
			case TypeErrNum:
				m.data.ErrNum +=1
			case TypeHandleLine:
				m.data.HandleLine +=1
			}
		}
	}()
	tps_ticker = time.NewTicker(time.Second*5)
	go func() {
		for {
			<- tps_ticker.C
			m.tpsSlice = append(m.tpsSlice, m.data.HandleLine)
			if len(m.tpsSlice) > 2 {
				m.tpsSlice = m.tpsSlice[1:]
			}
		}
	}()

	http.HandleFunc("/monitor", func(writer http.ResponseWriter, request *http.Request) {
		m.data.RunTime = time.Now().Sub(m.startTime).String()
		m.data.ReadChanLen = len(lp.rc)
		m.data.WriteChanLen = len(lp.wc)
		if len(m.tpsSlice) >=2 {m.data.Tps = float64(m.tpsSlice[1]-m.tpsSlice[0])/5}
		if req_data,err = json.MarshalIndent(m.data,"","\t");err!=nil{fmt.Println("monitor start fail")}
	    io.WriteString(writer,string(req_data))
	})
	http.ListenAndServe(":9192",nil)
}
func main() {
	var (
		r = &ReadFromFile{
			FilePath: "C:/Users/hzy/Desktop/笔记记录/go_damo/bash_demo/accesslog_parse///access.log",
		}
		w = &WriteToInfluxDB{
			InfluxDBsn: "http://114.67.110.204:8086",
			OrgName:    "admin",
			BucketName: "my-bucket",
			Token:      "MRfX3p6TwjBTArxfGy9BPysdK8zm6bHwQcYg98-Cq599hPVnuXGJg9_5NdItPremWtjwq_dRycoQJ5_anQTAdA==",
		}
		m = &Monitor{
			startTime: time.Now(),
		    data: SystemInfo{},
		}
		rc = make(chan []byte,200)
		wc = make(chan *Message,200)
	)
	lp := &LogProcess{
		rc:    rc,
		wc:    wc,
		read:  r,
		write: w,
	}
	go lp.read.Read(lp.rc)

	for i :=0; i<5;i++{
		go lp.Process()
	}
	for i :=0; i<9;i++{
		go lp.write.Writer(lp.wc)
	}
	m.start(lp)
}
