# Finance Server Implementation - WIP

## Feature set for Minimum Viable Product Release
* Creation of user accounts
* Creation of money accounts *(not user accounts in this context!)*. Accounts are little, separated finance data collections, which should represent real world money reservoirs as bank, wallet, credit card and debit card.
* Categories either private or global, which should be used to assign money to them.
  * All accounts can assign to global categories.
  * Only the parent account can assign to private categories.
* Regular adds or removes should be used for expenses or earnings that happen on a daily/weekly/monthly/yearly bases, like Spotify accounts, Netflix, income etc.

## API
API documentation is [here](https://github.com/trackmon/trackmon-api)

## Internal Server Implementation
The server is written in Go and tries to stick to the standard library as much as possible, with a few exceptions, such as gorilla/mux.
PostgreSQL is used as the database server, which is configured by the installer with the database layout below.

## Database Layout
Not yet determined.
