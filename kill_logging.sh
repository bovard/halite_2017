cd src
find . -type f -name '*.go' -exec sed -i '' s/log.Println/\\/\\/log.Println/ {} +
find . -type f -name '*.go' -exec sed -i '' s/\"log\"/\\/\\/\"log\"/ {} +
cd -
