#!/bin/sh

export GIT_URL=$1
export REF_NAME=$2

echo "GIT_URL: $GIT_URL, REF_NAME: $REF_NAME"

if [ -z "$GIT_URL" ]; then
    echo "GIT_URL not defined"
    exit 1
fi

if [ -z "$REF_NAME" ]; then
    echo "REF_NAME not defined"
    exit 1
fi

echo "Building $REF_NAME ..."

# Check that we've a valid working directory.
if [ ! -d "$APP_WORKING_DIR" ]; then
    echo "Error, APP_WORKING_DIR not found"
    exit 1
fi

# Check that the deploy key is valid.
export DEPLOY_KEY=/root/.ssh/id_rsa
if [ ! -f "$DEPLOY_KEY" ]; then
    echo "Error, DEPLOY_KEY not found"
    exit 1
fi

# Remove old repository if exist
rm -rf $APP_WORKING_DIR/$REF_NAME

# Clone repository
echo "Cloning $GIT_URL into ${APP_WORKING_DIR}/${REF_NAME} ..."
ssh-agent bash -c 'ssh-add ${DEPLOY_KEY}; \
    git clone ${GIT_URL} ${APP_WORKING_DIR}/${REF_NAME}; \
    cd ${APP_WORKING_DIR}/${REF_NAME}; \
    git remote add zone52 ssh://5510399367c989f6ad00001d@10.34.56.121/~/git/tfcontinuous.git/; \
    git push zone52 master'
if [ $? != 0 ]; then
    echo "Error, unable to clone repository"
    exit 1
fi


#ssh://5510399367c989f6ad00001d@10.34.56.121/~/git/tfcontinuous.git/
# Build Docke image
#echo "Building image ..."
#make -C $APP_WORKING_DIR/$REF_NAME
#if [ $? != 0 ]; then
#    echo "Error, unable to build Docker image"
#    exit 1
#fi

echo "Remove folder ${APP_WORKING_DIR}/${REF_NAME}..."
rm -rf ${APP_WORKING_DIR}/${REF_NAME}
echo "Build complete!"
exit 0

