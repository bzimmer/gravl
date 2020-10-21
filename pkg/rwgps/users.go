package rwgps

import "context"

// UsersService .
type UsersService service

// Paginator .
type Paginator struct {
	offset int
	limit  int
}

// AuthenticatedUser .
func (UsersService) AuthenticatedUser(ctx context.Context) (*User, error) {
	return nil, nil
}

// Routes .
func (UsersService) Routes(ctx context.Context, userID int64, page *Paginator) error {
	return nil
}

// Trips .
func (UsersService) Trips(ctx context.Context, userID int64, page *Paginator) error {
	return nil
}
