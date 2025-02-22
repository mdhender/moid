#!/bin/bash
##############################################################################
# This script builds the executables and pushes them to the server.
##############################################################################
# Confirm that we're running from the root of the repository.
[ -f go.mod -a -f bin/push_build.sh ] || {
  echo error: must run from the root of the repository
  exit 2
}
##############################################################################
# user must specify dev, test, or prod on the command line
buildEnvironment=
case "${1}" in
  dev|development)
    echo " info: building for development"
    buildEnvironment=development ;;
  test)
    echo " info: building for test"
    buildEnvironment=test ;;
  prd|prod|production)
    echo " info: building for production"
    buildEnvironment=production ;;
  *)
    echo "error: invalid environment: '${1}'"
    exit 2 ;;
esac
##############################################################################
# create the build directory if it doesn't exist
mkdir -p build
##############################################################################
# clean out the build directory
rm -rf build/*
mkdir -p build/{bin,data,logs,services} build/ui/{assets,views}
##############################################################################
# build the local executable to get the version number
mkdir -p build/bin
LOCAL_EXE=build/bin/epimethean
echo " info: building local executable..."
go build -o "${LOCAL_EXE}" || {
  echo "error: unable to build local executable"
  exit 2
}
VERSION=$( "${LOCAL_EXE}" version )
if [ -z "${VERSION}" ]; then
  echo "error: '${LOCAL_EXE} version' seems to have failed"
  exit 2
fi

##############################################################################
# build the linux executable that we'll push to the server
echo " info: building executables for version '${VERSION}'"
echo " info: building linux executable..."
LINUX_EXE=build/bin/epimethean
GOOS=linux GOARCH=amd64 go build -o "${LINUX_EXE}" || exit 2

##############################################################################
# copy the configuration files that we'll push to the server
echo " info: copying configuration files..."
cp -p testdata/dev.epimethean/.env* build/services/ || exit 2

##############################################################################
# copy the ui files that we'll push to the server
echo " info: copying ui files..."
cp -pr ui/assets build/ui || exit 2
cp -pr ui/views  build/ui || exit 2

##############################################################################
# create a compressed tarball of the assets and components.
# ensure that it doesn't contain Mac junk attributes.
echo " info: creating tarball..."
cd build || exit 2
TARBALL=epimethean-${VERSION}.tgz
tar -cz --no-xattrs --no-mac-metadata -f ${TARBALL} bin data services ui/assets ui/views || exit 2
cd - || exit 2

##############################################################################
# push the tarball to our development server
echo " info: pushing tarball to development web server..."
scp build/${TARBALL} epimethean@epimethean:/var/www/dev.epimethean/build/ || {
  echo "error: failed to copy the tarball to the development web server"
  exit 2
}

###############################################################################
## clean up the build directory
#echo " info: removing build files..."
#echo rm -rf build/*

##############################################################################
#
echo " info: push to development web server succeeded"
exit 0
