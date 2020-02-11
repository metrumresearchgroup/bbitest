## Babylon should allow passage of some NMFE options directly to NonMem
**Product Risk**: low

### Summary
Babylon acts as an abstraction layer for NonMem, but it may be necessary to pass some options (besides parafile) such 
as license, or compilation options directly to nonmem. Some of these are exposed by Babylon such that they are 
expressed in the final NMFE call.


#### Tests

##### Automated Tests

Test | Test Name | Count
-----|-----------|-------
bbi_local_test.go| TestNMFEOptionsEndInScript |1 
