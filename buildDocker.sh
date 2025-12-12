docker buildx create --name MultiBuilder --use
docker buildx build --platform=linux/amd64,linux/arm64 -f "Dockerfile" -t kevincfechtel/imaparchive:latest --push "." 
docker buildx rm MultiBuilder