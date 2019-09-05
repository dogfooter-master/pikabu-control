package main

import (
	"bufio"
	"bytes"
	"pikabu-control/control/pkg/endpoint"
	"pikabu-control/control/pkg/service"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/go-kit/kit/log"
	"gopkg.in/mgo.v2/bson"
	"io"
	"io/ioutil"
	stdLog "log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var fs = flag.NewFlagSet("http-client", flag.ExitOnError)
var hostsAddr = fs.String("hosts", "", "server address for testing")

var scheme string
var tokens []string
var accessToken string
var services []string
var ids map[int][]string
var logger log.Logger
var pipeData Response

func main() {
	fs.Parse(os.Args[1:])
	if err := service.LoadConfig(); err != nil {
		panic(err)
	}

	//cfg.HttpHosts = "pikabu.io:7977"
	if len(*hostsAddr) == 0 {
		scheme = "http://" + service.GetConfigClientControlHttp() + "/api"
	} else {
		scheme = *hostsAddr + "/api"
	}
	fmt.Fprintf(os.Stdout, "%v\n", scheme)

	services = append(services, "SignUp")
	services = append(services, "SignUpPassword")
	services = append(services, "SignUpComplete")
	services = append(services, "VerifyCertificationCode")
	services = append(services, "ResendCertificationCode")
	services = append(services, "UpdatePassword")
	services = append(services, "SignIn")
	services = append(services, "UpdateUserInformation")
	services = append(services, "UpdateAccessToken")
	services = append(services, "PrepareAvatar")
	services = append(services, "GetAvatarUri")
	services = append(services, "GetAgentList")

	services = append(services, "MakeTestCase")

	showUsage()

	lastServiceName := "0"
	for {
		var serviceName string
		fmt.Fprintf(os.Stdout, "%v > ", accessToken)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		line := scanner.Text()
		line = strings.TrimLeft(line, " ")
		if len(line) == 0 {
			serviceName = lastServiceName
		} else {
			tokens = strings.Split(line, " ")
			serviceName = tokens[0]
			lastServiceName = serviceName
		}

		//if n, err := fmt.Scanln(&args); err == nil {
		//	fmt.Fprintf(os.Stderr, "DEBUG:%v %v\n", n, args)
		//	tokens = strings.Split(args, " ")
		//	serviceName, _ = strconv.Atoi(tokens[0])
		//	lastServiceName = serviceName
		//} else {
		//	fmt.Fprintf(os.Stderr, "DEBUG ERR:%v\n", err)
		//	serviceName = lastServiceName
		//}
		serviceName = lastServiceName
		fmt.Fprintf(os.Stdout, "\n")

		switch serviceName {
		case "GetPlatform":
			GetPlatform()
		case "GetSystemType":
			GetSystemType()
		case "GetUserStatus":
			GetUserStatus()
		case "SignUp":
			SignUp()
		case "SignUpPassword":
			SignUpPassword()
		case "VerifyCertificationCode":
			VerifyCertificationCode()
		case "ResendCertificationCode":
			ResendCertificationCode()
		case "UpdatePassword":
			UpdatePassword()
		case "SignIn":
			SignIn()
		case "UpdateUserInformation":
			UpdateUserInformation()
		case "UpdateAccessToken":
			UpdateAccessToken()
		case "PrepareAvatar":
			PrepareAvatar()
		case "GetAvatarUri":
			GetAvatarUri()
		case "GetAgentList":
			GetAgentList()

		case "MakeTestCase":
			MakeTestCase()

		case "h":
			showUsage()
		}

		for _, e := range tokens {
			fmt.Fprintf(os.Stderr, "%v ", e)
		}
		fmt.Fprintf(os.Stderr, "\n")
	}
}


func GetAgentList () {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s:<service> <access_token>\n", GetFunctionName())
		return
	}

	message := service.Payload{
		Category: "private",
		Service: "GetAgentList",
		AccessToken: tokens[1],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
	return

	return
}

func GetImages () {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s:<service> <access_token> <image_id>...\n", GetFunctionName())
		return
	}

	message := service.Payload{
		Category: "image",
		Service: "GetImages",
		AccessToken: tokens[1],
	}

	var imageIdList []service.ImageObject
	for i := 2; i < len(tokens); i++ {
		imageIdList = append(imageIdList,
			service.ImageObject{
				ImageId: tokens[i],
			})
	}
	message.ImageList = imageIdList

	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
	return

	return
}

func DeleteImagesInLibrary () {
	if len(tokens) < 4 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <hospital_id> <image_id>...\n", GetFunctionName())
		return
	}

	message := service.Payload{
		Category: "image",
		Service: "DeleteImagesInLibrary",
		AccessToken: tokens[1],
		HospitalId: tokens[2],
		Status: "deleted",
	}

	for i := 3; i < len(tokens); i++ {
		message.ImageList = append(message.ImageList,
			service.ImageObject{
				ImageId: tokens[i],
			})
	}

	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
	return
}

func DeleteDate() {
	if len(tokens) < 4 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <diagnosis_id> <date_id>\n", GetFunctionName())
		return
	}

	message := service.Payload{
		Category: "date",
		Service: "DeleteDate",
		AccessToken: tokens[1],
	}

	message.Diagnosis = &service.DiagnosisObject{
		DiagnosisId: tokens[2],
	}
	message.Date = &service.DateObject{
		DateId: tokens[3],
	}

	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
	return
}

func UpdateDiagnosis() {
	if len(tokens) < 5 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <patient_id> <diagnosis_id> <date_id> <type> <visit_date> <note> <treatment> <tag_list>...\n", GetFunctionName())
		return
	}

	message := service.Payload{
		Category:    "diagnosis",
		Service:     "UpdateDiagnosis",
		AccessToken: tokens[1],
	}
	message.Patient = &service.PatientObject{
		PatientId: tokens[2],
	}
	message.Diagnosis = &service.DiagnosisObject{
		DiagnosisId: tokens[3],
	}
	message.Date = &service.DateObject{
		DateId: tokens[4],
	}

	index := 5
	if len(tokens) > index && tokens[index] != "0" {
		message.Diagnosis.Type = tokens[index]
		index++
	}

	if len(tokens) > index && tokens[index] != "0" {
		message.Date.VisitDate = tokens[index]
		index++
	}

	if len(tokens) > index && tokens[index] != "0" {
		message.Date.Note = tokens[index]
		index++
	}

	if len(tokens) > index && tokens[index] != "0" {
		message.Date.Treatment = tokens[index]
		index++
	}

	var tags []string
	for i := index; i < len(tokens); i++ {
		tags = append(tags, tokens[index])
	}
	message.Date.TagList = tags
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
	return
}

