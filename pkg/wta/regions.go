package wta

// Regions .
var Regions = &[11]Region{
	Region{ID: "b4845d8a21ad6a202944425c86b6e85f", Name: "Central Cascades",
		Subregions: &[]Region{
			Region{ID: "2b0ca41464d9baca77ced16fb4d40760", Name: "Entiat Mountains/Lake Chelan"},
			Region{ID: "637634387ca38685f89162475c7fc1d2", Name: "Stevens Pass - West"},
			Region{ID: "684278bc46c11ebe3c5b7212b6f8e486", Name: "Leavenworth Area"},
			Region{ID: "83c2ab06fbf236015c8848042f706d58", Name: "Blewett Pass"},
		}},
	Region{ID: "41f702968848492db697e10b14c14060", Name: "Central Washington",
		Subregions: &[]Region{
			Region{ID: "39c3197019b541c4b4b2970281fd85ed", Name: "Grand Coulee"},
			Region{ID: "3a415482d61dc1f893288c1bcf5cd8ae", Name: "Potholes Region"},
			Region{ID: "6febb0079f87770b5e790a71aafa3770", Name: "Tri-Cities"},
			Region{ID: "7a4397f2ffde7d0b490ff1ca77cceb9e", Name: "Yakima"},
			Region{ID: "b4be8a42f05d2054cbb2ca031b9a6a03", Name: "Wenatchee"},
		}},
	Region{ID: "9d321b42e903a3224fd4fef44af9bee3", Name: "Eastern Washington",
		Subregions: &[]Region{
			Region{ID: "3eff611193d7d4b57590df1f40b48800", Name: "Palouse and Blue Mountains"},
			Region{ID: "bec6f9858a88f32a0912ed21d9c63b51", Name: "Spokane Area/Coeur d'Alene"},
			Region{ID: "d305615b5db417f18661c5233d2ce950", Name: "Selkirk Range"},
			Region{ID: "fe742c316d095b81d23d712efa977d3d", Name: "Okanogan Highlands/Kettle River Range"},
		}},
	Region{ID: "592fcc9afd9208db3b81fdf93dada567", Name: "Issaquah Alps",
		Subregions: &[]Region{
			Region{ID: "325fdb0c3072b1b9acca522fb9e69ec2", Name: "Cougar Mountain"},
			Region{ID: "70056b3c13ba158deec7750ef9701a94", Name: "Squak Mountain"},
			Region{ID: "98ea7186f4ff4bbba5329613f9a89bfb", Name: "Taylor Mountain"},
			Region{ID: "9f13d8a3fcd2e1ab7a5b5aaab5997a9e", Name: "Tiger Mountain"},
		}},
	Region{ID: "344281caae0d5e845a5003400c0be9ef", Name: "Mount Rainier Area",
		Subregions: &[]Region{
			Region{ID: "3b53cfc78db378ecf7599df0fa14a51c", Name: "SW - Longmire/Paradise"},
			Region{ID: "883e708ab442592f904fd87c1c909f6b", Name: "Chinook Pass - Hwy 410"},
			Region{ID: "c8814620167ff2c018e9b0d6e961f0c1", Name: "NE - Sunrise/White River"},
			Region{ID: "cbe4acbaa2c01f9a5dbf4deece4e6ad9", Name: "NW - Carbon River/Mowich"},
			Region{ID: "cc329f21ff637826168e61bc9db77d65", Name: "SE - Cayuse Pass/Stevens Canyon"},
		}},
	Region{ID: "49aff77512c523f32ae13d889f6969c9", Name: "North Cascades",
		Subregions: &[]Region{
			Region{ID: "425fd9e8fd7edb23fc53782f16c2ea05", Name: "Pasayten"},
			Region{ID: "5674705352f9b856f2df1da7cbb8e0b1", Name: "Mount Baker Area"},
			Region{ID: "5952810559bc6d85f808011a53ea6fcf", Name: "Methow/Sawtooth"},
			Region{ID: "6194b417d1ae41b1ecd0d297b3fd2dea", Name: "Mountain Loop Highway"},
			Region{ID: "b52b426625b55325e408adfacae3b6c5", Name: "North Cascades Highway - Hwy 20"},
		}},
	Region{ID: "922e688d784aa95dfb80047d2d79dcf6", Name: "Olympic Peninsula",
		Subregions: &[]Region{
			Region{ID: "3b2e0197b8be6c19a273c919b3301405", Name: "Kitsap Peninsula"},
			Region{ID: "3ca3cd096bfedde6ff95b0859278cc75", Name: "Olympia"},
			Region{ID: "6135b6a861b5ac0c9b17a2f9b60c9295", Name: "Pacific Coast"},
			Region{ID: "bfbc0abe0fd04783aaad717ea2699866", Name: "Hood Canal"},
			Region{ID: "e4421728558408ef04e0a46afb2aa7ea", Name: "Northern Coast"},
		}},
	Region{ID: "0c1d82b18f8023acb08e4daf03173e94", Name: "Puget Sound and Islands",
		Subregions: &[]Region{
			Region{ID: "7e0a6ce03ba1204d6bd3fdf64d3ad805", Name: "San Juan Islands"},
			Region{ID: "92528ff6af30075eec65f35159defc50", Name: "Bellingham Area"},
			Region{ID: "d9dddf65d66479f065d40c1aeac18da3", Name: "Whidbey Island"},
			Region{ID: "df2c2da1637452abe74a5d10837c2e03", Name: "Seattle-Tacoma Area"},
		}},
	Region{ID: "04d37e830680c65b61df474e7e655d64", Name: "Snoqualmie Region",
		Subregions: &[]Region{
			Region{ID: "5d45ee6e4b5b077d069382b6aac9d388", Name: "Snoqualmie Pass"},
			Region{ID: "767d751df8b6495999e96486d4d32d49", Name: "Cle Elum Area"},
			Region{ID: "db086e5e85941a02ae188f726f7e9e2c", Name: "North Bend Area"},
			Region{ID: "f06510bd295c2d640ee2594d1b7a2ff6", Name: "Salmon La Sac/Teanaway"},
		}},
	Region{ID: "8a977ce4bf0528f4f833743e22acae5d", Name: "South Cascades",
		Subregions: &[]Region{
			Region{ID: "17dcd22410be73abfd45d2703e123a35", Name: "Mount St. Helens"},
			Region{ID: "6f227fc5711324cee6170aa6d4b52cec", Name: "White Pass/Cowlitz River Valley"},
			Region{ID: "73b109ca9145e4433f3089a1789d29bf", Name: "Mount Adams Area"},
			Region{ID: "ac28fa6c89800fca796a2e61b879f416", Name: "Dark Divide"},
			Region{ID: "b1376aba679a6bf3d6402cf91f16a44e", Name: "Goat Rocks"},
		}},
	Region{ID: "2b6f1470ed0a4735a4fc9c74e25096e0", Name: "Southwest Washington",
		Subregions: &[]Region{
			Region{ID: "35adb6fb84290947f778381d9d24a470", Name: "Lewis River Region"},
			Region{ID: "74fe67d98acee7d0decda17c2441a4d2", Name: "Columbia River Gorge - WA"},
			Region{ID: "7f27f833998f4bf38a1c816fdb4cac51", Name: "Columbia River Gorge - OR"},
			Region{ID: "99c63956051e489f96527d6e4ed0915c", Name: "Vancouver Area"},
			Region{ID: "f0095f8e8f394f4210f50999bf8abf2c", Name: "Long Beach Area"},
		}},
}

// RegionsService .
type RegionsService service
