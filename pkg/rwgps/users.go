package rwgps

import (
	"context"
)

// UsersService .
type UsersService service

// AuthenticatedUser .
func (UsersService) AuthenticatedUser(ctx context.Context) (*User, error) {
	return nil, nil
}

// // Routes .
// func (s *TripsService) Routes(ctx context.Context, userID int64, page *Paginator) error {
// 	return nil
// }

// // Trips .
// func (s *TripsService) Trips(ctx context.Context, userID int64, page *Paginator) error {
// 	return nil
// }
