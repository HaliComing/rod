// The .github/workflows/check-issues.yml will use it as an github action
// To test it locally, you can generate a personal github token: https://github.com/settings/tokens
// Then run this:
//   GITHUB_TOKEN=your_token go run ./lib/utils/check-issue

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/halicoming/rod/lib/utils"
	"github.com/ysmood/gson"
)

func main() {
	id, body := getIssue()

	log.Println("check issue", id)

	msg := check(body)

	deleteComments(id)

	if msg != "" {
		sendComment(id, msg)
	}
}

func check(body string) string {
	msg := []string{}

	for _, check := range []func(string) error{
		checkVersion, checkMarkdown, checkGoCode,
	} {
		err := check(body)
		if err != nil {
			msg = append(msg, err.Error())
		}
	}
	if len(msg) != 0 {
		return strings.Join(msg, "\n\n")
	}
	return ""
}

func getIssue() (int, string) {
	data, err := os.Open(os.Getenv("GITHUB_EVENT_PATH"))
	utils.E(err)

	issue := gson.New(data).Get("issue")

	id := issue.Get("number").Int()
	body := issue.Get("body").Str()

	return id, body
}

func sendComment(id int, msg string) {
	msg += fmt.Sprintf(
		"\n\n_<sub>generated by [check-issue](https://github.com/halicoming/rod/actions/runs/%s)</sub>_",
		os.Getenv("GITHUB_RUN_ID"),
	)

	q := req(fmt.Sprintf("/repos/go-rod/rod/issues/%d/comments", id))
	q.Method = http.MethodPost
	q.Body = io.NopCloser(bytes.NewBuffer(utils.MustToJSONBytes(map[string]string{"body": msg})))
	res, err := http.DefaultClient.Do(q)
	utils.E(err)
	defer func() { _ = res.Body.Close() }()
	resE(res)
}

func deleteComments(id int) {
	q := req(fmt.Sprintf("/repos/go-rod/rod/issues/%d/comments", id))
	res, err := http.DefaultClient.Do(q)
	utils.E(err)
	resE(res)

	list := gson.New(res.Body)

	for _, c := range list.Arr() {
		if c.Get("user.login").Str() == "github-actions[bot]" &&
			strings.Contains(c.Get("body").Str(), "[check-issue]") {
			iid := c.Get("id").Int()
			q := req(fmt.Sprintf("/repos/go-rod/rod/issues/comments/%d", iid))
			q.Method = http.MethodDelete
			res, err := http.DefaultClient.Do(q)
			utils.E(err)
			resE(res)
		}
	}
}

func req(u string) *http.Request {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		panic("missing github token")
	}

	r, err := http.NewRequest(http.MethodGet, "https://api.github.com"+u, nil)
	utils.E(err)
	r.Header.Add("Authorization", "token "+token)
	return r
}

func resE(res *http.Response) {
	if res.StatusCode >= 400 {
		str, err := io.ReadAll(res.Body)
		utils.E(err)
		panic(string(str))
	}
}
