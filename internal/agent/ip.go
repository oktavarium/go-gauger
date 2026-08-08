package agent

import (
	"fmt"
	"net"
)

func getLocalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("error on getting interfaces: %w", err)
	}

	for _, iface := range ifaces {
		// Проверяем, что интерфейс активен и не является loopback
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue // Интерфейс не активен или является loopback
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue // Пропускаем интерфейс, если возникла ошибка при получении адресов
		}

		for _, addr := range addrs {
			// Проверяем, что адрес является сетевым адресом IPv4
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet.IP.To4() == nil {
				continue // Это не IPv4
			}

			// Исключаем локальный адрес обратной связи
			if ipNet.IP.IsLoopback() {
				continue
			}

			return ipNet.IP.String(), nil // Возвращаем первый подходящий адрес
		}
	}

	return "", fmt.Errorf("local IP address not found")
}
