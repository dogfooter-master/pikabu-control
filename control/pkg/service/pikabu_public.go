package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type PikabuPublic struct {
}

func (s *PikabuPublic) Service(ctx context.Context, req Payload) (res Payload, err error) {
	switch req.Service {
	case "SignUp":
		res, err = s.SignUp(ctx, req)
	case "VerifyCertificationCode":
		res, err = s.VerifyCertificationCode(ctx, req)
	case "SignUpPassword":
		res, err = s.SignUpPassword(ctx, req)
	case "SignIn":
		// TODO: OAuth(우선순위 낮음)
		res, err = s.SignIn(ctx, req)
	case "GetUserStatus":
		res, err = s.GetUserStatus(ctx, req)
		/*
	case "ResendCertificationCode":
		res, err = s.ResendCertificationCode(ctx, req)
	case "UpdatePassword":
		res, err = s.UpdatePassword(ctx, req)
	case "SignIn":
		// TODO: OAuth(우선순위 낮음)
		res, err = s.SignIn(ctx, req)
	case "GetCertificationCode":
		res, err = s.GetCertificationCode(ctx, req)
	case "GetUserStatus":
		res, err = s.GetUserStatus(ctx, req)
	case "GetDiskUsage":
		res, err = s.GetDiskUsage(ctx, req)
		*/
	default:
		err = fmt.Errorf("unknown service '%v' in category: '%v'", req.Service, req.Category)
	}
	return
}

