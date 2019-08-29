package service

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

type TimeLogObject struct {
	LoginTime  time.Time `bson:"login_time,omitempty"`
	CreateTime time.Time `bson:"create_time,omitempty"`
	UpdateTime time.Time `bson:"update_time,omitempty"`
}

func (t *TimeLogObject) Initialize() {
	t.CreateTime = time.Now().UTC()
	t.UpdateTime = t.CreateTime
	t.LoginTime = t.CreateTime
}

func (t *TimeLogObject) Login() {
	t.LoginTime = time.Now().UTC()
}

func (t *TimeLogObject) Update() {
	t.UpdateTime = time.Now().UTC()
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Fprintf(os.Stderr, "%s took %v\n", name, elapsed.Nanoseconds())
}
func GetFunctionName() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])

	tokens := strings.Split(f.Name(), ".")

	return tokens[len(tokens)-1]
}
