#!/bin/bash

set -o errexit

E2E_DIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
PROJECT_ROOT=$E2E_DIR/../..

check_bin() {
  hash $1 2> /dev/null
}

check_bin_or_exit() {
  if ! check_bin $1
  then
    echo -e "error: command \"$1\" is needed but could not be found.\n"
    exit 1
  fi
}

install_k8s_e2e_test_pkg() {
  echo -e "installing k8s e2e test package of version: $k8s_server_version ...\n"
  curl -L https://storage.googleapis.com/kubernetes-release/release/${k8s_server_version}/kubernetes-test-linux-amd64.tar.gz --output e2e-tests.tar.gz
  tar -xf e2e-tests.tar.gz --directory=$E2E_DIR && rm e2e-tests.tar.gz
}

install_helm() {
  echo -e "installing helm...\n"
  curl -L https://get.helm.sh/helm-v3.7.1-linux-amd64.tar.gz --output helm.tar.gz
  tar -xf helm.tar.gz --directory=$E2E_DIR && rm helm.tar.gz
  mv $E2E_DIR/linux-amd64/helm /usr/local/bin
}

check_requires() {
  for pkg in git kubectl awk head curl tar docker sed
  do
    check_bin_or_exit $pkg
  done

  for path in /etc/qingcloud/access_key_id /etc/qingcloud/secret_access_key /etc/qingcloud/zone
  do
    if [[ ! -f $path ]]
    then
      echo -e "file $path is needed but could not be found.\n"
      exit 1
    fi
  done
}

prepare_packages() {
  if check_bin go
  then
    # in case it's not already set
    export PATH=$PATH:$(go env GOPATH)/bin
  fi

  # in case this script has already installed the test pkg
  # also, give this path more priority
  export PATH=$E2E_DIR/kubernetes/test/bin:$PATH

  if ! check_bin ginkgo
  then
    echo -e "command \"ginkgo\" not found, will install it...\n"
    need_install_ginkgo=true
  fi

  if ! check_bin e2e.test
  then
    echo -e "kubernetes e2e test package \"e2e.test\" not found in PATH, will install it.\n"
    need_install_e2e_test=true
  elif [[ $(e2e.test -version) == "$k8s_server_version"* ]]
  then
    need_install_e2e_test=false
  else
    echo -e "existing e2e test package of version: $(e2e.test -version) is incompatible with k8s server version: $k8s_server_version, will install it.\n"
    need_install_e2e_test=true
  fi

  if [[ "$need_install_e2e_test" == true ]] || [[ "$need_install_ginkgo" == true ]]
  then
    install_k8s_e2e_test_pkg
  fi

  if ! check_bin helm
  then
    echo -e "command \"helm\" not found, will install it..."
    install_helm
  fi
}

build_image() {
  image_url=${image_repo}:${commit_tag}

  if docker image inspect ${image_url} &> /dev/null
  then
    echo -e "docker image ${image_url} already exists, skip building.\n"
    return 0
  fi

  echo -e "building local image ${image_url} for local e2e test...\n\n"

  docker build -t ${image_url} -f "$PROJECT_ROOT/deploy/disk/docker/Dockerfile" $PROJECT_ROOT

  echo -e "\n"

  if [[ "${#nodes[@]}" -gt 1 ]]
  then
    echo -e "there are multiple nodes in this cluster\n"
    echo -e "try to push the built image to every node\n"
    echo -e "you might need to type in the ssh password if public key authentication isn't configured\n"
    for node in "${nodes[@]}"
    do
      echo -e "pushing image to node $node ...\n"
      docker save ${image_url} | ssh $node docker load
    done
  fi

  echo ""
}

install_helm_chart() {
  namespace="csi-qingcloud-e2e-test-${commit_tag}"
  name_template="csi-qingcloud-${commit_tag}"

  releases=($(helm list -n ${namespace} | awk 'NR>1 { print $1 }'))

  for rls in "${releases[@]}"
  do
      if [ "$rls" == "$name_template" ] ; then
          echo -e "Found existing helm release: $rls, uninstall it first\n"
          helm uninstall -n $namespace $rls
          echo ""
          break
      fi
  done

  echo -e "installing csi-qingcloud helm chart for local e2e test...\n"
  echo -e "will add commit tag to relevant values to avoid conflict\n"

  helm repo add ks-test https://charts.kubesphere.io/test

  helm install ks-test/csi-qingcloud \
  --namespace ${namespace} \
  --create-namespace \
  --name-template ${name_template} \
  --set config.qy_access_key_id=`cat /etc/qingcloud/access_key_id` \
  --set config.qy_secret_access_key=`cat /etc/qingcloud/secret_access_key` \
  --set config.zone=`cat /etc/qingcloud/zone` \
  --set driver.name=${commit_tag}.disk.csi.qingcloud.com \
  --set driver.repository=${image_repo} \
  --set driver.tag=${commit_tag} \
  --set sc.name=csi-qingcloud-${commit_tag}

  echo -e "\n"

}

install_snapshotclass() {
  echo -e "installing csi-qingcloud volumesnapshotclass for local e2e test...\n"
  sed -e "s/{commit_tag}/${commit_tag}/g" $E2E_DIR/snapshotclass_template.yaml | kubectl apply -f -
  echo ""
}

fmt_testdriver_file() {
  echo -e "generating testdriver.yaml for local e2e test...\n"
  sed -e "s/{commit_tag}/${commit_tag}/g" $E2E_DIR/testdriver_template.yaml > $E2E_DIR/testdriver.yaml
  if [[ "${#nodes[@]}" -gt 1 ]]
  then
    sed -i "s/{single_node}/true/g" $E2E_DIR/testdriver.yaml
  else
    echo -e "there is only one node in this cluster, disabling multi-node tests...\n"
    sed -i "s/{single_node}/false/g" $E2E_DIR/testdriver.yaml
  fi
}

run_test() {
  echo -e "begin running e2e test against k8s cluster...\n"

  # wait some time for the driver to finish initialization
  sleep 10s

  # we don't parallelize the runnings because that's more likely to trigger issues on the iaas layer
  # but what we need to test is the correct functioning of csi code
  # same reason for the skipped disruptive/stress test, which may be added in another test later

  # we explicitly skip the volume online expansion test
  # because in some versions it's not skipped even though the onlineExpansion cap is set to false
  logfile="${E2E_DIR}/e2e_test_${commit_tag}.log"

  echo "" > $logfile

  # the e2e tests take a very long time to run,
  # and the ssh connection breaks almost every time
  # so we use nohup here to prevent the test process from being terminated as well
  nohup ginkgo -focus='External.Storage.*' \
  -skip='(.*Disruptive.*|.*stress.*|.*should resize volume when PVC is edited while pod is using it.*)' \  $(which e2e.test) \
  -- -storage.testdriver="${E2E_DIR}/testdriver.yaml" \
  -kubeconfig="${HOME}/.kube/config" &> $logfile &

  tail -f $logfile
}

main() {
  check_requires

  k8s_server_version=$(kubectl version --short | awk '/Server/ {print $3}')
  nodes=($(kubectl get nodes | awk 'NR>1 { print $1 }'))

  image_repo="local-e2e-test/csi-qingcloud"

  commit_tag="git$(git rev-parse HEAD | head -c 6)"

  echo -e "will run e2e test at git commit ${commit_tag}:\n"
  echo -e "$(git log -1 --pretty=%B | cat)\n"

  prepare_packages

  build_image

  install_helm_chart

  install_snapshotclass

  fmt_testdriver_file

  run_test

}

main
