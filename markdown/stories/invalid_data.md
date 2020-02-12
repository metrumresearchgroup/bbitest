## Babylon should notify users of issues with the data referenced int he control stream
**Product Risk**: low

### Summary
If a user targets a NonMem control stream with Babylon, but the data file referenced therein cannot be located, 
Babylon should stop execution and notify the user. 

#### Tests

##### Automated Tests

Test | Test Name | Count
-----|-----------|-------
data_test.go| <ul><li>TestHasValidPathForCTL</li><li>TestHasInvalidDataPath</li></ul> |2

##### Automated Test Explanations
* TestHasValidPathForCTL: This test takes a model (*.ctl) file and changes nothing. This ensures that the data is in the
correct place. This is expected to execute completely and that no errors are thrown regarding location of the data file
* TestHasInvalidDataPath: A model file is updated programmatically to refer to a file that __does not exist__. This test
is expected to generate an error as Babylon will error indicating the file cannot be located.
