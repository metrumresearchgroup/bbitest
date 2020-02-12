# Top Level Expectations
1. Babylon should allow users to run NonMem Jobs:
    1. Locally
    1. On the Grid
1. Babylon should let users know if their control stream targets a data file that doesn't exist
1. Babylon should be able to initialize a project with minimum configs required for execution
1. Babylon should allow for passage of NMFE options to NonMem directly
1. Babylon should capture all configurations and render them into a file that can be stored in version control
1. Babylon should allow execution with NMQual (autolog.pl)

## Glossary

* __Output Directory__ : When a model such as `001.ctl` is targeted with Babylon, Babylon will create a new directory
called `001`, copy the model into it, and do the execution directory. This directory called `001` is referred to
as the __Output Directory__ for that model. 