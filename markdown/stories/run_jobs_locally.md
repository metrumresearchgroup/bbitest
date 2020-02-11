## Babylon should allow users to run NonMem jobs locally
**Product Risk**: High

### Summary
If NonMem is installed on the system, users should be able to use it to execute jobs without needing the grid. 

#### Tests

##### Automated Tests

Test | Test Name | Count
-----|-----------|-------
bbi_local_test.go| <ul><li>TestBabylonCompletesLocalExecution</li><li>TestBabylonParallelExecution</li></ul> |2
