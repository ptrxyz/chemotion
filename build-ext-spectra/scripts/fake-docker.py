#!/usr/bin/env python3

# example command would look like this:
# docker exec -d msconvert_docker wine msconvert /data/{}/{} -o /data/{}
#   --32 --zlib --filter "peakPicking true 1-" --filter
#   "zeroSamples removeExtra" --ignoreUnknownInstrumentError

import sys
import requests
import os


def toBool(s):
    """
    'Safely' converts a string to bool.
    """
    try:
        if type(s) == bool and s is True:
            return True
        if s.strip().lower() == "true":
            return True
    except Exception:
        pass
    return False


host = os.getenv("MSC_HOST", "msconvert")
port = os.getenv("MSC_PORT", "8088")
validate = toBool(os.getenv("MSC_VALIDATE", "False"))

url = f"http://{host}:{port}/run"

args = list(sys.argv[1:])
print(args)

try:
    if validate:
        assert args[0] == "exec"
        assert args[1] == "-d"
        assert args[2] == "msconvert_docker"
        assert args[3] == "wine"
    idx = args.index("wine")
    args = args[idx:]
except (ValueError, AssertionError):
    sys.exit(1)

print("Using args: ", args)
r = requests.post(url, json=args)
print(r.text, r.status_code)
