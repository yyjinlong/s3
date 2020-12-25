// Copyright @ 2020 OPS Inc.
//
// Author: Jinlong Yang
//

package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"s3/ops/db"
	"s3/ops/objects"
)

const (
	HEADER_AUTHORIZATION = "authorization"
	AUTH_VERSION         = "s3-auth-v1"
	S3_PREFIX            = "x-s3-"
)

func CheckSign(r *S3Request) (bool, *db.AuthMeta) {
	headers := r.Headers
	auth := headers[HEADER_AUTHORIZATION]

	//0 s3-auth-v1
	//1 {accessKeyId}
	//2 {timestamp}
	//3 {expirationPeriodInSeconds}
	//4 {signedHeaders}
	//5 {signature}
	authList := strings.Split(auth, "/")
	if len(authList) != 6 {
		log.Warnf("[logid:%s] authorization format error.", r.RequestID)
		return false, nil
	}

	if authList[0] != AUTH_VERSION {
		log.Warnf("[logid:%s] authorization version is not match.", r.RequestID)
		return false, nil
	}

	timestamp := authList[2]
	createTime, err := time.Parse("2006-01-02T15:04:05Z", timestamp)
	if err != nil {
		log.Warnf("[logid:%s] parse authorization time error: %s", r.RequestID, err)
		return false, nil
	}

	ttl := authList[3]
	duration, _ := strconv.Atoi(ttl)
	expireTime := createTime.Add(time.Second * time.Duration(duration))

	if (time.Now().UTC()).After(expireTime) {
		log.Warnf("[logid:%s] authorization is expired.", r.RequestID)
		return false, nil
	}

	ak := authList[1]
	meta, err := objects.GetAuthMetaByAk(ak)
	if err != nil {
		return false, nil
	}
	sk := meta.SK

	// step1, 获取认证字符串前缀
	urlPrefix := getAuthPrefix(ak, timestamp, duration)

	// step2, 获取Signing key
	signKey := getSignKey(sk, urlPrefix)

	// step3, 获取CanonicalRequest
	canonicalRequest := getCanonicalRequest(r.HTTPMethod, r.URL, r.Query, headers)

	// step4, 获取signature
	signature := getSignature(signKey, canonicalRequest)

	sign := authList[5]
	if sign != signature {
		log.Warnf("[logid:%s] sign:%s is not match.", r.RequestID, signature)
		return false, nil
	}
	return true, meta
}

func getAuthPrefix(ak, timestamp string, ttl int) string {
	// 格式: s3-auth-v1/{ak}/{timestamp}/ttl/
	return fmt.Sprintf("%s/%s/%s/%d", AUTH_VERSION, ak, timestamp, ttl)
}

func getSignKey(sk, authPrefix string) string {
	// 使用SHA256算法, 对 [SK] 和 [认证字符串前缀部分] 进行加密
	return hmacSha256Hex(sk, authPrefix)
}

func getCanonicalQuery(param string) string {
	// 获取规范化的url参数
	if param == "" {
		return ""
	}

	result := make([]string, 0)
	paramList := strings.Split(param, "&")
	for _, item := range paramList {
		itemList := strings.Split(item, "=")
		k := itemList[0]
		if strings.ToLower(k) == HEADER_AUTHORIZATION {
			continue
		}
		if len(itemList) == 1 {
			result = append(result, fmt.Sprintf("%s=", url.QueryEscape(k)))
			continue
		}
		v := itemList[1]
		result = append(result, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
	}
	sort.Strings(result)
	return strings.Join(result, "&")
}

func getCanonicalHeaders(headers map[string]string) string {
	// 获取规范化的headers
	result := make([]string, 0)
	for k, v := range headers {
		headKey := strings.ToLower(k)
		if strings.HasPrefix(headKey, S3_PREFIX) {
			headVal := strings.TrimSpace(v)
			result = append(result, fmt.Sprintf("%s:%s", url.QueryEscape(headKey), url.QueryEscape(headVal)))
		}
	}
	sort.Strings(result)
	return strings.Join(result, "\n")
}

func getCanonicalRequest(method, uri, param string, headers map[string]string) string {
	// 格式: HTTP Method + "\n" + CanonicalURI + "\n" + CanonicalQueryString + "\n" + CanonicalHeaders
	result := make([]string, 0)
	result = append(result, method)
	result = append(result, uri)

	canonicalQuery := getCanonicalQuery(param)
	result = append(result, canonicalQuery)

	canonicalHeadrs := getCanonicalHeaders(headers)
	result = append(result, canonicalHeadrs)
	return strings.Join(result, "\n")
}

func getSignature(signKey, canonicalRequest string) string {
	// 使用SHA256算法, 对 [Signingkey] 和 [CanonicalRequest] 进行加密
	return hmacSha256Hex(signKey, canonicalRequest)
}

func hmacSha256Hex(key, msg string) string {
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write([]byte(msg))
	return hex.EncodeToString(hash.Sum(nil))
}
