import sys
import os
from subprocess import call
import urllib.request
import json
#from pprint import pprint

# User needs to install postgres first
trackmon_server_api_info = "https://api.github.com/repos/paulkramme/roverpi/releases/latest"

def download(url, path):
    with urllib.request.urlopen(url) as response, open(path, 'wb') as output:
        shutil.copyfileobj(response, output)

def get_dl_from_gh_api(url):
    response = urllib.request.urlopen(url)
    data = response.read()
    jsonresp = json.loads(data.decode('utf-8'))
    #pprint(json)
    for asset in jsonresp["assets"]:
        print(str(asset["name"])) # BUG: Nothing prints here...
    print("Done.")

def main():
    if "-install" in sys.argv:
        print("Installing everything")
        # TODO: Verify that postgres exist
        # TODO: Download trackmon server
        get_dl_from_gh_api(trackmon_server_api_info)
    elif "-installapi" in sys.argv:
        print("Installing API backend only")
        # TODO: Download trackmon server
    elif "-installdb" in sys.argv:
        print("Installing database only")
        # TODO: Verify that postgres exist
    elif "-installfrontend" in sys.argv:
        print("Installing frontend only")
        # TODO: Later...

    elif "-update" in sys.argv:
        print("Updating components")

if __name__ == "__main__":
    main()
