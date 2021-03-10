package main

import "time"

// message는 단일 메시지를 나타낸다.(JSON을 보냄)
// 메시지 문자열 자체를 캡슐화한다.
type message struct {
	Name    string
	Message string
	When    time.Time
}
