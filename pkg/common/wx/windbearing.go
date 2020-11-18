package wx

import "fmt"

// WindBearing .
func WindBearing(bearing string) (float64, error) { // nolint:funlen,gocyclo
	switch bearing {
	case "N": // North
		return 0.0, nil
	case "NbE": // North by east
		return 11.25, nil
	case "NNE": // North-northeast
		return 22.5, nil
	case "NEbN": // Northeast by north
		return 33.75, nil
	case "NE": // Northeast
		return 45.0, nil
	case "NEbE": // Northeast by east
		return 56.25, nil
	case "ENE": // East-northeast
		return 67.5, nil
	case "EbN": // East by north
		return 78.75, nil
	case "E": // East
		return 90.0, nil
	case "EbS": // East by south
		return 101.25, nil
	case "ESE": // East-southeast
		return 112.5, nil
	case "SEbE": // Southeast by east
		return 123.75, nil
	case "SE": // Southeast
		return 135.0, nil
	case "SEbS": // Southeast by south
		return 146.25, nil
	case "SSE": // South-southeast
		return 157.5, nil
	case "SbE": // South by east
		return 168.75, nil
	case "S": // South
		return 180.0, nil
	case "SbW": // South by west
		return 191.25, nil
	case "SSW": // South-southwest
		return 202.5, nil
	case "SWbS": // Southwest by south
		return 213.75, nil
	case "SW": // Southwest
		return 225.0, nil
	case "SWbW": // Southwest by west
		return 236.25, nil
	case "WSW": // West-southwest
		return 247.5, nil
	case "WbS": // West by south
		return 258.75, nil
	case "W": // West
		return 270.0, nil
	case "WbN": // West by north
		return 281.25, nil
	case "WNW": // West-northwest
		return 292.5, nil
	case "NWbW": // Northwest by west
		return 303.75, nil
	case "NW": // Northwest
		return 315.0, nil
	case "NWbN": // Northwest by north
		return 326.25, nil
	case "NNW": // North-northwest
		return 337.5, nil
	case "NbW": // North by west
		return 348.75, nil
	default:
		return 0.0, fmt.Errorf("unknown bearing: %s", bearing)
	}
}
