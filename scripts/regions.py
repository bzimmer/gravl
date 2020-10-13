import re
import pprint

# Use the source from https://www.wta.org/go-outside/trip-reports

c = re.compile('<option value="(.*?)">(.*?)</option>')
regions = """
                <option value="b4845d8a21ad6a202944425c86b6e85f">Central Cascades</option>
              
                <option value="41f702968848492db697e10b14c14060">Central Washington</option>
              
                <option value="9d321b42e903a3224fd4fef44af9bee3">Eastern Washington</option>
              
                <option value="592fcc9afd9208db3b81fdf93dada567">Issaquah Alps</option>
              
                <option value="344281caae0d5e845a5003400c0be9ef">Mount Rainier Area</option>
              
                <option value="49aff77512c523f32ae13d889f6969c9">North Cascades</option>
              
                <option value="922e688d784aa95dfb80047d2d79dcf6">Olympic Peninsula</option>
              
                <option value="0c1d82b18f8023acb08e4daf03173e94">Puget Sound and Islands</option>
              
                <option value="04d37e830680c65b61df474e7e655d64">Snoqualmie Region</option>
              
                <option value="8a977ce4bf0528f4f833743e22acae5d">South Cascades</option>
              
                <option value="2b6f1470ed0a4735a4fc9c74e25096e0">Southwest Washington</option>
"""
subregions = {
'04d37e830680c65b61df474e7e655d64': {'5d45ee6e4b5b077d069382b6aac9d388': 'Snoqualmie Pass',
                                      '767d751df8b6495999e96486d4d32d49': 'Cle Elum Area',
                                      'db086e5e85941a02ae188f726f7e9e2c': 'North Bend Area',
                                      'f06510bd295c2d640ee2594d1b7a2ff6': 'Salmon La Sac/Teanaway'},
 '0c1d82b18f8023acb08e4daf03173e94': {'7e0a6ce03ba1204d6bd3fdf64d3ad805': 'San Juan Islands',
                                      '92528ff6af30075eec65f35159defc50': 'Bellingham Area',
                                      'd9dddf65d66479f065d40c1aeac18da3': 'Whidbey Island',
                                      'df2c2da1637452abe74a5d10837c2e03': 'Seattle-Tacoma Area'},
 '2b6f1470ed0a4735a4fc9c74e25096e0': {'35adb6fb84290947f778381d9d24a470': 'Lewis River Region',
                                      '74fe67d98acee7d0decda17c2441a4d2': 'Columbia River Gorge - WA',
                                      '7f27f833998f4bf38a1c816fdb4cac51': 'Columbia River Gorge - OR',
                                      '99c63956051e489f96527d6e4ed0915c': 'Vancouver Area',
                                      'f0095f8e8f394f4210f50999bf8abf2c': 'Long Beach Area'},
 '344281caae0d5e845a5003400c0be9ef': {'3b53cfc78db378ecf7599df0fa14a51c': 'SW - Longmire/Paradise',
                                      '883e708ab442592f904fd87c1c909f6b': 'Chinook Pass - Hwy 410',
                                      'c8814620167ff2c018e9b0d6e961f0c1': 'NE - Sunrise/White River',
                                      'cbe4acbaa2c01f9a5dbf4deece4e6ad9': 'NW - Carbon River/Mowich',
                                      'cc329f21ff637826168e61bc9db77d65': 'SE - Cayuse Pass/Stevens Canyon'},
 '41f702968848492db697e10b14c14060': {'39c3197019b541c4b4b2970281fd85ed': 'Grand Coulee',
                                      '3a415482d61dc1f893288c1bcf5cd8ae': 'Potholes Region',
                                      '6febb0079f87770b5e790a71aafa3770': 'Tri-Cities',
                                      '7a4397f2ffde7d0b490ff1ca77cceb9e': 'Yakima',
                                      'b4be8a42f05d2054cbb2ca031b9a6a03': 'Wenatchee'},
 '49aff77512c523f32ae13d889f6969c9': {'425fd9e8fd7edb23fc53782f16c2ea05': 'Pasayten',
                                      '5674705352f9b856f2df1da7cbb8e0b1': 'Mount Baker Area',
                                      '5952810559bc6d85f808011a53ea6fcf': 'Methow/Sawtooth',
                                      '6194b417d1ae41b1ecd0d297b3fd2dea': 'Mountain Loop Highway',
                                      'b52b426625b55325e408adfacae3b6c5': 'North Cascades Highway - Hwy 20'},
 '592fcc9afd9208db3b81fdf93dada567': {'325fdb0c3072b1b9acca522fb9e69ec2': 'Cougar Mountain',
                                      '70056b3c13ba158deec7750ef9701a94': 'Squak Mountain',
                                      '98ea7186f4ff4bbba5329613f9a89bfb': 'Taylor Mountain',
                                      '9f13d8a3fcd2e1ab7a5b5aaab5997a9e': 'Tiger Mountain'},
 '8a977ce4bf0528f4f833743e22acae5d': {'17dcd22410be73abfd45d2703e123a35': 'Mount St. Helens',
                                      '6f227fc5711324cee6170aa6d4b52cec': 'White Pass/Cowlitz River Valley',
                                      '73b109ca9145e4433f3089a1789d29bf': 'Mount Adams Area',
                                      'ac28fa6c89800fca796a2e61b879f416': 'Dark Divide',
                                      'b1376aba679a6bf3d6402cf91f16a44e': 'Goat Rocks'},
 '922e688d784aa95dfb80047d2d79dcf6': {'3b2e0197b8be6c19a273c919b3301405': 'Kitsap Peninsula',
                                      '3ca3cd096bfedde6ff95b0859278cc75': 'Olympia',
                                      '6135b6a861b5ac0c9b17a2f9b60c9295': 'Pacific Coast',
                                      'bfbc0abe0fd04783aaad717ea2699866': 'Hood Canal',
                                      'e4421728558408ef04e0a46afb2aa7ea': 'Northern Coast'},
 '9d321b42e903a3224fd4fef44af9bee3': {'3eff611193d7d4b57590df1f40b48800': 'Palouse and Blue Mountains',
                                      'bec6f9858a88f32a0912ed21d9c63b51': "Spokane Area/Coeur d'Alene",
                                      'd305615b5db417f18661c5233d2ce950': 'Selkirk Range',
                                      'fe742c316d095b81d23d712efa977d3d': 'Okanogan Highlands/Kettle River Range'},
 'b4845d8a21ad6a202944425c86b6e85f': {'2b0ca41464d9baca77ced16fb4d40760': 'Entiat Mountains/Lake Chelan',
                                      '637634387ca38685f89162475c7fc1d2': 'Stevens Pass - West',
                                      '684278bc46c11ebe3c5b7212b6f8e486': 'Leavenworth Area',
                                      '83c2ab06fbf236015c8848042f706d58': 'Blewett Pass'},
 }


