package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Payload struct {
	Category            string                      `json:"category,omitempty"`
	Service             string                      `json:"service,omitempty"`
	Token               string                      `json:"token,omitempty"`
	AccessToken         string                      `json:"access_token,omitempty"`
	AccessCode          string                      `json:"access_code,omitempty"`
	Account             string                      `json:"account,omitempty"`
	Password            string                      `json:"password,omitempty"`
	NewPassword         string                      `json:"new_password,omitempty"`
	Avatar              string                      `json:"avatar,omitempty"`
	Name                string                      `json:"name,omitempty"`
	Nickname            string                      `json:"nickname,omitempty"`
	LastLoginTime       string                      `json:"last_login_time,omitempty"`
	Code                string                      `json:"code,omitempty"`
	ImageFileName       string                      `json:"image_file_name,omitempty"`
	ImageId             string                      `json:"image_id,omitempty"`
	Uri                 *UriObject                  `json:"uri,omitempty"`
	RelativePath        string                      `json:"relative_path,omitempty"`
	Platform            string                      `json:"platform,omitempty"`
	Width               int32                       `json:"width,omitempty"`
	Height              int32                       `json:"height,omitempty"`
	Skip                int32                       `json:"skip,omitempty"`
	Limit               int32                       `json:"limit,omitempty"`
	PatientId           string                      `json:"patient_id,omitempty"`
	VisitDate           string                      `json:"visit_date,omitempty"`
	DateId              string                      `json:"date_id,omitempty"`
	Note                string                      `json:"note,omitempty"`
	KeywordList         []string                    `json:"keyword_list,omitempty"`
	PatientKeywordList  []string                    `json:"patient_keyword_list,omitempty"`
	DiseaseKeywordList  []string                    `json:"disease_keyword_list,omitempty"`
	DateKeywordList     []string                    `json:"date_keyword_list,omitempty"`
	List                []string                    `json:"list,omitempty"`
	Default             string                      `json:"default,omitempty"`
	Count               int32                       `json:"count,omitempty"`
	PatientList         []PatientObject             `json:"patient_list,omitempty"`
	Patient             *PatientObject              `json:"patient,omitempty"`
	PatientMap          map[string]*PatientObject   `json:"patient_map,omitempty"`
	DiagnosisList       []DiagnosisObject           `json:"diagnosis_list,omitempty"`
	Diagnosis           *DiagnosisObject            `json:"diagnosis,omitempty"`
	DiagnosisMap        map[string]*DiagnosisObject `json:"diagnosis_map,omitempty"`
	DateList            []DateObject                `json:"date_list,omitempty"`
	Date                *DateObject                 `json:"date,omitempty"`
	DateMap             map[string]*DateObject      `json:"date_map,omitempty"`
	ImageList           []ImageObject               `json:"image_list,omitempty"`
	Image               *ImageObject                `json:"image,omitempty"`
	FileServerUsage     *DiskUsageObject            `json:"disk_usage,omitempty"`
	HospitalName        string                      `json:"hospital_name,omitempty"`
	HospitalId          string                      `json:"hospital_id,omitempty"`
	HospitalCreatedBy   string                      `json:"hospital_created_by,omitempty"`
	Address             string                      `json:"address,omitempty"`
	UserId              string                      `json:"user_id,omitempty"`
	Pid                 string                      `json:"pid,omitempty"`
	DiagnosisId         string                      `json:"diagnosis_id,omitempty"`
	Status              string                      `json:"status,omitempty"`
	DiagnosisNumber     int32                       `json:"diagnosis_number,omitempty"`
	DateNumber          int32                       `json:"date_number,omitempty"`
	ImageNumber         int32                       `json:"image_number,omitempty"`
	PatientNumber       int32                       `json:"patient_number,omitempty"`
	Ethnicity           string                      `json:"ethnicity,omitempty"`
	Country             string                      `json:"country,omitempty"`
	Skin                string                      `json:"skin,omitempty"`
	Disease             string                      `json:"disease,omitempty"`
	Location            string                      `json:"location,omitempty"`
	Gender              string                      `json:"gender,omitempty"`
	Tag                 string                      `json:"tag,omitempty"`
	TagList             []string                    `json:"tag_list,omitempty"`
	Ip                  string                      `json:"ip,omitempty"`
	Port                string                      `json:"port,omitempty"`
	SystemType          string                      `json:"system_type,omitempty"`
	SecretToken         string                      `json:"secret_token,omitempty"`
	MateToken           string                      `json:"mate_token,omitempty"`
	ClientToken         string                      `json:"client_token,omitempty"`
	OpponentClientToken string                      `json:"opponent_client_token,omitempty"`
	ClientType          string                      `json:"client_type,omitempty"`
	Order               string                      `json:"order,omitempty"`
	OrderBy             string                      `json:"order_by,omitempty"`
	Sdp                 map[string]interface{}      `json:"sdp,omitempty"`
	Candidate           map[string]interface{}      `json:"candidate,omitempty"`
	Host                string                      `json:"host,omitempty"`
	LiveId              string                      `json:"live_id,omitempty"`
	Three               *ThreeObject                `json:"three,omitempty"`
}

