package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/azuvicenna/selaras-services/bff-mobile/internal/response"
)

type DeviceInfo struct {
	DeviceID   string `json:"device_id"`
	Platform   string `json:"platform"` // ios, android
	AppVersion string `json:"app_version"`
}

const DeviceContextKey contextKey = "mobile_device_info"

func DeviceMiddleware(minVersion, updateURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			deviceID := r.Header.Get("X-Device-Id")
			platform := strings.ToLower(r.Header.Get("X-Platform"))
			appVersion := r.Header.Get("X-App-Version")

			if appVersion != "" && minVersion != "" {
				if isOutdated(appVersion, minVersion) {
					response.WriteError(w, http.StatusUpgradeRequired, "App version is outdated. Please update your application.", map[string]string{
						"min_version": minVersion,
						"current":     appVersion,
						"update_url":  updateURL,
					})
					return
				}
			}

			devInfo := DeviceInfo{
				DeviceID:   deviceID,
				Platform:   platform,
				AppVersion: appVersion,
			}

			ctx := context.WithValue(r.Context(), DeviceContextKey, devInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetDeviceInfo(ctx context.Context) DeviceInfo {
	if val, ok := ctx.Value(DeviceContextKey).(DeviceInfo); ok {
		return val
	}
	return DeviceInfo{}
}

func isOutdated(current, min string) bool {
	currParts := strings.Split(current, ".")
	minParts := strings.Split(min, ".")

	for i := 0; i < len(minParts) && i < len(currParts); i++ {
		cNum, _ := strconv.Atoi(currParts[i])
		mNum, _ := strconv.Atoi(minParts[i])
		if cNum < mNum {
			return true
		}
		if cNum > mNum {
			return false
		}
	}
	return len(currParts) < len(minParts)
}
