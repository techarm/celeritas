package cache

import "testing"

func TestBadgerCache_Has(t *testing.T) {
	err := testBadgerCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}

	incache, err := testBadgerCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if incache {
		t.Error("foo found in cache, and it shouldn't be there")
	}

	err = testBadgerCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	incache, err = testBadgerCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if !incache {
		t.Error("foo not found in cache, and it should be there")
	}

	_ = testBadgerCache.Forget("foo")
}

func TestBadgerCache_Get(t *testing.T) {
	err := testBadgerCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	x, err := testBadgerCache.Get("foo")
	if err != nil {
		t.Error(err)
	}

	if x != "bar" {
		t.Error("dit not get correct value from cache")
	}
}

func TestBadgerCache_Forget(t *testing.T) {
	err := testBadgerCache.Set("alpha", "beta")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Forget("alpha")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("alpha")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("alpha found in cache, and it shouldn't be there")
	}
}

func TestBadgerCache_Empty(t *testing.T) {
	err := testBadgerCache.Set("alpha", "beta")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Empty()
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("alpha")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("alpha found in cache, and it shouldn't be there")
	}
}

func TestBadgerCache_EmptyByMatch(t *testing.T) {
	err := testBadgerCache.Set("alpha", "foo")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Set("alpha2", "foo")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Set("beta", "foo")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.EmptyByMatch("alpha")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("alpha1")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("alpha1 found in cache, and it shouldn't be there")
	}

	inCache, err = testBadgerCache.Has("alpha2")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("alpha2 found in cache, and it shouldn't be there")
	}

	inCache, err = testBadgerCache.Has("beta")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("beta not found in cache, and it should be there")
	}
}
