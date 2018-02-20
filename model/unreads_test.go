package model

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func beforeUnreadsTest(t *testing.T) (*Unread, *Unread, *Unread) {
	testMessage := mustMakeMessage(t)

	testUnread := &Unread{
		UserID:    testUserID,
		MessageID: testMessage.ID,
	}
	emptyUserIDUnread := &Unread{
		MessageID: testMessage.ID,
	}
	emptyMessageIDUnread := &Unread{
		UserID: testUserID,
	}
	return testUnread, emptyUserIDUnread, emptyMessageIDUnread
}

func TestUnreadTableName(t *testing.T) {
	assert.Equal(t, "unreads", (&Unread{}).TableName())
}

func TestUnreadCreate(t *testing.T) {
	beforeTest(t)
	testUnread, emptyUserIDUnread, emptyMessageIDUnread := beforeUnreadsTest(t)

	// 正常系
	assert.NoError(t, testUnread.Create())
	unreads, err := GetUnreadsByUserID(testUnread.UserID)
	require.NoError(t, err)
	assert.Len(t, unreads, 1)
	assert.Equal(t, *unreads[0], *testUnread)

	// 異常系
	assert.Error(t, emptyUserIDUnread.Create())
	assert.Error(t, emptyMessageIDUnread.Create())
}

func TestUnreadDelete(t *testing.T) {
	beforeTest(t)
	testUnread, emptyUserIDUnread, emptyMessageIDUnread := beforeUnreadsTest(t)

	// 正常系
	require.NoError(t, testUnread.Create())
	assert.NoError(t, testUnread.Delete())
	unreads, err := GetUnreadsByUserID(testUnread.UserID)
	require.NoError(t, err)
	assert.Len(t, unreads, 0)

	// 異常系
	assert.Error(t, emptyUserIDUnread.Delete())
	assert.Error(t, emptyMessageIDUnread.Delete())
}

func TestGetUnreadsByUserID(t *testing.T) {
	beforeTest(t)
	testUnread, _, _ := beforeUnreadsTest(t)

	// 正常系
	require.NoError(t, testUnread.Create())
	unreads, err := GetUnreadsByUserID(testUnread.UserID)
	assert.NoError(t, err)
	assert.Len(t, unreads, 1)
	assert.Equal(t, *unreads[0], *testUnread)

	// 異常系
	_, emptyErr := GetUnreadsByUserID("")
	assert.Error(t, emptyErr)
	nobodyUnreads, nobodyErr := GetUnreadsByUserID(nobodyID)
	assert.NoError(t, nobodyErr)
	assert.Len(t, nobodyUnreads, 0)
}

func TestDeleteUnreadsByMessageID(t *testing.T) {
	beforeTest(t)
	assert := assert.New(t)

	testMessage := mustMakeMessage(t)
	testMessage2 := mustMakeMessage(t)

	for i := 0; i < 10; i++ {
		mustMakeMessageUnread(t, CreateUUID(), testMessage.ID)
		mustMakeMessageUnread(t, CreateUUID(), testMessage2.ID)
	}

	// 正常系
	if assert.NoError(DeleteUnreadsByMessageID(testMessage.ID)) {
		if n, err := db.Count(&Unread{}); assert.NoError(err) {
			assert.EqualValues(10, n)
		}
	}
	if assert.NoError(DeleteUnreadsByMessageID(testMessage2.ID)) {
		if n, err := db.Count(&Unread{}); assert.NoError(err) {
			assert.EqualValues(0, n)
		}
	}

	// 異常系
	assert.Error(DeleteUnreadsByMessageID(""))
}
