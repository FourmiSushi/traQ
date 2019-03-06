package router

import (
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func TestHandlers_GetPublicUserIcon(t *testing.T) {
	t.Parallel()
	repo, server, _, _, _, _, testUser, _ := setupWithUsers(t, common5)

	t.Run("No name", func(t *testing.T) {
		t.Parallel()
		e := makeExp(t, server)
		e.GET("/api/1.0/public/icon/").
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("No user", func(t *testing.T) {
		t.Parallel()
		e := makeExp(t, server)
		e.GET("/api/1.0/public/icon/no+user").
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		_, require := assertAndRequire(t)

		_, src, err := repo.OpenFile(testUser.Icon)
		require.NoError(err)
		i, err := ioutil.ReadAll(src)
		require.NoError(err)

		e := makeExp(t, server)
		e.GET("/api/1.0/public/icon/{username}", testUser.Name).
			Expect().
			Status(http.StatusOK).
			Header(echo.HeaderContentLength).
			Equal(strconv.Itoa(len(i)))
	})

	t.Run("Success with thumbnail", func(t *testing.T) {
		t.Parallel()
		e := makeExp(t, server)
		e.GET("/api/1.0/public/icon/{username}", testUser.Name).
			WithQuery("thumb", "").
			Expect().
			Status(http.StatusOK)
	})
}

func TestHandlers_GetPublicEmojiJSON(t *testing.T) {
	t.Parallel()
	repo, server, _, _, _, _ := setup(t, s3)

	var stamps []interface{}
	for i := 0; i < 10; i++ {
		s := mustMakeStamp(t, repo, random, uuid.Nil)
		stamps = append(stamps, s.Name)
	}

	e := makeExp(t, server)
	e.GET("/api/1.0/public/emoji.json").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("all").
		Array().
		ContainsOnly(stamps...)
}

func TestHandlers_GetPublicEmojiCSS(t *testing.T) {
	t.Parallel()
	repo, server, _, _, _, _ := setup(t, s4)

	for i := 0; i < 10; i++ {
		mustMakeStamp(t, repo, random, uuid.Nil)
	}

	e := makeExp(t, server)
	e.GET("/api/1.0/public/emoji.css").
		Expect().
		Status(http.StatusOK).
		ContentType("text/css")
}

func TestHandlers_GetPublicEmojiImage(t *testing.T) {
	t.Parallel()
	repo, server, _, _, _, _ := setup(t, common5)

	s := mustMakeStamp(t, repo, random, uuid.Nil)

	t.Run("Not Found", func(t *testing.T) {
		t.Parallel()
		e := makeExp(t, server)
		e.GET("/api/1.0/public/emoji/{stampID}", uuid.NewV4()).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		e := makeExp(t, server)
		e.GET("/api/1.0/public/emoji/{stampID}", s.ID).
			Expect().
			Status(http.StatusOK)
	})
}