## Babylon should capture all configurations and render them into a file that can be stored in version control
**Product Risk**: high

### Summary
Due to the complexity of flags and configuration files available to Babylon, it is necessary to provide a record of 
what settings were used to execute a job for the sake of reproducibility. As such, at the end of every execution, a 
`bbi_config.json` file exists which contains the merged configurations between any flags provided, configuration files 
and default values to indicate exactly how the model was executed. 

#### Tests

##### Automated Tests

Test | Test Name | Count
-----|-----------|-------
config_test.go| TestBBIConfigJSONCreated |1 

##### Automated Test Explanations
* TestBBIConfigJSONCreated : Runs through a single scenario and verifies that Nonmem completes. After completion of
work, the test verifies that a `bbi_config.json` file exists within the model output directory. This file is then
re-processed to verify that it contains the `nmVersion` specified during the test execution.