#!/bin/bash

GetPackages() {
  PKGS_2_CHECK=""
  BASE_PACKAGE=""
 case $1 in 

 "test/brandstest")
     #echo "brand"
     BASE_PACKAGE="src/amenities/brands"
     ;;

 "test/categoriestest")
     BASE_PACKAGE="src/amenities/categories"
     ;;
     
 "test/productstest")
     BASE_PACKAGE="src/amenities/products"
     ;;

 "test/standardsizetest")
     BASE_PACKAGE="src/amenities/standardsize"
     ;;

 *)
    #echo "sdsd"
    return
    ;;
 esac	

for dir in $(find $BASE_PACKAGE -maxdepth 2 -mindepth 1 -type d);
do
 if [ -z "$(ls -A $dir | grep .go)" ]
then
  continue
 fi 

 if [ "$PKGS_2_CHECK" = "" ]
  then
   PKGS_2_CHECK="./$dir"
  else
   PKGS_2_CHECK="$PKGS_2_CHECK,./$dir"
 fi
   
done
echo $PKGS_2_CHECK
}


TEST_DIR="src/test"
TEMP_COVERAGE_FILE="tmp.cov"
COVERAGE_FILE="coverage.cov"
COVERAGE_HTML="coverage.html"

echo 'mode: count' > $COVERAGE_FILE #generate coverage file
for dir in $(find $TEST_DIR -maxdepth 1 -mindepth 1 -type d);
do
  sedDIR=`echo $dir|sed 's/src\///'`
    

  pck=$(GetPackages $sedDIR)
  if [ "$pck" = "" ] 
    then 
    continue
  fi  
  # echo $sedDIR 
  # echo $pck
  go test -c $sedDIR -covermode=count -coverpkg `echo $pck`
  name=`echo $dir|sed 's#.*/##'` #test name
  if [ -f "$name.test" ]
  then
    ./"$name.test" -test.coverprofile $TEMP_COVERAGE_FILE 
    sed -i '/_\/.*main.go/d' $TEMP_COVERAGE_FILE # remove illegal main file
    tail -q -n +2 $TEMP_COVERAGE_FILE >> $COVERAGE_FILE # append individual coverage
    rm "$name.test"
  else
    echo "no test file found in $dir"
  fi

done

rm $TEMP_COVERAGE_FILE # remove temp coverage
mv $COVERAGE_FILE bin/
cd bin
go tool cover -html=$COVERAGE_FILE -o $COVERAGE_HTML
