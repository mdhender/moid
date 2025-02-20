#!/bin/bash

# confirm that we're running from the root of the repository
[ -d bin ] || {
  echo error: must run from the root of the repository
  exit 2
}
[ -d build ] || {
  echo error: must run from the root of the repository
  exit 2
}

# clean out the build directory
rm -rf build && mkdir build || exit 2

# build the executable
echo " info: building executable..."
GOOS=linux GOARCH=amd64 go build -o build/epimethean || exit 2

# create a compressed tarball of the assets and components.
# ensure that it doesn't contain Mac junk attributes.
echo " info: creating tarball..."
tar -cz --no-xattrs --no-mac-metadata -f build/assets.tgz ui/assets || exit 2
tar -cz --no-xattrs --no-mac-metadata -f build/views.tgz  ui/views  || exit 2

# push the tarbals file to our production server
echo " info: pushing tarball..."
scp build/epimethean epimethean@epimethean:/home/epimethean/dev/build || exit 2
scp build/assets.tgz epimethean@epimethean:/home/epimethean/dev/epimethean.tgz || exit 2
scp epimethean.tgz epimethean@epimethean:/var/www/dev.epimethean/epimethean.tgz || exit 2

# execute the installation script
echo " info: executing the installation script..."
ssh epimethean@epimethean /home/epimethean/dev/bin/install.sh || {
  echo "error: installation script failed"
  exit 2
}

# next
echo " info: if this succeeded, you should restart the services"
echo "       ssh epimethean systemctl restart epimethean.service"
echo "       ssh epimethean systemctl status  epimethean.service"
echo "       ssh epimethean journalctl -f -u  epimethean.service"

exit 0