func UpdateImage3D() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <image_id>\n", GetFunctionName())
		return
	}

	message := service.Payload{
		Category:    "image",
		Service:     "UpdateImage3D",
		AccessToken: tokens[1],
	}

	message.Image = &service.ImageObject{
		ImageId: tokens[2],
	}
	message.Image.Three = &service.ThreeObject{}
	message.Image.Three.MeshList = append(message.Image.Three.MeshList,
		service.MeshObject{
			IntersectionObjectName: randomdata.Adjective(),
			LocationName:           randomdata.Adjective(),
			IntersectionPoint: service.Vector3{
				X: RandomFloat(),
				Y: RandomFloat(),
				Z: RandomFloat(),
			},
			Orientation: service.Vector3{
				X: RandomFloat(),
				Y: RandomFloat(),
				Z: RandomFloat(),
			},
			RayDirection: service.Vector3{
				X: RandomFloat(),
				Y: RandomFloat(),
				Z: RandomFloat(),
			},
		})
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func GetImagesInLive() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <live_id>\n", GetFunctionName())
		return
	}

	message := service.Payload{
		Category:    "image",
		Service:     "GetImagesInLive",
		AccessToken: tokens[1],
		LiveId:      tokens[2],
	}

	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func UpdatePatient() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <patient_id>\n", GetFunctionName())
		return
	}

	GetPatient()
	patient := pipeData.Data.Patient
	firstName, middleName, lastName := RandomFullName()
	patient.FirstName = firstName
	patient.MiddleName = middleName
	patient.LastName = lastName

	message := service.Payload{
		Category:    "patient",
		Service:     "UpdatePatient",
		AccessToken: tokens[1],
		Patient:     patient,
	}

	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}

func GetWebsocketHost() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "GetWebsocketHost",
		AccessToken: tokens[1],
	}

	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func UpdateImage() {
	if len(tokens) < 4 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <image_id> <type>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "image",
		Service:     "UpdateImage",
		AccessToken: tokens[1],
	}
	message.Image = &service.ImageObject{
		ImageId:      tokens[2],
		Type:         tokens[3],
		LocationList: RandomLocationList(),
		Note:         randomdata.Adjective(),
		//Treatment:    randomdata.Adjective(),
		TagList: RandomTagList(),
	}

	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func GetDiskUsage() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "public",
		Service:     "GetDiskUsage",
		AccessToken: tokens[1],
	}

	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}

func GetImage() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <image_id>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "image",
		Service:     "GetImage",
		AccessToken: tokens[1],
	}
	message.Image = &service.ImageObject{
		ImageId: tokens[2],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}

func GetLastDiagnosis() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <patient_id>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "patient",
		Service:     "GetLastDiagnosis",
		AccessToken: tokens[1],
	}

	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func GetLastDateInEachPatient() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <patient_id_list>...\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "date",
		Service:     "GetLastDateInEachPatient",
		AccessToken: tokens[1],
	}
	for i, e := range tokens {
		if i > 1 {
			message.PatientList = append(message.PatientList,
				service.PatientObject{
					PatientId: e,
				})
		}
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func GetLastDateInEachDiagnosis() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <diagnosis_id_list>...\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "date",
		Service:     "GetLastDateInEachDiagnosis",
		AccessToken: tokens[1],
	}
	for i, e := range tokens {
		if i > 1 {
			message.DiagnosisList = append(message.DiagnosisList,
				service.DiagnosisObject{
					DiagnosisId: e,
				})
		}
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func GetSlaveInformation() {
	if len(tokens) < 1 {
		fmt.Fprintf(os.Stderr, "%s: <service>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category: "slave",
		Service:  "GetSlaveInformation",
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func SignInSlave() {
	if len(tokens) < 11 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <access_code>(0: empty) <hospital_name> <name> <nickname> <user_id>(0: auto) <hospital_id>(0: auto)  <hospital_create_by>(0: auto) <account> <secret_password>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:     "slave",
		Service:      "SignInSlave",
		Account:      tokens[9],
		SecretToken:  tokens[10],
		AccessToken:  tokens[1],
		HospitalName: tokens[3],
		Name:         tokens[4],
		Nickname:     tokens[5],
	}
	if tokens[2] != "0" {
		message.AccessCode = tokens[2]
	}
	message.UserId = RandomId(0)
	if tokens[6] != "0" {
		message.UserId = tokens[6]
	}
	message.HospitalId = RandomId(1)
	if tokens[7] != "0" {
		message.HospitalId = tokens[7]
	}
	message.HospitalCreatedBy = RandomId(2)
	if tokens[8] != "0" {
		message.HospitalCreatedBy = tokens[8]
	}

	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}

func SetDefaultTag() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <tag>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "SetDefaultTag",
		AccessToken: tokens[1],
		Tag:         tokens[2],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func DelTag() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <tag>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "DelTag",
		AccessToken: tokens[1],
		Tag:         tokens[2],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func AddTag() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "AddTag",
		AccessToken: tokens[1],
		Tag:         RandomTag(),
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}

func AddTags() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <tags_number>\n", GetFunctionName())
		return
	}

	tagNum, err := strconv.Atoi(tokens[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <tags_number>\n", GetFunctionName())
		return
	}
	var tagList []string
	for i := 0; i < tagNum; i++ {
		tagList = append(tagList, RandomTag())
	}
	message := service.Payload{
		Category:    "private",
		Service:     "AddTags",
		AccessToken: tokens[1],
		TagList:     tagList,
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}

func ListTag() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "ListTag",
		AccessToken: tokens[1],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}

func SetDefaultGender() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <gender>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "SetDefaultGender",
		AccessToken: tokens[1],
		Gender:      tokens[2],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func DelGender() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <gender>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "DelGender",
		AccessToken: tokens[1],
		Gender:      tokens[2],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func AddGender() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "AddGender",
		AccessToken: tokens[1],
		Gender:      RandomGender(),
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func ListGender() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "ListGender",
		AccessToken: tokens[1],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}

func SetDefaultLocation() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <location>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "SetDefaultLocation",
		AccessToken: tokens[1],
		Location:    tokens[2],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func DelLocation() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <location>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "DelLocation",
		AccessToken: tokens[1],
		Location:    tokens[2],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func AddLocation() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "AddLocation",
		AccessToken: tokens[1],
		Location:    RandomLocation(),
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func ListLocation() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "ListLocation",
		AccessToken: tokens[1],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}

