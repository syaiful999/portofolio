package middleware

import (
	"context"
	"errors"

	cc "moyo-gateway-service/cache"
	"moyo-gateway-service/pkg/entities"
	"moyo-gateway-service/utils"
	"net/http"

	js "encoding/json"
)

var (
	ErrUnauthorized = errors.New("unauthorized. access denied")
)

func Authorization(ctx context.Context, token string, r *http.Request, cache cc.Cache) error {
	sub, ok := utils.GetSub(r.Header.Get("Authorization"))
	if !ok {
		return errors.New("invalid token.")
	}

	data, _, err := cache.Get(sub)
	if err != nil {
		return errors.New("token expired.")
	}

	var resRbac entities.RBAC
	if err := js.Unmarshal(data.([]byte), &resRbac); err != nil {
		return err
	}

	isAuthorized := false

	if len(resRbac.Rbac.Application) > 0 {
		for _, v := range resRbac.Rbac.Application {
			if len(v.AccessControl.Permission) > 0 {
				for _, x := range v.AccessControl.Permission {
					isAuthorized = checkPermission(x.Permission, r)
					if isAuthorized {
						break
					}
				}
			}
			if isAuthorized {
				break
			}
		}
	}
	if !isAuthorized {
		return ErrUnauthorized
	}

	return nil
}

func checkPermission(p []entities.Permission, r *http.Request) bool {
	for _, v := range p {
		if len(v.Permission) > 0 {
			return checkPermission(v.Permission, r)
		}
		if len(v.Operations) > 0 {
			for _, x := range v.Operations {
				if x.Path == string(r.URL.Path) {
					return true
				}
			}
		}
	}
	return false
}
