package repository

import (
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traQ/utils/optional"
	random2 "github.com/traPtitech/traQ/utils/random"
	"testing"
)

func TestRepositoryImpl_CreateStamp(t *testing.T) {
	t.Parallel()
	repo, _, require, user := setupWithUser(t, common2)

	fid, err := GenerateIconFile(repo, "stamp")
	require.NoError(err)

	t.Run("nil file id", func(t *testing.T) {
		t.Parallel()

		_, err := repo.CreateStamp(CreateStampArgs{Name: random2.AlphaNumeric(20), FileID: uuid.Nil, CreatorID: user.GetID()})
		assert.Error(t, err)
	})

	t.Run("invalid name", func(t *testing.T) {
		t.Parallel()

		_, err := repo.CreateStamp(CreateStampArgs{Name: "あ", FileID: fid, CreatorID: user.GetID()})
		assert.Error(t, err)
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()

		_, err := repo.CreateStamp(CreateStampArgs{Name: random2.AlphaNumeric(20), FileID: uuid.Must(uuid.NewV4()), CreatorID: user.GetID()})
		assert.Error(t, err)
	})

	t.Run("duplicate name", func(t *testing.T) {
		t.Parallel()
		s := mustMakeStamp(t, repo, rand, uuid.Nil)

		_, err := repo.CreateStamp(CreateStampArgs{Name: s.Name, FileID: fid, CreatorID: user.GetID()})
		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		name := random2.AlphaNumeric(20)
		s, err := repo.CreateStamp(CreateStampArgs{Name: name, FileID: fid, CreatorID: user.GetID()})
		if assert.NoError(err) {
			assert.NotEmpty(s.ID)
			assert.Equal(name, s.Name)
			assert.Equal(fid, s.FileID)
			assert.Equal(user.GetID(), s.CreatorID)
			assert.NotEmpty(s.CreatedAt)
			assert.NotEmpty(s.UpdatedAt)
			assert.Nil(s.DeletedAt)
		}
	})
}

func TestRepositoryImpl_UpdateStamp(t *testing.T) {
	t.Parallel()
	repo, _, _ := setup(t, common2)

	s := mustMakeStamp(t, repo, rand, uuid.Nil)

	t.Run("nil id", func(t *testing.T) {
		t.Parallel()

		assert.EqualError(t, repo.UpdateStamp(uuid.Nil, UpdateStampArgs{}), ErrNilID.Error())
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		assert.EqualError(t, repo.UpdateStamp(uuid.Must(uuid.NewV4()), UpdateStampArgs{}), ErrNotFound.Error())
	})

	t.Run("no change", func(t *testing.T) {
		t.Parallel()

		assert.NoError(t, repo.UpdateStamp(s.ID, UpdateStampArgs{}))
	})

	t.Run("invalid name", func(t *testing.T) {
		t.Parallel()

		assert.Error(t, repo.UpdateStamp(s.ID, UpdateStampArgs{Name: optional.StringFrom("あ")}))
	})

	t.Run("duplicate name", func(t *testing.T) {
		t.Parallel()
		s2 := mustMakeStamp(t, repo, rand, uuid.Nil)

		assert.Error(t, repo.UpdateStamp(s.ID, UpdateStampArgs{Name: optional.StringFrom(s2.Name)}))
	})

	t.Run("nil file id", func(t *testing.T) {
		t.Parallel()

		assert.Error(t, repo.UpdateStamp(s.ID, UpdateStampArgs{FileID: optional.UUIDFrom(uuid.Nil)}))
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()

		assert.Error(t, repo.UpdateStamp(s.ID, UpdateStampArgs{FileID: optional.UUIDFrom(uuid.Must(uuid.NewV4()))}))
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert, require := assertAndRequire(t)

		s := mustMakeStamp(t, repo, rand, uuid.Nil)
		newFile, err := GenerateIconFile(repo, "stamp")
		require.NoError(err)
		newName := random2.AlphaNumeric(20)

		if assert.NoError(repo.UpdateStamp(s.ID, UpdateStampArgs{
			Name:      optional.StringFrom(newName),
			FileID:    optional.UUIDFrom(newFile),
			CreatorID: optional.UUIDFrom(uuid.Nil),
		})) {
			a, err := repo.GetStamp(s.ID)
			require.NoError(err)
			assert.Equal(newFile, a.FileID)
			assert.Equal(newName, a.Name)
		}
	})
}

