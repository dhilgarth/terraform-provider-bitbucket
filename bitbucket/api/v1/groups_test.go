package v1

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/stretchr/testify/assert"
)

func TestGroups(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("ENV TF_ACC=1 not set")
	}

	c := NewClient()

	var groupResourceSlug string

	name := "tf-bb-group-test" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	owner := os.Getenv("BITBUCKET_WORKSPACE")

	t.Run("create", func(t *testing.T) {
		opt := &GroupOptions{
			OwnerUuid: owner,
			Name:      name,
		}

		group, err := c.Groups.Create(opt)

		assert.NoError(t, err)
		assert.Equal(t, name, group.Name)
		assert.Equal(t, "none", group.Permission)

		groupResourceSlug = group.Slug
	})

	t.Run("get", func(t *testing.T) {
		opt := &GroupOptions{
			OwnerUuid: owner,
			Slug:      groupResourceSlug,
		}
		group, err := c.Groups.Get(opt)

		assert.NoError(t, err)
		assert.Equal(t, name, group.Name)
		assert.Equal(t, "none", group.Permission)
		assert.Equal(t, groupResourceSlug, group.Slug)
	})

	t.Run("update", func(t *testing.T) {
		opt := &GroupOptions{
			OwnerUuid:  owner,
			Slug:       groupResourceSlug,
			Permission: "write",
		}
		group, err := c.Groups.Update(opt)

		assert.NoError(t, err)
		assert.Equal(t, name, group.Name)
		assert.Equal(t, "write", group.Permission)
		assert.Equal(t, groupResourceSlug, group.Slug)
	})

	t.Run("delete", func(t *testing.T) {
		opt := &GroupOptions{
			OwnerUuid: owner,
			Slug:      groupResourceSlug,
		}
		err := c.Groups.Delete(opt)
		assert.NoError(t, err)
	})
}

func TestGroupsGracefullyHandleNoReturnedGroupsForInvalidSlug(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("ENV TF_ACC=1 not set")
	}

	c := NewClient()

	var groupResourceSlug string

	name := "TF-BB-Group-Test"
	owner := os.Getenv("BITBUCKET_WORKSPACE")

	t.Run("create", func(t *testing.T) {
		opt := &GroupOptions{
			OwnerUuid: owner,
			Name:      name,
		}

		group, err := c.Groups.Create(opt)
		assert.NoError(t, err)

		assert.Equal(t, name, group.Name)
		assert.Equal(t, "none", group.Permission)

		groupResourceSlug = group.Slug
	})

	t.Run("get", func(t *testing.T) {
		opt := &GroupOptions{
			OwnerUuid: owner,
			Slug:      name, // Slugs are lowercase and the BB's API is case-sensitive, this will trigger a fail response
		}
		group, err := c.Groups.Get(opt)
		assert.Nil(t, group)
		assert.EqualError(t, err, "no group found")
	})

	t.Run("delete", func(t *testing.T) {
		opt := &GroupOptions{
			OwnerUuid: owner,
			Slug:      groupResourceSlug,
		}
		err := c.Groups.Delete(opt)
		assert.NoError(t, err)
	})
}
