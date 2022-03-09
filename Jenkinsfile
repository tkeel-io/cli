pipeline {
  agent {
    kubernetes {
      //自定义执行环境
      label 'cli-test'
      yaml """
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: test
    image: mx2542/cli-test:1.1
    command: ['cat']
    tty: true
"""
    }
  }

    parameters {
        string(name:'NAME_SPACES',defaultValue: 'keel-system',description:'平台的命名空间')
    }

    environment {
        KUBECONFIG_CREDENTIAL_ID = 'kubeconfig'
    }

    stages {
        stage ('checkout scm') {
            steps {
                checkout(scm)
            }
        }
 
        stage ('build tkeel') {
            steps {
                container ('test') {
                sh 'make build'
                sh './dist/linux_amd64/release/tkeel -h'
                }
            }
        }

        stage('init tkeel') {
          steps {
            container ('test') {
                withCredentials([
                    kubeconfigFile(
                    credentialsId: env.KUBECONFIG_CREDENTIAL_ID,
                    variable: 'KUBECONFIG')
                    ]) {
                      sh './dist/linux_amd64/release/tkeel doctor'
                      sh './dist/linux_amd64/release/tkeel init --wait --timeout 3000'
                }
            }
          }
        }



        // stage ('build tkeel') {
        //     steps {
        //       sh 'ls'
        //       sh 'echo $env.KUBECONFIG_CREDENTIAL_ID'
        //       sh './cli -h'
        //     }
        // }

        // stage('push latest'){
        //    when{
        //      branch 'master'
        //    }
        //    steps{
        //         container ('maven') {
        //           sh 'docker tag  $REGISTRY/$DOCKERHUB_NAMESPACE/$APP_NAME:SNAPSHOT-$BRANCH_NAME-$BUILD_NUMBER $REGISTRY/$DOCKERHUB_NAMESPACE/$APP_NAME:latest '
        //           sh 'docker push  $REGISTRY/$DOCKERHUB_NAMESPACE/$APP_NAME:latest '
        //         }
        //    }
        // }

        // stage('deploy to dev') {
        //   when{
        //     branch 'master'
        //   }
        //   steps {
        //     input(id: 'deploy-to-dev', message: 'deploy to dev?')
        //     container ('maven') {
        //         withCredentials([
        //             kubeconfigFile(
        //             credentialsId: env.KUBECONFIG_CREDENTIAL_ID,
        //             variable: 'KUBECONFIG')
        //             ]) {
        //             sh 'envsubst < deploy/dev-all-in-one/devops-sample.yaml | kubectl apply -f -'
        //         }
        //     }
        //   }
        // }

        // stage('push with tag'){
        //   when{
        //     expression{
        //       return params.TAG_NAME =~ /v.*/
        //     }
        //   }
        //   steps {
        //       container ('maven') {
        //         input(id: 'release-image-with-tag', message: 'release image with tag?')
        //           withCredentials([usernamePassword(credentialsId: "$GITHUB_CREDENTIAL_ID", passwordVariable: 'GIT_PASSWORD', usernameVariable: 'GIT_USERNAME')]) {
        //             sh 'git config --global user.email "kubesphere@yunify.com" '
        //             sh 'git config --global user.name "kubesphere" '
        //             sh 'git tag -a $TAG_NAME -m "$TAG_NAME" '
        //             sh 'git push http://$GIT_USERNAME:$GIT_PASSWORD@github.com/$GITHUB_ACCOUNT/devops-java-sample.git --tags --ipv4'
        //           }
        //         sh 'docker tag  $REGISTRY/$DOCKERHUB_NAMESPACE/$APP_NAME:SNAPSHOT-$BRANCH_NAME-$BUILD_NUMBER $REGISTRY/$DOCKERHUB_NAMESPACE/$APP_NAME:$TAG_NAME '
        //         sh 'docker push  $REGISTRY/$DOCKERHUB_NAMESPACE/$APP_NAME:$TAG_NAME '
        //   }
        //   }
        // }
        // stage('deploy to production') {
        //   when{
        //     expression{
        //       return params.TAG_NAME =~ /v.*/
        //     }
        //   }
        //   steps {
        //     input(id: 'deploy-to-production', message: 'deploy to production?')
        //     container ('maven') {
        //         withCredentials([
        //             kubeconfigFile(
        //             credentialsId: env.KUBECONFIG_CREDENTIAL_ID,
        //             variable: 'KUBECONFIG')
        //             ]) {
        //             sh 'envsubst < deploy/prod-all-in-one/devops-sample.yaml | kubectl apply -f -'
        //         }
        //     }
        //   }
        // }
    }
}