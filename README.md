# LotterySystem
Loterry System

A basic system which operates on the following rules:
Each ticket has a N amount of lines on it.
Each line is 3 numbers between 0 - 2 inclusive.
You can append lines to the ticket until its status has been checked.

Result for lines are :
10 pts - The sum of all numbers is 2 (0,1,1)
5 pts - If all numbers are the same (0,0,0)
1 pt -- If the 2nd and 3rd number are different from 1st (0,2,2)
0 pt - All other lines (0,1,0)

# Running Instructions
1. Open command prompt in the current directory
2. Run "./goServer/goServer.exe" command in the shell
3. Run "./goClient/goClient.exe" to get the text based UI for using the interface

# API Methods
1.(POST) -  "/ticket" - Creates a new ticket
2.(GET) -  "/ticket" -  Retrieves all tickets
3.(GET) -  "/ticket/{id}" - Retrieves a single ticket using the ID
4.(PUT) -  "/ticket/{id}/{lines}" - Adds specified amount of lines to the ticket with the ID
5.(PUT) -  "/status/{id}" - Checks the status of the ticket by ID  

# Testing
A base ticket ID : 1234567 is already provided in the code to check any functions required
