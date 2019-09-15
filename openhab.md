# OpenHab integration using `exec` binding

## generate `token.json`

TODO

## Define Shutters including Alexa integration

```
// South
Group:Rollershutter:AVG gShuttersSouth
  "Rollläden Süd [%d%%]"
  <rollershutter>
  (gshutters)
  ["OpenLevel", "Blinds"]
  {alexa="PercentageController.percentage" [category="SWITCH"]}

// Office
Group:Rollershutter:AVG gShuttersOffice
  "Rollläden Büro [%d%%]"
  <rollershutter>
  (Office, gshutters)
  ["OpenLevel", "Blinds"]
  {alexa="PercentageController.percentage" [category="SWITCH"]}
Rollershutter OfficeLeft_Rollershutter
  "Rollladen Büro Links [%d%%]"
  (Office, gshutters, gShuttersOffice, gShuttersSouth)
  ["OpenLevel", "Blinds"]
  {alexa="PercentageController.percentage" [category="SWITCH"]}
Rollershutter OfficeRight_Rollershutter
  "Rollladen Büro Rechts [%d%%]"
  (Office, gshutters, gShuttersOffice, gShuttersSouth)
  ["OpenLevel", "Blinds"]
  {alexa="PercentageController.percentage" [category="SWITCH"]}
```

## Update item state

## file `things/velux.things`

```
Thing exec:command:velux [command="/openhab/conf/velux-cli dump -outfile - -tokenfile /openhab/conf/token.json", interval=30, timeout=20]
```

### file `items/trigger.items`

```
String VeluxExec "velux status update" {channel="exec:command:velux:output"}
```

### file `rules/velux.rules`

```
val mapRollerShutterState = [ String veluxName, String state |
  100-Integer::parseInt(transform("JSONPATH", "$.ShutterStatus['" + veluxName + "']", state))
]

rule "velux item update"
when
   Item VeluxExec changed
then
  val jsonstate = triggeringItem.state.toString
  logInfo("velux item update", jsonstate)

  val temp = transform("JSONPATH", "$.Temperature.Bedroom", jsonstate)
  VeluxBedroomTemperature.postUpdate(Integer::parseInt(temp)/10.0)

  val co2 = transform("JSONPATH", "$.Co2.Bedroom", jsonstate)
  VeluxBedroomCo2.postUpdate(Integer::parseInt(co2))

  val humidity = transform("JSONPATH", "$.Humidity.Bedroom", jsonstate)
  VeluxBedroomHumidity.postUpdate(Integer::parseInt(humidity))

  val airQuality = transform("JSONPATH", "$.AirQuality.Bedroom", jsonstate)
  VeluxBedroomAirQuality.postUpdate(Integer::parseInt(airQuality))

  val lux = transform("JSONPATH", "$.Lux.Bedroom", jsonstate)
  VeluxBedroomLux.postUpdate(Integer::parseInt(lux))

  OfficeLeft_Rollershutter.postUpdate(mapRollerShutterState.apply("Büro links", jsonstate))
  OfficeRight_Rollershutter.postUpdate(mapRollerShutterState.apply("Büro rechts", jsonstate))

  # ...

  LivingroomLeft_Rollershutter.postUpdate(mapRollerShutterState.apply("Wohnzimmer links", jsonstate))
  LivingroomRight_Rollershutter.postUpdate(mapRollerShutterState.apply("Wohnzimmer rechts", jsonstate))

  val batterySensor = transform("JSONPATH", "$.BatteryPercent['Sensor switch 1']", jsonstate)
  VeluxSensor_BatteryLevel.postUpdate(Integer::parseInt(batterySensor))
  if (Integer::parseInt(batterySensor) < 10) {
    VeluxSensor_BatteryLow.postUpdate(ON)
  } else {
    VeluxSensor_BatteryLow.postUpdate(OFF)
  }

  val batteryDeparture = transform("JSONPATH", "$.BatteryPercent['Departure switch 1']", jsonstate)
  VeluxDeparture_BatteryLevel.postUpdate(Integer::parseInt(batteryDeparture))
  if (Integer::parseInt(batteryDeparture) < 10) {
    VeluxDeparture_BatteryLow.postUpdate(ON)
  } else {
    VeluxDeparture_BatteryLow.postUpdate(OFF)
  }
end
```



## Control Shutters

### file `rules/velux.rules`

```
val mapRollerShutterState = [ String veluxName, String state |
  100-Integer::parseInt(transform("JSONPATH", "$.ShutterStatus['" + veluxName + "']", state))
]

rule "shutters command"
when
  Member of gshutters received command
then
  logInfo("shutters", triggeringItem.name + ": " + triggeringItem.state.toString + " - " + receivedCommand)

  if (triggeringItem instanceof GroupItem) {
    return
  } else {
    var int position
    if (receivedCommand == UP) {
      position = 100
    } else if (receivedCommand == DOWN) {
      position = 0
    } else if (receivedCommand == STOP) {
      position = 100
    } else if (receivedCommand instanceof Number) {
      position = 100 - receivedCommand.intValue
    }

    var String veluxShutter
    switch (triggeringItem) {
      case OfficeLeft_Rollershutter: veluxShutter = "Büro links"
      case OfficeRight_Rollershutter: veluxShutter = "Büro rechts"

      # ...

      case LivingroomRight_Rollershutter: veluxShutter = "Wohnzimmer rechts"
      case LivingroomLeft_Rollershutter: veluxShutter = "Wohnzimmer links"
    }

    logInfo("shutters",  "moving " + veluxShutter + " to " + position)

    logInfo("shutters", "/openhab/conf/velux-cli@@moveShutters@@-tokenfile@@/openhab/conf/token.json@@-shutters@@" + veluxShutter + "@@-pos@@" + position, 1000)
    val String Answer = executeCommandLine("/openhab/conf/velux-cli@@moveShutters@@-tokenfile@@/openhab/conf/token.json@@-shutters@@" + veluxShutter + "@@-pos@@" + position, 1000)
    logInfo("shutters", "Result:" + Answer)
  }
end
```

## Usage in rules

```
var tempThreshold = 28|"℃"

rule "Close south shutters to 50% at 8:30 when too warm"
when
  Time cron "0 30 8 ? * * *"
then
  if (gWeatherForecastTodayMaxTemperature.state >= tempThreshold) {
    if (gShuttersSouth.state == 0) {
      gShuttersSouth.sendCommand(50)
    } else {
      logInfo("shutters", "some shutters already down: " + gShuttersSouth.state)
    }
  } else {
    logInfo("shutters", "keeping shutters open temperture below " + tempThreshold + ": " + gWeatherForecastTodayMaxTemperature.state)
  }
end

rule "Open south shutters at 18:00"
when
  Time cron "0 0 18 ? * * *"
then
  if (gShuttersSouth.state == 50 || gShuttersSouth.state == 70) {
    gShuttersSouth.sendCommand(0)
  } else {
    logInfo("shutters", "some shutters already down up: " + gShuttersSouth.state)
  }
end
```
