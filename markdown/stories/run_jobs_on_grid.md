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


##### Manual Tests

Test | Test Name | Count
-----|-----------|-------
Run BBI In SGE Mode | qsub | 1
Verify Job On Queue | qstat | 1
Verify Nonmem Output | validate files | 1


