go get

if [ "$BUILD" = false ] ; then
  go get github.com/cortesi/modd/cmd/modd
  modd
fi

if [ "$BUILD" = true ] ; then
  go build -o pdfformgenerator .
  ./pdfformgenerator
fi