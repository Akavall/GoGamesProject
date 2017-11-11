echo "Running tests"
go test ./...
#echo "setting up temporary deployment folder"
#mkdir -p ../zombie_dice_deployment
#echo "building executable"
#go build ../zombie_dice_server/basic_server.go
#echo "moving executable"
#mv basic_server ../zombie_dice_deployment/
#echo "modifying html"
#sed "s/localhost/"$1"/g" ../zombie_dice_server/zombie_dice.html > ../zombie_dice_deployment/zombie_dice.html
#echo "deploying to" $2":"$3
#scp -r ../zombie_dice_deployment $2":"$3
#echo "done deploying to" $2":"$3
## https://stackoverflow.com/questions/29142/getting-ssh-to-execute-a-command-in-the-background-on-target-machine
#ssh -n -f $2 "sh -c 'cd ~/zombie_dice_deployment; nohup ./basic_server > /dev/null 2>&1 &'"
#echo "lunched the server"
#rm -rf ../zombie_dice_deployment
#echo "deleted temporary deployment folder"
#

