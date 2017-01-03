# Release instructions #

We are hosting the release binaries on Google Cloud Storage service.
Here is how to push the compiled binaries to the bucket:

    gsutil cp bazel-bin/espcli gs://endpoints-release/$VERSION/bin/$OSTYPE/amd64/
    gsutil acl set public-read gs://endpoints-release/$VERSION/bin/$OSTYPE/amd64/espcli
    
Make sure your version matches git tag as well as the name in `version.go` source file.
ESP CLI Version follows semantic version convention (for example, v1.0.3).
$OSTYPE should be either `darwin` or `linux`. Make sure to compile the static binary on
the destination OS before pushing the binaries.


