package rwgps_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bzimmer/gravl/pkg/rwgps"
	"github.com/stretchr/testify/assert"
)

func Test_User(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c, err := newClient(http.StatusOK, "rwgps_users_1122.json")
	a.NoError(err)
	a.NotNil(c)

	ctx := context.Background()
	user, err := c.Users.AuthenticatedUser(ctx)
	a.NoError(err)
	a.NotNil(user)
	a.Equal(rwgps.UserID(1122), user.ID)
}
