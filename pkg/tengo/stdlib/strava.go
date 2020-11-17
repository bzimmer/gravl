package stdlib

import (
	"context"
	"encoding/json"

	"github.com/d5/tengo/v2"
	tjson "github.com/d5/tengo/v2/stdlib/json"

	"github.com/bzimmer/gravl/pkg/strava"
)

var stravaModule = map[string]tengo.Object{
	"service":    &tengo.UserFunction{Name: "service", Value: service},
	"athlete":    &tengo.UserFunction{Name: "athlete", Value: athlete},
	"activities": &tengo.UserFunction{Name: "activities", Value: activities},
	"activity":   &tengo.UserFunction{Name: "activity", Value: activity},
	"routes":     &tengo.UserFunction{Name: "routes", Value: routes},
}

type Strava struct {
	tengo.ObjectImpl
	Value *strava.Client
}

func (s *Strava) String() string {
	return "strava"
}

func (s *Strava) TypeName() string {
	return s.String()
}

func ToService(o tengo.Object) (v *strava.Client, ok bool) {
	switch o := o.(type) {
	case *Strava:
		v = o.Value
		ok = true
	}
	return
}

func ToTengo(v interface{}) (ret tengo.Object, err error) {
	b, err := json.Marshal(v)
	if err != nil {
		return
	}
	ret, err = tjson.Decode(b)
	return
}

func service(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) == 0 {
		err = tengo.ErrWrongNumArguments
		return
	}
	access, ok := tengo.ToString(args[0])
	if !ok {
		err = tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
		return
	}
	client, err := strava.NewClient(
		strava.WithAPICredentials(access, ""))
	if err != nil {
		return nil, err
	}
	return &Strava{Value: client}, nil
}

func athlete(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 0 {
		err = tengo.ErrWrongNumArguments
		return
	}
	service, ok := ToService(args[0])
	if !ok {
		err = tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "strava",
			Found:    args[0].TypeName(),
		}
		return
	}
	ath, err := service.Athlete.Athlete(context.Background())
	if err != nil {
		return nil, err
	}
	ret, err = ToTengo(ath)
	return
}

func activities(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 0 {
		err = tengo.ErrWrongNumArguments
		return
	}
	service, ok := ToService(args[0])
	if !ok {
		err = tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "strava",
			Found:    args[0].TypeName(),
		}
		return
	}
	acts, err := service.Activity.Activities(context.Background(), 10)
	if err != nil {
		return nil, err
	}
	ret, err = ToTengo(acts)
	return
}

func activity(args ...tengo.Object) (ret tengo.Object, err error) { // nolint
	if len(args) != 2 {
		err = tengo.ErrWrongNumArguments
		return
	}
	service, ok := ToService(args[0])
	if !ok {
		err = tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "strava",
			Found:    args[0].TypeName(),
		}
		return
	}
	id, ok := tengo.ToInt64(args[1])
	if !ok {
		err = tengo.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "int64",
			Found:    args[1].TypeName(),
		}
		return
	}
	act, err := service.Activity.Activity(context.Background(), id)
	if err != nil {
		return nil, err
	}
	ret, err = ToTengo(act)
	return
}

func routes(args ...tengo.Object) (ret tengo.Object, err error) { // nolint
	if len(args) != 2 {
		err = tengo.ErrWrongNumArguments
		return
	}
	service, ok := ToService(args[0])
	if !ok {
		err = tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "strava",
			Found:    args[0].TypeName(),
		}
		return
	}
	athleteID, ok := tengo.ToInt(args[1])
	if !ok {
		err = tengo.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "int",
			Found:    args[1].TypeName(),
		}
		return
	}
	rts, err := service.Route.Routes(context.Background(), athleteID)
	if err != nil {
		return nil, err
	}
	ret, err = ToTengo(rts)
	return
}
