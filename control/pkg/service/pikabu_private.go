package service

import (
	"context"
	"fmt"
)

type PikabuPrivate struct {
}

func (s *PikabuPrivate) Service(ctx context.Context, req Payload) (res Payload, err error) {
	//TODO: 인증토큰 체크 로직
	accessToken := SecretTokenObject{
		Token: req.AccessToken,
	}
	var user UserObject
	if user, err = accessToken.Authenticate(); err != nil {
		return
	}
	switch req.Service {
	case "GetUserInformation":
		res, err = s.GetUserInformation(ctx, req, user)
	/*
	case "UpdateUserInformation":
		res, err = s.UpdateUserInformation(ctx, req, user)
	case "UpdateAccessToken":
		res, err = s.UpdateAccessToken(ctx, req, user)
	case "PrepareAvatar":
		res, err = s.PrepareAvatar(ctx, req, user)
	case "GetAvatarUri":
		res, err = s.GetAvatarUri(ctx, req, user)
	case "CreateHospital":
		res, err = s.CreateHospital(ctx, req, user)
	case "SignUpComplete":
		res, err = s.SignUpComplete(ctx, req, user)
	case "GetUserInformation":
		res, err = s.GetUserInformation(ctx, req, user)
	case "ListEthnicity":
		res, err = s.ListEthnicity(ctx, req, user)
	case "AddEthnicity":
		res, err = s.AddEthnicity(ctx, req, user)
	case "DelEthnicity":
		res, err = s.DelEthnicity(ctx, req, user)
	case "SetDefaultEthnicity":
		res, err = s.SetDefaultEthnicity(ctx, req, user)
	case "ListCountry":
		res, err = s.ListCountry(ctx, req, user)
	case "AddCountry":
		res, err = s.AddCountry(ctx, req, user)
	case "DelCountry":
		res, err = s.DelCountry(ctx, req, user)
	case "SetDefaultCountry":
		res, err = s.SetDefaultCountry(ctx, req, user)
	case "ListSkin":
		res, err = s.ListSkin(ctx, req, user)
	case "AddSkin":
		res, err = s.AddSkin(ctx, req, user)
	case "DelSkin":
		res, err = s.DelSkin(ctx, req, user)
	case "SetDefaultSkin":
		res, err = s.SetDefaultSkin(ctx, req, user)
	case "ListDisease":
		res, err = s.ListDisease(ctx, req, user)
	case "AddDisease":
		res, err = s.AddDisease(ctx, req, user)
	case "DelDisease":
		res, err = s.DelDisease(ctx, req, user)
	case "SetDefaultDisease":
		res, err = s.SetDefaultDisease(ctx, req, user)
	case "ListLocation":
		res, err = s.ListLocation(ctx, req, user)
	case "AddLocation":
		res, err = s.AddLocation(ctx, req, user)
	case "DelLocation":
		res, err = s.DelLocation(ctx, req, user)
	case "SetDefaultLocation":
		res, err = s.SetDefaultLocation(ctx, req, user)
	case "ListGender":
		res, err = s.ListGender(ctx, req, user)
	case "AddGender":
		res, err = s.AddGender(ctx, req, user)
	case "DelGender":
		res, err = s.DelGender(ctx, req, user)
	case "SetDefaultGender":
		res, err = s.SetDefaultGender(ctx, req, user)
	case "ListTag":
		res, err = s.ListTag(ctx, req, user)
	case "AddTag":
		res, err = s.AddTag(ctx, req, user)
	case "AddTags":
		res, err = s.AddTags(ctx, req, user)
	case "DelTag":
		res, err = s.DelTag(ctx, req, user)
	case "SetDefaultTag":
		res, err = s.SetDefaultTag(ctx, req, user)
	case "GetWebsocketHost":
		res, err = s.GetWebsocketHost(ctx, req, user)
	case "SignOut":
	*/

	default:
		err = fmt.Errorf("unknown service '%v' in category: '%v'", req.Service, req.Category)
	}
	return
}