func (s *PikabuPublic) GetUserStatus(ctx context.Context, req Payload) (res Payload, err error) {
	if len(req.Account) == 0 {
		err = errors.New("'account' is mandatory")
		return
	}

	do := UserObject{
		Login: LoginObject{
			Account: req.Account,
		},
	}
	// 계정으로 검색
	if ro, err2 := do.Read(); err2 != nil {
		err = fmt.Errorf("%v", err2)
		return
	} else {
		res = Payload{
			Account: ro.Login.Account,
			Status:  ro.Status,
		}
	}
	return
}
func (s *PikabuPublic) SignIn(ctx context.Context, req Payload) (res Payload, err error) {
	TimeTrack(time.Now(), GetFunctionName())
	// 로그인 정보 체크
	if len(req.Account) == 0 {
		err = errors.New("'account' is mandatory")
		return
	}
	do := UserObject{
		Login: LoginObject{
			Account: req.Account,
		},
	}
	// 계정으로 검색
	if ro, err2 := do.Read(); err2 != nil {
		err = fmt.Errorf("%v", err2)
		return
	} else {
		// 인증
		if err = ro.Login.Auth(ro.Id.Hex(), req.Password); err != nil {
			return
		}
		if ro.Status == "active" && len(ro.SecretToken.Token) > 0 {
			do.SecretToken = ro.SecretToken
		} else {
			// 인증토큰 발급
			if err = do.SecretToken.GenerateAccessToken(); err != nil {
				return
			}
		}
		// 사용자 정보 업데이트
		do.Id = ro.Id
		do.Time.Login()
		if err = do.Update(); err != nil {
			return
		}
		// TODO: 필요한 정보들 더
		res = Payload{
			Account:       do.Login.Account,
			AccessToken:   do.SecretToken.Token,
			LastLoginTime: ro.Time.LoginTime.Format(time.RFC3339),
			UserId:        do.Id.Hex(),
			Status:        do.Status,
		}
	}

	return
}
func (s *PikabuPublic) UpdatePassword(ctx context.Context, req Payload) (res Payload, err error) {
	TimeTrack(time.Now(), GetFunctionName())
	// 인자 체크
	if len(req.Account) == 0 {
		err = errors.New("'account' is mandatory")
		return
	}
	// 패스워드 체크
	if len(req.NewPassword) == 0 {
		err = errors.New("'new_password' is mandatory")
		return
	}
	do := UserObject{
		Login: LoginObject{
			Account: req.Account,
		},
	}
	// account 가 존재하는 지 검색
	if ro, err2 := do.Read(); err2 != nil {
		err = fmt.Errorf("%v", err2)
		return
	} else {
		// 패스워드 정합성 체크
		if ro.Status != "password" {
			if strings.Compare(req.NewPassword, req.Password) == 0 {
				err = errors.New("password is invalid")
				return
			}
			if err = ro.Login.Auth(ro.Id.Hex(), req.Password); err != nil {
				return
			}
		}
		do.Id = ro.Id
		// 패스워드 유효성 체크
		//if err = do.Login.ValidatePassword(req.NewPassword); err != nil {
		//	return
		//}
		// status 변경
		if ro.Status == "password" {
			do.Status = "information"
		}
		// 패스워드 업데이트
		do.Login.EncodePassword(do.Id.Hex(), req.NewPassword)
		if err = do.Update(); err != nil {
			return
		}
		res = Payload{
			Account: do.Login.Account,
		}
	}
	return
}
func (s *PikabuPublic) SignUpPassword(ctx context.Context, req Payload) (res Payload, err error) {
	TimeTrack(time.Now(), GetFunctionName())
	if len(req.Account) == 0 {
		err = errors.New("'account' is mandatory")
		return
	}
	if len(req.Password) == 0 {
		err = errors.New("'password' is mandatory")
		return
	}
	req.NewPassword = req.Password
	if res, err = s.UpdatePassword(ctx, req); err != nil {
		return
	}
	return s.SignIn(ctx, req)
}
func (s *PikabuPublic) VerifyCertificationCode(ctx context.Context, req Payload) (res Payload, err error) {
	TimeTrack(time.Now(), GetFunctionName())
	if len(req.Account) == 0 {
		err = errors.New("'account' is mandatory")
		return
	}
	if len(req.Code) == 0 {
		err = errors.New("'code' is mandatory")
		return
	}

	do := UserObject{
		Login: LoginObject{
			Account: req.Account,
		},
	}
	// 계정으로 검색
	if ro, err2 := do.Read(); err2 != nil {
		err = fmt.Errorf("%v", err2)
		return
	} else {
		// 계정 status 확인(verifying 이 아니면 에러 발생)
		if ro.Status != "verifying" {
			err = errors.New("invalid status")
			return
		}
		// Code 가 맞지않다면 Fail
		if ro.SecretToken.Token != req.Code {
			err = errors.New("fail to verify code")
			return
		}
		// Code 유효 기간이 지났다면 Fail
		if ro.SecretToken.IsExpired() {
			err = errors.New("expired")
			return
		}
		// status 를 password 상태로 변경
		do.Id = ro.Id
		do.Status = "password"
		if err = do.Update(); err != nil {
			return
		}
		res = Payload{
			Account: do.Login.Account,
		}
	}
	return
}
func (s *PikabuPublic) SignUp(ctx context.Context, req Payload) (res Payload, err error) {
	TimeTrack(time.Now(), GetFunctionName())
	// TODO: 회원가입 절차
	// 1. 이메일 입력 -> 이미 가입된 이메일인지 체크
	// 2. 이메일로 인증 메일 발송(인증 번호) -> 모바일에서는 인증 번호를 입력하는 화면으로 대기 중
	// -> 웹에서는 인증 번호를 입력하는 화면으로 대기 또는 인증 메일에서 클릭
	// 3. 인증 메일 클릭하면 성공적으로 인증됐다고 나오면서 패스워드 입력 화면 출력
	// 4. 이름, 닉네임 입력 화면 출력

	// 로그인 정보 체크
	if len(req.Account) == 0 {
		err = errors.New("'account' is mandatory")
		return
	}
	login := LoginObject{
		Account: req.Account,
	}
	if err = login.ValidateAccount(req.Account); err != nil {
		return
	}
	do := UserObject{
		Login: login,
	}
	if ro, err2 := do.Read(); err2 == nil {
		// TODO: 토큰이 발급한 지 하루가 지났으면 재발급
		// TODO: 아니라면 'verifying' 에러
		if ro.Status == "active" {
			err = errors.New("duplicated")
		} else {
			err = errors.New(ro.Status)
		}
		return
	}
	//// 사용자 정보 체크
	//if err = do.Validate(); err != nil {
	//	return
	//}
	// 사용자 생성
	do.Status = "verifying"
	// 인증 번호 생성
	do.SecretToken.GenerateCertificationNumber()
	if err = do.Create(); err != nil {
		return
	}
	// 인증 번호 전송 횟수 초기화
	do.SecretToken.Count = 0
	// 이메일 전송
	email := Email{
		To:                  do.Login.Account,
		CertificationNumber: do.SecretToken.Token,
	}
	if err = email.SendEmailCertificationNumber(); err != nil {
		return
	}
	res = Payload{
		Account: do.Login.Account,
	}
	return
}

