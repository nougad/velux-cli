# velux-cli
Go Client for Velux Active KIX 300

## Install

```
go get github.com/nougad/velux-cli
# will fail with: package github.com/nougad/velux-cli/client: cannot find package "github.com/nougad/velux-cli/client" in any of:
# go then into `src/github.com/nougad/velux-cli` to generate the swagger client
$ cd src/github.com/nougad/velux-cli
# make changes to the `config.go` file and add velux client id and client secret
# and run make
$ make
```

## Auth

Login to Velux using username and password

```
CLIENT_ID=".."
CLIENT_SECRET=".."
USERNAME="email"
PASSWORD=""
curl -v -d "grant_type=password&client_id=${CLIENT_ID}&client_secret=${CLIENT_SECRET}&username=${USERNAME}&password=${PASSWORD}&user_prefix=velux" https://app.velux-active.com/oauth2/token
```

Then create a file named `token.json`

```
{
 "token": {
  "access_token": "5...d|5...0",
  "refresh_token": "5...d|c...4",
  "scope": [
   "all_scopes"
  ],
  "expires_in": 10800
 },
 "refreshed": "2019-08-11T21:42:51.597296912+02:00"
}
```

## Velux protocol

Documented in [velux-protocol.md](velux-protocol.md).

## OpenHab integration using exec plugin

Can be found in [openhab.md](openhab.md) file.

The current OpenHab netatmo binding *might* be able to use the reading from
`GetMeasure` but it does not allow to override the API endpoint so I could not
try.

There is currently no OpenHab binding since it would require disclosing client
secret and client id reverse engineered by Android App. More details at
https://community.openhab.org/t/connecting-velux-active-kix-300/75696

## Swagger definition

Can be found in [swagger.yaml](swagger.yaml). Only relevant methods added.

Due to missing union types in of OpenAPI 2.0 and missing code generators in
OpenHab 3.0 the `setState` only takes percentage right now.

## Open issues

* `stopall` is missing
* Only rollershutters supported right now
