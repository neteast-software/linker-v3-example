package captcha

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/neteast-software/go-module/http/gateway"
)

func TestCaptchaFilterConsumesChallengeOnce(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	const salt = "example-captcha-salt"
	challenge, err := New("challenge-one", "9271", salt, now.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	store := Memory()
	store.now = func() time.Time { return now }
	if err = store.Save(context.Background(), challenge); err != nil {
		t.Fatal(err)
	}
	factory := Filter(store, salt)
	factory.now = func() time.Time { return now }
	filter, err := factory.Build(nil)
	if err != nil {
		t.Fatal(err)
	}
	before := filter.(gateway.BeforeFilter)

	request, _ := http.NewRequest(http.MethodPost, "http://gateway.local/login", nil)
	request.Header.Set("X-Captcha-Challenge", challenge.ID)
	request.Header.Set("X-Captcha-Answer", "9271")
	response, err := before.Before(context.Background(), request)
	if err != nil || response != nil {
		t.Fatalf("first = (%v, %v)", response, err)
	}
	if request.Header.Get("X-Captcha-Answer") != "" {
		t.Fatal("验证码答案被转发到 upstream")
	}

	request, _ = http.NewRequest(http.MethodPost, "http://gateway.local/login", nil)
	request.Header.Set("X-Captcha-Challenge", challenge.ID)
	request.Header.Set("X-Captcha-Answer", "9271")
	response, err = before.Before(context.Background(), request)
	if err != nil || response.StatusCode != http.StatusBadRequest {
		t.Fatalf("second = (%v, %v)", response, err)
	}
}

func TestCaptchaFilterRejectsOversizedInputWithoutStoreRead(t *testing.T) {
	filter, err := Filter(Memory(), "example-captcha-salt").Build(nil)
	if err != nil {
		t.Fatal(err)
	}
	request, _ := http.NewRequest(http.MethodPost, "http://gateway.local/login", nil)
	request.Header.Set("X-Captcha-Challenge", "challenge")
	request.Header.Set("X-Captcha-Answer", string(make([]byte, 257)))
	response, err := filter.(gateway.BeforeFilter).Before(context.Background(), request)
	if err != nil || response.StatusCode != http.StatusBadRequest {
		t.Fatalf("response = (%v, %v)", response, err)
	}
}