/*
func (s *PikabuPublic) GetCertificationCode(ctx context.Context, req Payload) (res Payload, err error) {
	if len(req.Account) == 0 {
		err = errors.New("'account' is mandatory")
		return
	}

	do := UserObject{
		Login: LoginObject{
			Account: req.Account,
		},
	}
	// 계정으로 검색
	if ro, err2 := do.Read(); err2 != nil {
		err = fmt.Errorf("%v", err2)
		return
	} else {
		res = Payload{
			Account: ro.Login.Account,
			Code:    ro.SecretToken.Token,
		}
	}
	return
}
func (s *PikabuPublic) GetPlatform(ctx context.Context, req Payload) (res Payload, err error) {
	res = Payload{
		Platform: envOs,
	}
	return
}
func (s *PikabuPublic) GetSystemType(ctx context.Context, req Payload) (res Payload, err error) {
	res = Payload{
		SystemType: GetConfigSystemType(),
	}
	return
}
func (s *PikabuPublic) SignUp(ctx context.Context, req Payload) (res Payload, err error) {
	// TODO: 회원가입 절차
	// 1. 이메일 입력 -> 이미 가입된 이메일인지 체크
	// 2. 이메일로 인증 메일 발송(인증 번호) -> 모바일에서는 인증 번호를 입력하는 화면으로 대기 중
	// -> 웹에서는 인증 번호를 입력하는 화면으로 대기 또는 인증 메일에서 클릭
	// 3. 인증 메일 클릭하면 성공적으로 인증됐다고 나오면서 패스워드 입력 화면 출력
	// 4. 이름, 닉네임 입력 화면 출력

	// 로그인 정보 체크
	if len(req.Account) == 0 {
		err = errors.New("'account' is mandatory")
		return
	}
	login := LoginObject{
		Account: req.Account,
	}
	if err = login.ValidateAccount(req.Account); err != nil {
		return
	}
	do := UserObject{
		Login: login,
	}
	if ro, err2 := do.Read(); err2 == nil {
		// TODO: 토큰이 발급한 지 하루가 지났으면 재발급
		// TODO: 아니라면 'verifying' 에러
		if ro.Status == "verifying" {
			err = errors.New("verifying")
		} else {
			err = errors.New("duplicated")
		}
		return
	}
	//// 사용자 정보 체크
	//if err = do.Validate(); err != nil {
	//	return
	//}
	// 사용자 생성
	do.Status = "verifying"
	// 인증 번호 생성
	do.SecretToken.GenerateCertificationNumber()
	if err = do.Create(); err != nil {
		return
	}
	// 인증 번호 전송 횟수 초기화
	do.SecretToken.Count = 0
	// 이메일 전송
	email := Email{
		To:                  do.Login.Account,
		CertificationNumber: do.SecretToken.Token,
	}
	if err = email.SendEmailCertificationNumber(); err != nil {
		return
	}
	res = Payload{
		Account: do.Login.Account,
	}
	return
}
func (s *PikabuPublic) ResendCertificationCode(ctx context.Context, req Payload) (res Payload, err error) {
	// 인자 체크
	if len(req.Account) == 0 {
		err = errors.New("'account' is mandatory")
		return
	}
	do := UserObject{
		Login: LoginObject{
			Account: req.Account,
		},
	}
	// 계정으로 검색
	if ro, err2 := do.Read(); err2 != nil {
		err = fmt.Errorf("%v", err2)
		return
	} else {
		// status 가 verifying 이 아니면 에러
		if ro.Status != "verifying" {
			err = errors.New("invalid status")
			return
		}
		// 하루 이상 지났으면 재전송 초기화
		if ro.SecretToken.IsExpired() {
			ro.SecretToken.Count = 0
		}
		// 3회 이상 재전송은 블럭
		if ro.SecretToken.Count > 1 {
			err = errors.New("too many resend")
			return
		}
		// 인증 번호 생성
		do.SecretToken.GenerateCertificationNumber()
		// 인증 번호 전송 횟수 증가
		do.SecretToken.Count = ro.SecretToken.Count + 1

		do.Id = ro.Id
		if err = do.Update(); err != nil {
			return
		}
		// 이메일 전송
		email := Email{
			To:                  ro.Login.Account,
			ToAlias:             ro.Nickname,
			CertificationNumber: do.SecretToken.Token,
		}
		if err = email.SendEmailCertificationNumber(); err != nil {
			return
		}
		res = Payload{
			Account: do.Login.Account,
		}
	}
	return
}

func (s *PikabuPublic) GetDiskUsage(ctx context.Context, req Payload) (res Payload, err error) {
	fmt.Fprintf(os.Stdout, "GetDiskUsage\n")
	if !withoutFileServer {
		usage, err2 := FileServerUsage()
		if err2 != nil {
			err = fmt.Errorf("%v", err2)
			return
		}

		res.FileServerUsage = &usage
	} else {

		message := Payload{
			Category:    "public",
			Service:     "GetDiskUsage",
			AccessToken: req.AccessToken,
		}
		response, err2 := HttpRequestFileServer(message)
		if err2 != nil {
			err = fmt.Errorf("%v", err2)
			return
		}
		res = response
	}
	return
}
*/
