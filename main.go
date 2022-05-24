package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/de4dwood/pushjob/job"
	"github.com/de4dwood/pushjob/push"
)

var (
	inputLabel flagArray
	PJS        push.PushJobStatus
	cmdfull    string
	outputlog  string
)

func init() {
	flag.Var(&inputLabel, "l", "label for metrics")
	flag.StringVar(&PJS.PushGatewayUrl, "h", "pushgateway.example.local", "pushgateway host")
	flag.StringVar(&PJS.JobName, "j", "pushGateway", "job name in pushgateway ")
	flag.StringVar(&cmdfull, "c", "echo nothing", "command to run ")
	flag.StringVar(&outputlog, "o", "/var/log/pushjob/log.log", "output log file")
	flag.Parse()
	for _, s := range inputLabel {
		output := strings.Split(s, "=")
		l := push.Label{Key: output[0], Value: output[1]}
		PJS.AddLabel(l)
	}

}

func main() {
	file, err := os.OpenFile(outputlog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()
	CheckErr(err)
	log.SetOutput(file)
	PJS.StartTime = time.Now().Local()
	log.Println("-----------------jobStart-----------------")
	output, statuCode, err := job.Command(cmdfull)
	CheckErr(err)
	log.Println(string(output))
	PJS.EndTime = time.Now().Local()
	PJS.StatusCode = statuCode
	if err = PJS.Push(); err != nil {
		log.Println(err)
	} else {
		log.Println(fmt.Sprintf("notification successfully sent for job:%s labels:%s", PJS.JobName, PJS.GetLabels()))
	}
}

func CheckErr(err error) {
	if err != nil {
		log.Println(err)

	}
}

type flagArray []string

func (i *flagArray) String() string {
	return "list of lablels"
}

func (a *flagArray) Set(label string) error {
	if label == "" {
		return fmt.Errorf("labesl cant be empty")
	}
	*a = append(*a, label)
	return nil
}
