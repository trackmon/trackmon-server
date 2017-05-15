#/bin/bash

if (whiptail --title "Trackmon Database Setup" --yesno "Do you want to setup the Trackmon Server Database? This requires postgreSQL to be installed!" 8 78) then
  echo "Continuing setup"
else
  echo "User selected No, aborting setup."
  exit 1
fi


echo "Creating new user trackmon"
sudo adduser --disabled-login --gecos 'Trackmon' trackmon
echo "Trying to create postgres user trackmon"
if (whiptail --title "Trackmon Database Setup" --yesno "Do you want to set a database password? Required when trackmon accesses the database remotely!" 8 78) then
  whiptail --passwordbox "Enter your database password. REMEMBER THIS!" 8 78 --title "Trackmon Database Setup" 2> password
  PASSWORD=$(cat password)
  rm password
  sudo -u postgres psql -d template1 -c "CREATE USER trackmon CREATEDB PASSWORD '$PASSWORD';"
else
  sudo -u postgres psql -d template1 -c "CREATE USER trackmon CREATEDB;"
fi
echo "Trying to create database trackmon_server_production"
sudo -u postgres psql -d template1 -c "CREATE DATABASE trackmon_server_production OWNER trackmon;"

whiptail --title "Trackmon Database Setup" --infobox "Setup finished." 8 78
