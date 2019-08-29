package service

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"
)

type FileObject struct {
	Path       string    `bson:"path,omitempty"`
	CreateTime time.Time `bson:"create_time,omitempty"`
}

// TODO: UTC()를 써야할 때가 있고 아닐 때가 있다. 어떻게 해야 로컬 시간으로 할 수가 있지??
// => Local() 로 해결
func (f *FileObject) PrepareAvatarPath(userId string, fileName string) {
	f.CreateTime = time.Now().UTC()
	f.Path = "u/" + userId + "/p/a/" + f.CreateTime.Local().Format("20060102150405") + "/" + fileName
}
func (f *FileObject) GetFileUri(accessToken string) (uri string) {
	if len(f.Path) == 0 {
		return ""
	}
	host := GetConfigClientFileHttp()
	//fmt.Fprintf(os.Stderr, "%v\n", host)
	if host[len(host)-1] != '/' {
		host = host + "/"
	}
	// TODO: 파일 서버가 분리되도 DB 는 하나를 바라보게 해서 AccessToken 을 공유하자
	uri = host + "files/" + accessToken + "/" + f.Path
	return
}

// 실제 경로는 파일 서버만 안다.
// PikabuImage 가 이미지 파일들을 관리하기는 하지만 그것은 DB 에서 관리되는 상대 경로 및 관계 정보이다.
// Static 한 파일들은 모두 파일 서버(=PikabuControl)이 관리하도록 한다.
func StaticDataFilePath(account string, id string, datetime time.Time, fileName string) (filePath string) {

	filePath = os.Getenv("PIKABU_DATA")
	filePath += "/" + datetime.Format("2006-01-02")
	filePath += "/" + id + "_" + datetime.Format("20060102150405") + "_" + fileName

	return
}
func StaticSystemFilePath(separators []string, id string, datetime time.Time, fileName string) (filePath string) {

	filePath = os.Getenv("PIKABU_HOME") + "/system"
	filePath += "/" + separators[0]
	filePath += "/" + id
	filePath += "/" + separators[1]
	filePath += "/" + separators[2]
	filePath += "/" + datetime.Format("20060102150405") + "_" + fileName

	return
}
func StaticProfileAvatarFilePath(account string, datetime time.Time, fileName string) (filePath string) {

	filePath = os.Getenv("PIKABU_HOME") + "/system/account"
	filePath += "/" + account
	filePath += "/profile/avatar"
	filePath += "/" + datetime.Format("20060102150405") + "_" + fileName

	return
}
// TODO: 파일쓰기 작업은 go routine 으로 해야하지 않을까?
func FilePost(r *http.Request, filePath string) (err error) {
	if err = filePostWithForm(r, filePath); err == nil {
		return
	}
	CreateDirectory(filePath)
	f, err := os.Create(filePath)
	if err != nil {
		err = fmt.Errorf("fail to create file: %s", err)
		return
	}
	defer f.Close()

	buf, _ := ioutil.ReadAll(r.Body)
	if _, err = f.Write(buf); err != nil {
		err = fmt.Errorf("fail to write file: %s", err)
		return
	}

	return
}
func filePostWithForm(r *http.Request, filePath string) (err error) {
	if err = r.ParseForm(); err != nil {
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		err = fmt.Errorf("fail to read form: %v", err)
		return
	}
	defer file.Close()

	CreateDirectory(filePath)
	f, err := os.Create(filePath)
	if err != nil {
		err = fmt.Errorf("fail to create file: %s", err)
		return
	}
	defer f.Close()

	io.Copy(f, file)

	return
}
func CreateDirectory(filePath string) bool {
	dirName := path.Dir(filePath)
	src, err := os.Stat(dirName)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirName, 0755)
		if errDir != nil {
			panic(err)
		}
		return true
	}

	if src.Mode().IsRegular() {
		return false
	}

	return false
}
