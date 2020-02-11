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
