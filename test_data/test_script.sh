#! /usr/bin/env /bin/bash

test_md5_file="md5sums.txt"
target_md5_file="target_md5sums.txt"
exit_code=0

go run main.go format --test --write
find ./test_data/features -type f -exec md5sum '{}' \; > $test_md5_file
find ./test_data/target_files/features -type f -exec md5sum '{}' \; > $target_md5_file

md5_sum=""
md5_file=""

for word in $(cat $test_md5_file); do
  is_file=$(echo $word | grep -c "test_data")
  if [ $is_file -eq "1" ]; then
    md5_file=$word
    target_md5_sum_line=$(grep -E "/$(basename $md5_file)$" $target_md5_file)
    target_md5_sum_file=$(echo $target_md5_sum_line | awk '{ print $2 }')
    target_md5_sum=$(echo $target_md5_sum_line | awk '{ print $1 }')
    if [ $md5_sum != $target_md5_sum ]; then
      echo "------------------- NO MATCH -------------------"
      echo "Target MD5 sum file: $target_md5_sum_file" 
      echo "MD5 sum file: $md5_file"
      echo "Target MD5 sum: $target_md5_sum"
      echo "MD5 sum: $md5_sum"
      exit_code=1
    fi
  else
    md5_sum=$word
  fi
done

if [ -z $RESTORE_FILES ]; then
  git restore test_data
  rm $test_md5_file $target_md5_file
fi

exit $exit_code
