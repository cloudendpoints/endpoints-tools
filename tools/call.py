#!/usr/bin/python -u

import argparse
import urllib2
import json
import time
import uuid

"""
This script will repeatly call Google service control service on every second
to detect big latency.  Its usage

./call.py --api_key=$API_KEY  --sevice=$SERVICE

SERVICE: an endpoint service for service control Check and Report.
API_KEY: an api key of the project.

"""

def fetch_access_token():
  url = "http://169.254.169.254/computeMetadata/v1/instance/service-accounts/default/token"
  headers = {"Metadata-Flavor": "Google"}
  request = urllib2.Request(url, None, headers)
  response = urllib2.urlopen(request)
  contents = response.read()
  return json.loads(contents)["access_token"]


def call_report(access_token, operation_id, args):
  headers = {"Authorization": "Bearer {}".format(access_token),
             "X-Cloud-Trace-Context": "{};o=1".format(operation_id),
             "Content-Type": "application/json"}
  url = "https://servicecontrol.googleapis.com/v1/services/{}:report".format(args.service_name)
  data_obj = {"service_name": args.service_name,
          "operations": [{
              "operation_id": operation_id,
              "operation_name": "/echo",
              "consumer_id": "api_key:{}".format(args.api_key),
              "start_time": {
                "seconds": int(time.time())
              },
              "end_time": {
                "seconds": int(time.time())
              }
           }]
         }
  data = json.dumps(data_obj)
  t0 = time.time()
  try:
    request = urllib2.Request(url, data, headers)
    response = urllib2.urlopen(request)
    trace_id = response.info().getheader("X-GOOG-TRACE-ID")
#    print "response: {}".format(response.info())
  except urllib2.HTTPError as e:
    print "{} Check failed code: {},  error {}".format(time.ctime(), e.code, e.reason)
    return
  latency = time.time() - t0
  if trace_id and (latency >= 1.0):
    print "{}: report big latency {}, trace_id: {} operation_id: {}".format(time.ctime(), latency, trace_id, operation_id)


def call_check(access_token, operation_id, args):
  headers = {"Authorization": "Bearer {}".format(access_token),
             "X-Cloud-Trace-Context": "{};o=1".format(operation_id),
             "Content-Type": "application/json"}
  url = "https://servicecontrol.googleapis.com/v1/services/{}:check".format(args.service_name)
  data_obj = {"service_name": args.service_name,
          "operation": {
              "operation_id": operation_id,
              "operation_name": "/echo",
              "consumer_id": "api_key:{}".format(args.api_key),
              "start_time": {
                "seconds": int(time.time())
              }
           }
         }
  data = json.dumps(data_obj)
  t0 = time.time()
  try:
    request = urllib2.Request(url, data, headers)
    response = urllib2.urlopen(request)
    trace_id = response.info().getheader("X-GOOG-TRACE-ID")
#    print "response: {}".format(response.info())
  except urllib2.HTTPError as e:
    print "{} Check failed code: {},  error {}".format(time.ctime(), e.code, e.reason)
    return
  latency = time.time() - t0
  if trace_id and (latency >= 1.0):
    print "{}: check big latency {}, trace_id: {} operation_id: {}".format(time.ctime(), latency, trace_id, operation_id)

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--api_key', help='Project api_key to access service.')
    parser.add_argument('--service_name', help='service name.')
    args = parser.parse_args()

    print "Start repeatedly call google service control"
    cnt = 0
    while True:
      token = fetch_access_token()
      operation_id = uuid.uuid4().hex
      call_report(token, operation_id, args)
      call_check(token, operation_id, args)
      time.sleep(1)
      cnt += 1
      if (cnt % 10 == 0):
        print "{}: Number of calls: {}".format(time.ctime(), cnt)
