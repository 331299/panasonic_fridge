package panasonicfridge

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"log"
	"net"
	"strings"
)

func getMacAddress() string {
	interfaces, err := net.Interfaces()
	if nil != err {
		log.Println("unable to get mac address", err)
		return ""
	}
	for _, info := range interfaces {
		mac := info.HardwareAddr.String()
		if mac != "" {
			return mac
		}
	}
	return ""
}

func encodePassword(pwd string, phone string, token string) string {
	res := getMd5(pwd)
	res = getMd5(res + phone)
	if token != "" {
		res = getMd5(res + token)
	}
	return res
}

func getMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	res := hex.EncodeToString(h.Sum(nil))
	return strings.ToUpper(res)
}

func getSha512(src string) string {
	h := sha512.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}

func getSToken(deviceId string) string {
	parts := strings.SplitN(deviceId, "_", 3)
	sToken := parts[0][6:] + "_" + parts[1] + "_" + parts[0][:6]
	return getSha512(getSha512(sToken) + "_" + parts[2])
}
