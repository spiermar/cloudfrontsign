package main

import (
  "github.com/aws/aws-sdk-go/service/cloudfront/sign"
  "os"
  "log"
  "time"
  "strings"
  "gopkg.in/gin-gonic/gin.v1"
)

type SignResponse struct{
  SignedURL string
}

// Sign URL to be valid for 1 hour from now.
func getCloudfrontSignedUrl(rawURL string) string {
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

func signHandler(c *gin.Context) {
  rawURL := c.Query("url")
  s := SignResponse{SignedURL: getCloudfrontSignedUrl(rawURL)}
  c.JSON(200, gin.H{
      "SignedURL": s,
  })
}

func main() {
  port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

  r := gin.Default()
  r.GET("/sign", signHandler)
  r.Run(":" + port)
}
