#/bin/bash

echo "TRACKMON INSTALLER"
echo "This script installs TRACKMON SERVER for you."
sleep 2s

if (whiptail --title "Trackmon Installer" --yesno "Do you want to install Trackmon Server? By installing you agree to the license.\n\
It can be found at https://github.com/trackmon/trackmon-server" 8 78) then
  echo "Continuing installation"
else
  echo "User selected No, aborting installation."
  exit 1
fi
echo -n "Getting TRACKMON MANAGER program... "
wget -q https://raw.githubusercontent.com/trackmon/trackmon-server/master/manager/trackmon_manager.py
echo "done."
if (whiptail --title "Trackmon Installer" --yes-button "Fresh Install" --no-button "Update" --yesno "Is this a fresh install or an update" 8 78) then
  echo "Making a fresh install"
  whiptail --title "Trackmon Installer" --menu "What is the Linux distro of your choice?" 25 78 16 "install" "Full Install" "installapi" "Application Server Only" "installdb" "Database Only" "installfrontend" "Website Only" 2>installchoice
  INSTALLTYPE=$(cat installchoice)
  rm installchoice
  exitstatus=$?
  if [ $exitstatus = 0 ]; then
    echo "Launching MANAGER with $INSTALLTYPE flag"
    python3 trackmon_manager.py -$INSTALLTYPE
  else
    echo "Aborting the installation."
  fi
else
  echo "Making an update"
fi
