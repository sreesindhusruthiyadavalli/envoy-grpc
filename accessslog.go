package main

import (
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	ald "github.com/envoyproxy/go-control-plane/envoy/data/accesslog/v2"
	als "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v2"
)


type AccessLogService struct {
	entries []string
	mu      sync.Mutex
}

func (svc *AccessLogService) log(entry string) {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	svc.entries = append(svc.entries, entry)
	log.Printf("AccessLog:  " + entry)

}

// Dump releases the collected log entries and clears the log entry list.
func (svc *AccessLogService) Dump(f func(string)) {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	for _, entry := range svc.entries {
		f(entry)
	}
	svc.entries = nil
}

// StreamAccessLogs implements the access log service.
func (svc *AccessLogService) StreamAccessLogs(stream als.AccessLogService_StreamAccessLogsServer) error {
	var logName string
	log.Printf("In StreamAccesslogs method")
	for {
		msg, err := stream.Recv()
		if err != nil {
			continue
		}
		if msg.Identifier != nil {
			logName = msg.Identifier.LogName
		}
		switch entries := msg.LogEntries.(type) {
		case *als.StreamAccessLogsMessage_HttpLogs:
			log.Printf("http logs")
			for _, entry := range entries.HttpLogs.LogEntry {
				if entry != nil {
					common := entry.CommonProperties
					req := entry.Request
					resp := entry.Response
					if common == nil {
						common = &ald.AccessLogCommon{}
					}
					if req == nil {
						req = &ald.HTTPRequestProperties{}
					}
					if resp == nil {
						resp = &ald.HTTPResponseProperties{}
					}
					svc.log(fmt.Sprintf("[%s%s] common: %+v, requestProperties:%+v,responseProperties:%+v",
						logName, time.Now().Format(time.RFC3339), common, req, resp))
				}
			}
		case *als.StreamAccessLogsMessage_TcpLogs:
			log.Printf("tcp logs  access message")
			for _, entry := range entries.TcpLogs.LogEntry {
				if entry != nil {
					common := entry.CommonProperties
					if common == nil {
						common = &ald.AccessLogCommon{}
					}
					svc.log(fmt.Sprintf("[%s%s] tcp common:%+v",
						logName, time.Now().Format(time.RFC3339), common))
				}
			}
		}
	}
}
