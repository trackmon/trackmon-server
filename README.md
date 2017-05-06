# Trackmon Server ![Build Status][build]

## Installing
1. Download the correct version for your OS from the release section
2. Install `postgresql` and optionally `screen` with your package manager
3. Run trackmon_db_setup.sh on the computer where you want to install the database:
  * The setup assumes that the database server is the same as the server where you want to install trackmon!
  If this is not the case, you have to configure postgreSQL to accept connections from the network.
  You can do this by changing these files:  
  `pg_hba.conf`:  
  ```
  host all all 0.0.0.0/0 md5
  ```
  `postgresql.conf`  
  ```
  listen_addresses='*'
  ```
  After that you simply run `service postgresql restart` and you should be good to go. Do not forget to change trackmon servers configfile to your database location. **Do not do this if the database is on the same machine as trackmon!**  

4. Generate trackmon servers configfile by starting trackmon with the `-createconf` flag. There should be a `trackmonserv.conf` file in the directory. If the database is not on the same device, you have to edit the configfile.
5. Start trackmon server. It is recommended that you do this inside a screen (or similar) so that the server continues to run when you close the session.

[build]: https://api.travis-ci.org/trackmon/trackmon-server.svg?branch=master
