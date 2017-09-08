# Generate custom JWT for service-to-service authentication
There are two forms of JSON Web Token (JWT) signed by Google: A Google ID token or a custom JWT that 
is signed only by the service account of the caller. The following two endpoint security definitions
can be used:
```
securityDefinitions:
  [ISSUER_NAME]:
    authorizationUrl: ""
    flow: "implicit"
    type: "oauth2"
    x-google-issuer: "[YOUR-SERVICE-ACCOUNT-EMAIL]"
    x-google-jwks_uri: "https://www.googleapis.com/robot/v1/metadata/x509/[YOUR-SERVICE-ACCOUNT-EMAIL]"
```
```
securityDefinitions:
  google_id_token:
    authorizationUrl: ""
    flow: "implicit"
    type: "oauth2"
    x-google-issuer: "https://accounts.google.com"
```
This is based on the following endpoint documentation (https://cloud.google.com/endpoints/docs/service-to-service-auth).


## Setup
Before using the script please run the following command to install python dependences:
```
$pip install google-cloud
```
You will also need to create and save the service account json file.


## Usage
To generate a custom JWT that is signed only by the service account use the following:
```
$ python generate-jwt.py -h
usage: generate-jwt.py [-h] [-e EMAIL] [-g GROUPID] [-iss ISSUER]
                       aud service_account_file

Python script generates a signed JWT token based on the input payload

positional arguments:
  aud                   Audience . This must match 'audience' in the security
                        configuration in the swagger spec. It can be any
                        string
  service_account_file  The path to your service account json file.

optional arguments:
  -h, --help            show this help message and exit
  -e EMAIL, --email EMAIL
                        Email claim in JWT
  -g GROUPID, --groupId GROUPID
                        GroupId claim in JWT
  -iss ISSUER, --issuer ISSUER
                        Issuer claim. This will also be used for sub claim
```
To generate a Google ID token JWT use the following:
```
$ python generate-google-id-jwt.py -h
usage: generate-google-id-jwt.py [-h] [-iss ISSUER]
                       aud service_account_file

Python script generates a signed Google ID JWT token based on the input payload

positional arguments:
  aud                   Audience . This must match 'audience' in the security
                        configuration in the swagger spec. It can be any
                        string
  service_account_file  The path to your service account json file.

optional arguments:
  -h, --help            show this help message and exit
  -iss ISSUER, --issuer ISSUER
                        Issuer claim. This will also be used for sub claim
```


## Examples
1. Generate JWT token without any custom claims
```
$ python generate-jwt.py <YOUR-AUDIENCE> /path/to/service_account.json
```
2. Generate a Google ID JWT token without any custom claims
```
$ python generate-google-id-jwt.py <YOUR-AUDIENCE> /path/to/service_account.json
```
3. Generate JWT token with email claim
```
$ python generate-jwt.py -e alice@yahoo.com <YOUR-AUDIENCE> /path/to/service_account.json
```
4. Generate JWT token with groupId claim
```
$ python generate-jwt.py -g acme <YOUR-AUDIENCE> /path/to/service_account.json
```
5. Generate JWT token with both email and groupId claim
```
$ python generate-jwt.py -e alice@yahoo.com -g acme <YOUR-AUDIENCE> /path/to/service_account.json
```
