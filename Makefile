PROJECT_NAME := tsp-hardware-test
PROGRAM_NAME := hwt
DEPLOY_PATH := /mnt/SDCARD/Apps/HardwareTest

IP := 192.168.0.105
USN := root
PWD := tina

all: clean docker deploy

clean:
	rm bin/${PROGRAM_NAME} -f

docker:
	#docker run -d --name arkos-sdk -c 1024 -it --volume=/home/vitaly/GolandProjects/:/work/ --workdir=/work/ arkos-sdk
	docker exec arkos-sdk /bin/bash -c 'cd ${PROJECT_NAME} && make build'

build:
	go build -o bin/${PROGRAM_NAME} ${PROJECT_NAME}/src/

deploy:
	sshpass -p ${PWD} ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ${USN}@${IP} "rm ${DEPLOY_PATH}/${PROGRAM_NAME} -f"
	sshpass -p ${PWD} scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null bin/${PROGRAM_NAME} ${USN}@${IP}:${DEPLOY_PATH}
	sshpass -p ${PWD} ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ${USN}@${IP} "chmod 777 ${DEPLOY_PATH}/${PROGRAM_NAME}"
	sshpass -p ${PWD} ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ${USN}@${IP} "if pgrep ${PROGRAM_NAME}; then pkill -f ${PROGRAM_NAME}; fi"
	sshpass -p ${PWD} ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ${USN}@${IP} "sh -c 'cd /tmp; ${DEPLOY_PATH}/${PROGRAM_NAME}'" &
