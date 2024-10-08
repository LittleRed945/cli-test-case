Feature: Data Product ruleset update

Scenario:
Given NATS has been opened
Given Dispatcher has been opened
#Scenario
	Scenario: Success scenario for updating data product ruleset
	Given Create data product "'drink'" and enabled is "'[true]'"
    Given Create "'drink'" ruleset "'drinkCreated'" and enabled is "'<GivenRSEnabled>'"
	When Update data product "'<ProductName>'" ruleset "'<Ruleset>'" using parameters "'<Method>'" "'<Event>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'"
	Then Check updating ruleset success
	And Use NATS jetstream to query the "'drink'" "'drinkCreated'" update successfully and parameters are "'<Method>'" "'<Event>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'"
	Examples:
    |  ID  | ProductName | Ruleset       | Event         | Method    | 		Schema          | 		Handler_script	   | Pk       | Desc          | Enabled  | GivenRSEnabled |
	| M(1) | drink       | drinkCreated  | drinkCreated  | create    | ./assets/schema.json |  ./assets/handler.js     | id       |  description  | [true]   |   [ignore]   |
	| M(2) | drink       | drinkCreated  | drinkCreated  | [ignore]  | 		[ignore]        | 		  [ignore]         | [ignore] | [ignore]      | [ignore] |   [ignore]   |
    #Error occurred while using alone update method: Invalid method
    | M(3) | drink       | drinkCreated  | [ignore]      | create    | 		[ignore]        | 		  [ignore]         | [ignore] | [ignore]      | [ignore] |   [ignore]   |
    | M(4) | drink       | drinkCreated  | [ignore]      | [ignore]  |./assets/schema.json  | 		  [ignore]         | [ignore] | [ignore]      | [ignore] |   [ignore]   |
    | M(5) | drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 	./assets/handler.js    | [ignore] | [ignore]      | [ignore] |   [ignore]   |
    | M(6) | drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 		  [ignore]         | id       | [ignore]      | [ignore] |   [ignore]   |
	| M(7) | drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 		  [ignore]         | "id, num"| [ignore]      | [ignore] |   [ignore]   |
    | M(8) | drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 		  [ignore]         | ""       | [ignore]      | [ignore] |   [ignore]   |
    | M(9) | drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 		  [ignore]         | [ignore] | description   | [ignore] |   [ignore]   |
    | M(10)| drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 		  [ignore]         | [ignore] | ""            | [ignore] |   [ignore]   |
    | M(11)| drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 		  [ignore]         | [ignore] | " "           | [ignore] |   [ignore]   |
    | M(12)| drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 		  [ignore]         | [ignore] | [ignore]      | [true]   |   [ignore]   |
	| M(13)| drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 		  [ignore]         | [ignore] | [ignore]      | [false]  |   [true]     |
	| M(14)| drink       | drinkCreated  | [ignore]      | [ignore]  | 		[ignore]        | 		  [ignore]         | [ignore] | [ignore]      | [ignore] |   [ignore]   |

#Scenario
	Scenario Outline: Fail scenario for updating data product ruleset
	Given Create data product "'drink'" and enabled is "'[true]'"
    Given Create "'drink'" ruleset "'drinkCreated'" and enabled is "'[ignore]'"
	Given Store NATS copy of existing data product "'drink'" ruleset "'drinkCreated'"
	When Update data product "'<ProductName>'" ruleset "'<Ruleset>'" using parameters "'<Method>'" "'<Event>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'"
	Then CLI returns exit code 1
	# And The error message should be "'<Error_message>'"
	And Use NATS jetstream to query the "'drink'" "'drinkCreated'" without changing parameters
	Examples:
	|  ID   | ProductName | Ruleset       | Event         | Method    | 		Schema         	 	 | 		Handler_script	   | Pk       | Desc          | Enabled  | Error_message |
	| E1(1) | [null]	  | [null]		  | [ignore]	  | [ignore]  | 		[ignore]        	 | 		  [ignore]         | [ignore] | [ignore]      | [ignore] | 		         |
	| E1(2) | drink       | [null]		  | drinkCreated  | create    |./assets/schema.json  		 | 	./assets/handler.js    | id       | description   | [true]   | 		         | 
	| E1(3) | [null]      | drinkCreated  | drinkCreated  | create    |./assets/schema.json  		 | 	./assets/handler.js    | id       | description   | [true]   | 		         |
	| E1(4) | drink       | not_exist	  | drinkCreated  | create    |./assets/schema.json  		 | 	./assets/handler.js    | id       | description   | [true]   | 		         |
	| E1(5) | not_exist   | drinkCreated  | drinkCreated  | create    |./assets/schema.json  		 | 	./assets/handler.js    | id       | description   | [true]   | 		         |
	| E1(6) | drink       | drinkCreated  | drinkCreated  | create    |		"not_exist.json" 	  	 | 	./assets/handler.js    | id       | description   | [true]   | 		         |
	| E1(7) | drink       | drinkCreated  | drinkCreated  | create    |		"" 						 | 	./assets/handler.js    | id       | description   | [true]   | 		         |	
	| E1(8) | drink       | drinkCreated  | drinkCreated  | create    |./assets/fail_schema.json	 | 	./assets/handler.js    | id       | description   | [true]   | 		         |
	| E1(9) | drink       | drinkCreated  | drinkCreated  | create    |./assets/schema.json  		 | 	not_exist.js		   | id       | description   | [true]   | 		         |
	| E1(10)| drink       | drinkCreated  | drinkCreated  | create    |./assets/schema.json  		 | 	""					   | id       | description   | [true]   | 		         |
	| E1(11)| drink  	  | drinkCreated  | ""			  | create    |./assets/schema.json  		 | 	./assets/handler.js    | id       | description   | [true]   | 		         |
	| E1(12)| drink  	  | drinkCreated  | drinkCreated  | "" 	      |./assets/schema.json  		 | 	./assets/handler.js    | id       | description   | [true]   | 		         |
	| E1(13)| drink  	  | drinkCreated  | " "			  | create    |./assets/schema.json  		 | 	./assets/handler.js    | id       | description   | [true]   | 		         |
	| E1(14)| drink  	  | drinkCreated  | drinkCreated  | " "       |./assets/schema.json  		 | 	./assets/handler.js    | id       | description   | [true]   | 		         |