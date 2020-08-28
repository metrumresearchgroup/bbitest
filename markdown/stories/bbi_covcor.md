## Babylon should allow users to parse the .cov and .cor files
**Product Risk**: low

### Summary
NONMEM writes information about the covariance and correlation matrices into files with the extensions `.cov`
and `.cor` respectively in the output folder. The user should be able to parse the contents of these files to
a `.json` structure.

#### Tests

##### Automated Tests

Test | Test Name | Count
-----|-----------|-------
bbi_covcor_test.go| TestCovCorHappyPath |1
bbi_covcor_test.go| TestCovCorErrors |1

##### Automated Test Explanations
* TestCovCorHappyPath : This test runs `bbi nonmem covcor` on a number of models (stored in `testdata/bbi_summary/`) 
that should complete successfully and compares the output to corresponding golden files stored in 
`testdata/bbi_summary/aa_golden_files/`.
* TestCovCorErrors : This test runs several scenarios that should cause `bbi nonmem covcor` to error, for
example passing model directory with no `.cov` or `.cor` files, and checks for the expected error messages 
to be returned.
