/*
Package postman implements a library that configures tester.Tester from Postman collections.

Instructions:
	1. Create a new collection in Postman
	2. Create your requests in Postman. Request name will be imported as tester.TestCategory
	3. Execute your requests and save as Example. These examples will become tester.TestScenario
	4. Export your collection from Postman (right click/Export) as Collection v2.1

This library will replace variables in the request/example with variables defined on the collection.
*/
package postman