# Want this:
# var Regions *[]Region = *[]Region


# From this:
# [{'id': 'b4845d8a21ad6a202944425c86b6e85f',
#   'name': 'Central Cascades',
#   'subregions': [{'id': '2b0ca41464d9baca77ced16fb4d40760', 'name': 'Entiat Mountains/Lake Chelan'},
#                  {'id': '637634387ca38685f89162475c7fc1d2', 'name': 'Stevens Pass - West'},
#                  {'id': '684278bc46c11ebe3c5b7212b6f8e486', 'name': 'Leavenworth Area'},
#                  {'id': '83c2ab06fbf236015c8848042f706d58', 'name': 'Blewett Pass'}]},
#  {'id': '41f702968848492db697e10b14c14060',
#   'name': 'Central Washington',
#   'subregions': [{'id': '39c3197019b541c4b4b2970281fd85ed', 'name': 'Grand Coulee'},
#                  {'id': '3a415482d61dc1f893288c1bcf5cd8ae', 'name': 'Potholes Region'},
#                  {'id': '6febb0079f87770b5e790a71aafa3770', 'name': 'Tri-Cities'},
#                  {'id': '7a4397f2ffde7d0b490ff1ca77cceb9e', 'name': 'Yakima'},
#                  {'id': 'b4be8a42f05d2054cbb2ca031b9a6a03', 'name': 'Wenatchee'}]}]

# To this:
# var Regions = &[11]Region{
# 	Region{ID: "b4845d8a21ad6a202944425c86b6e85f", Name: "Central Cascades",
# 		Subregions: &[]Region{
# 			Region{ID: "2b0ca41464d9baca77ced16fb4d40760", Name: "Entiat Mountains/Lake Chelan"},
# 			Region{ID: "637634387ca38685f89162475c7fc1d2", Name: "Stevens Pass - West"},
# 			Region{ID: "684278bc46c11ebe3c5b7212b6f8e486", Name: "Leavenworth Area"},
# 			Region{ID: "83c2ab06fbf236015c8848042f706d58", Name: "Blewett Pass"},
# 		}},
# 	Region{ID: "41f702968848492db697e10b14c14060", Name: "Central Washington",
# 		Subregions: &[]Region{
# 			Region{ID: "39c3197019b541c4b4b2970281fd85ed", Name: "Grand Coulee"},
# 			Region{ID: "3a415482d61dc1f893288c1bcf5cd8ae", Name: "Potholes Region"},
# 			Region{ID: "6febb0079f87770b5e790a71aafa3770", Name: "Tri-Cities"},
# 			Region{ID: "7a4397f2ffde7d0b490ff1ca77cceb9e", Name: "Yakima"},
# 			Region{ID: "b4be8a42f05d2054cbb2ca031b9a6a03", Name: "Wenatchee"},
# 		}},
# }

m = []
for key, name in c.findall(regions):
    sr = [{"id":k, "name":v} for k, v in subregions[key].items()]
    x = {"id": key, "name": name, "subregions":sr}
    m.append(x)
# pprint.pprint(m, width=150)

print("var Regions = &[%d]Region{" % len(m))
for x in m:
    print('Region{ID:"%s", Name: "%s",' % (x["id"], x["name"]))
    print('  Subregions: &[]Region{')
    for s in x["subregions"]:
        print('    Region{ID:"%s", Name: "%s"},' % (s["id"], s["name"]))
    print("}},")
print("}")