func SetDefaultDisease() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <disease>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "SetDefaultDisease",
		AccessToken: tokens[1],
		Disease:     tokens[2],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func DelDisease() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <disease>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "DelDisease",
		AccessToken: tokens[1],
		Disease:     tokens[2],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func AddDisease() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "AddDisease",
		AccessToken: tokens[1],
		Disease:     RandomDiseaseType(),
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func ListDisease() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "ListDisease",
		AccessToken: tokens[1],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}

func SetDefaultSkin() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <skin>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "SetDefaultSkin",
		AccessToken: tokens[1],
		Skin:        tokens[2],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func DelSkin() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <skin>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "DelSkin",
		AccessToken: tokens[1],
		Skin:        tokens[2],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func AddSkin() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "AddSkin",
		AccessToken: tokens[1],
		Skin:        RandomSkin(),
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func ListSkin() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "ListSkin",
		AccessToken: tokens[1],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}

func SetDefaultCountry() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <country>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "SetDefaultCountry",
		AccessToken: tokens[1],
		Country:     tokens[2],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func DelCountry() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <country>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "DelCountry",
		AccessToken: tokens[1],
		Country:     tokens[2],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func AddCountry() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "AddCountry",
		AccessToken: tokens[1],
		Country:     randomdata.Country(randomdata.FullCountry),
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func ListCountry() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "ListCountry",
		AccessToken: tokens[1],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}

func SetDefaultEthnicity() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <ethnicity>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "SetDefaultEthnicity",
		AccessToken: tokens[1],
		Ethnicity:   tokens[2],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func DelEthnicity() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <ethnicity>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "DelEthnicity",
		AccessToken: tokens[1],
		Ethnicity:   tokens[2],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func AddEthnicity() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "AddEthnicity",
		AccessToken: tokens[1],
		Ethnicity:   RandomEthnicity(),
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}