func TestRepositoryImpl_GetStamp(t *testing.T) {
	t.Parallel()
	repo, _, _ := setup(t, common2)

	t.Run("nil id", func(t *testing.T) {
		t.Parallel()

		_, err := repo.GetStamp(uuid.Nil)
		assert.EqualError(t, err, ErrNotFound.Error())
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		_, err := repo.GetStamp(uuid.Must(uuid.NewV4()))
		assert.EqualError(t, err, ErrNotFound.Error())
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		a := mustMakeStamp(t, repo, rand, uuid.Nil)

		s, err := repo.GetStamp(a.ID)
		if assert.NoError(err) {
			assert.Equal(a.ID, s.ID)
			assert.Equal(a.Name, s.Name)
			assert.Equal(a.FileID, s.FileID)
			assert.Equal(a.CreatorID, s.CreatorID)
		}
	})
}

func TestRepositoryImpl_DeleteStamp(t *testing.T) {
	t.Parallel()
	repo, _, _ := setup(t, common2)

	t.Run("nil id", func(t *testing.T) {
		t.Parallel()

		assert.EqualError(t, repo.DeleteStamp(uuid.Nil), ErrNilID.Error())
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		assert.EqualError(t, repo.DeleteStamp(uuid.Must(uuid.NewV4())), ErrNotFound.Error())
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		s := mustMakeStamp(t, repo, rand, uuid.Nil)
		if assert.NoError(repo.DeleteStamp(s.ID)) {
			_, err := repo.GetStamp(s.ID)
			assert.EqualError(err, ErrNotFound.Error())
		}
	})
}

func TestRepositoryImpl_GetAllStamps(t *testing.T) {
	t.Parallel()
	repo, assert, _ := setup(t, ex1)

	n := 10
	for i := 0; i < 10; i++ {
		mustMakeStamp(t, repo, rand, uuid.Nil)
	}

	arr, err := repo.GetAllStamps(false)
	if assert.NoError(err) {
		assert.Len(arr, n)
	}
}

func TestRepositoryImpl_StampExists(t *testing.T) {
	t.Parallel()
	repo, _, _ := setup(t, common2)

	s := mustMakeStamp(t, repo, rand, uuid.Nil)

	t.Run("nil id", func(t *testing.T) {
		t.Parallel()

		ok, err := repo.StampExists(uuid.Nil)
		if assert.NoError(t, err) {
			assert.False(t, ok)
		}
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		ok, err := repo.StampExists(uuid.Must(uuid.NewV4()))
		if assert.NoError(t, err) {
			assert.False(t, ok)
		}
	})

	t.Run("found", func(t *testing.T) {
		t.Parallel()

		ok, err := repo.StampExists(s.ID)
		if assert.NoError(t, err) {
			assert.True(t, ok)
		}
	})
}

func TestRepositoryImpl_ExistStamps(t *testing.T) {
	t.Parallel()
	repo, _, _ := setup(t, common2)

	stampIDs := make([]uuid.UUID, 0, 10)

	for i := 0; i < 10; i++ {
		s := mustMakeStamp(t, repo, rand, uuid.Nil)
		stampIDs = append(stampIDs, s.ID)
	}

	t.Run("argument err", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		stampIDsCopy := make([]uuid.UUID, len(stampIDs), cap(stampIDs))
		_ = copy(stampIDsCopy, stampIDs)
		if assert.True(len(stampIDsCopy) > 0) {
			stampIDsCopy[0] = uuid.Must(uuid.NewV4())
		}
		assert.Error(repo.ExistStamps(stampIDsCopy))
	})

	t.Run("sucess", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		assert.NoError(repo.ExistStamps(stampIDs))
	})
}

func TestRepositoryImpl_GetUserStampHistory(t *testing.T) {
	t.Parallel()
	repo, _, _, user, channel := setupWithUserAndChannel(t, common2)

	message := mustMakeMessage(t, repo, user.GetID(), channel.ID)
	stamp1 := mustMakeStamp(t, repo, rand, uuid.Nil)
	stamp2 := mustMakeStamp(t, repo, rand, uuid.Nil)
	stamp3 := mustMakeStamp(t, repo, rand, uuid.Nil)
	mustAddMessageStamp(t, repo, message.ID, stamp1.ID, user.GetID())
	mustAddMessageStamp(t, repo, message.ID, stamp3.ID, user.GetID())
	mustAddMessageStamp(t, repo, message.ID, stamp2.ID, user.GetID())

	t.Run("Nil id", func(t *testing.T) {
		t.Parallel()
		ms, err := repo.GetUserStampHistory(uuid.Nil, 0)
		if assert.NoError(t, err) {
			assert.Empty(t, ms)
		}
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ms, err := repo.GetUserStampHistory(user.GetID(), 0)
		if assert.NoError(t, err) && assert.Len(t, ms, 3) {
			assert.Equal(t, ms[0].StampID, stamp2.ID)
			assert.Equal(t, ms[1].StampID, stamp3.ID)
			assert.Equal(t, ms[2].StampID, stamp1.ID)
		}
	})
}
