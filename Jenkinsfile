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
        string(name:'GITHUB_ACCOUNT',defaultValue: 'lunz1207',description:'默认的 chart 仓库')
        string(name:'NAME_SPACES',defaultValue: 'testing',description:'平台的命名空间')
        string(name:'TKEEL_VERSION',defaultValue: '0.4.2-testing',description:'平台版本')
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
                      sh './dist/linux_amd64/release/tkeel init --repo-url https://$GITHUB_ACCOUNT.github.io/helm-charts/  --repo-name $GITHUB_ACCOUNT --runtime-version $TKEEL_VERSION --wait --timeout 3000'
                      // sh './dist/linux_amd64/release/tkeel init --wait --timeout 3000'
                }
            }
          }
        }
    }
}