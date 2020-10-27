DRONE_BUILD_NUMBER="SethTest"

#mkdir -p /data/${DRONE_BUILD_NUMBER}/apps; chmod -R 0755 /data/${DRONE_BUILD_NUMBER}/apps
export PATH=$PATH:/data/${DRONE_BUILD_NUMBER}/apps
#go build -o /data/${DRONE_BUILD_NUMBER}/apps/bbi ../babylon/cmd/bbi/main.go

export ROOT_EXECUTION_DIR=/data/${DRONE_BUILD_NUMBER}
export BABYLON_GRID_NAME_PREFIX="drone_${DRONE_BUILD_NUMBER}"

export MPIEXEC_PATH="/usr/bin/mpiexec"
export NONMEMROOT="/opt/NONMEM"
export NMVERSION="nm74gf"
export SGE="true"
export POST_EXECUTION="true"
export NMQUAL="true"
export LOCAL="true"

#Run test suite and copy results to s3
#bbi init --dir /opt/NONMEM
#go test ./... -v --json -timeout 30m | tee test_output.json

go test ./... -v -timeout 30m
#go test -v -run TestHasValidDataPathForCTL
#go test -v -run TestSum
