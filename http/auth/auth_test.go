package auth

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnauthorized(t *testing.T) {
	setup(t)
	defer teardown()
	a := assert.New(t)

	resp := testClient.GetBody("/profile")
	a.Equal(http.StatusUnauthorized, resp.Code)
	a.NotEqual(0, resp.Body.Len())
}

func TestLogin_emptyVals(t *testing.T) {
	setup(t)
	defer teardown()
	a := assert.New(t)

	vals := url.Values{}
	resp := testClient.PostForm("/login", vals)
	a.Equal(http.StatusBadRequest, resp.Code)
}

func TestLogin_badLogin(t *testing.T) {
	setup(t)
	defer teardown()
	a := assert.New(t)

	vals := url.Values{
		"user": {"false"},
		"pass": {"false"},
	}
	called := false
	testAuthProvider.checkMock = func(u, p string) (interface{}, error) {
		called = true
		return nil, ErrBadLogin
	}
	resp := testClient.PostForm("/login", vals)
	a.Equal(http.StatusBadRequest, resp.Code)
	a.True(called)
	a.Contains(resp.Body.String(), ErrBadLogin.Error())
}

func TestLogin_workingLogin(t *testing.T) {
	setup(t)
	defer teardown()
	a := assert.New(t)

	vals := url.Values{
		"user": {"testUser"},
		"pass": {"testPassw"},
	}
	called := false
	testAuthProvider.checkMock = func(u, p string) (interface{}, error) {
		called = true
		if !(u == "testUser" && p == "testPassw") {
			return nil, ErrBadLogin
		}
		return 23, nil
	}
	resp := testClient.PostForm("/login", vals)
	a.Equal(http.StatusSeeOther, resp.Code)
	a.Equal("/landingRedir", resp.Header().Get("Location"))
	a.True(called)
	newCookie := resp.Header().Get("Set-Cookie")
	a.Contains(newCookie, defaultSessionName)
}

func TestLogin_workingLoginAndRestrictedAcc(t *testing.T) {
	setup(t)
	defer teardown()
	a := assert.New(t)

	vals := url.Values{
		"user": {"testUser"},
		"pass": {"testPassw"},
	}
	called := false
	testAuthProvider.checkMock = func(u, p string) (interface{}, error) {
		called = true
		if !(u == "testUser" && p == "testPassw") {
			return nil, ErrBadLogin
		}
		return 23, nil
	}
	resp := testClient.PostForm("/login", vals)
	a.Equal(http.StatusSeeOther, resp.Code)
	a.True(called)
	newCookie := resp.Header().Get("Set-Cookie")
	a.Contains(newCookie, defaultSessionName)

	testClient.SetHeaders(http.Header{"Cookie": []string{newCookie}})
	resp2 := testClient.GetBody("/profile")
	a.Equal(http.StatusOK, resp2.Code)
}

func TestLogin_workingLoginAndLogout(t *testing.T) {
	setup(t)
	defer teardown()
	a := assert.New(t)

	vals := url.Values{
		"user": {"testUser"},
		"pass": {"testPassw"},
	}
	called := false
	testAuthProvider.checkMock = func(u, p string) (interface{}, error) {
		called = true
		if !(u == "testUser" && p == "testPassw") {
			return nil, ErrBadLogin
		}
		return 23, nil
	}
	resp := testClient.PostForm("/login", vals)
	a.Equal(http.StatusSeeOther, resp.Code)
	a.True(called)
	newCookie := resp.Header().Get("Set-Cookie")
	a.Contains(newCookie, defaultSessionName)

	testClient.SetHeaders(http.Header{"Cookie": []string{newCookie}})
	resp2 := testClient.GetBody("/logout")
	logoutCookie := resp2.Header().Get("Set-Cookie")
	a.Equal("/landingRedir", resp2.Header().Get("Location"))
	a.NotEqual("", logoutCookie)
	a.NotEqual(newCookie, logoutCookie)

	testClient.ClearHeaders()
	testClient.SetHeaders(http.Header{"Cookie": []string{logoutCookie}})
	resp3 := testClient.GetBody("/profile")
	a.Equal(http.StatusUnauthorized, resp3.Code)
	a.Equal("Not Authorized\n", resp3.Body.String(), "Body %q", resp3.Body.String())
}
