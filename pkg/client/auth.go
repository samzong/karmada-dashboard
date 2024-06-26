package client

import (
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"net/http"
	"strings"
)

const (
	// authorizationHeader is the default authorization header name.
	authorizationHeader = "Authorization"
	// authorizationTokenPrefix is the default bearer token prefix.
	authorizationTokenPrefix = "Bearer "
)

func karmadaConfigFromRequest(request *http.Request) (*rest.Config, error) {
	authInfo, err := buildAuthInfo(request)
	if err != nil {
		return nil, err
	}

	return buildConfigFromAuthInfo(authInfo)
}

func buildConfigFromAuthInfo(authInfo *clientcmdapi.AuthInfo) (*rest.Config, error) {
	cmdCfg := clientcmdapi.NewConfig()

	cmdCfg.Clusters[DefaultCmdConfigName] = &clientcmdapi.Cluster{
		Server:                   karmadaRestConfig.Host,
		CertificateAuthority:     karmadaRestConfig.TLSClientConfig.CAFile,
		CertificateAuthorityData: karmadaRestConfig.TLSClientConfig.CAData,
		InsecureSkipTLSVerify:    karmadaRestConfig.TLSClientConfig.Insecure,
	}

	cmdCfg.AuthInfos[DefaultCmdConfigName] = authInfo

	cmdCfg.Contexts[DefaultCmdConfigName] = &clientcmdapi.Context{
		Cluster:  DefaultCmdConfigName,
		AuthInfo: DefaultCmdConfigName,
	}

	cmdCfg.CurrentContext = DefaultCmdConfigName

	return clientcmd.NewDefaultClientConfig(
		*cmdCfg,
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
}

func buildAuthInfo(request *http.Request) (*clientcmdapi.AuthInfo, error) {
	if !HasAuthorizationHeader(request) {
		return nil, k8serrors.NewUnauthorized("MSG_LOGIN_UNAUTHORIZED_ERROR")
	}

	token := GetBearerToken(request)
	authInfo := &clientcmdapi.AuthInfo{
		Token:                token,
		ImpersonateUserExtra: make(map[string][]string),
	}

	handleImpersonation(authInfo, request)
	return authInfo, nil
}

func HasAuthorizationHeader(req *http.Request) bool {
	header := req.Header.Get(authorizationHeader)
	if len(header) == 0 {
		return false
	}

	token := extractBearerToken(header)
	return strings.HasPrefix(header, authorizationTokenPrefix) && len(token) > 0
}

func GetBearerToken(req *http.Request) string {
	header := req.Header.Get(authorizationHeader)
	return extractBearerToken(header)
}

func SetAuthorizationHeader(req *http.Request, token string) {
	req.Header.Set(authorizationHeader, authorizationTokenPrefix+token)
}

func extractBearerToken(header string) string {
	return strings.TrimPrefix(header, authorizationTokenPrefix)
}

func handleImpersonation(authInfo *clientcmdapi.AuthInfo, request *http.Request) {
	user := request.Header.Get(ImpersonateUserHeader)
	groups := request.Header[ImpersonateGroupHeader]

	if len(user) == 0 {
		return
	}

	// Impersonate user
	authInfo.Impersonate = user

	// Impersonate groups if available
	if len(groups) > 0 {
		authInfo.ImpersonateGroups = groups
	}

	// Add extra impersonation fields if available
	for name, values := range request.Header {
		if strings.HasPrefix(name, ImpersonateUserExtraHeader) {
			extraName := strings.TrimPrefix(name, ImpersonateUserExtraHeader)
			authInfo.ImpersonateUserExtra[extraName] = values
		}
	}
}
