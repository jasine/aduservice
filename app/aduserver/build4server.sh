
green='\e[0;32m' # '\e[1;32m' is too bright for white bg.
red='\e[04;31m' 
endColor='\e[0m'

RunCommand() {
	$@

	if [ $? != 0 ]; then
		printf "${red}Failed when executing command: '$@' in `pwd`\n${endColor}"
		exit $ERROR_CODE
	else
		printf "${green}Successfully run '$@' in `pwd`\n${endColor}"
	fi
}

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 RunCommand  go build -o ./aduservice.linux  ./aduserver.go
echo "==>adu.linux"




