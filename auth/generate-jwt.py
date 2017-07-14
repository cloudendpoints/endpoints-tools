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

"""Python script generates a signed JWT token based on the input payload"""

import argparse
import time

import oauth2client.crypt
from oauth2client.service_account import ServiceAccountCredentials

def main(args):
  """Generates a signed JSON Web Token using a Google API Service Account."""
  credentials = ServiceAccountCredentials.from_json_keyfile_name(
      args.service_account_file)

  now = int(time.time())

  payload = {
        "exp": now + credentials.MAX_TOKEN_LIFETIME_SECS,
        "iat": now,
        "aud": args.aud,
    }

  if args.email:
    payload["email"] = args.email
  if args.groupId:
    payload["groupId"] = args.groupId

  if args.issuer:
    payload["iss"] = args.issuer
    payload["sub"] = args.issuer
  else:
    payload["iss"] = credentials.service_account_email
    payload["sub"] = credentials.service_account_email

  signed_jwt = oauth2client.crypt.make_signed_jwt(
        credentials._signer, payload, key_id=credentials._private_key_id)

  return signed_jwt

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
  parser.add_argument("-e", "--email", help="Email claim in JWT")
  parser.add_argument("-g", "--groupId", help="GroupId claim in JWT")
  parser.add_argument("-iss", "--issuer", help="Issuer claim. This will also be used for sub claim")
  print main(parser.parse_args())
