## Babylon should be able to initialize a project with minimum configs required for execution
**Product Risk**: medium

### Summary
Babylon requires a `babylon.yaml` file to be in the project directory primarily because it’s necessary for users to be 
able to select which version of nonmem to run against. Curating this from scratch is painful and error prone, so the 
`bbi init` command should be able to do this automatically.


#### Tests

##### Automated Tests

Test | Test Name | Count
-----|-----------|-------
init_test.go| TestInitialization |1 


##### Automated Test Explanations
* TestInitialization : This test copies the scenarios down to the temporary working directory and initializes each
one. After `bbi init` is executed, the test verifies that a `babylon.yaml` file has been created and that it does
contain at least one NonMem key (which is required for execution).