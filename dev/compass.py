import re
import csv
import io

rose = """No.	Compass point	Abbreviation	Traditional wind point	Minimum	Middle azimuth	Maximum
0	North	N	Tramontana	354.375	0.000	5.625
1	North by east	NbE	Quarto di Tramontana verso Greco	5.625	11.250	16.875
2	North-northeast	NNE	Greco-Tramontana	16.875	22.500	28.125
3	Northeast by north	NEbN	Quarto di Greco verso Tramontana	28.125	33.750	39.375
4	Northeast	NE	Greco	39.375	45.000	50.625
5	Northeast by east	NEbE	Quarto di Greco verso Levante	50.625	56.250	61.875
6	East-northeast	ENE	Greco-Levante	61.875	67.500	73.125
7	East by north	EbN	Quarto di Levante verso Greco	73.125	78.750	84.375
8	East	E	Levante	84.375	90.000	95.625
9	East by south	EbS	Quarto di Levante verso Scirocco	95.625	101.250	106.875
10	East-southeast	ESE	Levante-Scirocco	106.875	112.500	118.125
11	Southeast by east	SEbE	Quarto di Scirocco verso Levante	118.125	123.750	129.375
12	Southeast	SE	Scirocco	129.375	135.000	140.625
13	Southeast by south	SEbS	Quarto di Scirocco verso Ostro	140.625	146.250	151.875
14	South-southeast	SSE	Ostro-Scirocco	151.875	157.500	163.125
15	South by east	SbE	Quarto di Ostro verso Scirocco	163.125	168.750	174.375
16	South	S	Ostro	174.375	180.000	185.625
17	South by west	SbW	Quarto di Ostro verso Libeccio	185.625	191.250	196.875
18	South-southwest	SSW	Ostro-Libeccio	196.875	202.500	208.125
19	Southwest by south	SWbS	Quarto di Libeccio verso Ostro	208.125	213.750	219.375
20	Southwest	SW	Libeccio	219.375	225.000	230.625
21	Southwest by west	SWbW	Quarto di Libeccio verso Ponente	230.625	236.250	241.875
22	West-southwest	WSW	Ponente-Libeccio	241.875	247.500	253.125
23	West by south	WbS	Quarto di Ponente verso Libeccio	253.125	258.750	264.375
24	West	W	Ponente	264.375	270.000	275.625
25	West by north	WbN	Quarto di Ponente verso Maestro	275.625	281.250	286.875
26	West-northwest	WNW	Maestro-Ponente	286.875	292.500	298.125
27	Northwest by west	NWbW	Quarto di Maestro verso Ponente	298.125	303.750	309.375
28	Northwest	NW	Maestro	309.375	315.000	320.625
29	Northwest by north	NWbN	Quarto di Maestro verso Tramontana	320.625	326.250	331.875
30	North-northwest	NNW	Maestro-Tramontana	331.875	337.500	343.125
31	North by west	NbW	Quarto di Tramontana verso Maestro	343.125	348.750	354.375
"""

rows = iter(csv.reader(io.StringIO(rose), delimiter="\t"))
next(rows)

print("// WindBearing .")
print("func WindBearing(bearing string)(float64, error) {")
print("switch(bearing) {")
for row in rows:
    name = row[2]
    bounds = [float(x) for x in row[4:6]]
    print(" case \"%s\": // %s" % (name, row[1]))
    print(" return %s, nil" % bounds[1])
print(" default:")
print(' return 0.0, fmt.Errorf("unknown bearing: %s", bearing)')
print("}")
print("}")
