#!/bin/bash

if [ -d "${PWD}/configFiles" ]; then
    KUBECONFIG_FOLDER=${PWD}/configFiles
else
    echo "Configuration files are not found."
    exit
fi

# Copy the required chaincode into volume
echo -e "\nCreating Update chaincode job."
echo "Running: kubectl create -f ${KUBECONFIG_FOLDER}/updateChaincodeJob.yaml"
kubectl create -f ${KUBECONFIG_FOLDER}/updateChaincodeJob.yaml

pod=$(kubectl get pods --selector=job-name=updatechaincode --output=jsonpath={.items..metadata.name})

podSTATUS=$(kubectl get pods --selector=job-name=updatechaincode --output=jsonpath={.items..phase})

while [ "${podSTATUS}" != "Running" ]; do
    echo "Wating for container of updateChaincode pod to run. Current status of ${pod} is ${podSTATUS}"
    sleep 5;
    if [ "${podSTATUS}" == "Error" ]; then
        echo "There is an error in updatechaincode job. Please check logs."
        exit 1
    fi
    podSTATUS=$(kubectl get pods --selector=job-name=updatechaincode --output=jsonpath={.items..phase})
done

echo -e "pod \"${pod}\" status is ${podSTATUS}"
echo -e "\nStarting to copy chaincode in persistent volume."

#fix for this script to work on icp and ICS
kubectl cp $pod:/shared/crypto-config ./crypto-config

echo "Waiting for 10 more seconds for copying artifacts to avoid any network delay"
sleep 10
JOBSTATUS=$(kubectl get jobs --selector=job-name=updatechaincode --output=jsonpath={.items..status.succeeded})
while [ "${JOBSTATUS}" != "1" ]; do
    echo "Waiting for updatechaincode job to complete"
    sleep 1;
    PODSTATUS=$(kubectl get pods | grep "updatechaincode" | awk '{print $2}')
        if [ "${PODSTATUS}" == "Error" ]; then
            echo "There is an error in updateChaincodeJob. Please check logs."
            exit 1
        fi
    JOBSTATUS=$(kubectl get jobs --selector=job-name=updatechaincode --output=jsonpath={.items..status.succeeded})
done
echo "Update chaincode job completed"

kubectl delete -f ${KUBECONFIG_FOLDER}/updateChaincodeJob.yaml
