package main

import (
  "github.com/aws/aws-sdk-go/service/cloudfront/sign"
  "os"
  "log"
  "time"
  "net/http"
  "encoding/json"
  "strings"
)

type SignResponse struct{
  SignedURL string
}

// Sign URL to be valid for 1 hour from now.
func getCloudfrontSignedUrl(rawURL string) string {)
  privKeyReader := strings.NewReader(os.Getenv("AWS_CLOUDFRONT_PRIVATE_KEY"))
  privKey, err := sign.LoadPEMPrivKey(privKeyReader)
  if err != nil {
      log.Fatalf("Failed load private key, err: %s\n", err.Error())
  }
  signer := sign.NewURLSigner(os.Getenv("AWS_CLOUDFRONT_KEYPAIR_ID"), privKey)
  signedURL, err := signer.Sign(rawURL, time.Now().Add(1*time.Hour))
  if err != nil {
      log.Fatalf("Failed to sign url, err: %s\n", err.Error())
  }

  return signedURL
}

func signHandler(w http.ResponseWriter, r *http.Request) {
  rawURL := r.URL.Query().Get("url")
  s := SignResponse{SignedURL: getCloudfrontSignedUrl(rawURL)}
  enc := json.NewEncoder(w)
  enc.SetEscapeHTML(false)
  enc.Encode(s)
}

func main() {
  port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/sign/", signHandler)
	http.ListenAndServe(":" + port, nil)
}
