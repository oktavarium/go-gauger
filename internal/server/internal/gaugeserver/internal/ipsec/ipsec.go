package ipsec

import (
	"fmt"
	"net"
	"net/http"

	"github.com/oktavarium/go-gauger/internal/server/internal/logger"
	"go.uber.org/zap"
)

type securityProvider struct {
	trustedSubnet *net.IPNet
}

func NewIPSec(subnet string) (securityProvider, error) {
	if len(subnet) == 0 {
		return securityProvider{
			trustedSubnet: nil,
		}, nil
	}

	_, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		return securityProvider{}, fmt.Errorf("error parsing subnet addr: %w", err)
	}

	return securityProvider{
		trustedSubnet: ipnet,
	}, nil
}

// IPSecMiddleware - метод проверки доверенных сетей клиентов
func (sec securityProvider) IPSecMiddleware(next http.Handler) http.Handler {
	hf := func(w http.ResponseWriter, r *http.Request) {
		if sec.trustedSubnet != nil {
			clientIP := r.Header.Get("X-Real-IP")
			if len(clientIP) != 0 {
				ipaddr, _, err := net.ParseCIDR(clientIP)
				if err != nil {
					logger.Logger().Error("error",
						zap.String("func", "IPSecMiddleware"),
						zap.Error(err))

					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if !sec.trustedSubnet.Contains(ipaddr) {
					w.WriteHeader(http.StatusForbidden)
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(hf)
}