func ListEthnicity() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "ListEthnicity",
		AccessToken: tokens[1],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func GetLibraryImages() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "image",
		Service:     "GetLibraryImages",
		AccessToken: tokens[1],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func GetUserInformation() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "private",
		Service:     "GetUserInformation",
		AccessToken: tokens[1],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func GetPlaygroundDate() {
	if len(tokens) < 4 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <skip> <limit> <patient_id_list>x2 <diagnosis_id_list>x2 <date_id_list>x2 \n", GetFunctionName())
		return
	}
	offset := 2
	skip, _ := strconv.Atoi(tokens[offset])
	offset += 1
	limit, _ := strconv.Atoi(tokens[offset])
	offset += 1
	message := service.Payload{
		Category:    "date",
		Service:     "GetPlaygroundDate",
		AccessToken: tokens[1],
		Skip:        int32(skip),
		Limit:       int32(limit),
	}
	for i, e := range tokens {
		if i < offset {
			continue
		} else if i < 6 {
			if e != "0" {
				message.PatientList = append(message.PatientList,
					service.PatientObject{
						PatientId: e,
					})
			}
		} else if i < 8 {
			if e != "0" {
				message.DiagnosisList = append(message.DiagnosisList,
					service.DiagnosisObject{
						DiagnosisId: e,
					})
			}
		} else if i < 10 {
			if e != "0" {
				message.DateList = append(message.DateList,
					service.DateObject{
						DateId: e,
					})
			}
		}
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func GetUserStatus() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <account>\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category: "public",
		Service:  "GetUserStatus",
		Account:  tokens[1],
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	request.Req.Debug("> " + GetFunctionName() + " Request")

	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
}
func MoveImagesToDate() {
	if len(tokens) < 5 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <diagnosis_id> <visit_date> <image_id_list>...\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "image",
		Service:     "MoveImagesToDate",
		AccessToken: tokens[1],
	}
	for i, e := range tokens {
		if i < 2 {
			continue
		} else if i == 2 {
			message.Diagnosis = &service.DiagnosisObject{
				DiagnosisId: tokens[i],
			}
		} else if i == 3 {
			message.Date = &service.DateObject{
				VisitDate: tokens[i],
			}
		} else {
			message.ImageList = append(message.ImageList,
				service.ImageObject{
					ImageId: e,
				})
		}
	}

	request := endpoint.ApiRequest{
		Req: message,
	}

	request.Req.Debug("> " + GetFunctionName() + " Request")

	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
}
func MoveImagesToLibrary() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <image_id_list>...\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "image",
		Service:     "MoveImagesToLibrary",
		AccessToken: tokens[1],
	}
	for i, e := range tokens {
		if i < 2 {
			continue
		} else {
			message.ImageList = append(message.ImageList,
				service.ImageObject{
					ImageId: e,
				})
		}
	}

	request := endpoint.ApiRequest{
		Req: message,
	}

	request.Req.Debug("> " + GetFunctionName() + " Request")

	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
}
func MakeTestCase() {
	if len(tokens) < 4 {
		fmt.Fprintf(os.Stderr, "%s: <service> <account> <password> <patient_number>\n", GetFunctionName())
		return
	}

	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	account := tokens[1]
	password := tokens[2]
	patientNumber := tokens[3]

	testAccessToken := ""
	tokens = []string{"SignIn", account, password}
	SignIn()
	if len(accessToken) == 0 {
		tokens = []string{"SignUp", account}
		SignUp()

		tokens = []string{"GetCertificationCode", account}
		code := GetCertificationCode()

		tokens = []string{"VerifyCertificationCode", account, code}
		VerifyCertificationCode()

		tokens = []string{"SignUpPassword", account, password}
		testAccessToken = SignUpPassword()

		tokens = []string{"SignUpComplete", testAccessToken, RandomKoreanLastName() + RandomKoreanFirstName(), randomdata.Adjective(), randomdata.Adjective()}
		SignUpComplete()
	} else {
		testAccessToken = accessToken
	}

	tokens = []string{"AddPatient", testAccessToken, patientNumber}
	patientIdList := RepeatAddPatient()

	for _, ep := range patientIdList {
		tokens = []string{"AddDiagnosis", testAccessToken, ep, "0", RandomLocation(), RandomLocation(), RandomTag(), RandomTag(), RandomDate().Format(service.GetTimeFormat()), strconv.Itoa(randomdata.Number(1, 5))}
		fmt.Fprintf(os.Stderr, "=> %v\n", tokens)
		diagnosisIdList := RepeatAddDiagnosis()
		for _, edg := range diagnosisIdList {

			tokens = []string{"GetDiagnosis", testAccessToken, edg}
			diagnosis := GetDiagnosis()
			var additionalList []string
			for _, e2 := range diagnosis.DateList {
				additionalList = append(additionalList, e2.DateId)
			}

			tokens = []string{"AddDate", testAccessToken, ep, edg, "0", strconv.Itoa(randomdata.Number(3, 10))}
			fmt.Fprintf(os.Stderr, "=> %v\n", tokens)
			dateIdList := RepeatAddDate()
			dateIdList = append(dateIdList, additionalList...)
			for _, ed := range dateIdList {
				n := rand.Intn(15) + 1
				for i := 0; i < n; i += 1 {
					fileName := RandomTestImageFileName()
					tokens = []string{"PrepareImage", testAccessToken, fileName, ed}
					fmt.Fprintf(os.Stderr, "=> %v\n", tokens)
					image := PrepareImage()
					tokens = []string{"UpdateImage", testAccessToken, image.ImageId, diagnosis.Type}
					UpdateImage()
					go UploadFileAndProcess(testAccessToken, fileName, image)
				}
				n = rand.Intn(1) + 1
				for i := 0; i < n; i += 1 {
					fileName := RandomTestImageFileName()
					tokens = []string{"PrepareImage", testAccessToken, fileName}
					fmt.Fprintf(os.Stderr, "=> %v\n", tokens)
					image := PrepareImage()
					tokens = []string{"UpdateImage", testAccessToken, image.ImageId, diagnosis.Type}
					UpdateImage()
					go UploadFileAndProcess(testAccessToken, fileName, image)
				}
			}
		}
	}
}
func UploadFileAndProcess(testAccessToken string, fileName string, image *service.ImageObject) {
	UploadFile(image.Uri.Image, "sample/"+fileName, "file")
	fmt.Fprintf(os.Stderr, "=> %v %v\n", testAccessToken, image.ImageId)
	ProcessImage(testAccessToken, image.ImageId)
}
func GetImagesByDate() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <date_id_list>...\n", GetFunctionName())
		return
	}
	message := service.Payload{
		Category:    "image",
		Service:     "GetImagesInEachDate",
		AccessToken: tokens[1],
	}
	for i, e := range tokens {
		if i < 2 {
			continue
		} else {
			message.DateList = append(message.DateList,
				service.DateObject{
					DateId: e,
				})
		}
	}

	request := endpoint.ApiRequest{
		Req: message,
	}

	request.Req.Debug("> " + GetFunctionName() + " Request")

	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
}
func GetDate() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <date_id>\n", GetFunctionName())
		return
	}
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category:    "date",
			Service:     "GetDate",
			AccessToken: tokens[1],
			Date: &service.DateObject{
				DateId: tokens[2],
			},
		},
	}

	request.Req.Debug("> " + GetFunctionName() + " Request")

	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
}
func SearchDate() {
	if len(tokens) < 6 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <skip> <limit> <order> <order_by> <patient_id> <diagnosis_id> <date_keyword_list>x2  <patient_keyword_list>x2 <disease_keyword_list>x2 \n", GetFunctionName())
		return
	}
	offset := 2
	skip, _ := strconv.Atoi(tokens[offset])
	offset += 1
	limit, _ := strconv.Atoi(tokens[offset])
	offset += 1
	message := service.Payload{
		Category:    "date",
		Service:     "SearchDate",
		AccessToken: tokens[1],
		Skip:        int32(skip),
		Limit:       int32(limit),
	}
	for i, e := range tokens {
		if i < offset {
			continue
		} else if i == 4 {
			message.Order = e
		} else if i == 5 {
			message.OrderBy = e
		} else if i == 6 {
			if e != "0" {
				message.Patient = &service.PatientObject{
					PatientId: e,
				}
			}
		} else if i == 7 {
			if e != "0" {
				message.Diagnosis = &service.DiagnosisObject{
					DiagnosisId: e,
				}
			}
		} else if i < 10 {
			message.DateKeywordList = append(message.DateKeywordList, e)
		} else if i < 12 {
			message.PatientKeywordList = append(message.PatientKeywordList, e)
		} else if i < 14 {
			message.DiseaseKeywordList = append(message.DiseaseKeywordList, e)
		}
	}
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func GetDiagnosis() (diagnosis *service.DiagnosisObject) {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <diagnosis_id>\n", GetFunctionName())
		return
	}
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category:    "diagnosis",
			Service:     "GetDiagnosis",
			AccessToken: tokens[1],
			Diagnosis: &service.DiagnosisObject{
				DiagnosisId: tokens[2],
			},
		},
	}

	request.Req.Debug("> " + GetFunctionName() + " Request")

	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	diagnosis = respBody.Data.Diagnosis
	respBody.Debug("< " + GetFunctionName() + " Response")

	return
}
func SearchDisease() {
	if len(tokens) < 6 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <skip> <limit> <order> <order_by> <patient_id>(0: empty) <disease_keyword_list>x2 <patient_keyword_list>x2 <date_keyword_list>x2\n", GetFunctionName())
		return
	}
	offset := 2
	skip, _ := strconv.Atoi(tokens[offset])
	offset += 1
	limit, _ := strconv.Atoi(tokens[offset])
	offset += 1
	message := service.Payload{
		Category:    "diagnosis",
		Service:     "SearchDisease",
		AccessToken: tokens[1],
		Skip:        int32(skip),
		Limit:       int32(limit),
	}
	for i, e := range tokens {
		if i < offset {
			continue
		} else if i == 4 {
			message.Order = e
		} else if i == 5 {
			message.OrderBy = e
		} else if i == 6 {
			if e != "0" {
				message.Patient = &service.PatientObject{
					PatientId: e,
				}
			}
		} else if i < 9 {
			message.DiseaseKeywordList = append(message.DiseaseKeywordList, e)
		} else if i < 11 {
			message.PatientKeywordList = append(message.PatientKeywordList, e)
		} else if i < 13 {
			message.DateKeywordList = append(message.DateKeywordList, e)
		}
	}
	defer timeTrack(time.Now(), GetFunctionName())
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")

	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func RepeatAddDate() (dateIdList []string) {
	repeat := 1
	if len(tokens) > 5 {
		repeat, _ = strconv.Atoi(tokens[5])
	}
	defer timeTrack(time.Now(), GetFunctionName())
	for i := 0; i < repeat; i++ {
		dateId := AddDate()
		dateIdList = append(dateIdList, dateId)
	}
	return
}
func AddDate() (dateId string) {
	if len(tokens) < 5 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <patient_id>(0: auto) <diagnosis_id>(0: auto) <location>(0: auto)\n", GetFunctionName())
		return
	}

	patientId := RandomId(1)
	if tokens[2] != "0" {
		patientId = tokens[2]
	}
	patient := service.PatientObject{
		PatientId: patientId,
	}

	diagnosisId := RandomId(2)
	if tokens[3] != "0" {
		diagnosisId = tokens[3]
	}
	diagnosis := service.DiagnosisObject{
		DiagnosisId: diagnosisId,
	}

	date := service.DateObject{
		VisitDate: RandomDate().Format(service.GetTimeFormat()),
		Note:      randomdata.Adjective(),
		TagList:   RandomTagList(),
		Treatment: randomdata.Adjective(),
	}

	if tokens[4] == "0" {
		date.LocationList = append(date.LocationList, RandomLocation())
		date.LocationList = append(date.LocationList, RandomLocation())
	}

	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category:    "date",
			Service:     "AddDate",
			AccessToken: tokens[1],
			Patient:     &patient,
			Diagnosis:   &diagnosis,
			Date:        &date,
		},
	}

	request.Req.Debug("> " + GetFunctionName() + " Request")

	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")

	dateId = respBody.Data.Date.DateId
	return
}
func GetPatient() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <patient_id>\n", GetFunctionName())
		return
	}
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category:    "patient",
			Service:     "GetPatient",
			AccessToken: tokens[1],
			Patient: &service.PatientObject{
				PatientId: tokens[2],
			},
		},
	}

	request.Req.Debug("> " + GetFunctionName() + " Request")

	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	defer timeTrack(time.Now(), GetFunctionName())
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	pipeData = respBody
}
func RepeatAddDiagnosis() (diagnosisIdList []string) {
	repeat := 1
	if len(tokens) > 9 {
		repeat, _ = strconv.Atoi(tokens[9])
	}
	defer timeTrack(time.Now(), GetFunctionName())
	for i := 0; i < repeat; i++ {
		diagnosisId := AddDiagnosis()
		diagnosisIdList = append(diagnosisIdList, diagnosisId)
	}
	return
}
func AddDiagnosis() (diagnosisId string) {
	if len(tokens) < 9 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <patient_id>(0: empty) <type>(0: empty) <location>x2 <tag_list>x2 <visit_date>(0: empty) <image_id>\n", GetFunctionName())
		return
	}

	fmt.Fprintf(os.Stdout, "AddDiagnosis")
	message := service.Payload{
		Category:    "diagnosis",
		Service:     "AddDiagnosis",
		AccessToken: tokens[1],
	}
	message.Patient = &service.PatientObject{
		PatientId: RandomId(1),
	}
	if tokens[2] != "0" {
		message.Patient.PatientId = tokens[2]
	}
	message.Diagnosis = &service.DiagnosisObject{
		Type: RandomDiseaseType(),
	}
	if tokens[3] != "0" {
		message.Diagnosis.Type = tokens[3]
	}

	//message.Diagnosis.LocationList = append(message.Diagnosis.LocationList, tokens[4])
	//message.Diagnosis.LocationList = append(message.Diagnosis.LocationList, tokens[5])
	//changed location => three
	message.Three = &service.ThreeObject{}
	message.Three.MeshList = append(message.Three.MeshList,
		service.MeshObject{
			IntersectionObjectName: randomdata.Adjective(),
			LocationName:           randomdata.Adjective(),
			IntersectionPoint: service.Vector3{
				X: RandomFloat(),
				Y: RandomFloat(),
				Z: RandomFloat(),
			},
			Orientation: service.Vector3{
				X: RandomFloat(),
				Y: RandomFloat(),
				Z: RandomFloat(),
			},
			RayDirection: service.Vector3{
				X: RandomFloat(),
				Y: RandomFloat(),
				Z: RandomFloat(),
			},
		})

	message.Diagnosis.TagList = append(message.Diagnosis.TagList, tokens[6])
	message.Diagnosis.TagList = append(message.Diagnosis.TagList, tokens[7])

	date := RandomDate()
	message.Diagnosis.Date = date.Format(service.GetTimeFormat())
	if tokens[8] != "0" {
		message.Diagnosis.Date = tokens[8]
	}
	message.Diagnosis.Note = randomdata.Adjective()
	message.Diagnosis.Treatment = randomdata.Adjective()

	request := endpoint.ApiRequest{
		Req: message,
	}
	request.Req.Debug("> " + GetFunctionName() + " Request")

	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer timeTrack(time.Now(), GetFunctionName())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")

	if len(respBody.Err) == 0 {
		diagnosisId = respBody.Data.Diagnosis.DiagnosisId
	}
	return
}
func RepeatAddPatient() (patientIdList []string) {
	repeat := 1
	if len(tokens) > 2 {
		repeat, _ = strconv.Atoi(tokens[2])
	}
	defer timeTrack(time.Now(), GetFunctionName())
	for i := 0; i < repeat; i++ {
		patientId := AddPatient()
		patientIdList = append(patientIdList, patientId)
	}
	return
}
func AddPatient() (patientId string) {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}

	birthDate := RandomDate()
	firstName, middleName, lastName := RandomFullName()

	patient := service.PatientObject{
		FirstName:               firstName,
		MiddleName:              middleName,
		LastName:                lastName,
		Phone:                   GeneratePhone(),
		Pid:                     GeneratePid(),
		Country:                 randomdata.Country(randomdata.FullCountry),
		Address:                 randomdata.Address(),
		Gender:                  RandomGender(),
		Ethnicity:               RandomEthnicity(),
		BirthDate:               birthDate.Format(service.GetTimeFormat()),
		SkinType:                RandomSkin(),
		Description:             randomdata.Adjective(),
		SkinCancerHistory:       RandomBoolean(),
		FamilySkinCancerHistory: RandomBoolean(),
	}
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category:    "patient",
			Service:     "AddPatient",
			AccessToken: tokens[1],
			Patient:     &patient,
		},
	}

	request.Req.Debug("> " + GetFunctionName() + " Request")

	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	patientId = respBody.Data.Patient.PatientId
	return
}
func SearchPatient() {
	if len(tokens) < 6 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <skip> <limit> <order> <order_by> <patient_keyword_list>x2 <disease_keyword_list>x2 <date_keyword_list>x2\n", GetFunctionName())
		return
	}
	offset := 2
	skip, _ := strconv.Atoi(tokens[offset])
	offset += 1
	limit, _ := strconv.Atoi(tokens[offset])
	offset += 1
	message := service.Payload{
		Category:    "patient",
		Service:     "SearchPatient",
		AccessToken: tokens[1],
		Skip:        int32(skip),
		Limit:       int32(limit),
	}
	for i, e := range tokens {
		if i < offset {
			continue
		} else if i == 4 {
			message.Order = e
		} else if i == 5 {
			message.OrderBy = e
		} else if i < 8 {
			message.PatientKeywordList = append(message.PatientKeywordList, e)
		} else if i < 10 {
			message.DiseaseKeywordList = append(message.DiseaseKeywordList, e)
		} else if i < 12 {
			message.DateKeywordList = append(message.DateKeywordList, e)
		}
	}
	defer timeTrack(time.Now(), GetFunctionName())
	request := endpoint.ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")

	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	fmt.Fprintf(os.Stderr, "Result: %v\n", respBody.Data.Count)
}
func SignUpComplete() (hospitalId string) {
	if len(tokens) < 5 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <name> <nickname> <hospital_name>\n", GetFunctionName())
		return
	}
	defer timeTrack(time.Now(), GetFunctionName())
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category:     "private",
			Service:      "SignUpComplete",
			AccessToken:  tokens[1],
			Name:         tokens[2],
			Nickname:     tokens[3],
			HospitalName: tokens[4],
		},
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	hospitalId = respBody.Data.HospitalId

	return
}
func SignUpPassword() (accessToken string) {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <account> <password>\n", GetFunctionName())
		return
	}
	defer timeTrack(time.Now(), GetFunctionName())
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category: "public",
			Service:  "SignUpPassword",
			Account:  tokens[1],
			Password: tokens[2],
		},
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	accessToken = respBody.Data.AccessToken
	return
}
func CreateHospital() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <hospital_name>\n", GetFunctionName())
		return
	}
	defer timeTrack(time.Now(), GetFunctionName())
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category:     "private",
			Service:      "CreateHospital",
			AccessToken:  tokens[1],
			HospitalName: tokens[2],
		},
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
}
func GetThumbnailUri() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <image_id>\n", GetFunctionName())
		return
	}
	defer timeTrack(time.Now(), GetFunctionName())
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category:    "image",
			Service:     "GetThumbnailUri",
			AccessToken: tokens[1],
			Image: &service.ImageObject{
				ImageId: tokens[2],
			},
		},
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	respBody.Debug("< " + GetFunctionName() + " Response")
}
func ProcessImage(accessToken string, imageId string) {
	message := service.Payload{
		Category: "image",
		Service:  "ProcessImage",
	}

	if len(accessToken) > 0 {
		message.AccessToken = accessToken
		message.Image = &service.ImageObject{
			ImageId: imageId,
		}
	} else {
		if len(tokens) < 3 {
			fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <image_id>\n", GetFunctionName())
			return
		}
		message.AccessToken = tokens[1]
		message.Image = &service.ImageObject{
			ImageId: tokens[2],
		}
	}

	defer timeTrack(time.Now(), GetFunctionName())
	//width, _ := strconv.Atoi(tokens[3])
	//height, _ := strconv.Atoi(tokens[4])
	request := endpoint.ApiRequest{
		Req: message,
	}

	request.Req.Debug("> " + GetFunctionName() + " Request")

	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	respBody.Debug("< " + GetFunctionName() + " Response")
}
func GetImageUri() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <image_id>\n", GetFunctionName())
		return
	}
	defer timeTrack(time.Now(), GetFunctionName())
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category:    "image",
			Service:     "GetImageUri",
			AccessToken: tokens[1],
			Image: &service.ImageObject{
				ImageId: tokens[2],
			},
		},
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
}
func GetSystemType() {
	defer timeTrack(time.Now(), GetFunctionName())
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category: "public",
			Service:  "GetSystemType",
		},
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
}
func GetPlatform() {
	defer timeTrack(time.Now(), GetFunctionName())
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category: "public",
			Service:  "GetPlatform",
		},
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
}
func UploadImage() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <uri> <file_path>\n", GetFunctionName())
		return
	}
	UploadFile(tokens[1], tokens[2], "file")
}
func PrepareImage() (image *service.ImageObject) {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <image_file_name> <date_id> <live_id>\n", GetFunctionName())
		return
	}
	defer timeTrack(time.Now(), GetFunctionName())
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category:    "image",
			Service:     "PrepareImage",
			AccessToken: tokens[1],
			Image: &service.ImageObject{
				FileName: tokens[2],
			},
		},
	}
	if len(tokens) == 4 {
		request.Req.Date = &service.DateObject{
			DateId: tokens[3],
		}
	}
	if len(tokens) == 5 {
		request.Req.LiveId = tokens[4]
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")

	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")

	image = respBody.Data.Image
	return
}
func SignIn() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <account> <password>\n", GetFunctionName())
		return
	}
	defer timeTrack(time.Now(), GetFunctionName())
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category: "public",
			Service:  "SignIn",
			Account:  tokens[1],
			Password: tokens[2],
		},
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	accessToken = respBody.Data.AccessToken
	respBody.Debug("< " + GetFunctionName() + " Response")
}
func SignUp() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <account>\n", GetFunctionName())
		return
	}
	defer timeTrack(time.Now(), GetFunctionName())

	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category: "public",
			Service:  "SignUp",
			Account:  tokens[1],
		},
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
}
func GetCertificationCode() (code string) {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <account>\n", GetFunctionName())
		return
	}
	defer timeTrack(time.Now(), GetFunctionName())
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category: "public",
			Service:  "GetCertificationCode",
			Account:  tokens[1],
		},
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
	code = respBody.Data.Code

	return
}
func VerifyCertificationCode() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <account> <code>\n", GetFunctionName())
		return
	}
	defer timeTrack(time.Now(), GetFunctionName())
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category: "public",
			Service:  "VerifyCertificationCode",
			Account:  tokens[1],
			Code:     tokens[2],
		},
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
}
func ResendCertificationCode() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <account>\n", GetFunctionName())
		return
	}
	defer timeTrack(time.Now(), GetFunctionName())
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category: "public",
			Service:  "ResendCertificationCode",
			Account:  tokens[1],
		},
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
}
func UpdatePassword() {
	if len(tokens) < 4 {
		fmt.Fprintf(os.Stderr, "%s: <service> <account> <password> <new_password>\n", GetFunctionName())
		return
	}
	defer timeTrack(time.Now(), GetFunctionName())
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category:    "public",
			Service:     "UpdatePassword",
			Account:     tokens[1],
			Password:    tokens[2],
			NewPassword: tokens[3],
		},
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
}
func UpdateUserInformation() {
	if len(tokens) < 5 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <name> <nickname> <hospital_id>\n", GetFunctionName())
		return
	}
	defer timeTrack(time.Now(), GetFunctionName())
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category:    "private",
			Service:     "UpdateUserInformation",
			AccessToken: tokens[1],
			Name:        tokens[2],
			Nickname:    tokens[3],
			HospitalId:  tokens[4],
		},
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
}
func UpdateAccessToken() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	defer timeTrack(time.Now(), GetFunctionName())
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category:    "private",
			Service:     "UpdateAccessToken",
			AccessToken: tokens[1],
		},
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
}
func PrepareAvatar() {
	if len(tokens) < 3 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token> <avatar_file_name>\n", GetFunctionName())
		return
	}
	defer timeTrack(time.Now(), GetFunctionName())
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category:    "private",
			Service:     "PrepareAvatar",
			AccessToken: tokens[1],
			Avatar:      tokens[2],
		},
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
}
func GetAvatarUri() {
	if len(tokens) < 2 {
		fmt.Fprintf(os.Stderr, "%s: <service> <access_token>\n", GetFunctionName())
		return
	}
	defer timeTrack(time.Now(), GetFunctionName())
	request := endpoint.ApiRequest{
		Req: service.Payload{
			Category:    "private",
			Service:     "GetAvatarUri",
			AccessToken: tokens[1],
		},
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	request.Req.Debug("> " + GetFunctionName() + " Request")
	req, err := http.NewRequest("POST", scheme, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	respBody.Debug("< " + GetFunctionName() + " Response")
}

type Response struct {
	Data service.Payload `json:"data"`
	Err  string          `json:"err"`
}

func (p *Response) Debug(prefix string) {
	fmt.Fprintf(os.Stderr, "%v: \n", prefix)
	if j, err2 := json.MarshalIndent(p, "", " "); err2 != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err2)
	} else {
		fmt.Fprintf(os.Stderr, "%v\n", string(j))
	}
	fmt.Fprintf(os.Stdout, "\n")
}

