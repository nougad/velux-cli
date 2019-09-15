# Velux Protocol

The Velux API is mostly compatible with the [Netatmo API](https://dev.netatmo.com/en-US/resources/technical/introduction) (swagger can be found [here](https://cbornet.github.io/netatmo-swagger-decl/)) but all Velux relevant calls are in a second API (`/syncapi/v1/setstate`).

## Initial login

Request:

```
curl -d "grant_type=password&client_id=${CLIENT_ID}&client_secret=${CLIENT_SECRET}&username=${USERNAME}&password=${PASSWORD}&user_prefix=velux" https://app.velux-active.com/oauth2/token
```

Response:

```
{
  "access_token": "...",
  "refresh_token": "...",
  "scope": [
    "all_scopes"
  ],
  "expires_in": 10800,
  "expire_in": 10800
}
```

Available scopes:

* `all_scopes`
* `access_velux`
* `read_velux`
* `write_velux`
* most likely more

Available device types:

* NXO - RollerShutter
* NXG - Bridge
* NXD - Depature switch
* NXS - Sensor

## Refresh token

Request:

```
curl -d "grant_type=refresh_token&refresh_token=${REFRESH_TOKEN}&client_id=${CLIENT_ID}&client_secret=${CLIENT_SECRET}" https://app.velux-active.com/oauth2/token
```

Response:

```
{
  "access_token": "...",
  "refresh_token": "...",
  "scope": [
    "all_scopes"
  ],
  "expires_in": 10800,
  "expire_in": 10800
}
```

## GetUser

Request:

```
curl -d "access_token=${YOUR_TOKEN}" https://app.velux-active.com/api/getuser
```

Response:

```
{
  "body": {
    "_id": "${ID}",
    "not_associable": false,
    "mail": "${EMAIL}",
    "account_validation": {
      "validated_mail": true,
      "validation_date": 1662521387
    },
    "fb_chatbot_available": true,
    "app_telemetry": true,
    "administrative": {
      "unit": 0,
      "windunit": 0,
      "pressureunit": 0,
      "feel_like_algo": 0,
      "reg_locale": "en-GB",
      "lang": "en"
    }
  },
  "status": "ok",
  "time_exec": 0.020205020904541,
  "time_server": 1571429393
}
```

## GetHomeData

Request:

```
curl -v -X POST -d "access_token=${YOUR_TOKEN}" https://app.velux-active.com/api/gethomedata
```

Response:

```
{
  "body": {
    "homes": [
      {
        "id": "${HOMEID}",
        "name": "${HOMENAME}",
        "share_info": false,
        "gone_after": 14400,
        "smart_notifs": true,
        "notify_movements": "empty",
        "record_movements": "empty",
        "notify_unknowns": "empty",
        "record_alarms": "always",
        "record_animals": true,
        "notify_animals": true,
        "events_ttl": "never",
        "persons": [],
        "record_humans": "empty",
        "notify_humans": "empty",
        "presence_record_humans": "record_and_notify",
        "presence_record_vehicles": "record_and_notify",
        "presence_record_animals": "record",
        "presence_record_alarms": "record",
        "presence_record_movements": "record",
        "presence_notify_from": 0,
        "presence_notify_to": 86399,
        "presence_enable_notify_from_to": "empty",
        "place": {
          "altitude": 42,
          "city": "Berlin",
          "country": "DE",
          "location": [
            13.3,
            52.5
          ],
          "timezone": "Europe/Berlin",
          "trust_location": true
        },
        "cameras": [],
        "smokedetectors": [],
        "admin_access_code": [
          "....."
        ]
      }
    ],
    "user": {
      "reg_locale": "en-GB",
      "lang": "en",
      "country": "DE",
      "mail": "${EMAIL}",
      "fb_chatbot_available": true,
      "app_telemetry": true
    },
    "global_info": {
      "show_tags": true
    }
  },
  "status": "ok",
  "time_exec": 0.013638019561768,
  "time_server": 1561413743
}
```

## HomesData

Request:

```
curl -v -X POST -d "access_token=${YOUR_TOKEN}&gateway_types=[NXG]" https://app.velux-active.com/api/homesdata
```

## HomeStatus

Request:

```
curl -v -X POST -d "access_token=${YOUR_TOKEN}&home_id=${HOMEID}" https://app.velux-active.com/api/homestatus
```

alias for `/syncapi/v1/homestatus`

## GetHomeUsers

Request:

```
curl -v -X POST -d "access_token=${YOUR_TOKEN}&home_id=${HOMEID}" https://app.netatmo.net/api/gethomeusers
```

## SetState

```
curl -s -H "Content-Type: application/json; charset=utf-8" -H "Authorization: Bearer $YOUR_TOKEN" -X POST -d "${JSON}" "https://app.velux-active.com/syncapi/v1/setstate"
```

### Shutter position

```
{
    "home": {
        "id": "${HOMEID}",
        "modules": [
            {
                "bridge": "${BRIDGE}",
                "id": "${MODULEID}",
                "target_position": 100
            }
        ]
    }
}
```
### Retrieve Key

```
{
    "home": {
        "id": "${HOMEID}",
        "modules": [
            {
                "id": "${BRIDGE}",
                "retrieve_key": true
            }
        ]
    }
}
```

### Stop Movements

```
{
    "home": {
        "id": "${HOMEID}",
        "modules": [
            {
                "id": "${BRIDGE}",
                "stop_movements": "all"
            }
        ]
    }
}
```

### WakeUp

```
{
    "home": {
        "id": "${HOMEID}",
        "modules": [
            {
                "id": "${BRIDGE}",
                "scenario": "wake_up"
            }
        ]
    }
}
```

### BedTime

```
{
    "home": {
        "id": "${HOMEID}",
        "modules": [
            {
                "id": "${BRIDGE}",
                "scenario": "bedtime"
            }
        ]
    }
}
```

### Identity

Make device blink

```
{
    "home": {
        "id": "${HOMEID}",
        "modules": [
            {
                "id": "${BRIDGE}",
                "identify": true
            }
        ]
    }
}
```

## GetMeasure

sensors: https://dev.netatmo.com/resources/technical/reference/common/getmeasure

all values:

```
curl -v -X POST -d "access_token=${YOUR_TOKEN}&device_id=${BRIDGE}&module_id=${DEV1}&scale=30min&type=Temperature,CO2,Humidity,Pressure,Noise,Rain,WindStrength,WindAngle,Guststrength,GustAngle,Sp_Temperature,BoilerOn,BoilerOff,min_temp,max_temp,min_hum,max_hum,min_pressure,max_pressure,min_noise,max_noise,sum_rain,sum_boiler_on,sum_boiler_off" https://app.velux-active.com/api/getmeasure
```

only works for Humidity, CO2:

```
curl -v -X POST -d "access_token=${YOUR_TOKEN}&device_id=${BRIDGE}&module_id=${DEV1}&scale=30min&type=CO2,Humidity&date_begin=1561834800" https://app.velux-active.com/api/getmeasure
```

and when used without a module for Pressure:

```
curl -v -X POST -d "access_token=${YOUR_TOKEN}&device_id=${BRIDGE}&scale=30min&type=Pressure&date_begin=1561834800" https://app.velux-active.com/api/getmeasure
```

## Other Calls

### SetPersonsAway

Request:

```
curl -v -X POST -d "access_token=${YOUR_TOKEN}&home_id=${HOMEID}" https://app.velux-active.com/api/setpersonsaway
```

not sure what it does

### AddPushContext

https://app.velux-active.com/api/addpushcontext

### CheckSession

https://app.velux-active.com/api/checksession

### CreateNewhomeSchedule

https://app.velux-active.com/api/createnewhomeschedule

### GenerateHomeAdminAccesscode

https://app.velux-active.com/api/generatehomeadminaccesscode

### GetAlgorithmsActivity

https://app.velux-active.com/api/getalgorithmsactivity

## GetConfigs

Request:

```
curl -v -X POST -d "access_token=${YOUR_TOKEN}&home_id=${HOMEID}" https://app.netatmo.net/syncapi/v1/getconfigs
```

### GetNotificationSettings

https://app.velux-active.com/api/getnotificationsettings

## GetScenarios

Request:

```
curl -v -X POST -d "access_token=${YOUR_TOKEN}&home_id=${HOMEID}" https://app.netatmo.net/syncapi/v1/getscenarios
```

### MailNotificationConfiguration

https://app.velux-active.com/api/mailnotificationconfiguration

### ModifyUser

https://app.velux-active.com/api/modifyuser

### RemoveAccessTokens

https://app.velux-active.com/api/removeaccesstokens

### RemoveUserAccessToHome

https://app.velux-active.com/api/removeuseraccesstohome

### SendMailApproveMail

https://app.velux-active.com/api/sendmailapprovemail

### SetConfigs

https://app.velux-active.com/syncapi/v1/setconfigs

### SetNotificationSettings

https://app.velux-active.com/api/setnotificationsettings

### SetScenarios

* https://app.velux-active.com/syncapi/v1/setscenarios

### SwitchHomeSchedule

* https://app.velux-active.com/api/switchhomeschedule

### SynchomeSchedule

* https://app.velux-active.com/api/synchomeschedule

### UpdateHomePlace

* https://app.velux-active.com/api/updatehomeplace

## Netatmo calls not working with velux

* https://app.velux-active.com/api/gethomecoachsdata
* https://app.velux-active.com/api/getstationsdata
* https://app.velux-active.com/api/getthermstate
* https://app.velux-active.com/api/devicelist
* https://app.velux-active.com/api/partnerdevices - works but empty
* https://app.velux-active.com/api/getthermostatsdata

# Settings Web UI

Parts of the APP is just a WEB UI loaded:

domain: https://settings.velux-active.com/#/

The page is empty per default because no access token is found. Set the cookie via and reload the page

```
document.cookie = "veluxactivecomaccess_token=" + encodeURIComponent(access_token)
```

Relevant routes:

```
path: "test"
  path: "homecontrol/:home_id"
  path: "home/:home_id"
  path: "account"
  path: "users/:home_id"
  path: "scenario/:home_id"
  path: "notifications/:home_id"
  path: "actuator/:home_id/:module_id"
  path: "calibration/:home_id/:module_id"
  path: "activity-report/:home_id"

path: "roomcontrol/:home_id/:room_id"
path: "roomcontrol/:home_id/:room_id/advanced"
path: "roomcontrol/:home_id/:room_id/schedule/:index"
path: "location/:home_id"
path: "timezone/:home_id"
path: "test/:home_id"
path: "test/:home_id/:module_id"
path: "modify_mail"
path: "modify_password"
path: "mail_preferences"
```
