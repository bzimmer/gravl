package zwift

// Fault represents a Zwift error
type Fault struct {
	Message string `json:"message"`
}

func (f *Fault) Error() string {
	return f.Message
}

type Privacy struct {
	ApprovalRequired             bool   `json:"approvalRequired"`
	DisplayWeight                bool   `json:"displayWeight"`
	Minor                        bool   `json:"minor"`
	PrivateMessaging             bool   `json:"privateMessaging"`
	DefaultFitnessDataPrivacy    bool   `json:"defaultFitnessDataPrivacy"`
	SuppressFollowerNotification bool   `json:"suppressFollowerNotification"`
	DisplayAge                   bool   `json:"displayAge"`
	DefaultActivityPrivacy       string `json:"defaultActivityPrivacy"`
}

type SocialFacts struct {
	ProfileID                           int64  `json:"profileId"`
	FollowersCount                      int    `json:"followersCount"`
	FolloweesCount                      int    `json:"followeesCount"`
	FolloweesInCommonWithLoggedInPlayer int    `json:"followeesInCommonWithLoggedInPlayer"`
	FollowerStatusOfLoggedInPlayer      string `json:"followerStatusOfLoggedInPlayer"`
	FolloweeStatusOfLoggedInPlayer      string `json:"followeeStatusOfLoggedInPlayer"`
	IsFavoriteOfLoggedInPlayer          bool   `json:"isFavoriteOfLoggedInPlayer"`
}

type Profile struct {
	ID                           int64        `json:"id"`
	PublicID                     string       `json:"publicId"`
	FirstName                    string       `json:"firstName"`
	LastName                     string       `json:"lastName"`
	Male                         bool         `json:"male"`
	ImageSrc                     string       `json:"imageSrc"`
	ImageSrcLarge                string       `json:"imageSrcLarge"`
	PlayerType                   string       `json:"playerType"`
	CountryAlpha3                string       `json:"countryAlpha3"`
	CountryCode                  int          `json:"countryCode"`
	UseMetric                    bool         `json:"useMetric"`
	Riding                       bool         `json:"riding"`
	Privacy                      *Privacy     `json:"privacy"`
	SocialFacts                  *SocialFacts `json:"socialFacts"`
	WorldID                      int64        `json:"worldId"`
	EnrolledZwiftAcademy         bool         `json:"enrolledZwiftAcademy"`
	PlayerTypeID                 int64        `json:"playerTypeId"`
	PlayerSubTypeID              int64        `json:"playerSubTypeId"`
	CurrentActivityID            int64        `json:"currentActivityId"`
	Address                      string       `json:"address"`
	Age                          int          `json:"age"`
	BodyType                     int          `json:"bodyType"`
	ConnectedToStrava            bool         `json:"connectedToStrava"`
	ConnectedToTrainingPeaks     bool         `json:"connectedToTrainingPeaks"`
	ConnectedToTodaysPlan        bool         `json:"connectedToTodaysPlan"`
	ConnectedToUnderArmour       bool         `json:"connectedToUnderArmour"`
	ConnectedToWithings          bool         `json:"connectedToWithings"`
	ConnectedToFitbit            bool         `json:"connectedToFitbit"`
	ConnectedToGarmin            bool         `json:"connectedToGarmin"`
	ConnectedToRuntastic         bool         `json:"connectedToRuntastic"`
	ConnectedToZwiftPower        bool         `json:"connectedToZwiftPower"`
	StravaPremium                bool         `json:"stravaPremium"`
	Bt                           string       `json:"bt"`
	BirthDate                    string       `json:"dob"`
	EmailAddress                 string       `json:"emailAddress"`
	Height                       int          `json:"height"`
	Location                     string       `json:"location"`
	PreferredLanguage            string       `json:"preferredLanguage"`
	MixpanelDistinctID           string       `json:"mixpanelDistinctId"`
	ProfileChanges               bool         `json:"profileChanges"`
	Weight                       int          `json:"weight"`
	B                            bool         `json:"b"`
	CreatedOn                    string       `json:"createdOn"`
	Source                       string       `json:"source"`
	Origin                       string       `json:"origin"`
	LaunchedGameClient           string       `json:"launchedGameClient"`
	FTP                          int          `json:"ftp"`
	UserAgent                    string       `json:"userAgent"`
	RunTime1MiInSeconds          int          `json:"runTime1miInSeconds"`
	RunTime5KmInSeconds          int          `json:"runTime5kmInSeconds"`
	RunTime10KmInSeconds         int          `json:"runTime10kmInSeconds"`
	RunTimeHalfMarathonInSeconds int          `json:"runTimeHalfMarathonInSeconds"`
	RunTimeFullMarathonInSeconds int          `json:"runTimeFullMarathonInSeconds"`
	CyclingOrganization          string       `json:"cyclingOrganization"`
	LicenseNumber                string       `json:"licenseNumber"`
	BigCommerceID                string       `json:"bigCommerceId"`
	AchievementLevel             int          `json:"achievementLevel"`
	TotalDistance                int          `json:"totalDistance"`
	TotalDistanceClimbed         int          `json:"totalDistanceClimbed"`
	TotalTimeInMinutes           int          `json:"totalTimeInMinutes"`
	TotalInKOMJersey             int          `json:"totalInKomJersey"`
	TotalInSprintersJersey       int          `json:"totalInSprintersJersey"`
	TotalInOrangeJersey          int          `json:"totalInOrangeJersey"`
	TotalWattHours               int          `json:"totalWattHours"`
	TotalExperiencePoints        int          `json:"totalExperiencePoints"`
	TotalGold                    int          `json:"totalGold"`
	RunAchievementLevel          int          `json:"runAchievementLevel"`
	TotalRunDistance             int          `json:"totalRunDistance"`
	TotalRunTimeInMinutes        int          `json:"totalRunTimeInMinutes"`
	TotalRunExperiencePoints     int          `json:"totalRunExperiencePoints"`
	TotalRunCalories             int          `json:"totalRunCalories"`
	PowerSourceType              string       `json:"powerSourceType"`
	PowerSourceModel             string       `json:"powerSourceModel"`
	VirtualBikeModel             string       `json:"virtualBikeModel"`
	NumberOfFolloweesInCommon    int          `json:"numberOfFolloweesInCommon"`
	Affiliate                    string       `json:"affiliate"`
	AvantlinkID                  string       `json:"avantlinkId"`
	FundraiserID                 string       `json:"fundraiserId"`
	// PublicAttributes             interface{} `json:"publicAttributes"`
	// PrivateAttributes            interface{} `json:"privateAttributes"`
}

