import pandas as pd
import json

headings = ["id","ident","type","elevation_ft","continent","iso_country","iso_region","municipality","scheduled_service","gps_code","home_link","wikipedia_link","keywords"]
df = pd.read_csv('airports.csv')
for h in headings:
    df = df.drop(h, 1)

airports = {}
for index, row in df.iterrows():
    a = {
        'name': row['name'],
        'latitude_deg': row['latitude_deg'],
        'longitude_deg': row['longitude_deg'],
    }
    airports[row['iata_code']] = a

with open('airports.json', 'w') as w_file:
     w_file.write(json.dumps(airports))