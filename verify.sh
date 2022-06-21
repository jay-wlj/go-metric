
path=`pwd`
echo "$path"

for subdir in "example/beego" "example/beegorm" "example/gin" "example/go-redis/v8/" "example/gorm/" "example/iris/" ; do
  echo $subdir
  cd $subdir
  go mod verify
  cd "$path"
done