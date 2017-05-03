# Finance Server Implementation

The finance server should be the server implementation for [SirWindfields finance application](https://github.com/SirWindfield/theclassic-desktop).  

This document should define this server implementation, so that integration tests and/or a detailed development plan can be created.

# Featureset for Minimum Viable Product Release
* Creation of money accounts *(not user accounts!)*. Accounts are little, seperated finance datacollections, which should represent real world money reservoas as bank, wallet, credit card and debit card.
* Categories either private or global, which should be used to assign money to them.
  * All accounts can assign to global categories.
  * Only the parent account can assign to private categories.
* Regular adds or removes should be used for expenses or earnings that happen on a daily/weekly/monthly/yearly bases, like Spotify accounts, Netflix, income etc.

# API
The product should interface with the clients over a http based REST-API.

# Database Layout

# Internal Server Implementaion

# Additional features