type Activity struct {
	IDStr                string   `json:"id_str"`
	ID                   int64    `json:"id"`
	ProfileID            int64    `json:"profileId"`
	Profile              *Profile `json:"profile"`
	WorldID              int64    `json:"worldId"`
	Name                 string   `json:"name"`
	Description          string   `json:"description"`
	PrivateActivity      bool     `json:"privateActivity"`
	Sport                string   `json:"sport"`
	StartDate            string   `json:"startDate"`
	EndDate              string   `json:"endDate"`
	LastSaveDate         string   `json:"lastSaveDate"`
	AutoClosed           bool     `json:"autoClosed"`
	Duration             string   `json:"duration"`
	DistanceInMeters     float64  `json:"distanceInMeters"`
	FitFileBucket        string   `json:"fitFileBucket"`
	FitFileKey           string   `json:"fitFileKey"`
	TotalElevation       float64  `json:"totalElevation"`
	AvgWatts             float64  `json:"avgWatts"`
	RideOnGiven          bool     `json:"rideOnGiven"`
	ActivityRideOnCount  int      `json:"activityRideOnCount"`
	ActivityCommentCount int      `json:"activityCommentCount"`
	Calories             float64  `json:"calories"`
	PrimaryImageURL      string   `json:"primaryImageUrl"`
	MovingTimeInMillis   int      `json:"movingTimeInMs"`
	Privacy              string   `json:"privacy"`
	// SnapshotList         interface{} `json:"snapshotList"`
}
