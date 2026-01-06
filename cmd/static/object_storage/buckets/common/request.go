package common

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	authPkg "github.com/magaluCloud/mgccli/cmd/common/auth"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
)

type SignatureParameters struct {
	Algorithm     string
	AccessKey     string
	Credential    string
	Scope         string
	Date          string
	ShortDate     string
	PayloadHash   string
	SignedHeaders []string
}

type SignatureContext struct {
	Parameters SignatureParameters

	HTTPMethod       string
	CanonicalURI     string
	CanonicalQuery   string
	CanonicalHeaders string

	SignedHeaders string

	Signature string
}

var excludedHeaders = map[string]struct{}{
	http.CanonicalHeaderKey("Authorization"):         {},
	http.CanonicalHeaderKey("Accept-Encoding"):       {},
	http.CanonicalHeaderKey("Amz-Sdk-Invocation-Id"): {},
	http.CanonicalHeaderKey("Amz-Sdk-Request"):       {},
	http.CanonicalHeaderKey("User-Agent"):            {},
	http.CanonicalHeaderKey("X-Amzn-Trace-Id"):       {},
	http.CanonicalHeaderKey("Expect"):                {},
	http.CanonicalHeaderKey("Content-Length"):        {},
}

func SendRequest(ctx context.Context, req *http.Request, region string) (*http.Response, error) {
	auth := ctx.Value(cmdutils.CTX_AUTH_KEY).(authPkg.Auth)
	config := auth.GetConfig()

	client, err := authPkg.NewOAuthClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth client: %w", err)
	}

	httpClient := client.AuthenticatedHttpClientFromContext(ctx)
	if httpClient == nil {
		return nil, fmt.Errorf("programming error: unable to get HTTP Client from context")
	}

	accessKeyID := auth.GetAccessKeyID()
	secretAccessKey := auth.GetSecretAccessKey()

	if accessKeyID == "" || secretAccessKey == "" {
		return nil, fmt.Errorf("api-key not set, see how to set it with \"object-storage api-key -h\"")
	}

	var unsignedPayload bool
	if req.Method == http.MethodPut {
		unsignedPayload = true
	}

	req.Header.Set("Host", req.Host)

	if err := setMD5Checksum(req); err != nil {
		return nil, fmt.Errorf("unable to compute checksum of the body content: %w", err)
	}

	signingTime := time.Now().UTC()

	req.Header.Set(headerDateKey, signingTime.Format(longTimeFormat))

	payloadHash, err := setContentHeader(req, unsignedPayload)
	if err != nil {
		return nil, err
	}

	signedHeaders := getSignedHeaders(req, excludedHeaders)

	params := newSignatureParameters(accessKeyID, signingTime, payloadHash, signedHeaders, region)

	signatureCtx := NewSignatureContext(params, req)
	if err = sign(signatureCtx, secretAccessKey, region); err != nil {
		return nil, err
	}

	authorization := fmt.Sprintf(
		"%s Credential=%s, SignedHeaders=%s, Signature=%s",
		signingAlgorithm,
		signatureCtx.Parameters.Credential,
		signatureCtx.SignedHeaders,
		signatureCtx.Signature,
	)
	req.Header.Set("Authorization", authorization)

	return httpClient.Do(req)
}

func setMD5Checksum(req *http.Request) error {
	if req.Body == nil {
		return nil
	}

	if v := req.Header.Get(contentMD5Header); len(v) != 0 {
		return nil
	}
	if req.GetBody == nil {
		return fmt.Errorf("programming error: object storage operation must define a GetBody function in the request to set the MD5 Checksum")
	}

	body, err := req.GetBody()
	if err != nil {
		return err
	}

	defer body.Close()

	h := md5.New()

	_, err = io.Copy(h, body)

	if err != nil {
		return err
	}
	checksum := base64.StdEncoding.EncodeToString(h.Sum(nil))
	req.Header.Set(contentMD5Header, checksum)

	return nil
}

func setContentHeader(req *http.Request, unsigned bool) (payloadHash string, err error) {
	if unsigned {
		req.Header.Set(contentSHAKey, unsignedPayloadHeader)
		return unsignedPayloadHeader, nil
	}

	payloadHash, err = getPayloadHash(req)
	if err != nil {
		return "", err
	}
	req.Header.Set(contentSHAKey, payloadHash)
	return payloadHash, nil
}