func GetFunctionName() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])

	tokens := strings.Split(f.Name(), ".")

	return tokens[len(tokens)-1]
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Fprintf(os.Stderr, "%s took %v\n", name, elapsed.Seconds())
}

func showUsage() {
	fmt.Fprintf(os.Stdout, "\n")
	for i, e := range services {
		fmt.Fprintf(os.Stdout, "%3d) %s\n", i, e)
	}
	fmt.Fprintf(os.Stdout, "\n")
}
func GeneratePhone() string {
	n := rand.Intn(100)
	if n < 70 {
		return "010-" + strconv.Itoa(rand.Intn(8999)+1000) + "-" + strconv.Itoa(rand.Intn(8999)+1000)
	}
	return randomdata.PhoneNumber()
}
func GeneratePid() string {
	n := rand.Intn(100)
	if n < 10 {
		return ""
	}

	randomPid := randomdata.Alphanumeric(5) + bson.NewObjectIdWithTime(time.Now()).Hex()[3:5]
	return randomPid
}
func RandomGender() string {
	n := rand.Intn(100)
	if n < 50 {
		return "M"
	}
	return "F"
}
func RandomFloat() float32 {
	n := rand.Float32()
	return n
}
func RandomBoolean() string {
	n := rand.Intn(100)
	if n < 50 {
		return "No"
	}
	return "Yes"
}
func RandomEthnicity() string {
	ethnicityList := []string{
		"Caucasian", "Latino", "Asian", "African", "Arab",
		"Caucasian", "Latino", "Asian", "African", "Arab",
		"Caucasian", "Latino", "Asian", "African", "Arab",
	}
	n := rand.Intn(len(ethnicityList) - 1)
	return ethnicityList[n]
}
func RandomKoreanFirstName() string {
	firstNameList := []string{
		"민준", "서준", "예준", "주원", "도윤", "시우", "지후", "현우", "준우", "지훈", "도현", "건우", "우진", "민재", "현준",
		"선우", "서진", "연우", "정우", "유준", "승현", "준혁", "승우", "지환", "승민", "시윤", "지우", "민성", "유찬", "준영", "진우",
		"시후", "지원", "은우", "윤우", "수현", "동현", "재윤", "민규", "시현", "태윤", "재원", "민우", "재민", "은찬", "한결", "윤호",
		"민찬", "시원", "성민", "성현", "수호", "준호", "승준", "현서", "재현", "시온", "지성", "태민", "태현", "민혁", "예성", "민호",
		"하율", "지안", "성준", "우빈", "지율", "정민", "규민", "윤성", "지한", "민석", "지민", "이준", "준", "준수", "서우", "은호",
		"은성", "예찬", "이안", "윤재", "율", "하람", "태양", "준희", "준성", "하진", "현수", "도훈", "승원", "정현", "건", "지완",
		"선우", "서진", "연우", "정우", "유준", "승현", "준혁", "승우", "지환", "승민", "시윤", "지우", "민성", "유찬", "준영", "진우",
		"시후", "지원", "은우", "윤우", "수현", "동현", "재윤", "민규", "시현", "태윤", "재원", "민우", "재민", "은찬", "한결", "윤호",
		"민찬", "시원", "성민", "성현", "수호", "준호", "승준", "현서", "재현", "시온", "지성", "태민", "태현", "민혁", "예성", "민호",
		"하율", "지안", "성준", "우빈", "지율", "정민", "규민", "윤성", "지한", "민석", "지민", "이준", "준", "준수", "서우", "은호",
		"은성", "예찬", "이안", "윤재", "율", "하람", "태양", "준희", "준성", "하진", "현수", "도훈", "승원", "정현", "건", "지완",
		"선우", "서진", "연우", "정우", "유준", "승현", "준혁", "승우", "지환", "승민", "시윤", "지우", "민성", "유찬", "준영", "진우",
		"시후", "지원", "은우", "윤우", "수현", "동현", "재윤", "민규", "시현", "태윤", "재원", "민우", "재민", "은찬", "한결", "윤호",
		"민찬", "시원", "성민", "성현", "수호", "준호", "승준", "현서", "재현", "시온", "지성", "태민", "태현", "민혁", "예성", "민호",
		"하율", "지안", "성준", "우빈", "지율", "정민", "규민", "윤성", "지한", "민석", "지민", "이준", "준", "준수", "서우", "은호",
		"은성", "예찬", "이안", "윤재", "율", "하람", "태양", "준희", "준성", "하진", "현수", "도훈", "승원", "정현", "건", "지완",
		"강민", "승호", "율", "준", "윤", "건", "봄", "현", "솔", "산", "별", "찬", "민", "설", "진", "원", "결", "환", "강", "은", "린", "훈",
		"겸", "혁", "단", "한", "슬", "빈", "선", "호", "수", "담", "유", "범", "연", "희", "신", "휘", "정", "온", "훤", "안", "비", "권", "영",
		"도", "완", "인", "운", "후", "헌", "솜", "랑", "성", "주", "우", "률", "경", "엘", "란", "레", "승", "송", "든", "리", "샘", "늘", "룩",
		"웅", "을", "용", "하", "들", "림", "본", "석", "빛", "욱", "명", "해", "상", "금", "람", "홍", "탄", "이", "아", "미", "나", "철", "륜",
		"국", "서", "규", "룬", "루", "곤", "름", "언", "지", "휼", "효",
	}
	n := rand.Intn(len(firstNameList) - 1)
	return firstNameList[n]
}
func RandomKoreanLastName() string {
	lastNameList := []string{
		"김", "이", "박", "최", "정", "강", "조", "윤", "장", "임", "한", "오", "서", "신", "권", "황", "안", "송", "전",
		"홍", "유", "고", "문", "양", "손", "배", "조", "백", "허", "유", "남", "심", "노", "정", "하", "곽", "성", "차",
		"주", "우", "구", "민", "유", "류", "나", "엄", "원", "천", "방", "공", "남궁", "황보", "모", "장", "기", "반",
		"명", "맹", "제", "탁", "국", "여", "어", "은", "구", "석", "사", "가", "시", "갈", "호", "설", "팽", "사공", "음",
	}
	n := rand.Intn(len(lastNameList) - 1)
	return lastNameList[n]
}
func RandomFullName() (string, string, string) {
	n := rand.Intn(100)
	if n < 50 {
		return RandomKoreanFirstName(), "", RandomKoreanLastName()
	}

	return randomdata.FirstName(randomdata.RandomGender), randomdata.FirstName(randomdata.RandomGender), randomdata.LastName()
}

