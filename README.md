# Assignment 1
-   Author:         Marco Ip
-   Date:           25-2-2021
-   Worked with:    Michael-Angelo Karpowicz and Julian Kragset.

The program has all of its functions and code inside of main.go, with propper commenting, to enhace readability. 
The code consists of three enpoints: 
/exchange/v1/exchangehistory/
/exchange/v1/exchangeborder/
/exchange/v1/diag/
Which should all work as intended and described in the assignment and can be accessed both through the Heroku deployment: https://afternoon-cliffs-01801.herokuapp.com/
or through lanching it locally with localhost, where the enviroment port can be set, or defaulted to 8080.

The exchangehistory endpoint takes in two parameters: {country_name}/{begin_date-end_date}
Example: http://localhost/exchange/v1/exchangehistory/norway/2020-01-01-2020-01-10
This should allow you to see the exchange history of a given county's currency compared to euro (EUR) from the given star date to end date.
(if the currency in a country ie euro it will compare the currency to us dollars (USD))

The exchangeborder endpoint takes in one parameter: {country_name}
Example: http://localhost/exchange/v1/exchangeborder/norway
This allows the user to see all the neighbouring countries and their currency value compared to the country given (with countries using euro as an exeption)

The diag endpoint takes in no parameter and only shows some diagnostics for the apis used and the uptime.

Api's used:
https://exchangeratesapi.io/
https://restcountries.eu/

NOTE:
As mentioned, we did work on this project together, so there main problemsolving and solutions might be identical but other than that we would format the code ourselves, 
with some deviation here and there.  