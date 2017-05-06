import sys
import os
from subprocess import call
import urllib.request
import json
from pprint import pprint

# User needs to install postgres first
trackmon_server_api_info = "https://api.github.com/repos/atom/atom/releases/latest"

current_os = ""
current_arch = ""
if sys.platform.startswith('linux'):
    current_os = "linux"
elif sys.platform.startswith('win32'):
    current_os = "windows"
elif sys.platform.startswith('darwin'):
    current_os = "darwin"
else:
    print("Your system is not supported by this installer.")
    sys.exit(0)

def is_os_64bit():
    return platform.machine().endswith('64')

if is_os_64bit == True:
    current_arch = 64

class color:
	HEADER = '\033[95m'
	OKBLUE = '\033[94m'
	OKGREEN = '\033[92m'
	WARNING = '\033[93m'
	FAIL = '\033[91m'
	ENDC = '\033[0m'
	BOLD = '\033[1m'
	UNDERLINE = '\033[4m'

def split(string, splitters): #MAY RESOLVE ALL PROBLEMS WITH CSV
	final = [string]
	for x in splitters:
		for i,s in enumerate(final):
			if x in s and x != s:
				left, right = s.split(x, 1)
				final[i] = left
				final.insert(i + 1, x)
				final.insert(i + 2, right)
	return final

def download(url, path):
    with urllib.request.urlopen(url) as response, open(path, 'wb') as output:
        shutil.copyfileobj(response, output)

def get_dl_from_gh_api(url):
    response = urllib.request.urlopen(url)
    data = response.read()
    jsonresp = json.loads(data.decode('utf-8'))
    #pprint(jsonresp)
    #print(jsonresp["assets"])
    for asset in jsonresp["assets"]:
        assetname = str(asset["name"])
        splitted_assetname = split(assetname, "_")
        sys_and_arch = split(splitted_assetname[2], ["-", ".") # Dot for windows versions
        if sys_and_arch[0] == current_os and sys_and_arch[2] == current_arch:
            print("Downloading server...", end='')
            download(str(asset["browser_download_url"]), assetname)
            print("done.")
            return
    print("Didn't find any fitting version, you might have to download it manually")

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
