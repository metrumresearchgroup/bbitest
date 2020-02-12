## Babylon should allow users to run NonMem jobs locally
**Product Risk**: High

### Summary
If NonMem is installed on the system, users should be able to use it to execute jobs without needing the grid. 

#### Tests

##### Automated Tests

Test | Test Name | Count
-----|-----------|-------
bbi_local_test.go| <ul><li>TestBabylonCompletesLocalExecution</li><li>TestBabylonParallelExecution</li></ul> |2

##### Automated Test Explanations
* TestBabylonCompletesLocalExecution : This is an end-to-end test of Babylon in local execution mode. Here, Babylon uses
a local NonMem instance to run all models in all three test scenarios. Each scenario is evaluated to make sure:
    * NonMem completes execution: The generated LST file is evaluated for a completion time
    * NonMem creates output files: The test makes sure that certain expected files (*.cpu, *.grd) exist in the output
    directory after completion of Nonmem
    * Babylon creates a shell script: In each output directory should be a shell script <`modelname`>.sh which is
    executed to perform the actual work.
* TestBabylonParallelExecution : This is an end-to-end test of Babylon in local execution mode across multiple
compute nodes in parallel. Here, Babylon uses a local NonMem instance to run all models in all three test scenarios. 
Each scenario is evaluated to make sure:
    * NonMem completes execution: The generated LST file is evaluated for a completion time
    * NonMem creates output files: The test makes sure that certain expected files (*.cpu, *.grd) exist in the output
    directory after completion of Nonmem
    * Babylon creates a shell script: In each output directory should be a shell script <`modelname`>.sh which is
    executed to perform the actual work.
    * The Nonmem Output (lst) file contains a reference to the `parafile` indicating that the work was done in parallel