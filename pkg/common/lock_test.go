package common

import "testing"

func TestResourceLocks_TryAcquire(t *testing.T) {
	lock := NewResourceLocks()
	tests := []struct {
		name       string
		resourceId string
		isAcquire  bool
	}{
		{
			name:       "first acquire",
			resourceId: "vol-123456",
			isAcquire:  true,
		},
		{
			name:       "acquire failed",
			resourceId: "vol-123456",
			isAcquire:  false,
		},
		{
			name:       "acquire another",
			resourceId: "vol-234567",
			isAcquire:  true,
		},
	}
	for _, test := range tests {
		res := lock.TryAcquire(test.resourceId)
		if test.isAcquire != res {
			t.Errorf("name %s: expect %t, but actually %t", test.name, test.isAcquire, res)
		}
	}
}

func TestResourceLocks_Release(t *testing.T) {
	lock := NewResourceLocks()
	resourceId1 := "vol-00001"
	resourceId2 := "vol-00002"
	// release a not exist resource
	lock.Release(resourceId1)
	// release a exist resource
	if res := lock.TryAcquire(resourceId2); !res {
		t.Errorf("try acquire %s failed", resourceId2)
	}
	lock.Release(resourceId2)
	// release multiple times
	lock.Release(resourceId2)
}
