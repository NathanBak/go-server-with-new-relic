#! /bin/sh

# Delete existing working directory and create new one
# Note that some files may not be deleted if not run as root
rm -rf working
mkdir working

# Copy swagger doc into place
cp ../api/swagger.yaml working/swagger.yaml

# Generate client in working directory
# Pinned to 3.0.46 until https://github.com/swagger-api/swagger-codegen/issues/12321 is resolved
docker run --rm -v $(pwd):/local swaggerapi/swagger-codegen-cli-v3:3.0.46 generate -i /local/working/swagger.yaml -l go -DpackageName=widgetclient -o /local/working/

# Use int instead of int32
find $(pwd)/working -type f -name "*.go" -exec sed -i.bak 's/int32/int/g' {} \;
rm -f $(pwd)/working/*.bak

# Format before removing header comments
gofmt -l -e -w $(pwd)/working

# Remove header comments
find $(pwd)/working -type f -name "*.go" -exec sed -i.bak '1,8d' {} \;
rm -f $(pwd)/working/*.bak

# Reformat again just in case
gofmt -l -e -w $(pwd)/working

# Clean package and stage in fresh files
rm $(pwd)/../generated-client/*
mkdir $(pwd)/../generated-client
cp working/*.go $(pwd)/../generated-client/