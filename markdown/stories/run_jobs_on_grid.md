## Babylon should allow users to run NonMem jobs on the Grid
**Product Risk**: High

### Summary
Babylon offers an execution mode called sge: `bbi nonmem run sge <modelfile>` that operates as a simple qsub wrapper for 
Babylon itself. It creates a `grid.sh` file that literally calls babylon’s binary to execute from the context of a SGE 
worker node. This has both automated and manual testing because the automated test suite allows Babylon to create the 
script that SGE would take and run, but then manually executes it in an emulated fashion. We also manually test this 
on a metworx instance to make sure there are no errors based on the environment.

#### Tests

##### Automated Tests

Test | Test Name | Count
-----|-----------|-------
bbi_sge_test.go| <ul><li>TestBabylonCompletesSGEExecution</li><li>TestBabylonCompletesParallelSGEExecution</li></ul> |2

##### Automated Test Explanations
* TestBabylonCompletesSGEExecution : For automated testing, a fake `qsub` binary is created by the test that just
exits with a code of 0 (normally). This way, Babylon in SGE mode doesn't execute automatically, but creates the `grid.sh`
that would be executed by a remote SGE worker. The test then locates the shell script and executes it to verify functionality.
Each scenario is evaluated to see:
    * NonMem creates a `grid.sh` script that would be executed by a remote worker.
    * NonMem completes execution: The generated LST file is evaluated for a completion time
    * NonMem creates output files: The test makes sure that certain expected files (*.cpu, *.grd) exist in the output
    directory after completion of Nonmem
    * After executing `grid.sh`, Babylon creates a shell script: In each output directory should be a shell script <`modelname`>.sh which is
    executed to perform the actual work.
* TestBabylonCompletesParallelSGEExecution : Similar to the above with the exception that parallel mode is requested for execution.
This means that not only does the output file contain the grid reference, but also secondarily that the node count on the grid.
    * NonMem creates a `grid.sh` script that would be executed by a remote worker.
    * NonMem completes execution: The generated LST file is evaluated for a completion time
    * NonMem creates output files: The test makes sure that certain expected files (*.cpu, *.grd) exist in the output
    directory after completion of Nonmem
    * After executing `grid.sh`, Babylon creates a shell script: In each output directory should be a shell script <`modelname`>.sh which is
    executed to perform the actual work.

##### Manual Tests

Test | Test Name | Count
-----|-----------|-------
Run BBI In SGE Mode | qsub | 1
Verify Job On Queue | qstat | 1
Verify Nonmem Output | validate files | 1