func newSignatureParameters(
	accessKey string,
	signingTime time.Time,
	payloadHash string,
	signedHeaders []string,
	region string,
) SignatureParameters {
	shortDate := signingTime.Format(shortTimeFormat)
	scope := strings.Join([]string{
		shortDate,
		region,
		signingService,
		requestSuffix,
	}, "/")

	return SignatureParameters{
		Algorithm:     signingAlgorithm,
		AccessKey:     accessKey,
		Credential:    fmt.Sprintf("%s/%s", accessKey, scope),
		Scope:         scope,
		Date:          signingTime.Format(longTimeFormat),
		ShortDate:     shortDate,
		PayloadHash:   payloadHash,
		SignedHeaders: signedHeaders,
	}
}

func getSignedHeaders(req *http.Request, ignoredHeaders map[string]struct{}) []string {
	signedHeaders := make([]string, 0, len(req.Header))

	for k := range req.Header {
		if _, ok := ignoredHeaders[k]; ok {
			continue
		}
		signedHeaders = append(signedHeaders, strings.ToLower(k))
	}

	slices.Sort(signedHeaders)
	return signedHeaders
}

func NewSignatureContext(
	params SignatureParameters,
	req *http.Request,
) *SignatureContext {
	canonicalHeaders := buildCanonicalHeaders(req, params.SignedHeaders)

	return &SignatureContext{
		Parameters: params,

		HTTPMethod:       req.Method,
		CanonicalURI:     req.URL.EscapedPath(),
		CanonicalQuery:   req.URL.RawQuery,
		CanonicalHeaders: canonicalHeaders,

		SignedHeaders: strings.Join(params.SignedHeaders, ";"),
	}
}

func buildCanonicalHeaders(req *http.Request, signedHeaders []string) (canonicalHeaders string) {
	for _, k := range signedHeaders {
		v := req.Header.Values(k)

		line := fmt.Sprintf("%s:%s", strings.ToLower(k), strings.Join(v, ","))
		canonicalHeaders = fmt.Sprintf("%s%s\n", canonicalHeaders, line)
	}
	return
}

func sign(ctx *SignatureContext, secretKey, region string) (err error) {
	strToSign, err := buildStringToSign(ctx)
	if err != nil {
		return
	}
	signKey := deriveKey(secretKey, ctx.Parameters.ShortDate, region)
	ctx.Signature = hex.EncodeToString(HMACSHA256String(signKey, strToSign))
	return
}

func getPayloadHash(req *http.Request) (string, error) {
	if req.Body == nil {
		return emptyStringSHA256, nil
	}
	bodyReader, err := req.GetBody()
	if err != nil {
		return "", err
	}

	defer bodyReader.Close()
	return SHA256Hex(bodyReader)
}

func deriveKey(secretKey, shortTime, region string) []byte {
	hmacDate := HMACSHA256String([]byte(secretPrefix+secretKey), shortTime)
	hmacRegion := HMACSHA256String(hmacDate, region)
	hmacService := HMACSHA256String(hmacRegion, signingService)
	return HMACSHA256String(hmacService, requestSuffix)
}

func HMACSHA256(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

func HMACSHA256String(key []byte, data string) []byte {
	return HMACSHA256(key, []byte(data))
}

func SHA256Hex(reader io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func buildStringToSign(ctx *SignatureContext) (string, error) {
	canonicalStr := buildCanonicalString(ctx)
	canonicalSHA, err := SHA256Hex(bytes.NewReader([]byte(canonicalStr)))
	if err != nil {
		return "", fmt.Errorf("failed to compute SHA from canonical str: %w", err)
	}
	return strings.Join([]string{
		ctx.Parameters.Algorithm,
		ctx.Parameters.Date,
		ctx.Parameters.Scope,
		canonicalSHA,
	}, "\n"), nil
}

func buildCanonicalString(ctx *SignatureContext) string {
	return strings.Join([]string{
		ctx.HTTPMethod,
		ctx.CanonicalURI,
		ctx.CanonicalQuery,
		ctx.CanonicalHeaders,
		ctx.SignedHeaders,
		ctx.Parameters.PayloadHash,
	}, "\n")
}
