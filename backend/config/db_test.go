package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"path"
	"testing"
)

func newTestDb(t *testing.T) *Db {
	db := NewDb(path.Join(t.TempDir(), "db"), log.Default())
	assert.NoError(t, db.Init())
	return db
}

func TestDb_Init_CreatesConfigTable(t *testing.T) {
	db := newTestDb(t)
	db.Upsert("k", "v")
	assert.Equal(t, "v", db.Get("k", ""))
}

func TestDb_Upsert_OverwritesExistingValue(t *testing.T) {
	db := newTestDb(t)
	db.Upsert("k", "first")
	db.Upsert("k", "second")
	assert.Equal(t, "second", db.Get("k", ""))
}

func TestDb_Delete_RemovesKey(t *testing.T) {
	db := newTestDb(t)
	db.Upsert("k", "v")
	db.Delete("k")
	assert.Nil(t, db.GetOrNilString("k"))
}

func TestDb_Get_ReturnsDefaultWhenMissing(t *testing.T) {
	db := newTestDb(t)
	assert.Equal(t, "fallback", db.Get("missing", "fallback"))
}

func TestDb_GetOrNilString_ReturnsNilWhenMissing(t *testing.T) {
	db := newTestDb(t)
	assert.Nil(t, db.GetOrNilString("missing"))
}

func TestDb_GetStringOrError_ErrorsWhenMissing(t *testing.T) {
	db := newTestDb(t)
	_, err := db.GetStringOrError("missing")
	assert.Error(t, err)
}

func TestDb_GetStringOrError_ReturnsValue(t *testing.T) {
	db := newTestDb(t)
	db.Upsert("k", "v")
	v, err := db.GetStringOrError("k")
	assert.NoError(t, err)
	assert.Equal(t, "v", v)
}

func TestDb_GetOrDefaultInt(t *testing.T) {
	db := newTestDb(t)
	assert.Equal(t, 0, db.GetOrDefaultInt("unknown", 0))
	db.Upsert("unknown", "1")
	assert.Equal(t, 1, db.GetOrDefaultInt("unknown", 0))
}

func TestDb_GetOrDefaultInt_FallsBackOnNonNumeric(t *testing.T) {
	db := newTestDb(t)
	db.Upsert("k", "not-a-number")
	assert.Equal(t, 7, db.GetOrDefaultInt("k", 7))
}

func TestDb_GetOrDefaultInt64(t *testing.T) {
	db := newTestDb(t)
	assert.Equal(t, int64(42), db.GetOrDefaultInt64("unknown", 42))
	db.Upsert("k", "9000")
	assert.Equal(t, int64(9000), db.GetOrDefaultInt64("k", 0))
}

func TestDb_GetOrDefaultString(t *testing.T) {
	db := newTestDb(t)
	assert.Equal(t, "default", db.GetOrDefaultString("unknown", "default"))
	db.Upsert("unknown", "test")
	assert.Equal(t, "test", db.GetOrDefaultString("unknown", "default"))
}

func TestDb_List(t *testing.T) {
	db := newTestDb(t)
	db.Upsert("a", "1")
	db.Upsert("b", "2")
	all := db.List()
	assert.Equal(t, "1", all["a"])
	assert.Equal(t, "2", all["b"])
	assert.Len(t, all, 2)
}

func TestDb_Exec_AllowsArbitrarySchema(t *testing.T) {
	db := newTestDb(t)
	_, err := db.Exec("create table extra (k varchar primary key)")
	assert.NoError(t, err)
	_, err = db.Exec("insert into extra values ('x')")
	assert.NoError(t, err)
}