func (s *PikabuPrivate) GetUserInformation(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	uri := do.Avatar.GetFileUri(do.SecretToken.Token)
	res = Payload{
		AccessToken:       req.AccessToken,
		Account:           do.Login.Account,
		SecretToken:       do.Login.Password,
		UserId:            do.Id.Hex(),
		HospitalId:        do.Relation.HospitalId.Hex(),
		Name:              do.Name,
		Nickname:          do.Nickname,
	}
	if len(uri) > 0 {
		res.Uri = &UriObject{
			Avatar: uri,
		}
	}
	return
}
/*
func (s *PikabuPrivate) GetWebsocketHost(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Host: GetConfigWebsocketHosts(),
	}
	return
}
func (s *PikabuPrivate) SetDefaultTag(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.SetDefaultCustomTag(req.Tag); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) DelTag(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.ExcludeCustomTag(req.Tag); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) AddTag(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.AddCustomTag(req.Tag); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) AddTags(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.AddCustomTags(req.TagList); err != nil {
		return
	}

	return
}
func (s *PikabuPrivate) ListTag(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	res.Default, res.List, _ = do.CustomConfig.Tag.Apply(GetDefaultTag(), GetTagList())
	return
}

func (s *PikabuPrivate) SetDefaultGender(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.SetDefaultCustomGender(req.Gender); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) DelGender(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.ExcludeCustomGender(req.Gender); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) AddGender(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.AddCustomGender(req.Gender); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) ListGender(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	res.Default, res.List, _ = do.CustomConfig.Gender.Apply(GetDefaultGender(), GetGenderList())
	return
}

func (s *PikabuPrivate) SetDefaultCountry(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.SetDefaultCustomCountry(req.Country); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) DelCountry(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.ExcludeCustomCountry(req.Country); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) AddCountry(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.AddCustomCountry(req.Country); err != nil {
		return
	}
	return
}

func (s *PikabuPrivate) SetDefaultDisease(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.SetDefaultCustomDisease(req.Disease); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) DelDisease(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.ExcludeCustomDisease(req.Disease); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) AddDisease(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.AddCustomDisease(req.Disease); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) SetDefaultLocation(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.SetDefaultCustomLocation(req.Location); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) DelLocation(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.ExcludeCustomLocation(req.Location); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) AddLocation(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.AddCustomLocation(req.Location); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) SetDefaultSkin(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.SetDefaultCustomSkin(req.Skin); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) DelSkin(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.ExcludeCustomSkin(req.Skin); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) AddSkin(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.AddCustomSkin(req.Skin); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) SetDefaultEthnicity(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.SetDefaultCustomEthnicity(req.Ethnicity); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) DelEthnicity(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.ExcludeCustomEthnicity(req.Ethnicity); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) AddEthnicity(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.AddCustomEthnicity(req.Ethnicity); err != nil {
		return
	}
	return
}
func (s *PikabuPrivate) ListLocation(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	res.Default, res.List, _ = do.CustomConfig.Location.Apply(GetDefaultLocation(), GetLocationList())
	return
}
func (s *PikabuPrivate) ListDisease(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	res.Default, res.List, _ = do.CustomConfig.Disease.Apply(GetDefaultDisease(), GetDiseaseList())
	return
}
func (s *PikabuPrivate) ListSkin(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	res.Default, res.List, _ = do.CustomConfig.Skin.Apply(GetDefaultSkin(), GetSkinList())
	return
}
func (s *PikabuPrivate) ListCountry(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	res.Default, res.List, _ = do.CustomConfig.Country.Apply(GetDefaultCountry(), GetCountryList())
	return
}
func (s *PikabuPrivate) ListEthnicity(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {

	res = Payload{
		Account: do.Login.Account,
	}
	res.Default, res.List, _ = do.CustomConfig.Ethnicity.Apply(GetDefaultEthnicity(), GetEthnicityList())

	return
}
func (s *PikabuPrivate) GetUserInformation(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	uri := do.Avatar.GetFileUri(do.SecretToken.Token)

	req.HospitalId = do.Relation.HospitalId.Hex()
	if res, err = s.GetHospital(ctx, req, do); err != nil {
		return
	}
	hospitalName := res.HospitalName
	hospitalCreatedBy := res.HospitalCreatedBy
	res = Payload{
		AccessToken:       req.AccessToken,
		Account:           do.Login.Account,
		SecretToken:       do.Login.Password,
		UserId:            do.Id.Hex(),
		HospitalId:        do.Relation.HospitalId.Hex(),
		HospitalName:      hospitalName,
		HospitalCreatedBy: hospitalCreatedBy,
		Name:              do.Name,
		Nickname:          do.Nickname,
	}
	if len(uri) > 0 {
		res.Uri = &UriObject{
			Avatar: uri,
		}
	}
	return
}
func (s *PikabuPrivate) SignUpComplete(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	if len(req.HospitalName) == 0 {
		err = errors.New("'hospital_name' is mandatory")
		return
	}
	if len(req.Name) == 0 {
		err = errors.New("'name' is mandatory")
		return
	}
	if len(req.Nickname) == 0 {
		err = errors.New("'nickname' is mandatory")
		return
	}
	if res, err = s.CreateHospital(ctx, req, do); err != nil {
		return
	}
	req.HospitalId = res.HospitalId

	return s.UpdateUserInformation(ctx, req, do)
}
func (s *PikabuPrivate) GetHospital(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	if len(req.HospitalId) == 0 {
		err = errors.New("'hospital_id' is mandatory")
		return
	}
	if bson.IsObjectIdHex(req.HospitalId) == false {
		err = errors.New("'hospital_id' is invalid")
		return
	}
	ho := HospitalObject{
		Id: bson.ObjectIdHex(req.HospitalId),
	}
	var ro HospitalObject
	if ro, err = ho.Read(); err != nil {
		return
	}
	res = Payload{
		HospitalId:        ro.Id.Hex(),
		HospitalName:      ro.Name,
		HospitalCreatedBy: ro.CreatedBy.Hex(),
	}
	return
}
func (s *PikabuPrivate) CreateHospital(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	if len(req.HospitalName) == 0 {
		err = errors.New("'hospital_name' is mandatory")
		return
	}
	ho := HospitalObject{
		CreatedBy: do.Id,
		Name:      req.HospitalName,
	}
	if err = ho.Create(); err != nil {
		return
	}
	res = Payload{
		HospitalId: ho.Id.Hex(),
	}
	return
}
func (s *PikabuPrivate) UpdateUserInformation(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	do.Name = req.Name
	do.Nickname = req.Nickname
	if len(req.HospitalId) > 0 && bson.IsObjectIdHex(req.HospitalId) == true {
		do.Relation.HospitalId = bson.ObjectIdHex(req.HospitalId)
	}
	if do.Status == "information" {
		do.Status = "active"
	}
	if err = do.Update(); err != nil {
		return
	}
	res = Payload{
		Account: do.Login.Account,
	}
	if len(req.HospitalId) > 0 {
		res.HospitalId = req.HospitalId
	}
	return
}
func (s *PikabuPrivate) UpdateAccessToken(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	do.SecretToken.Refresh()
	if err = do.Update(); err != nil {
		return
	}
	res = Payload{
		Account: do.Login.Account,
	}
	return
}

// TODO: 아바타 파일 저장하는 순서
// 1. PrepareAvatar 를 콜해서 저장할 경로를 얻어온다.
// 2. 해당 경로에 POST 한다.
func (s *PikabuPrivate) PrepareAvatar(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	if len(req.Avatar) == 0 {
		err = errors.New("'avatar' is mandatory")
		return
	}
	do.Avatar.PrepareAvatarPath(do.Id.Hex(), req.Avatar)
	if err = do.Update(); err != nil {
		return
	}
	uri := do.Avatar.GetFileUri(do.SecretToken.Token)
	res = Payload{
		Account: do.Login.Account,
		Avatar:  uri,
	}
	return
}
func (s *PikabuPrivate) GetAvatarUri(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	uri := do.Avatar.GetFileUri(do.SecretToken.Token)
	res = Payload{
		Account: do.Login.Account,
		Uri: &UriObject{
			Avatar: uri,
		},
	}
	return
}
*/