func RandomDate() time.Time {
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func RandomDiseaseType() string {
	diseaseTypeList := []string{
		"Acne", "ActinicKeratosis", "Birthmarks", "Blisters", "CherryAngiomas", "ColdSores", "DrySkin", "Eczema",
		"FungalNail", "Melasma", "Moles", "Psoriasis", "Rashes", "Rosacea", "Scabies", "Scars",
		"Shingles", "SkinAllergies", "SkinCancer", "Irritants",
		"Vitiligo", "Warts",
		"Acne", "ActinicKeratosis", "Birthmarks", "Blisters", "CherryAngiomas", "ColdSores", "DrySkin", "Eczema",
		"FungalNail", "Melasma", "Moles", "Psoriasis", "Rashes", "Rosacea", "Scabies", "Scars",
		"Shingles", "SkinAllergies", "SkinCancer", "Irritants",
		"Vitiligo", "Warts",
	}
	n := rand.Intn(len(diseaseTypeList) - 1)
	return diseaseTypeList[n]
}
func RandomLocation() string {
	locationList := []string{
		"Wrists", "Forearms", "Genitals", "Legs", "Face", "Eyes", "Mouth", "Forehead", "Nose", "Ear", "Groin", "Breasts",
		"Wrists", "Forearms", "Genitals", "Legs", "Face", "Eyes", "Mouth", "Forehead", "Nose", "Ear", "Groin", "Breasts",
	}
	n := rand.Intn(len(locationList) - 1)
	return locationList[n]
}
func RandomSkin() string {
	skinList := []string{
		"Dry", "Wet", "Normal", "Oily", "Sensitive", "Aging",
		"Dry", "Wet", "Normal", "Oily", "Sensitive", "Aging",
	}
	n := rand.Intn(len(skinList) - 1)
	return skinList[n]
}
func RandomTestImageFileName() string {
	files, err := ioutil.ReadDir("sample")
	if err != nil {
		stdLog.Fatal(err)
	}

	var fileNameList []string
	for _, f := range files {
		fileNameList = append(fileNameList, f.Name())
	}

	n := rand.Intn(len(fileNameList) - 1)

	return fileNameList[n]
}
func RandomTag() string {
	tagList := []string{
		"10", "20", "30", "40", "50", "60", "70", "남자", "여자", "심각함", "암",
	}
	n := rand.Intn(len(tagList) - 1)

	return tagList[n]
}
func RandomTagList() []string {
	tagList := []string{
		"10", "20", "30", "40", "50", "60", "70", "남자", "여자", "심각함", "암", "중증", "심신미약", "critical",
	}
	c := rand.Intn(5) + 1
	var t map[string]bool
	var l []string
	for i := 0; i < c; i++ {
		n := rand.Intn(len(tagList) - 1)
		if t[tagList[n]] == false {
			l = append(l, tagList[n])
		}
	}
	return l
}
func RandomDiseaseTypeList() []string {
	var typeList []string
	n := rand.Intn(5) + 1
	for i := 0; i < n; i++ {
		typeList = append(typeList, RandomDiseaseType())
	}
	return typeList
}
func RandomLocationList() []string {
	var l []string
	n := rand.Intn(1) + 1
	for i := 0; i < n; i++ {
		l = append(l, RandomLocation())
	}
	return l
}
func RandomId(i int) string {
	if ids == nil {
		ids = make(map[int][]string)
	}
	c := rand.Intn(100)
	id := ""
	if len(ids[i]) == 0 || c < 80 {
		id = bson.NewObjectId().Hex()
		ids[i] = append(ids[i], id)
	} else {
		n := rand.Intn(len(ids[i]))
		id = ids[i][n]
	}

	return id
}
func UploadFile(url string, filename string, filetype string) []byte {
	file, err := os.Open(filename)

	if err != nil {
		stdLog.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(filetype, filepath.Base(file.Name()))

	if err != nil {
		stdLog.Fatal(err)
	}

	io.Copy(part, file)
	writer.Close()
	request, err := http.NewRequest("POST", url, body)

	if err != nil {
		stdLog.Fatal(err)
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		stdLog.Fatal(err)
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		stdLog.Fatal(err)
	}
	fmt.Fprintf(os.Stderr, "< "+GetFunctionName()+" Ok\n")

	return content
}
