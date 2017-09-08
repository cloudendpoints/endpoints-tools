#!/usr/bin/env python
# Copyright 2017 Google Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Python script generates a Google ID token based on the input payload"""

import base64
import httplib
import json
import argparse
import time
import urllib

import oauth2client.crypt
import googleapiclient.discovery
from oauth2client.service_account import ServiceAccountCredentials
from oauth2client.client import GoogleCredentials

def generate_jwt(args):
    """Generates a signed JSON Web Token using a service account. Based on https://cloud.google.com/endpoints/docs/service-to-service-auth"""
    # Make sure the service account has "Service Account Token Creator" permissions in Google IAM
    credentials = ServiceAccountCredentials.from_json_keyfile_name(
      args.service_account_file).create_scoped(['https://www.googleapis.com/auth/cloud-platform'])

    service = googleapiclient.discovery.build(
        serviceName='iam', version='v1', credentials=credentials)

    now = int(time.time())
    header_json = json.dumps({
        "typ": "JWT",
        "alg": "RS256"})

    payload_json = json.dumps({
        'iat': now,
        "exp": now + 3600,
        'iss': args.issuer if args.issuer else credentials.service_account_email,
        "target_audience": 'https://' + args.aud,
        "aud": "https://www.googleapis.com/oauth2/v4/token"
    })

    header_and_payload = '{}.{}'.format(
        base64.urlsafe_b64encode(header_json),
        base64.urlsafe_b64encode(payload_json))
    slist = service.projects().serviceAccounts().signBlob(
        name="projects/-/serviceAccounts/" + credentials.service_account_email,
        body={'bytesToSign': base64.b64encode(header_and_payload)})
    res = slist.execute()
    signature = base64.urlsafe_b64encode(
        base64.decodestring(res['signature']))
    signed_jwt = '{}.{}'.format(header_and_payload, signature)

    return signed_jwt

def main(args):
    """Request a Google ID token using a JWT."""
    params = urllib.urlencode({
        'grant_type': 'urn:ietf:params:oauth:grant-type:jwt-bearer',
        'assertion': generate_jwt(args)})
    headers = {"Content-Type": "application/x-www-form-urlencoded"}
    conn = httplib.HTTPSConnection("www.googleapis.com")
    conn.request("POST", "/oauth2/v4/token", params, headers)
    res = json.loads(conn.getresponse().read())
    conn.close()
    return res['id_token']

if __name__ == '__main__':
  parser = argparse.ArgumentParser(
      description=__doc__,
      formatter_class=argparse.RawDescriptionHelpFormatter)
  # positional arguments
  parser.add_argument(
      "aud",
      help="Audience . This must match 'audience' in the security configuration"
      " in the swagger spec. It can be any string")
  parser.add_argument(
        'service_account_file',
        help='The path to your service account json file.')

  #optional arguments
  parser.add_argument("-iss", "--issuer", help="Issuer claim. This will also be used for sub claim")
  print main(parser.parse_args())