func (p *Payload) Debug(prefix string) {
	fmt.Fprintf(os.Stderr, "%v: \n", prefix)
	if j, err2 := json.MarshalIndent(p, "", " "); err2 != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err2)
	} else {
		fmt.Fprintf(os.Stderr, "%v\n", string(j))
	}
}

type UriObject struct {
	Avatar    string `json:"avatar,omitempty"`
	Image     string `json:"image,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

type FileVars struct {
	Method         string
	Id             string
	AccessToken    string
	File           string
	Separators     []string
	StaticFilePath string
	Date           time.Time
	R              *http.Request
}

type PatientObject struct {
	Pid                     string            `json:"pid,omitempty"`
	PatientId               string            `json:"patient_id,omitempty"`
	Ethnicity               string            `json:"ethnicity,omitempty"`
	FirstName               string            `json:"first_name,omitempty"`
	MiddleName              string            `json:"middle_name,omitempty"`
	LastName                string            `json:"last_name,omitempty"`
	Name                    string            `json:"name,omitempty"`
	Address                 string            `json:"address,omitempty"`
	Country                 string            `json:"country,omitempty"`
	Phone                   string            `json:"phone,omitempty"`
	Gender                  string            `json:"gender,omitempty"`
	BirthDate               string            `json:"birth_date,omitempty"`
	VisitDate               string            `json:"visit_date,omitempty"`
	HospitalId              string            `json:"hospital_id,omitempty"`
	UserId                  string            `json:"user_id,omitempty"`
	Sequence                int32             `json:"sequence,omitempty"`
	LastDate                *DateObject       `json:"last_date,omitempty"`
	LastDiagnosis           *DiagnosisObject  `json:"last_diagnosis,omitempty"`
	DiagnosisList           []DiagnosisObject `json:"diagnosis_list,omitempty"`
	DateList                []DateObject      `json:"date_list,omitempty"`
	Description             string            `json:"description,omitempty"`
	SkinType                string            `json:"skin_type,omitempty"`
	SkinCancerHistory       string            `json:"skin_cancer_history,omitempty"`
	FamilySkinCancerHistory string            `json:"family_skin_cancer_history,omitempty"`
}
type DiagnosisObject struct {
	DiagnosisId  string         `json:"diagnosis_id,omitempty"`
	Type         string         `json:"type,omitempty"`
	LocationList []string       `json:"location_list,omitempty"`
	TagList      []string       `json:"tag_list,omitempty"`
	Date         string         `json:"date,omitempty"`
	PatientId    string         `json:"patient_id,omitempty"`
	HospitalId   string         `json:"hospital_id,omitempty"`
	UserId       string         `json:"user_id,omitempty"`
	Patient      *PatientObject `json:"patient,omitempty"`
	DateList     []DateObject   `json:"date_list,omitempty"`
	LastDate     *DateObject    `json:"last_date,omitempty"`
	DateNumber   int32          `json:"date_number,omitempty"`
	Treatment    string         `json:"treatment,omitempty"`
	Note         string         `json:"note,omitempty"`
}

type DateObject struct {
	DateId           string           `json:"date_id,omitempty"`
	Type             string           `json:"type,omitempty"`
	TypeList         []string         `json:"type_list,omitempty"`
	LocationList     []string         `json:"location_list,omitempty"`
	TagList          []string         `json:"tag_list,omitempty"`
	TagListUpdated   bool             `json:"tag_list_updated,omitempty"`
	Note             string           `json:"note,omitempty"`
	NoteUpdated      bool             `json:"note_updated,omitempty"`
	VisitDate        string           `json:"visit_date,omitempty"`
	DiagnosisId      string           `json:"diagnosis_id,omitempty"`
	PatientId        string           `json:"patient_id,omitempty"`
	HospitalId       string           `json:"hospital_id,omitempty"`
	UserId           string           `json:"user_id,omitempty"`
	Treatment        string           `json:"treatment,omitempty"`
	TreatmentUpdated bool             `json:"treatment_updated,omitempty"`
	TreatmentList    []string         `json:"treatment_list,omitempty"`
	Patient          *PatientObject   `json:"patient,omitempty"`
	Diagnosis        *DiagnosisObject `json:"diagnosis,omitempty"`
	ImageList        []ImageObject    `json:"image_list,omitempty"`
	ImageNumber      int32            `json:"image_number,omitempty"`
}
type ImageObject struct {
	ImageId          string           `json:"image_id,omitempty"`
	LiveId           string           `json:"live_id,omitempty"`
	FileName         string           `json:"file_name,omitempty"`
	Width            int32            `json:"width,omitempty"`
	Height           int32            `json:"height,omitempty"`
	Uri              *UriObject       `json:"uri,omitempty"`
	Patient          *PatientObject   `json:"patient,omitempty"`
	Diagnosis        *DiagnosisObject `json:"diagnosis,omitempty"`
	Date             *DateObject      `json:"date,omitempty"`
	LocationList     []string         `json:"location_list,omitempty"`
	TagList          []string         `json:"tag_list,omitempty"`
	TagListUpdated   bool             `json:"tag_list_updated,omitempty"`
	UploadTime       string           `json:"upload_time,omitempty"`
	Treatment        string           `json:"treatment,omitempty"`
	TreatmentUpdated bool             `json:"treatment_updated,omitempty"`
	Type             string           `json:"type,omitempty"`
	TypeList         []string         `json:"type_list,omitempty"`
	VisitDate        string           `json:"visit_date,omitempty"`
	Note             string           `json:"note,omitempty"`
	NoteUpdated      bool             `json:"note_updated,omitempty"`
	Three            *ThreeObject     `json:"three,omitempty"`
}

type ThreeObject struct {
	MeshNumber int32        `json:"mesh_number,omitempty"`
	MeshList   []MeshObject `json:"mesh_list,omitempty"`
}

type MeshObject struct {
	IntersectionObjectName string  `json:"intersection_object_name,omitempty"`
	LocationName           string  `json:"location_name,omitempty"`
	IntersectionPoint      Vector3 `json:"intersection_point,omitempty"`
	Orientation            Vector3 `json:"orientation,omitempty"`
	RayDirection           Vector3 `json:"ray_direction,omitempty"`
}

type IntersectionObject struct {
	Name  string  `json:"name,omitempty"`
	Point Vector3 `json:"point,omitempty"`
}

type Vector3 struct {
	X float32 `json:"x,omitempty"`
	Y float32 `json:"y,omitempty"`
	Z float32 `json:"z,omitempty"`
}

type DiskUsageObject struct {
	Total       uint64  `json:"total,omitempty"`
	Used        uint64  `json:"used,omitempty"`
	Available   uint64  `json:"available,omitempty"`
	UsedPercent float64 `json:"used_percent,omitempty"`
}
