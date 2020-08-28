## Babylon should allow users to parse the model output folder
**Product Risk**: medium

### Summary
Babylon creates a number of files in an output folder when a model is run. The `bbi nonmem summary` command should
be able to parse a specified subset of these files into either a human-readable summary table, or a machine-readable
`.json` structure.

#### Tests

##### Automated Tests

Test | Test Name | Count
-----|-----------|-------
bbi_summary_test.go| TestSummaryHappyPath |1
bbi_summary_test.go| TestSummaryArgs |1
bbi_summary_test.go| TestSummaryErrors |1
bbi_summary_test.go| TestSummaryHappyPathNoExtension |1 

##### Automated Test Explanations
* TestSummaryHappyPath : This test runs `bbi nonmem summary`, both with and without the `--json` flag, on a 
number of models (stored in `testdata/bbi_summary/`) that should complete successfully and compares the output 
to corresponding golden files stored in `testdata/bbi_summary/aa_golden_files/`.
* TestSummaryArgs : This test does the same thing as `TestSummaryHappyPath` except that it is testing several
models that need some specific argument passed to `bbi nonmem summary` for them to parse successfully. They
are first run _without_ the argument and the error is checked to be sure it is as expected. Then they are run 
_with_ the necessary argument and the output is compared to the golden file. 
* TestSummaryErrors : This test runs several scenarios that should cause `bbi nonmem summary` to error, for
example passing a directory path that doesn't exist, and checks for the expected error messages to be returned.
* TestSummaryHappyPathNoExtension : This test runs one of the models from `TestSummaryHappyPath` but does not 
specify any file extension when pointing to the output folder. This is to test that the program correctly
infers that an `.lst` file is needed and finds the relevant file.
