package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/Southclaws/fault/fctx"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
)

func TestWithMeta(t *testing.T) {
	ctx := context.Background()
	ctx = fctx.WithMeta(ctx, "key", "value")

	err := fctx.Wrap(errors.New("a problem"), ctx)
	data := fctx.Unwrap(err)

	assert.Equal(t, map[string]string{"key": "value"}, data)
}

func TestWithMetaAdditional(t *testing.T) {
	ctx := context.Background()
	ctx = fctx.WithMeta(ctx, "key", "value")

	err := fctx.Wrap(errors.New("a problem"), ctx, "additional", "value")
	data := fctx.Unwrap(err)

	assert.Equal(t, map[string]string{
		"key":        "value",
		"additional": "value",
	}, data)
}

func TestWithMetaOverwrite(t *testing.T) {
	ctx := context.Background()
	ctx = fctx.WithMeta(ctx, "key", "value")
	ctx = fctx.WithMeta(ctx, "key", "value2")

	err := fctx.Wrap(errors.New("a problem"), ctx)
	data := fctx.Unwrap(err)

	assert.Equal(t, map[string]string{"key": "value2"}, data)
}

func TestWithMetaNested(t *testing.T) {
	ctx := context.Background()
	ctx = fctx.WithMeta(ctx, "key", "value")
	ctx = fctx.WithMeta(ctx, "key", "value2")
	ctx = context.WithValue(ctx, "some other", "stuff")
	ctx = fctx.WithMeta(ctx, "key", "value3")

	err := fctx.Wrap(errors.New("a problem"), ctx)
	data := fctx.Unwrap(err)

	assert.Equal(t, map[string]string{"key": "value3"}, data)
}

func TestWithMetaNestedManyKeys(t *testing.T) {
	ctx := context.Background()
	ctx = fctx.WithMeta(ctx, "key1", "value1")
	ctx = context.WithValue(ctx, "some other", "stuff")
	ctx = fctx.WithMeta(ctx, "key2", "value2")
	ctx = fctx.WithMeta(ctx, "key3", "value3", "key4", "value4")

	err := fctx.Wrap(errors.New("a problem"), ctx)
	data := fctx.Unwrap(err)

	assert.Equal(t, map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
	}, data)
}

func TestWithMetaNestedManyKeysPlusExtraWrappedKV(t *testing.T) {
	ctx := context.Background()
	ctx = fctx.WithMeta(ctx, "key1", "value1")
	ctx = context.WithValue(ctx, "some other", "stuff")
	ctx = fctx.WithMeta(ctx, "key2", "value2")
	ctx = fctx.WithMeta(ctx, "key3", "value3", "key4", "value4")

	err := fctx.Wrap(errors.New("a problem"), ctx, "extra1", "extravalue1", "extra2", "extravalue2")
	data := fctx.Unwrap(err)

	assert.Equal(t, map[string]string{
		"key1":   "value1",
		"key2":   "value2",
		"key3":   "value3",
		"key4":   "value4",
		"extra1": "extravalue1",
		"extra2": "extravalue2",
	}, data)
}

func TestWithMetaOddNumberKV(t *testing.T) {
	ctx := context.Background()
	ctx = fctx.WithMeta(ctx, "key", "value", "ignored")

	err := fctx.Wrap(errors.New("a problem"), ctx)
	data := fctx.Unwrap(err)

	assert.Equal(t, map[string]string{"key": "value"}, data)
}

func TestWithMetaOddNumberWrapKV(t *testing.T) {
	ctx := context.Background()
	ctx = fctx.WithMeta(ctx, "key", "value", "ignored")

	err := fctx.Wrap(errors.New("a problem"), ctx, "wrapkey", "wrapvalue", "ignored")
	data := fctx.Unwrap(err)

	assert.Equal(t, map[string]string{"key": "value", "wrapkey": "wrapvalue"}, data)
}

func TestWithMetaOneValueKV(t *testing.T) {
	ctx := context.Background()
	ctx = fctx.WithMeta(ctx, "ignored")

	err := fctx.Wrap(errors.New("a problem"), ctx)
	data := fctx.Unwrap(err)

	assert.Nil(t, data)
}

func TestWithMetaOneValueWrapKV(t *testing.T) {
	ctx := context.Background()
	ctx = fctx.WithMeta(ctx, "ignored")

	err := fctx.Wrap(errors.New("a problem"), ctx, "wrapkey", "wrapvalue", "ignored")
	data := fctx.Unwrap(err)

	assert.Equal(t, map[string]string{"wrapkey": "wrapvalue"}, data)
}

func TestWithMetaOneValueEmptyWrapKV(t *testing.T) {
	ctx := context.Background()
	ctx = fctx.WithMeta(ctx, "ignored")

	err := fctx.Wrap(errors.New("a problem"), ctx, "ignored")
	data := fctx.Unwrap(err)

	assert.Nil(t, data)
}

func TestWithMetaEmpty(t *testing.T) {
	err := errors.New("a problem")
	data := fctx.Unwrap(err)

	assert.Nil(t, data)
}

func TestWithMetaDifferentMapAddress(t *testing.T) {
	ctx := context.Background()
	err := errors.New("a problem")

	ctx1 := fctx.WithMeta(ctx, "key1", "value1")
	err1 := fctx.Wrap(err, ctx1)

	ctx2 := fctx.WithMeta(ctx1, "key2", "value2")
	err2 := fctx.Wrap(err1, ctx2)

	data1 := fctx.Unwrap(err1)
	data2 := fctx.Unwrap(err2)

	assert.Equal(t,
		`&context.valueCtx{
    Context: &context.emptyCtx(0),
    key:     fctx.contextKey{},
    val:     map[string]string{"key1":"value1"},
}`,
		pretty.Sprint(ctx1),
		"The map from the first context should be left unmodified by the second call to fctx.WithMeta",
	)

	assert.Equal(t,
		`&context.valueCtx{
    Context: &context.valueCtx{
        Context: &context.emptyCtx(0),
        key:     fctx.contextKey{},
        val:     map[string]string{"key1":"value1"},
    },
    key: fctx.contextKey{},
    val: map[string]string{"key1":"value1", "key2":"value2"},
}`,
		pretty.Sprint(ctx2),
		"The second context value should contain both maps and the first only contains the first key-value pair.",
	)

	assert.Equal(t,
		`map[string]string{"key1":"value1"}`,
		pretty.Sprint(data1),
		"The first unwrap result is only the key-value pair from the first wrap.",
	)

	assert.Equal(t,
		`map[string]string{"key1":"value1", "key2":"value2"}`,
		pretty.Sprint(data2),
		"The second unwrap result contains all the data merged together.",
	)
}
