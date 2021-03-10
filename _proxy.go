// package main

// import (
// 	"golang.org/x/net/proxy"
// )

// // NewClientFromEnv http client that leverages a SOCKS5 proxy and a DialContext
// func NewClientFromEnv() (*http.Client, error) {
// 	proxyHost := "176.113.73.97:3128" // os.Getenv("PROXY_HOST")

// 	baseDialer := &net.Dialer{
// 		Timeout:   30 * time.Second,
// 		KeepAlive: 30 * time.Second,
// 	}
// 	if proxyHost != "" {
// 		fmt.Println("Using SOCKS5 proxy for http client: " + proxyHost)
// 		dialSocksProxy, err := proxy.SOCKS5("tcp", proxyHost, nil, baseDialer)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "Error creating SOCKS5 proxy")
// 		}
// 		if contextDialer, ok := dialSocksProxy.(proxy.ContextDialer); ok {
// 			return &http.Client{
// 				Transport: &http.Transport{
// 					Proxy:                 http.ProxyFromEnvironment,
// 					DialContext:           contextDialer.DialContext,
// 					MaxIdleConns:          10,
// 					IdleConnTimeout:       60 * time.Second,
// 					TLSHandshakeTimeout:   10 * time.Second,
// 					ExpectContinueTimeout: 1 * time.Second,
// 					MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
// 				},
// 			}, nil
// 		}
// 		return nil, errors.New("Failed type assertion to DialContext")

// 		// logger.Debug("Using SOCKS5 proxy for http client",
// 		// 	zap.String("host", proxyHost),
// 		// )
// 	}
// 	return &http.Client{
// 		Transport: &http.Transport{
// 			Proxy:                 http.ProxyFromEnvironment,
// 			DialContext:           (baseDialer).DialContext,
// 			MaxIdleConns:          10,
// 			IdleConnTimeout:       60 * time.Second,
// 			TLSHandshakeTimeout:   10 * time.Second,
// 			ExpectContinueTimeout: 1 * time.Second,
// 			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
// 		},
// 	}, nil

// 	// return nil, nil
// }

// // NewClientFromEnv2 def
// func NewClientFromEnv2() (*http.Client, error) {
// 	// att := &proxy.Auth{%Login%, %Password%}
// 	var att *proxy.Auth // := nil
// 	proxyHost := "176.113.73.97:3128"
// 	dialSocks5, err := proxy.SOCKS5("tcp", proxyHost, att, proxy.Direct)
// 	if err != nil {
// 		return nil, err
// 	}

// 	transport := &http.Transport{Dial: dialSocks5.Dial}
// 	client := &http.Client{}
// 	client.Transport = transport

// 	return client, nil
// }
