DRONE_BUILD_NUMBER="SethTest"

export ROOT_EXECUTION_DIR=/data/${DRONE_BUILD_NUMBER}
export BABYLON_GRID_NAME_PREFIX="drone_${DRONE_BUILD_NUMBER}"

export MPIEXEC_PATH="/usr/bin/mpiexec"
export NONMEMROOT="/opt/NONMEM"
export NMVERSION="nm74gf"
export SGE="true"
export POST_EXECUTION="true"
export NMQUAL="true"
export LOCAL="true"

#go test -v -run TestHasValidDataPathForCTL
go test -v -run TestBabylonCompletesLocalExecution
