#!/bin/bash -x

function  replaceClientConfigFile() {
    # sed on MacOSX does not support -i flag with a null extension. We will use
  # 't' for our back-up's extension and delete it at the end of the function
  ARCH=$(uname -s | grep Darwin)
  if [ "$ARCH" == "Darwin" ]; then
    OPTS="-it"
  else
    OPTS="-i"
  fi

  CERT_DIR=$PWD
  CONF_DIR=$PWD/../../
  # The next steps will replace the template's contents with the
  # actual values of the private key file names for the two CAs.
  cd "$CERT_DIR"
  cd  crypto-config/peerOrganizations/org1.lzawt.com/users/Admin@org1.lzawt.com/msp/keystore
  FILE=$(ls)
  ADM1_KEY=`cat $FILE`
  cd "$CERT_DIR"
  cd  crypto-config/peerOrganizations/org1.lzawt.com/users/Admin@org1.lzawt.com/msp/signcerts
  FILE=$(ls)
  ADM1_CERT=`cat $FILE`

  cd "$CERT_DIR"
  cd  crypto-config/peerOrganizations/org2.lzawt.com/users/Admin@org2.lzawt.com/msp/keystore
  FILE=$(ls)
  ADM2_KEY=`cat $FILE`
  cd "$CERT_DIR"
  cd  crypto-config/peerOrganizations/org2.lzawt.com/users/Admin@org2.lzawt.com/msp/signcerts
  FILE=$(ls)
  ADM2_CERT=`cat $FILE`

  cd "$CONF_DIR"
  # Copy the template to the file that will be modified to add the private key
  cp config-con-template.yaml config-conn.yaml

  sed $OPTS "s/ADMIN_ORG1_KEY/${ADM1_KEY}/g" config-conn.yaml
  sed $OPTS "s/ADMIN_ORG1_CERT/${ADM1_CERT}/g" config-conn.yaml

  sed $OPTS "s/ADMIN_ORG2_KEY/${ADM2_KEY}/g" config-conn.yaml
  sed $OPTS "s/ADMIN_ORG2_CERT/${ADM2_CERT}/g" config-conn.yaml

  # If MacOSX, remove the temporary backup of the docker-compose file
  if [ "$ARCH" == "Darwin" ]; then
    rm docker-compose.yamlt
  fi
}

replaceClientConfigFile
