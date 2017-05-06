#/bin/bash

echo "TRACKMON DATABASE SETUP"
echo "This script setups the TRACKMON SERVER DATABASE for you."
sleep 2s

if (whiptail --title "Trackmon Database Setup" --yesno "Do you want to setup the Trackmon Server Database? This requires postgreSQL to be installed!" 8 78) then
  echo "Continuing setup"
else
  echo "User selected No, aborting setup."
  exit 1
fi

echo "Creating new user trackmon"
sudo adduser --disabled-login --gecos 'Trackmon' trackmon
echo "Trying to create postgres user trackmon"
sudo -u postgres psql -d template1 -c "CREATE USER trackmon CREATEDB;"
echo "Trying to create database trackmon_server_production"
sudo -u postgres psql -d template1 -c "CREATE DATABASE trackmon_server_production OWNER trackmon;"
